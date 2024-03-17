package microbatching

type batcher struct {
	batchSize int
	jobs      <-chan job
	p         BatchProcessor
}

// batch reads the jobs from the queue, groups them in batches of the defined size and sends them to the processor.
func (b *batcher) batch() {
	batch := make([]Job, 0, b.batchSize)

	queueSize := len(b.jobs)

	for i := 0; i < queueSize; i++ {
		if len(batch) == b.batchSize {
			b.p.Process(ProcessProps{batch})
			batch = make([]Job, 0, b.batchSize)
		}

		j := <-b.jobs
		batch = append(batch, j.J)
	}

	if len(batch) > 0 {
		b.p.Process(ProcessProps{batch})
	}
}
