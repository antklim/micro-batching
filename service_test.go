package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

type TestBP struct{}

func (bp *TestBP) Process(jobs []mb.Job) []mb.JobResult {
	results := make([]mb.JobResult, 0, len(jobs))

	for _, j := range jobs {
		results = append(results, j.Do())
	}

	return results
}

var _ mb.BatchProcessor = (*TestBP)(nil)

type TestJob struct {
	id string
}

func NewTestJob() *TestJob {
	return &TestJob{id: ulid.Make().String()}
}

func (j *TestJob) ID() string {
	return j.id
}

func (j *TestJob) Do() mb.JobResult {
	return mb.JobResult{JobID: j.id, Result: "OK"}
}

var _ mb.Job = (*TestJob)(nil)

func TestServiceInit(t *testing.T) {
	testCases := []struct {
		desc              string
		opts              []mb.ServiceOption
		expectedBatchSize int
		expectedFrequency time.Duration
		expectedQueueSize int
	}{
		{
			desc:              "inits service with default options",
			expectedBatchSize: 3,
			expectedFrequency: time.Second,
			expectedQueueSize: 100,
		},
		{
			desc: "inits service with custom options",
			opts: []mb.ServiceOption{
				mb.WithBatchSize(15),
				mb.WithFrequency(5 * time.Minute),
				mb.WithQueueSize(200),
			},
			expectedBatchSize: 15,
			expectedFrequency: 5 * time.Minute,
			expectedQueueSize: 200,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := mb.NewService(&TestBP{}, tC.opts...)

			assert.Equal(t, tC.expectedBatchSize, srv.BatchSize())
			assert.Equal(t, tC.expectedFrequency, srv.Frequency())
		})
	}
}

func TestServiceJobResult(t *testing.T) {
	srv := mb.NewService(&TestBP{})

	testJob := NewTestJob()

	err := srv.AddJob(testJob)
	assert.NoError(t, err)

	jobState, jobResult, err := srv.JobResult(testJob.ID())
	assert.NoError(t, err)

	assert.Equal(t, mb.Submitted, jobState)
	assert.Equal(t, mb.JobResult{JobID: testJob.ID(), Err: nil, Result: nil}, jobResult)
}

func TestServiceJobResultWhenJobIsNotFound(t *testing.T) {
	srv := mb.NewService(&TestBP{})

	_, _, err := srv.JobResult("non-existing-job-id")
	assert.Equal(t, mb.ErrJobNotFound, err)
}

func TestServiceAddJobWhenShuttingDown(t *testing.T) {
	srv := mb.NewService(&TestBP{})
	srv.Shutdown()

	err := srv.AddJob(NewTestJob())
	assert.Equal(t, mb.ErrServiceClosed, err)
}

func TestServiceProcessJobs(t *testing.T) {
	srv := mb.NewService(&TestBP{}, mb.WithBatchSize(3), mb.WithFrequency(10*time.Millisecond))

	testJob := NewTestJob()

	err := srv.AddJob(testJob)
	assert.NoError(t, err)

	jobState, jobResult, err := srv.JobResult(testJob.ID())
	assert.NoError(t, err)

	assert.Equal(t, mb.Submitted, jobState)
	assert.Equal(t, mb.JobResult{JobID: testJob.id, Err: nil, Result: nil}, jobResult)

	time.Sleep(50 * time.Millisecond)

	jobState, jobResult, err = srv.JobResult(testJob.ID())
	assert.NoError(t, err)

	assert.Equal(t, mb.Completed, jobState)
	assert.Equal(t, mb.JobResult{JobID: testJob.id, Err: nil, Result: "OK"}, jobResult)
}
