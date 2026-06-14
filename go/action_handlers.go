// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/netip"
	"sort"
	"strconv"
	"time"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	dnspkg "dappco.re/go/lns/pkg/dns"
	"dappco.re/go/lns/pkg/primitives"
)

const (
	actionNameRegister       = "lns.name.register"
	actionNameRenew          = "lns.name.renew"
	actionNameTransfer       = "lns.name.transfer"
	actionNameUpdate         = "lns.name.update"
	actionNameRevoke         = "lns.name.revoke"
	actionNameDNSServe       = "lns.dns.serve"
	actionNameDNSDiscover    = "lns.dns.discover"
	actionNameWalletSend     = "lns.wallet.send"
	actionNameVPNDiscover    = "lns.vpn.discover"
	actionNameVPNConnect     = "lns.vpn.connect"
	actionNamePeerRegister   = "lns.peer.register"
	actionNameGossipReceive  = "lns.gossip.receive"
	actionNameGossipAnnounce = "lns.gossip.announce"
	actionOpBridge           = "lns.bridge"
	actionOpResource         = "lns.action.resource"
	actionOpHash             = "lns.action.hash"

	errInvalidName          = "invalid name"
	errNameRequired         = "name is required"
	errResourceRequired     = "resource is required"
	errWalletNotInitialized = "wallet is not initialized"
	vpnModeLocalSession     = "local-session"
)

var actionNames = []string{
	actionNameRegister,
	actionNameRenew,
	actionNameTransfer,
	actionNameUpdate,
	actionNameRevoke,
	"lns.name.lookup",
	"lns.name.auction",
	"lns.name.available",
	"lns.dns.resolve",
	"lns.dns.resolve.txt",
	"lns.dns.resolve.all",
	"lns.dns.reverse",
	actionNameDNSServe,
	"lns.dns.health",
	actionNameDNSDiscover,
	"lns.wallet.create",
	"lns.wallet.balance",
	actionNameWalletSend,
	"lns.wallet.receive",
	"lns.wallet.history",
	"lns.wallet.import",
	"lns.wallet.export",
	"lns.chain.height",
	"lns.chain.block",
	"lns.chain.sync",
	"lns.chain.status",
	"lns.chain.peers",
	actionNameVPNDiscover,
	actionNameVPNConnect,
	"lns.vpn.status",
	"lns.vpn.disconnect",
	actionNamePeerRegister,
	"lns.peer.lookup",
	"lns.peer.list",
	actionNameGossipReceive,
	actionNameGossipAnnounce,
}

type nameLookupResult struct {
	Input              string          `json:"input"`
	Canonical          string          `json:"canonical"`
	Hash               primitives.Hash `json:"hash"`
	Valid              bool            `json:"valid"`
	Reserved           bool            `json:"reserved"`
	Locked             bool            `json:"locked"`
	ReservedAtHeight   bool            `json:"reservedAtHeight"`
	LockedAtHeight     bool            `json:"lockedAtHeight"`
	Available          bool            `json:"available"`
	AvailabilityMode   string          `json:"availabilityMode"`
	AvailabilityReason string          `json:"availabilityReason"`
	Height             uint32          `json:"height"`
	RolloutStart       uint32          `json:"rolloutStart"`
	RolloutWeek        uint32          `json:"rolloutWeek"`
	HasRollout         bool            `json:"hasRollout"`
	Rules              *NameRules      `json:"rules,omitempty"`
	ReservedName       *ReservedName   `json:"reservedName,omitempty"`
	LockedName         *LockedName     `json:"lockedName,omitempty"`
}

type nameAvailabilityResult struct {
	Input              string `json:"input"`
	Canonical          string `json:"canonical"`
	Valid              bool   `json:"valid"`
	Available          bool   `json:"available"`
	AvailabilityMode   string `json:"availabilityMode"`
	AvailabilityReason string `json:"availabilityReason"`
	Reserved           bool   `json:"reserved"`
	Locked             bool   `json:"locked"`
	ReservedAtHeight   bool   `json:"reservedAtHeight"`
	LockedAtHeight     bool   `json:"lockedAtHeight"`
	Height             uint32 `json:"height"`
	HasRollout         bool   `json:"hasRollout"`
}

type nameAuctionResult struct {
	Input            string          `json:"input"`
	Canonical        string          `json:"canonical"`
	Hash             primitives.Hash `json:"hash"`
	Valid            bool            `json:"valid"`
	Height           uint32          `json:"height"`
	RolloutStart     uint32          `json:"rolloutStart"`
	RolloutWeek      uint32          `json:"rolloutWeek"`
	HasRollout       bool            `json:"hasRollout"`
	ReservedAtHeight bool            `json:"reservedAtHeight"`
	LockedAtHeight   bool            `json:"lockedAtHeight"`
	Available        bool            `json:"available"`
}

type dnsResolveResult struct {
	Zone      string        `json:"zone"`
	Host      string        `json:"host"`
	Records4  []GLUE4Record `json:"records4"`
	Records6  []GLUE6Record `json:"records6"`
	Addresses []string      `json:"addresses"`
}

type dnsTXTResult struct {
	Zone    string      `json:"zone"`
	Records []TXTRecord `json:"records"`
}

type dnsAllResult struct {
	Zone    string `json:"zone"`
	Records []any  `json:"records"`
}

type dnsReverseResult struct {
	Address string   `json:"address"`
	Names   []string `json:"names"`
}

type dnsHealthResult struct {
	OK           bool          `json:"ok"`
	Live         bool          `json:"live"`
	Service      string        `json:"service"`
	Actions      []string      `json:"actions"`
	Capabilities []string      `json:"capabilities"`
	Notes        []string      `json:"notes"`
	Bridge       *BridgeHealth `json:"bridge,omitempty"`
}

type nameActionPlanStep struct {
	Phase        string           `json:"phase"`
	CovenantType CovenantType     `json:"covenantType"`
	CovenantName string           `json:"covenantName"`
	Covenant     Covenant         `json:"covenant"`
	Draft        *nameDraftTx     `json:"draft,omitempty"`
	BidValue     uint64           `json:"bidValue,omitempty"`
	Nonce        *primitives.Hash `json:"nonce,omitempty"`
	Blind        *primitives.Hash `json:"blind,omitempty"`
	BlockHash    *primitives.Hash `json:"blockHash,omitempty"`
	Output       *Address         `json:"output,omitempty"`
	ResourceSize int              `json:"resourceSize,omitempty"`
	Notes        []string         `json:"notes,omitempty"`
}

type nameActionPlanResult struct {
	Action             string               `json:"action"`
	Input              string               `json:"input"`
	Canonical          string               `json:"canonical"`
	Hash               primitives.Hash      `json:"hash"`
	Height             uint32               `json:"height"`
	Available          bool                 `json:"available"`
	AvailabilityMode   string               `json:"availabilityMode"`
	AvailabilityReason string               `json:"availabilityReason"`
	Resource           *ResourceJSON        `json:"resource,omitempty"`
	ResourceSize       int                  `json:"resourceSize,omitempty"`
	Address            *Address             `json:"address,omitempty"`
	BlockHash          *primitives.Hash     `json:"blockHash,omitempty"`
	Steps              []nameActionPlanStep `json:"steps"`
	Warnings           []string             `json:"warnings,omitempty"`
}

type nameDraftTx struct {
	Transaction         Transaction `json:"transaction"`
	EstimatedSize       int         `json:"estimatedSize"`
	PlaceholderInputs   uint32      `json:"placeholderInputs"`
	PlaceholderOutputs  uint32      `json:"placeholderOutputs"`
	Fee                 uint64      `json:"fee"`
	FeeRate             uint64      `json:"feeRate"`
	PrimaryOutputValue  uint64      `json:"primaryOutputValue"`
	ChangeOutputValue   uint64      `json:"changeOutputValue,omitempty"`
	ChangeOutputPresent bool        `json:"changeOutputPresent"`
	Notes               []string    `json:"notes,omitempty"`
}

