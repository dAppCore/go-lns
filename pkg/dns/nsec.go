// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

// NSECRecord is a minimal representation of the reference NSEC helper output.
//
// The JS reference helper returns a DNS record object with the name, next
// domain, type bitmap, and TTL already populated. The Go port keeps the data
// accessible without introducing a heavier DNS wire dependency.
type NSECRecord struct {
	Name       string
	NextDomain string
	TypeBitmap []byte
	TTL        int
}

// Create constructs a reference NSEC record container.
//
//	name := Create(".", NextName("."), TYPE_MAP_ROOT)
func Create(name, nextDomain string, typeBitmap []byte) NSECRecord {
	return NSECRecord{
		Name:       fqdn(name),
		NextDomain: fqdn(nextDomain),
		TypeBitmap: append([]byte(nil), typeBitmap...),
		TTL:        DEFAULT_TTL,
	}
}

// GetCreate is an alias for Create.
func GetCreate(name, nextDomain string, typeBitmap []byte) NSECRecord {
	return Create(name, nextDomain, typeBitmap)
}

// NextName returns the canonical successor for a top-level domain name.
//
// The helper lowercases and trims any trailing dot before applying the
// reference ordering logic.
func NextName(tld string) string {
	tld = trimFQDN(core.Lower(tld))

	if len(tld) == 63 {
		last := tld[62] + 1
		return tld[:62] + string([]byte{last}) + "."
	}

	return tld + "\000."
}

// GetNextName is an alias for NextName.
func GetNextName(tld string) string {
	return NextName(tld)
}

// PrevName returns the canonical predecessor for a top-level domain name.
//
// The helper lowercases and trims any trailing dot before applying the
// reference ordering logic.
func PrevName(tld string) string {
	tld = trimFQDN(core.Lower(tld))
	if len(tld) == 0 {
		return "."
	}

	last := tld[len(tld)-1] - 1
	tld = tld[:len(tld)-1] + string([]byte{last})

	if len(tld) < 63 {
		tld += "\xff"
	}

	return fqdn(tld)
}

// GetPrevName is an alias for PrevName.
func GetPrevName(tld string) string {
	return PrevName(tld)
}

func trimFQDN(name string) string {
	for core.HasSuffix(name, ".") {
		name = core.TrimSuffix(name, ".")
	}
	return name
}

func fqdn(name string) string {
	if name == "" {
		return "."
	}

	if core.HasSuffix(name, ".") {
		return name
	}

	return name + "."
}
