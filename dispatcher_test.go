package queue

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Job struct {
	msg string
}

func (j *Job) getMessage() string {
	return j.msg
}

type Processor struct {
	afterFunc func()
}

func (p Processor) Handle(job interface{}) {
	j := job.(struct {
		wait *sync.WaitGroup
	})
	j.wait.Done()

	if p.afterFunc != nil {
		p.afterFunc()
	}
}

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(&Processor{})
	assert.IsType(t, &Dispatcher{}, dispatcher)
}

func TestRun(t *testing.T) {
	t.Run("1 worker", func(t *testing.T) {
		dispatcher := NewDispatcher(&Processor{})
		dispatcher.MaxWorkers = 1
		dispatcher.Run()
		assert.Equal(t, len(dispatcher.Workers), 1)
	})

	t.Run("N workers", func(t *testing.T) {
		dispatcher := NewDispatcher(&Processor{})
		dispatcher.MaxWorkers = 0
		dispatcher.Run()
		assert.Equal(t, len(dispatcher.Workers), runtime.NumCPU())
	})
}

func TestDispatcherStop(t *testing.T) {
	dispatcher := NewDispatcher(&Processor{})
	dispatcher.MaxWorkers = 1
	dispatcher.Run()

	dispatcher.Stop()
	assert.Nil(t, dispatcher.WorkerPool)
}

func TestDispatch(t *testing.T) {
	dispatcher := NewDispatcher(&Processor{})
	dispatcher.MaxWorkers = 1
	dispatcher.Run()

	assert.Equal(t, len(dispatcher.Workers), 1)

	var wg sync.WaitGroup
	wg.Add(1)

	dispatcher.Queue(struct {
		wait *sync.WaitGroup
	}{
		wait: &wg,
	})

	wg.Wait()

	dispatcher.Stop()
}

func TestQueueAfter(t *testing.T) {
	p := &Processor{}

	dispatcher := NewDispatcher(p)
	dispatcher.MaxWorkers = 1
	dispatcher.Run()

	assert.Equal(t, len(dispatcher.Workers), 1)

	var wg sync.WaitGroup
	wg.Add(1)

	sTime := time.Now()
	var tDiff time.Duration
	p.afterFunc = func() {
		tDiff = time.Since(sTime)
	}

	dispatcher.QueueAfter(struct {
		wait *sync.WaitGroup
	}{
		wait: &wg,
	}, 2*time.Millisecond)

	wg.Wait()

	assert.Equal(t, tDiff.Milliseconds(), int64(2))
}

func DispatcherNWorkers(b *testing.B, workers int) {
	dispatcher := NewDispatcher(&Processor{})
	dispatcher.MaxWorkers = workers
	dispatcher.Run()

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		dispatcher.Queue(struct {
			wait *sync.WaitGroup
		}{
			wait: &wg,
		})
	}

	wg.Wait()

	dispatcher.Stop()
}

func BenchmarkDispatcher(b *testing.B) {
	for n := 1; n < 10; n += 2 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			DispatcherNWorkers(b, n)
		})
	}
}
