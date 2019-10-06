package workerpool

import (
	"context"
	"fmt"
)

// PoolWithContext is a worker group that runs a number of tasks at a
// configured concurrency.
type PoolWithContext struct {
	Tasks          []*TaskWithContext
	concurrency    int
	tasksChan      chan *TaskWithContext
	CancelHandler  *cancelHandler
	IsPoolCanceled bool
}

type cancelHandler struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// NewPoolWithContext initializes a new pool with the given tasks and
// at the given concurrency.
func NewPoolWithContext(tasks []*TaskWithContext, concurrency int) *PoolWithContext {
	cntx, cancelFunction := context.WithCancel(context.Background())

	obj := cancelHandler{}

	obj.Ctx = cntx
	obj.CancelFunc = cancelFunction

	return &PoolWithContext{
		Tasks:         tasks,
		concurrency:   concurrency,
		tasksChan:     make(chan *TaskWithContext),
		CancelHandler: &obj,
		// isTaskChannelClosed: make(chan bool),
		// isTaskChannelClosed: make(chan bool, len(tasks)),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *PoolWithContext) Run() {

	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	for _, task := range p.Tasks {

		if p.IsPoolCanceled {
			return
		} else {
			p.tasksChan <- task
		}

	}

	close(p.tasksChan)

}

// The work loop for any single goroutine.
func (p *PoolWithContext) work() {
	for task := range p.tasksChan {
		task.Run(p)
	}
}

// TaskWithContext encapsulates a work item that should go in a work
// pool.
type TaskWithContext struct {
	// Err holds an error that occurred during a task. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Err error

	Data interface{}

	f func(data interface{}) error
}

// NewTaskWithContext initializes a new task based on a given work
// function.
func NewTaskWithContext(d interface{}, f func(data interface{}) error) *TaskWithContext {
	return &TaskWithContext{Data: d, f: f}
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (t *TaskWithContext) Run(p *PoolWithContext) {

	for {
		select {

		case <-p.CancelHandler.Ctx.Done():

			fmt.Println("Exist Task:", t.Data)
			// fmt.Println("after channel")
			// close(p.tasksChan)
			// p.isTaskChannelClosed <- true
			return

		default:

			t.Err = t.f(t.Data)

			return
		}
	}

}

//Cancel all tasks
func (p *PoolWithContext) Cancel() {
	if p != nil {
		p.CancelHandler.CancelFunc()
		p.IsPoolCanceled = true
	}
}
