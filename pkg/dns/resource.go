// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
)

const resourceEncodingVersion = 0

// ResourceRecord is a DNS resource payload item.
type ResourceRecord interface {
	Type() HSType
}

// Resource mirrors the JS DNS resource container used for covenant payloads.
type Resource struct {
	TTL     int
	Records []ResourceRecord
}

// DSRecord mirrors the JS DS payload entry.
type DSRecord struct {
	KeyTag     uint16
	Algorithm  uint8
	DigestType uint8
	Digest     []byte
}

// NSRecord mirrors the JS NS payload entry.
type NSRecord struct {
	NS string
}

// GLUE4Record mirrors the JS IPv4 glue payload entry.
type GLUE4Record struct {
	NS      string
	Address netip.Addr
}

// GLUE6Record mirrors the JS IPv6 glue payload entry.
type GLUE6Record struct {
	NS      string
	Address netip.Addr
}

// SYNTH4Record mirrors the JS synthesized IPv4 payload entry.
type SYNTH4Record struct {
	Address netip.Addr
}

// SYNTH6Record mirrors the JS synthesized IPv6 payload entry.
type SYNTH6Record struct {
	Address netip.Addr
}

// TXTRecord mirrors the JS TXT payload entry.
type TXTRecord struct {
	Entries []string
}

// NewResource constructs a resource with the reference default TTL.
func NewResource() *Resource {
	return &Resource{TTL: DEFAULT_TTL}
}

// GetNewResource is an alias for NewResource.
func GetNewResource() *Resource {
	return NewResource()
}

// Type returns the DNS record type.
func (DSRecord) Type() HSType { return HSTypeDS }

// Type returns the DNS record type.
func (NSRecord) Type() HSType { return HSTypeNS }

// Type returns the DNS record type.
func (GLUE4Record) Type() HSType { return HSTypeGLUE4 }

// Type returns the DNS record type.
func (GLUE6Record) Type() HSType { return HSTypeGLUE6 }

// Type returns the DNS record type.
func (SYNTH4Record) Type() HSType { return HSTypeSYNTH4 }

// Type returns the DNS record type.
func (SYNTH6Record) Type() HSType { return HSTypeSYNTH6 }

// Type returns the DNS record type.
func (TXTRecord) Type() HSType { return HSTypeTXT }

// NS returns the synthesized target host for a SYNTH4 record.
func (r SYNTH4Record) NS() string {
	if !r.Address.Is4() {
		return "._synth."
	}

	addr := r.Address.As4()
	return synthName(addr[:])
}

// NS returns the synthesized target host for a SYNTH6 record.
func (r SYNTH6Record) NS() string {
	if !r.Address.Is6() {
		return "._synth."
	}

	addr := r.Address.As16()
	return synthName(addr[:])
}

func synthName(ip []byte) string {
	enc := base32.HexEncoding.WithPadding(base32.NoPadding)
	return "_" + strings.ToLower(enc.EncodeToString(ip)) + "._synth."
}

// HasType reports whether the resource contains at least one record of type t.
func (r *Resource) HasType(t HSType) bool {
	if r == nil {
		return false
	}

	for _, record := range r.Records {
		if record != nil && record.Type() == t {
			return true
		}
	}

	return false
}

// GetHasType is an alias for HasType.
func (r *Resource) GetHasType(t HSType) bool {
	return r.HasType(t)
}

// HasNS reports whether the resource contains NS-like authority data.
func (r *Resource) HasNS() bool {
	if r == nil {
		return false
	}

	for _, record := range r.Records {
		if record == nil {
			continue
		}

		switch record.Type() {
		case HSTypeNS, HSTypeGLUE4, HSTypeGLUE6, HSTypeSYNTH4, HSTypeSYNTH6:
			return true
		}
	}

	return false
}

// GetHasNS is an alias for HasNS.
func (r *Resource) GetHasNS() bool {
	return r.HasNS()
}

// HasDS reports whether the resource contains DS records.
func (r *Resource) HasDS() bool {
	return r.HasType(HSTypeDS)
}

// GetHasDS is an alias for HasDS.
func (r *Resource) GetHasDS() bool {
	return r.HasDS()
}

