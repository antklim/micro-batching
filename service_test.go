package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/stretchr/testify/assert"
)

type TestBP struct{}

func (bp *TestBP) Process(jobs []mb.Job) []mb.JobResult {
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
	}{
		{
			desc:              "inits service with default options",
			expectedBatchSize: 3,
			expectedFrequency: time.Second,
		},
		{
			desc: "inits service with custom options",
			opts: []mb.ServiceOption{
				mb.WithBatchSize(15),
				mb.WithFrequency(5 * time.Minute),
			},
			expectedBatchSize: 15,
			expectedFrequency: 5 * time.Minute,
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

func TestServiceAddJobWhenShuttingDown(t *testing.T) {
	srv := mb.NewService(&TestBP{})
	srv.Shutdown()

	_, err := srv.AddJob(&TestJob{})
	assert.Equal(t, mb.ErrServiceClosed, err)
}

func TestServiceProcessJobs(t *testing.T) {
	// Processes jobs in the queue.
	// Calls the processor with the batch of jobs every X interval.
}
