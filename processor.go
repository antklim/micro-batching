package microbatching

import (
	"fmt"
	"time"
)

type ProcessProps struct {
	jobs      []Job
	startTime time.Time
}

type BatchProcessor interface {
	Process(props ProcessProps) []JobResult
}

type batcherProps struct {
	batchSize int
	jobs      <-chan job
	p         BatchProcessor
	t         time.Time
}

func batcher(props batcherProps) {
	batch := make([]Job, 0, props.batchSize)

	for i := 0; i < len(props.jobs); i++ {
		if len(batch) == props.batchSize {
			props.p.Process(ProcessProps{batch, props.t})
			batch = make([]Job, 0, props.batchSize)
		}

		j := <-props.jobs
		batch = append(batch, j.J)
	}

	if len(batch) > 0 {
		props.p.Process(ProcessProps{batch, props.t})
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
				t:         time.Now(),
			})

			return
		case t := <-ticker.C:
			fmt.Println(">>> Tick ....")
			batcher(batcherProps{
				batchSize: props.batchSize,
				jobs:      props.jobs,
				p:         props.processor,
				t:         t,
			})
		}
	}
}
