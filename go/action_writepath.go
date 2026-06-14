// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	"dappco.re/go/lns/pkg/primitives"
	walletpkg "dappco.re/go/lns/pkg/wallet"
)

// This file is the glue between the action layer and the write-path: it parses
// broadcast intent and funding data from action options, runs the
// build→sign→broadcast pipeline, and shapes the action response. The write-path
// itself (coin selection, assembly, signing, broadcast) lives in writepath.go.

// walletSendResult is the live (broadcast) result of lns.wallet.send: the
// recipient and amount, the broadcast outcome (txid, raw hex, fee), and the
// updated local balance view.
type walletSendResult struct {
	SeedID    string             `json:"seedId"`
	Address   primitives.Address `json:"address"`
	Amount    uint64             `json:"amount"`
	Broadcast *broadcastOutcome  `json:"broadcast"`
	Confirmed uint64             `json:"confirmed"`
	Locked    uint64             `json:"locked"`
}

// nameBroadcastResult is the live result of a broadcast lns.name.* action: the
// canonical name plan plus the broadcast outcome for the covenant transaction.
type nameBroadcastResult struct {
	Action    string            `json:"action"`
	Input     string            `json:"input"`
	Canonical string            `json:"canonical"`
	Covenant  CovenantType      `json:"covenantType"`
	Broadcast *broadcastOutcome `json:"broadcast"`
}

// walletSendShouldBroadcast reports whether the caller asked lns.wallet.send to
// build, sign, and broadcast a real transaction rather than return a draft. The
// draft path remains the default so existing planning callers are unaffected.
func walletSendShouldBroadcast(options core.Options) bool {
	if value, ok := actionBoolValue(options, "broadcast", "send", "live"); ok {
		return value
	}
	return false
}

// broadcastWalletSend runs the write-path for a plain P2WPKH spend to a recipient
// and returns the broadcast outcome. It is the live counterpart of the
// draft-building wallet send.
func (s *Service) broadcastWalletSend(wallet *walletpkg.Wallet, address primitives.Address, amount uint64, options core.Options) core.Result {
	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(core.E(actionNameWalletSend, "chain_url is required to broadcast", err))
	}

	req, err := spendRequestFromOptions(wallet, client, options)
	if err != nil {
		return core.Fail(core.E(actionNameWalletSend, "build spend request", err))
	}

	changeAddress, changePresent := actionAddressValue(options, "change_address")
	if !changePresent || !changeAddress.IsValid() {
		changeAddress = wallet.ReceiveAddress()
	}

	req.outputs = []primitives.Output{{Value: amount, Address: address}}
	req.changeAddress = changeAddress
	req.allowChange = !boolOption(options, false, "no_change", "suppress_change", "exact_send")

	outcome, err := s.execSpend(req)
	if err != nil {
		return core.Fail(core.E(actionNameWalletSend, "broadcast wallet send", err))
	}

	confirmed, locked := wallet.Balance()
	return core.Ok(walletSendResult{
		SeedID:    wallet.SeedID(),
		Address:   address,
		Amount:    amount,
		Broadcast: outcome,
		Confirmed: confirmed,
		Locked:    locked,
	})
}