type vpnGatewayResult struct {
	Name         string            `json:"name"`
	FQDN         string            `json:"fqdn"`
	Source       string            `json:"source,omitempty"`
	Type         string            `json:"type,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Addresses    []string          `json:"addresses,omitempty"`
	TXT          []string          `json:"txt,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type vpnDiscoverResult struct {
	Gateways []vpnGatewayResult `json:"gateways"`
}

type vpnStatusResult struct {
	Connected bool              `json:"connected"`
	Gateway   *vpnGatewayResult `json:"gateway,omitempty"`
	Since     string            `json:"since,omitempty"`
	Mode      string            `json:"mode"`
	Notes     []string          `json:"notes,omitempty"`
}

type vpnConnectionState struct {
	Gateway     vpnGatewayResult
	ConnectedAt string
}

func (s *Service) registerActions() core.Result {
	if s == nil || s.Core() == nil {
		return okResult()
	}

	c := s.Core()
	c.Action(actionNameRegister, s.handleNameRegister).Description = "Prepare an LNS name registration request"
	c.Action(actionNameRenew, s.handleNameRenew).Description = "Prepare an LNS name renewal request"
	c.Action(actionNameTransfer, s.handleNameTransfer).Description = "Prepare an LNS name transfer request"
	c.Action(actionNameUpdate, s.handleNameUpdate).Description = "Prepare an LNS DNS resource update request"
	c.Action(actionNameRevoke, s.handleNameRevoke).Description = "Prepare an LNS name revocation request"
	c.Action("lns.name.lookup", s.handleNameLookup).Description = "Inspect canonical name status and catalog availability"
	c.Action("lns.name.auction", s.handleNameAuction).Description = "Report rollout and auction gating data for a name"
	c.Action("lns.name.available", s.handleNameAvailable).Description = "Check whether a name is available from the current local rule view"
	c.Action("lns.dns.resolve", s.handleDNSResolve).Description = "Project GLUE and SYNTH address records for a host"
	c.Action("lns.dns.resolve.txt", s.handleDNSResolveTXT).Description = "Project TXT records from an LNS resource payload"
	c.Action("lns.dns.resolve.all", s.handleDNSResolveAll).Description = "Project the full zone-order record set from an LNS resource payload"
	c.Action("lns.dns.reverse", s.handleDNSReverse).Description = "Project PTR-style reverse matches from local resource payloads"
	c.Action(actionNameDNSServe, s.handleDNSServe).Description = "Start the live DNS bridge service"
	c.Action("lns.dns.health", s.handleDNSHealth).Description = "Report local DNS action availability and implementation status"
	c.Action(actionNameDNSDiscover, s.handleDNSDiscover).Description = "Refresh the live DNS cache from the chain"
	c.Action("lns.wallet.create", s.handleWalletCreate).Description = "Create a deterministic local LNS wallet seed"
	c.Action("lns.wallet.balance", s.handleWalletBalance).Description = "Return the local LNS wallet balance view"
	c.Action(actionNameWalletSend, s.handleWalletSend).Description = "Construct and submit a sidechain wallet spend"
	c.Action("lns.wallet.receive", s.handleWalletReceive).Description = "Return the current or next receive address"
	c.Action("lns.wallet.history", s.handleWalletHistory).Description = "Return the local LNS wallet history"
	c.Action("lns.wallet.import", s.handleWalletImport).Description = "Import a deterministic LNS wallet seed"
	c.Action("lns.wallet.export", s.handleWalletExport).Description = "Export the current deterministic LNS wallet seed"
	c.Action("lns.chain.height", s.handleChainHeight).Description = "Return the sidechain height"
	c.Action("lns.chain.block", s.handleChainBlock).Description = "Fetch a raw sidechain block by hash or height"
	c.Action("lns.chain.sync", s.handleChainSync).Description = "Refresh the live LNS cache from alias discovery and sidechain resources"
	c.Action("lns.chain.status", s.handleChainStatus).Description = "Return sidechain and cache status"
	c.Action("lns.chain.peers", s.handleChainPeers).Description = "Return the raw sidechain peer list"
	c.Action(actionNameVPNDiscover, s.handleVPNDiscover).Description = "Discover VPN-capable gateways from alias metadata and TXT records"
	c.Action(actionNameVPNConnect, s.handleVPNConnect).Description = "Select a VPN gateway for the local LNS session"
	c.Action("lns.vpn.status", s.handleVPNStatus).Description = "Report the local VPN gateway session state"
	c.Action("lns.vpn.disconnect", s.handleVPNDisconnect).Description = "Clear the local VPN gateway session state"
	c.Action(actionNamePeerRegister, s.handlePeerRegister).Description = "Register a peer record and score it for the local registry"
	c.Action("lns.peer.lookup", s.handlePeerLookup).Description = "Look up a registered peer by bare name"
	c.Action("lns.peer.list", s.handlePeerList).Description = "List registered peers, optionally filtered by capability"
	c.Action(actionNameGossipReceive, s.handleGossipReceive).Description = "Verify and merge an inbound signed gossip delta into the local cache"
	c.Action(actionNameGossipAnnounce, s.handleGossipAnnounce).Description = "Broadcast a signed gossip delta to connected peers over the transport"
	return okResult()
}

func (s *Service) handleNameLookup(_ context.Context, options core.Options) core.Result {
	lookup, err := s.lookupNameAction(options)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(lookup)
}

func (s *Service) handleNameAvailable(_ context.Context, options core.Options) core.Result {
	lookup, err := s.lookupNameAction(options)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(
		nameAvailabilityResult{
			Input:              lookup.Input,
			Canonical:          lookup.Canonical,
			Valid:              lookup.Valid,
			Available:          lookup.Available,
			AvailabilityMode:   lookup.AvailabilityMode,
			AvailabilityReason: lookup.AvailabilityReason,
			Reserved:           lookup.Reserved,
			Locked:             lookup.Locked,
			ReservedAtHeight:   lookup.ReservedAtHeight,
			LockedAtHeight:     lookup.LockedAtHeight,
			Height:             lookup.Height,
			HasRollout:         lookup.HasRollout,
		})

}

func (s *Service) handleNameAuction(_ context.Context, options core.Options) core.Result {
	lookup, err := s.lookupNameAction(options)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(
		nameAuctionResult{
			Input:            lookup.Input,
			Canonical:        lookup.Canonical,
			Hash:             lookup.Hash,
			Valid:            lookup.Valid,
			Height:           lookup.Height,
			RolloutStart:     lookup.RolloutStart,
			RolloutWeek:      lookup.RolloutWeek,
			HasRollout:       lookup.HasRollout,
			ReservedAtHeight: lookup.ReservedAtHeight,
			LockedAtHeight:   lookup.LockedAtHeight,
			Available:        lookup.Available,
		})

}

