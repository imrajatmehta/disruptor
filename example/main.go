package main

import (
	"sync"

	"github.com/imrajatmehta/disruptor"
)

func main() {
	newDisruptor := disruptor.New(
		disruptor.WithCapacity(BUFFERSIZE),
		disruptor.WithConsumerGroups(NewConsumer{}, NewConsumer{}))
	go publish(newDisruptor)
	read(newDisruptor)

}
func publish(newD disruptor.Disruptor) {

	for sequence := int64(0); sequence < 1_000_000_00; {
		sequence = newD.Reserve(RESERVATION)
		// if RESERVATION is high then it will try to insert in batches, as we will getting the sequence upto which we can insert from Reserve()
		for lower := sequence - RESERVATION + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BUFFEMASK] = lower
		}

		newD.Commit(sequence-RESERVATION+1, sequence) //Batch commit
	}
	for _, reader := range newD.Readers {
		reader.Close()
	}

}
func read(newDisruptor disruptor.Disruptor) {
	wg := sync.WaitGroup{}
	wg.Add(len(newDisruptor.Readers))
	for _, reader := range newDisruptor.Readers {
		go func(obj disruptor.Reader) {
			obj.Read()
			wg.Done()
		}(reader)
	}
	wg.Wait()
}

type NewConsumer struct{}

func (obj NewConsumer) Consume(lower, upper int64) {
	for ; lower <= upper; lower++ {
		message := ringBuffer[lower&BUFFEMASK]

		if message != lower {
			panic("race condition occured")
		}
	}
}

const (
	BUFFERSIZE  = 1024 * 64
	BUFFEMASK   = BUFFERSIZE - 1
	RESERVATION = 1
)

var ringBuffer = [BUFFERSIZE]int64{}
