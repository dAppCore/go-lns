// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/hex"

	core "github.com/dappcore/go"
)

// CovenantJSON represents the JSON form of a covenant.
type CovenantJSON struct {
	Type   uint8    `json:"type"`
	Action string   `json:"action"`
	Items  []string `json:"items"`
}

func covenantTypeName(t uint8) string {
	switch t {
	case covenantTypeNone:
		return "NONE"
	case covenantTypeClaim:
		return "CLAIM"
	case covenantTypeOpen:
		return "OPEN"
	case covenantTypeBid:
		return "BID"
	case covenantTypeReveal:
		return "REVEAL"
	case covenantTypeRedeem:
		return "REDEEM"
	case covenantTypeRegister:
		return "REGISTER"
	case covenantTypeUpdate:
		return "UPDATE"
	case covenantTypeRenew:
		return "RENEW"
	case covenantTypeTransfer:
		return "TRANSFER"
	case covenantTypeFinalize:
		return "FINALIZE"
	case covenantTypeRevoke:
		return "REVOKE"
	default:
		return "UNKNOWN"
	}
}

// GetJSON converts the covenant into its JSON representation.
func (c Covenant) GetJSON() CovenantJSON {
	items := make([]string, len(c.Items))
	for i, item := range c.Items {
		items[i] = hex.EncodeToString(item)
	}

	return CovenantJSON{
		Type:   c.Type,
		Action: covenantTypeName(c.Type),
		Items:  items,
	}
}

// FromJSON populates the covenant from its JSON representation.
func (c *Covenant) FromJSON(json CovenantJSON) error {
	c.Type = json.Type
	c.Items = nil

	if len(json.Items) == 0 {
		return nil
	}

	c.Items = make([][]byte, 0, len(json.Items))
	for _, str := range json.Items {
		item, err := hex.DecodeString(str)
		if err != nil {
			return core.E("primitives.Covenant.FromJSON", "invalid item encoding", err)
		}

		c.Items = append(c.Items, item)
	}

	return nil
}
