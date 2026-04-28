// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"context"
	"net"
	"net/netip"
	"sort"
	"strings"
	"sync"
	"time"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	miekgdns "github.com/miekg/dns"
)

var defaultKnownNames = []string{
	"lethean", "snider", "darbs", "charon", "cladius",
	"explorer", "testnet", "builder", "develop", "miners",
	"relayer", "gateway", "monitor", "network", "storage",
	"support", "trade", "trading",
}

// NameRecord is a cached DNS view for a single .lthn name.
type NameRecord struct {
	Name     string
	Resource *Resource
	Source   string
	Comment  string
}

// BridgeConfig configures the live DNS bridge.
type BridgeConfig struct {
	AliasURL      string
	ChainURL      string
	APIKey        string
	DNSAddr       string
	CheckInterval time.Duration
	KnownNames    []string
}

// SyncReport reports the result of a discovery run.
type SyncReport struct {
	TreeRoot string   `json:"treeRoot"`
	Names    int      `json:"names"`
	Aliases  int      `json:"aliases"`
	Fallback bool     `json:"fallback"`
	Updated  []string `json:"updated"`
}

// BridgeHealth reports the current state of the bridge cache and server.
type BridgeHealth struct {
	Live          bool          `json:"live"`
	Serving       bool          `json:"serving"`
	Names         int           `json:"names"`
	TreeRoot      string        `json:"treeRoot"`
	LastSync      time.Time     `json:"lastSync"`
	DNSAddr       string        `json:"dnsAddr"`
	CheckInterval time.Duration `json:"checkInterval"`
	AliasURL      string        `json:"aliasUrl"`
	ChainURL      string        `json:"chainUrl"`
}

// Bridge syncs aliases and sidechain resources into a live DNS cache.
type Bridge struct {
	cfg         BridgeConfig
	aliasClient *chainpkg.AliasClient
	chainClient *chainpkg.Client

	mu          sync.RWMutex
	cache       map[string]*NameRecord
	lastRoot    string
	lastSync    time.Time
	serving     bool
	serveCancel context.CancelFunc
	udpServer   *miekgdns.Server
	tcpServer   *miekgdns.Server
}

// NewBridge constructs a DNS bridge from chain endpoints.
func NewBridge(cfg BridgeConfig) *Bridge {
	if cfg.DNSAddr == "" {
		cfg.DNSAddr = ":5354"
	}
	if cfg.CheckInterval <= 0 {
		cfg.CheckInterval = 15 * time.Second
	}
	if len(cfg.KnownNames) == 0 {
		cfg.KnownNames = append([]string(nil), defaultKnownNames...)
	}

	bridge := &Bridge{
		cfg:   cfg,
		cache: make(map[string]*NameRecord),
	}

	if cfg.AliasURL != "" {
		bridge.aliasClient = chainpkg.NewAliasClient(cfg.AliasURL)
	}
	if cfg.ChainURL != "" {
		bridge.chainClient = chainpkg.NewClient(cfg.ChainURL, chainpkg.WithAPIKey(cfg.APIKey))
	}

	return bridge
}

