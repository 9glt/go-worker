package worker

import (
	"sync"
)

// New pool
func New() *Pool {
	mu := &sync.Mutex{}
	p := &Pool{
		make(chan Payload), make(chan struct{}), mu, sync.NewCond(mu), make(chan Results), &sync.WaitGroup{},
	}
	p.wg.Add(1)
	return p
}

type Results struct {
	Result Result
	p      *Pool
}

func (r *Results) Done() {
	r.p.wg.Done()
}

// Payload interface
type Payload interface{}

// Result result
type Result interface{}

// Pool .
type Pool struct {
	ch       chan Payload
	wait     chan struct{}
	mu       *sync.Mutex
	cond     *sync.Cond
	resultCh chan Results
	wg       *sync.WaitGroup
}

// Start function
func (p *Pool) Start(fn func(Payload) Result, concurrency int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			for payload := range p.ch {
				p.resultCh <- Results{fn(payload), p}
			}
		}()
	}
	wg.Wait()
}

// Push payload
func (p *Pool) Push(pl Payload) {
	p.wg.Add(1)
	p.ch <- pl
}

// Close queue
func (p *Pool) Close() {
	p.wg.Done()
}

// Wait to finish
func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.ch)
	close(p.resultCh)
}

// Results channel
func (p *Pool) Results() chan Results {
	return p.resultCh
}
