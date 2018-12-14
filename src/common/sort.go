package common

import "sort"

type Int64Slice []int64

func (a Int64Slice) Len() int      { return len(a) }
func (a Int64Slice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Int64Slice) Less(i, j int) bool {
	return a[i] < a[j]
}

// Sort is a convenience method.
func (a Int64Slice) Sort() { sort.Sort(a) }