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

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, j := range jobs {
				if err := j(); err != nil {
					ch <- fmt.Errorf("error running job %d: %v", i, err)
				}
			}
		}()
	}

	for _, path := range jobs {
		jobCh <- path
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
