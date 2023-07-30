package disruptor_test

import (
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

func TestCreate(t *testing.T) {
	seq := disruptor.NewSequence()
	if seq.Load() != -1 {
		t.Fatal("sequence not initilized !")
	}
}

func TestSetandGet(t *testing.T) {
	seq := disruptor.NewSequence()
	seq.Store(10)
	if seq.Load() != 10 {
		t.Fatal("sequencer set/get doesn't work !")
	}
}
