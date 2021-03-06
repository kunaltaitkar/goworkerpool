package main

import (
	"fmt"
	"goworkerpool/workerpool"
	"sync"
	"time"
)

var count int
var countMutex = &sync.Mutex{}

func main() {

	count = 0

	// preparing tasks
	tasks := []*workerpool.Task{}

	for index := 0; index < 1000; index++ {
		tasks = append(tasks, workerpool.NewTask(index, incrementCount))
	}

	//initialize pool
	pool := workerpool.NewPoolWithContext(tasks, 10)

	// to check count value
	ticker := time.NewTicker(1 * time.Millisecond)

	// cancel all workers when count is more than 500
	go func() {
		for range ticker.C {
			if count > 500 {
				fmt.Println("cancelling tasks...")
				pool.Cancel()
				return
			}
		}
	}()

	// run pool
	pool.Run()

	time.Sleep(10 * time.Second)
}

//incrementCount- increment count by 1
func incrementCount(data interface{}) error {

	countMutex.Lock()
	count++
	countMutex.Unlock()

	return nil
}
