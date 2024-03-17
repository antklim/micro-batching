package microbatching

import (
	"fmt"
	"time"
)

type batcherProps struct {
	batchSize int
	jobs      <-chan job
	p         BatchProcessor
}

func batcher(props batcherProps) {
	batch := make([]Job, 0, props.batchSize)

	size := len(props.jobs)

	for i := 0; i < size; i++ {
		if len(batch) == props.batchSize {
			props.p.Process(ProcessProps{batch})
			batch = make([]Job, 0, props.batchSize)
		}

		j := <-props.jobs
		batch = append(batch, j.J)
	}

	if len(batch) > 0 {
		props.p.Process(ProcessProps{batch})
	}
}

type batchRunnerProps struct {
	processor BatchProcessor
	batchSize int
	frequency time.Duration
	jobs      <-chan job
	done      <-chan bool
}

func batchRunner(props batchRunnerProps) {
	ticker := time.NewTicker(props.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-props.done:
			fmt.Println(">>> Abort ....")
			batcher(batcherProps{
				batchSize: props.batchSize,
				jobs:      props.jobs,
				p:         props.processor,
			})

			return
		case <-ticker.C:
			fmt.Println(">>> Tick ....")
			batcher(batcherProps{
				batchSize: props.batchSize,
				jobs:      props.jobs,
				p:         props.processor,
			})
		}
	}
}
