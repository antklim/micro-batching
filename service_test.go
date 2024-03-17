package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

type TestBP struct{}

func (bp *TestBP) Process(props mb.ProcessProps) []mb.JobResult {
	return nil
}

var _ mb.BatchProcessor = (*TestBP)(nil)

type TestJob struct{}

func (j *TestJob) Do() mb.JobResult {
	return mb.JobResult{Result: "OK"}
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

func TestServiceAddJob(t *testing.T) {
	srv := mb.NewService(&TestBP{})

	assert.Equal(t, 0, srv.JobsQueueSize())

	jobsNum := 5
	for i := 0; i < jobsNum; i++ {
		_, err := srv.AddJob(&TestJob{})
		assert.NoError(t, err)
	}

	assert.Equal(t, jobsNum, srv.JobsQueueSize())
}

func TestServiceJobResult(t *testing.T) {
	srv := mb.NewService(&TestBP{})

	jobID, err := srv.AddJob(&TestJob{})
	assert.NoError(t, err)

	jobResult, err := srv.JobResult(jobID)
	assert.NoError(t, err)

	assert.Equal(t, mb.JobResult{Err: nil, Result: nil}, jobResult)
}

func TestServiceJobResultWhenJobIsNotFound(t *testing.T) {
	srv := mb.NewService(&TestBP{})

	_, err := srv.JobResult(ulid.ULID{})
	assert.Equal(t, mb.ErrJobNotFound, err)
}

func TestServiceAddJobWhenShuttingDown(t *testing.T) {
	srv := mb.NewService(&TestBP{})
	srv.Shutdown()

	_, err := srv.AddJob(&TestJob{})
	assert.Equal(t, mb.ErrServiceClosed, err)
}

func TestServiceProcessJobs(t *testing.T) {
	t.Skip("not implemented")
	// Processes jobs in the queue.
	// Calls the processor with the batch of jobs every X interval.
}
