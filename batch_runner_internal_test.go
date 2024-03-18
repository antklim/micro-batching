package microbatching

// import (
// 	"sync/atomic"
// 	"testing"
// 	"time"

// 	"github.com/oklog/ulid/v2"
// 	"github.com/stretchr/testify/assert"
// )

// type TestJob struct {
// 	id string
// }

// func NewTestJob() *TestJob {
// 	return &TestJob{id: ulid.Make().String()}
// }

// func (j *TestJob) ID() string {
// 	return j.id
// }

// func (j *TestJob) Do() JobResult {
// 	return JobResult{Result: "OK"}
// }

// var _ Job = (*TestJob)(nil)

// type TestBP struct {
// 	counter atomic.Uint32
// }

// func (bp *TestBP) Process(jobs []Job) []JobResult {
// 	bp.counter.Add(1)

// 	// fmt.Printf("Processing batch #%d\n", bp.counter.Load())

// 	for _, j := range jobs {
// 		// fmt.Printf("Processing job #%d\n", j.(*TestJob).ID)
// 		j.Do()
// 	}

// 	return nil
// }

// func (bp *TestBP) Counter() uint32 {
// 	return bp.counter.Load()
// }

// var _ BatchProcessor = (*TestBP)(nil)

// func TestBatchRunnerInvokesByFrequency(t *testing.T) {
// 	done := make(chan bool)
// 	defer close(done)

// 	bp := &TestBP{}

// 	br := batchRunner{
// 		processor: bp,
// 		batchSize: 3,
// 		frequency: 10 * time.Millisecond,
// 		done:      done,
// 	}

// 	go br.run()

// 	time.Sleep(52 * time.Millisecond)

// 	done <- true

// 	assert.Equal(t, uint32(5), br.ticks.Load())
// }

// func TestBatchRunnerDoesNotCallBatchProcessorWithoutJobs(t *testing.T) {
// 	done := make(chan bool)
// 	defer close(done)

// 	bp := &TestBP{}

// 	br := batchRunner{
// 		processor: bp,
// 		batchSize: 3,
// 		frequency: 10 * time.Millisecond,
// 		done:      done,
// 	}

// 	go br.run()

// 	time.Sleep(52 * time.Millisecond)

// 	done <- true

// 	assert.Equal(t, uint32(0), bp.Counter())
// }

// func TestBatchRunnerJobProcessing(t *testing.T) {
// 	jobs := make(chan job, 100)
// 	jobNotifications := make(chan jobNotification, 100)
// 	done := make(chan bool)
// 	defer close(jobs)
// 	defer close(jobNotifications)
// 	defer close(done)

// 	bp := &TestBP{}

// 	br := batchRunner{
// 		processor:        bp,
// 		batchSize:        3,
// 		frequency:        10 * time.Millisecond,
// 		jobs:             jobs,
// 		jobNotifications: jobNotifications,
// 		done:             done,
// 	}

// 	for i := 0; i < 11; i++ {
// 		jobs <- job{Job: &TestJob{}}
// 	}

// 	// make sure the jobs are in the queue before starting the batch runner
// 	assert.Equal(t, 11, len(jobs))

// 	go br.run()

// 	time.Sleep(52 * time.Millisecond)

// 	assert.Equal(t, 0, len(jobs))

// 	for i := 11; i < 22; i++ {
// 		jobs <- job{Job: &TestJob{}}
// 	}

// 	time.Sleep(52 * time.Millisecond)

// 	done <- true

// 	assert.Equal(t, 0, len(jobs))

// 	assert.Equal(t, uint32(8), bp.Counter())
// }

// func TestBatchRunnerAfterAbortSignal(t *testing.T) {
// 	jobs := make(chan job, 100)
// 	jobNotifications := make(chan jobNotification, 100)
// 	done := make(chan bool)
// 	defer close(jobs)
// 	defer close(jobNotifications)
// 	defer close(done)

// 	bp := &TestBP{}
// 	br := batchRunner{
// 		processor:        bp,
// 		batchSize:        3,
// 		frequency:        10 * time.Millisecond,
// 		jobs:             jobs,
// 		jobNotifications: jobNotifications,
// 		done:             done,
// 	}

// 	for i := 0; i < 11; i++ {
// 		jobs <- job{Job: &TestJob{}}
// 	}

// 	assert.Equal(t, 11, len(jobs))

// 	go br.run()

// 	done <- true

// 	time.Sleep(52 * time.Millisecond)

// 	assert.Equal(t, 0, len(jobs))

// 	// Make sure we haven't process anything within the ticker interval
// 	assert.Equal(t, uint32(0), br.ticks.Load())

// 	// 11 items goes into 4 batches of size 3
// 	assert.Equal(t, uint32(4), bp.Counter())
// }
