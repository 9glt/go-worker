package worker

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// New pool
func New() *Pool {
	mu := &sync.Mutex{}
	return &Pool{
		0, make(chan Payload), make(chan struct{}), mu, sync.NewCond(mu), make(chan Results),
	}
}

type Results struct {
	Result Result
	p      *Pool
}

func (r *Results) Done() {
	atomic.AddInt32(&r.p.cnt, -1)
}

// Payload interface
type Payload interface{}

// Result result
type Result interface{}

// Pool .
type Pool struct {
	cnt      int32
	ch       chan Payload
	wait     chan struct{}
	mu       *sync.Mutex
	cond     *sync.Cond
	resultCh chan Results
}

// Start function
func (p *Pool) Start(fn func(Payload) Result, concurrency int) {
	var wg sync.WaitGroup
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
	atomic.AddInt32(&p.cnt, 1)
	p.ch <- pl
}

// Close queue
func (p *Pool) Close() {
	atomic.AddInt32(&p.cnt, -1)
}

// Wait to finish
func (p *Pool) Wait() {
	for atomic.LoadInt32(&p.cnt) > -1 {
		runtime.Gosched()
	}
	close(p.ch)
	close(p.resultCh)
}

// Results channel
func (p *Pool) Results() chan Results {
	return p.resultCh
}
