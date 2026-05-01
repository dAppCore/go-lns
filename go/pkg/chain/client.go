// SPDX-License-Identifier: EUPL-1.2

// Package chain provides lightweight RPC helpers for the LNS sidechain and
// the main-chain alias discovery path used by the DNS bridge.
package chain

import (
	"io"
	"net/http"
	"time"

	core "dappco.re/go"
)

const opRawCall = "chain.Client.RawCall"

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
