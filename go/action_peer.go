// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"encoding/hex"
	"strconv"

	core "github.com/dappcore/go"
	netpkg "github.com/dappcore/go-lns/go/pkg/net"
)

const (
	errPeerNameRequired = "peer name is required"
	errPublicKeyMissing = "public_key is required"
	errGossipMissing    = "gossip, message, or payload is required"
)

// peerRecordResult is the wire view of a registered peer record returned by the
// lns.peer.* actions.
type peerRecordResult struct {
	Name       string   `json:"name"`
	PublicKey  string   `json:"publicKey,omitempty"`
	Addresses  []string `json:"addresses,omitempty"`
	Cap        []string `json:"cap,omitempty"`
	TrustScore float64  `json:"trustScore"`
}

// peerLookupResult reports whether a peer lookup matched and the matching record.
type peerLookupResult struct {
	Name  string            `json:"name"`
	Found bool              `json:"found"`
	Peer  *peerRecordResult `json:"peer,omitempty"`
}

// peerListResult is the response for lns.peer.list: the matching peers ordered
// by descending trust score when a capability filter is supplied.
type peerListResult struct {
	Capability string             `json:"capability,omitempty"`
	Count      int                `json:"count"`
	Peers      []peerRecordResult `json:"peers"`
}

// gossipReceiveResult reports the outcome of merging an inbound gossip delta
// into the local cache (RFC §12.1). The counts mirror net.MergeResult.
type gossipReceiveResult struct {
	Peer                 string `json:"peer"`
	Merged               int    `json:"merged"`
	RejectedUnauthorised int    `json:"rejectedUnauthorised"`
	RejectedInvalid      int    `json:"rejectedInvalid"`
	SkippedHSD           int    `json:"skippedHsd"`
	CacheSize            int    `json:"cacheSize"`
}

// gossipAnnounceResult reports the outcome of an outbound gossip broadcast
// (RFC §12.1). Broadcast is false when no transport is configured (offline).
type gossipAnnounceResult struct {
	Peer      string   `json:"peer"`
	Records   int      `json:"records"`
	Broadcast bool     `json:"broadcast"`
	Notes     []string `json:"notes,omitempty"`
}

// peerRecordToResult projects a net.PeerRecord into its hex-encoded wire view.
func peerRecordToResult(rec netpkg.PeerRecord) peerRecordResult {
	out := peerRecordResult{
		Name:       rec.Name,
		Addresses:  rec.Addresses,
		Cap:        rec.Cap,
		TrustScore: rec.TrustScore,
	}
	if len(rec.PublicKey) > 0 {
		out.PublicKey = hex.EncodeToString(rec.PublicKey)
	}
	return out
}

func (s *Service) handlePeerRegister(_ context.Context, options core.Options) core.Result {
	name := core.Trim(actionStringFromAny(actionOptionValue(options, "name", "peer", "_input")))
	if name == "" {
		return core.Fail(core.E(actionNamePeerRegister, errPeerNameRequired, nil))
	}

	publicKey, ok := actionBytesValue(options, "public_key", "publicKey", "pubkey")
	if !ok || len(publicKey) == 0 {
		return core.Fail(core.E(actionNamePeerRegister, errPublicKeyMissing, nil))
	}

	record := netpkg.PeerRecord{
		Name:      name,
		PublicKey: publicKey,
		Addresses: actionStringSliceValue(actionOptionValue(options, "addresses", "address")),
		Cap:       actionStringSliceValue(actionOptionValue(options, "cap", "capabilities")),
	}

	registered, err := s.peerRegistryHandle().Register(record, peerVectorFromOptions(options))
	if err != nil {
		return core.Fail(core.E(actionNamePeerRegister, "register peer", err))
	}

	return core.Ok(peerRecordToResult(registered))
}

func (s *Service) handlePeerLookup(_ context.Context, options core.Options) core.Result {
	name := core.Trim(actionStringFromAny(actionOptionValue(options, "name", "peer", "_input")))
	if name == "" {
		return core.Fail(core.E("lns.peer.lookup", errPeerNameRequired, nil))
	}

	rec, found := s.peerRegistryHandle().Lookup(name)
	result := peerLookupResult{Name: name, Found: found}
	if found {
		view := peerRecordToResult(rec)
		result.Peer = &view
	}

	return core.Ok(result)
}

