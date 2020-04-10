package basiccheck

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/eventcheck/tests"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
)

// TestBasicCheck is a main testing func
func TestBasicCheck(t *testing.T) {
	lachesisConfigs := tests.DagConfigs()

	for _, lachesisConfig := range lachesisConfigs {
		checker := New(lachesisConfig)
		testValidateTx(t, checker)
		testCheckTxs(t, checker)
		testCheckerValidate(t, checker)
	}
}

// testCheckerValidate runs tests for Validate func
func testCheckerValidate(t *testing.T, c *Checker) {
	testEvents := tests.Events()
	for _, event := range testEvents {
		if event != nil && event.GasPowerUsed == 0 && c != nil && c.config != nil {
			event.GasPowerUsed = CalcGasPowerUsed(event, c.config)
			//log.Println("event.GasPowerUsed", event.GasPowerUsed)
		}

		err := c.Validate(event)
		checkValidateTxError(t, event, err, c)
	}
}

// checkValidateTxError checks that returned error is expected
func checkValidateTxError(t *testing.T, event *inter.Event, err error, c *Checker) {
	if event == nil {
		require.Equal(t, err, ErrUninitialisedEvent)
		return
	}
	if event.Version != 0 {
		require.Equal(t, err, ErrVersion)
		return
	}
	if len(event.Extra) > params.MaxExtraData {
		require.Equal(t, err, ErrExtraTooLarge)
		return
	}
	if len(event.Parents) > c.config.MaxParents {
		require.Equal(t, err, ErrTooManyParents)
		return
	}
	if event.Seq >= math.MaxInt32/2 || event.Epoch >= math.MaxInt32/2 || event.Frame >= math.MaxInt32/2 ||
		event.Lamport >= math.MaxInt32/2 || event.GasPowerUsed >= math.MaxInt64/2 || event.GasPowerLeft.Max() >= math.MaxInt64/2 {
		require.Equal(t, err, ErrHugeValue)
		return
	}
	if event.Seq == 0 || event.Epoch == 0 || event.Frame == 0 || event.Lamport == 0 {
		require.Equal(t, err, ErrNotInited)
		return
	}

	if event.ClaimedTime == 0 {
		require.Equal(t, err, ErrZeroTime)
		return
	}
	if event.Seq > 1 && len(event.Parents) == 0 {
		require.Equal(t, err, ErrNoParents)
		return
	}
	if len(event.Sig) != validSigLength {
		require.Equal(t, err, ErrSigMalformed)
		return
	}
	if event.GasPowerUsed > params.MaxGasPowerUsed {
		require.Equal(t, err, ErrTooBigGasUsed)
		return
	}
	if event.GasPowerUsed != CalcGasPowerUsed(event, c.config) {
		require.Equal(t, err, ErrWrongGasUsed)
		return
	}
}

// testCheckTxs runs checkTxs func of a checker and validates result
func testCheckTxs(t *testing.T, c *Checker) {
	event := inter.Event{
		Transactions: tests.Transactions(),
	}
	err := c.checkTxs(&event)
	require.NotNil(t, err)

	event = inter.Event{
		Transactions: tests.ValidTransactions(),
	}
	err = c.checkTxs(&event)
	require.Nil(t, err)
}

// testValidateTx runs validateTx func of a checker and validates result
func testValidateTx(t *testing.T, c *Checker) {
	transactions := tests.Transactions()
	for _, tx := range transactions {
		err := c.validateTx(tx)
		handleCheckTxError(t, tx, err)
	}
}

// handleCheckTxError checks that returned error is expected
func handleCheckTxError(t *testing.T, tx *types.Transaction, returnedError error) {
	if tx.Value().Sign() < 0 || tx.GasPrice().Sign() < 0 {
		require.Equal(t, returnedError, ErrNegativeValue)
		return
	}

	intrGas, err2 := evmcore.IntrinsicGas(tx.Data(), tx.To() == nil, true)
	if err2 != nil {
		require.Equal(t, returnedError, err2)
		return
	}

	if tx.Gas() < intrGas {
		require.Equal(t, returnedError, ErrIntrinsicGas)
		return
	}
	if tx.GasPrice().Cmp(params.MinGasPrice) < 0 {
		require.Equal(t, returnedError, ErrUnderpriced)
		return
	}
}