// Discover refreshes the bridge cache from alias discovery and chain resources.
func (b *Bridge) Discover(context.Context) (SyncReport, error) {
	if b == nil || b.chainClient == nil {
		return SyncReport{}, core.E("dns.Bridge.Discover", "chain client is required", nil)
	}

	blockchainInfo, err := b.chainClient.GetBlockchainInfo()
	if err != nil {
		return SyncReport{}, err
	}

	aliases := []chainpkg.Alias(nil)
	if b.aliasClient != nil {
		aliases, _ = b.aliasClient.GetAliases()
	}

	names := chainpkg.DiscoverNames(aliases, b.cfg.KnownNames)
	aliasByName := make(map[string]chainpkg.Alias, len(aliases))
	for _, alias := range aliases {
		aliasByName[chainpkg.HNSName(alias)] = alias
	}

	nextCache := make(map[string]*NameRecord, len(names))
	updated := make([]string, 0, len(names))

	for _, name := range names {
		resource, err := b.chainClient.GetNameResource(name)
		if err == nil {
			projected, err := resourceFromChainResource(resource)
			if err == nil && len(projected.Records) > 0 {
				record := &NameRecord{Name: name, Resource: projected, Source: "chain"}
				if alias, ok := aliasByName[name]; ok {
					record.Comment = alias.Comment
				}
				nextCache[name] = record
				updated = append(updated, name)
				continue
			}
		}

		if alias, ok := aliasByName[name]; ok {
			resource := resourceFromAliasComment(name, alias.Comment)
			if resource != nil {
				nextCache[name] = &NameRecord{Name: name, Resource: resource, Source: "alias", Comment: alias.Comment}
				updated = append(updated, name)
			}
		}
	}

	b.mu.Lock()
	b.cache = nextCache
	b.lastRoot = blockchainInfo.TreeRoot
	b.lastSync = time.Now().UTC()
	b.mu.Unlock()

	return SyncReport{
		TreeRoot: blockchainInfo.TreeRoot,
		Names:    len(nextCache),
		Aliases:  len(aliases),
		Fallback: len(aliases) == 0,
		Updated:  updated,
	}, nil
}

// Lookup returns a cached name record.
func (b *Bridge) Lookup(name string) (*NameRecord, bool) {
	if b == nil {
		return nil, false
	}

	name = strings.TrimSuffix(strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), ".lthn."), ".lthn")
	if name == "" {
		return nil, false
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	record, ok := b.cache[name]
	return record, ok
}

