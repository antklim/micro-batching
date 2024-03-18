package processor

import (
	"fmt"
	"sync"
	"time"
)

type ProcessResult struct {
	ID     int
	Result int
}

type BatchProcessor interface {
	Process(jobs []int) []ProcessResult
}

type Processor struct {
	sync.RWMutex

	batchProcessor BatchProcessor
	bChan          <-chan []int
	frequency      time.Duration
	batches        [][]int
}

func NewProcessor(bp BatchProcessor, batches <-chan []int, freq time.Duration) *Processor {
	return &Processor{
		batchProcessor: bp,
		bChan:          batches,
		frequency:      freq,
		batches:        make([][]int, 0),
	}
}

func (p *Processor) Start() {
	ticker := time.NewTicker(p.frequency)

	for {
		select {
		case batch, ok := <-p.bChan:
			fmt.Printf("Batch %v, %v\n", batch, ok)

			if !ok {
				ticker.Stop()
				return
			}

			p.Lock()
			p.batches = append(p.batches, batch)
			p.Unlock()
		case <-ticker.C:
			fmt.Println("Ticker ...")
			p.Lock()

			for _, batch := range p.batches {
				result := p.batchProcessor.Process(batch)
				fmt.Printf("Processed batch: %v\n", result)
			}

			p.batches = nil

			p.Unlock()
		}
	}
}