func (s *Service) handleNameRegister(_ context.Context, options core.Options) core.Result {
	lookup, err := s.lookupNameAction(options)
	if err != nil {
		return core.Fail(err)
	}
	if !lookup.Valid {
		return core.Fail(core.E(actionNameRegister, errInvalidName, nil))
	}
	if !lookup.Available {
		return core.Fail(core.E(actionNameRegister, "name is not available", nil))
	}

	plan := newNameActionPlan(actionNameRegister, lookup)
	rawName := []byte(lookup.Canonical)
	ownerAddress, ownerWarning := s.nameDraftAddress(options, "address", "owner", "owner_address")
	if ownerWarning != "" {
		plan.Warnings = append(plan.Warnings, ownerWarning)
	}
	plan.Address = &ownerAddress

	resource, jsonView, size, err := actionOptionalResourceValue(options)
	if err != nil {
		return core.Fail(core.E(actionNameRegister, "invalid resource", err))
	}
	plan.Resource = jsonView
	plan.ResourceSize = size

	bidValue, hasBid := actionUint64Value(options, "bid_value", "bid", "value")
	if !hasBid {
		plan.Warnings = append(plan.Warnings, "bid_value was not provided; the draft blind commitment uses 0.")
	}

	nonce, ok, err := actionHashValue(options, "nonce")
	if err != nil {
		return core.Fail(core.E(actionNameRegister, "invalid nonce", err))
	}
	if !ok {
		nonce, err = randomActionHash()
		if err != nil {
			return core.Fail(err)
		}
		plan.Warnings = append(plan.Warnings, "nonce was generated automatically for this local draft.")
	}

	blind, err := s.Blind(bidValue, nonce)
	if err != nil {
		return core.Fail(err)
	}

	blockHash, blockHashSet, err := actionHashValue(options, "block_hash", "blockHash")
	if err != nil {
		return core.Fail(core.E(actionNameRegister, "invalid block hash", err))
	}
	if blockHashSet {
		plan.BlockHash = &blockHash
	} else {
		plan.Warnings = append(plan.Warnings, "block_hash was not provided; the REGISTER draft uses a zero hash placeholder.")
	}

	var open Covenant
	open.SetOpen(plan.Hash, rawName)
	var bid Covenant
	bid.SetBid(plan.Hash, plan.Height, rawName, blind)
	var reveal Covenant
	reveal.SetReveal(plan.Hash, plan.Height, nonce)
	var register Covenant
	register.SetRegister(plan.Hash, plan.Height, resource, blockHash)

	openValue, _ := actionUint64Value(options, "open_value")
	bidValueDraft, bidValueWarning := actionUint64Value(options, "bid_lockup", "lockup_value")
	if !bidValueWarning {
		bidValueDraft = bidValue
	}
	registerValue, registerValueSet := actionUint64Value(options, "register_value")
	if !registerValueSet {
		registerValue = bidValueDraft
	}
	revealValue, revealValueSet := actionUint64Value(options, "reveal_value")
	if !revealValueSet {
		revealValue = bidValueDraft
	}

	if nameWantsBroadcast(options) {
		phaseCovenants := map[string]struct {
			covenant Covenant
			value    uint64
		}{
			"open":     {open, openValue},
			"bid":      {bid, bidValueDraft},
			"reveal":   {reveal, revealValue},
			"register": {register, registerValue},
		}
		phase := core.Lower(stringOption(options, "phase", "step"))
		if phase == "" {
			phase = "open"
		}
		chosen, known := phaseCovenants[phase]
		if !known {
			return core.Fail(core.E(actionNameRegister, "phase must be one of open, bid, reveal, register", nil))
		}
		return s.broadcastNameCovenant(actionNameRegister, lookup, chosen.covenant, ownerAddress, chosen.value, options)
	}

	plan.Steps = []nameActionPlanStep{
		{
			Phase:        "open",
			CovenantType: TypeOpen,
			CovenantName: TypeName(TypeOpen),
			Covenant:     open,
			Draft:        buildNameDraftTx(options, open, ownerAddress, openValue, plan.Height),
			Output:       &ownerAddress,
			Notes:        []string{"Starts the auction for the name."},
		},
		{
			Phase:        "bid",
			CovenantType: TypeBid,
			CovenantName: TypeName(TypeBid),
			Covenant:     bid,
			Draft:        buildNameDraftTx(options, bid, ownerAddress, bidValueDraft, plan.Height),
			BidValue:     bidValue,
			Blind:        &blind,
			Output:       &ownerAddress,
			Notes:        []string{"Places a sealed bid commitment for the requested amount."},
		},
		{
			Phase:        "reveal",
			CovenantType: TypeReveal,
			CovenantName: TypeName(TypeReveal),
			Covenant:     reveal,
			Draft:        buildNameDraftTx(options, reveal, ownerAddress, revealValue, plan.Height),
			BidValue:     bidValue,
			Nonce:        &nonce,
			Output:       &ownerAddress,
			Notes:        []string{"Reveals the sealed bid nonce after the bidding window closes."},
		},
		{
			Phase:        "register",
			CovenantType: TypeRegister,
			CovenantName: TypeName(TypeRegister),
			Covenant:     register,
			Draft:        buildNameDraftTx(options, register, ownerAddress, registerValue, plan.Height),
			BlockHash:    &blockHash,
			Output:       &ownerAddress,
			ResourceSize: size,
			Notes:        []string{"Publishes the final DNS resource once the auction is won."},
		},
	}

	return core.Ok(plan)
}

func (s *Service) handleNameRenew(_ context.Context, options core.Options) core.Result {
	plan, err := s.prepareLinkedNameAction(actionNameRenew, options)
	if err != nil {
		return core.Fail(err)
	}

	blockHash, ok, err := actionHashValue(options, "block_hash", "blockHash")
	if err != nil {
		return core.Fail(core.E(actionNameRenew, "invalid block hash", err))
	}
	if !ok {
		plan.Warnings = append(plan.Warnings, "block_hash was not provided; the RENEW draft uses a zero hash placeholder.")
	}
	plan.BlockHash = &blockHash
	ownerAddress, ownerWarning := s.nameDraftAddress(options, "address", "owner", "owner_address")
	if ownerWarning != "" {
		plan.Warnings = append(plan.Warnings, ownerWarning)
	}
	plan.Address = &ownerAddress
	outputValue, _ := actionUint64Value(options, "renew_value", "value")

	var renew Covenant
	renew.SetRenew(plan.Hash, plan.Height, blockHash)

	if nameWantsBroadcast(options) {
		lookup, lookupErr := s.lookupNameAction(options)
		if lookupErr != nil {
			return core.Fail(core.E(actionNameRenew, "resolve name", lookupErr))
		}
		return s.broadcastNameCovenant(actionNameRenew, lookup, renew, ownerAddress, outputValue, options)
	}

	plan.Steps = []nameActionPlanStep{
		{
			Phase:        "renew",
			CovenantType: TypeRenew,
			CovenantName: TypeName(TypeRenew),
			Covenant:     renew,
			Draft:        buildNameDraftTx(options, renew, ownerAddress, outputValue, plan.Height),
			BlockHash:    &blockHash,
			Output:       &ownerAddress,
			Notes:        []string{"Extends the current registration period."},
		},
	}

	return core.Ok(plan)
}

func (s *Service) handleNameTransfer(_ context.Context, options core.Options) core.Result {
	plan, err := s.prepareLinkedNameAction(actionNameTransfer, options)
	if err != nil {
		return core.Fail(err)
	}

	address, ok := actionAddressValue(options, "address", "recipient", "to")
	if !ok {
		return core.Fail(core.E(actionNameTransfer, "recipient address is required", nil))
	}
	if !address.IsValid() {
		return core.Fail(core.E(actionNameTransfer, "recipient address is invalid", nil))
	}
	plan.Address = &address
	outputValue, _ := actionUint64Value(options, "transfer_value", "value")

	var transfer Covenant
	transfer.SetTransfer(plan.Hash, plan.Height, address)

	if nameWantsBroadcast(options) {
		lookup, lookupErr := s.lookupNameAction(options)
		if lookupErr != nil {
			return core.Fail(core.E(actionNameTransfer, "resolve name", lookupErr))
		}
		ownerAddress, _ := s.nameDraftAddress(options, "owner", "owner_address", "from")
		return s.broadcastNameCovenant(actionNameTransfer, lookup, transfer, ownerAddress, outputValue, options)
	}

	plan.Steps = []nameActionPlanStep{
		{
			Phase:        "transfer",
			CovenantType: TypeTransfer,
			CovenantName: TypeName(TypeTransfer),
			Covenant:     transfer,
			Draft:        buildNameDraftTx(options, transfer, address, outputValue, plan.Height),
			Output:       &address,
			Notes:        []string{"Begins transfer lockup to the destination address."},
		},
	}

	return core.Ok(plan)
}

