package worker

import (
	"sync/atomic"
	"testing"
)

func TestMain(t *testing.T) {
	var cnt int32
	fn1 := New()
	fn1.Start(func(p Payload) Result {
		return nil
	}, 2)

	go func() {
		for r := range fn1.Results() {
			atomic.AddInt32(&cnt, 1)
			r.Done()
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			fn1.Push(i)
		}
		fn1.Close()
	}()

	fn1.Wait()
	if cnt != 10000 {
		t.Fatal()
	}

}
