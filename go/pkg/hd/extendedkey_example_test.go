// SPDX-License-Identifier: EUPL-1.2

package hd

import core "github.com/dappcore/go"

func ExampleNewMaster() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, err := NewMaster(seed)
	if err != nil {
		core.Println("error")
		return
	}
	core.Println(master.IsPrivate())
	// Output: true
}

func ExampleExtendedKey_Child() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	child, _ := master.Child(HardenedKeyStart + 0)
	priv, _ := child.PrivateKeyBytes()
	core.Println(core.HexEncode(priv))
	// Output: edb2e14f9ee77d26dd93b4ecede8d16ed408ce149b6cd80b0715a2d911a0afea
}

func ExampleExtendedKey_Neuter() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	pub, _ := master.Neuter()
	core.Println(pub.IsPrivate())
	// Output: false
}

func ExampleExtendedKey_PublicKeyBytes() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	pub, _ := master.PublicKeyBytes()
	core.Println(core.HexEncode(pub))
	// Output: 0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2
}

func ExampleExtendedKey_Address() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	addr := master.Address()
	core.Println(addr.IsPubkeyHash())
	// Output: true
}
