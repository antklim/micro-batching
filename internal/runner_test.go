package internal_test

import (
	"strconv"
	"testing"
	"time"

	mb "github.com/antklim/micro-batching"
	internal "github.com/antklim/micro-batching/internal"
)

func TestRunner(t *testing.T) {
	batches := make(chan []mb.Job)

	runner := internal.NewRunner(&mockBatchProcessor{}, batches, nil, 10*time.Millisecond)

	go func() {
		for i := 0; i < 11; i++ {
			batches <- []mb.Job{
				newMockJob(strconv.Itoa(i)),
				newMockJob(strconv.Itoa(i + 1)),
				newMockJob(strconv.Itoa(i + 2)),
			}
		}

		time.Sleep(50 * time.Millisecond)

		for i := 11; i < 22; i++ {
			batches <- []mb.Job{
				newMockJob(strconv.Itoa(i)),
				newMockJob(strconv.Itoa(i + 1)),
				newMockJob(strconv.Itoa(i + 2)),
			}
		}

		time.Sleep(50 * time.Millisecond)

		close(batches)
	}()

	go runner.Run()

	time.Sleep(200 * time.Millisecond)
}
