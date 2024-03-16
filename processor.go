package microbatching

type BatchProcessor interface {
	Process(jobs []Job) []JobResult
}