// Records returns a stable snapshot of the cached name records.
func (b *Bridge) Records() []NameRecord {
	if b == nil {
		return nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	names := make([]string, 0, len(b.cache))
	for name := range b.cache {
		names = append(names, name)
	}
	sort.Strings(names)

	records := make([]NameRecord, 0, len(names))
	for _, name := range names {
		record := b.cache[name]
		if record == nil {
			continue
		}
		records = append(records, *record)
	}

	return records
}

// Reverse returns all cached .lthn names that point at the address.
func (b *Bridge) Reverse(address netip.Addr) []string {
	if b == nil {
		return nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	var names []string
	for _, record := range b.cache {
		for _, candidate := range record.Addresses4() {
			if candidate == address.String() {
				names = append(names, record.FQDN())
				break
			}
		}
		for _, candidate := range record.Addresses6() {
			if candidate == address.String() {
				names = append(names, record.FQDN())
				break
			}
		}
	}
	return names
}

// Health returns the current bridge state.
func (b *Bridge) Health() BridgeHealth {
	if b == nil {
		return BridgeHealth{}
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return BridgeHealth{
		Live:          b.chainClient != nil,
		Serving:       b.serving,
		Names:         len(b.cache),
		TreeRoot:      b.lastRoot,
		LastSync:      b.lastSync,
		DNSAddr:       b.cfg.DNSAddr,
		CheckInterval: b.cfg.CheckInterval,
		AliasURL:      b.cfg.AliasURL,
		ChainURL:      b.cfg.ChainURL,
	}
}

// Serve starts authoritative UDP and TCP DNS listeners for the cached zone.
func (b *Bridge) Serve(ctx context.Context) error {
	if b == nil {
		return core.E("dns.Bridge.Serve", "bridge is required", nil)
	}

	b.mu.RLock()
	if b.serving {
		b.mu.RUnlock()
		return nil
	}
	b.mu.RUnlock()

	if b.Health().Names == 0 {
		if _, err := b.Discover(ctx); err != nil {
			return err
		}
	}

	handler := miekgdns.HandlerFunc(func(writer miekgdns.ResponseWriter, request *miekgdns.Msg) {
		b.handleDNS(writer, request)
	})

	packetConn, err := net.ListenPacket("udp", b.cfg.DNSAddr)
	if err != nil {
		return core.E("dns.Bridge.Serve", "bind udp", err)
	}
	packetConn.Close()

	udpServer := &miekgdns.Server{Addr: b.cfg.DNSAddr, Net: "udp", Handler: handler}
	tcpServer := &miekgdns.Server{Addr: b.cfg.DNSAddr, Net: "tcp", Handler: handler}
	serveCtx, serveCancel := context.WithCancel(context.Background())

	b.mu.Lock()
	if b.serving {
		b.mu.Unlock()
		serveCancel()
		return nil
	}
	b.udpServer = udpServer
	b.tcpServer = tcpServer
	b.serving = true
	b.serveCancel = serveCancel
	b.mu.Unlock()

	go func() { _ = udpServer.ListenAndServe() }()
	go func() { _ = tcpServer.ListenAndServe() }()
	if b.chainClient != nil && b.cfg.CheckInterval > 0 {
		go b.refreshLoop(serveCtx)
	}

	if ctx != nil {
		go func() {
			<-ctx.Done()
			if err := b.Shutdown(); err != nil {
				core.Warn("DNS bridge shutdown failed", "err", err)
			}
		}()
	}

	return nil
}

// Shutdown stops the live DNS listeners.
func (b *Bridge) Shutdown() error {
	if b == nil {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.serveCancel != nil {
		b.serveCancel()
		b.serveCancel = nil
	}

	var firstErr error
	if b.udpServer != nil {
		if err := b.udpServer.Shutdown(); err != nil && firstErr == nil && !strings.Contains(err.Error(), "not started") {
			firstErr = err
		}
	}
	if b.tcpServer != nil {
		if err := b.tcpServer.Shutdown(); err != nil && firstErr == nil && !strings.Contains(err.Error(), "not started") {
			firstErr = err
		}
	}
	b.udpServer = nil
	b.tcpServer = nil
	b.serving = false
	return firstErr
}

func (b *Bridge) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(b.cfg.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = b.refreshIfChanged(ctx)
		}
	}
}

func (b *Bridge) refreshIfChanged(ctx context.Context) (bool, error) {
	if b == nil || b.chainClient == nil {
		return false, nil
	}

	info, err := b.chainClient.GetBlockchainInfo()
	if err != nil {
		return false, err
	}

	b.mu.RLock()
	currentRoot := b.lastRoot
	b.mu.RUnlock()

	if info.TreeRoot != "" && info.TreeRoot == currentRoot {
		return false, nil
	}

	if _, err := b.Discover(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// FQDN returns the canonical `.lthn.` host name for the record.
func (record *NameRecord) FQDN() string {
	if record == nil || record.Name == "" {
		return "lthn."
	}
	return strings.ToLower(record.Name) + ".lthn."
}

// Addresses4 returns IPv4 glue addresses for the record.
func (record *NameRecord) Addresses4() []string {
	if record == nil || record.Resource == nil {
		return nil
	}

	var addresses []string
	for _, rr := range record.Resource.ToGlue(record.FQDN()) {
		if glue, ok := rr.(GLUE4Record); ok {
			addresses = append(addresses, glue.Address.String())
		}
	}
	return addresses
}

// Addresses6 returns IPv6 glue addresses for the record.
func (record *NameRecord) Addresses6() []string {
	if record == nil || record.Resource == nil {
		return nil
	}

	var addresses []string
	for _, rr := range record.Resource.ToGlue(record.FQDN()) {
		if glue, ok := rr.(GLUE6Record); ok {
			addresses = append(addresses, glue.Address.String())
		}
	}
	return addresses
}

// TXT returns TXT entries for the record.
func (record *NameRecord) TXT() []string {
	if record == nil || record.Resource == nil {
		return nil
	}

	var txt []string
	for _, rr := range record.Resource.ToTXT(record.FQDN()) {
		txt = append(txt, rr.Entries...)
	}
	return txt
}

// NS returns nameserver targets for the record.
func (record *NameRecord) NS() []string {
	if record == nil || record.Resource == nil {
		return nil
	}

	var nameservers []string
	for _, rr := range record.Resource.ToNS(record.FQDN()) {
		nameservers = append(nameservers, rr.NS)
	}
	return nameservers
}

func resourceFromAliasComment(name string, comment string) *Resource {
	parsed := chainpkg.ParseAliasComment(comment)
	dnsField := strings.TrimSpace(parsed["dns"])
	if dnsField == "" {
		return nil
	}

	resource := NewResource()
	var entries []string
	if strings.Contains(dnsField, "|") {
		entries = strings.Split(dnsField, "|")
	} else {
		entries = strings.Split(dnsField, ",")
	}

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, ":", 3)
		if len(parts) < 3 {
			continue
		}

		switch strings.ToUpper(parts[0]) {
		case "A":
			if address, err := netip.ParseAddr(strings.TrimSpace(parts[2])); err == nil {
				resource.Records = append(resource.Records, GLUE4Record{NS: aliasRecordHost(name, parts[1]), Address: address})
			}
		case "AAAA":
			if address, err := netip.ParseAddr(strings.TrimSpace(parts[2])); err == nil {
				resource.Records = append(resource.Records, GLUE6Record{NS: aliasRecordHost(name, parts[1]), Address: address})
			}
		case "TXT":
			resource.Records = append(resource.Records, TXTRecord{Entries: []string{strings.TrimSpace(parts[2])}})
		case "NS":
			resource.Records = append(resource.Records, NSRecord{NS: dnsFQDN(parts[2])})
		}
	}

	if len(resource.Records) == 0 {
		return nil
	}

	return resource
}

func dnsFQDN(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" || name == "@" {
		return "."
	}
	if strings.HasSuffix(name, ".") {
		return name
	}
	return name + "."
}

func aliasRecordHost(name, host string) string {
	host = strings.TrimSpace(host)
	if host == "" || host == "@" {
		return strings.ToLower(name) + ".lthn."
	}
	if strings.HasSuffix(strings.ToLower(host), ".lthn.") {
		return strings.ToLower(host)
	}
	return dnsFQDN(host)
}

func parseLTHNName(queryName string) (string, bool) {
	queryName = miekgdns.Fqdn(queryName)
	if !miekgdns.IsSubDomain("lthn.", queryName) {
		return "", false
	}
	if queryName == "lthn." {
		return "", true
	}
	label := strings.TrimSuffix(queryName, "lthn.")
	label = strings.TrimSuffix(label, ".")
	return label, label != ""
}

func (b *Bridge) handleDNS(writer miekgdns.ResponseWriter, request *miekgdns.Msg) {
	reply := new(miekgdns.Msg)
	reply.SetReply(request)
	reply.Authoritative = true
	reply.RecursionAvailable = false

	if len(request.Question) == 0 {
		reply.SetRcode(request, miekgdns.RcodeFormatError)
		if err := writer.WriteMsg(reply); err != nil {
			core.Warn("DNS response write failed", "err", err)
		}
		return
	}

	question := request.Question[0]
	if question.Qtype == miekgdns.TypePTR && miekgdns.IsSubDomain("in-addr.arpa.", miekgdns.Fqdn(question.Name)) {
		b.handlePTR(writer, reply, question)
		return
	}

	label, ok := parseLTHNName(question.Name)
	if !ok {
		reply.SetRcode(request, miekgdns.RcodeRefused)
		if err := writer.WriteMsg(reply); err != nil {
			core.Warn("DNS response write failed", "err", err)
		}
		return
	}

	if label == "" {
		if question.Qtype == miekgdns.TypeSOA || question.Qtype == miekgdns.TypeANY {
			reply.Answer = append(reply.Answer, b.soaRecord())
		} else {
			reply.Ns = []miekgdns.RR{b.soaRecord()}
		}
		if err := writer.WriteMsg(reply); err != nil {
			core.Warn("DNS response write failed", "err", err)
		}
		return
	}

	record, found := b.Lookup(label)
	if !found {
		reply.SetRcode(request, miekgdns.RcodeNameError)
		reply.Ns = []miekgdns.RR{b.soaRecord()}
		if err := writer.WriteMsg(reply); err != nil {
			core.Warn("DNS response write failed", "err", err)
		}
		return
	}

	switch question.Qtype {
	case miekgdns.TypeA:
		for _, address := range record.Addresses4() {
			if ip := net.ParseIP(address); ip != nil {
				reply.Answer = append(reply.Answer, &miekgdns.A{
					Hdr: miekgdns.RR_Header{Name: question.Name, Rrtype: miekgdns.TypeA, Class: miekgdns.ClassINET, Ttl: 300},
					A:   ip.To4(),
				})
			}
		}
	case miekgdns.TypeAAAA:
		for _, address := range record.Addresses6() {
			if ip := net.ParseIP(address); ip != nil {
				reply.Answer = append(reply.Answer, &miekgdns.AAAA{
					Hdr:  miekgdns.RR_Header{Name: question.Name, Rrtype: miekgdns.TypeAAAA, Class: miekgdns.ClassINET, Ttl: 300},
					AAAA: ip.To16(),
				})
			}
		}
	case miekgdns.TypeTXT:
		txt := record.TXT()
		if len(txt) > 0 {
			reply.Answer = append(reply.Answer, &miekgdns.TXT{
				Hdr: miekgdns.RR_Header{Name: question.Name, Rrtype: miekgdns.TypeTXT, Class: miekgdns.ClassINET, Ttl: 300},
				Txt: txt,
			})
		}
	case miekgdns.TypeNS:
		for _, ns := range record.NS() {
			reply.Answer = append(reply.Answer, &miekgdns.NS{
				Hdr: miekgdns.RR_Header{Name: question.Name, Rrtype: miekgdns.TypeNS, Class: miekgdns.ClassINET, Ttl: 300},
				Ns:  miekgdns.Fqdn(ns),
			})
		}
	case miekgdns.TypeSOA:
		reply.Answer = append(reply.Answer, b.soaRecord())
	default:
		reply.Ns = []miekgdns.RR{b.soaRecord()}
	}

	if len(reply.Answer) == 0 && question.Qtype != miekgdns.TypeSOA {
		reply.Ns = []miekgdns.RR{b.soaRecord()}
	}

	if err := writer.WriteMsg(reply); err != nil {
		core.Warn("DNS response write failed", "err", err)
	}
}

func (b *Bridge) handlePTR(writer miekgdns.ResponseWriter, reply *miekgdns.Msg, question miekgdns.Question) {
	ip := arpaToIP(miekgdns.Fqdn(question.Name))
	if ip == "" {
		reply.SetRcode(reply, miekgdns.RcodeNameError)
		if err := writer.WriteMsg(reply); err != nil {
			core.Warn("DNS response write failed", "err", err)
		}
		return
	}

	for _, name := range b.Reverse(netip.MustParseAddr(ip)) {
		reply.Answer = append(reply.Answer, &miekgdns.PTR{
			Hdr: miekgdns.RR_Header{Name: question.Name, Rrtype: miekgdns.TypePTR, Class: miekgdns.ClassINET, Ttl: 300},
			Ptr: miekgdns.Fqdn(name),
		})
	}

	if len(reply.Answer) == 0 {
		reply.SetRcode(reply, miekgdns.RcodeNameError)
	}

	if err := writer.WriteMsg(reply); err != nil {
		core.Warn("DNS response write failed", "err", err)
	}
}

func (b *Bridge) soaRecord() miekgdns.RR {
	return &miekgdns.SOA{
		Hdr:     miekgdns.RR_Header{Name: "lthn.", Rrtype: miekgdns.TypeSOA, Class: miekgdns.ClassINET, Ttl: 300},
		Ns:      "ns1.lthn.",
		Mbox:    "hostmaster.lthn.",
		Serial:  uint32(time.Now().Unix()),
		Refresh: 3600,
		Retry:   600,
		Expire:  86400,
		Minttl:  300,
	}
}

func arpaToIP(arpa string) string {
	const suffix = ".in-addr.arpa."
	if !strings.HasSuffix(arpa, suffix) {
		return ""
	}

	parts := strings.Split(strings.TrimSuffix(arpa, suffix), ".")
	if len(parts) != 4 {
		return ""
	}
	return parts[3] + "." + parts[2] + "." + parts[1] + "." + parts[0]
}
