package disruptor_test

import (
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

func TestCreateProducer(t *testing.T) {

	producerSeq := disruptor.NewSequence()
	consumerSeq := disruptor.NewSequence()

	producer := *disruptor.NewProducer(producerSeq, consumerSeq, 1024)
	if producer == (disruptor.Producer{}) {
		t.Fatal("empty producer recieved.")
	}
}
func TestProducerReserve(t *testing.T) {

	producerSeq := disruptor.NewSequence()
	consumerSeq := disruptor.NewSequence()

	producer := disruptor.NewProducer(producerSeq, disruptor.NewCompositeConsumerBarrier(consumerSeq), 1024)
	upper := producer.Reserve(1024)
	producer.Commit(upper, upper)
	if upper != 1023 {
		t.Fatal("producer current sequence is not at 1023.")
	}
	consumerSeq.Store(consumerSeq.Load() + 10)
	upper = producer.Reserve(5)
	producer.Commit(upper, upper)
	if upper != 1028 {
		t.Fatal("producer current sequence is not at 1028.")
	}

}
func TestProducerCommit(t *testing.T) {

	producerSeq := disruptor.NewSequence()
	consumerSeq := disruptor.NewSequence()

	producer := disruptor.NewProducer(producerSeq, disruptor.NewCompositeConsumerBarrier(consumerSeq), 1024)
	upper := producer.Reserve(10)
	producer.Commit(upper, upper)
	if producer.Reserve(1) != 10 {
		t.Fatal("producer sequence has not been commited 10.")
	}

}
