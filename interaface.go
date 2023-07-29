package disruptor

import "errors"

// Act as a wall b/w user and consumer functionalities.
type Reader interface {
	Read()
	Close()
}

// Act as a wall b/w user and producer functionalities.
type Writer interface {
	Reserve(int64) int64
	Commit(int64, int64)
}

// Act as wall b/w producer and consumer to check the current sequence of consumer, and producer.
type Barrier interface {
	Load() int64
}

// User will provide its custom consumer struct{}, as a input we will take this as this interface so that user will have to implement Consume on its side.
type CustomConsumer interface {
	Consume(int64, int64)
}

type Waiter interface {
	Idle()
}

var errMinimumReservationSize = errors.New("the minimum reservation size is 1 slot")
