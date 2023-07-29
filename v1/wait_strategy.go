package disruptor

import "time"

type WaitStrategy struct{}

func NewWaitStrategy() WaitStrategy {
	return WaitStrategy{}
}
func (w WaitStrategy) Idle() { time.Sleep(time.Microsecond * 50) }
