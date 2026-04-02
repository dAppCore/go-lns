// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"crypto/sha3"
	"testing"
)

func TestResolve(t *testing.T) {
	svc := NewService()

	got, err := svc.Resolve("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("Resolve returned %x, want %x", got, want)
	}
}

func TestResolveRejectsInvalidNames(t *testing.T) {
	svc := NewService()

	cases := []string{
		"",
		"foo.bar.lthn",
		"foo-",
		"test.lthn",
	}

	for _, name := range cases {
		if _, err := svc.Resolve(name); err == nil {
			t.Fatalf("Resolve(%q) should reject the name", name)
		}
	}
}
