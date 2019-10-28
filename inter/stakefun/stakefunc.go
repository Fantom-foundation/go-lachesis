package pos

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

// this class provides several implementations
// to compute a stake value from an account's balance

// account balance. Can be positive, zero or negative (err)
type Balance int64

// the stake value computed from account balance
type Stake float64

// main interface
type BalanceToStake interface {
	func GetStake(balance Balance) Stake
	func GetFuncName() string
}

// a 1:1 calculation of balance to stake
type BalanceToStakeSame struct {
	BalanceToStake
}

func (b2s *BalanceToStakeSame) GetStake(balance Balance) Stake {
	return balance
}

func (b2s *BalanceToStakeSame) GetFuncName() string {
	return "BalanceToStakeSame"
}

// the accummulated Saga points
type Saga uint64

// Calculation stake from balance and Saga points
type BalanceToStake_Saga struct {
	// the accummulated Saga points of an account
	Sp Saga

	BalanceToStake
}

func NewBalanceToStake_Saga(s Saga) {
	return &BalanceToStake_Saga {
		Sp : s
	}
}

func (b2s *BalanceToStake_Saga) GetStake(balance Balance) Stake {
	return balance * Sp
}

func (b2s *BalanceToStake_Saga) GetFuncName() string {
	return "BalanceToStake_Saga"
}