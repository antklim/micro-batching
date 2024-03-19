package microbatching

import "time"

// Batch takes values from an 'in' channel and sends them in batches to an 'out' channel.
func Batch[V any](batchSize int, in <-chan V, out chan<- []V, freq time.Duration) {
	ticker := time.NewTicker(freq)

	var batch []V

	for {
		select {
		case v, ok := <-in:
			if !ok {
				if len(batch) > 0 {
					out <- batch
				}

				ticker.Stop()

				return
			}

			batch = append(batch, v)

			if len(batch) == batchSize {
				out <- batch
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				out <- batch
				batch = nil
			}
		}
	}
}