// Encode serializes the resource payload in the reference binary format.
func (r *Resource) Encode() ([]byte, error) {
	if r == nil {
		return nil, core.E("dns.Resource.Encode", "resource is required", nil)
	}

	var buf bytes.Buffer
	buf.WriteByte(resourceEncodingVersion)

	for _, record := range r.Records {
		if record == nil {
			return nil, core.E("dns.Resource.Encode", "record is required", nil)
		}

		buf.WriteByte(byte(record.Type()))
		if err := encodeRecord(&buf, record); err != nil {
			return nil, core.E("dns.Resource.Encode", err.Error(), nil)
		}
	}

	if buf.Len() > covenant.MaxResourceSize {
		return nil, core.E("dns.Resource.Encode", "resource exceeds maximum size", nil)
	}

	return buf.Bytes(), nil
}

// GetEncode is an alias for Encode.
func (r *Resource) GetEncode() ([]byte, error) {
	return r.Encode()
}

// Decode populates the resource from the reference binary format.
func (r *Resource) Decode(raw []byte) error {
	if r == nil {
		return core.E("dns.Resource.Decode", "resource is required", nil)
	}

	if len(raw) == 0 {
		return core.E("dns.Resource.Decode", "resource payload is required", nil)
	}

	if len(raw) > covenant.MaxResourceSize {
		return core.E("dns.Resource.Decode", "resource exceeds maximum size", nil)
	}

	if raw[0] != resourceEncodingVersion {
		return core.E("dns.Resource.Decode", fmt.Sprintf("unknown serialization version: %d", raw[0]), nil)
	}

	r.TTL = DEFAULT_TTL
	r.Records = r.Records[:0]

	pos := 1
	for pos < len(raw) {
		t := HSType(raw[pos])
		pos++

		record, next, known, err := decodeRecord(raw, pos, t)
		if err != nil {
			return core.E("dns.Resource.Decode", err.Error(), nil)
		}
		if !known {
			break
		}

		r.Records = append(r.Records, record)
		pos = next
	}

	return nil
}

// GetDecode is an alias for Decode.
func (r *Resource) GetDecode(raw []byte) error {
	return r.Decode(raw)
}

// DecodeResource decodes a raw DNS resource payload.
func DecodeResource(raw []byte) (*Resource, error) {
	resource := NewResource()
	if err := resource.Decode(raw); err != nil {
		return nil, err
	}

	return resource, nil
}

// GetDecodeResource is an alias for DecodeResource.
func GetDecodeResource(raw []byte) (*Resource, error) {
	return DecodeResource(raw)
}

func encodeRecord(buf *bytes.Buffer, record ResourceRecord) error {
	switch rr := record.(type) {
	case DSRecord:
		return encodeDSRecord(buf, rr)
	case *DSRecord:
		return encodeDSRecord(buf, *rr)
	case NSRecord:
		return encodeNSRecord(buf, rr)
	case *NSRecord:
		return encodeNSRecord(buf, *rr)
	case GLUE4Record:
		return encodeGLUE4Record(buf, rr)
	case *GLUE4Record:
		return encodeGLUE4Record(buf, *rr)
	case GLUE6Record:
		return encodeGLUE6Record(buf, rr)
	case *GLUE6Record:
		return encodeGLUE6Record(buf, *rr)
	case SYNTH4Record:
		return encodeSYNTH4Record(buf, rr)
	case *SYNTH4Record:
		return encodeSYNTH4Record(buf, *rr)
	case SYNTH6Record:
		return encodeSYNTH6Record(buf, rr)
	case *SYNTH6Record:
		return encodeSYNTH6Record(buf, *rr)
	case TXTRecord:
		return encodeTXTRecord(buf, rr)
	case *TXTRecord:
		return encodeTXTRecord(buf, *rr)
	default:
		return errors.New("unknown record type")
	}
}

func encodeDSRecord(buf *bytes.Buffer, record DSRecord) error {
	if len(record.Digest) > 255 {
		return errors.New("ds digest exceeds maximum size")
	}

	if err := binary.Write(buf, binary.BigEndian, record.KeyTag); err != nil {
		return err
	}

	buf.WriteByte(record.Algorithm)
	buf.WriteByte(record.DigestType)
	buf.WriteByte(byte(len(record.Digest)))
	buf.Write(record.Digest)
	return nil
}

