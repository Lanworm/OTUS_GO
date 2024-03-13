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

func taskLimiter(tasks []Task, done chan<- error, maxErrorCount int, maxGoCount int) {
	errorCh := make(chan struct{}, maxGoCount*2)
	successCh := make(chan struct{}, maxGoCount*2)
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(errorCh)
		close(successCh)
	}()
	nbTask := 0
	successCount := 0
	errorCount := 0
	taskSize := len(tasks)
	for i := 0; i < maxGoCount && i < taskSize; i++ {
		wg.Add(1)
		go runTask(tasks[nbTask], errorCh, successCh, &wg)
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
			wg.Add(1)
			go runTask(tasks[nbTask], errorCh, successCh, &wg)
			nbTask++
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	resultChan := make(chan error)
	defer func() {
		close(resultChan)
	}()
	go taskLimiter(tasks, resultChan, m, n)
	x := <-resultChan
	return x
}
