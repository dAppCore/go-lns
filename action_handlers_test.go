// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"net/netip"
	"testing"

	core "dappco.re/go"
)

func TestServiceOnStartupRegistersActionSlice(t *testing.T) {
	c := core.New()
	svc := NewService(c)

	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	for _, name := range []string{
		"lns.name.lookup",
		"lns.name.available",
		"lns.name.auction",
		"lns.name.register",
		"lns.dns.resolve",
		"lns.dns.resolve.txt",
		"lns.dns.resolve.all",
		"lns.dns.reverse",
		"lns.dns.health",
		"lns.dns.discover",
		"lns.vpn.discover",
		"lns.vpn.connect",
		"lns.vpn.status",
		"lns.vpn.disconnect",
	} {
		if !c.Action(name).Exists() {
			t.Fatalf("expected action %q to be registered", name)
		}
	}
}

func TestNameActionsExposeCatalogAndRuleState(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	lookupResult := c.Action("lns.name.lookup").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "reserved.lthn"},
	))
	if !lookupResult.OK {
		t.Fatalf("lns.name.lookup returned %#v, want OK", lookupResult)
	}

	lookup, ok := lookupResult.Value.(nameLookupResult)
	if !ok {
		t.Fatalf("lns.name.lookup returned %T, want nameLookupResult", lookupResult.Value)
	}

	if !lookup.Valid || !lookup.Reserved || lookup.Available {
		t.Fatalf("lookup = %#v, want valid reserved unavailable result", lookup)
	}

	availableResult := c.Action("lns.name.available").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "free-name.lthn"},
	))
	if !availableResult.OK {
		t.Fatalf("lns.name.available returned %#v, want OK", availableResult)
	}

	available, ok := availableResult.Value.(nameAvailabilityResult)
	if !ok {
		t.Fatalf("lns.name.available returned %T, want nameAvailabilityResult", availableResult.Value)
	}

	if !available.Valid || !available.Available || available.AvailabilityReason != "available" {
		t.Fatalf("available = %#v, want valid available result", available)
	}

	auctionResult := c.Action("lns.name.auction").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "nec.lthn"},
		core.Option{Key: "height", Value: 150},
		core.Option{Key: "rules", Value: NameRules{
			NoRollout:         true,
			ClaimPeriod:       100,
			AlexaLockupPeriod: 250,
		}},
	))
	if !auctionResult.OK {
		t.Fatalf("lns.name.auction returned %#v, want OK", auctionResult)
	}

	auction, ok := auctionResult.Value.(nameAuctionResult)
	if !ok {
		t.Fatalf("lns.name.auction returned %T, want nameAuctionResult", auctionResult.Value)
	}

	if !auction.Valid || !auction.HasRollout || !auction.LockedAtHeight || auction.Available {
		t.Fatalf("auction = %#v, want rolled-out locked unavailable result", auction)
	}
}

