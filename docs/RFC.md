---
module: github.com/dappcore/go-lns/go
repo: core/go-lns
lang: go
tier: lib
depends:
  - code/core/go
tags:
  - lns
  - naming
  - resolution
  - dns
---
# go-lns — Lethean Name System as CoreGO Native

> Convert `lthn/itns-sidechain` (118k lines JS/HSD fork) into native CoreGO.
> Every subsystem becomes a CoreApp service with actions → automatic MCP/API/CLI.
> Any CoreApp includes LNS with a build tag. ISP-in-a-box level.

**Module:** `github.com/dappcore/go-lns/go`
**Repository:** `core/go-lns`
**Source:** `ssh://git@forge.lthn.ai:2223/lthn/itns-sidechain.git` (HSD 8.99.0 fork)
**Existing Go:** `~/Code/lthn/lns/` (672 lines, DNS resolution working on testnet)

---

## 1. Overview

The ITNS sidechain is a Handshake (HSD) fork that powers .lthn name registration, DNS resolution, and VPN gateway discovery. Currently 118k lines of JavaScript. This RFC converts it into native CoreGO packages that integrate as CoreApp services.

The conversion is NOT a line-by-line port. It's a semantic extraction: take what each JS module DOES and implement it using Core primitives.

---

## 2. Source Subsystems (JS → Go)

| JS Module | Lines | CoreGO Package | Actions Prefix | Priority |
|-----------|-------|---------------|----------------|----------|
| `lib/dns/` | ~1,200 | `pkg/dns` | `lns.dns.*` | P1 — already 672 lines Go |
| `lib/covenants/` | ~8,000 | `pkg/covenant` | `lns.name.*` | P1 — name rules |
| `lib/wallet/` | ~12,000 | `pkg/wallet` | `lns.wallet.*` | P1 — key management |
| `lib/blockchain/` | ~6,000 | `pkg/chain` | `lns.chain.*` | P1 — sync + validation |
| `lib/blockstore/` | ~2,000 | `pkg/blockstore` | `lns.store.*` | P2 — via go-store |
| `lib/net/` | ~4,000 | `pkg/net` | `lns.peer.*` | P2 — via go-p2p |
| `lib/primitives/` | ~5,000 | `pkg/primitives` | — (types only) | P1 — shared types |
| `lib/mining/` | ~3,000 | `pkg/mining` | `lns.mine.*` | P3 — block production |
| `lib/hd/` | ~2,000 | `pkg/hd` | — (key derivation) | P2 |
| `lib/script/` | ~4,000 | `pkg/script` | — (tx scripting) | P2 |
| `lib/mempool/` | ~3,000 | `pkg/mempool` | `lns.mempool.*` | P3 |
| `lib/ui/` | ~500 | `ui/` | — (CoreGUI elements) | P2 |

---

## 3. Core Actions

### 3.1 Name Management (from covenants/)

```go
c.Action("lns.name.register", ...)    // OPEN + BID + REVEAL + REGISTER
c.Action("lns.name.renew", ...)       // Renew before expiry
c.Action("lns.name.transfer", ...)    // Transfer to new owner
c.Action("lns.name.update", ...)      // Update DNS records
c.Action("lns.name.revoke", ...)      // Revoke name
c.Action("lns.name.lookup", ...)      // Check name state
c.Action("lns.name.auction", ...)     // View auction state
c.Action("lns.name.available", ...)   // Check availability
```

### 3.2 DNS Resolution (from dns/ + existing LNS)

```go
c.Action("lns.dns.resolve", ...)      // A record
c.Action("lns.dns.resolve.txt", ...)  // TXT record
c.Action("lns.dns.resolve.all", ...)  // All record types
c.Action("lns.dns.reverse", ...)      // PTR
c.Action("lns.dns.serve", ...)        // Start DNS server
c.Action("lns.dns.health", ...)       // Health check
c.Action("lns.dns.discover", ...)     // Rescan chain
```

### 3.3 Wallet (from wallet/)

```go
c.Action("lns.wallet.create", ...)    // Create HD wallet
c.Action("lns.wallet.balance", ...)   // Get balance
c.Action("lns.wallet.send", ...)      // Send HNS
c.Action("lns.wallet.receive", ...)   // Generate address
c.Action("lns.wallet.history", ...)   // Transaction history
c.Action("lns.wallet.import", ...)    // Import seed/key
c.Action("lns.wallet.export", ...)    // Export seed
```

### 3.4 Chain (from blockchain/)

