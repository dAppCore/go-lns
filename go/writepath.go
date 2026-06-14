// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"encoding/hex"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	"dappco.re/go/lns/pkg/mempool"
	"dappco.re/go/lns/pkg/primitives"
	"dappco.re/go/lns/pkg/script"
	walletpkg "dappco.re/go/lns/pkg/wallet"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// The write-path turns the unsigned drafts produced by the action layer into
// real signed transactions and broadcasts them through the daemon RPC. It is the
// shared machinery behind lns.wallet.send and the lns.name.* covenant actions:
//
//	coins → select → assemble (inputs + outputs + change) → sign every input →
//	verify each input → serialise with witnesses → sendrawtransaction → txid
//
// Broadcast goes through chain.Client.SendRawTransaction (the HSD daemon's
// sendrawtransaction RPC), not p2p gossip. The daemon validates and relays the
// transaction; native local-mempool relay over go-p2p is a separate future path.
const (
	opWritePath  = "lns.writepath"
	opBuildTx    = "lns.writepath.build"
	opSignTx     = "lns.writepath.sign"
	opSelectCoin = "lns.writepath.select"

	// defaultFeeRate is the fee charged per serialised byte (dollarydoos/byte)
	// when the caller supplies neither an explicit fee nor a fee_rate. It
	// matches the planning default the draft builders already use. The daemon
	// applies its own relay-fee floor; a transaction the daemon rejects as
	// underpaying surfaces as a broadcast error rather than a silent success.
	defaultFeeRate uint64 = 1

	// dustThreshold is the minimum value of a change output. Change below this
	// is dropped into the fee instead of creating an uneconomical output. It is
	// a conservative local policy, not a consensus value.
	dustThreshold uint64 = 200

	// sequenceFinal marks an input final (locktime disabled for that input).
	sequenceFinal uint32 = 0xffffffff

	// witnessPlaceholderSize is the per-input witness byte cost used when
	// sizing the fee before signatures exist: a P2WPKH witness is
	// varint(2) + varint(65) + 65 + varint(33) + 33 = 101 bytes.
	witnessPlaceholderSize = 101
)

// broadcastOutcome captures the result of pushing a signed transaction to the
// daemon: the canonical txid (BLAKE2b-256 over the non-witness serialisation),
// the daemon-reported id from sendrawtransaction, the raw hex that was sent, and
// the fee paid.
type broadcastOutcome struct {
	TxID         string `json:"txid"`
	DaemonTxID   string `json:"daemonTxid,omitempty"`
	RawTx        string `json:"rawTx"`
	Fee          uint64 `json:"fee"`
	FeeRate      uint64 `json:"feeRate"`
	InputCount   int    `json:"inputCount"`
	OutputCount  int    `json:"outputCount"`
	ChangeValue  uint64 `json:"changeValue,omitempty"`
	ChangeOutput bool   `json:"changeOutput"`
	TotalInput   uint64 `json:"totalInput"`
	TotalOutput  uint64 `json:"totalOutput"`
}

// spendRequest is the input to the write-path: the covenant outputs to create
// (recipient/owner outputs already carrying their covenants), the funding wallet,
// the daemon client, the change address, the fee controls, and any required
// inputs that must be spent regardless of selection (a name's current owner
// coin for linked covenants).
type spendRequest struct {
	wallet        *walletpkg.Wallet
	client        *chainpkg.Client
	outputs       []primitives.Output
	changeAddress primitives.Address
	lockTime      uint32
	feeRate       uint64
	explicitFee   uint64
	hasFee        bool
	allowChange   bool
	// addressPrefix is the network bech32 human-readable prefix used to render
	// wallet addresses for the daemon's getcoinsbyaddress query. The daemon keys
	// coins by bech32 address and the prefix is network-specific, so the caller
	// supplies it (chain context this package cannot infer).
	addressPrefix string
	// suppliedCoins lets the caller pass the spendable coin set directly,
	// bypassing the daemon query. It is the offline/testable funding path and
	// the source used when the daemon address query is unavailable.
	suppliedCoins []selectedCoin
	// requiredInputs are outpoints that must be included as inputs (e.g. a
	// name's owner coin). Their values are looked up from the available coin set.
	requiredInputs []primitives.Outpoint
}