func (s *Service) handleNameUpdate(_ context.Context, options core.Options) core.Result {
	plan, err := s.prepareLinkedNameAction(actionNameUpdate, options)
	if err != nil {
		return core.Fail(err)
	}

	resource, jsonView, size, err := actionOptionalResourceValue(options)
	if err != nil {
		return core.Fail(core.E(actionNameUpdate, "invalid resource", err))
	}
	plan.Resource = jsonView
	plan.ResourceSize = size
	ownerAddress, ownerWarning := s.nameDraftAddress(options, "address", "owner", "owner_address")
	if ownerWarning != "" {
		plan.Warnings = append(plan.Warnings, ownerWarning)
	}
	plan.Address = &ownerAddress
	outputValue, _ := actionUint64Value(options, "update_value", "value")

	var update Covenant
	update.SetUpdate(plan.Hash, plan.Height, resource)

	if nameWantsBroadcast(options) {
		lookup, lookupErr := s.lookupNameAction(options)
		if lookupErr != nil {
			return core.Fail(core.E(actionNameUpdate, "resolve name", lookupErr))
		}
		return s.broadcastNameCovenant(actionNameUpdate, lookup, update, ownerAddress, outputValue, options)
	}

	plan.Steps = []nameActionPlanStep{
		{
			Phase:        "update",
			CovenantType: TypeUpdate,
			CovenantName: TypeName(TypeUpdate),
			Covenant:     update,
			Draft:        buildNameDraftTx(options, update, ownerAddress, outputValue, plan.Height),
			Output:       &ownerAddress,
			ResourceSize: size,
			Notes:        []string{"Updates the DNS resource payload for the name."},
		},
	}

	return core.Ok(plan)
}

func (s *Service) handleNameRevoke(_ context.Context, options core.Options) core.Result {
	plan, err := s.prepareLinkedNameAction(actionNameRevoke, options)
	if err != nil {
		return core.Fail(err)
	}
	ownerAddress, ownerWarning := s.nameDraftAddress(options, "address", "owner", "owner_address")
	if ownerWarning != "" {
		plan.Warnings = append(plan.Warnings, ownerWarning)
	}
	plan.Address = &ownerAddress
	outputValue, _ := actionUint64Value(options, "revoke_value", "value")

	var revoke Covenant
	revoke.SetRevoke(plan.Hash, plan.Height)

	if nameWantsBroadcast(options) {
		lookup, lookupErr := s.lookupNameAction(options)
		if lookupErr != nil {
			return core.Fail(core.E(actionNameRevoke, "resolve name", lookupErr))
		}
		return s.broadcastNameCovenant(actionNameRevoke, lookup, revoke, ownerAddress, outputValue, options)
	}

	plan.Steps = []nameActionPlanStep{
		{
			Phase:        "revoke",
			CovenantType: TypeRevoke,
			CovenantName: TypeName(TypeRevoke),
			Covenant:     revoke,
			Draft:        buildNameDraftTx(options, revoke, ownerAddress, outputValue, plan.Height),
			Output:       &ownerAddress,
			Notes:        []string{"Burns the name and makes it unavailable until re-auctioned."},
		},
	}

	return core.Ok(plan)
}

func (s *Service) handleDNSResolve(_ context.Context, options core.Options) core.Result {
	zone, host, resource, err := s.dnsActionInputs(context.Background(), options)
	if err != nil {
		return core.Fail(err)
	}

	glue := resource.ToGlue(zone)
	result := dnsResolveResult{
		Zone: zone,
		Host: host,
	}

	for _, record := range glue {
		switch rr := record.(type) {
		case GLUE4Record:
			if !dnsNameMatches(rr.NS, host) {
				continue
			}
			result.Records4 = append(result.Records4, rr)
			result.Addresses = append(result.Addresses, rr.Address.String())
		case GLUE6Record:
			if !dnsNameMatches(rr.NS, host) {
				continue
			}
			result.Records6 = append(result.Records6, rr)
			result.Addresses = append(result.Addresses, rr.Address.String())
		}
	}

	return core.Ok(result)
}

func (s *Service) handleDNSResolveTXT(_ context.Context, options core.Options) core.Result {
	zone, _, resource, err := s.dnsActionInputs(context.Background(), options)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(
		dnsTXTResult{
			Zone:    zone,
			Records: resource.ToTXT(zone),
		})

}

func (s *Service) handleDNSResolveAll(_ context.Context, options core.Options) core.Result {
	zone, _, resource, err := s.dnsActionInputs(context.Background(), options)
	if err != nil {
		return core.Fail(err)
	}

	records := resource.ToZone(zone)
	projected := make([]any, 0, len(records))
	for _, record := range records {
		projected = append(projected, record)
	}

	return core.Ok(
		dnsAllResult{
			Zone:    zone,
			Records: projected,
		})

}

func (s *Service) handleDNSReverse(ctx context.Context, options core.Options) core.Result {
	address, err := dnsReverseAddress(options)
	if err != nil {
		return core.Fail(err)
	}

	if actionOptionValue(options, "resources", "resource", "resource_json") == nil {
		if names, ok := s.reverseNamesFromBridge(ctx, options, address); ok {
			return core.Ok(dnsReverseResult{
				Address: address.String(),
				Names:   names,
			})
		}
	}

	resources, err := dnsReverseResources(options)
	if err != nil {
		return core.Fail(err)
	}

	seen := map[string]struct{}{}
	names := make([]string, 0)

	for zone, resource := range resources {
		for _, record := range resource.ToGlue(zone) {
			names = appendReverseMatch(names, seen, record, address)
		}
	}

	sort.Strings(names)

	return core.Ok(
		dnsReverseResult{
			Address: address.String(),
			Names:   names,
		})

}

func (s *Service) reverseNamesFromBridge(ctx context.Context, options core.Options, address netip.Addr) ([]string, bool) {
	bridge, err := s.bridgeFromOptions(options)
	if err != nil || bridge == nil {
		return nil, false
	}
	if bridge.Health().Names == 0 {
		_, _ = bridge.Discover(ctx)
	}

	names := bridge.Reverse(address)
	sort.Strings(names)
	return names, true
}

func appendReverseMatch(names []string, seen map[string]struct{}, record ResourceRecord, address netip.Addr) []string {
	ns, ok := glueRecordNameForAddress(record, address)
	if !ok {
		return names
	}
	if _, exists := seen[ns]; exists {
		return names
	}
	seen[ns] = struct{}{}
	return append(names, ns)
}

func glueRecordNameForAddress(record ResourceRecord, address netip.Addr) (string, bool) {
	switch rr := record.(type) {
	case GLUE4Record:
		return rr.NS, rr.Address == address
	case GLUE6Record:
		return rr.NS, rr.Address == address
	default:
		return "", false
	}
}

func (s *Service) handleDNSHealth(_ context.Context, options core.Options) core.Result {
	capabilities := []string{
		"local-resource-projection",
		"catalog-backed-name-inspection",
	}

	result := dnsHealthResult{
		OK:           true,
		Live:         false,
		Service:      ServiceName,
		Actions:      append([]string(nil), actionNames...),
		Capabilities: capabilities,
		Notes: []string{
			"Chain-backed HSD discovery is not implemented yet.",
			"DNS actions currently project caller-supplied resource payloads.",
		},
	}

	if bridge, err := s.bridgeFromOptions(options); err == nil && bridge != nil {
		health := bridge.Health()
		result.Live = health.Live
		result.Bridge = &health
		result.Notes = []string{
			"Live DNS bridge is available when chain_url is configured.",
			"Bridge refresh uses tree-root invalidation and optional alias fallback.",
		}
		if health.Live {
			result.Capabilities = append(result.Capabilities, "live-chain-discovery")
		}
		if health.Serving {
			result.Capabilities = append(result.Capabilities, "authoritative-serve")
		}
	}

	return core.Ok(result)
}

