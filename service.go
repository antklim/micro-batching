package microbatching

import (
	"errors"
	"sync/atomic"
	"time"
)

// ErrServiceClosed is returned by the [Service.AddJob] methods after a call to [Service.Shutdown].
var ErrServiceClosed = errors.New("microbatching: Service closed")

type serviceOptions struct {
	batchSize int
	frequency time.Duration
}

var defaultOptions = serviceOptions{
	batchSize: 3,
	frequency: time.Second,
}

// Service is a micro-batching service that processes jobs in batches.
type Service struct {
	processor  BatchProcessor
	opts       serviceOptions
	inShutdown atomic.Bool

	jobs []Job
}

func NewService(processor BatchProcessor, opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	// TODO: start ticker with the frequency

	return &Service{
		processor: processor,
		opts:      opts,
		jobs:      make([]Job, 0, opts.batchSize),
	}
}

func (s *Service) BatchSize() int {
	return s.opts.batchSize
}

func (s *Service) Frequency() time.Duration {
	return s.opts.frequency
}

// JobsQueueSize returns the number of jobs in the queue.
func (s *Service) JobsQueueSize() int {
	return len(s.jobs)
}

// AddJob adds a job to the queue.
func (s *Service) AddJob(job Job) (JobResult, error) {
	if s.shuttingDown() {
		return JobResult{}, ErrServiceClosed
	}

	s.jobs = append(s.jobs, job)
	return JobResult{}, nil
}

func (s *Service) shuttingDown() bool {
	return s.inShutdown.Load()
}

func (s *Service) Shutdown() {
	s.inShutdown.Store(true)
}
