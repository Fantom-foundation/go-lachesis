package client

import (
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
)

// Balance type
type Balance big.Int

// Stake type
type Stake big.Int

// Member is a member interface.
type Member interface {
	// GetAccount returns the underlying Account
	GetAccount() *genesis.Account
	// GetBalance returns the amount of tokens in the account
	GetBalance() Balance
	// GetStake returns the stake calculated from the account's balance
	GetStake() Stake
	// IsValid returns whether this member is not banned.
	IsValid() bool
	// IsBanned returns whether this member has been banned in this epoch
	IsBanned() bool
	// GetLastEpoch returns the last epoch that this member was participating and not banned in the epoch
	GetLastEpoch()
}