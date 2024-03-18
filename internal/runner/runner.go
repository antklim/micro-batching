package runner

import (
	"fmt"
	"time"
)

type ProcessResult struct {
	ID     int
	Result int
}

type BatchProcessor interface {
	Process(jobs []int) []ProcessResult
}

type Runner struct {
	batchProcessor BatchProcessor
	bc             <-chan []int
	freq           time.Duration
	batches        [][]int
}

func NewRunner(bp BatchProcessor, bc <-chan []int, freq time.Duration) *Runner {
	return &Runner{
		batchProcessor: bp,
		bc:             bc,
		freq:           freq,
		batches:        make([][]int, 0),
	}
}

func (r *Runner) Run() {
	ticker := time.NewTicker(r.freq)

	for {
		select {
		case batch, ok := <-r.bc:
			fmt.Printf("Batch %v, %v\n", batch, ok)

			if !ok {
				ticker.Stop()
				return
			}

			r.batches = append(r.batches, batch)
		case <-ticker.C:
			fmt.Println("Ticker ...")

			for _, batch := range r.batches {
				result := r.batchProcessor.Process(batch)
				fmt.Printf("Processed batch: %v\n", result)
			}

			r.batches = nil
		}
	}
}