func (s *Service) handleDNSServe(_ context.Context, options core.Options) core.Result {
	bridge, err := s.bridgeFromOptions(options)
	if err != nil {
		return core.Fail(core.E(actionNameDNSServe, "chain_url is required", err))
	}
	if err := bridge.Serve(context.Background()); err != nil {
		return core.Fail(core.E(actionNameDNSServe, "start live DNS bridge", err))
	}
	return core.Ok(bridge.Health())
}

func (s *Service) handleDNSDiscover(ctx context.Context, options core.Options) core.Result {
	bridge, err := s.bridgeFromOptions(options)
	if err != nil {
		return core.Fail(core.E(actionNameDNSDiscover, "chain_url is required", err))
	}

	report, err := bridge.Discover(ctx)
	if err != nil {
		return core.Fail(core.E(actionNameDNSDiscover, "refresh live DNS cache", err))
	}

	return core.Ok(report)
}

func (s *Service) handleVPNDiscover(ctx context.Context, options core.Options) core.Result {
	gateways, err := s.vpnGatewaysFromOptions(ctx, options)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(vpnDiscoverResult{Gateways: gateways})
}

func (s *Service) handleVPNConnect(ctx context.Context, options core.Options) core.Result {
	gateways, err := s.vpnGatewaysFromOptions(ctx, options)
	if err != nil {
		return core.Fail(err)
	}

	filter := actionStringFromAny(actionOptionValue(options, "name", "gateway", "_input"))
	gateway, ok := selectVPNGateway(gateways, filter)
	if !ok {
		return core.Fail(core.E(actionNameVPNConnect, "gateway was not found", nil))
	}

	status := s.setVPNConnection(gateway)
	status.Notes = append(status.Notes, "Connection state is local to this service instance; no OS tunnel is created yet.")
	return core.Ok(status)
}

func (s *Service) handleVPNStatus(_ context.Context, _ core.Options) core.Result {
	return core.Ok(s.currentVPNStatus())
}

func (s *Service) handleVPNDisconnect(_ context.Context, _ core.Options) core.Result {
	return core.Ok(s.clearVPNConnection())
}

func (s *Service) lookupNameAction(options core.Options) (nameLookupResult, error) {
	input := actionOptionValue(options, "name", "label", "_input")
	if input == nil {
		return nameLookupResult{}, core.E("lns.action.name", errNameRequired, nil)
	}

	result := nameLookupResult{
		Input: actionStringFromAny(input),
		Valid: s.Verify(input),
	}

	if canonical, ok := canonicalizeName(input); ok {
		result.Canonical = canonical
	}

	if !result.Valid {
		result.AvailabilityMode = "catalog"
		result.AvailabilityReason = "invalid"
		return result, nil
	}

	hash, err := s.Resolve(input)
	if err != nil {
		return nameLookupResult{}, err
	}

	result.Hash = hash

	if reserved, ok := s.GetReserved(input); ok {
		result.Reserved = true
		result.ReservedName = &reserved
	}

	if locked, ok := s.GetLocked(input); ok {
		result.Locked = true
		result.LockedName = &locked
	}

	height, hasHeight := actionUint32Value(options, "height")
	rules, hasRules := actionNameRulesValue(options)

	if hasHeight {
		result.Height = height
	}

	if hasRules {
		result.Rules = &rules
		result.RolloutStart, result.RolloutWeek = s.GetRollout(hash, rules)
		result.HasRollout = !hasHeight || s.HasRollout(hash, height, rules)
		if hasHeight {
			result.ReservedAtHeight = s.IsReserved(hash, height, rules)
			result.LockedAtHeight = s.IsLockedUp(hash, height, rules)
		}
		result.AvailabilityMode = "height"
		result.Available = result.HasRollout && !result.ReservedAtHeight && !result.LockedAtHeight
		result.AvailabilityReason = nameAvailabilityReason(result, true)
		return result, nil
	}

	result.AvailabilityMode = "catalog"
	result.Available = !result.Reserved && !result.Locked
	result.AvailabilityReason = nameAvailabilityReason(result, false)
	return result, nil
}

func newNameActionPlan(action string, lookup nameLookupResult) nameActionPlanResult {
	plan := nameActionPlanResult{
		Action:             action,
		Input:              lookup.Input,
		Canonical:          lookup.Canonical,
		Hash:               lookup.Hash,
		Height:             lookup.Height,
		Available:          lookup.Available,
		AvailabilityMode:   lookup.AvailabilityMode,
		AvailabilityReason: lookup.AvailabilityReason,
		Warnings: []string{
			"This is the planning view: it does not build inputs, fees, signatures, or broadcast transactions.",
			"Pass broadcast=true (with chain_url and, for owned-name covenants, owner_outpoint) to build, sign, and broadcast a real transaction.",
		},
	}

	if lookup.Height == 0 {
		plan.Warnings = append(plan.Warnings, "height was not provided; covenant drafts default to height 0.")
	}

	return plan
}

func (s *Service) prepareLinkedNameAction(operation string, options core.Options) (nameActionPlanResult, error) {
	lookup, err := s.lookupNameAction(options)
	if err != nil {
		return nameActionPlanResult{}, err
	}
	if !lookup.Valid {
		return nameActionPlanResult{}, core.E(operation, errInvalidName, nil)
	}

	plan := newNameActionPlan(operation, lookup)
	plan.Warnings = append(plan.Warnings, "Ownership and current chain state are not validated by this local draft.")
	return plan, nil
}

func (s *Service) nameDraftAddress(options core.Options, keys ...string) (Address, string) {
	if address, ok := actionAddressValue(options, keys...); ok && address.IsValid() {
		return address, ""
	}

	if wallet, ok := s.currentWallet(); ok {
		return wallet.ReceiveAddress(), ""
	}

	return Address{
		Version: 0,
		Hash:    make([]byte, 20),
	}, "no owner address or wallet was provided; the draft uses a zero-hash P2WPKH placeholder address."
}

func buildNameDraftTx(options core.Options, covenant Covenant, address Address, primaryValue uint64, lockTime uint32) *nameDraftTx {
	outputs := []Output{{
		Value:    primaryValue,
		Address:  address,
		Covenant: covenant,
	}}

	changeValue, hasChangeValue := actionUint64Value(options, "change_value")
	changeAddress, hasChangeAddress := actionAddressValue(options, "change_address")
	changePresent := hasChangeValue && hasChangeAddress && changeAddress.IsValid()
	if changePresent {
		outputs = append(outputs, Output{
			Value:   changeValue,
			Address: changeAddress,
		})
	}

	tx := Transaction{
		Version:  0,
		Outputs:  outputs,
		LockTime: lockTime,
	}

	placeholderInputs, ok := actionUint32Value(options, "input_count", "inputs")
	if !ok || placeholderInputs == 0 {
		placeholderInputs = 1
	}

	feeRate, ok := actionUint64Value(options, "fee_rate")
	if !ok || feeRate == 0 {
		feeRate = 1
	}

	estimatedSize := estimateNameDraftTxSize(tx, placeholderInputs)
	fee, ok := actionUint64Value(options, "fee")
	if !ok {
		fee = uint64(estimatedSize) * feeRate
	}

	notes := []string{
		"Inputs and witnesses are placeholders only; signable funding data is not assembled here.",
		"Fee estimation is approximate and intended for planning only.",
	}
	if !changePresent && hasChangeValue {
		notes = append(notes, "change_value was ignored because no valid change_address was provided.")
	}

	return &nameDraftTx{
		Transaction:         tx,
		EstimatedSize:       estimatedSize,
		PlaceholderInputs:   placeholderInputs,
		PlaceholderOutputs:  uint32(len(tx.Outputs)),
		Fee:                 fee,
		FeeRate:             feeRate,
		PrimaryOutputValue:  primaryValue,
		ChangeOutputValue:   changeValue,
		ChangeOutputPresent: changePresent,
		Notes:               notes,
	}
}

