package microbatching_test

import (
	"fmt"
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/stretchr/testify/assert"
)

func TestServiceInit(t *testing.T) {
	testCases := []struct {
		desc string
		opts []mb.ServiceOption
	}{
		{
			desc: "inits service with default options",
		},
		{
			desc: "inits service with custom options",
			opts: []mb.ServiceOption{
				mb.WithBatchSize(15),
				mb.WithFrequency(5 * time.Minute),
				mb.WithQueueSize(200),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := mb.NewService(tC.opts...)

			assert.NotEmpty(t, srv)
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
	srv.Run(&mockBatchProcessor{})
	srv.Shutdown()

	err := srv.AddJob(newMockJob("test-job-id"))
	assert.Equal(t, mb.ErrServiceClosed, err)
}

func TestShutdownWaitsForJobsToBeDone(t *testing.T) {
	srv := mb.NewService(mb.WithFrequency(10 * time.Millisecond))
	srv.Run(&mockBatchProcessor{})

	for i := 0; i < 10; i++ {
		err := srv.AddJob(newMockJob(fmt.Sprintf("test-job-id-%d", i)))
		assert.NoError(t, err)
	}

	srv.Shutdown()

	for i := 0; i < 10; i++ {
		result, err := srv.JobResult(fmt.Sprintf("test-job-id-%d", i))
		assert.NoError(t, err)
		assert.Equal(t, mb.Completed, result.State)
	}
}
