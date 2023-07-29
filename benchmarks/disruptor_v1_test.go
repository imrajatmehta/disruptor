package benchmarks

import (
	"log"
	"sync"
	"testing"

	"github.com/imrajatmehta/disruptor/v1"
)

func BenchmarkReserveOneOneConsumer(b *testing.B) {
	benchMarkDisruptor(b, ReserveOne, SampleConsumer{})
}
func BenchmarkReserveManyOneConsumer(b *testing.B) {
	benchMarkDisruptor(b, ReserveMany, SampleConsumer{})
}
func BenchmarkReserveOneMultipleConsumer(b *testing.B) {
	benchMarkDisruptor(b, ReserveOne, SampleConsumer{}, SampleConsumer{})
}
func BenchmarkReserveManyMultipleConsumer(b *testing.B) {
	benchMarkDisruptor(b, ReserveMany, SampleConsumer{}, SampleConsumer{})

}
func benchMarkDisruptor(b *testing.B, count int64, consumers ...disruptor.CustomConsumer) {
	iterations := int64(b.N)

	dis := build(consumers...)
	go func() {
		b.ReportAllocs()
		b.ResetTimer()
		var seq int64 = -1
		for seq < iterations {
			seq = dis.Reserve(count)
			for i := seq - (count - 1); i <= seq; i++ {
				ringBuffer[i&RingBufferMask] = i
			}
			dis.Commit(seq, seq)
		}
		for _, reader := range dis.Readers {
			reader.Close()
		}
	}()
	read(dis)
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

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	var message int64
	for lower <= upper {
		message = ringBuffer[lower&RingBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}

func build(consumers ...disruptor.CustomConsumer) disruptor.Disruptor {
	return disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithConsumerGroups(consumers...))
}

const (
	RingBufferSize = 1024 * 64
	RingBufferMask = RingBufferSize - 1
	ReserveOne     = 1
	ReserveMany    = 16
)

var ringBuffer = [RingBufferSize]int64{}
