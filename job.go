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

type JobNotification struct {
	JobID string
	State JobState
	JobResult
}

type job struct {
	Job    Job
	State  JobState
	Result JobResult
}

type jobNotification struct {
	JobID  string
	State  JobState
	Result JobResult
}
