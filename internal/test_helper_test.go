package internal_test

import (
	"fmt"

	mb "github.com/antklim/micro-batching"
)

// mocks batch processor for testing purposes
type mockBatchProcessor struct{}

func (m *mockBatchProcessor) Process(jobs []mb.Job) []mb.ProcessingResult {
	var result []mb.ProcessingResult

	fmt.Printf("Processing batch: %v\n", jobs)

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