func encodeNSRecord(buf *bytes.Buffer, record NSRecord) error {
	return writeName(buf, record.NS)
}

func encodeGLUE4Record(buf *bytes.Buffer, record GLUE4Record) error {
	if !record.Address.Is4() {
		return errors.New("glue4 address must be ipv4")
	}

	if err := writeName(buf, record.NS); err != nil {
		return err
	}

	addr := record.Address.As4()
	buf.Write(addr[:])
	return nil
}

func encodeGLUE6Record(buf *bytes.Buffer, record GLUE6Record) error {
	if !record.Address.Is6() {
		return errors.New("glue6 address must be ipv6")
	}

	if err := writeName(buf, record.NS); err != nil {
		return err
	}

	addr := record.Address.As16()
	buf.Write(addr[:])
	return nil
}

func encodeSYNTH4Record(buf *bytes.Buffer, record SYNTH4Record) error {
	if !record.Address.Is4() {
		return errors.New("synth4 address must be ipv4")
	}

	addr := record.Address.As4()
	buf.Write(addr[:])
	return nil
}

func encodeSYNTH6Record(buf *bytes.Buffer, record SYNTH6Record) error {
	if !record.Address.Is6() {
		return errors.New("synth6 address must be ipv6")
	}

	addr := record.Address.As16()
	buf.Write(addr[:])
	return nil
}

func encodeTXTRecord(buf *bytes.Buffer, record TXTRecord) error {
	if len(record.Entries) > 255 {
		return errors.New("txt entry count exceeds maximum size")
	}

	buf.WriteByte(byte(len(record.Entries)))

	for _, entry := range record.Entries {
		if len(entry) > 255 {
			return errors.New("txt entry exceeds maximum size")
		}

		buf.WriteByte(byte(len(entry)))
		buf.WriteString(entry)
	}

	return nil
}

func decodeRecord(raw []byte, pos int, t HSType) (ResourceRecord, int, bool, error) {
	switch t {
	case HSTypeDS:
		return decodeDSRecord(raw, pos)
	case HSTypeNS:
		return decodeNSRecord(raw, pos)
	case HSTypeGLUE4:
		return decodeGLUE4Record(raw, pos)
	case HSTypeGLUE6:
		return decodeGLUE6Record(raw, pos)
	case HSTypeSYNTH4:
		return decodeSYNTH4Record(raw, pos)
	case HSTypeSYNTH6:
		return decodeSYNTH6Record(raw, pos)
	case HSTypeTXT:
		return decodeTXTRecord(raw, pos)
	default:
		return nil, pos, false, nil
	}
}

func decodeDSRecord(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	if len(raw)-pos < 5 {
		return nil, pos, true, errors.New("truncated ds record")
	}

	digestLen := int(raw[pos+4])
	if len(raw)-pos < 5+digestLen {
		return nil, pos, true, errors.New("truncated ds digest")
	}

	record := DSRecord{
		KeyTag:     binary.BigEndian.Uint16(raw[pos : pos+2]),
		Algorithm:  raw[pos+2],
		DigestType: raw[pos+3],
		Digest:     append([]byte(nil), raw[pos+5:pos+5+digestLen]...),
	}
	return record, pos + 5 + digestLen, true, nil
}

func decodeNSRecord(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	name, next, err := readName(raw, pos)
	if err != nil {
		return nil, pos, true, err
	}

	return NSRecord{NS: name}, next, true, nil
}

func decodeGLUE4Record(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	name, next, err := readName(raw, pos)
	if err != nil {
		return nil, pos, true, err
	}

	if len(raw)-next < 4 {
		return nil, pos, true, errors.New("truncated glue4 address")
	}

	addr, ok := netip.AddrFromSlice(raw[next : next+4])
	if !ok {
		return nil, pos, true, errors.New("invalid glue4 address")
	}

	return GLUE4Record{NS: name, Address: addr}, next + 4, true, nil
}

func decodeGLUE6Record(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	name, next, err := readName(raw, pos)
	if err != nil {
		return nil, pos, true, err
	}

	if len(raw)-next < 16 {
		return nil, pos, true, errors.New("truncated glue6 address")
	}

	addr, ok := netip.AddrFromSlice(raw[next : next+16])
	if !ok {
		return nil, pos, true, errors.New("invalid glue6 address")
	}

	return GLUE6Record{NS: name, Address: addr}, next + 16, true, nil
}

