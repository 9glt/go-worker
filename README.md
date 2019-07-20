# go-worker
Simple goroutines worker pool

```go
package main

import (
	"fmt"
	goworker "github.com/9glt/go-worker"
)

func main() {

    fn1 := goworker.New()
    
    // start 32 go routines
	fn1.Start(func(p goworker.Payload) goworker.Result {
		return fmt.Sprintf("test payload %d", p)
	}, 32)

    // read results
	go func() {
		for res := range fn1.Results() {
			println(res.(string))
		}
    }()
    
    // push results and close after done
	go func() {
		for i := 0; i < 100; i++ {
			fn1.Push(i)
		}
		fn1.Close()
    }()
    
	fn1.Wait()
}
```