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

	jobs             chan job
	jobNotifications chan jobNotification
	jobResults       map[string]job
}

func NewService(processor BatchProcessor, opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	jobs := make(chan job, opts.queueSize)
	jobNotifications := make(chan jobNotification, opts.queueSize)
	pDone := make(chan bool)
	jobResults := make(map[string]job)

	br := batchRunner{
		batchSize:        opts.batchSize,
		frequency:        opts.frequency,
		processor:        processor,
		jobs:             jobs,
		jobNotifications: jobNotifications,
		done:             pDone,
	}

	go br.run()

	go func() {
		for jn := range jobNotifications {
			result := jobResults[jn.JobID]

			newResult := job{
				Job:   result.Job,
				State: jn.State,
			}

			if jn.State == Completed {
				newResult.Result = jn.Result
			}

			jobResults[jn.JobID] = newResult
		}
	}()

	return &Service{
		processor:        processor,
		opts:             opts,
		jobs:             jobs,
		jobNotifications: jobNotifications,
		jobResults:       jobResults,
		pDone:            pDone,
	}
}

func (s *Service) BatchSize() int {
	return s.opts.batchSize
}

func (s *Service) Frequency() time.Duration {
	return s.opts.frequency
}

// AddJob adds a job to the queue. It returns an error if the service is closed.
func (s *Service) AddJob(j Job) error {
	if s.shuttingDown() {
		return ErrServiceClosed
	}

	newJob := job{j, Submitted, JobResult{JobID: j.ID()}}

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
	close(s.jobNotifications)
}
