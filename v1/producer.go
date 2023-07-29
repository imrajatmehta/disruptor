package disruptor

import "runtime"

type Producer struct {
	current  *Sequence
	upstream Barrier
	previous int64
	capacity int64
	gate     int64
}

func NewProducer(current *Sequence, upstream Barrier, capacity int64) *Producer {
	return &Producer{
		current:  current,
		upstream: upstream,
		previous: defaultSequence,
		capacity: capacity,
		gate:     defaultSequence,
	}
}
func (p *Producer) Reserve(count int64) int64 {
	if count <= 0 {
		panic(errMinimumReservationSize)
	}
	p.previous += count //NOTE: mutual exlusion is not handled for this variable when called Reserve() parrallely.
	for spin := int64(0); p.gate < p.previous-p.capacity; spin++ {
		if spin&SPINMASK == 0 { //TO minimise the goSchedule
			runtime.Gosched()
		}
		p.gate = p.upstream.Load()
	}
	return p.previous
}
func (p *Producer) Commit(_, upper int64) {
	p.current.Store(upper)
}

const SPINMASK = 1024*16 - 1