func estimateNameDraftTxSize(tx Transaction, placeholderInputs uint32) int {
	size := 4 + 1 + 1 + 4
	size += int(placeholderInputs) * 180

	for _, output := range tx.Outputs {
		size += 8
		size += output.Address.GetSize()
		size += output.Covenant.GetVarSize()
	}

	return size
}

func actionOptionalResourceValue(options core.Options) ([]byte, *ResourceJSON, int, error) {
	value := actionOptionValue(options, "resource", "resource_json")
	if value == nil {
		resource := NewResource()
		encoded, err := resource.Encode()
		if err != nil {
			return nil, nil, 0, err
		}
		jsonView := resource.GetJSON()
		return encoded, &jsonView, len(encoded), nil
	}

	resource, err := resourceFromAny(value)
	if err != nil {
		return nil, nil, 0, err
	}

	encoded, err := resource.Encode()
	if err != nil {
		return nil, nil, 0, err
	}
	jsonView := resource.GetJSON()
	return encoded, &jsonView, len(encoded), nil
}

func nameAvailabilityReason(result nameLookupResult, usingRules bool) string {
	if !result.Valid {
		return "invalid"
	}
	if usingRules {
		if !result.HasRollout {
			return "rollout"
		}
		if result.ReservedAtHeight {
			return "reserved"
		}
		if result.LockedAtHeight {
			return "locked"
		}
		return "available"
	}
	if result.Reserved {
		return "reserved"
	}
	if result.Locked {
		return "locked"
	}
	return "available"
}

func dnsLookupInputs(options core.Options) (string, string, *Resource, error) {
	zone, host, err := dnsLookupNames(options)
	if err != nil {
		return "", "", nil, err
	}

	resource, err := actionResourceValue(options, "resource", "resource_json")
	if err != nil {
		return "", "", nil, err
	}

	return zone, host, resource, nil
}

func dnsLookupNames(options core.Options) (string, string, error) {
	zoneValue := actionOptionValue(options, "name", "zone", "_input")
	if zoneValue == nil {
		return "", "", core.E("lns.action.dns", errNameRequired, nil)
	}

	zone := dnsFQDN(actionStringFromAny(zoneValue))
	if zone == "." {
		return "", "", core.E("lns.action.dns", errNameRequired, nil)
	}

	host := zone
	if hostValue := actionOptionValue(options, "host"); hostValue != nil {
		host = dnsFQDN(actionStringFromAny(hostValue))
	}

	return zone, host, nil
}

func (s *Service) dnsActionInputs(ctx context.Context, options core.Options) (string, string, *Resource, error) {
	zone, host, err := dnsLookupNames(options)
	if err != nil {
		return "", "", nil, err
	}

	if value := actionOptionValue(options, "resource", "resource_json"); value != nil {
		resource, err := resourceFromAny(value)
		if err != nil {
			return "", "", nil, err
		}
		return zone, host, resource, nil
	}

	bridge, err := s.bridgeFromOptions(options)
	if err != nil {
		return "", "", nil, core.E(actionOpResource, errResourceRequired, nil)
	}

	if bridge.Health().Names == 0 {
		_, _ = bridge.Discover(ctx)
	}

	label := core.TrimSuffix(core.TrimSuffix(core.Lower(zone), ".lthn."), ".lthn")
	if record, ok := bridge.Lookup(label); ok && record != nil && record.Resource != nil {
		return zone, host, record.Resource, nil
	}

	return zone, host, NewResource(), nil
}

func dnsReverseAddress(options core.Options) (netip.Addr, error) {
	value := actionOptionValue(options, "ip", "address", "_input")
	if value == nil {
		return netip.Addr{}, core.E("lns.action.dns.reverse", "ip is required", nil)
	}

	text := actionStringFromAny(value)
	addr, err := netip.ParseAddr(text)
	if err != nil {
		return netip.Addr{}, core.E("lns.action.dns.reverse", "invalid ip", err)
	}

	return addr, nil
}

func dnsReverseResources(options core.Options) (map[string]*Resource, error) {
	if value := actionOptionValue(options, "resources"); value != nil {
		return actionResourcesValue(value)
	}

	zone, _, resource, err := dnsLookupInputs(options)
	if err != nil {
		return nil, err
	}

	return map[string]*Resource{zone: resource}, nil
}

func actionResourceValue(options core.Options, keys ...string) (*Resource, error) {
	value := actionOptionValue(options, keys...)
	if value == nil {
		return nil, core.E(actionOpResource, errResourceRequired, nil)
	}

	return resourceFromAny(value)
}

func resourceFromAny(value any) (*Resource, error) {
	switch typed := value.(type) {
	case *Resource:
		if typed == nil {
			return nil, core.E(actionOpResource, errResourceRequired, nil)
		}
		return typed, nil
	case Resource:
		resource := typed
		return &resource, nil
	case ResourceJSON:
		resource := NewResource()
		if err := resource.FromJSON(typed); err != nil {
			return nil, err
		}
		return resource, nil
	case []byte:
		return DecodeResource(typed)
	case string:
		return resourceFromJSONBytes([]byte(typed))
	case map[string]any:
		raw := core.JSONMarshal(typed)
		if !raw.OK {
			return nil, core.E(actionOpResource, "invalid resource JSON", raw.Value.(error))
		}
		return resourceFromJSONBytes(raw.Value.([]byte))
	default:
		return nil, core.E(actionOpResource, "unsupported resource value", nil)
	}
}

func resourceFromJSONBytes(raw []byte) (*Resource, error) {
	var jsonView ResourceJSON
	if r := core.JSONUnmarshal(raw, &jsonView); !r.OK {
		return nil, core.E(actionOpResource, "invalid resource JSON", r.Value.(error))
	}

	resource := NewResource()
	if err := resource.FromJSON(jsonView); err != nil {
		return nil, err
	}

	return resource, nil
}

func actionResourcesValue(value any) (map[string]*Resource, error) {
	switch typed := value.(type) {
	case map[string]*Resource:
		return typed, nil
	case map[string]Resource:
		out := make(map[string]*Resource, len(typed))
		for key, resource := range typed {
			resourceCopy := resource
			out[dnsFQDN(key)] = &resourceCopy
		}
		return out, nil
	case map[string]any:
		out := make(map[string]*Resource, len(typed))
		for key, raw := range typed {
			resource, err := resourceFromAny(raw)
			if err != nil {
				return nil, err
			}
			out[dnsFQDN(key)] = resource
		}
		return out, nil
	default:
		return nil, core.E("lns.action.resources", "resources must be a map of zone names to resource payloads", nil)
	}
}

func actionHashValue(options core.Options, keys ...string) (primitives.Hash, bool, error) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			hash, ok, err := hashFromAny(result.Value)
			return hash, ok, err
		}
	}
	return primitives.Hash{}, false, nil
}