func decodeSYNTH4Record(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	if len(raw)-pos < 4 {
		return nil, pos, true, errors.New("truncated synth4 address")
	}

	addr, ok := netip.AddrFromSlice(raw[pos : pos+4])
	if !ok {
		return nil, pos, true, errors.New("invalid synth4 address")
	}

	return SYNTH4Record{Address: addr}, pos + 4, true, nil
}

func decodeSYNTH6Record(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	if len(raw)-pos < 16 {
		return nil, pos, true, errors.New("truncated synth6 address")
	}

	addr, ok := netip.AddrFromSlice(raw[pos : pos+16])
	if !ok {
		return nil, pos, true, errors.New("invalid synth6 address")
	}

	return SYNTH6Record{Address: addr}, pos + 16, true, nil
}

func decodeTXTRecord(raw []byte, pos int) (ResourceRecord, int, bool, error) {
	if len(raw)-pos < 1 {
		return nil, pos, true, errors.New("truncated txt record")
	}

	count := int(raw[pos])
	pos++
	record := TXTRecord{Entries: make([]string, 0, count)}

	for i := 0; i < count; i++ {
		if len(raw)-pos < 1 {
			return nil, pos, true, errors.New("truncated txt entry")
		}

		size := int(raw[pos])
		pos++
		if len(raw)-pos < size {
			return nil, pos, true, errors.New("truncated txt value")
		}

		record.Entries = append(record.Entries, string(raw[pos:pos+size]))
		pos += size
	}

	return record, pos, true, nil
}

func writeName(buf *bytes.Buffer, name string) error {
	name = fqdn(name)
	if name == "." {
		buf.WriteByte(0)
		return nil
	}

	trimmed := trimFQDN(name)
	if trimmed == "" {
		buf.WriteByte(0)
		return nil
	}

	labels := strings.Split(trimmed, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > covenant.MaxNameSize {
			return fmt.Errorf("invalid dns label %q", label)
		}

		buf.WriteByte(byte(len(label)))
		buf.WriteString(label)
	}

	buf.WriteByte(0)
	return nil
}

func readName(raw []byte, pos int) (string, int, error) {
	labels, next, err := readNameLabels(raw, pos, 0, map[int]struct{}{})
	if err != nil {
		return "", pos, err
	}

	if len(labels) == 0 {
		return ".", next, nil
	}

	return strings.Join(labels, ".") + ".", next, nil
}

func readNameLabels(raw []byte, pos, depth int, seen map[int]struct{}) ([]string, int, error) {
	if depth > 16 {
		return nil, pos, errors.New("dns name pointer depth exceeded")
	}

	if pos >= len(raw) {
		return nil, pos, errors.New("truncated dns name")
	}

	if _, ok := seen[pos]; ok {
		return nil, pos, errors.New("cyclic dns name pointer")
	}
	seen[pos] = struct{}{}

	labels := make([]string, 0, 4)
	origPos := pos

	for {
		if pos >= len(raw) {
			return nil, origPos, errors.New("truncated dns name")
		}

		size := raw[pos]
		pos++

		switch size & 0xc0 {
		case 0xc0:
			if pos >= len(raw) {
				return nil, origPos, errors.New("truncated dns name pointer")
			}

			offset := int(size&0x3f)<<8 | int(raw[pos])
			pos++

			ptrLabels, _, err := readNameLabels(raw, offset, depth+1, seen)
			if err != nil {
				return nil, origPos, err
			}

			labels = append(labels, ptrLabels...)
			return labels, pos, nil
		case 0x00:
			if size == 0 {
				return labels, pos, nil
			}
		default:
			return nil, origPos, errors.New("unsupported dns name label encoding")
		}

		if size > covenant.MaxNameSize {
			return nil, origPos, errors.New("dns label exceeds maximum size")
		}

		if len(raw)-pos < int(size) {
			return nil, origPos, errors.New("truncated dns label")
		}

		labels = append(labels, string(raw[pos:pos+int(size)]))
		pos += int(size)
	}
}
