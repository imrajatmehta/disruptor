package disruptor

import "sync/atomic"

type Sequence [8]int64 //TODO: Add comments

var defaultSequenceValue int64 = -1

func NewSequence() *Sequence {
	seq := Sequence{}
	seq[0] = defaultSequenceValue
	return &seq
}
func (cur *Sequence) Get() int64 {
	return atomic.LoadInt64(&cur[0])
}
func (cur *Sequence) Set(data int64) {
	atomic.StoreInt64(&cur[0], data)
}
