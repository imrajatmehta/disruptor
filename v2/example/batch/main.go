package main

import (
	"fmt"

	"github.com/imrajatmehta/disruptor/v2"
)

func main() {
	d := disruptor.New(
		disruptor.WithCapacity(BUFFERSIZE),
		disruptor.WithBatchSize(RESERVATION),
		disruptor.WithConsumerGroups(NewConsumer{}, NewConsumer{}))
	go func() {
		for seq := int64(0); seq+RESERVATION < 1_000; seq += RESERVATION {
			d.PublishBatch([]int64{seq + 1, seq + 2, seq + 3, seq + 4, seq + 5, seq + 6, seq + 7, seq + 8, seq + 9, seq + 10, seq + 11, seq + 12, seq + 13, seq + 14, seq + 15, seq + 16})
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
	RESERVATION = 16
)
