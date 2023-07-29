package disruptor

type Disruptor struct {
	Readers []Reader
	Writer
}

func NewDisruptor(reader []Reader, writer Writer) Disruptor {
	return Disruptor{reader, writer}
}
