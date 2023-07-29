package disruptor

import (
	"runtime"
	"sync/atomic"
)

const (
	STATERUNNING = iota
	STATECLOSED
)

type Consumer struct {
	current        *Sequence // this consumer has processed upto this sequence/cursor/index.
	upstream       Barrier   // this ring buffer has been writtern up to this sequence or producer current sequence/index/cursor.
	waiter         Waiter
	state          int64
	callerConsumer CustomConsumer //this will given by disruptor package caller, as caller will implementing the Consume(lower,upper) implementation.
}

func NewConsumer(current *Sequence, upstream Barrier, callerConsumer CustomConsumer, waiter Waiter) *Consumer {
	return &Consumer{
		current:        current,
		upstream:       upstream,
		state:          STATERUNNING,
		callerConsumer: callerConsumer,
		waiter:         waiter,
	}
}
func (c *Consumer) Read() {
	var lower, upper int64
	var current = c.current.Load()
	for {
		lower = current + 1
		upper = c.upstream.Load()
		if lower <= upper { //Consumer is before producer.
			c.callerConsumer.Consume(lower, upper)
			c.current.Store(upper)
			current = upper
		} else if upper = c.upstream.Load(); lower <= upper {
			runtime.Gosched() //avoiding continour lookup to producer/upstream cursor as sequence take lock on CPU level, this will help reducing CPU cycle.
		} else if atomic.LoadInt64(&c.state) == STATERUNNING { //checking if consumer is in running state.
			c.waiter.Idle()
		} else { //consumer is stoped/closed.
			break
		}
	}
}
func (c *Consumer) Close() {
	atomic.StoreInt64(&c.state, STATECLOSED)
}
