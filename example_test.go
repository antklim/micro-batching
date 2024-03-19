package microbatching_test

import (
	"fmt"
	"time"

	mb "github.com/antklim/micro-batching"
)

func ExampleService() {
	srv := mb.NewService(mb.WithFrequency(10 * time.Millisecond))
	srv.Run(&mockBatchProcessor{})

	jobsSize := 7
	testJobs := makeMockJobs(jobsSize)

	for _, j := range testJobs {
		srv.AddJob(j)
	}

	time.Sleep(20 * time.Millisecond)

	for _, j := range testJobs {
		r, err := srv.JobResult(j.ID())

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Job ID: %s, State: %s\n", r.JobID, r.State)
		}

	}

	srv.Shutdown()

	// Output:
	// Job ID: 0, State: Completed
	// Job ID: 1, State: Completed
	// Job ID: 2, State: Completed
	// Job ID: 3, State: Completed
	// Job ID: 4, State: Completed
	// Job ID: 5, State: Completed
	// Job ID: 6, State: Completed
}
