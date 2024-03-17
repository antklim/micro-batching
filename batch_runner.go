package microbatching

import (
	"fmt"
	"time"
)

type batchRunner struct {
	processor BatchProcessor
	batchSize int
	frequency time.Duration
	jobs      <-chan job
	done      <-chan bool
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
			fmt.Println(">>> Abort ....")
			batcher.batch()

			return
		case <-ticker.C:
			fmt.Println(">>> Tick ....")
			batcher.batch()
		}
	}
}