func hashFromAny(value any) (primitives.Hash, bool, error) {
	switch typed := value.(type) {
	case primitives.Hash:
		return typed, true, nil
	case [32]byte:
		return primitives.Hash(typed), true, nil
	case []byte:
		if len(typed) == 32 {
			var hash primitives.Hash
			copy(hash[:], typed)
			return hash, true, nil
		}
		raw, err := hex.DecodeString(core.Trim(string(typed)))
		if err != nil || len(raw) != 32 {
			return primitives.Hash{}, false, core.E(actionOpHash, "invalid hash encoding", err)
		}
		var hash primitives.Hash
		copy(hash[:], raw)
		return hash, true, nil
	case string:
		raw, err := hex.DecodeString(core.Trim(typed))
		if err != nil || len(raw) != 32 {
			return primitives.Hash{}, false, core.E(actionOpHash, "invalid hash encoding", err)
		}
		var hash primitives.Hash
		copy(hash[:], raw)
		return hash, true, nil
	default:
		if value == nil {
			return primitives.Hash{}, false, nil
		}
		return primitives.Hash{}, false, core.E(actionOpHash, "unsupported hash value", nil)
	}
}

func randomActionHash() (primitives.Hash, error) {
	var hash primitives.Hash
	if _, err := rand.Read(hash[:]); err != nil {
		return primitives.Hash{}, core.E(actionOpHash, "random source failure", err)
	}
	return hash, nil
}

func actionAddressValue(options core.Options, keys ...string) (Address, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return addressFromAny(result.Value)
		}
	}

	hashValue := actionOptionValue(options, "address_hash")
	if hashValue == nil {
		return Address{}, false
	}

	hash, ok := addressHashFromAny(hashValue)
	if !ok {
		return Address{}, false
	}
	version, _ := actionUint32Value(options, "address_version")
	return Address{Version: uint8(version), Hash: hash}, true
}

func addressFromAny(value any) (Address, bool) {
	switch typed := value.(type) {
	case Address:
		return typed, true
	case *Address:
		if typed == nil {
			return Address{}, false
		}
		return *typed, true
	case map[string]any:
		hash, ok := addressHashFromAny(typed["hash"])
		if !ok {
			return Address{}, false
		}
		version, ok := actionUint32FromAny(typed["version"])
		if !ok {
			version = 0
		}
		return Address{Version: uint8(version), Hash: hash}, true
	default:
		return Address{}, false
	}
}

func addressHashFromAny(value any) ([]byte, bool) {
	switch typed := value.(type) {
	case []byte:
		if len(typed) == 0 {
			return nil, false
		}
		return append([]byte(nil), typed...), true
	case string:
		raw, err := hex.DecodeString(core.Trim(typed))
		if err != nil || len(raw) == 0 {
			return nil, false
		}
		return raw, true
	default:
		return nil, false
	}
}

func actionNameRulesValue(options core.Options) (NameRules, bool) {
	rules := NameRules{}
	hasRules := false

	if value := actionOptionValue(options, "rules"); value != nil {
		parsed, ok := nameRulesFromAny(value)
		if ok {
			rules = parsed
			hasRules = true
		}
	}

	if value, ok := actionUint32Value(options, "rollout_interval", "rolloutInterval"); ok {
		rules.RolloutInterval = value
		hasRules = true
	}
	if value, ok := actionUint32Value(options, "auction_start", "auctionStart"); ok {
		rules.AuctionStart = value
		hasRules = true
	}
	if value, ok := actionUint32Value(options, "claim_period", "claimPeriod"); ok {
		rules.ClaimPeriod = value
		hasRules = true
	}
	if value, ok := actionUint32Value(options, "alexa_lockup_period", "alexaLockupPeriod"); ok {
		rules.AlexaLockupPeriod = value
		hasRules = true
	}
	if value, ok := actionBoolValue(options, "no_rollout", "noRollout"); ok {
		rules.NoRollout = value
		hasRules = true
	}
	if value, ok := actionBoolValue(options, "no_reserved", "noReserved"); ok {
		rules.NoReserved = value
		hasRules = true
	}

	return rules, hasRules
}

func nameRulesFromAny(value any) (NameRules, bool) {
	switch typed := value.(type) {
	case NameRules:
		return typed, true
	case map[string]any:
		rules := NameRules{}
		rules.RolloutInterval = actionUint32MapValue(typed, "rollout_interval", "rolloutInterval")
		rules.AuctionStart = actionUint32MapValue(typed, "auction_start", "auctionStart")
		rules.ClaimPeriod = actionUint32MapValue(typed, "claim_period", "claimPeriod")
		rules.AlexaLockupPeriod = actionUint32MapValue(typed, "alexa_lockup_period", "alexaLockupPeriod")
		rules.NoRollout = actionBoolMapValue(typed, "no_rollout", "noRollout")
		rules.NoReserved = actionBoolMapValue(typed, "no_reserved", "noReserved")
		return rules, true
	default:
		return NameRules{}, false
	}
}

func actionUint32MapValue(values map[string]any, keys ...string) uint32 {
	for _, key := range keys {
		if v, ok := actionUint32FromAny(values[key]); ok {
			return v
		}
	}
	return 0
}

func actionBoolMapValue(values map[string]any, keys ...string) bool {
	for _, key := range keys {
		if v, ok := actionBoolFromAny(values[key]); ok {
			return v
		}
	}
	return false
}

func actionOptionValue(options core.Options, keys ...string) any {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return result.Value
		}
	}
	return nil
}

func actionStringFromAny(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		if typed != nil {
			return core.Trim(core.Sprint(typed))
		}
		return ""
	}
}

func actionUint32Value(options core.Options, keys ...string) (uint32, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return actionUint32FromAny(result.Value)
		}
	}
	return 0, false
}

func actionUint64Value(options core.Options, keys ...string) (uint64, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return actionUint64FromAny(result.Value)
		}
	}
	return 0, false
}

func actionUint32FromAny(value any) (uint32, bool) {
	switch typed := value.(type) {
	case uint32:
		return typed, true
	case uint:
		return uint32(typed), true
	case int:
		if typed < 0 {
			return 0, false
		}
		return uint32(typed), true
	case int64:
		if typed < 0 {
			return 0, false
		}
		return uint32(typed), true
	case float64:
		if typed < 0 {
			return 0, false
		}
		return uint32(typed), true
	case string:
		if parsed, err := strconv.ParseUint(core.Trim(typed), 10, 32); err == nil {
			return uint32(parsed), true
		}
	}
	return 0, false
}

func actionUint64FromAny(value any) (uint64, bool) {
	switch typed := value.(type) {
	case uint64:
		return typed, true
	case uint32:
		return uint64(typed), true
	case uint:
		return uint64(typed), true
	case int:
		if typed < 0 {
			return 0, false
		}
		return uint64(typed), true
	case int64:
		if typed < 0 {
			return 0, false
		}
		return uint64(typed), true
	case float64:
		if typed < 0 {
			return 0, false
		}
		return uint64(typed), true
	case string:
		if parsed, err := strconv.ParseUint(core.Trim(typed), 10, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func actionBoolValue(options core.Options, keys ...string) (bool, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return actionBoolFromAny(result.Value)
		}
	}
	return false, false
}

func actionBoolFromAny(value any) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		parsed, err := strconv.ParseBool(core.Trim(typed))
		if err != nil {
			return false, false
		}
		return parsed, true
	default:
		return false, false
	}
}

func dnsFQDN(name string) string {
	name = core.Lower(core.Trim(name))
	if name == "" {
		return "."
	}
	if core.HasSuffix(name, ".") {
		return name
	}
	return name + "."
}

func dnsNameMatches(recordName string, want string) bool {
	return dnsFQDN(recordName) == dnsFQDN(want)
}

