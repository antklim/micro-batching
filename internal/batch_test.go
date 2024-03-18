package internal_test

import (
	"testing"

	mb "github.com/antklim/micro-batching/internal"
	"github.com/stretchr/testify/assert"
)

func TestBatchGroupsJobsIntoBatches(t *testing.T) {
	jobs := make(chan int)
	batches := make(chan []int)

	go func() {
		for i := 0; i < 11; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	go mb.Batch(3, jobs, batches)

	var bb [][]int
	for b := range batches {
		bb = append(bb, b)
	}

	assert.Equal(t, 4, len(bb))
	expectedBatches := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{9, 10},
	}
	assert.Equal(t, expectedBatches, bb)
}
