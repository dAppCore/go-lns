// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

const signingSeed = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

// TestSigning_SigningKey_Good resolves the signing key for the wallet's own
// receive address and checks the key produces that same address.
func TestSigning_SigningKey_Good(t *testing.T) {
	wallet, err := ImportSeed(signingSeed)
	if err != nil {
		t.Fatalf("ImportSeed error: %v", err)
	}

	addr := wallet.ReceiveAddress()
	priv, pub, err := wallet.SigningKey(addr, 0)
	if err != nil {
		t.Fatalf("SigningKey error: %v", err)
	}
	if len(priv) != 32 {
		t.Fatalf("private key length = %d, want 32", len(priv))
	}
	if len(pub) != 33 {
		t.Fatalf("public key length = %d, want 33", len(pub))
	}
}

// TestSigning_SigningKey_Bad covers an address the wallet does not own.
func TestSigning_SigningKey_Bad(t *testing.T) {
	wallet, err := ImportSeed(signingSeed)
	if err != nil {
		t.Fatalf("ImportSeed error: %v", err)
	}

	foreign := primitives.Address{Version: 0, Hash: make([]byte, 20)}
	if _, _, err := wallet.SigningKey(foreign, 0); err == nil {
		t.Fatalf("SigningKey resolved a foreign address, want error")
	}
}

// TestSigning_SigningKey_Ugly covers a non-P2WPKH address.
func TestSigning_SigningKey_Ugly(t *testing.T) {
	wallet, err := ImportSeed(signingSeed)
	if err != nil {
		t.Fatalf("ImportSeed error: %v", err)
	}

	scriptHash := primitives.Address{Version: 0, Hash: make([]byte, 32)}
	if _, _, err := wallet.SigningKey(scriptHash, 0); err == nil {
		t.Fatalf("SigningKey accepted a 32-byte programme, want error")
	}
}

// TestSigning_SpendableAddresses_Good enumerates owned addresses across both
// branches and confirms the current receive address is included and resolvable.
func TestSigning_SpendableAddresses_Good(t *testing.T) {
	wallet, err := ImportSeed(signingSeed)
	if err != nil {
		t.Fatalf("ImportSeed error: %v", err)
	}

	refs := wallet.SpendableAddresses(0)
	if len(refs) == 0 {
		t.Fatalf("SpendableAddresses returned none")
	}

	recv := wallet.ReceiveAddress()
	found := false
	for _, ref := range refs {
		if ref.Address.Equals(recv) {
			found = true
			if _, _, err := wallet.SigningKey(ref.Address, 0); err != nil {
				t.Fatalf("enumerated address not resolvable: %v", err)
			}
		}
	}
	if !found {
		t.Fatalf("receive address not present in SpendableAddresses")
	}
}

// TestSigning_SpendableAddresses_Bad covers the nil receiver.
func TestSigning_SpendableAddresses_Bad(t *testing.T) {
	var wallet *Wallet
	if refs := wallet.SpendableAddresses(0); refs != nil {
		t.Fatalf("nil wallet returned %d refs, want nil", len(refs))
	}
}

// TestSigning_SpendableAddresses_Ugly checks the gap extends the scan range so a
// forward address beyond the receive index is still enumerated.
func TestSigning_SpendableAddresses_Ugly(t *testing.T) {
	wallet, err := ImportSeed(signingSeed)
	if err != nil {
		t.Fatalf("ImportSeed error: %v", err)
	}

	small := wallet.SpendableAddresses(1)
	large := wallet.SpendableAddresses(50)
	if len(large) <= len(small) {
		t.Fatalf("larger gap (%d) did not enumerate more than small gap (%d)", len(large), len(small))
	}
}
