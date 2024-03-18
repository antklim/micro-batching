package internal

import (
	"fmt"
	"time"

	mb "github.com/antklim/micro-batching"
)

type Runner struct {
	batchProcessor mb.BatchProcessor
	bc             <-chan []mb.Job
	freq           time.Duration
	batches        [][]mb.Job
}

func NewRunner(bp mb.BatchProcessor, bc <-chan []mb.Job, freq time.Duration) *Runner {
	return &Runner{
		batchProcessor: bp,
		bc:             bc,
		freq:           freq,
		batches:        make([][]mb.Job, 0),
	}
}

func (r *Runner) Run() {
	ticker := time.NewTicker(r.freq)

	for {
		select {
		case batch, ok := <-r.bc:
			if !ok {
				ticker.Stop()
				return
			}

			r.batches = append(r.batches, batch)
		case <-ticker.C:
			for _, batch := range r.batches {
				result := r.batchProcessor.Process(batch)
				fmt.Printf("Processed batch: %v\n", result)
			}

			r.batches = nil
		}
	}
}
