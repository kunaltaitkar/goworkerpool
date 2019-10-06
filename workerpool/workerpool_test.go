package workerpool

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var count int
var countMutex = &sync.Mutex{}
var cancelTriggered bool

//TestWorkerPoolWithContext - test for cancel trigger if count >= 500
func TestWorkerPoolWithContext(t *testing.T) {
	count = 0

	tasks := []*TaskWithContext{}

	for index := 0; index < 1000; index++ {
		tasks = append(tasks, NewTaskWithContext(index, incrementCount))
	}

	pool := NewPoolWithContext(tasks, 10)

	ticker := time.NewTicker(1 * time.Millisecond)

	go func() {
		for range ticker.C {
			if count > 500 {
				fmt.Println("cancelling tasks...")
				pool.Cancel()
				return
			}
		}
	}()

	pool.Run()

	assert.GreaterOrEqual(t, count, 500, "Count be greater than or equals to 500")

}

//TestWorkerpoolWithoutCancel - test without cancel trigger
func TestWorkerpoolWithoutCancel(t *testing.T) {
	count = 0

	tasks := []*TaskWithContext{}

	for index := 0; index < 1000; index++ {
		tasks = append(tasks, NewTaskWithContext(index, incrementCount))
	}

	pool := NewPoolWithContext(tasks, 10)

	pool.Run()
	fmt.Println("pool 2 completed....")
	assert.Equal(t, count, 1000, "Count should be equals to 1000")
}

//incrementCount- increment count by 1
func incrementCount(data interface{}) error {

	countMutex.Lock()
	count++
	countMutex.Unlock()

	return nil
}
