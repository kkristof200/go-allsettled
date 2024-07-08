package allSettled

import "sync"

func AllSettled[T any](
	tasks []func() (*T, error),
	config *AllSettledConfig[T],
) []Result {
	if config == nil {
		config = &AllSettledConfig[T]{MaxConcurrency: -1}
	} else if config.MaxConcurrency == 0 {
		config.MaxConcurrency = -1
	}

	var wg sync.WaitGroup
	results := make([]Result, len(tasks))
	ch := make(chan struct {
		index  int
		result Result
	}, len(tasks))

	concurrencyLimit := config.MaxConcurrency
	semaphore := make(chan struct{}, len(tasks))
	if concurrencyLimit > 0 {
		semaphore = make(chan struct{}, concurrencyLimit)
	}

	wg.Add(len(tasks))
	for i, task := range tasks {
		go func(i int, task func() (*T, error)) {
			if concurrencyLimit > 0 {
				semaphore <- struct{}{}
			}

			defer func() {
				if concurrencyLimit > 0 {
					<-semaphore
				}
				wg.Done()
			}()

			value, err := task()
			success := err == nil
			result := Result{Success: success}
			ch <- struct {
				index  int
				result Result
			}{i, result}

			if success {
				if config.OnSuccess != nil {
					config.OnSuccess(i, SuccessfulResult[T]{Result: result, Value: value})
				}
			} else {
				if config.OnFailure != nil {
					config.OnFailure(i, FailedResult{Result: result, Reason: err})
				}
			}
		}(i, task)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		results[res.index] = res.result
	}

	return results
}
