package microbatching

import "time"

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
