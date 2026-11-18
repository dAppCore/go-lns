// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/hex"

	core "github.com/dappcore/go"
)

// OutpointJSON represents the JSON form of an outpoint.
type OutpointJSON struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

// GetJSON converts the outpoint into its JSON representation.
func (o Outpoint) GetJSON() OutpointJSON {
	return OutpointJSON{
		Hash:  hex.EncodeToString(o.TxHash[:]),
		Index: o.Index,
	}
}

// FromJSON populates the outpoint from its JSON representation.
func (o *Outpoint) FromJSON(json OutpointJSON) error {
	if len(json.Hash) != len(Hash{})*2 {
		return core.E(opOutpointFromJSON, "invalid hash length", nil)
	}

	raw, err := hex.DecodeString(json.Hash)
	if err != nil {
		return core.E(opOutpointFromJSON, "invalid hash encoding", err)
	}

	if len(raw) != len(Hash{}) {
		return core.E(opOutpointFromJSON, "invalid hash length", nil)
	}

	copy(o.TxHash[:], raw)
	o.Index = json.Index
	return nil
}