// broadcastNameCovenant builds, signs, and broadcasts a single covenant
// transaction for a name action and returns the live result. The covenant output
// (the owner/recipient output already carrying the covenant) is created by the
// action's draft builder; this funds it from the wallet, attaches the name's
// owner coin as a required input for linked covenants, signs every input, and
// broadcasts through the daemon. A linked covenant with no resolvable owner coin
// is an error, not a silent draft.
func (s *Service) broadcastNameCovenant(
	action string,
	lookup nameLookupResult,
	covenant Covenant,
	outputAddress primitives.Address,
	outputValue uint64,
	options core.Options,
) core.Result {
	wallet, ok := s.currentWallet()
	if !ok {
		return core.Fail(core.E(action, errWalletNotInitialized, nil))
	}

	client, err := chainClientFromOptions(options)
	if err != nil {
		return core.Fail(core.E(action, "chain_url is required to broadcast", err))
	}

	req, err := spendRequestFromOptions(wallet, client, options)
	if err != nil {
		return core.Fail(core.E(action, "build spend request", err))
	}

	// A linked covenant (REVEAL, REGISTER, RENEW, TRANSFER, UPDATE, FINALIZE,
	// REVOKE) spends the name's current owner coin. The owner outpoint is chain
	// state: the caller supplies it (owner_outpoint) or it is resolved from the
	// daemon's name info. Without it the transaction would be invalid, so reject
	// rather than broadcast a malformed transaction.
	if covenant.IsLinked() && len(req.requiredInputs) == 0 {
		owner, err := resolveNameOwner(client, lookup.Canonical, options)
		if err != nil {
			return core.Fail(core.E(action,
				"the name owner coin is required to broadcast this covenant; supply owner_outpoint or a reachable chain_url", err))
		}
		req.requiredInputs = []primitives.Outpoint{owner}
	}

	changeAddress, changePresent := actionAddressValue(options, "change_address")
	if !changePresent || !changeAddress.IsValid() {
		changeAddress = wallet.ReceiveAddress()
	}

	req.outputs = []primitives.Output{{Value: outputValue, Address: outputAddress, Covenant: covenant}}
	req.changeAddress = changeAddress
	req.allowChange = !boolOption(options, false, "no_change", "suppress_change", "exact_send")

	outcome, err := s.execSpend(req)
	if err != nil {
		return core.Fail(core.E(action, "broadcast covenant transaction", err))
	}

	return core.Ok(nameBroadcastResult{
		Action:    action,
		Input:     lookup.Input,
		Canonical: lookup.Canonical,
		Covenant:  CovenantType(covenant.Type),
		Broadcast: outcome,
	})
}

// resolveNameOwner returns the outpoint of a name's current owner coin. It first
// honours a caller-supplied owner_outpoint, then queries the daemon name info.
// It fails when neither source yields an owner (an unowned or unknown name, or an
// unreachable daemon).
func resolveNameOwner(client *chainpkg.Client, canonical string, options core.Options) (primitives.Outpoint, error) {
	if single := actionOptionValue(options, "owner_outpoint"); single != nil {
		return outpointFromAny(single)
	}

	owner, err := client.GetNameOwner(canonical)
	if err != nil {
		return primitives.Outpoint{}, err
	}
	if owner.Hash == "" {
		return primitives.Outpoint{}, core.E("lns.writepath.owner", "name has no owner coin", nil)
	}

	hash, err := hashFromHex(owner.Hash)
	if err != nil {
		return primitives.Outpoint{}, core.E("lns.writepath.owner", "decode owner txid", err)
	}
	return primitives.Outpoint{TxHash: hash, Index: owner.Index}, nil
}

// nameWantsBroadcast reports whether a name action should build, sign, and
// broadcast a covenant transaction rather than return a planning draft. The
// draft path stays the default.
func nameWantsBroadcast(options core.Options) bool {
	if value, ok := actionBoolValue(options, "broadcast", "send", "live"); ok {
		return value
	}
	return false
}

// spendRequestFromOptions assembles the shared parts of a spend request from
// action options: the fee controls, the bech32 address prefix for daemon coin
// queries, and any caller-supplied coins. The caller fills in outputs, change,
// and required inputs.
func spendRequestFromOptions(wallet *walletpkg.Wallet, client *chainpkg.Client, options core.Options) (spendRequest, error) {
	req := spendRequest{
		wallet:        wallet,
		client:        client,
		lockTime:      uint32Option(options, 0, "lock_time", "locktime"),
		addressPrefix: stringOption(options, "address_prefix", "hrp", "prefix"),
	}

	if feeRate, ok := actionUint64Value(options, "fee_rate"); ok && feeRate > 0 {
		req.feeRate = feeRate
	}
	if fee, ok := actionUint64Value(options, "fee"); ok {
		req.explicitFee = fee
		req.hasFee = true
	}

	coins, err := coinsFromOptions(options)
	if err != nil {
		return spendRequest{}, err
	}
	req.suppliedCoins = coins

	required, err := requiredInputsFromOptions(options)
	if err != nil {
		return spendRequest{}, err
	}
	req.requiredInputs = required

	return req, nil
}