```go
c.Action("lns.chain.height", ...)     // Current height
c.Action("lns.chain.block", ...)      // Get block by height/hash
c.Action("lns.chain.sync", ...)       // Start sync
c.Action("lns.chain.status", ...)     // Sync status
c.Action("lns.chain.peers", ...)      // Connected peers
```

### 3.5 VPN Gateway (from dns/ TXT records)

```go
c.Action("lns.vpn.discover", ...)     // Find VPN gateways via TXT records
c.Action("lns.vpn.connect", ...)      // Connect to gateway
c.Action("lns.vpn.status", ...)       // Connection status
c.Action("lns.vpn.disconnect", ...)   // Disconnect
```

---

## 4. CoreGUI Elements

Custom HTML elements for any CoreApp:

```html
<!-- Wallet widget -->
<core-lns-wallet></core-lns-wallet>

<!-- Name registration form -->
<core-lns-register></core-lns-register>

<!-- Name lookup -->
<core-lns-lookup name="charon.lthn"></core-lns-lookup>

<!-- VPN gateway selector -->
<core-lns-vpn></core-lns-vpn>

<!-- DNS status -->
<core-lns-dns-status></core-lns-dns-status>
```

These are CoreTS components rendered in CoreGUI's WebView (Wails). They call the Go actions via the WebView bridge.

---

## 5. Conversion Strategy

### 5.1 NOT Line-by-Line

The JS code has Node.js-isms (EventEmitter, streams, callbacks) that don't translate. Instead:

1. **Read the JS module** — understand what it DOES
2. **Identify the data types** — extract structs from JS objects
3. **Identify the operations** — extract functions as Core actions
4. **Implement in Go** — using Core primitives (c.Fs, c.Process, core.E, go-store, go-p2p)
5. **Test against testnet** — the C++ seed node validates correctness

### 5.2 What Already Exists

| Component | Status | Location |
|-----------|--------|----------|
| DNS resolution | Working (672 lines Go) | `~/Code/lthn/lns/` |
| P2P networking | go-p2p exists | `github.com/dappcore/go-p2p/go` |
| Block storage | go-store exists | `github.com/dappcore/go-store/go` |
| Chain sync | go-blockchain has it | Charon's repo |
| Wallet | Charon implementing today | go-blockchain |
| Covenants | JS only | Needs conversion |
| Script engine | JS only | Needs conversion |
| HD key derivation | JS only | Needs conversion |

### 5.3 Dispatch Prompt

```
Read the JS source at ~/ref/itns-sidechain/lib/{module}/.
Read the Go code at ~/repo/pkg/{module}/.
Find ONE feature implemented in JS that is NOT in the Go code.
Implement it using Core primitives. Commit when done.
```

The JS source is the spec. The Go code is the implementation. Same "find ONE missing feature" prompt, different source material.

---

## 6. CoreApp Integration

Any CoreApp includes LNS with a build tag:

```go
//go:build lns

import _ "github.com/dappcore/go-lns/go"
```

This auto-registers all `lns.*` actions. The app gets:
- Name registration/lookup
- DNS resolution
- Wallet management
- VPN gateway discovery
- Chain sync

Without the build tag: no LNS code in the binary. AX-1 (predictable names), §7.1 (build variants).

---

## 7. CorePlay Integration

For the Stronghold Kingdoms-inspired game:

```go
// Team messaging via chain
c.Action("lns.message.send", ...)     // Send encrypted message via TX memo
c.Action("lns.message.inbox", ...)    // Read received messages

// Notification settings per player
c.Action("lns.notify.set", ...)       // "wake me up if attacked"
c.Action("lns.notify.level", ...)     // urgent/normal/silent

// Team discovery
c.Action("lns.team.register", ...)    // Register team name (e.g. "squirrels.lthn")
c.Action("lns.team.members", ...)     // List team members
c.Action("lns.team.invite", ...)      // Invite via name
```

The sleeping channel from SHK W3 — configurable notification levels backed by chain state. Your team name IS your .lthn name. Messages are encrypted TX memos. The chain is the message bus.

---

## 8. Implementation Priority

1. `pkg/primitives` — shared types (Block, TX, Address, Name, Covenant)
2. `pkg/covenant` — name rules (OPEN, BID, REVEAL, REGISTER, UPDATE, RENEW, TRANSFER, REVOKE)
3. `pkg/dns` — merge existing 672-line LNS + extend
4. `pkg/wallet` — HD wallet, address gen, balance, send/receive
5. `pkg/chain` — sync, validation, block storage via go-store
6. `pkg/script` — transaction scripting
7. `pkg/hd` — HD key derivation (BIP32/44)
8. `pkg/net` — P2P via go-p2p
9. `pkg/blockstore` — raw block storage via go-io
10. `pkg/mining` — block production
11. `pkg/mempool` — transaction pool
12. `ui/` — CoreGUI custom elements

