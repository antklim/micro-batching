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

// JobResult describes job result.
type JobResult struct {
	Err    error
	Result interface{}
}

// Job describes a job interface. Job is a unit of work executed by the batch processor.
type Job interface {
	ID() string
	Do() JobResult
}

// JobExtendedResult describes job result with state.
type JobExtendedResult struct {
	JobID string
	State JobState
	JobResult
}
