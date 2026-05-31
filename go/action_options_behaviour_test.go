// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"time"

	core "dappco.re/go"
	dnspkg "dappco.re/go/lns/pkg/dns"
)

func TestActionOptionsBehaviour_chainClientFromOptions_Good(t *core.T) {
	opts := core.NewOptions()
	opts.Set("chain_url", "http://127.0.0.1:14037")
	opts.Set("api_key", "secret")
	client, err := chainClientFromOptions(opts)
	core.AssertNoError(t, err)
	core.AssertNotNil(t, client)
}

func TestActionOptionsBehaviour_chainClientFromOptions_Bad(t *core.T) {
	// Missing chain_url is rejected.
	client, err := chainClientFromOptions(core.NewOptions())
	core.AssertError(t, err)
	core.AssertNil(t, client)
}

func TestActionOptionsBehaviour_bridgeConfigFromOptions(t *core.T) {
	// No chain URL -> not configured.
	_, ok := bridgeConfigFromOptions(core.NewOptions())
	core.AssertFalse(t, ok)

	opts := core.NewOptions()
	opts.Set("chain_url", "http://127.0.0.1:14037")
	opts.Set("alias_url", "http://127.0.0.1:14039")
	opts.Set("dns_port", "5355")
	opts.Set("known_names", []string{"charon", "trade"})
	cfg, ok := bridgeConfigFromOptions(opts)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "http://127.0.0.1:14037", cfg.ChainURL)
	core.AssertEqual(t, ":5355", cfg.DNSAddr)
	core.AssertLen(t, cfg.KnownNames, 2)
}

func TestActionOptionsBehaviour_dnsListenAddr(t *core.T) {
	// Explicit address wins.
	opts := core.NewOptions()
	opts.Set("dns_addr", "0.0.0.0:53")
	core.AssertEqual(t, "0.0.0.0:53", dnsListenAddr(opts))

	// Bare port is prefixed with a colon.
	port := core.NewOptions()
	port.Set("dns_port", "5300")
	core.AssertEqual(t, ":5300", dnsListenAddr(port))

	// Already-colonised port is left as-is.
	colon := core.NewOptions()
	colon.Set("dns_port", ":5301")
	core.AssertEqual(t, ":5301", dnsListenAddr(colon))

	// Default when nothing is set.
	core.AssertEqual(t, ":5354", dnsListenAddr(core.NewOptions()))
}

func TestActionOptionsBehaviour_bridgeCheckInterval(t *core.T) {
	// Numeric seconds option.
	secs := core.NewOptions()
	secs.Set("check_interval", uint32(30))
	core.AssertEqual(t, 30*time.Second, bridgeCheckInterval(secs))

	// Duration string.
	dur := core.NewOptions()
	dur.Set("check_interval", "2m")
	core.AssertEqual(t, 2*time.Minute, bridgeCheckInterval(dur))

	// Default when unset.
	core.AssertEqual(t, 15*time.Second, bridgeCheckInterval(core.NewOptions()))
}

func TestActionOptionsBehaviour_actionStringSliceValue(t *core.T) {
	core.AssertEqual(t, []string{"a", "b"}, actionStringSliceValue([]string{"a", "b"}))
	core.AssertEqual(t, []string{"a", "b"}, actionStringSliceValue([]any{"a", " b ", ""}))
	core.AssertEqual(t, []string{"a", "b"}, actionStringSliceValue("a, b ,"))
	core.AssertNil(t, actionStringSliceValue(""))
	core.AssertNil(t, actionStringSliceValue(42))
}

func TestActionOptionsBehaviour_bridgeConfigsEqual(t *core.T) {
	left := dnspkg.BridgeConfig{ChainURL: "u", KnownNames: []string{"a"}}
	right := dnspkg.BridgeConfig{ChainURL: "u", KnownNames: []string{"a"}}
	core.AssertTrue(t, bridgeConfigsEqual(left, right))

	right.KnownNames = []string{"b"}
	core.AssertFalse(t, bridgeConfigsEqual(left, right))
}

func TestActionOptionsBehaviour_decodeRawJSON(t *core.T) {
	core.AssertNil(t, decodeRawJSON(nil))

	decoded := decodeRawJSON([]byte(`{"a":1}`))
	asMap, ok := decoded.(map[string]any)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, float64(1), asMap["a"])

	// Non-JSON falls back to the raw string.
	core.AssertEqual(t, "not json", decodeRawJSON([]byte("not json")))
}
