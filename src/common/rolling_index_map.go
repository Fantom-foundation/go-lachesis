package common

import (
	"fmt"
)

type RollingIndexMap struct {
	name    string
	size    int
	keys    []uint32
	mapping map[uint32]*RollingIndex
}

func NewRollingIndexMap(name string, size int, keys []uint32) *RollingIndexMap {
	items := make(map[uint32]*RollingIndex)
	for _, key := range keys {
		items[key] = NewRollingIndex(fmt.Sprintf("%s[%d]", name, key), size)
	}
	return &RollingIndexMap{
		name:    name,
		size:    size,
		keys:    keys,
		mapping: items,
	}
}

func (rim *RollingIndexMap) AddKey(key uint32) error {
	if _, ok := rim.mapping[key]; ok {
		return NewStoreErr(rim.name, KeyAlreadyExists, fmt.Sprintf("%d", key))
	}
	rim.keys = append(rim.keys, key)
	rim.mapping[key] = NewRollingIndex(fmt.Sprintf("%s[%d]", rim.name, key), rim.size)
	return nil
}

//return key items with index > skip
func (rim *RollingIndexMap) Get(key uint32, skipIndex int64) ([]interface{}, error) {
	items, ok := rim.mapping[key]
	if !ok {
		return nil, NewStoreErr(rim.name, KeyNotFound, fmt.Sprint(key))
	}

	cached, err := items.Get(skipIndex)
	if err != nil {
		return nil, err
	}

	return cached, nil
}

func (rim *RollingIndexMap) GetItem(key uint32, index int64) (interface{}, error) {
	return rim.mapping[key].GetItem(index)
}

func (rim *RollingIndexMap) GetLast(key uint32) (interface{}, error) {
	pe, ok := rim.mapping[key]
	if !ok {
		return nil, NewStoreErr(rim.name, KeyNotFound, fmt.Sprint(key))
	}
	cached, _ := pe.GetLastWindow()
	if len(cached) == 0 {
		return "", NewStoreErr(rim.name, Empty, "")
	}
	return cached[len(cached)-1], nil
}

func (rim *RollingIndexMap) Set(key uint32, item interface{}, index int64) error {
	items, ok := rim.mapping[key]
	if !ok {
		items = NewRollingIndex(fmt.Sprintf("%s[%d]", rim.name, key), rim.size)
		rim.mapping[key] = items
	}
	return items.Set(item, index)
}

//returns [key] => lastKnownIndex
func (rim *RollingIndexMap) Known() map[uint32]int64 {
	known := make(map[uint32]int64)
	for k, items := range rim.mapping {
		_, lastIndex := items.GetLastWindow()
		known[k] = lastIndex
	}
	return known
}

func (rim *RollingIndexMap) Reset() error {
	items := make(map[uint32]*RollingIndex)
	for _, key := range rim.keys {
		items[key] = NewRollingIndex(fmt.Sprintf("%s[%d]", rim.name, key), rim.size)
	}
	rim.mapping = items
	return nil
}

func (rim *RollingIndexMap) Import(other *RollingIndexMap) {
	for _, key := range other.keys {
		rim.mapping[key] = NewRollingIndex(fmt.Sprintf("%s[%d]", rim.name, key), rim.size)
 		rim.mapping[key].lastIndex = other.mapping[key].lastIndex
		rim.mapping[key].items = other.mapping[key].items
 		// copy(rim.mapping[key].items[:len(other.mapping[key].items)], other.mapping[key].items)
	}
}
