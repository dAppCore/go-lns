// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"net"
	"net/netip"
	"testing"

	core "dappco.re/go"
	miekgdns "github.com/miekg/dns"
)

// captureWriter is a miekgdns.ResponseWriter that records the reply message.
type captureWriter struct {
	msg *miekgdns.Msg
}

func (w *captureWriter) LocalAddr() net.Addr  { return &net.UDPAddr{} }
func (w *captureWriter) RemoteAddr() net.Addr { return &net.UDPAddr{} }
func (w *captureWriter) WriteMsg(m *miekgdns.Msg) error {
	w.msg = m
	return nil
}
func (w *captureWriter) Write(b []byte) (int, error) { return len(b), nil }
func (w *captureWriter) Close() error                { return nil }
func (w *captureWriter) TsigStatus() error           { return nil }
func (w *captureWriter) TsigTimersOnly(bool)         {}
func (w *captureWriter) Hijack()                     {}

// bridgeWithRecord builds a bridge whose cache holds a single populated name.
func bridgeWithRecord(t *testing.T, name, comment string) *Bridge {
	t.Helper()
	resource := resourceFromAliasComment(name, comment)
	if resource == nil {
		t.Fatalf("resourceFromAliasComment(%q,%q) returned nil", name, comment)
	}
	return &Bridge{
		cache: map[string]*NameRecord{
			name: {Name: name, Resource: resource, Source: "test", Comment: comment},
		},
	}
}

func TestBridgeBehaviour_dnsFQDN_Good(t *testing.T) {
	core.AssertEqual(t, "charon.", dnsFQDN("charon"))
	core.AssertEqual(t, "charon.", dnsFQDN("CHARON."))
	core.AssertEqual(t, "charon.lthn.", dnsFQDN("Charon.LTHN."))
}

func TestBridgeBehaviour_dnsFQDN_Bad(t *testing.T) {
	core.AssertEqual(t, ".", dnsFQDN(""))
	core.AssertEqual(t, ".", dnsFQDN("@"))
}

func TestBridgeBehaviour_dnsFQDN_Ugly(t *testing.T) {
	// Whitespace is trimmed before the suffix decision.
	core.AssertEqual(t, "charon.", dnsFQDN("  charon  "))
}

func TestBridgeBehaviour_aliasRecordHost_Good(t *testing.T) {
	core.AssertEqual(t, "charon.lthn.", aliasRecordHost("charon", "@"))
	core.AssertEqual(t, "charon.lthn.", aliasRecordHost("charon", ""))
	core.AssertEqual(t, "ns1.charon.lthn.", aliasRecordHost("charon", "ns1.charon.lthn."))
}

func TestBridgeBehaviour_aliasRecordHost_Bad(t *testing.T) {
	// A bare host without the .lthn. suffix is FQDN-normalised.
	core.AssertEqual(t, "ns1.example.com.", aliasRecordHost("charon", "ns1.example.com"))
}

func TestBridgeBehaviour_aliasRecordHost_Ugly(t *testing.T) {
	// Mixed-case host carrying the suffix is lowercased and kept as-is.
	core.AssertEqual(t, "ns1.charon.lthn.", aliasRecordHost("charon", "NS1.Charon.LTHN."))
}

func TestBridgeBehaviour_parseLTHNName_Good(t *testing.T) {
	label, ok := parseLTHNName("charon.lthn.")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", label)

	// The apex returns an empty label but reports ok.
	apex, okApex := parseLTHNName("lthn.")
	core.AssertTrue(t, okApex)
	core.AssertEqual(t, "", apex)
}

func TestBridgeBehaviour_parseLTHNName_Bad(t *testing.T) {
	label, ok := parseLTHNName("example.com.")
	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", label)
}

func TestBridgeBehaviour_parseLTHNName_Ugly(t *testing.T) {
	// Unqualified input is FQDN-promoted before the subdomain test.
	label, ok := parseLTHNName("charon.lthn")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", label)
}

func TestBridgeBehaviour_arpaToIP_Good(t *testing.T) {
	core.AssertEqual(t, "192.0.2.10", arpaToIP("10.2.0.192.in-addr.arpa."))
}

func TestBridgeBehaviour_arpaToIP_Bad(t *testing.T) {
	// Wrong suffix is rejected.
	core.AssertEqual(t, "", arpaToIP("10.2.0.192.ip6.arpa."))
	// Wrong octet count is rejected.
	core.AssertEqual(t, "", arpaToIP("2.0.192.in-addr.arpa."))
}

func TestBridgeBehaviour_arpaToIP_Ugly(t *testing.T) {
	// Empty input lacks the suffix entirely.
	core.AssertEqual(t, "", arpaToIP(""))
}

func TestBridgeBehaviour_soaRecord_Good(t *testing.T) {
	bridge := NewBridge(BridgeConfig{})
	soa, ok := bridge.soaRecord().(*miekgdns.SOA)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "lthn.", soa.Hdr.Name)
	core.AssertEqual(t, "ns1.lthn.", soa.Ns)
	core.AssertEqual(t, "hostmaster.lthn.", soa.Mbox)
}

func TestBridgeBehaviour_handleDNS_Good(t *testing.T) {
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=A:@:192.0.2.10|TXT:@:hello")

	request := new(miekgdns.Msg)
	request.SetQuestion("charon.lthn.", miekgdns.TypeA)
	writer := &captureWriter{}
	bridge.handleDNS(writer, request)

	core.AssertNotNil(t, writer.msg)
	core.AssertEqual(t, miekgdns.RcodeSuccess, writer.msg.Rcode)
	core.AssertLen(t, writer.msg.Answer, 1)
	answer, ok := writer.msg.Answer[0].(*miekgdns.A)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "192.0.2.10", answer.A.String())
}

