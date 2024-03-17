package microbatching

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

type TestJob struct {
	ID int
	T  time.Time
}

func (j *TestJob) Do() JobResult {
	return JobResult{Result: "OK"}
}

var _ Job = (*TestJob)(nil)

type TestBP struct {
	counter atomic.Uint32
}

func (bp *TestBP) Process(props ProcessProps) []JobResult {
	bp.counter.Add(1)

	for _, j := range props.jobs {
		j.Do()
	}

	return nil
}

func (bp *TestBP) Counter() uint32 {
	return bp.counter.Load()
}

var _ BatchProcessor = (*TestBP)(nil)

func TestBatchRunnerInvokesByFrequency(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	bp := &TestBP{}

	br := batchRunner{
		processor: bp,
		batchSize: 3,
		frequency: 10 * time.Millisecond,
		done:      done,
	}

	go br.run()

	time.Sleep(52 * time.Millisecond)

	done <- true

	assert.Equal(t, uint32(5), br.ticks.Load())
}

func TestBatchRunnerDoesNotCallBatchProcessorWithoutJobs(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	bp := &TestBP{}

	br := batchRunner{
		processor: bp,
		batchSize: 3,
		frequency: 10 * time.Millisecond,
		done:      done,
	}

	go br.run()

	time.Sleep(52 * time.Millisecond)

	done <- true

	assert.Equal(t, uint32(0), bp.Counter())
}

func TestBatchRunnerJobProcessing(t *testing.T) {
	jobs := make(chan job, 100)
	done := make(chan bool)
	defer close(jobs)
	defer close(done)

	bp := &TestBP{}

	br := batchRunner{
		processor: bp,
		batchSize: 3,
		frequency: 10 * time.Millisecond,
		jobs:      jobs,
		done:      done,
	}

	for i := 0; i < 11; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	// make sure the jobs are in the queue before starting the batch runner
	assert.Equal(t, 11, len(jobs))

	go br.run()

	time.Sleep(52 * time.Millisecond)

	assert.Equal(t, 0, len(jobs))

	for i := 11; i < 22; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	time.Sleep(52 * time.Millisecond)

	done <- true

	assert.Equal(t, 0, len(jobs))

	assert.Equal(t, uint32(8), bp.Counter())
}

func TestBatchRunnerAfterAbortSignal(t *testing.T) {
	jobs := make(chan job, 100)
	done := make(chan bool)
	defer close(jobs)
	defer close(done)

	bp := &TestBP{}
	br := batchRunner{
		processor: bp,
		batchSize: 3,
		frequency: 10 * time.Millisecond,
		jobs:      jobs,
		done:      done,
	}

	for i := 0; i < 11; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	assert.Equal(t, 11, len(jobs))

	go br.run()

	done <- true

	time.Sleep(52 * time.Millisecond)

	assert.Equal(t, 0, len(jobs))

	// Make sure we haven't process anything within the ticker interval
	assert.Equal(t, uint32(0), br.ticks.Load())

	// 11 items goes into 4 batches of size 3
	assert.Equal(t, uint32(4), bp.Counter())
}
