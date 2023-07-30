package disruptor

import (
	"sync"
)

type Disruptor struct {
	ringBuffer []int64  // using slice as golang dont support this [SIZE]int{}, (comparison b/w slice and array speed also most same in updating the index.)
	capacity   int64    // ring buffer size i.e 2^n
	bufferMask int64    // ring buffer size - 1
	batchSize  int64    // it store the batch size for publishing msgs in batches.
	Readers    []Reader // contains all custom consumers interface given by user. Reader interface to interact with consumer.
	Writer     Writer   // writer interface to interact with producer.
}

func NewDisruptor(readers []Reader, writer Writer, capacity, batchSize int64) Disruptor {
	return Disruptor{Readers: readers, Writer: writer, bufferMask: capacity - 1, batchSize: batchSize, capacity: capacity, ringBuffer: make([]int64, capacity)}
}

func (d Disruptor) Publish(msg int64) {
	sequence := d.Writer.Reserve(1)
	d.ringBuffer[sequence&d.bufferMask] = msg
	d.Writer.Commit(sequence-1, sequence)
}

func (d Disruptor) PublishBatch(msgs []int64) {
	msgLen := len(msgs)
	msgInd := 0
	for msgInd < msgLen {
		var maxSequenceInserted int64
		lastSequence := d.Writer.Reserve(d.batchSize) // Reserve() will give the batch index till upto we can commit.
		for lower := lastSequence - d.batchSize + 1; lower <= lastSequence && msgInd < msgLen; lower++ {
			d.ringBuffer[lower&d.bufferMask] = msgs[msgInd]
			maxSequenceInserted = lower
			msgInd++
		}
		d.Writer.Commit(lastSequence-d.batchSize+1, maxSequenceInserted) //Batch commit
	}
}

func (d Disruptor) Read() {
	wg := sync.WaitGroup{}
	wg.Add(len(d.Readers))
	for _, reader := range d.Readers {
		go func(reader Reader, ringBuffer []int64, bufferMask int64) {
			reader.Read(ringBuffer, bufferMask)
			wg.Done()
		}(reader, d.ringBuffer, d.bufferMask)
	}
	wg.Wait()
}

func (d Disruptor) Close() {
	for _, reader := range d.Readers {
		reader.Close()
	}
}
