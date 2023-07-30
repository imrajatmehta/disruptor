package disruptor

import "errors"

type Option func(*Wireup)

type Wireup struct {
	waitStrategy   Waiter
	capacity       int64
	batchSize      int64
	consumerGroups [][]EventHandler //Having user create custom consumers
}

func WithConsumerGroups(c ...EventHandler) Option {
	return func(w *Wireup) { w.consumerGroups = append(w.consumerGroups, c) }
}
func WithCapacity(value int64) Option      { return func(w *Wireup) { w.capacity = value } }
func WithWaitStrategy(waitS Waiter) Option { return func(w *Wireup) { w.waitStrategy = waitS } }
func WithBatchSize(v int64) Option         { return func(w *Wireup) { w.batchSize = v } }

func New(options ...Option) Disruptor {
	if w, err := NewWireUp(options...); err != nil {
		panic(err)
	} else {
		consumers, producer := w.build()
		return NewDisruptor(consumers, producer, w.capacity, w.batchSize)
	}
}
func NewWireUp(options ...Option) (*Wireup, error) {
	w := &Wireup{}
	WithWaitStrategy(NewWaitStrategy())(w)
	for _, option := range options {
		option(w)
	}
	if err := w.validate(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Wireup) build() ([]Reader, Writer) {
	writerSequence := NewSequence()
	consumers, consumerBarrier := w.buildConsumers(writerSequence)
	return consumers, NewProducer(writerSequence, consumerBarrier, w.capacity)

}
func (w *Wireup) validate() error {

	if w.capacity < 1 {
		return errors.New("the capacity must be at least 1")
	}
	if w.capacity&(w.capacity-1) != 0 {
		return errors.New("the capacity is not a power of 2")
	}
	if w.batchSize < 1 {
		return errors.New("the batchSize must be at least 1")
	}
	if w.batchSize >= w.capacity {
		return errors.New("the batchSize must be smaller than capacity")
	}
	if len(w.consumerGroups) == 0 {
		return errors.New("the consumer group dont have any consumers")
	}
	for _, consumerGroup := range w.consumerGroups {
		if len(consumerGroup) == 0 {
			return errors.New("the consumer group does not have any consumers")
		}
		for _, consumer := range consumerGroup {
			if consumer == nil {
				return errors.New("an empty consumer was specified in consumer group")
			}
		}
	}
	return nil
}

func (w *Wireup) buildConsumers(writerSequence *Sequence) (readers []Reader, upstream Barrier) {
	var consumerSequences []*Sequence
	for _, consumerGroup := range w.consumerGroups {
		for _, callerConsumer := range consumerGroup {
			sequence := NewSequence()
			readers = append(readers, NewConsumer(sequence, writerSequence, callerConsumer, w.waitStrategy))
			consumerSequences = append(consumerSequences, sequence)
		}
		upstream = NewCompositeConsumerBarrier(consumerSequences...)
	}
	return
}