// selectedCoin is a coin chosen as a transaction input, carrying the value and
// holding address needed to size the fee, sign, and verify.
type selectedCoin struct {
	outpoint primitives.Outpoint
	value    uint64
	address  primitives.Address
}

// execSpend runs the full write-path for a spend request and returns the
// broadcast outcome. It selects funding coins, assembles the transaction with a
// fee-and-change pass, signs and self-verifies every input, serialises with
// witnesses, and broadcasts through the daemon. Any failure (insufficient funds,
// signing error, self-verification failure, broadcast rejection) is returned as
// a wrapped error; the function never reports a fake success.
func (s *Service) execSpend(req spendRequest) (*broadcastOutcome, error) {
	if req.wallet == nil {
		return nil, core.E(opWritePath, errWalletNotInitialized, nil)
	}
	if req.client == nil {
		return nil, core.E(opWritePath, "chain_url is required to broadcast", nil)
	}
	if len(req.outputs) == 0 {
		return nil, core.E(opWritePath, "at least one output is required", nil)
	}

	feeRate := req.feeRate
	if feeRate == 0 {
		feeRate = defaultFeeRate
	}

	// Gather the wallet's spendable coins from the daemon.
	available, err := s.gatherCoins(req)
	if err != nil {
		return nil, err
	}

	// Pull in any required inputs (name owner coin) first; they are mandatory.
	mandatory, mandatoryValue, err := resolveRequiredInputs(req.requiredInputs, available)
	if err != nil {
		return nil, err
	}

	outputTotal := uint64(0)
	for i := range req.outputs {
		outputTotal += req.outputs[i].Value
	}

	tx, selected, fee, changeValue, err := assembleTransaction(req, available, mandatory, mandatoryValue, outputTotal, feeRate)
	if err != nil {
		return nil, err
	}

	// Sign every input: the mandatory (owner) coins and the selected funding
	// coins together form the complete spent-coin set.
	spent := make([]selectedCoin, 0, len(mandatory)+len(selected))
	spent = append(spent, mandatory...)
	spent = append(spent, selected...)
	if err := s.signTransaction(&tx, req.wallet, spent); err != nil {
		return nil, err
	}

	raw, err := mempool.SerializeWithWitness(tx)
	if err != nil {
		return nil, core.E(opWritePath, "serialise signed transaction", err)
	}

	txid, err := mempool.ComputeTxID(tx)
	if err != nil {
		return nil, core.E(opWritePath, "compute txid", err)
	}

	rawHex := hex.EncodeToString(raw)
	daemonTxID, err := req.client.SendRawTransaction(rawHex)
	if err != nil {
		return nil, core.E(opWritePath, "broadcast transaction", err)
	}

	totalInput := mandatoryValue
	for i := range selected {
		totalInput += selected[i].value
	}

	return &broadcastOutcome{
		TxID:         hex.EncodeToString(txid[:]),
		DaemonTxID:   daemonTxID,
		RawTx:        rawHex,
		Fee:          fee,
		FeeRate:      feeRate,
		InputCount:   len(tx.Inputs),
		OutputCount:  len(tx.Outputs),
		ChangeValue:  changeValue,
		ChangeOutput: changeValue > 0,
		TotalInput:   totalInput,
		TotalOutput:  outputTotal + changeValue,
	}, nil
}

