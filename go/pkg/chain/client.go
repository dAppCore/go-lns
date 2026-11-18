// SPDX-License-Identifier: EUPL-1.2

// Package chain provides lightweight RPC helpers for the LNS sidechain and
// the main-chain alias discovery path used by the DNS bridge.
package chain

import (
	"io"
	"net/http"
	"time"

	core "github.com/dappcore/go"
)

const (
	opRawCall            = "chain.Client.RawCall"
	opSendRawTransaction = "chain.Client.SendRawTransaction"
	opGetCoinsByAddress  = "chain.Client.GetCoinsByAddress"
	opGetNameOwner       = "chain.Client.GetNameOwner"
	methodSendRawTx      = "sendrawtransaction"
	methodGetCoinsByAddr = "getcoinsbyaddress"
	methodGetNameInfo    = "getnameinfo"
)

// Client is a JSON-RPC client for the LNS sidechain daemon.
type Client struct {
	url        string
	apiKey     string
	httpClient *http.Client
}

// ClientOption configures a sidechain RPC client.
type ClientOption func(*Client)

// WithAPIKey configures HTTP Basic auth for the sidechain RPC.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		if c == nil {
			return
		}
		c.apiKey = core.Trim(apiKey)
	}
}

// WithHTTPClient sets the HTTP transport used by the client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if c == nil || httpClient == nil {
			return
		}
		c.httpClient = httpClient
	}
}

// NewClient constructs a sidechain RPC client.
func NewClient(url string, opts ...ClientOption) *Client {
	client := &Client{
		url:        core.Trim(url),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	return client
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type rpcResponse struct {
	Result jsonraw     `json:"result"`
	Error  *rpcfailure `json:"error,omitempty"`
}

type rpcfailure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcfailure) Error() string {
	return core.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

type jsonraw []byte

func (m jsonraw) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return []byte(m), nil
}

func (m *jsonraw) UnmarshalJSON(data []byte) error {
	if m == nil {
		return core.NewError("chain.jsonraw: nil receiver")
	}
	*m = append((*m)[:0], data...)
	return nil
}

// BlockchainInfo contains the sidechain cache invalidation metadata.
type BlockchainInfo struct {
	Blocks   uint64 `json:"blocks"`
	Headers  uint64 `json:"headers,omitempty"`
	TreeRoot string `json:"treeroot"`
}

// Info contains a lightweight daemon status view.
type Info struct {
	Version     string `json:"version,omitempty"`
	Network     string `json:"network,omitempty"`
	Blocks      uint64 `json:"blocks,omitempty"`
	Height      uint64 `json:"height,omitempty"`
	Connections int    `json:"connections,omitempty"`
}

// Record mirrors the `getnameresource` RPC JSON payload.
type Record struct {
	Type       string   `json:"type"`
	NS         string   `json:"ns,omitempty"`
	Address    string   `json:"address,omitempty"`
	TXT        []string `json:"txt,omitempty"`
	KeyTag     uint16   `json:"keyTag,omitempty"`
	Algorithm  uint8    `json:"algorithm,omitempty"`
	DigestType uint8    `json:"digestType,omitempty"`
	Digest     string   `json:"digest,omitempty"`
}

// NameResource is the RPC representation of a name resource payload.
type NameResource struct {
	Records []Record `json:"records"`
}

// RawCall executes a JSON-RPC call and returns the raw `result` payload.
func (c *Client) RawCall(method string, params any) ([]byte, error) {
	if c == nil || c.url == "" {
		return nil, core.E(opRawCall, "url is required", nil)
	}

	body := core.JSONMarshal(rpcRequest{
		JSONRPC: "2.0",
		ID:      "0",
		Method:  method,
		Params:  params,
	})
	if !body.OK {
		return nil, core.E(opRawCall, "marshal request", body.Value.(error))
	}

	request, err := http.NewRequest(http.MethodPost, c.url, core.NewBuffer(body.Value.([]byte)))
	if err != nil {
		return nil, core.E(opRawCall, "build request", err)
	}

	request.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		request.SetBasicAuth("x", c.apiKey)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, core.E(opRawCall, method, err)
	}
	defer func() { _ = response.Body.Close() }()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, core.E(opRawCall, "read response", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, core.E(opRawCall, core.Sprintf("%s http %d", method, response.StatusCode), nil)
	}

	var decoded rpcResponse
	if r := core.JSONUnmarshal(raw, &decoded); !r.OK {
		return nil, core.E(opRawCall, "decode response", r.Value.(error))
	}

	if decoded.Error != nil {
		return nil, core.E(opRawCall, method, decoded.Error)
	}

	return []byte(decoded.Result), nil
}

