package microbatching

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestBP struct {
	counter atomic.Uint32
}

func (bp *TestBP) Process(props ProcessProps) []JobResult {
	bp.counter.Add(1)
	return nil
}

func (bp *TestBP) Counter() uint32 {
	return bp.counter.Load()
}

var _ BatchProcessor = (*TestBP)(nil)

func TestBatchRunnerCallsBatchProcessorByFrequency(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	bp := &TestBP{}

	go batchRunner(batchRunnerProps{
		processor: bp,
		batchSize: 3,
		frequency: 100 * time.Millisecond,
		done:      done,
	})

	time.Sleep(510 * time.Millisecond)

	done <- true

	assert.Equal(t, uint32(5), bp.Counter())
}
