package errpool // import "github.com/stuarthicks/errpool"

import (
	"sync"
)

// Group is similar to a sync.ErrGroup, but contains a worker pool to constrain how many goroutines run concurrently processing work
type Group struct {
	NumberOfWorkers int
	jobsChan        chan func() error
	errChan         chan error
	wg              *sync.WaitGroup
	errors          []error
}

// StartWorkers begins a worker pool using the specified number of goroutines
func (g *Group) StartWorkers(n int) {
	g.wg = &sync.WaitGroup{}
	g.jobsChan = make(chan func() error)
	g.errChan = make(chan error)
	g.errors = make([]error, 0)
	for w := 1; w <= n; w++ {
		go g.worker()
	}
	go g.listenForErrors()
}

// Run adds a new job to the queue. Jobs are simply functions with no arguments, that return an error.
func (g *Group) Run(f func() error) {
	g.wg.Add(1)
	g.jobsChan <- f
}

// Wait blocks until all submitted jobs have completed, and any errors collected. It returns a slice of any errors returned by any jobs.
func (g *Group) Wait() []error {
	close(g.jobsChan)
	g.wg.Wait()
	close(g.errChan)
	return g.errors
}

func (g *Group) listenForErrors() {
	for {
		select {
		case err := <-g.errChan:
			g.errors = append(g.errors, err)
		}
	}
}

// worker is a dumb worker that executes job functions and sends their errors to a channel.
func (g *Group) worker() {
	for j := range g.jobsChan {
		defer g.wg.Done()
		if err := j(); err != nil {
			g.errChan <- err
		}
	}
}
