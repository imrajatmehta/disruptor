package disruptor

import "time"

type DefaultWaitStrategy struct{}

func NewWaitStrategy() DefaultWaitStrategy { return DefaultWaitStrategy{} }
func (this DefaultWaitStrategy) Gate() {
	time.Sleep(time.Nanosecond)
}
func (this DefaultWaitStrategy) Idle() {
	time.Sleep(time.Microsecond * 50)
}
