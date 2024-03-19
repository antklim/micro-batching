package microbatching

import (
	"time"
)

// Runner is a micro-batching runner. It reads batches from the channel and stores them into a queue.
// It processes the queue in a batch when the ticker ticks. It notifies the results to the notification channel.
type Runner struct {
	batchProcessor BatchProcessor
	bc             <-chan []Job
	nc             chan<- JobExtendedResult
	freq           time.Duration
	queue          [][]Job
	shutdown       <-chan bool
}

func NewRunner(
	bp BatchProcessor,
	bc <-chan []Job,
	nc chan<- JobExtendedResult,
	freq time.Duration,
	shutdown <-chan bool,
) *Runner {
	return &Runner{
		batchProcessor: bp,
		bc:             bc,
		nc:             nc,
		freq:           freq,
		queue:          make([][]Job, 0),
		shutdown:       shutdown,
	}
}

func (r *Runner) Run() {
	ticker := time.NewTicker(r.freq)

	for {
		select {
		case <-r.shutdown:
			ticker.Stop()

			for _, batch := range r.queue {
				r.notifyProcessing(batch)
				result := r.batchProcessor.Process(batch)
				r.notifyCompleted(result)
			}

			r.queue = nil

			close(r.nc)

			return
		case batch := <-r.bc:
			r.queue = append(r.queue, batch)
		case <-ticker.C:
			for _, batch := range r.queue {
				r.notifyProcessing(batch)
				result := r.batchProcessor.Process(batch)
				r.notifyCompleted(result)
			}

			r.queue = nil
		}
	}
}

func (r *Runner) notifyProcessing(batch []Job) {
	for _, job := range batch {
		r.nc <- JobExtendedResult{
			JobID: job.ID(),
			State: Processing,
		}
	}
}

func (r *Runner) notifyCompleted(results []ProcessingResult) {
	for _, result := range results {
		r.nc <- JobExtendedResult{
			JobID:     result.JobID,
			State:     Completed,
			JobResult: result.JobResult,
		}
	}
}
