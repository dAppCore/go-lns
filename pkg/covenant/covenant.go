// SPDX-License-Identifier: EUPL-1.2

// Package covenant defines the covenant type constants for the Lethean Name System.
// Covenants encode name-system operations (auctions, registrations, transfers)
// as typed payloads attached to transaction outputs.
//
// Check a covenant type:
//
//	if covType == covenant.TypeRegister { /* name registration */ }
//
// Look up a type name:
//
//	name := covenant.TypeName(covenant.TypeBid)  // "BID"
package covenant

// CovenantType represents the type of name-system operation encoded in a covenant.
// Values mirror the HSD rules.js type enumeration.
//
//	var ct covenant.CovenantType = covenant.TypeOpen
type CovenantType uint8

const (
	// TypeNone is a plain transfer with no name-system operation.
	TypeNone CovenantType = 0

	// TypeClaim imports an existing ICANN/Alexa name into the system.
	TypeClaim CovenantType = 1

	// TypeOpen begins a name auction, making the name available for bidding.
	TypeOpen CovenantType = 2

	// TypeBid places a sealed bid during the auction's bidding phase.
	TypeBid CovenantType = 3

	// TypeReveal reveals a previously sealed bid after the bidding phase closes.
	TypeReveal CovenantType = 4

	// TypeRedeem reclaims coins from a losing bid after the reveal phase.
	TypeRedeem CovenantType = 5

	// TypeRegister completes name registration after winning an auction.
	TypeRegister CovenantType = 6

	// TypeUpdate modifies the DNS resource data for a registered name.
	TypeUpdate CovenantType = 7

	// TypeRenew extends the registration period for a name.
	TypeRenew CovenantType = 8

	// TypeTransfer initiates a name transfer to a new address.
	TypeTransfer CovenantType = 9

	// TypeFinalize completes a name transfer after the lock-up period.
	TypeFinalize CovenantType = 10

	// TypeRevoke permanently burns a name, making it unspendable and
	// available for re-auction.
	TypeRevoke CovenantType = 11
)

// Types mirrors the JS rules.types lookup table.
//
// The map is intentionally exported so callers can discover covenant types
// without hard-coding the numeric values.
var Types = map[string]CovenantType{
	"NONE":     TypeNone,
	"CLAIM":    TypeClaim,
	"OPEN":     TypeOpen,
	"BID":      TypeBid,
	"REVEAL":   TypeReveal,
	"REDEEM":   TypeRedeem,
	"REGISTER": TypeRegister,
	"UPDATE":   TypeUpdate,
	"RENEW":    TypeRenew,
	"TRANSFER": TypeTransfer,
	"FINALIZE": TypeFinalize,
	"REVOKE":   TypeRevoke,
}

// GetTypes returns the covenant type lookup table.
func GetTypes() map[string]CovenantType {
	return Types
}

// TypesByVal mirrors the JS rules.typesByVal lookup table.
var TypesByVal = map[CovenantType]string{
	TypeNone:     "NONE",
	TypeClaim:    "CLAIM",
	TypeOpen:     "OPEN",
	TypeBid:      "BID",
	TypeReveal:   "REVEAL",
	TypeRedeem:   "REDEEM",
	TypeRegister: "REGISTER",
	TypeUpdate:   "UPDATE",
	TypeRenew:    "RENEW",
	TypeTransfer: "TRANSFER",
	TypeFinalize: "FINALIZE",
	TypeRevoke:   "REVOKE",
}

// GetTypesByVal returns the reverse covenant type lookup table.
func GetTypesByVal() map[CovenantType]string {
	return TypesByVal
}

// TypeName returns the human-readable name for a covenant type.
// Returns "UNKNOWN" for unrecognised type values.
//
//	covenant.TypeName(covenant.TypeBid)    // "BID"
//	covenant.TypeName(covenant.TypeRevoke) // "REVOKE"
//	covenant.TypeName(99)                 // "UNKNOWN"
func TypeName(ct CovenantType) string {
	if name, ok := TypesByVal[ct]; ok {
		return name
	}
	return "UNKNOWN"
}

// GetTypeName is an alias for TypeName.
//
//	name := covenant.GetTypeName(covenant.TypeBid)
func GetTypeName(ct CovenantType) string {
	return TypeName(ct)
}

// IsName returns true if the covenant type is a name-related operation
// (anything from CLAIM through REVOKE, inclusive).
//
//	covenant.TypeOpen.IsName()  // true
//	covenant.TypeNone.IsName()  // false
func (ct CovenantType) IsName() bool {
	return ct >= TypeClaim && ct <= TypeRevoke
}

// IsKnown returns true if the covenant type is a recognised type
// (NONE through REVOKE, inclusive).
//
//	covenant.TypeBid.IsKnown()       // true
//	covenant.CovenantType(99).IsKnown() // false
func (ct CovenantType) IsKnown() bool {
	return ct <= TypeRevoke
}

// IsLinked returns true if the covenant type is part of the linked set
// (REVEAL through REVOKE, inclusive).
//
//	covenant.TypeReveal.IsLinked() // true
//	covenant.TypeBid.IsLinked()    // false
func (ct CovenantType) IsLinked() bool {
	return ct >= TypeReveal && ct <= TypeRevoke
}

// String returns the human-readable name for a covenant type.
// Satisfies the fmt.Stringer interface.
//
//	ct := covenant.TypeRegister
//	println(ct.String())  // "REGISTER"
func (ct CovenantType) String() string {
	return TypeName(ct)
}
