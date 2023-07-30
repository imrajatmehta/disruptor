package disruptor_test

import (
	"testing"

	"github.com/imrajatmehta/disruptor/v2"
)

func TestCompositeConsumerBarrier(t *testing.T) {
	s1 := disruptor.NewSequence()
	s2 := disruptor.NewSequence()
	obj := disruptor.NewCompositeConsumerBarrier(s1, s2)
	s2.Store(2)
	s1.Store(5)
	if obj.Load() != 2 {
		t.Fatal("incorrect slow consumer index/sequence")
	}
}