func (s *Service) handlePeerList(_ context.Context, options core.Options) core.Result {
	registry := s.peerRegistryHandle()

	capability := core.Trim(actionStringFromAny(actionOptionValue(options, "cap", "capability")))

	var records []netpkg.PeerRecord
	if capability != "" {
		records = registry.WithCapability(capability)
	} else {
		for rec := range registry.Records() {
			records = append(records, rec)
		}
	}

	peers := make([]peerRecordResult, 0, len(records))
	for _, rec := range records {
		peers = append(peers, peerRecordToResult(rec))
	}

	return core.Ok(peerListResult{
		Capability: capability,
		Count:      len(peers),
		Peers:      peers,
	})
}

func (s *Service) handleGossipReceive(_ context.Context, options core.Options) core.Result {
	gossip, err := gossipFromOptions(actionNameGossipReceive, options)
	if err != nil {
		return core.Fail(err)
	}

	authority, err := gossipAuthorityFromOptions(options, gossip)
	if err != nil {
		return core.Fail(err)
	}

	registry := s.peerRegistryHandle()
	merge, err := registry.Receive(gossip, authority)
	if err != nil {
		return core.Fail(core.E(actionNameGossipReceive, "merge gossip", err))
	}

	return core.Ok(gossipReceiveResult{
		Peer:                 gossip.Peer,
		Merged:               merge.Merged,
		RejectedUnauthorised: merge.RejectedUnauthorised,
		RejectedInvalid:      merge.RejectedInvalid,
		SkippedHSD:           merge.SkippedHSD,
		CacheSize:            registry.Cache().Len(),
	})
}

func (s *Service) handleGossipAnnounce(_ context.Context, options core.Options) core.Result {
	gossip, err := gossipFromOptions(actionNameGossipAnnounce, options)
	if err != nil {
		return core.Fail(err)
	}

	result := gossipAnnounceResult{
		Peer:    gossip.Peer,
		Records: len(gossip.Records),
	}

	registry := s.peerRegistryHandle()
	announce := registry.Announce(gossip)
	if !announce.OK {
		return announce
	}

	if registry.Transport() == nil {
		result.Broadcast = false
		result.Notes = append(result.Notes,
			"no go-p2p transport is configured; the signed delta was accepted locally but not broadcast (offline mode).")
		return core.Ok(result)
	}

	result.Broadcast = true
	return core.Ok(result)
}

// peerVectorFromOptions parses the observable trust-score inputs (RFC §11.1)
// from action options, defaulting every field to its zero value.
func peerVectorFromOptions(options core.Options) netpkg.PeerVector {
	vector := netpkg.PeerVector{}

	if value, ok := actionFloatValue(options, "uptime", "uptime_seconds"); ok {
		vector.Uptime = value
	}
	if value, ok := actionBoolValue(options, "verified_chain", "verifiedChain"); ok {
		vector.VerifiedChain = value
	}
	if value, ok := actionBoolValue(options, "verified_hns", "verifiedHns"); ok {
		vector.VerifiedHNS = value
	}
	if value, ok := actionFloatValue(options, "response_time", "response_ms"); ok {
		vector.ResponseTime = value
	}
	if value, ok := actionUint32Value(options, "consistency", "consistency_n"); ok {
		vector.ConsistencyN = int(value)
	}

	return vector
}

// gossipFromOptions extracts a signed gossip message from action options. It
// accepts a full transport envelope (message), a decoded gossip object/string
// (gossip), or the raw JSON-encoded gossip payload (payload).
func gossipFromOptions(op string, options core.Options) (*netpkg.Gossip, error) {
	if value := actionOptionValue(options, "message"); value != nil {
		message, ok := gossipMessageFromAny(value)
		if !ok {
			return nil, core.E(op, "invalid gossip message envelope", nil)
		}
		gossip, err := netpkg.DecodeGossipMessage(message)
		if err != nil {
			return nil, core.E(op, "decode gossip message", err)
		}
		return gossip, nil
	}

	if value := actionOptionValue(options, "gossip"); value != nil {
		return gossipFromAny(op, value)
	}

	if value := actionOptionValue(options, "payload"); value != nil {
		raw, ok := actionRawBytes(value)
		if !ok || len(raw) == 0 {
			return nil, core.E(op, "invalid gossip payload", nil)
		}
		var gossip netpkg.Gossip
		if res := core.JSONUnmarshal(raw, &gossip); !res.OK {
			err, _ := res.Value.(error)
			return nil, core.E(op, "decode gossip payload", err)
		}
		return &gossip, nil
	}

	return nil, core.E(op, errGossipMissing, nil)
}

