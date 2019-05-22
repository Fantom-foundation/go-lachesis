package poset

type CountMap map[EventHash]map[int64]uint64

// NewCountMap creates new empty CountMap
func NewCountMap() CountMap {
	return CountMap(make(map[EventHash]map[int64]uint64))
}

func (cm *CountMap) Inc(hash EventHash, lamportTime int64) {
	_, exists := (*cm)[hash]
	if exists {
		cval, exists2 := (*cm)[hash][lamportTime]
		if exists2 {
			(*cm)[hash][lamportTime] = 1 + cval
		} else {
			(*cm)[hash][lamportTime] = 1
		}
	} else {
		(*cm)[hash] = make(map[int64]uint64)
		(*cm)[hash][lamportTime] = 1
	}
}
