package routines

import (
	"fmt"
	"runtime"
	"sync"
)

func StartWork(workers int, jobs []func() error) error {
	var wg sync.WaitGroup
	l := len(jobs)
	jobCh := make(chan func() error, l)
	ch := make(chan error, l)

	if workers < 1 {
		workers = runtime.NumCPU()
	}

	wg.Add(workers)
	for i := range workers {
		go func() {
			defer wg.Done()
			for job := range jobCh {
				if err := job(); err != nil {
					ch <- fmt.Errorf("error running job %d: %v", i, err)
				}
			}
		}()
	}

	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		if err != nil {
			return err
		}
	}

	return nil
}
