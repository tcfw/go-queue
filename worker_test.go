package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	workPool := make(chan chan interface{}, 1)
	worker := NewWorker(workPool)

	assert.IsType(t, Worker{}, worker)
}

func TestReady(t *testing.T) {
	workPool := make(chan chan interface{}, 1)
	worker := NewWorker(workPool)

	assert.Empty(t, workPool)
	worker.Ready()
	assert.NotEmpty(t, workPool)
}

func TestStop(t *testing.T) {
	workPool := make(chan chan interface{}, 1)
	worker := NewWorker(workPool)

	worker.Stop()

	assert.NotEmpty(t, worker.quit)
}

type TestWorkerProcessor struct{}

func (p TestWorkerProcessor) Handle(job interface{}) {
	job.(struct {
		wait *sync.WaitGroup
	}).wait.Done()
}

func TestStart(t *testing.T) {
	workPool := make(chan chan interface{}, 1)
	worker := NewWorker(workPool)
	worker.Processor = &Processor{}

	go worker.Start()

	jobChan := <-workPool

	var wg sync.WaitGroup
	wg.Add(1)

	jobChan <- struct{ wait *sync.WaitGroup }{
		wait: &wg,
	}

	wg.Wait()

	worker.Stop()
}