// gatherCoins resolves the wallet's spendable coin set. When the caller supplies
// coins directly they are used as-is (the offline/testable funding path).
// Otherwise it queries the daemon's getcoinsbyaddress for each derived wallet
// address. Only version-0 P2WPKH coins with a spendable covenant are returned:
// NONE (plain), OPEN, and REDEEM outputs can be spent for funding; outputs locked
// by an active name covenant (BID/REVEAL/REGISTER/…) are excluded from generic
// funding and only spent when named explicitly as a required input.
func (s *Service) gatherCoins(req spendRequest) ([]selectedCoin, error) {
	if len(req.suppliedCoins) > 0 {
		return dedupeCoins(req.suppliedCoins), nil
	}

	if core.Trim(req.addressPrefix) == "" {
		return nil, core.E(opSelectCoin,
			"address_prefix (network bech32 HRP) is required to query coins by address; supply coins explicitly otherwise", nil)
	}

	refs := req.wallet.SpendableAddresses(0)
	if len(refs) == 0 {
		return nil, core.E(opSelectCoin, "wallet has no derivable addresses", nil)
	}

	seen := make(map[string]struct{})
	coins := make([]selectedCoin, 0, len(refs))
	for i := range refs {
		addr := refs[i].Address
		encoded, err := addr.Encode(req.addressPrefix)
		if err != nil {
			return nil, core.E(opSelectCoin, "encode wallet address", err)
		}
		daemonCoins, err := req.client.GetCoinsByAddress(encoded)
		if err != nil {
			return nil, core.E(opSelectCoin, "load coins for address", err)
		}
		for j := range daemonCoins {
			coin, ok, err := coinFromDaemon(daemonCoins[j], addr)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			key := outpointKey(coin.outpoint)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			coins = append(coins, coin)
		}
	}
	return coins, nil
}

// dedupeCoins removes duplicate coins (by outpoint) from a caller-supplied set,
// keeping the first occurrence.
func dedupeCoins(coins []selectedCoin) []selectedCoin {
	seen := make(map[string]struct{}, len(coins))
	out := make([]selectedCoin, 0, len(coins))
	for i := range coins {
		key := outpointKey(coins[i].outpoint)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, coins[i])
	}
	return out
}

// coinFromDaemon converts a daemon coin payload into a selectedCoin, keeping only
// fundable version-0 P2WPKH coins. It returns ok=false for coins that are not
// generic funding candidates (non-P2WPKH programmes or outputs still locked by an
// active name covenant). The funding address is supplied by the caller because
// the daemon reports the address as a bech32 string this package does not decode;
// the coin is matched to the wallet address that was queried.
func coinFromDaemon(coin chainpkg.Coin, addr primitives.Address) (selectedCoin, bool, error) {
	txHash, err := hashFromHex(coin.Hash)
	if err != nil {
		return selectedCoin{}, false, core.E(opSelectCoin, "decode coin txid", err)
	}
	if coin.Value == 0 {
		return selectedCoin{}, false, nil
	}
	if !addr.IsPubkeyHash() {
		return selectedCoin{}, false, nil
	}
	if coin.Covenant != nil && !covenantIsFundable(coin.Covenant.Type) {
		return selectedCoin{}, false, nil
	}

	return selectedCoin{
		outpoint: primitives.Outpoint{TxHash: txHash, Index: coin.Index},
		value:    coin.Value,
		address:  addr,
	}, true, nil
}

// covenantIsFundable reports whether a coin carrying the given covenant type may
// be spent as generic funding. NONE (0), OPEN (2), and REDEEM (5) leave the coin
// spendable; every other type locks it to a name action and must not be consumed
// for fees. This mirrors primitives.Covenant.IsNonspendable.
func covenantIsFundable(covType uint8) bool {
	switch covType {
	case 0, 2, 5:
		return true
	default:
		return false
	}
}

// resolveRequiredInputs matches each required outpoint to a coin in the
// available set so its value is known for signing and fee accounting. A required
// input that is not present in the wallet's coin set is an error: the
// transaction cannot spend a coin the wallet does not own.
func resolveRequiredInputs(required []primitives.Outpoint, available []selectedCoin) ([]selectedCoin, uint64, error) {
	if len(required) == 0 {
		return nil, 0, nil
	}

	index := make(map[string]selectedCoin, len(available))
	for i := range available {
		index[outpointKey(available[i].outpoint)] = available[i]
	}

	out := make([]selectedCoin, 0, len(required))
	total := uint64(0)
	for _, outpoint := range required {
		coin, ok := index[outpointKey(outpoint)]
		if !ok {
			return nil, 0, core.E(opSelectCoin, "required input coin is not owned by the wallet", nil)
		}
		out = append(out, coin)
		total += coin.value
	}
	return out, total, nil
}

