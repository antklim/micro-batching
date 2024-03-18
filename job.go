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

type Job interface {
	ID() string
	Do() JobResult
}

type JobResult struct {
	JobID  string
	Err    error
	Result interface{}
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
