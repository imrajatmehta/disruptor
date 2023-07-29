package disruptor

import "errors"

type Option func(*Wireup)

type Wireup struct {
	waitStrategy   Waiter
	capacity       int64
	consumerGroups [][]CustomConsumer
}

func WithCapacity(value int64) Option { return func(this *Wireup) { this.capacity = value } }

func WithConsumerGroups(c ...CustomConsumer) Option {
	return func(this *Wireup) { this.consumerGroups = append(this.consumerGroups, c) }
}

func WithWaitStrategy(w Waiter) Option { return func(this *Wireup) { this.waitStrategy = w } }

func New(options ...Option) Disruptor {
	if this, err := NewWireUp(options...); err != nil {
		panic(err)
	} else {
		return NewDisruptor(this.build())
	}
}
func NewWireUp(options ...Option) (*Wireup, error) {
	this := &Wireup{}
	WithWaitStrategy(NewWaitStrategy())(this)
	for _, option := range options {
		option(this)
	}
	if err := this.validate(); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Wireup) build() ([]Reader, Writer) {
	writerSequence := NewSequence()
	consumers, consumerBarrier := this.buildConsumers(writerSequence)
	return consumers, NewProducer(writerSequence, consumerBarrier, this.capacity)

}
func (this *Wireup) validate() error {

	if this.capacity < 1 {
		return errors.New("the capacity must be at least 1")
	}
	if this.capacity&(this.capacity-1) != 0 {
		return errors.New("the capacity is not a power of 2")
	}
	if len(this.consumerGroups) == 0 {
		return errors.New("the consumer group dont have any consumers")
	}
	for _, consumerGroup := range this.consumerGroups {
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

func (this *Wireup) buildConsumers(writerSequence *Sequence) (readers []Reader, upstream Barrier) {
	var consumerSequences []*Sequence
	for _, consumerGroup := range this.consumerGroups {
		for _, callerConsumer := range consumerGroup {
			sequence := NewSequence()
			readers = append(readers, NewConsumer(sequence, writerSequence, callerConsumer, this.waitStrategy))
			consumerSequences = append(consumerSequences, sequence)
		}
		upstream = NewCompositeConsumerBarrier(consumerSequences...)
	}
	return
}
