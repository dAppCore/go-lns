// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"crypto/sha3"
	"encoding/hex"
	"slices"
	"strconv"
	"time"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	dnspkg "dappco.re/go/lns/pkg/dns"
	"dappco.re/go/lns/pkg/primitives"
	walletpkg "dappco.re/go/lns/pkg/wallet"
)

type walletBalanceResult struct {
	Confirmed uint64 `json:"confirmed"`
	Locked    uint64 `json:"locked"`
}

type walletReceiveResult struct {
	Address primitives.Address `json:"address"`
	SeedID  string             `json:"seedId"`
	Index   uint32             `json:"index"`
}

type walletSendPlanResult struct {
	SeedID          string              `json:"seedId"`
	DraftID         string              `json:"draftId,omitempty"`
	Address         primitives.Address  `json:"address"`
	Amount          uint64              `json:"amount"`
	Confirmed       uint64              `json:"confirmed"`
	Locked          uint64              `json:"locked"`
	SufficientFunds bool                `json:"sufficientFunds"`
	ChangeAddress   *primitives.Address `json:"changeAddress,omitempty"`
	Recorded        bool                `json:"recorded"`
	PendingEntry    *walletpkg.Entry    `json:"pendingEntry,omitempty"`
	Draft           *nameDraftTx        `json:"draft,omitempty"`
	Warnings        []string            `json:"warnings,omitempty"`
}

type chainHeightResult struct {
	Height   uint64 `json:"height"`
	TreeRoot string `json:"treeRoot"`
}

type chainStatusResult struct {
	Info           *chainpkg.Info           `json:"info,omitempty"`
	BlockchainInfo *chainpkg.BlockchainInfo `json:"blockchainInfo,omitempty"`
	Bridge         *dnspkg.BridgeHealth     `json:"bridge,omitempty"`
}

type chainBlockResult struct {
	Block any `json:"block"`
}

func (s *Service) handleWalletCreate(_ context.Context, _ core.Options) core.Result {
	wallet, err := walletpkg.Create()
	if err != nil {
		return core.Fail(err)
	}

	s.walletMu.Lock()
	s.wallet = wallet
	s.walletMu.Unlock()

	return core.Ok(wallet.Snapshot())
}

func (s *Service) handleWalletImport(_ context.Context, options core.Options) core.Result {
	seed := actionStringFromAny(actionOptionValue(options, "seed", "_input"))
	if seed == "" {
		return core.Fail(core.E("lns.wallet.import", "seed is required", nil))
	}

	wallet, err := walletpkg.ImportSeed(seed)
	if err != nil {
		return core.Fail(err)
	}

	s.walletMu.Lock()
	s.wallet = wallet
	s.walletMu.Unlock()

	return core.Ok(wallet.Snapshot())
}

func (s *Service) handleWalletExport(_ context.Context, _ core.Options) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E("lns.wallet.export", errWalletNotInitialized, nil))
	}
	return core.Ok(wallet.Snapshot())
}

func (s *Service) handleWalletReceive(_ context.Context, options core.Options) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E("lns.wallet.receive", errWalletNotInitialized, nil))
	}

	next, _ := actionBoolValue(options, "next")
	var address primitives.Address
	if next {
		address = wallet.NextReceiveAddress()
	} else {
		address = wallet.ReceiveAddress()
	}

	return core.Ok(
		walletReceiveResult{
			Address: address,
			SeedID:  wallet.SeedID(),
			Index:   wallet.Snapshot().ReceiveIndex,
		})

}

func (s *Service) handleWalletBalance(_ context.Context, _ core.Options) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E("lns.wallet.balance", errWalletNotInitialized, nil))
	}

	confirmed, locked := wallet.Balance()
	return core.Ok(walletBalanceResult{Confirmed: confirmed, Locked: locked})
}

func (s *Service) handleWalletHistory(_ context.Context, _ core.Options) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E("lns.wallet.history", errWalletNotInitialized, nil))
	}

	return core.Ok(wallet.History())
}

