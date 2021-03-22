package gossip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func isEmptyDB(db ethdb.Iteratee) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	return !it.Next()
}

func (s *Store) Migrate() error {
	versions := migration.NewKvdbIDStore(s.table.Version)
	if isEmptyDB(s.mainDb) && isEmptyDB(s.async.mainDb) {
		// short circuit if empty DB
		versions.SetID(s.migrations().ID())
		return nil
	}
	return s.migrations().Exec(versions)
}

func (s *Store) migrations() *migration.Migration {
	return migration.
		Begin("lachesis-gossip-store").
		Next("remove async data from sync DBs",
			func() error {
				s.rmPrefix(s.table.PackInfos, "serverPool")
				s.rmPrefix(s.mainDb, "Z")
				return nil
			}).
		Next("remove legacy genesis field",
			s.app.MigrateEraseGenesisField).
		Next("multi-delegations",
			s.app.MigrateMultiDelegations).
		Next("adjustable offline pruning time",
			s.app.MigrateAdjustableOfflinePeriod).
		Next("adjustable minimum gas price",
			s.app.MigrateAdjustableMinGasPrice).
		Next("backported topicsdb from go-opera",
			func() error {
				commitIfNeeded := func() {
					_ = s.Commit(nil, false)
				}
				_ = s.app.DropTopicsdb(commitIfNeeded)
				i := 1
				s.ForEachBlock(func(blockIdx idx.Block, block *inter.Block) {
					receipts := s.app.GetReceipts(blockIdx)
					if len(receipts) != 0 {
						allTxs := s.GetBlockTransactions(block)
						logIdx := uint(0)
						for i, r := range receipts {
							for _, l := range r.Logs {
								l.BlockNumber = uint64(blockIdx)
								l.TxHash = allTxs[i].Hash()
								l.Index = logIdx
								l.TxIndex = uint(i)
								l.BlockHash = common.Hash(block.Atropos)
								logIdx++
							}
							s.app.IndexLogs(r.Logs...)
						}
					}
					if i%10000 == 0 {
						commitIfNeeded()
					}
					i++
				})
				return nil
			})
}
