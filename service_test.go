package microbatching_test

import (
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	"github.com/stretchr/testify/assert"
)

type TestBP struct{}

func (bp *TestBP) Process(jobs []mb.Job) []mb.JobResult {
	return nil
}

var _ mb.BatchProcessor = (*TestBP)(nil)

func TestServiceInit(t *testing.T) {
	testCases := []struct {
		desc              string
		opts              []mb.ServiceOption
		expectedBatchSize int
		expectedFrequency time.Duration
	}{
		{
			desc:              "inits service with default options",
			expectedBatchSize: 3,
			expectedFrequency: time.Second,
		},
		{
			desc: "inits service with custom options",
			opts: []mb.ServiceOption{
				mb.WithBatchSize(15),
				mb.WithFrequency(5 * time.Minute),
			},
			expectedBatchSize: 15,
			expectedFrequency: 5 * time.Minute,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := mb.NewService(&TestBP{}, tC.opts...)

			assert.Equal(t, tC.expectedBatchSize, srv.BatchSize())
			assert.Equal(t, tC.expectedFrequency, srv.Frequency())
		})
	}
}
