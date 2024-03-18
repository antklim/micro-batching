package internal_test

import (
	"fmt"
	"testing"
	"time"

	mb "github.com/antklim/micro-batching/internal"
)

type mockBatchProcessor struct{}

func (m *mockBatchProcessor) Process(jobs []int) []mb.ProcessResult {
	var result []mb.ProcessResult

	fmt.Printf("Processing batch: %v\n", jobs)

	for _, j := range jobs {
		result = append(result, mb.ProcessResult{ID: j, Result: j})
	}

	return result
}

var _ mb.BatchProcessor = (*mockBatchProcessor)(nil)

func TestRunner(t *testing.T) {
	batches := make(chan []int)

	runner := mb.NewRunner(&mockBatchProcessor{}, batches, 10*time.Millisecond)

	go func() {
		for i := 0; i < 11; i++ {
			batches <- []int{i, i + 1, i + 2}
		}

		time.Sleep(50 * time.Millisecond)

		for i := 11; i < 22; i++ {
			batches <- []int{i, i + 1, i + 2}
		}

		time.Sleep(50 * time.Millisecond)

		close(batches)
	}()

	go runner.Run()

	time.Sleep(200 * time.Millisecond)
}
