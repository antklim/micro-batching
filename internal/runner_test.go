package internal_test

import (
	"strconv"
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	internal "github.com/antklim/micro-batching/internal"
)

func makeMockJobs(n int) []mb.Job {
	var jobs []mb.Job

	for i := 0; i < n; i++ {
		jobs = append(jobs, newMockJob(strconv.Itoa(i)))
	}

	return jobs
}

func makeMockBatches(jobs []mb.Job, batchSize int) [][]mb.Job {
	var batches [][]mb.Job
	var batch []mb.Job

	for i, job := range jobs {
		batch = append(batch, job)

		if (i+1)%batchSize == 0 {
			batches = append(batches, batch)
			batch = nil
		}
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

func TestRunner(t *testing.T) {
	batches := make(chan []mb.Job)

	runner := internal.NewRunner(&mockBatchProcessor{}, batches, nil, 10*time.Millisecond)

	testJobs := makeMockJobs(22)
	testBatches := makeMockBatches(testJobs, 3)

	go func() {
		// send first part of the batches
		for i := 0; i < 4; i++ {
			batches <- testBatches[i]
		}

		time.Sleep(50 * time.Millisecond)

		// send the remaining batches
		for i := 4; i < len(testBatches); i++ {
			batches <- testBatches[i]
		}

		time.Sleep(50 * time.Millisecond)

		close(batches)
	}()

	go runner.Run()

	time.Sleep(200 * time.Millisecond)
}
