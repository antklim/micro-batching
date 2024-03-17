package microbatching

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	jobs := make(chan job, 100)
	defer close(jobs)

	bp := &TestBP{}

	b := batcher{
		batchSize: 3,
		jobs:      jobs,
		p:         bp,
	}

	for i := 0; i < 11; i++ {
		jobs <- job{Job: &TestJob{}}
	}

	// make sure the jobs are in the queue before calling batch
	assert.Equal(t, 11, len(jobs))

	b.batch()

	assert.Equal(t, 0, len(jobs))

	// batch processor should have been called 4 times (3 + 3 + 3 + 2)
	assert.Equal(t, uint32(4), bp.Counter())
}
