// SPDX-License-Identifier: EUPL-1.2

package primitives

import "testing"

func TestAddressHelpers(t *testing.T) {
	tests := []struct {
		name          string
		address       Address
		isNull        bool
		isUnspendable bool
		isPkh         bool
		isSh          bool
		isNulldata    bool
		isUnknown     bool
		isValid       bool
	}{
		{
			name:          "zero value",
			address:       Address{},
			isNull:        true,
			isUnspendable: false,
			isPkh:         false,
			isSh:          false,
			isNulldata:    false,
			isUnknown:     true,
			isValid:       false,
		},
		{
			name:          "pubkey hash",
			address:       Address{Version: 0, Hash: make([]byte, 20)},
			isNull:        true,
			isUnspendable: false,
			isPkh:         true,
			isSh:          false,
			isNulldata:    false,
			isUnknown:     false,
			isValid:       true,
		},
		{
			name:          "scripthash",
			address:       Address{Version: 0, Hash: make([]byte, 32)},
			isNull:        true,
			isUnspendable: false,
			isPkh:         false,
			isSh:          true,
			isNulldata:    false,
			isUnknown:     false,
			isValid:       true,
		},
		{
			name:          "nulldata",
			address:       Address{Version: 31, Hash: make([]byte, 2)},
			isNull:        true,
			isUnspendable: true,
			isPkh:         false,
			isSh:          false,
			isNulldata:    true,
			isUnknown:     false,
			isValid:       true,
		},
		{
			name:          "unknown programme",
			address:       Address{Version: 1, Hash: []byte{1, 2, 3}},
			isNull:        false,
			isUnspendable: false,
			isPkh:         false,
			isSh:          false,
			isNulldata:    false,
			isUnknown:     true,
			isValid:       true,
		},
	}

	for _, tc := range tests {
		if got := tc.address.IsNull(); got != tc.isNull {
			t.Fatalf("%s: IsNull() = %v, want %v", tc.name, got, tc.isNull)
		}

		if got := tc.address.IsUnspendable(); got != tc.isUnspendable {
			t.Fatalf("%s: IsUnspendable() = %v, want %v", tc.name, got, tc.isUnspendable)
		}

		if got := tc.address.IsPubkeyHash(); got != tc.isPkh {
			t.Fatalf("%s: IsPubkeyHash() = %v, want %v", tc.name, got, tc.isPkh)
		}

		if got := tc.address.IsScripthash(); got != tc.isSh {
			t.Fatalf("%s: IsScripthash() = %v, want %v", tc.name, got, tc.isSh)
		}

		if got := tc.address.IsNulldata(); got != tc.isNulldata {
			t.Fatalf("%s: IsNulldata() = %v, want %v", tc.name, got, tc.isNulldata)
		}

		if got := tc.address.IsUnknown(); got != tc.isUnknown {
			t.Fatalf("%s: IsUnknown() = %v, want %v", tc.name, got, tc.isUnknown)
		}

		if got := tc.address.IsValid(); got != tc.isValid {
			t.Fatalf("%s: IsValid() = %v, want %v", tc.name, got, tc.isValid)
		}
	}

	left := Address{Version: 0, Hash: []byte{1, 2, 3}}
	right := left.Clone()

	if !left.Equals(right) {
		t.Fatal("Clone() should produce an equal address")
	}

	if got := left.Compare(right); got != 0 {
		t.Fatalf("Compare() = %d, want 0", got)
	}

	if got := (Address{Version: 0, Hash: []byte{1, 2, 3}}).Compare(Address{Version: 1, Hash: []byte{1, 2, 3}}); got >= 0 {
		t.Fatalf("Compare() = %d, want negative", got)
	}

	if got := (Address{Version: 1, Hash: []byte{1, 2, 3}}).Compare(Address{Version: 0, Hash: []byte{1, 2, 3}}); got <= 0 {
		t.Fatalf("Compare() = %d, want positive", got)
	}

	var injected Address
	if got := injected.Inject(left); got != &injected {
		t.Fatal("Inject() should return the receiver")
	}

	if !injected.Equals(left) {
		t.Fatal("Inject() should copy the source address")
	}
}
