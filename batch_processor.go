package microbatching

// ProcessingResult describes job processing result by the batch processor.
type ProcessingResult struct {
	JobID string
	JobResult
}

// BatchProcessor describes a batch processor interface.
type BatchProcessor interface {
	Process(jobs []Job) []ProcessingResult
}
