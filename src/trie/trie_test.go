package trie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/kvdb"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = false
}

// Used for testing
func newEmpty() *Trie {
	trie, _ := New(common.Hash{}, NewDatabase(kvdb.NewMemDatabase()))
	return trie
}

func TestEmptyTrie(t *testing.T) {
	var trie Trie
	res := trie.Hash()
	exp := emptyRoot
	if res != common.Hash(exp) {
		t.Errorf("expected %x got %x", exp, res)
	}
}

func TestNull(t *testing.T) {
	var trie Trie
	key := make([]byte, 32)
	value := []byte("test")
	trie.Update(key, value)
	if !bytes.Equal(trie.Get(key), value) {
		t.Fatal("wrong value")
	}
}

func TestMissingRoot(t *testing.T) {
	trie, err := New(common.HexToHash("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"),
		NewDatabase(kvdb.NewMemDatabase()))
	if trie != nil {
		t.Error("New returned non-nil trie for invalid root")
	}
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("New returned wrong error: %v", err)
	}
}

func TestMissingNodeDisk(t *testing.T)    { testMissingNode(t, false) }
func TestMissingNodeMemonly(t *testing.T) { testMissingNode(t, true) }

func testMissingNode(t *testing.T, memonly bool) {
	diskdb := kvdb.NewMemDatabase()
	triedb := NewDatabase(diskdb)

	trie, _ := New(common.Hash{}, triedb)
	updateString(trie, "120000", "qwerqwerqwerqwerqwerqwerqwerqwer")
	updateString(trie, "123456", "asdfasdfasdfasdfasdfasdfasdfasdf")
	root, _ := trie.Commit(nil)
	if !memonly {
		err := triedb.Commit(root, true)
		if err != nil {
			t.Fatal(err)
		}
	}

	trie, _ = New(root, triedb)
	_, err := trie.TryGet([]byte("120000"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120099"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcvzxcvzxcvzxcvzxcvzxcvzxcvzxcv"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryDelete([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	hash_ := common.HexToHash("2db7571d4823d62e2897c63e29e8360954fa43808feba89e85470fe9b2c10f83")

	if memonly {
		delete(triedb.dirties, hash_)
	} else {
		err = diskdb.Delete(hash_[:])
		if err != nil {
			t.Fatal(err)
		}
	}

	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120000"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120099"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcv"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryDelete([]byte("123456"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestInsert(t *testing.T) {
	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := common.HexToHash("e2310e8a1602e925553df6758ddd71feaf3e8cecedee87eb04c6ce496f68fe52")
	root := trie.Hash()
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}

	trie = newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	exp = common.HexToHash("b57a2c2e580634d1f28565817d4b90ea63521be9a2677fe94afa628bb5646946")
	root, err := trie.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}
}

func TestGet(t *testing.T) {
	trie := newEmpty()
	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	for i := 0; i < 2; i++ {
		res := getString(trie, "dog")
		if !bytes.Equal(res, []byte("puppy")) {
			t.Errorf("expected puppy got %x", res)
		}

		unknown := getString(trie, "unknown")
		if unknown != nil {
			t.Errorf("expected nil got %x", unknown)
		}

		if i == 1 {
			return
		}
		_, err := trie.Commit(nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDelete(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			updateString(trie, val.k, val.v)
		} else {
			deleteString(trie, val.k)
		}
	}

	h := trie.Hash()
	exp := common.HexToHash("8d7bf42c3e38c024b3e932fa8170bca210706a60a68d47d66411c1bb7bf4c891")
	if h != exp {
		t.Errorf("expected %x got %x", exp, h)
	}
}

func TestEmptyValues(t *testing.T) {
	trie := newEmpty()

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}

	h := trie.Hash()
	exp := common.HexToHash("8d7bf42c3e38c024b3e932fa8170bca210706a60a68d47d66411c1bb7bf4c891")
	if h != exp {
		t.Errorf("expected %x got %x", exp, h)
	}
}

func TestReplication(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}
	exp, err := trie.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}

	// create a new trie on top of the database and check that lookups work.
	trie2, err := New(exp, trie.db)
	if err != nil {
		t.Fatalf("can't recreate trie at %x: %v", exp, err)
	}
	for _, kv := range vals {
		if string(getString(trie2, kv.k)) != kv.v {
			t.Errorf("trie2 doesn't have %q => %q", kv.k, kv.v)
		}
	}
	trie2CommittedHash, err := trie2.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if trie2CommittedHash != exp {
		t.Errorf("root failure. expected %x got %x", exp, trie2CommittedHash)
	}

	// perform some insertions on the new trie.
	vals2 := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		// {"shaman", "horse"},
		// {"doge", "coin"},
		// {"ether", ""},
		// {"dog", "puppy"},
		// {"somethingveryoddindeedthis is", "myothernodedata"},
		// {"shaman", ""},
	}
	for _, val := range vals2 {
		updateString(trie2, val.k, val.v)
	}
	if trie2Hash := trie2.Hash(); trie2Hash != exp {
		t.Errorf("root failure. expected %x got %x", exp, trie2Hash)
	}
}

