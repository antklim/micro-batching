package microbatching

import "github.com/oklog/ulid/v2"

type Job interface {
	Do() JobResult
}

type JobResult struct {
	Err    error
	Result interface{}
}

type job struct {
	ID ulid.ULID
	J  Job
}
