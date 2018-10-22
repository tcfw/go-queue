# Go-Queue

A multi-threaded chan based queue worker 

## License
Please refer to [LICENSE.md](https://github.com/tcfw/go-queue/LICENSE.md)

## Examples

### Simple example
```
package main 

import (
	queue "github.com/tcfw/go-queue"
)

type Processor struct {}

func (p *Processor) Handle(job interface{}) {
	//Handle job...
}

func main() {
	processor := &Processor{}

	dispatcher := queue.NewDispatcher(processor)
	dispatcher.Run()
}

```

### Specify number of workers
```
package main 

import (
	queue "github.com/tcfw/go-queue"
)

type Processor struct {}

func (p *Processor) Handle(job interface{}) {
	//Handle job...
}

func main() {
	processor := &Processor{}

	dispatcher := queue.NewDispatcher(processor)
	//20 workers will be created 
	dispatcher.MaxWorkers = 20
	dispatcher.Run()
}

```