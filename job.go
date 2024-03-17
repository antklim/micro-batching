package microbatching

import "github.com/oklog/ulid/v2"

type JobState int

const (
	Submitted JobState = iota
	Processing
	Completed
)

func (s JobState) String() string {
	return [...]string{"Submitted", "Processing", "Completed"}[s]
}

type Job interface {
	Do() JobResult
}

type JobResult struct {
	Err    error
	Result interface{}
	State  JobState
}

type job struct {
	ID ulid.ULID
	J  Job
}
