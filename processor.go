package microbatching

import (
	"time"
)

type ProcessProps struct {
	jobs      []Job
	startTime time.Time
}

type BatchProcessor interface {
	Process(props ProcessProps) []JobResult
}

type batchRunnerProps struct {
	processor BatchProcessor
	batchSize int
	frequency time.Duration
	done      <-chan bool
}

func batchRunner(props batchRunnerProps) {
	ticker := time.NewTicker(props.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-props.done:
			return
		case t := <-ticker.C:
			props.processor.Process(ProcessProps{nil, t})
		}
	}
}
