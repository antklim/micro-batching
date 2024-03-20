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

```
+---------+  jobs channel  +-------+  batches channel  +-----------------+
| Service | -------------> | Batch | ----------------> | Batch processor |
+---------+                +-------+                   +-----------------+
    ^                                                           |
    |                  notifications channel                    |
    +-----------------------------------------------------------+
```

### Processing jobs

When `AddJob` is called, the job is added to the jobs channel. The batching go-routine listens to the jobs channel and creates batches of jobs. When the batch is ready or when the frequency time ticks, the batches are sent to the batch runner go-routine.


The batch runner listens to the batch channel and accumulates the batches. When the frequency timer ticks the batches are sent to the batch processor. The results of the processed batches are sent to the notifications channel.


The service listens to the notifications channel and updates the job results.

### Job results
Job results are available via the `JobResult` method.

### Service shutdown

When the `Shutdown` method is called, the service closes the jobs channel and waits for the submitted jobs to be processed.


Closing jobs channel causes batch go-routine to send all the remaining batches to the batch runner. The batch go-routine closes the batches channel.


Closing the batches channel causes the batch runner to send all the remaining batches to the batch processor. The batch runner closes the notifications channel.


Closing the notifications channel causes the service to stop listening to the notifications channel and send the final event to the `done` channel.


The `Shutdown` method is waiting for the event in `done` channel or for the timeout to occur.

## Testing

To test the library, run the following command:
```bash
go test -v
```