// SPDX-License-Identifier: EUPL-1.2

package dns

// DUMMY mirrors the JS zero-length buffer export used by the DNS reference
// helpers.
var DUMMY = []byte{}

// GetDummy returns the zero-length DNS reference buffer.
func GetDummy() []byte {
	return DUMMY
}

// DEFAULT_TTL mirrors the reference DNS TTL used for synthesized records.
const DEFAULT_TTL = 21600

// GetDefaultTTL returns the reference DNS TTL used for synthesized records.
func GetDefaultTTL() int {
	return DEFAULT_TTL
}

// Type map bitmaps mirror the JS DNS reference constants.
var (
	TYPE_MAP_ROOT  = []byte{0x00, 0x07, 0x22, 0x00, 0x00, 0x00, 0x00, 0x03, 0x80}
	TYPE_MAP_EMPTY = []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}
	TYPE_MAP_NS    = []byte{0x00, 0x06, 0x20, 0x00, 0x00, 0x00, 0x00, 0x03}
	TYPE_MAP_TXT   = []byte{0x00, 0x06, 0x00, 0x00, 0x80, 0x00, 0x00, 0x03}
	TYPE_MAP_A     = []byte{0x00, 0x06, 0x40, 0x00, 0x00, 0x00, 0x00, 0x03}
	TYPE_MAP_AAAA  = []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x08, 0x00, 0x03}
)

// GetTypeMapRoot returns the reference type bitmap for the root zone.
func GetTypeMapRoot() []byte {
	return TYPE_MAP_ROOT
}

// GetTypeMapEmpty returns the reference type bitmap for an empty response.
func GetTypeMapEmpty() []byte {
	return TYPE_MAP_EMPTY
}

// GetTypeMapNS returns the reference type bitmap for NS responses.
func GetTypeMapNS() []byte {
	return TYPE_MAP_NS
}

// GetTypeMapTXT returns the reference type bitmap for TXT responses.
func GetTypeMapTXT() []byte {
	return TYPE_MAP_TXT
}

// GetTypeMapA returns the reference type bitmap for A responses.
func GetTypeMapA() []byte {
	return TYPE_MAP_A
}

// GetTypeMapAAAA returns the reference type bitmap for AAAA responses.
func GetTypeMapAAAA() []byte {
	return TYPE_MAP_AAAA
}

// HSType mirrors the JS hsTypes record-type enum.
type HSType uint8

const (
	// HSTypeDS corresponds to DS records.
	HSTypeDS HSType = iota
	// HSTypeNS corresponds to NS records.
	HSTypeNS
	// HSTypeGLUE4 corresponds to IPv4 glue records.
	HSTypeGLUE4
	// HSTypeGLUE6 corresponds to IPv6 glue records.
	HSTypeGLUE6
	// HSTypeSYNTH4 corresponds to synthesized IPv4 records.
	HSTypeSYNTH4
	// HSTypeSYNTH6 corresponds to synthesized IPv6 records.
	HSTypeSYNTH6
	// HSTypeTXT corresponds to TXT records.
	HSTypeTXT
)

// HSTypes mirrors the JS hsTypes lookup table.
var HSTypes = map[string]HSType{
	"DS":     HSTypeDS,
	"NS":     HSTypeNS,
	"GLUE4":  HSTypeGLUE4,
	"GLUE6":  HSTypeGLUE6,
	"SYNTH4": HSTypeSYNTH4,
	"SYNTH6": HSTypeSYNTH6,
	"TXT":    HSTypeTXT,
}

// HSTypesByVal mirrors the JS hsTypesByVal reverse lookup table.
var HSTypesByVal = map[HSType]string{
	HSTypeDS:     "DS",
	HSTypeNS:     "NS",
	HSTypeGLUE4:  "GLUE4",
	HSTypeGLUE6:  "GLUE6",
	HSTypeSYNTH4: "SYNTH4",
	HSTypeSYNTH6: "SYNTH6",
	HSTypeTXT:    "TXT",
}

// GetHSTypes returns the DNS record-type lookup table.
func GetHSTypes() map[string]HSType {
	return HSTypes
}

// GetHSTypesByVal returns the reverse DNS record-type lookup table.
func GetHSTypesByVal() map[HSType]string {
	return HSTypesByVal
}