// assembleTransaction selects funding inputs to cover the output total plus fee,
// adds a change output when economical, and returns the unsigned transaction
// (witnesses empty), the funding coins selected, the fee paid, and the change
// value. Mandatory inputs are always present and their value counts toward the
// target. Fee is recomputed as inputs and the change output change the size, so
// the final fee reflects the actual transaction shape. Insufficient funds is an
// error.
func assembleTransaction(
	req spendRequest,
	available, mandatory []selectedCoin,
	mandatoryValue, outputTotal, feeRate uint64,
) (primitives.Transaction, []selectedCoin, uint64, uint64, error) {
	// Candidates exclude the mandatory inputs (already committed) and are
	// ordered largest-first so the input count stays small.
	candidates := excludeCoins(available, mandatory)
	sortCoinsDesc(candidates)

	selected := make([]selectedCoin, 0, len(candidates))
	selectedValue := uint64(0)

	build := func(includeChange bool) (primitives.Transaction, uint64, uint64) {
		inputs := make([]primitives.Input, 0, len(mandatory)+len(selected))
		for i := range mandatory {
			inputs = append(inputs, primitives.Input{Prevout: mandatory[i].outpoint, Sequence: sequenceFinal})
		}
		for i := range selected {
			inputs = append(inputs, primitives.Input{Prevout: selected[i].outpoint, Sequence: sequenceFinal})
		}

		outputs := make([]primitives.Output, len(req.outputs))
		copy(outputs, req.outputs)

		totalIn := mandatoryValue + selectedValue
		fee := req.explicitFee
		change := uint64(0)
		if includeChange && req.allowChange {
			outputs = append(outputs, primitives.Output{Value: 0, Address: req.changeAddress})
		}

		tx := primitives.Transaction{Version: 0, Inputs: inputs, Outputs: outputs, LockTime: req.lockTime}
		if !req.hasFee {
			fee = estimateFee(tx, len(inputs), feeRate)
		}

		if includeChange && req.allowChange {
			if totalIn > outputTotal+fee {
				change = totalIn - outputTotal - fee
				if change < dustThreshold {
					// Drop uneconomical change back into the fee.
					tx.Outputs = tx.Outputs[:len(tx.Outputs)-1]
					change = 0
					if !req.hasFee {
						fee = estimateFee(tx, len(inputs), feeRate)
					}
				} else {
					tx.Outputs[len(tx.Outputs)-1].Value = change
				}
			} else {
				tx.Outputs = tx.Outputs[:len(tx.Outputs)-1]
				change = 0
			}
		}
		return tx, fee, change
	}

	// Grow the selection until the inputs cover outputs + fee.
	for {
		tx, fee, change := build(true)
		totalIn := mandatoryValue + selectedValue
		if totalIn >= outputTotal+fee {
			return tx, selected, fee, change, nil
		}
		if len(candidates) == 0 {
			return primitives.Transaction{}, nil, 0, 0, core.E(
				opBuildTx,
				core.Sprintf("insufficient funds: have %d, need %d", totalIn, outputTotal+fee),
				nil,
			)
		}
		next := candidates[0]
		candidates = candidates[1:]
		selected = append(selected, next)
		selectedValue += next.value
	}
}

// estimateFee returns the fee for a transaction at the given byte rate, sizing
// the witnesses with a P2WPKH placeholder for each input (the real witnesses are
// attached after signing but cost the same number of bytes). It is the fee the
// write-path commits to before signing.
func estimateFee(tx primitives.Transaction, inputCount int, feeRate uint64) uint64 {
	base, err := mempool.SerializeNoWitness(tx)
	if err != nil {
		// A serialisation failure here means malformed outputs; charge a
		// non-zero fee so the caller's sufficiency check still runs and the
		// daemon ultimately rejects a malformed transaction.
		return feeRate
	}
	size := len(base) + inputCount*witnessPlaceholderSize
	return uint64(size) * feeRate
}