func (s *Service) handleWalletSend(_ context.Context, options core.Options) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E(actionNameWalletSend, errWalletNotInitialized, nil))
	}

	address, ok := actionAddressValue(options, "address", "to", "recipient")
	if !ok || !address.IsValid() {
		return core.Fail(core.E(actionNameWalletSend, "recipient address is required", nil))
	}

	amount, ok := actionUint64Value(options, "amount", "value")
	if !ok || amount == 0 {
		return core.Fail(core.E(actionNameWalletSend, "amount is required", nil))
	}

	confirmed, locked := wallet.Balance()
	changeAddress, changePresent := actionAddressValue(options, "change_address")
	if !changePresent || !changeAddress.IsValid() {
		changeAddress = wallet.ReceiveAddress()
		changePresent = true
	}
	suppressAutoChange, _ := actionBoolValue(options, "no_change", "suppress_change", "exact_send")

	placeholderInputs, ok := actionUint32Value(options, "input_count", "inputs")
	if !ok || placeholderInputs == 0 {
		placeholderInputs = 1
	}

	feeRate, ok := actionUint64Value(options, "fee_rate")
	if !ok || feeRate == 0 {
		feeRate = 1
	}

	lockTime, _ := actionUint32Value(options, "lock_time", "locktime")
	notes := []string{
		"Inputs and witnesses are placeholders only; no spendable UTXO selection or signing is performed.",
		"Fee estimation is approximate and intended for planning only.",
	}

	tx := Transaction{
		Version: 0,
		Outputs: []Output{{
			Value:   amount,
			Address: address,
		}},
		LockTime: lockTime,
	}

	fee, feeSet := actionUint64Value(options, "fee")
	if !feeSet {
		fee = uint64(estimateNameDraftTxSize(tx, placeholderInputs)) * feeRate
	}

	autoChange := false
	changeValue, explicitChange := actionUint64Value(options, "change_value")
	if !explicitChange && !suppressAutoChange && confirmed > amount+fee {
		autoChange = true
		changeValue = confirmed - amount - fee
	}

	if changePresent && changeValue > 0 {
		tx.Outputs = append(tx.Outputs, Output{
			Value:   changeValue,
			Address: changeAddress,
		})

		if !feeSet {
			fee = uint64(estimateNameDraftTxSize(tx, placeholderInputs)) * feeRate
			if autoChange {
				if confirmed > amount+fee {
					changeValue = confirmed - amount - fee
					tx.Outputs[1].Value = changeValue
				} else {
					tx.Outputs = tx.Outputs[:1]
					changeValue = 0
				}
			}
		}
	}

	required := amount + fee
	if explicitChange {
		required += changeValue
	}
	sufficientFunds := confirmed >= required
	if !sufficientFunds {
		notes = append(notes, "Local wallet balance does not cover the draft amount plus placeholder fee.")
	}
	if suppressAutoChange && !explicitChange && confirmed > amount+fee {
		notes = append(notes, "Auto-change output was suppressed; leftover confirmed funds remain unassigned in this draft.")
	}

	draft := &nameDraftTx{
		Transaction:         tx,
		EstimatedSize:       estimateNameDraftTxSize(tx, placeholderInputs),
		PlaceholderInputs:   placeholderInputs,
		PlaceholderOutputs:  uint32(len(tx.Outputs)),
		Fee:                 fee,
		FeeRate:             feeRate,
		PrimaryOutputValue:  amount,
		ChangeOutputValue:   changeValue,
		ChangeOutputPresent: len(tx.Outputs) > 1,
		Notes:               notes,
	}

	result := walletSendPlanResult{
		SeedID:          wallet.SeedID(),
		Address:         address,
		Amount:          amount,
		Confirmed:       confirmed,
		Locked:          locked,
		SufficientFunds: sufficientFunds,
		Draft:           draft,
	}
	draftHash, draftID, err := walletDraftID(wallet, draft, noteForWalletSend(options))
	if err != nil {
		return core.Fail(err)
	}
	result.DraftID = draftID

	if len(tx.Outputs) > 1 {
		change := tx.Outputs[1].Address
		result.ChangeAddress = &change
	}
	if !sufficientFunds {
		result.Warnings = append(result.Warnings, "Draft exceeds the locally tracked confirmed wallet balance.")
	}

	recordHistory, _ := actionBoolValue(options, "record_history", "record")
	if recordHistory {
		if !sufficientFunds {
			result.Warnings = append(result.Warnings, "Draft history entry was not recorded because the local balance is insufficient.")
			return core.Ok(result)
		}

		note := noteForWalletSend(options)
		if note == "" {
			note = "draft wallet send"
		}

		entry := wallet.PendingDebit(amount+fee, draftHash, note)
		result.Recorded = true
		result.PendingEntry = &entry
	}

	return core.Ok(result)
}

