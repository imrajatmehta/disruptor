// Below code is implemented coz for a ease of producer.
// As producer will need to check the slow pointer from consumers before writting to buffer,
// if we dont implement it saperately then producer will need to write
// the below Load Fn logic in there side.
package disruptor

import "math"

type compositeConsumerBarrier []*Sequence

func NewCompositeConsumerBarrier(sequences ...*Sequence) Barrier {
	return compositeConsumerBarrier(sequences)
}

// Return the slow sequence/cursor/index from all consumers.
func (sequences compositeConsumerBarrier) Load() int64 {
	var min int64 = math.MaxInt64
	for i := range sequences {
		if seqIndex := sequences[i].Load(); seqIndex < min {
			min = seqIndex
		}
	}
	return min
}
