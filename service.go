package microbatching

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/oklog/ulid/v2"
)

// ErrServiceClosed is returned by the [Service.AddJob] methods after a call to [Service.Shutdown].
var ErrServiceClosed = errors.New("microbatching: Service closed")

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

	jobs chan job
}

func NewService(processor BatchProcessor, opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	pDone := make(chan bool)

	go batchRunner(batchRunnerProps{
		batchSize: opts.batchSize,
		frequency: opts.frequency,
		processor: processor,
		done:      pDone,
	})

	return &Service{
		processor: processor,
		opts:      opts,
		jobs:      make(chan job, opts.queueSize),
		pDone:     pDone,
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
func (s *Service) AddJob(j Job) (ulid.ULID, error) {
	if s.shuttingDown() {
		return ulid.ULID{}, ErrServiceClosed
	}

	jobID := ulid.Make()
	newJob := job{ID: jobID, J: j}

	s.jobs <- newJob

	return jobID, nil
}

func (s *Service) shuttingDown() bool {
	return s.inShutdown.Load()
}

func (s *Service) Shutdown() {
	s.inShutdown.Store(true)
	s.pDone <- true
}
