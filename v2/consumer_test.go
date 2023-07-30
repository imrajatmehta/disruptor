package disruptor_test

import (
	"sync"
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

func TestConsumerCircular(t *testing.T) {
	consSeq := disruptor.NewSequence()
	prodSeq := disruptor.NewSequence()
	waiter := disruptor.NewWaitStrategy()
	cons := disruptor.NewConsumer(consSeq, disruptor.Barrier(prodSeq), NewTestConsumer{}, disruptor.Waiter(waiter))
	ringBuffer := make([]int64, 1024)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cons.Read(ringBuffer, int64(len(ringBuffer)-1))
		wg.Done()
	}()
	cons.Close()
	wg.Wait()
	if consSeq.Load() != -1 {
		t.Fatal("consumer is not at initialise place i.e -1")
	}

	prodSeq.Store(5)
	wg.Add(1)
	go func() {
		cons.Read(ringBuffer, int64(len(ringBuffer)-1))
		wg.Done()
	}()
	cons.Close()
	wg.Wait()
	if consSeq.Load() != 5 {
		t.Fatal("consumer is not at 5, as producer is just at 5 commited")
	}
	prodSeq.Store(1026)
	wg.Add(1)
	go func() {
		cons.Read(ringBuffer, int64(len(ringBuffer)-1))
		wg.Done()
	}()
	cons.Close()
	wg.Wait()
	if consSeq.Load() != 1026 {
		t.Fatal("consumer not able to reach at 1026 seq.")
	}
}
