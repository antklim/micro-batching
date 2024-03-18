package microbatching

// type batcher struct {
// 	batchSize        int
// 	jobs             <-chan job
// 	jobNotifications chan<- jobNotification
// 	p                BatchProcessor
// }

// // batch reads the jobs from the queue, groups them in batches of the defined size and sends them to the processor.
// func (b *batcher) batch() {
// 	batch := make([]Job, 0, b.batchSize)

// 	queueSize := len(b.jobs)

// 	for i := 0; i < queueSize; i++ {
// 		if len(batch) == b.batchSize {
// 			jobResults := b.p.Process(batch)

// 			for _, r := range jobResults {
// 				b.jobNotifications <- jobNotification{
// 					JobID:  r.JobID,
// 					State:  Completed,
// 					Result: r,
// 				}
// 			}

// 			batch = make([]Job, 0, b.batchSize)
// 		}

// 		j := <-b.jobs

// 		b.jobNotifications <- jobNotification{
// 			JobID: j.Job.ID(),
// 			State: Processing,
// 		}

// 		batch = append(batch, j.Job)
// 	}

// 	if len(batch) > 0 {
// 		jobResults := b.p.Process(batch)

// 		for _, r := range jobResults {
// 			b.jobNotifications <- jobNotification{
// 				JobID:  r.JobID,
// 				State:  Completed,
// 				Result: r,
// 			}
// 		}
// 	}
// }
