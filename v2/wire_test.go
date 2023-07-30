package disruptor_test

import (
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

type NewConsumer struct{}

func (obj NewConsumer) OnEvent(msg int64) {}
func TestWireCapacity(t *testing.T) {
	_, err := disruptor.New(
		disruptor.WithCapacity(10),
		disruptor.WithBatchSize(5),
		disruptor.WithConsumerGroups(NewConsumer{}))
	if err == nil {
		t.Fatal("capacity is not the power of 2, but it still works.")
	}
	_, err = disruptor.New(
		disruptor.WithCapacity(8),
		disruptor.WithBatchSize(5),
		disruptor.WithConsumerGroups(NewConsumer{}))
	if err != nil {
		t.Fatal("capacity is not the power of 2, but it still works.")
	}
}
func TestWireBatchSize(t *testing.T) {
	_, err := disruptor.New(
		disruptor.WithCapacity(64),
		disruptor.WithBatchSize(70),
		disruptor.WithConsumerGroups(NewConsumer{}))
	if err == nil {
		t.Fatal("batch size is greater than capacity but still it works.")
	}
	_, err = disruptor.New(
		disruptor.WithCapacity(64),
		disruptor.WithBatchSize(5),
		disruptor.WithConsumerGroups(NewConsumer{}))
	if err != nil {
		t.Fatal("batch size is smaller than capacity still it doesn't works.")
	}
}
func TestWireConsumerGroup(t *testing.T) {
	_, err := disruptor.New(
		disruptor.WithCapacity(64),
		disruptor.WithBatchSize(8),
		disruptor.WithConsumerGroups())
	if err == nil {
		t.Fatal("no consumer was provided, but still it works.")
	}
	_, err = disruptor.New(
		disruptor.WithCapacity(64),
		disruptor.WithBatchSize(5),
		disruptor.WithConsumerGroups(NewConsumer{}))
	if err != nil {
		t.Fatal("consumer was provided still it doesn't works.")
	}
}