func walletDraftID(wallet *walletpkg.Wallet, draft *nameDraftTx, note string) (primitives.Hash, string, error) {
	if wallet == nil {
		return primitives.Hash{}, "", core.E(actionNameWalletSend, errWalletNotInitialized, nil)
	}
	if draft == nil {
		return primitives.Hash{}, "", core.E(actionNameWalletSend, "draft transaction is required", nil)
	}

	payload := struct {
		SeedID string       `json:"seedId"`
		Note   string       `json:"note,omitempty"`
		Draft  *nameDraftTx `json:"draft"`
	}{
		SeedID: wallet.SeedID(),
		Note:   note,
		Draft:  draft,
	}

	raw := core.JSONMarshal(payload)
	if !raw.OK {
		return primitives.Hash{}, "", core.E(actionNameWalletSend, "marshal draft id payload", raw.Value.(error))
	}

	hash := primitives.Hash(sha3.Sum256(raw.Value.([]byte)))
	return hash, hex.EncodeToString(hash[:]), nil
}

func noteForWalletSend(options core.Options) string {
	return core.Trim(actionStringFromAny(actionOptionValue(options, "note")))
}

func (s *Service) handleChainHeight(_ context.Context, options core.Options) core.Result {
	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(err)
	}

	info, err := client.GetBlockchainInfo()
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(chainHeightResult{Height: info.Blocks, TreeRoot: info.TreeRoot})
}

func (s *Service) handleChainStatus(_ context.Context, options core.Options) core.Result {
	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(err)
	}

	status := chainStatusResult{}
	if info, err := client.GetInfo(); err == nil {
		status.Info = info
	}
	if info, err := client.GetBlockchainInfo(); err == nil {
		status.BlockchainInfo = info
	}
	if bridge, err := s.bridgeFromOptions(options); err == nil && bridge != nil {
		health := bridge.Health()
		status.Bridge = &health
	}

	return core.Ok(status)
}

func (s *Service) handleChainSync(ctx context.Context, options core.Options) core.Result {
	bridge, err := s.bridgeFromOptions(options)
	if err != nil {
		return core.Fail(err)
	}

	report, err := bridge.Discover(ctx)
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(report)
}

func (s *Service) handleChainPeers(_ context.Context, options core.Options) core.Result {
	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(err)
	}

	raw, err := client.GetPeerInfo()
	if err != nil {
		return core.Fail(err)
	}

	return core.Ok(decodeRawJSON(raw))
}

func (s *Service) handleChainBlock(_ context.Context, options core.Options) core.Result {
	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(err)
	}

	if hash := actionStringFromAny(actionOptionValue(options, "hash")); hash != "" {
		raw, err := client.GetBlock(hash)
		if err != nil {
			return core.Fail(err)
		}
		return core.Ok(chainBlockResult{Block: decodeRawJSON(raw)})
	}

	if height, ok := actionUint32Value(options, "height"); ok {
		raw, err := client.GetBlockByHeight(uint64(height))
		if err != nil {
			return core.Fail(err)
		}
		return core.Ok(chainBlockResult{Block: decodeRawJSON(raw)})
	}

	return core.Fail(core.E("lns.chain.block", "hash or height is required", nil))
}

func (s *Service) currentWallet() (*walletpkg.Wallet, bool) {
	if s == nil {
		return nil, false
	}

	s.walletMu.Lock()
	defer s.walletMu.Unlock()

	if s.wallet == nil {
		return nil, false
	}

	return s.wallet, true
}

