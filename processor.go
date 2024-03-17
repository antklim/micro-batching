package microbatching

import "time"

type ProcessProps struct {
	jobs      []Job
	startTime time.Time
}

type BatchProcessor interface {
	Process(props ProcessProps) []JobResult
}
