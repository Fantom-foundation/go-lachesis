package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/src/kvdb"
)

func TestBalanceState(t *testing.T) {
	assertar := assert.New(t)

	var aa = []common.Address{
		hash.FakePeer(),
		hash.FakePeer(),
		hash.FakePeer(),
	}

	mem := kvdb.NewMemDatabase()
	store := NewDatabase(mem)

	stateAt := func(point common.Hash) *DB {
		db, err := New(point, store)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return db
	}

	checkBalance := func(point common.Hash, addr common.Address, balance pos.Stake) {
		db := stateAt(point)
		got := db.FreeBalance(addr)
		if !assertar.Equalf(balance, got, "unexpected balance") {
			t.FailNow()
		}
	}

	commit := func(db *DB) common.Hash {
		root, err := db.Commit(true)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return root
	}

	// empty
	for _, a := range aa {
		checkBalance(common.Hash{}, a, 0)
	}

	// root
	db := stateAt(common.Hash{})
	db.SetBalance(aa[0], 10)
	db.SetBalance(aa[1], 10)
	db.SetBalance(aa[2], 10)
	root := commit(db)

	checkBalance(root, aa[0], 10)
	checkBalance(root, aa[1], 10)
	checkBalance(root, aa[2], 10)

	// fork 1
	db = stateAt(root)
	db.Transfer(aa[0], aa[1], 1)
	if !assertar.Equalf(pos.Stake(9), db.FreeBalance(aa[0]), "before commit") ||
		!assertar.Equalf(pos.Stake(11), db.FreeBalance(aa[1]), "before commit") {
		return
	}
	fork1 := commit(db)

	checkBalance(fork1, aa[0], 9)
	checkBalance(fork1, aa[1], 11)
	checkBalance(fork1, aa[2], 10)

	// fork 2
	db = stateAt(root)
	db.Transfer(aa[0], aa[2], 5)
	fork2 := commit(db)

	checkBalance(fork2, aa[0], 5)
	checkBalance(fork2, aa[1], 10)
	checkBalance(fork2, aa[2], 15)
}

func TestDelegationState(t *testing.T) {
	assertar := assert.New(t)

	const __ uint64 = 0

	var aa = []common.Address{
		hash.FakePeer(),
		hash.FakePeer(),
		hash.FakePeer(),
	}

	mem := kvdb.NewMemDatabase()
	store := NewDatabase(mem)

	stateAt := func(point common.Hash) *DB {
		db, err := New(point, store)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return db
	}

	check := func(x direction, point common.Hash, addr common.Address, amount ...uint64) {
		db := stateAt(point)

		var j int
		for j = range aa {
			if aa[j] == addr {
				break
			}
		}

		var dir string
		if x == TO {
			dir = "-->"
		} else {
			dir = "<--"
		}

		for i, exp := range amount {
			p := aa[i]
			if p == addr {
				continue
			}
			got := db.GetDelegations(addr)[x][p]
			if !assertar.Equalf(exp, got, "unexpected delegation amount: aa[%d] %s aa[%d]", j, dir, i) {
				t.FailNow()
			}
		}
	}

	commit := func(db *DB) common.Hash {
		root, err := db.Commit(true)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return root
	}

	// step 0
	db := stateAt(common.Hash{})
	db.SetBalance(aa[0], 100)
	db.SetBalance(aa[1], 100)
	db.SetBalance(aa[2], 100)
	root := commit(db)

	// step 1
	db = stateAt(root)
	db.Delegate(aa[1], aa[0], 10, 1)
	db.Delegate(aa[1], aa[2], 20, 1)
	root = commit(db)

	check(TO, root, aa[0], __, 00, 00)
	check(TO, root, aa[1], 10, __, 20)
	check(TO, root, aa[2], 00, 00, __)

	check(FROM, root, aa[0], __, 10, 00)
	check(FROM, root, aa[1], 00, __, 00)
	check(FROM, root, aa[2], 00, 20, __)

	// step 2
	db = stateAt(root)
	db.Delegate(aa[0], aa[1], 15, 2)
	db.Delegate(aa[1], aa[2], 25, 2)
	root = commit(db)

	check(TO, root, aa[0], __, 15, 00)
	check(TO, root, aa[1], 10, __, 45)
	check(TO, root, aa[2], 00, 00, __)

	check(FROM, root, aa[0], __, 10, 00)
	check(FROM, root, aa[1], 15, __, 00)
	check(FROM, root, aa[2], 00, 45, __)

	// step 3
	db = stateAt(root)
	db.ExpireDelegations(aa[0], 1)
	db.ExpireDelegations(aa[1], 1)
	db.ExpireDelegations(aa[2], 1)
	root = commit(db)

	check(TO, root, aa[0], __, 15, 00)
	check(TO, root, aa[1], 00, __, 25)
	check(TO, root, aa[2], 00, 00, __)

	check(FROM, root, aa[0], __, 00, 00)
	check(FROM, root, aa[1], 15, __, 00)
	check(FROM, root, aa[2], 00, 25, __)
}

func TestIdempotency(t *testing.T) {
	assertar := assert.New(t)

	mem := kvdb.NewMemDatabase()
	store := NewDatabase(mem)

	stateAt := func(point common.Hash) *DB {
		db, err := New(point, store)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return db
	}

	checkBalance := func(point common.Hash, addr common.Address, balance pos.Stake) {
		db := stateAt(point)
		got := db.FreeBalance(addr)
		if !assertar.Equalf(balance, got, "unexpected balance") {
			t.FailNow()
		}
	}

	commit := func(db *DB) common.Hash {
		root, err := db.Commit(true)
		if !assertar.NoError(err) {
			t.FailNow()
		}
		return root
	}

	aa := hash.FakePeer()

	db := stateAt(common.Hash{})
	db.SetBalance(aa, 10)
	root1 := commit(db)

	db = stateAt(common.Hash{})
	db.SetBalance(aa, 10)
	root2 := commit(db)

	if !assertar.Equal(root1, root2) {
		return
	}

	checkBalance(root1, aa, 10)
	checkBalance(root2, aa, 10)
}
