package utils

type RingBuff struct {
	buff []float64
	size uint64
	pos  uint64
}

func NewRingBuff(size int) *RingBuff {
	return &RingBuff{
		buff: make([]float64, size),
		size: uint64(size),
	}
}

func (b *RingBuff) Push(v float64) {
	i := b.pos % b.size
	b.buff[i] = v
	b.pos++
}

func (b *RingBuff) Avg() float64 {
	count := b.pos
	if count > b.size {
		count = b.size
	}

	var sum float64
	for i := uint64(0); i < count; i++ {
		sum += b.buff[i]
	}

	return sum / float64(count)
}
