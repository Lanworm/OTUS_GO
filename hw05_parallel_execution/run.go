package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(ch chan Task, errorCount *int64, wg *sync.WaitGroup) {
	for {
		t, ok := <-ch
		if !ok {
			break
		}
		wg.Add(1)
		err := t()
		if err != nil {
			atomic.AddInt64(errorCount, 1)
		}
		wg.Done()
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	var result error
	var errorCount int64
	tasksCh := make(chan Task)
	defer close(tasksCh)
	for i := 0; i < n; i++ {
		go worker(tasksCh, &errorCount, &wg)

	}
	for _, t := range tasks {
		if int(errorCount) != 0 && int(errorCount) == m {
			return ErrErrorsLimitExceeded
		} else {
			result = nil
		}
		tasksCh <- t
	}
	wg.Wait()
	return result
}