func TestBridgeBehaviour_handleDNS_Bad(t *testing.T) {
	bridge := NewBridge(BridgeConfig{})

	// No question at all yields a format error.
	empty := new(miekgdns.Msg)
	writer := &captureWriter{}
	bridge.handleDNS(writer, empty)
	core.AssertNotNil(t, writer.msg)
	core.AssertEqual(t, miekgdns.RcodeFormatError, writer.msg.Rcode)

	// A non-.lthn name is refused.
	offZone := new(miekgdns.Msg)
	offZone.SetQuestion("example.com.", miekgdns.TypeA)
	writer2 := &captureWriter{}
	bridge.handleDNS(writer2, offZone)
	core.AssertEqual(t, miekgdns.RcodeRefused, writer2.msg.Rcode)

	// A known-zone but unknown name yields NXDOMAIN with an SOA authority.
	missing := new(miekgdns.Msg)
	missing.SetQuestion("nope.lthn.", miekgdns.TypeA)
	writer3 := &captureWriter{}
	bridge.handleDNS(writer3, missing)
	core.AssertEqual(t, miekgdns.RcodeNameError, writer3.msg.Rcode)
	core.AssertLen(t, writer3.msg.Ns, 1)
}

func TestBridgeBehaviour_handleDNS_Ugly(t *testing.T) {
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=A:@:192.0.2.10")

	// Apex SOA query returns the SOA in the answer section.
	apex := new(miekgdns.Msg)
	apex.SetQuestion("lthn.", miekgdns.TypeSOA)
	writer := &captureWriter{}
	bridge.handleDNS(writer, apex)
	core.AssertLen(t, writer.msg.Answer, 1)
	_, ok := writer.msg.Answer[0].(*miekgdns.SOA)
	core.AssertTrue(t, ok)

	// A query type with no matching records falls back to an SOA authority.
	noData := new(miekgdns.Msg)
	noData.SetQuestion("charon.lthn.", miekgdns.TypeAAAA)
	writer2 := &captureWriter{}
	bridge.handleDNS(writer2, noData)
	core.AssertLen(t, writer2.msg.Answer, 0)
	core.AssertLen(t, writer2.msg.Ns, 1)
}

func TestBridgeBehaviour_handlePTR_Good(t *testing.T) {
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=A:@:192.0.2.10")

	request := new(miekgdns.Msg)
	request.SetQuestion("10.2.0.192.in-addr.arpa.", miekgdns.TypePTR)
	writer := &captureWriter{}
	bridge.handlePTR(writer, func() *miekgdns.Msg { m := new(miekgdns.Msg); m.SetReply(request); return m }(), request.Question[0])

	core.AssertNotNil(t, writer.msg)
	core.AssertLen(t, writer.msg.Answer, 1)
	ptr, ok := writer.msg.Answer[0].(*miekgdns.PTR)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon.lthn.", ptr.Ptr)
}

func TestBridgeBehaviour_handlePTR_Bad(t *testing.T) {
	bridge := NewBridge(BridgeConfig{})

	request := new(miekgdns.Msg)
	request.SetQuestion("not.an.arpa.name.", miekgdns.TypePTR)
	reply := new(miekgdns.Msg)
	reply.SetReply(request)
	writer := &captureWriter{}
	bridge.handlePTR(writer, reply, request.Question[0])
	core.AssertEqual(t, miekgdns.RcodeNameError, writer.msg.Rcode)
}

func TestBridgeBehaviour_handlePTR_Ugly(t *testing.T) {
	// Valid arpa name but no cached record pointing at it: NXDOMAIN.
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=A:@:192.0.2.10")

	request := new(miekgdns.Msg)
	request.SetQuestion("99.99.99.99.in-addr.arpa.", miekgdns.TypePTR)
	reply := new(miekgdns.Msg)
	reply.SetReply(request)
	writer := &captureWriter{}
	bridge.handlePTR(writer, reply, request.Question[0])
	core.AssertLen(t, writer.msg.Answer, 0)
	core.AssertEqual(t, miekgdns.RcodeNameError, writer.msg.Rcode)
}

func TestBridgeBehaviour_handleDNS_TXTandNS(t *testing.T) {
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=TXT:@:hello|NS:@:ns1.example.com")

	txtReq := new(miekgdns.Msg)
	txtReq.SetQuestion("charon.lthn.", miekgdns.TypeTXT)
	w1 := &captureWriter{}
	bridge.handleDNS(w1, txtReq)
	core.AssertLen(t, w1.msg.Answer, 1)
	txt, ok := w1.msg.Answer[0].(*miekgdns.TXT)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "hello", txt.Txt[0])

	nsReq := new(miekgdns.Msg)
	nsReq.SetQuestion("charon.lthn.", miekgdns.TypeNS)
	w2 := &captureWriter{}
	bridge.handleDNS(w2, nsReq)
	core.AssertLen(t, w2.msg.Answer, 1)
	ns, ok := w2.msg.Answer[0].(*miekgdns.NS)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "ns1.example.com.", ns.Ns)
}

func TestBridgeBehaviour_Reverse_IPv6(t *testing.T) {
	bridge := bridgeWithRecord(t, "charon", "v=lthn1;dns=AAAA:@:2001:db8::1")
	names := bridge.Reverse(netip.MustParseAddr("2001:db8::1"))
	core.AssertLen(t, names, 1)
	core.AssertEqual(t, "charon.lthn.", names[0])
}
