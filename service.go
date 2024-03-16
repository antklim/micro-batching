package microbatching

import "time"

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
	processor BatchProcessor
	opts      serviceOptions
	jobs      []Job
}

func NewService(processor BatchProcessor, opt ...ServiceOption) *Service {
	opts := defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

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

// ServiceOption sets service options such as batch size and frequency.
type ServiceOption interface {
	apply(*serviceOptions)
}

type funcOption struct {
	f func(*serviceOptions)
}

func (fo *funcOption) apply(o *serviceOptions) {
	fo.f(o)
}

func newFuncOption(f func(*serviceOptions)) *funcOption {
	return &funcOption{f}
}

// WithBatchSize returns a ServiceOption that sets batch size.
func WithBatchSize(v int) ServiceOption {
	return newFuncOption(func(o *serviceOptions) {
		o.batchSize = v
	})
}

// WithBatchSize returns a ServiceOption that sets frequency.
func WithFrequency(v time.Duration) ServiceOption {
	return newFuncOption(func(o *serviceOptions) {
		o.frequency = v
	})
}