// gossipFromAny coerces a gossip option value (a *Gossip, Gossip, JSON string,
// JSON bytes, or decoded map) into a *netpkg.Gossip.
func gossipFromAny(op string, value any) (*netpkg.Gossip, error) {
	switch typed := value.(type) {
	case *netpkg.Gossip:
		if typed == nil {
			return nil, core.E(op, errGossipMissing, nil)
		}
		return typed, nil
	case netpkg.Gossip:
		clone := typed
		return &clone, nil
	case string:
		return decodeGossipJSON(op, []byte(typed))
	case []byte:
		return decodeGossipJSON(op, typed)
	case map[string]any:
		encoded := core.JSONMarshal(typed)
		if !encoded.OK {
			return nil, core.E(op, "encode gossip object", nil)
		}
		raw, _ := encoded.Value.([]byte)
		return decodeGossipJSON(op, raw)
	default:
		return nil, core.E(op, "unsupported gossip value", nil)
	}
}

func decodeGossipJSON(op string, raw []byte) (*netpkg.Gossip, error) {
	if len(raw) == 0 {
		return nil, core.E(op, errGossipMissing, nil)
	}
	var gossip netpkg.Gossip
	if res := core.JSONUnmarshal(raw, &gossip); !res.OK {
		err, _ := res.Value.(error)
		return nil, core.E(op, "decode gossip", err)
	}
	return &gossip, nil
}

// gossipMessageFromAny coerces a message option value into a net.GossipMessage.
func gossipMessageFromAny(value any) (netpkg.GossipMessage, bool) {
	switch typed := value.(type) {
	case netpkg.GossipMessage:
		return typed, true
	case *netpkg.GossipMessage:
		if typed == nil {
			return netpkg.GossipMessage{}, false
		}
		return *typed, true
	case map[string]any:
		message := netpkg.GossipMessage{
			Type: actionStringFromAny(typed["type"]),
			From: actionStringFromAny(typed["from"]),
			To:   actionStringFromAny(typed["to"]),
		}
		if message.Type == "" {
			message.Type = netpkg.GossipMessageType
		}
		if raw, ok := actionRawBytes(typed["payload"]); ok {
			message.Payload = raw
		}
		return message, len(message.Payload) > 0
	default:
		return netpkg.GossipMessage{}, false
	}
}

// gossipAuthorityFromOptions builds the merge authority (RFC §12.1) from action
// options, defaulting the peer label to the gossip's own peer when omitted.
func gossipAuthorityFromOptions(options core.Options, gossip *netpkg.Gossip) (netpkg.GossipAuthority, error) {
	peer := core.Trim(actionStringFromAny(actionOptionValue(options, "peer", "from")))
	if peer == "" && gossip != nil {
		peer = gossip.Peer
	}
	if peer == "" {
		return netpkg.GossipAuthority{}, core.E(actionNameGossipReceive, errPeerNameRequired, nil)
	}

	publicKey, ok := actionBytesValue(options, "public_key", "publicKey", "pubkey")
	if !ok || len(publicKey) == 0 {
		return netpkg.GossipAuthority{}, core.E(actionNameGossipReceive, errPublicKeyMissing, nil)
	}

	return netpkg.GossipAuthority{
		Peer:            peer,
		PublicKey:       publicKey,
		Address:         core.Trim(actionStringFromAny(actionOptionValue(options, "address", "peer_address"))),
		AuthorizedNames: actionStringSliceValue(actionOptionValue(options, "authorized_names", "authorizedNames")),
	}, nil
}

// actionBytesValue reads a byte slice option, accepting raw bytes or a hex
// string under any of the supplied keys.
func actionBytesValue(options core.Options, keys ...string) ([]byte, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return actionRawBytes(result.Value)
		}
	}
	return nil, false
}

// actionRawBytes coerces an option value into bytes, decoding a hex string when
// the value is textual.
func actionRawBytes(value any) ([]byte, bool) {
	switch typed := value.(type) {
	case []byte:
		if len(typed) == 0 {
			return nil, false
		}
		return append([]byte(nil), typed...), true
	case string:
		trimmed := core.Trim(typed)
		if trimmed == "" {
			return nil, false
		}
		raw, err := hex.DecodeString(trimmed)
		if err != nil {
			return nil, false
		}
		return raw, true
	default:
		return nil, false
	}
}

// actionFloatValue reads a float64 option under any of the supplied keys.
func actionFloatValue(options core.Options, keys ...string) (float64, bool) {
	for _, key := range keys {
		if result := options.Get(key); result.OK {
			return actionFloatFromAny(result.Value)
		}
	}
	return 0, false
}

func actionFloatFromAny(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case string:
		if parsed, err := strconv.ParseFloat(core.Trim(typed), 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}
