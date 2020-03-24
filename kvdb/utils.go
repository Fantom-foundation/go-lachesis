package kvdb

// Move data from src to dst.
func Move(src, dst KeyValueStore, prefix []byte) (err error) {
	keys := make([][]byte, 0, 5000) // don't write during iteration

	it := src.NewIteratorWithPrefix(prefix)
	defer it.Release()

	for it.Next() {
		err = dst.Put(it.Key(), it.Value())
		if err != nil {
			return
		}
		keys = append(keys, it.Key())
	}

	err = it.Error()
	if err != nil {
		return
	}

	for _, key := range keys {
		err = src.Delete(key)
		if err != nil {
			return
		}
	}

	return nil
}

// Copy data from src to dst.
func Copy(src, dst KeyValueStore, prefix []byte) (err error) {
	it := src.NewIteratorWithPrefix(prefix)
	defer it.Release()

	for it.Next() {
		err = dst.Put(it.Key(), it.Value())
		if err != nil {
			return
		}
	}

	err = it.Error()
	if err != nil {
		return
	}

	return nil
}
