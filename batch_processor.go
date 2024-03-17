package microbatching

type ProcessProps struct {
	jobs []Job
}

type BatchProcessor interface {
	Process(props ProcessProps) []JobResult
}