func (s *Service) vpnGatewaysFromOptions(ctx context.Context, options core.Options) ([]vpnGatewayResult, error) {
	records, err := s.vpnRecordsFromOptions(ctx, options)
	if err != nil {
		return nil, err
	}

	filter := canonicalVPNFilter(actionStringFromAny(actionOptionValue(options, "name", "gateway", "_input")))
	gateways := make([]vpnGatewayResult, 0, len(records))
	for _, record := range records {
		gateway, ok := vpnGatewayFromRecord(record)
		if !ok {
			continue
		}
		if filter != "" && gateway.Name != filter && gateway.FQDN != dnsFQDN(filter) {
			continue
		}
		gateways = append(gateways, gateway)
	}

	sort.Slice(gateways, func(i, j int) bool {
		return gateways[i].Name < gateways[j].Name
	})

	return gateways, nil
}

func (s *Service) vpnRecordsFromOptions(ctx context.Context, options core.Options) ([]dnspkg.NameRecord, error) {
	if value := actionOptionValue(options, "resources"); value != nil {
		return vpnRecordsFromResourceMap(value)
	}

	if value := actionOptionValue(options, "resource", "resource_json"); value != nil {
		return vpnRecordFromResource(value, options)
	}

	bridge, err := s.bridgeFromOptions(options)
	if err != nil {
		return nil, core.E(actionNameVPNDiscover, "chain_url or resource data is required", err)
	}
	if bridge.Health().Names == 0 {
		if _, err := bridge.Discover(ctx); err != nil {
			return nil, err
		}
	}

	snapshot := bridge.Records()
	records := append([]dnspkg.NameRecord(nil), snapshot...)
	return records, nil
}

func vpnRecordsFromResourceMap(value any) ([]dnspkg.NameRecord, error) {
	resources, err := actionResourcesValue(value)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(resources))
	for zone := range resources {
		names = append(names, zone)
	}
	sort.Strings(names)

	records := make([]dnspkg.NameRecord, 0, len(names))
	for _, zone := range names {
		label := core.TrimSuffix(core.TrimSuffix(core.Lower(zone), ".lthn."), ".lthn")
		records = append(records, dnspkg.NameRecord{Name: label, Resource: resources[zone], Source: "resource"})
	}
	return records, nil
}

func vpnRecordFromResource(value any, options core.Options) ([]dnspkg.NameRecord, error) {
	resource, err := resourceFromAny(value)
	if err != nil {
		return nil, err
	}
	label, ok := canonicalizeName(actionOptionValue(options, "name", "zone", "_input"))
	if !ok {
		return nil, core.E(actionNameVPNDiscover, "name is required with resource", nil)
	}
	return []dnspkg.NameRecord{{Name: label, Resource: resource, Source: "resource"}}, nil
}

func vpnGatewayFromRecord(record dnspkg.NameRecord) (vpnGatewayResult, bool) {
	txt := record.TXT()
	metadata := chainpkg.ParseAliasComment(record.Comment)
	for key, value := range vpnMetadataFromTXT(txt) {
		if _, exists := metadata[key]; !exists {
			metadata[key] = value
		}
	}

	capabilities := uniqueStrings(splitCSV(metadata["cap"]))
	if metadata["vpn"] != "" && !containsString(capabilities, "vpn") {
		capabilities = append(capabilities, "vpn")
	}

	gatewayType := core.Lower(core.Trim(metadata["type"]))
	if gatewayType == "" && metadata["vpn"] != "" {
		gatewayType = "gateway"
	}

	isGateway := gatewayType == "gateway" || containsString(capabilities, "vpn") || txtLooksVPN(txt)
	if !isGateway {
		return vpnGatewayResult{}, false
	}

	addresses := append(record.Addresses4(), record.Addresses6()...)
	return vpnGatewayResult{
		Name:         core.Lower(record.Name),
		FQDN:         record.FQDN(),
		Source:       record.Source,
		Type:         gatewayType,
		Capabilities: capabilities,
		Addresses:    addresses,
		TXT:          txt,
		Metadata:     metadata,
	}, true
}

func vpnMetadataFromTXT(entries []string) map[string]string {
	metadata := make(map[string]string)
	for _, entry := range entries {
		for _, part := range splitMetadataFields(entry) {
			key, value, ok := cutKeyValue(core.Trim(part))
			if !ok {
				continue
			}
			key = core.Lower(core.Trim(key))
			value = core.Trim(value)
			if key == "" || value == "" {
				continue
			}
			metadata[key] = value
		}
	}
	return metadata
}

func splitMetadataFields(entry string) []string {
	var out []string
	start := 0
	for i, r := range entry {
		if r != ';' && r != ',' {
			continue
		}
		out = append(out, entry[start:i])
		start = i + len(string(r))
	}
	return append(out, entry[start:])
}

func cutKeyValue(value string) (string, string, bool) {
	for i, r := range value {
		if r == '=' {
			return value[:i], value[i+len(string(r)):], true
		}
	}
	return value, "", false
}

func txtLooksVPN(entries []string) bool {
	for _, entry := range entries {
		lower := core.Lower(core.Trim(entry))
		switch {
		case core.Contains(lower, "vpn"):
			return true
		case core.HasPrefix(lower, "wg="),
			core.HasPrefix(lower, "wireguard="),
			core.HasPrefix(lower, "openvpn="),
			core.HasPrefix(lower, "ovpn="),
			core.HasPrefix(lower, "endpoint="):
			return true
		}
	}
	return false
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := core.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = core.Lower(core.Trim(part))
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func containsString(values []string, want string) bool {
	want = core.Lower(core.Trim(want))
	for _, value := range values {
		if core.Lower(core.Trim(value)) == want {
			return true
		}
	}
	return false
}

func canonicalVPNFilter(value string) string {
	value = core.Lower(core.Trim(value))
	value = core.TrimSuffix(value, ".")
	value = core.TrimSuffix(value, ".lthn")
	return value
}

func selectVPNGateway(gateways []vpnGatewayResult, filter string) (vpnGatewayResult, bool) {
	if len(gateways) == 0 {
		return vpnGatewayResult{}, false
	}
	if filter == "" {
		return gateways[0], true
	}

	filter = canonicalVPNFilter(filter)
	for _, gateway := range gateways {
		if gateway.Name == filter || canonicalVPNFilter(gateway.FQDN) == filter {
			return gateway, true
		}
	}
	return vpnGatewayResult{}, false
}

func (s *Service) setVPNConnection(gateway vpnGatewayResult) vpnStatusResult {
	if s == nil {
		return vpnStatusResult{Connected: false, Mode: vpnModeLocalSession}
	}

	status := vpnStatusResult{
		Connected: true,
		Gateway:   &gateway,
		Since:     time.Now().UTC().Format(time.RFC3339),
		Mode:      vpnModeLocalSession,
	}

	s.vpnMu.Lock()
	s.vpnConnection = &vpnConnectionState{
		Gateway:     gateway,
		ConnectedAt: status.Since,
	}
	s.vpnMu.Unlock()
	return status
}

func (s *Service) currentVPNStatus() vpnStatusResult {
	status := vpnStatusResult{
		Connected: false,
		Mode:      vpnModeLocalSession,
		Notes:     []string{"No VPN gateway is currently selected."},
	}

	if s == nil {
		return status
	}

	s.vpnMu.Lock()
	defer s.vpnMu.Unlock()

	if s.vpnConnection == nil {
		return status
	}

	gateway := s.vpnConnection.Gateway
	status.Connected = true
	status.Gateway = &gateway
	status.Since = s.vpnConnection.ConnectedAt
	status.Notes = []string{"Gateway selection is local to this service instance."}
	return status
}

func (s *Service) clearVPNConnection() vpnStatusResult {
	status := vpnStatusResult{
		Connected: false,
		Mode:      vpnModeLocalSession,
		Notes:     []string{"VPN gateway selection was cleared."},
	}

	if s == nil {
		return status
	}

	s.vpnMu.Lock()
	s.vpnConnection = nil
	s.vpnMu.Unlock()
	return status
}
