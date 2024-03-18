package internal

func Batch[V any](batchSize int, values <-chan V, batches chan<- []V) {
	var batch []V

	for v := range values {
		batch = append(batch, v)

		if len(batch) == batchSize {
			batches <- batch
			batch = nil
		}
	}

	if len(batch) > 0 {
		batches <- batch
	}

	close(batches)
}