---

## 9. Testing Strategy

The C++ seed node on testnet IS the test harness:
- Go peer syncs from C++ seed → blocks match = correct
- Go wallet sends TX → C++ seed validates = correct
- Go name registration → C++ seed accepts = correct
- No unit test mocking needed — the real chain validates everything

---

## Changelog

- 2026-04-02: Initial RFC — convert ITNS sidechain (118k JS) to native CoreGO. ISP-in-a-box level CoreApp integration. Charon's real wishlist item.
# go-lns — Lethean Name Server (Sidechain DNS Bridge)

> The LNS bridges the ITNS sidechain (HSD fork) to standard DNS.
> It reads name records from the sidechain and serves them over port 53.
> Subject matter expert: Charon. Proven on testnet with 14 names.

**Module:** `github.com/dappcore/go-lns/go`
**Depends on:** `github.com/dappcore/go-dns/go` (record DTOs), `github.com/dappcore/go-blockchain/go` (alias discovery)
**Current implementation:** `~/Code/lthn/lns/` (672 lines, standalone)

## 1. Architecture

```
Main chain (aliases)
  │ c.Action("blockchain.chain.aliases")
  │ → 14 aliases with hns= comments
  │
  ▼
HSD sidechain (DNS records)
  │ hsd.GetNameResource("charon")
  │ → GLUE4, TXT, NS, DS records
  │
  ▼
Zone Cache (in-memory)
  │ tree-root invalidation (15s check)
  │ hash: getblockchaininfo → treeroot
  │
  ▼
DNS Server (miekg/dns)
  │ UDP + TCP, port 53
  │ A, AAAA, TXT, NS, SOA, PTR
  │
  ▼
Standard DNS clients
  dig @10.69.69.165 charon.lthn A
```

## 2. HSD RPC Client

The sidechain is an HSD fork (Handshake). Different RPC from the main chain:

| Method | Purpose | Port |
|--------|---------|------|
| `getblockchaininfo` | Tree root hash for cache invalidation | 14037 |
| `getnameresource` | DNS records for a name | 14037 |
| `getinfo` | Sidechain height | 14037 |

Auth: HTTP Basic (`x:{api_key}`). NOT the same as main chain RPC.

```go
type HSDClient struct {
    url    string
    apiKey string
}

func (c *HSDClient) GetBlockchainInfo() (*HSDBlockchainInfo, error)
func (c *HSDClient) GetNameResource(name string) (*HSDNameResource, error)
func (c *HSDClient) GetInfo() (*HSDInfo, error)
```

**This is the go-hsd client that should be a shared package.** The LNS, zone-cache daemon, and Bob Wallet all need it.

### 2.1 HSD Record Types

HSD returns records in its own format. The LNS translates to standard DNS:

| HSD Type | DNS Type | Translation |
|----------|----------|-------------|
| `GLUE4` | `A` | `{ns, address}` → A record |
| `GLUE6` | `AAAA` | `{ns, address}` → AAAA record |
| `TXT` | `TXT` | `{txt: [...]}` → TXT record |
| `NS` | `NS` | `{ns}` → NS delegation |
| `DS` | `DS` | `{keyTag, algorithm, digestType, digest}` → DNSSEC |

```go
// HSD record as returned by getnameresource RPC
type HSDRecord struct {
    Type    string   `json:"type"`    // GLUE4, GLUE6, TXT, NS, DS
    NS      string   `json:"ns"`      // nameserver (GLUE4/6, NS)
    Address string   `json:"address"` // IP address (GLUE4/6)
    TXT     []string `json:"txt"`     // text records
    KeyTag  uint16   `json:"keyTag"`  // DS
    Algo    uint8    `json:"algorithm"`
    DigestT uint8    `json:"digestType"`
    Digest  string   `json:"digest"`
}
```

## 3. Discovery Pipeline

### 3.1 Chain-Based Discovery (Primary)

