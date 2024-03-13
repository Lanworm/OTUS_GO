package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func runTask(t Task, errorCh chan<- struct{}, successCh chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	err := t()
	if err != nil {
		errorCh <- struct{}{}
	} else {
		successCh <- struct{}{}
	}
}

func taskLimiter(tasks []Task, done chan<- error, wg *sync.WaitGroup, maxErrorCount int, maxGoCount int) {
	errorCh := make(chan struct{}, maxGoCount*2)
	successCh := make(chan struct{}, maxGoCount*2)
	var taskWg sync.WaitGroup
	defer func() {
		taskWg.Wait()
		close(errorCh)
		close(successCh)
		wg.Done()
	}()
	nbTask := 0
	successCount := 0
	errorCount := 0
	taskSize := len(tasks)
	for i := 0; i < maxGoCount && i < taskSize; i++ {
		taskWg.Add(1)
		go runTask(tasks[nbTask], errorCh, successCh, &taskWg)
		nbTask++
	}
	for {
		select {
		case <-errorCh:
			errorCount++
		case <-successCh:
			successCount++
		}

		if errorCount == maxErrorCount {
			done <- ErrErrorsLimitExceeded
			return
		}

		if errorCount+successCount == taskSize {
			done <- nil
			return
		}

		if nbTask < taskSize {
			taskWg.Add(1)
			go runTask(tasks[nbTask], errorCh, successCh, &taskWg)
			nbTask++
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	wg.Add(1)
	resultChan := make(chan error)
	defer func() {
		close(resultChan)
	}()
	go taskLimiter(tasks, resultChan, &wg, m, n)
	x := <-resultChan
	wg.Wait()
	return x
}
