package disruptor

import "sync/atomic"

type Sequence [8]int64 //TODO: Add comments

func NewSequence() *Sequence {
	seq := &Sequence{}
	seq[0] = defaultSequence
	return seq
}

func (this *Sequence) Load() int64 {
	return atomic.LoadInt64(&this[0])
}

func (this *Sequence) Store(value int64) {
	atomic.StoreInt64(&this[0], value)
}

var defaultSequence int64 = -1
