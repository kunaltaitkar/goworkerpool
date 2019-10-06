package workerpool

import (
	"context"
)

// PoolWithContext is a worker group that runs a number of tasks at a
// configured concurrency.
type PoolWithContext struct {
	Tasks          []*Task
	concurrency    int
	tasksChan      chan *Task
	CancelHandler  *cancelHandler
	IsPoolCanceled bool
}

type cancelHandler struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// NewPoolWithContext initializes a new pool with the given tasks and
// at the given concurrency.
func NewPoolWithContext(tasks []*Task, concurrency int) *PoolWithContext {
	cntx, cancelFunction := context.WithCancel(context.Background())

	obj := cancelHandler{}

	obj.Ctx = cntx
	obj.CancelFunc = cancelFunction

	return &PoolWithContext{
		Tasks:         tasks,
		concurrency:   concurrency,
		tasksChan:     make(chan *Task),
		CancelHandler: &obj,
	}
}

// Run runs all work within the pool
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
type Task struct {
	// Err holds an error that occurred during a task. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Err error

	Data interface{}

	f func(data interface{}) error
}

// NewTaskWithContext initializes a new task based on a given work
// function.
func NewTask(d interface{}, f func(data interface{}) error) *Task {
	return &Task{Data: d, f: f}
}

// Run runs a Task
func (t *Task) Run(p *PoolWithContext) {

	for {
		select {

		case <-p.CancelHandler.Ctx.Done():

			//return from all tasks when cancel handler calls
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
