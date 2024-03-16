package microbatching

type Job interface {
	Do() JobResult
}

type JobResult struct {
	Err    error
	Result interface{}
}
