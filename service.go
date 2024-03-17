package microbatching

import (
	"errors"
	"sync/atomic"
	"time"
)

// ErrServiceClosed is returned by the [Service.AddJob] methods after a call to [Service.Shutdown].
var ErrServiceClosed = errors.New("microbatching: Service closed")

// ErrJobNotFound is returned by the [Service.JobResult] method when the job is not found.
var ErrJobNotFound = errors.New("microbatching: Job not found")

type serviceOptions struct {
	batchSize int
	frequency time.Duration
	queueSize int
}

var defaultOptions = serviceOptions{
	batchSize: 3,
	frequency: time.Second,
	queueSize: 100,
}

// Service is a micro-batching service that processes jobs in batches.
type Service struct {
	opts       serviceOptions
	inShutdown atomic.Bool

	processor BatchProcessor
	pDone     chan bool

	jobs       chan job
	jobResults map[string]job
}

func NewService(processor BatchProcessor, opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	jobs := make(chan job, opts.queueSize)
	pDone := make(chan bool)

	br := batchRunner{
		batchSize: opts.batchSize,
		frequency: opts.frequency,
		processor: processor,
		jobs:      jobs,
		done:      pDone,
	}

	go br.run()

	return &Service{
		processor:  processor,
		opts:       opts,
		jobs:       jobs,
		jobResults: make(map[string]job),
		pDone:      pDone,
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

// AddJob adds a job to the queue. It returns an error if the service is closed.
func (s *Service) AddJob(j Job) error {
	if s.shuttingDown() {
		return ErrServiceClosed
	}

	newJob := job{j, Submitted, JobResult{}}

	s.jobResults[j.ID()] = newJob
	s.jobs <- newJob

	return nil
}

// JobResult returns the result of a job. It returns an error if the job is not found.
func (s *Service) JobResult(jobID string) (JobState, JobResult, error) {
	result, ok := s.jobResults[jobID]

	if !ok {
		return Submitted, JobResult{}, ErrJobNotFound
	}

	return result.State, result.Result, nil
}

func (s *Service) shuttingDown() bool {
	return s.inShutdown.Load()
}

// Shutdown stops the service.
func (s *Service) Shutdown() {
	if s.shuttingDown() {
		return
	}

	s.inShutdown.Store(true)
	s.pDone <- true

	// Should wait until all jobs in the queue are processed.

	close(s.jobs)
}
