package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/stretchr/testify/assert"
)

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
				mb.WithQueueSize(200),
			},
			expectedBatchSize: 15,
			expectedFrequency: 5 * time.Minute,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := mb.NewService(tC.opts...)

			assert.Equal(t, tC.expectedBatchSize, srv.BatchSize())
			assert.Equal(t, tC.expectedFrequency, srv.Frequency())
		})
	}
}

func TestServiceJobResult(t *testing.T) {
	srv := mb.NewService(mb.WithFrequency(10 * time.Millisecond))
	srv.Run(&mockBatchProcessor{})

	testJobID := "test-job-id"
	testJob := newMockJob(testJobID)

	err := srv.AddJob(testJob)
	assert.NoError(t, err)

	jobResult, err := srv.JobResult(testJobID)
	assert.NoError(t, err)

	assert.Equal(t, mb.Submitted, jobResult.State)
	assert.Equal(
		t,
		mb.JobExtendedResult{
			JobID:     "test-job-id",
			State:     mb.Submitted,
			JobResult: mb.JobResult{Err: nil, Result: nil},
		},
		jobResult,
	)

	// wait for the job to be processed
	time.Sleep(50 * time.Millisecond)

	jobResult, err = srv.JobResult(testJobID)
	assert.NoError(t, err)

	assert.Equal(t, mb.Completed, jobResult.State)
}

func TestServiceJobResultWhenJobIsNotFound(t *testing.T) {
	srv := mb.NewService()
	srv.Run(&mockBatchProcessor{})

	_, err := srv.JobResult("non-existing-job-id")
	assert.Equal(t, mb.ErrJobNotFound, err)
}

func TestServiceAddJobWhenShuttingDown(t *testing.T) {
	srv := mb.NewService(mb.WithShutdownTimeout(10 * time.Millisecond))
	srv.Shutdown()

	err := srv.AddJob(newMockJob("test-job-id"))
	assert.Equal(t, mb.ErrServiceClosed, err)
}
