package microbatching

import (
	"time"
)

// Runner is a micro-batching runner. It reads batches from the channel and stores them into a queue.
// It processes the queue in a batch when the ticker ticks.
type Runner struct {
	batchProcessor BatchProcessor
	bc             <-chan []Job
	nc             chan<- JobNotification
	freq           time.Duration
	queue          [][]Job
}

func NewRunner(bp BatchProcessor, bc <-chan []Job, nc chan<- JobNotification, freq time.Duration) *Runner {
	return &Runner{
		batchProcessor: bp,
		bc:             bc,
		nc:             nc,
		freq:           freq,
		queue:          make([][]Job, 0),
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
				r.notify(result)
			}

			r.queue = nil
		}
	}
}

func (r *Runner) notify(results []ProcessingResult) {
	for _, result := range results {
		r.nc <- JobNotification{
			JobID:     result.JobID,
			State:     Completed,
			JobResult: result.JobResult,
		}
	}
}
