package microbatching

import (
	"fmt"
	"time"
)

type batchRunnerProps struct {
	processor BatchProcessor
	batchSize int
	frequency time.Duration
	jobs      <-chan job
	done      <-chan bool
}

func batchRunner(props batchRunnerProps) {
	batcher := batcher{
		batchSize: props.batchSize,
		jobs:      props.jobs,
		p:         props.processor,
	}

	ticker := time.NewTicker(props.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-props.done:
			fmt.Println(">>> Abort ....")
			batcher.batch()

			return
		case <-ticker.C:
			fmt.Println(">>> Tick ....")
			batcher.batch()
		}
	}
}
