# micro-batching
Micro batching library in Go

## Installation
```bash
go get github.com/antklim/micro-batching
```

## Usage
```go
package main

func main() {
    // create and run a micro-batching service
	srv := mb.NewService(mb.WithFrequency(10 * time.Millisecond))
	// provide a batch processor
    srv.Run(&BatchProcessor{})

	jobsSize := 7

    for i := 0; i < jobsSize; i++ {
        // add jobs to be processed
        srv.AddJob(Job{i})
    }

    // wait for the jobs to be processed
    time.Sleep(20 * time.Millisecond)

    // get the jobs results
	for i := 0; i < jobsSize; i++ {
		r, err := srv.JobResult(i)

		if err != nil {
			fmt.Println(err)
		} else {
            fmt.Printf("Job ID: %s, State: %s\n", r.JobID, r.State)
        }
	}

	srv.Shutdown()
}
```

## Service Options
The service can be configured with the following options:
- `WithBatchSize` - sets the batch size (default is 3)
- `WithFrequency` - sets the frequency of the calling the batch processor (default is 1 second)
- `WithLogger` - sets the logger for the service
- `WithQueueSize` - sets the queue size for the jobs (default is 100). In case the queue is full, the `AddJob` call will be blocked until the queue has space.
- `WithShutdownTimeout` - sets the timeout for the service shutdown (default is 5 seconds)

## API
The `NewService` function creates a new service with the provided options.

The service provides the following API:
- `AddJob` - adds a job to the queue
- `JobResult` - returns the result of the job by its ID
- `Run` - runs the service
- `Shutdown` - shuts down the service

## Architecture

## Testing

To test the library, run the following command:
```bash
go test -v
```