// signTransaction signs every input of tx in place with the P2WPKH witness for
// the coin it spends, then self-verifies each input with script.VerifyInput
// before the transaction is broadcast. The coins slice must be the inputs in
// transaction order: mandatory inputs first, then the selected funding coins.
func (s *Service) signTransaction(tx *primitives.Transaction, wallet *walletpkg.Wallet, selected []selectedCoin) error {
	if tx == nil {
		return core.E(opSignTx, "nil transaction", nil)
	}
	if len(tx.Inputs) == 0 {
		return core.E(opSignTx, "transaction has no inputs to sign", nil)
	}

	// Reconstruct the per-input coin list in input order. tx.Inputs already
	// holds mandatory inputs first, then selected; rebuild the spent-coin view
	// keyed by outpoint so signing tolerates any ordering.
	coinByOutpoint := make(map[string]selectedCoin, len(tx.Inputs))
	for i := range selected {
		coinByOutpoint[outpointKey(selected[i].outpoint)] = selected[i]
	}

	for index := range tx.Inputs {
		coin, ok := coinByOutpoint[outpointKey(tx.Inputs[index].Prevout)]
		if !ok {
			return core.E(opSignTx, "no spent-coin context for input", nil)
		}

		privBytes, pubBytes, err := wallet.SigningKey(coin.address, 0)
		if err != nil {
			return core.E(opSignTx, "resolve signing key", err)
		}

		prev, err := script.FromPubkeyhash(coin.address.Hash)
		if err != nil {
			return core.E(opSignTx, "build previous script", err)
		}

		priv := secp256k1.PrivKeyFromBytes(privBytes)
		sig, err := script.Sign(*tx, index, prev, coin.value, priv, script.SighashAll)
		if err != nil {
			return core.E(opSignTx, "sign input", err)
		}

		tx.Inputs[index].Witness = [][]byte{sig, pubBytes}

		prevOutput := primitives.Output{Value: coin.value, Address: coin.address}
		if verr := script.VerifyInput(prevOutput, *tx, index, script.MandatoryVerifyFlags); verr != nil {
			return core.E(opSignTx, "self-verify signed input", verr)
		}
	}

	return nil
}

// excludeCoins returns the coins in src that are not present in remove, compared
// by outpoint.
func excludeCoins(src, remove []selectedCoin) []selectedCoin {
	if len(remove) == 0 {
		return append([]selectedCoin(nil), src...)
	}
	removed := make(map[string]struct{}, len(remove))
	for i := range remove {
		removed[outpointKey(remove[i].outpoint)] = struct{}{}
	}
	out := make([]selectedCoin, 0, len(src))
	for i := range src {
		if _, skip := removed[outpointKey(src[i].outpoint)]; skip {
			continue
		}
		out = append(out, src[i])
	}
	return out
}

// sortCoinsDesc orders coins by value, largest first, so the selection loop
// reaches the target with the fewest inputs. It is a simple insertion-free sort
// over a small slice.
func sortCoinsDesc(coins []selectedCoin) {
	for i := 1; i < len(coins); i++ {
		for j := i; j > 0 && coins[j].value > coins[j-1].value; j-- {
			coins[j], coins[j-1] = coins[j-1], coins[j]
		}
	}
}

// outpointKey returns a stable map key for an outpoint: hex(txhash) + ":" +
// index.
func outpointKey(outpoint primitives.Outpoint) string {
	return hex.EncodeToString(outpoint.TxHash[:]) + ":" + core.Sprintf("%d", outpoint.Index)
}

// hashFromHex parses a 32-byte hash from the daemon's little-endian hex form into
// a primitives.Hash. HSD name-chain hashes are not reversed, so the bytes map
// directly.
func hashFromHex(value string) (primitives.Hash, error) {
	raw, err := hex.DecodeString(core.Trim(value))
	if err != nil {
		return primitives.Hash{}, core.E(opSelectCoin, "decode hash", err)
	}
	if len(raw) != 32 {
		return primitives.Hash{}, core.E(opSelectCoin, "hash must be 32 bytes", nil)
	}
	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}
