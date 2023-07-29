package disruptor

type Reader interface {
	Read()
	Close()
}
type Writer interface {
	Reserve(int64) int64
	Commit(int64, int64)
}
type Barrier interface {
	Load() int64
}

type CustomConsumer interface {
	Consume(int64, int64)
}
