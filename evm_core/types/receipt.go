package types

import (
	"github.com/Fantom-foundation/go-ethereum/core/types"
)

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

var (
	receiptStatusFailedRLP     = []byte{}
	receiptStatusSuccessfulRLP = []byte{0x01}
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)
)

// Receipt represents the results of a transaction.
type Receipt = types.Receipt

//type receiptMarshaling

// receiptRLP is the consensus encoding of a receipt.
//type receiptRLP

// storedReceiptRLP is the storage encoding of a receipt.
//type storedReceiptRLP

// v4StoredReceiptRLP is the storage encoding of a receipt used in database version 4.
//type v4StoredReceiptRLP

// v3StoredReceiptRLP is the original storage encoding of a receipt including some unnecessary fields.
//type v3StoredReceiptRLP


// ReceiptForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type ReceiptForStorage = types.ReceiptForStorage


// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts = types.Receipts