func (c *Client) call(method string, params any, out any) error {
	raw, err := c.RawCall(method, params)
	if err != nil {
		return err
	}

	if out == nil || len(raw) == 0 {
		return nil
	}

	if r := core.JSONUnmarshal(raw, out); !r.OK {
		return core.E("chain.Client.call", core.Sprintf("decode %s result", method), r.Value.(error))
	}

	return nil
}

// GetBlockchainInfo returns tree-root and height metadata for cache invalidation.
func (c *Client) GetBlockchainInfo() (*BlockchainInfo, error) {
	var info BlockchainInfo
	if err := c.call("getblockchaininfo", nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetInfo returns lightweight daemon status.
func (c *Client) GetInfo() (*Info, error) {
	var info Info
	if err := c.call("getinfo", nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetNameResource loads the DNS records for a single LNS name.
func (c *Client) GetNameResource(name string) (*NameResource, error) {
	var resource NameResource
	if err := c.call("getnameresource", []any{core.TrimSuffix(core.Lower(name), ".lthn")}, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetBlock returns the raw block payload for a block hash.
func (c *Client) GetBlock(hash string) ([]byte, error) {
	return c.RawCall("getblock", []any{hash})
}

// GetBlockByHeight returns the raw block payload for a block height.
func (c *Client) GetBlockByHeight(height uint64) ([]byte, error) {
	return c.RawCall("getblockbyheight", []any{height})
}

// GetPeerInfo returns the raw peer list payload.
func (c *Client) GetPeerInfo() ([]byte, error) {
	return c.RawCall("getpeerinfo", nil)
}

// Coin mirrors a single spendable output as returned by the HSD daemon's
// `getcoinsbyaddress` RPC (the Coin.getJSON shape in
// docs/js-primitives/coin.js): the funding transaction id and output index, the
// value in dollarydoos, the holding address, and the covenant carried by the
// output. The Hash field is the funding txid in the daemon's little-endian hex
// form (HSD does not reverse name-chain hashes, see
// docs/js-primitives/outpoint.js txid()), so it parses byte-for-byte into a
// primitives outpoint hash. Coinbase coins are flagged so the caller can apply
// the maturity rule.
type Coin struct {
	Version  uint32        `json:"version"`
	Height   int64         `json:"height"`
	Value    uint64        `json:"value"`
	Address  string        `json:"address"`
	Covenant *CoinCovenant `json:"covenant,omitempty"`
	Coinbase bool          `json:"coinbase"`
	Hash     string        `json:"hash"`
	Index    uint32        `json:"index"`
}

// CoinCovenant is the covenant tail of a coin as the daemon serialises it: a
// numeric type and a list of hex-encoded items (Covenant.toJSON in
// docs/js-primitives/covenant.js). A nil or zero-type covenant denotes a plain
// spendable output.
type CoinCovenant struct {
	Type   uint8    `json:"type"`
	Action string   `json:"action,omitempty"`
	Items  []string `json:"items,omitempty"`
}

// GetCoinsByAddress loads the spendable coins (UTXOs) held by an address from
// the daemon via `getcoinsbyaddress`. The daemon returns the address's current
// unspent set; the caller selects inputs from it and applies any
// covenant/maturity policy.
//
//	coins, err := client.GetCoinsByAddress("ts1q...")
func (c *Client) GetCoinsByAddress(address string) ([]Coin, error) {
	if c == nil || c.url == "" {
		return nil, core.E(opGetCoinsByAddress, "url is required", nil)
	}
	if core.Trim(address) == "" {
		return nil, core.E(opGetCoinsByAddress, "address is required", nil)
	}

	var coins []Coin
	if err := c.call(methodGetCoinsByAddr, []any{address}, &coins); err != nil {
		return nil, core.E(opGetCoinsByAddress, address, err)
	}
	return coins, nil
}

// NameOwner is the outpoint that currently owns a name: the funding txid and
// output index of the coin carrying the name's latest covenant. It is the input
// a linked covenant (REGISTER, RENEW, TRANSFER, …) must spend. The Hash is the
// daemon's little-endian txid hex, matching the coin and outpoint encoding.
type NameOwner struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

// nameInfoResult is the daemon's getnameinfo payload. HSD wraps the name state
// under an "info" object (lib/wallet/rpc getnameinfo); the owner outpoint lives
// there. Both the wrapped and a flat shape are decoded so different daemon
// builds work.
type nameInfoResult struct {
	Info *struct {
		Owner *NameOwner `json:"owner"`
	} `json:"info"`
	Owner *NameOwner `json:"owner"`
}

// GetNameOwner returns the current owner outpoint of a name via getnameinfo. It
// is how the write-path resolves the coin a linked covenant must spend when the
// caller does not supply one. A name with no owner (never registered, or expired
// back to an open auction) yields an empty owner.
//
//	owner, err := client.GetNameOwner("example")
func (c *Client) GetNameOwner(name string) (NameOwner, error) {
	if c == nil || c.url == "" {
		return NameOwner{}, core.E(opGetNameOwner, "url is required", nil)
	}
	label := core.TrimSuffix(core.Lower(core.Trim(name)), ".lthn")
	if label == "" {
		return NameOwner{}, core.E(opGetNameOwner, "name is required", nil)
	}

	var info nameInfoResult
	if err := c.call(methodGetNameInfo, []any{label}, &info); err != nil {
		return NameOwner{}, core.E(opGetNameOwner, label, err)
	}

	if info.Info != nil && info.Info.Owner != nil {
		return *info.Info.Owner, nil
	}
	if info.Owner != nil {
		return *info.Owner, nil
	}
	return NameOwner{}, nil
}

// SendRawTransaction broadcasts a fully-signed, hex-encoded transaction through
// the daemon's `sendrawtransaction` RPC and returns the accepted transaction id.
// The daemon validates the transaction (structure, scripts, fees, double-spend)
// against its mempool and consensus rules and relays it to peers; a rejection
// surfaces as the wrapped RPC error. This is the broadcast path the LNS uses
// instead of go-p2p gossip — native local-mempool relay is a separate, future
// path.
//
//	txid, err := client.SendRawTransaction(rawHex)
func (c *Client) SendRawTransaction(rawHex string) (string, error) {
	if c == nil || c.url == "" {
		return "", core.E(opSendRawTransaction, "url is required", nil)
	}
	if core.Trim(rawHex) == "" {
		return "", core.E(opSendRawTransaction, "raw transaction hex is required", nil)
	}

	raw, err := c.RawCall(methodSendRawTx, []any{rawHex})
	if err != nil {
		return "", core.E(opSendRawTransaction, "broadcast transaction", err)
	}

	// The daemon returns the accepted txid, either as a bare JSON string or, on
	// some builds, an object carrying it. Accept both rather than assume one.
	var txid string
	if r := core.JSONUnmarshal(raw, &txid); r.OK && core.Trim(txid) != "" {
		return txid, nil
	}

	var wrapped struct {
		Hash string `json:"hash"`
		TxID string `json:"txid"`
	}
	if r := core.JSONUnmarshal(raw, &wrapped); r.OK {
		if core.Trim(wrapped.Hash) != "" {
			return wrapped.Hash, nil
		}
		if core.Trim(wrapped.TxID) != "" {
			return wrapped.TxID, nil
		}
	}

	return "", core.E(opSendRawTransaction, "daemon returned no transaction id", nil)
}