func TestDNSActionsProjectResourcePayloads(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	resource := NewResource()
	resource.Records = []ResourceRecord{
		GLUE4Record{NS: "charon.lthn.", Address: netip.MustParseAddr("192.0.2.10")},
		GLUE6Record{NS: "charon.lthn.", Address: netip.MustParseAddr("2001:db8::10")},
		TXTRecord{Entries: []string{"hello"}},
	}

	resolveResult := c.Action("lns.dns.resolve").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "resource", Value: resource},
	))
	if !resolveResult.OK {
		t.Fatalf("lns.dns.resolve returned %#v, want OK", resolveResult)
	}

	resolve, ok := resolveResult.Value.(dnsResolveResult)
	if !ok {
		t.Fatalf("lns.dns.resolve returned %T, want dnsResolveResult", resolveResult.Value)
	}

	if len(resolve.Records4) != 1 || len(resolve.Records6) != 1 || len(resolve.Addresses) != 2 {
		t.Fatalf("resolve = %#v, want one IPv4, one IPv6, and two addresses", resolve)
	}

	txtResult := c.Action("lns.dns.resolve.txt").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "resource", Value: resource},
	))
	if !txtResult.OK {
		t.Fatalf("lns.dns.resolve.txt returned %#v, want OK", txtResult)
	}

	txt, ok := txtResult.Value.(dnsTXTResult)
	if !ok {
		t.Fatalf("lns.dns.resolve.txt returned %T, want dnsTXTResult", txtResult.Value)
	}

	if len(txt.Records) != 1 || len(txt.Records[0].Entries) != 1 || txt.Records[0].Entries[0] != "hello" {
		t.Fatalf("txt = %#v, want one TXT record", txt)
	}

	allResult := c.Action("lns.dns.resolve.all").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "resource", Value: resource},
	))
	if !allResult.OK {
		t.Fatalf("lns.dns.resolve.all returned %#v, want OK", allResult)
	}

	all, ok := allResult.Value.(dnsAllResult)
	if !ok {
		t.Fatalf("lns.dns.resolve.all returned %T, want dnsAllResult", allResult.Value)
	}

	if len(all.Records) != 4 {
		t.Fatalf("all = %#v, want four projected records", all)
	}

	reverseResult := c.Action("lns.dns.reverse").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "resource", Value: resource},
		core.Option{Key: "ip", Value: "192.0.2.10"},
	))
	if !reverseResult.OK {
		t.Fatalf("lns.dns.reverse returned %#v, want OK", reverseResult)
	}

	reverse, ok := reverseResult.Value.(dnsReverseResult)
	if !ok {
		t.Fatalf("lns.dns.reverse returned %T, want dnsReverseResult", reverseResult.Value)
	}

	if len(reverse.Names) != 1 || reverse.Names[0] != "charon.lthn." {
		t.Fatalf("reverse = %#v, want charon.lthn.", reverse)
	}

	healthResult := c.Action("lns.dns.health").Run(context.Background(), core.NewOptions())
	if !healthResult.OK {
		t.Fatalf("lns.dns.health returned %#v, want OK", healthResult)
	}

	health, ok := healthResult.Value.(dnsHealthResult)
	if !ok {
		t.Fatalf("lns.dns.health returned %T, want dnsHealthResult", healthResult.Value)
	}

	if health.Live {
		t.Fatalf("health = %#v, want Live=false until the HSD bridge exists", health)
	}
}

func TestNameRegisterActionReturnsDraftPlan(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	registerResult := c.Action("lns.name.register").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "free-name.lthn"},
		core.Option{Key: "height", Value: 100},
		core.Option{Key: "bid_value", Value: 5000},
		core.Option{Key: "nonce", Value: "1111111111111111111111111111111111111111111111111111111111111111"},
		core.Option{Key: "rules", Value: NameRules{
			NoRollout:         true,
			ClaimPeriod:       100,
			AlexaLockupPeriod: 250,
		}},
	))
	if !registerResult.OK {
		t.Fatalf("lns.name.register returned %#v, want OK", registerResult)
	}

	plan, ok := registerResult.Value.(nameActionPlanResult)
	if !ok {
		t.Fatalf("lns.name.register returned %T, want nameActionPlanResult", registerResult.Value)
	}

	if plan.Action != "lns.name.register" || !plan.Available || len(plan.Steps) != 4 {
		t.Fatalf("register plan = %#v, want available four-step draft", plan)
	}

	if plan.Steps[0].CovenantType != TypeOpen || plan.Steps[1].CovenantType != TypeBid || plan.Steps[2].CovenantType != TypeReveal || plan.Steps[3].CovenantType != TypeRegister {
		t.Fatalf("register steps = %#v, want OPEN/BID/REVEAL/REGISTER", plan.Steps)
	}

	for i, step := range plan.Steps {
		if step.Draft == nil || len(step.Draft.Transaction.Outputs) != 1 || step.Draft.Fee == 0 {
			t.Fatalf("register step %d draft = %#v, want one-output draft transaction with placeholder fee", i, step.Draft)
		}
	}
}

func TestUnsupportedDNSDiscoverStillReturnsExplicitErrorsWithoutSource(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	discoverResult := c.Action("lns.dns.discover").Run(context.Background(), core.NewOptions())
	if discoverResult.OK {
		t.Fatal("lns.dns.discover should fail without a chain_url")
	}

	err, ok := discoverResult.Value.(error)
	if !ok {
		t.Fatalf("lns.dns.discover returned %T, want error", discoverResult.Value)
	}

	if got := core.Operation(err); got != "lns.dns.discover" {
		t.Fatalf("lns.dns.discover error operation = %q, want %q", got, "lns.dns.discover")
	}
}
