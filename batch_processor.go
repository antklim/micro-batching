package microbatching

type ProcessingResult struct {
	JobID string
	JobResult
}

type BatchProcessor interface {
	Process(jobs []Job) []ProcessingResult
}
