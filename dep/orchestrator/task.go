package orchestrator

import (
	"sync"
)

type task struct {
	err error
	job func() error
	errFunc func(err error) error
}

func newTask(job func() error, errFunc func(err error) error) *task {
	return &task{job: job, errFunc: errFunc}
}

func (t *task) runTask(wg *sync.WaitGroup) {
	t.err = t.job()
	wg.Done()
}

type taskPool struct {
	tasks       []*task
	wg          sync.WaitGroup
	taskChannel chan *task
}

func newTaskPool(tasks []*task) *taskPool {
	return &taskPool{
		tasks:       tasks,
		taskChannel: make(chan *task),
	}
}

func (tp *taskPool) runAll() {
	for i := 0; i < 4; i++ {
		go tp.do()
	}
	tp.wg.Add(len(tp.tasks))
	for _, t := range tp.tasks {
		tp.taskChannel <- t
	}
	// close channel because that's all the job
	close(tp.taskChannel)
	tp.wg.Wait()
}

func (tp *taskPool) do() {
	for t := range tp.taskChannel {
		t.runTask(&tp.wg)
	}
}
