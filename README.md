# Go-Queue

[![GoDoc](https://godoc.org/github.com/tcfw/go-queue?status.svg)](https://godoc.org/github.com/tcfw/go-queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/tcfw/go-queue)](https://goreportcard.com/report/github.com/tcfw/go-queue)


A simple multi-threaded chan based queue worker 

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