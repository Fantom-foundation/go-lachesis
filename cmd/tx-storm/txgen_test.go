package main

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	t.Skip("example only")
	cfg := &Config{
		ChainId: 999,
		Accs: struct {
			Count  uint
			Offset uint
		}{
			Count:  10,
			Offset: 100,
		},
	}
	g := NewTxGenerator(cfg, 1, 1)
	for i := 0; i < 2*len(g.accs); i++ {
		tx := g.Yield()
		t.Log(tx.Info.String(), tx.Raw.Nonce(), tx.Raw.Value())
	}
}
