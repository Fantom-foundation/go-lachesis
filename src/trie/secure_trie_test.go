package trie

import (
	"context"
	"bytes"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/kvdb"
	"github.com/fortytw2/leaktest"
)

func newEmptySecure() *SecureTrie {
	trie, _ := NewSecure(hash.Hash{}, NewDatabase(kvdb.NewMemDatabase()), 0)
	return trie
}

// makeTestSecureTrie creates a large enough secure trie for testing.
func makeTestSecureTrie(t *testing.T) (*Database, *SecureTrie, map[string][]byte) {
	// Create an empty trie
	triedb := NewDatabase(kvdb.NewMemDatabase())

	trie, _ := NewSecure(hash.Hash{}, triedb, 0)

	// Fill it with some arbitrary data
	content := make(map[string][]byte)
	for i := byte(0); i < 255; i++ {
		// Map the same data under multiple keys
		key, val := common.LeftPadBytes([]byte{1, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		key, val = common.LeftPadBytes([]byte{2, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		// Add some other data to inflate the trie
		for j := byte(3); j < 13; j++ {
			key, val = common.LeftPadBytes([]byte{j, i}, 32), []byte{j, i}
			content[string(key)] = val
			trie.Update(key, val)
		}
	}
	_, err := trie.Commit(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Return the generated trie
	return triedb, trie, content
}

func TestSecureDelete(t *testing.T) {
	trie := newEmptySecure()
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
			trie.Update([]byte(val.k), []byte(val.v))
		} else {
			trie.Delete([]byte(val.k))
		}
	}

	h := trie.Hash()
	exp := hash.HexToHash("d42f28391c76efa9d3868aae7d0bf7c7010527d93bd19bf23a0dd95b9b31c7a6")
	if h != exp {
		t.Errorf("expected %x got %x", exp, h)
	}
}

func TestSecureGetKey(t *testing.T) {
	trie := newEmptySecure()
	trie.Update([]byte("foo"), []byte("bar"))

	key := []byte("foo")
	value := []byte("bar")
	seckey := crypto.Keccak256(key)

	if !bytes.Equal(trie.Get(key), value) {
		t.Errorf("Get did not return bar")
	}
	if k := trie.GetKey(seckey); !bytes.Equal(k, key) {
		t.Errorf("GetKey returned %q, want %q", k, key)
	}
}

func TestSecureTrieConcurrency(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    defer leaktest.CheckContext(ctx, t)()

	// Create an initial trie and copy if for concurrent access
	_, trie, _ := makeTestSecureTrie(t)

	threads := runtime.NumCPU()
	tries := make([]*SecureTrie, threads)
	for i := 0; i < threads; i++ {
		cpy := *trie
		tries[i] = &cpy
	}
	// Start a batch of goroutines interacting with the trie
	errs := make(chan error)
	pend := new(sync.WaitGroup)
	pend.Add(threads)
	go func() {
		pend.Wait()
		close(errs)
	}()
	for i := 0; i < threads; i++ {
		go func(index int) {
			defer pend.Done()

			for j := byte(0); j < 255; j++ {
				// Map the same data under multiple keys
				key, val := common.LeftPadBytes([]byte{byte(index), 1, j}, 32), []byte{j}
				tries[index].Update(key, val)

				key, val = common.LeftPadBytes([]byte{byte(index), 2, j}, 32), []byte{j}
				tries[index].Update(key, val)

				// Add some other data to inflate the trie
				for k := byte(3); k < 13; k++ {
					key, val = common.LeftPadBytes([]byte{byte(index), k, j}, 32), []byte{k, j}
					tries[index].Update(key, val)
				}
			}
			_, err := tries[index].Commit(nil)
			if err != nil {
				errs <- err
			}
		}(i)
	}
	// Wait for all threads to finish
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}