func TestLargeValue(t *testing.T) {
	trie := newEmpty()
	trie.Update([]byte("key1"), []byte{99, 99, 99, 99})
	trie.Update([]byte("key2"), bytes.Repeat([]byte{1}, 32))
	trie.Hash()
}

type countingDB struct {
	kvdb.Database
	gets map[string]int
}

func (db *countingDB) Get(key []byte) ([]byte, error) {
	db.gets[string(key)]++
	return db.Database.Get(key)
}

// TestCacheUnload checks that decoded nodes are unloaded after a
// certain number of commit operations.
func TestCacheUnload(t *testing.T) {
	// Create test trie with two branches.
	trie := newEmpty()
	key1 := "---------------------------------"
	key2 := "---some other branch"
	updateString(trie, key1, "this is the branch of key1.")
	updateString(trie, key2, "this is the branch of key2.")

	root, err := trie.Commit(nil)
	if err != nil {
		t.Fatal(err)
	}
	err = trie.db.Commit(root, true)
	if err != nil {
		t.Fatal(err)
	}

	// Commit the trie repeatedly and access key1.
	// The branch containing it is loaded from DB exactly two times:
	// in the 0th and 6th iteration.
	diskdb := &countingDB{Database: trie.db.diskdb, gets: make(map[string]int)}
	triedb := NewDatabase(diskdb)
	trie, _ = New(root, triedb)
	trie.SetCacheLimit(5)
	for i := 0; i < 12; i++ {
		getString(trie, key1)
		_, err = trie.Commit(nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	// Check that it got loaded two times.
	for dbkey, count := range diskdb.gets {
		if count != 2 {
			t.Errorf("db key %x loaded %d times, want %d times", []byte(dbkey), count, 2)
		}
	}
}

// randTest performs random trie operations.
// Instances of this test are created by Generate.
type randTest []randTestStep

type randTestStep struct {
	op    int
	key   []byte // for opUpdate, opDelete, opGet
	value []byte // for opUpdate
	err   error  // for debugging
}

const (
	opUpdate = iota
	opDelete
	opGet
	opCommit
	opHash
	opReset
	opItercheckhash
	opCheckCacheInvariant
	opMax // boundary value, not an actual op
)

func (randTest) Generate(r *rand.Rand, size int) reflect.Value {
	var allKeys [][]byte
	genKey := func() []byte {
		if len(allKeys) < 2 || r.Intn(100) < 10 {
			// new key
			key := make([]byte, r.Intn(50))
			r.Read(key)
			allKeys = append(allKeys, key)
			return key
		}
		// use existing key
		return allKeys[r.Intn(len(allKeys))]
	}

	var steps randTest
	for i := 0; i < size; i++ {
		step := randTestStep{op: r.Intn(opMax)}
		switch step.op {
		case opUpdate:
			step.key = genKey()
			step.value = make([]byte, 8)
			binary.BigEndian.PutUint64(step.value, uint64(i))
		case opGet, opDelete:
			step.key = genKey()
		}
		steps = append(steps, step)
	}
	return reflect.ValueOf(steps)
}

func runRandTest(rt randTest) bool {
	triedb := NewDatabase(kvdb.NewMemDatabase())

	tr, _ := New(common.Hash{}, triedb)
	values := make(map[string]string) // tracks content of the trie

	for i, step := range rt {
		switch step.op {
		case opUpdate:
			tr.Update(step.key, step.value)
			values[string(step.key)] = string(step.value)
		case opDelete:
			tr.Delete(step.key)
			delete(values, string(step.key))
		case opGet:
			v := tr.Get(step.key)
			want := values[string(step.key)]
			if string(v) != want {
				rt[i].err = fmt.Errorf("mismatch for key 0x%x, got 0x%x want 0x%x", step.key, v, want)
			}
		case opCommit:
			_, rt[i].err = tr.Commit(nil)
		case opHash:
			tr.Hash()
		case opReset:
			hash_, err := tr.Commit(nil)
			if err != nil {
				rt[i].err = err
				return false
			}
			newtr, err := New(hash_, triedb)
			if err != nil {
				rt[i].err = err
				return false
			}
			tr = newtr
		case opItercheckhash:
			checktr, _ := New(common.Hash{}, triedb)
			it := NewIterator(tr.NodeIterator(nil))
			for it.Next() {
				checktr.Update(it.Key, it.Value)
			}
			if tr.Hash() != checktr.Hash() {
				rt[i].err = fmt.Errorf("hash_ mismatch in opItercheckhash")
			}
		case opCheckCacheInvariant:
			rt[i].err = checkCacheInvariant(tr.root, nil, tr.cachegen, false, 0)
		}
		// Abort the test on error.
		if rt[i].err != nil {
			return false
		}
	}
	return true
}

func checkCacheInvariant(n, parent node, parentCachegen uint16, parentDirty bool, depth int) error {
	var children []node
	var flag nodeFlag
	switch n := n.(type) {
	case *shortNode:
		flag = n.flags
		children = []node{n.Val}
	case *fullNode:
		flag = n.flags
		children = n.Children[:]
	default:
		return nil
	}

	errorf := func(format string, args ...interface{}) error {
		msg := fmt.Sprintf(format, args...)
		msg += fmt.Sprintf("\nat depth %d node %s", depth, spew.Sdump(n))
		msg += fmt.Sprintf("parent: %s", spew.Sdump(parent))
		return errors.New(msg)
	}
	if flag.gen > parentCachegen {
		return errorf("cache invariant violation: %d > %d\n", flag.gen, parentCachegen)
	}
	if depth > 0 && !parentDirty && flag.dirty {
		return errorf("cache invariant violation: %d > %d\n", flag.gen, parentCachegen)
	}
	for _, child := range children {
		if err := checkCacheInvariant(child, n, flag.gen, flag.dirty, depth+1); err != nil {
			return err
		}
	}
	return nil
}

func TestRandom(t *testing.T) {
	if err := quick.Check(runRandTest, nil); err != nil {
		if cerr, ok := err.(*quick.CheckError); ok {
			t.Fatalf("random test iteration %d failed: %s", cerr.Count, spew.Sdump(cerr.In))
		}
		t.Fatal(err)
	}
}

func BenchmarkGet(b *testing.B)      { benchGet(b) }
func BenchmarkUpdateBE(b *testing.B) { benchUpdate(b, binary.BigEndian) }
func BenchmarkUpdateLE(b *testing.B) { benchUpdate(b, binary.LittleEndian) }

const benchElemCount = 20000

func benchGet(b *testing.B) {
	trie := new(Trie)
	k := make([]byte, 32)
	for i := 0; i < benchElemCount; i++ {
		binary.LittleEndian.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	binary.LittleEndian.PutUint64(k, benchElemCount/2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Get(k)
	}
	b.StopTimer()
}

func benchUpdate(b *testing.B, e binary.ByteOrder) *Trie {
	trie := newEmpty()
	k := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		e.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	return trie
}

// Benchmarks the trie hashing. Since the trie caches the result of any operation,
// we cannot use b.N as the number of hashing rounds, since all rounds apart from
// the first one will be NOOP. As such, we'll use b.N as the number of account to
// insert into the trie before measuring the hashing.
func BenchmarkHash(b *testing.B) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))

	// Create a realistic account trie to hash
	addresses := make([][20]byte, b.N)
	for i := 0; i < len(addresses); i++ {
		for j := 0; j < len(addresses[i]); j++ {
			addresses[i][j] = byte(random.Intn(256))
		}
	}
	accounts := make([][]byte, len(addresses))
	for i := 0; i < len(accounts); i++ {
		var (
			key  = crypto.Keccak256(nil)
			root = emptyRoot
		)
		accounts[i], _ = EncodeToBytes(&shortNode{
			Key: key,
			Val: &shortNode{
				Key: root.Bytes(),
			},
		})
	}
	// Insert the accounts into the trie and hash it
	trie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	trie.Hash()
}

func getString(trie *Trie, k string) []byte {
	return trie.Get([]byte(k))
}

func updateString(trie *Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}

func deleteString(trie *Trie, k string) {
	trie.Delete([]byte(k))
}
