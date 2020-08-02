package evmcore

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Transaction and additional fields
type Transaction struct {
	*types.Transaction
	batch struct {
		FirstTx common.Hash
		Count   uint
	}
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) (err error) {
	err = rlp.Encode(w, tx.Transaction)
	if err != nil {
		return
	}

	err = rlp.Encode(w, &tx.batch)
	return
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) (err error) {
	tx.Transaction = new(types.Transaction)
	err = s.Decode(tx.Transaction)
	if err != nil {
		return
	}

	err = s.Decode(&tx.batch)
	return
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

func transaction(tx *types.Transaction) *Transaction {
	return &Transaction{
		Transaction: tx,
	}
}

func transactions(txs []*types.Transaction) Transactions {
	res := make(Transactions, len(txs))
	for i, t := range txs {
		res[i] = transaction(t)
	}
	return res
}

func (txs Transactions) Transactions() types.Transactions {
	res := make(types.Transactions, len(txs))
	for i, t := range txs {
		res[i] = t.Transaction
	}
	return res
}

func (txs Transactions) Len() int {
	return len(txs)
}

// txByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type txByNonce Transactions

func (s txByNonce) Len() int           { return len(s) }
func (s txByNonce) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s txByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// txsDifference returns a new set which is the difference between a and b.
func txsDifference(a, b Transactions) Transactions {
	keep := make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}
