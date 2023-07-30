package benchmarks

import (
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

func BenchmarkReserveOneOneConsumer(b *testing.B) {
	benchMarkDisruptor(b, SampleConsumer{})
}
func BenchmarkReserveManyOneConsumer(b *testing.B) {
	benchMarkDisruptorBatch(b, SampleConsumer{})
}
func BenchmarkReserveOneMultipleConsumer(b *testing.B) {
	benchMarkDisruptor(b, SampleConsumer{}, SampleConsumer{})
}
func BenchmarkReserveManyMultipleConsumer(b *testing.B) {
	benchMarkDisruptorBatch(b, SampleConsumer{}, SampleConsumer{})

}
func benchMarkDisruptorBatch(b *testing.B, consumers ...disruptor.EventHandler) {
	iterations := int64(b.N)

	d, _ := disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithBatchSize(ReserveMany),
		disruptor.WithConsumerGroups(consumers...))
	go func() {
		b.ReportAllocs()
		b.ResetTimer()
		for seq := int64(0); seq < iterations; seq += 16 {
			d.PublishBatch([]int64{seq + 1, seq + 2, seq + 3, seq + 4, seq + 5, seq + 6, seq + 7, seq + 8, seq + 9, seq + 10, seq + 11, seq + 12, seq + 13, seq + 14, seq + 15, seq + 16})
		}
		d.Close()
	}()
	d.Read()
}

func benchMarkDisruptor(b *testing.B, consumers ...disruptor.EventHandler) {
	iterations := int64(b.N)

	d, _ := disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithBatchSize(ReserveOne),
		disruptor.WithConsumerGroups(consumers...))
	go func() {
		b.ReportAllocs()
		b.ResetTimer()
		for seq := int64(0); seq < iterations; seq++ {
			d.Publish(seq)
		}
		d.Close()
	}()
	d.Read()
}

type SampleConsumer struct{}

func (s SampleConsumer) OnEvent(msg int64) {}

const (
	RingBufferSize = 1024 * 64
	ReserveOne     = 1
	ReserveMany    = 16
)
