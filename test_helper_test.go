package microbatching_test

import (
	"strconv"

	mb "github.com/antklim/micro-batching"
)

// mocks batch processor for testing purposes
type mockBatchProcessor struct{}

func (m *mockBatchProcessor) Process(jobs []mb.Job) []mb.ProcessingResult {
	var result []mb.ProcessingResult

	for _, j := range jobs {
		jr := j.Do()
		result = append(result, mb.ProcessingResult{
			JobID:     j.ID(),
			JobResult: jr,
		})
	}

	return result
}

var _ mb.BatchProcessor = (*mockBatchProcessor)(nil)

// mocks job for testing purposes
type mockJob struct {
	id string
}

func newMockJob(id string) *mockJob {
	return &mockJob{id: id}
}

func (m *mockJob) Do() mb.JobResult {
	return mb.JobResult{Result: "OK"}
}

func (m *mockJob) ID() string {
	return m.id
}

var _ mb.Job = (*mockJob)(nil)

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
