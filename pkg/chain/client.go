// SPDX-License-Identifier: EUPL-1.2

// Package chain provides lightweight RPC helpers for the LNS sidechain and
// the main-chain alias discovery path used by the DNS bridge.
package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	core "dappco.re/go"
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
		c.apiKey = strings.TrimSpace(apiKey)
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
		url:        strings.TrimSpace(url),
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
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
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
func (c *Client) RawCall(method string, params any) (json.RawMessage, error) {
	if c == nil || c.url == "" {
		return nil, core.E("chain.Client.RawCall", "url is required", nil)
	}

	body, err := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		ID:      "0",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, core.E("chain.Client.RawCall", "marshal request", err)
	}

	request, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return nil, core.E("chain.Client.RawCall", "build request", err)
	}

	request.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		request.SetBasicAuth("x", c.apiKey)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, core.E("chain.Client.RawCall", method, err)
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, core.E("chain.Client.RawCall", "read response", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, core.E("chain.Client.RawCall", fmt.Sprintf("%s http %d", method, response.StatusCode), nil)
	}

	var decoded rpcResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, core.E("chain.Client.RawCall", "decode response", err)
	}

	if decoded.Error != nil {
		return nil, core.E("chain.Client.RawCall", method, decoded.Error)
	}

	return decoded.Result, nil
}

func (c *Client) call(method string, params any, out any) error {
	raw, err := c.RawCall(method, params)
	if err != nil {
		return err
	}

	if out == nil || len(raw) == 0 {
		return nil
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return core.E("chain.Client.call", fmt.Sprintf("decode %s result", method), err)
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
	if err := c.call("getnameresource", []any{strings.TrimSuffix(strings.ToLower(name), ".lthn")}, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetBlock returns the raw block payload for a block hash.
func (c *Client) GetBlock(hash string) (json.RawMessage, error) {
	return c.RawCall("getblock", []any{hash})
}

// GetBlockByHeight returns the raw block payload for a block height.
func (c *Client) GetBlockByHeight(height uint64) (json.RawMessage, error) {
	return c.RawCall("getblockbyheight", []any{height})
}

// GetPeerInfo returns the raw peer list payload.
func (c *Client) GetPeerInfo() (json.RawMessage, error) {
	return c.RawCall("getpeerinfo", nil)
}
