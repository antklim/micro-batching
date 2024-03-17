package microbatching

import (
	"sync/atomic"
	"time"
)

type batchRunner struct {
	processor BatchProcessor
	batchSize int
	frequency time.Duration
	jobs      <-chan job
	done      <-chan bool

	ticks atomic.Uint32
}

func (br *batchRunner) run() {
	batcher := batcher{
		batchSize: br.batchSize,
		jobs:      br.jobs,
		p:         br.processor,
	}

	ticker := time.NewTicker(br.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-br.done:
			batcher.batch()

			return
		case <-ticker.C:
			br.ticks.Add(1)
			batcher.batch()
		}
	}
}
