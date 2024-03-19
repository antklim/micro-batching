package microbatching

// Batch takes values from an 'in' channel and sends them in batches to an 'out' channel.
func Batch[V any](batchSize int, in <-chan V, out chan<- []V) {
	var batch []V

	for v := range in {
		batch = append(batch, v)

		if len(batch) == batchSize {
			out <- batch
			batch = nil
		}
	}

	if len(batch) > 0 {
		out <- batch
	}
}
