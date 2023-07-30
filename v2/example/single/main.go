package main

import (
	"fmt"

	"github.com/imrajatmehta/disruptor/v2"
)

func main() {
	d, _ := disruptor.New(
		disruptor.WithCapacity(BUFFERSIZE),
		disruptor.WithBatchSize(RESERVATION),
		disruptor.WithConsumerGroups(NewConsumer{}, NewConsumer{}))
	go func() {
		for i := int64(0); i < 1_000; i++ {
			d.Publish(i)
		}
		d.Close()
	}()
	d.Read()
}

type NewConsumer struct{}

func (obj NewConsumer) OnEvent(msg int64) {
	fmt.Println(msg)
}

const (
	BUFFERSIZE  = 1024 * 64
	RESERVATION = 1
)