func (s *Service) bridgeFromOptions(options core.Options) (*dnspkg.Bridge, error) {
	if s == nil {
		return nil, core.E(actionOpBridge, "service is required", nil)
	}

	config, ok := bridgeConfigFromOptions(options)
	if !ok {
		s.bridgeMu.Lock()
		defer s.bridgeMu.Unlock()
		if s.dnsBridge == nil {
			return nil, core.E(actionOpBridge, "chain_url is required", nil)
		}
		return s.dnsBridge, nil
	}

	s.bridgeMu.Lock()
	defer s.bridgeMu.Unlock()

	if s.dnsBridge == nil || !bridgeConfigsEqual(s.bridgeConfig, config) {
		if s.dnsBridge != nil {
			if err := s.dnsBridge.Shutdown(); err != nil {
				return nil, core.Wrap(err, actionOpBridge, "shutdown existing DNS bridge")
			}
		}
		s.bridgeConfig = config
		s.dnsBridge = dnspkg.NewBridge(config)
	}

	return s.dnsBridge, nil
}

func chainClientFromOptions(options core.Options) (*chainpkg.Client, error) {
	url := core.Trim(actionStringFromAny(actionOptionValue(options, "chain_url", "hsd_url", "url")))
	if url == "" {
		return nil, core.E("lns.chain", "chain_url is required", nil)
	}

	apiKey := core.Trim(actionStringFromAny(actionOptionValue(options, "api_key", "hsd_api_key")))
	return chainpkg.NewClient(url, chainpkg.WithAPIKey(apiKey)), nil
}

func bridgeConfigFromOptions(options core.Options) (dnspkg.BridgeConfig, bool) {
	chainURL := core.Trim(actionStringFromAny(actionOptionValue(options, "chain_url", "hsd_url", "url")))
	if chainURL == "" {
		return dnspkg.BridgeConfig{}, false
	}

	config := dnspkg.BridgeConfig{
		ChainURL:      chainURL,
		AliasURL:      core.Trim(actionStringFromAny(actionOptionValue(options, "alias_url", "daemon_url"))),
		APIKey:        core.Trim(actionStringFromAny(actionOptionValue(options, "api_key", "hsd_api_key"))),
		DNSAddr:       dnsListenAddr(options),
		CheckInterval: bridgeCheckInterval(options),
		KnownNames:    actionStringSliceValue(actionOptionValue(options, "known_names", "names")),
	}

	return config, true
}

func dnsListenAddr(options core.Options) string {
	if addr := core.Trim(actionStringFromAny(actionOptionValue(options, "dns_addr", "dns_listen"))); addr != "" {
		return addr
	}
	if port := core.Trim(actionStringFromAny(actionOptionValue(options, "dns_port"))); port != "" {
		if core.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":5354"
}

func bridgeCheckInterval(options core.Options) time.Duration {
	if seconds, ok := actionUint32Value(options, "check_interval", "check_seconds"); ok {
		return time.Duration(seconds) * time.Second
	}
	if text := core.Trim(actionStringFromAny(actionOptionValue(options, "check_interval"))); text != "" {
		if duration, err := time.ParseDuration(text); err == nil {
			return duration
		}
		if seconds, err := strconv.Atoi(text); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return 15 * time.Second
}

func actionStringSliceValue(value any) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if text := core.Trim(actionStringFromAny(item)); text != "" {
				out = append(out, text)
			}
		}
		return out
	case string:
		if typed == "" {
			return nil
		}
		parts := core.Split(typed, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			if text := core.Trim(part); text != "" {
				out = append(out, text)
			}
		}
		return out
	default:
		return nil
	}
}

func bridgeConfigsEqual(left, right dnspkg.BridgeConfig) bool {
	return left.AliasURL == right.AliasURL &&
		left.ChainURL == right.ChainURL &&
		left.APIKey == right.APIKey &&
		left.DNSAddr == right.DNSAddr &&
		left.CheckInterval == right.CheckInterval &&
		slices.Equal(left.KnownNames, right.KnownNames)
}

func decodeRawJSON(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}

	var decoded any
	if r := core.JSONUnmarshal(raw, &decoded); r.OK {
		return decoded
	}

	return string(raw)
}
