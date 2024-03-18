package internal

import (
	"fmt"
	"time"

	mb "github.com/antklim/micro-batching"
)

// Runner is a micro-batching runner. It reads batches from the channel and stores them into a queue.
// It processes the queue in a batch when the ticker ticks.
type Runner struct {
	batchProcessor mb.BatchProcessor
	bc             <-chan []mb.Job
	nc             chan<- mb.JobNotification
	freq           time.Duration
	queue          [][]mb.Job
}

func NewRunner(bp mb.BatchProcessor, bc <-chan []mb.Job, nc chan<- mb.JobNotification, freq time.Duration) *Runner {
	return &Runner{
		batchProcessor: bp,
		bc:             bc,
		nc:             nc,
		freq:           freq,
		queue:          make([][]mb.Job, 0),
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

			r.queue = append(r.queue, batch)
		case <-ticker.C:
			for _, batch := range r.queue {
				result := r.batchProcessor.Process(batch)
				fmt.Printf("Processed batch: %v\n", result)
			}

			r.queue = nil
		}
	}
}

func (r *Runner) notify(results []mb.ProcessingResult) {
	for _, result := range results {
		r.nc <- mb.JobNotification{
			JobID:     result.JobID,
			State:     mb.Completed,
			JobResult: result.JobResult,
		}
	}
}
