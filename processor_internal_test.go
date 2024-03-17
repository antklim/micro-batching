package microbatching

import (
	"fmt"
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
	fmt.Printf("job %d done\n", j.ID)
	return JobResult{Result: "OK"}
}

var _ Job = (*TestJob)(nil)

type TestBP struct {
	counter atomic.Uint32
}

func (bp *TestBP) Process(props ProcessProps) []JobResult {
	bp.counter.Add(1)

	fmt.Printf("Starting batch #%d, len %d\n", bp.counter.Load(), len(props.jobs))

	for _, j := range props.jobs {
		j.Do()
	}

	return nil
}

func (bp *TestBP) Counter() uint32 {
	return bp.counter.Load()
}

var _ BatchProcessor = (*TestBP)(nil)

func TestBatchRunnerDoesNotCallBatchProcessorWhenTheresNoJobs(t *testing.T) {
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

	assert.Equal(t, uint32(0), bp.Counter())
}

func TestBatchRunnerJobProcessing(t *testing.T) {
	jobs := make(chan job, 100)
	done := make(chan bool)
	defer close(jobs)
	defer close(done)

	bp := &TestBP{}

	for i := 0; i < 11; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	assert.Equal(t, 11, len(jobs))

	go batchRunner(batchRunnerProps{
		processor: bp,
		batchSize: 3,
		frequency: 100 * time.Millisecond,
		jobs:      jobs,
		done:      done,
	})

	time.Sleep(500 * time.Millisecond)

	for i := 11; i < 22; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	assert.Equal(t, 11, len(jobs))

	time.Sleep(500 * time.Millisecond)

	done <- true

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 0, len(jobs))
	assert.Equal(t, uint32(8), bp.Counter())
}

func TestBatchRunnerAfterAbortSignal(t *testing.T) {
	jobs := make(chan job, 100)
	done := make(chan bool)
	defer close(jobs)
	defer close(done)

	bp := &TestBP{}

	for i := 0; i < 11; i++ {
		jobs <- job{ID: ulid.Make(), J: &TestJob{ID: i}}
	}

	assert.Equal(t, 11, len(jobs))

	go batchRunner(batchRunnerProps{
		processor: bp,
		batchSize: 3,
		frequency: 100 * time.Millisecond,
		jobs:      jobs,
		done:      done,
	})

	done <- true

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 0, len(jobs))
}