// coinsFromOptions parses a caller-supplied coin (UTXO) set from the "coins" or
// "utxos" option. Each entry needs a funding txid, output index, value, and the
// holding wallet address (so the signing key can be resolved). It is the offline
// funding path and the source used when the daemon address query is unavailable.
func coinsFromOptions(options core.Options) ([]selectedCoin, error) {
	value := actionOptionValue(options, "coins", "utxos", "inputs_data")
	if value == nil {
		return nil, nil
	}

	items, ok := value.([]any)
	if !ok {
		return nil, core.E("lns.writepath.coins", "coins must be a list of {hash,index,value,address}", nil)
	}

	coins := make([]selectedCoin, 0, len(items))
	for _, raw := range items {
		coin, err := coinFromOptionEntry(raw)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

// coinFromOptionEntry parses one caller-supplied coin object.
func coinFromOptionEntry(raw any) (selectedCoin, error) {
	entry, ok := raw.(map[string]any)
	if !ok {
		return selectedCoin{}, core.E("lns.writepath.coins", "each coin must be an object", nil)
	}

	hashText := actionStringFromAny(entry["hash"])
	if hashText == "" {
		hashText = actionStringFromAny(entry["txid"])
	}
	hash, err := hashFromHex(hashText)
	if err != nil {
		return selectedCoin{}, core.E("lns.writepath.coins", "invalid coin hash", err)
	}

	index, _ := actionUint32FromAny(entry["index"])
	coinValue, ok := actionUint64FromAny(entry["value"])
	if !ok || coinValue == 0 {
		return selectedCoin{}, core.E("lns.writepath.coins", "coin value is required", nil)
	}

	addr, ok := addressFromAny(entry["address"])
	if !ok || !addr.IsPubkeyHash() {
		return selectedCoin{}, core.E("lns.writepath.coins", "coin address must be a version-0 pubkey hash", nil)
	}

	return selectedCoin{
		outpoint: primitives.Outpoint{TxHash: hash, Index: index},
		value:    coinValue,
		address:  addr,
	}, nil
}

// requiredInputsFromOptions parses outpoints the transaction must spend (a name's
// current owner coin) from the "owner_outpoint" / "required_inputs" option. A
// linked covenant action needs the name owner coin as an input; the caller
// supplies it (or it is resolved from the daemon name info when available).
func requiredInputsFromOptions(options core.Options) ([]primitives.Outpoint, error) {
	value := actionOptionValue(options, "required_inputs", "owner_outpoints")
	if value == nil {
		if single := actionOptionValue(options, "owner_outpoint"); single != nil {
			outpoint, err := outpointFromAny(single)
			if err != nil {
				return nil, err
			}
			return []primitives.Outpoint{outpoint}, nil
		}
		return nil, nil
	}

	items, ok := value.([]any)
	if !ok {
		return nil, core.E("lns.writepath.inputs", "required_inputs must be a list", nil)
	}

	out := make([]primitives.Outpoint, 0, len(items))
	for _, raw := range items {
		outpoint, err := outpointFromAny(raw)
		if err != nil {
			return nil, err
		}
		out = append(out, outpoint)
	}
	return out, nil
}

// outpointFromAny parses an outpoint from a {hash,index} object or the typed
// primitives.Outpoint.
func outpointFromAny(value any) (primitives.Outpoint, error) {
	switch typed := value.(type) {
	case primitives.Outpoint:
		return typed, nil
	case map[string]any:
		hashText := actionStringFromAny(typed["hash"])
		if hashText == "" {
			hashText = actionStringFromAny(typed["txid"])
		}
		hash, err := hashFromHex(hashText)
		if err != nil {
			return primitives.Outpoint{}, core.E("lns.writepath.inputs", "invalid outpoint hash", err)
		}
		index, _ := actionUint32FromAny(typed["index"])
		return primitives.Outpoint{TxHash: hash, Index: index}, nil
	default:
		return primitives.Outpoint{}, core.E("lns.writepath.inputs", "unsupported outpoint value", nil)
	}
}

// boolOption returns the first boolean option found, or fallback.
func boolOption(options core.Options, fallback bool, keys ...string) bool {
	if value, ok := actionBoolValue(options, keys...); ok {
		return value
	}
	return fallback
}

// uint32Option returns the first uint32 option found, or fallback.
func uint32Option(options core.Options, fallback uint32, keys ...string) uint32 {
	if value, ok := actionUint32Value(options, keys...); ok {
		return value
	}
	return fallback
}

// stringOption returns the first non-empty string option found.
func stringOption(options core.Options, keys ...string) string {
	return core.Trim(actionStringFromAny(actionOptionValue(options, keys...)))
}
