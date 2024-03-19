package microbatching

type JobState int

const (
	Submitted JobState = iota
	Processing
	Completed
)

func (s JobState) String() string {
	return [...]string{"Submitted", "Processing", "Completed"}[s]
}

type JobResult struct {
	Err    error
	Result interface{}
}

type Job interface {
	ID() string
	Do() JobResult
}

type JobExtendedResult struct {
	JobID string
	State JobState
	JobResult
}
