// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/hex"

	core "github.com/dappcore/go"
)

const nameStateJSONMaxSize = 63

// NameStateJSON represents the JSON form of a name state.
type NameStateJSON struct {
	Name       string          `json:"name"`
	NameHash   string          `json:"nameHash"`
	State      string          `json:"state,omitempty"`
	Height     uint32          `json:"height"`
	Renewal    uint32          `json:"renewal"`
	Owner      OutpointJSON    `json:"owner"`
	Value      uint64          `json:"value"`
	Highest    uint64          `json:"highest"`
	Data       string          `json:"data"`
	Transfer   uint32          `json:"transfer"`
	Revoked    uint32          `json:"revoked"`
	Claimed    uint32          `json:"claimed"`
	Renewals   uint32          `json:"renewals"`
	Registered bool            `json:"registered"`
	Expired    bool            `json:"expired"`
	Weak       bool            `json:"weak"`
	Stats      *NameStateStats `json:"stats,omitempty"`
}

// GetJSON converts the name state into its JSON representation at height.
func (ns NameState) GetJSON(height uint32, rules NameStateRules) NameStateJSON {
	return ns.Format(height, rules)
}

// Format converts the name state into its JSON representation at height.
func (ns NameState) Format(height uint32, rules NameStateRules) NameStateJSON {
	state := ns.State(height, rules).String()
	stats := ns.ToStats(height, rules)

	return NameStateJSON{
		Name:       string(ns.Name),
		NameHash:   hex.EncodeToString(ns.NameHash[:]),
		State:      state,
		Height:     ns.Height,
		Renewal:    ns.Renewal,
		Owner:      ns.Owner.GetJSON(),
		Value:      ns.Value,
		Highest:    ns.Highest,
		Data:       hex.EncodeToString(ns.Data),
		Transfer:   ns.Transfer,
		Revoked:    ns.Revoked,
		Claimed:    ns.Claimed,
		Renewals:   ns.Renewals,
		Registered: ns.Registered,
		Expired:    ns.Expired,
		Weak:       ns.Weak,
		Stats:      stats,
	}
}

// GetFormat is an alias for Format.
func (ns NameState) GetFormat(height uint32, rules NameStateRules) NameStateJSON {
	return ns.Format(height, rules)
}

// FromJSON populates the name state from its JSON representation.
func (ns *NameState) FromJSON(json NameStateJSON) error {
	if len(json.Name) > nameStateJSONMaxSize {
		return core.E(opNameStateFromJSON, "invalid name length", nil)
	}

	if len(json.NameHash) != len(Hash{})*2 {
		return core.E(opNameStateFromJSON, "invalid name hash length", nil)
	}

	nameHash, err := hex.DecodeString(json.NameHash)
	if err != nil {
		return core.E(opNameStateFromJSON, "invalid name hash encoding", err)
	}

	data, err := hex.DecodeString(json.Data)
	if err != nil {
		return core.E(opNameStateFromJSON, "invalid data encoding", err)
	}

	if err := ns.Owner.FromJSON(json.Owner); err != nil {
		return err
	}

	copy(ns.NameHash[:], nameHash)
	ns.Name = []byte(json.Name)
	ns.Height = json.Height
	ns.Renewal = json.Renewal
	ns.Value = json.Value
	ns.Highest = json.Highest
	ns.Data = data
	ns.Transfer = json.Transfer
	ns.Revoked = json.Revoked
	ns.Claimed = json.Claimed
	ns.Renewals = json.Renewals
	ns.Registered = json.Registered
	ns.Expired = json.Expired
	ns.Weak = json.Weak
	ns.delta = nil
	return nil
}

// GetFromJSON is an alias for FromJSON.
func (ns *NameState) GetFromJSON(json NameStateJSON) error {
	return ns.FromJSON(json)
}
