package microbatching

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// ErrServiceClosed is returned by the [Service.AddJob] methods after a call to [Service.Shutdown].
var ErrServiceClosed = errors.New("microbatching: Service closed")

// ErrJobNotFound is returned by the [Service.JobResult] method when the job is not found.
var ErrJobNotFound = errors.New("microbatching: Job not found")

type serviceOptions struct {
	batchSize       int
	frequency       time.Duration
	queueSize       int
	shutdownTimeout time.Duration
}

var defaultOptions = serviceOptions{
	batchSize:       3,
	frequency:       time.Second,
	queueSize:       100,
	shutdownTimeout: 5 * time.Second,
}

// Service is a micro-batching service that processes jobs in batches.
type Service struct {
	opts       serviceOptions
	inShutdown atomic.Bool

	jobs          chan Job
	batches       chan []Job
	notifications chan JobExtendedResult
	done          chan bool

	jobResults map[string]JobExtendedResult
}

func NewService(opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	return &Service{
		opts:          opts,
		jobs:          make(chan Job),
		batches:       make(chan []Job),
		notifications: make(chan JobExtendedResult),
		done:          make(chan bool),
		jobResults:    make(map[string]JobExtendedResult),
	}
}

func (s *Service) Run(bp BatchProcessor) {
	runner := NewRunner(bp, s.batches, s.notifications, s.opts.frequency)

	// group jobs into batches
	go Batch(s.opts.batchSize, s.jobs, s.batches, s.opts.frequency)

	// runs batches
	go runner.Run()

	// collect notifications
	go func() {
		for n := range s.notifications {
			s.jobResults[n.JobID] = n
		}

		s.done <- true
	}()
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

	s.jobResults[j.ID()] = JobExtendedResult{
		JobID: j.ID(),
		State: Submitted,
	}
	s.jobs <- j

	return nil
}

// JobResult returns the result of a job. It returns an error if the job is not found.
func (s *Service) JobResult(jobID string) (JobExtendedResult, error) {
	result, ok := s.jobResults[jobID]

	if !ok {
		return JobExtendedResult{}, ErrJobNotFound
	}

	return result, nil
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
	close(s.jobs)
	close(s.batches)
	close(s.notifications)

	select {
	case <-s.done:
		return
	case <-time.After(s.opts.shutdownTimeout):
		fmt.Println("microbatching: service shutdown timeout")
		return
	}
}
