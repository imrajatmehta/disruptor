package disruptor

import "math"

type compositeConsumerBarrier []*Sequence

func NewCompositeConsumerBarrier(sequences ...*Sequence) Barrier {
	return compositeConsumerBarrier(sequences)
}
func (sequences compositeConsumerBarrier) Load() int64 {
	var min int64 = math.MaxInt64
	for i := range sequences {
		if seqIndex := sequences[i].Load(); seqIndex < min {
			min = seqIndex
		}
	}
	return min
}
