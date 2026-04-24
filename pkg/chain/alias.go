// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	core "dappco.re/go/core"
)

// Alias describes a main-chain alias entry used for LNS discovery.
type Alias struct {
	Name    string
	Address string
	Comment string
}

// AliasClient loads main-chain aliases from the daemon JSON-RPC endpoint.
type AliasClient struct {
	url        string
	httpClient *http.Client
}

// AliasClientOption configures an alias client.
type AliasClientOption func(*AliasClient)

// WithAliasHTTPClient overrides the HTTP transport used by the alias client.
func WithAliasHTTPClient(httpClient *http.Client) AliasClientOption {
	return func(c *AliasClient) {
		if c == nil || httpClient == nil {
			return
		}
		c.httpClient = httpClient
	}
}

// NewAliasClient constructs a main-chain alias discovery client.
func NewAliasClient(url string, opts ...AliasClientOption) *AliasClient {
	url = strings.TrimSpace(url)
	url = strings.TrimRight(url, "/")
	if url != "" && !strings.HasSuffix(url, "/json_rpc") {
		url += "/json_rpc"
	}

	client := &AliasClient{
		url:        url,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	return client
}

// GetAliases returns all main-chain aliases.
func (c *AliasClient) GetAliases() ([]Alias, error) {
	if c == nil || c.url == "" {
		return nil, core.E("chain.AliasClient.GetAliases", "url is required", nil)
	}

	body := []byte(`{"jsonrpc":"2.0","id":"0","method":"get_all_alias_details","params":{}}`)
	request, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return nil, core.E("chain.AliasClient.GetAliases", "build request", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, core.E("chain.AliasClient.GetAliases", "perform request", err)
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, core.E("chain.AliasClient.GetAliases", "read response", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, core.E("chain.AliasClient.GetAliases", "unexpected status", nil)
	}

	var payload struct {
		Result struct {
			Aliases []struct {
				Alias   string `json:"alias"`
				Address string `json:"address"`
				Comment string `json:"comment"`
			} `json:"aliases"`
		} `json:"result"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, core.E("chain.AliasClient.GetAliases", "decode response", err)
	}

	aliases := make([]Alias, 0, len(payload.Result.Aliases))
	for _, alias := range payload.Result.Aliases {
		aliases = append(aliases, Alias{
			Name:    alias.Alias,
			Address: alias.Address,
			Comment: alias.Comment,
		})
	}

	return aliases, nil
}

// ParseAliasComment parses `key=value;...` alias comment metadata.
func ParseAliasComment(comment string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(comment, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		result[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return result
}

// HNSName resolves the LNS name for an alias comment, applying `hns=` overrides.
func HNSName(alias Alias) string {
	name := strings.TrimSpace(alias.Name)
	comment := ParseAliasComment(alias.Comment)
	if hns := strings.TrimSpace(comment["hns"]); hns != "" {
		name = strings.TrimSuffix(strings.TrimSuffix(strings.ToLower(hns), ".lthn."), ".lthn")
	}
	return strings.ToLower(name)
}

// DiscoverNames derives the set of distinct HNS names from aliases.
func DiscoverNames(aliases []Alias, fallback []string) []string {
	if len(aliases) == 0 {
		return append([]string(nil), fallback...)
	}

	seen := make(map[string]struct{}, len(aliases))
	names := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		name := HNSName(alias)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	return names
}
