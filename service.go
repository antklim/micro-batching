package microbatching

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/oklog/ulid/v2"
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
	jobResults map[ulid.ULID]JobResult
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
		jobResults: make(map[ulid.ULID]JobResult),
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
func (s *Service) AddJob(j Job) (ulid.ULID, error) {
	if s.shuttingDown() {
		return ulid.ULID{}, ErrServiceClosed
	}

	jobID := ulid.Make()
	newJob := job{ID: jobID, J: j}

	s.jobResults[jobID] = JobResult{}
	s.jobs <- newJob

	return jobID, nil
}

// JobResult returns the result of a job. It returns an error if the job is not found.
func (s *Service) JobResult(jobID ulid.ULID) (JobResult, error) {
	result, ok := s.jobResults[jobID]

	if !ok {
		return JobResult{}, ErrJobNotFound
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
	s.pDone <- true

	close(s.jobs)
}
