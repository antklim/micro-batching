package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/stretchr/testify/assert"
)

func batchSender(bc chan<- []mb.Job, testBatches [][]mb.Job) {
	// send first part of the batches
	for i := 0; i < 4; i++ {
		bc <- testBatches[i]
	}

	time.Sleep(50 * time.Millisecond)

	// send the remaining batches
	for i := 4; i < len(testBatches); i++ {
		bc <- testBatches[i]
	}

	time.Sleep(50 * time.Millisecond)

	close(bc)
}

func TestRunner(t *testing.T) {
	bc := make(chan []mb.Job)
	nc := make(chan mb.JobExtendedResult)

	jobsSize := 22
	batchSize := 3

	testJobs := makeMockJobs(jobsSize)
	testBatches := makeMockBatches(testJobs, batchSize)
	notifications := make([]mb.JobExtendedResult, 0)

	runner := mb.NewRunner(&mockBatchProcessor{}, bc, nc, 10*time.Millisecond)

	go batchSender(bc, testBatches)
	go func() {
		for n := range nc {
			notifications = append(notifications, n)
		}
	}()
	go runner.Run()

	time.Sleep(150 * time.Millisecond)

	close(nc)

	// should receive notifications for all jobs
	assert.Equal(t, jobsSize*2, len(notifications))

	for _, n := range notifications {
		assert.True(t, n.State == mb.Processing || n.State == mb.Completed)
	}
}
