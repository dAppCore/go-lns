// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"encoding/hex"
	"net/netip"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
)

func resourceFromChainResource(resource *chainpkg.NameResource) (*Resource, error) {
	result := NewResource()
	if resource == nil {
		return result, nil
	}

	for _, record := range resource.Records {
		switch core.Upper(record.Type) {
		case "GLUE4":
			address, err := netip.ParseAddr(record.Address)
			if err != nil {
				return nil, core.E(opResourceFromChain, "parse GLUE4 address", err)
			}
			result.Records = append(result.Records, GLUE4Record{NS: dnsFQDN(record.NS), Address: address})
		case "GLUE6":
			address, err := netip.ParseAddr(record.Address)
			if err != nil {
				return nil, core.E(opResourceFromChain, "parse GLUE6 address", err)
			}
			result.Records = append(result.Records, GLUE6Record{NS: dnsFQDN(record.NS), Address: address})
		case "TXT":
			result.Records = append(result.Records, TXTRecord{Entries: append([]string(nil), record.TXT...)})
		case "NS":
			result.Records = append(result.Records, NSRecord{NS: dnsFQDN(record.NS)})
		case "DS":
			digest, err := hex.DecodeString(record.Digest)
			if err != nil {
				return nil, core.E(opResourceFromChain, "decode DS digest", err)
			}
			result.Records = append(result.Records, DSRecord{
				KeyTag:     record.KeyTag,
				Algorithm:  record.Algorithm,
				DigestType: record.DigestType,
				Digest:     digest,
			})
		}
	}

	return result, nil
}