```go
// Step 1: Get aliases from main chain
aliases := c.Action("blockchain.chain.aliases").Run(ctx, opts)

// Step 2: For each alias, derive HNS name
for _, alias := range aliases {
    hnsName := alias.Name
    if comment := parseComment(alias.Comment); comment["hns"] != "" {
        hnsName = trimSuffix(comment["hns"], ".lthn")
    }
    
    // Step 3: Query HSD for DNS records
    resource := hsdClient.GetNameResource(hnsName)
    
    // Step 4: Translate HSD records → DNS cache
    cache.Set(hnsName, translateRecords(resource))
}
```

### 3.2 Tree-Root Invalidation

```go
// Every 15 seconds:
info := hsdClient.GetBlockchainInfo()
if info.TreeRoot != lastTreeRoot {
    // Sidechain state changed — regenerate zone
    rediscover()
    lastTreeRoot = info.TreeRoot
}
```

Only regenerate when the sidechain tree root changes. No polling individual names.

### 3.3 Alias Comment Format

```
v=lthn1;type=gateway;cap=vpn,dns,proxy;hns=charon.lthn
```

| Field | Purpose |
|-------|---------|
| `v` | Protocol version (lthn1) |
| `type` | Node type (gateway, service, root) |
| `cap` | Capabilities (comma-separated) |
| `hns` | Explicit HNS name override (default: alias name) |

## 4. DNS Server

Uses `github.com/miekg/dns` (battle-tested, full RFC compliance).

### 4.1 Supported Record Types

| Type | Handler | Source |
|------|---------|--------|
| A | Forward lookup | HSD GLUE4 |
| AAAA | Forward lookup | HSD GLUE6 |
| TXT | Forward lookup | HSD TXT |
| NS | Delegation | HSD NS |
| SOA | Zone apex | Generated |
| PTR | Reverse lookup | Built from A records |

### 4.2 Response Behaviour

| Query | Response |
|-------|----------|
| Known name, matching type | Answer section with records |
| Known name, no records of type | NOERROR, SOA in authority |
| Zone apex (lthn.) | SOA record |
| Unknown name (.lthn) | NXDOMAIN, SOA in authority |
| Non-.lthn query | REFUSED |

### 4.3 Reverse DNS Index

Built atomically when cache updates:
```go
type ReverseIndex map[string][]string // IP → []names
```

`dig -x 10.69.69.165` → returns all .lthn names pointing to that IP.

## 5. Core Service Registration

```go
lnsSvc := lns.NewService(c, lns.Options{
    HSDUrl:        "http://127.0.0.1:14037",
    HSDApiKey:     "testkey",
    DNSPort:       53,
    HTTPPort:      5554,
    CheckInterval: 15 * time.Second,
})

// Auto-registers actions:
// dns.resolve       — A record lookup
// dns.resolve.txt   — TXT record lookup
// dns.resolve.all   — all record types
// dns.reverse       — PTR lookup
// dns.health        — cache status
// dns.discover      — force re-scan
// dns.serve         — start DNS server
```

## 6. go-dns Integration

go-dns provides the record DTOs that go-lns populates:

```go
import "github.com/dappcore/go-dns/go"

// go-dns provides:
type ARecord struct { Address net.IP; TTL uint32 }
type AAAARecord struct { Address net.IP; TTL uint32 }
type TXTRecord struct { Text []string; TTL uint32 }
type NSRecord struct { Nameserver string; TTL uint32 }
type SOARecord struct { NS, Mbox string; Serial, Refresh, Retry, Expire, MinTTL uint32 }
type PTRRecord struct { Ptr string; TTL uint32 }

// go-lns populates these from HSD data and serves via miekg/dns
```

## 7. Entitlement Gating (Future)

```go
// Rate limit DNS queries per source IP
e := c.Entitled("dns.resolve")
if !e.Allowed || e.NearLimit(0.9) {
    // Return REFUSED or SERVFAIL
}
c.RecordUsage("dns.resolve", 1)
```

For hosted LNS nodes serving public DNS. Unlimited for self-hosted.

## 8. Proven (Testnet)

| Feature | Status |
|---------|--------|
| Chain discovery | 14 aliases → 14 HNS names |
| HSD record fetch | GLUE4, TXT for all 14 |
| Tree-root invalidation | 15s check, atomic update |
| DNS over UDP+TCP | A, AAAA, TXT, NS, SOA, PTR |
| Zone apex SOA | `dig lthn SOA` works |
| Reverse DNS | `dig -x 10.69.69.165` → 14 names |
| `hns=` override | `@trading` → `trade.lthn` |
| systemd service | `lns.service`, port 53, boot-enabled |
| Docker | 13MB image, `--network host` |

672 lines of Go. 14/14 names. Standard DNS on port 53.
