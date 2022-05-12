package workerpool

import "sync"

type WorkerPool[T any] struct {
	funcs []func() (T, error)
	err   error
	size  int
	cof   bool
}

// New creates a new WorkerPool with the given options.
// *size*: maximum worker to be run at once
// *continueOnFail*: if false, will stop sending new worker to the pool when one returns an error
func New[T any](size int, continueOnFail bool) *WorkerPool[T] {
	return &WorkerPool[T]{
		size: size,
		cof:  continueOnFail,
	}
}

// Add adds a new worker to the pool.
// Added workers are not yet executed.
func (wp *WorkerPool[T]) Add(f func() (T, error)) {
	wp.funcs = append(wp.funcs, f)
}

// Err returns the last reported error
func (wp *WorkerPool[_]) Err() error {
	return wp.err
}

// Exec executes the workers in dedicated go routines.
// Returns a channel which will receive the worker's results.
// Channel is closed when the last worker finishes.
func (wp *WorkerPool[T]) Exec() <-chan T {
	poolC := make(chan bool, wp.size)
	resultC := make(chan T, len(wp.funcs))
	wg := &sync.WaitGroup{}
	stop := false
	wg.Add(1) // Add 1 first to prevent Wait() to exit before the first routine starts
	go func() {
		defer wg.Done()
		for _, f := range wp.funcs {
			f := f
			if !wp.cof && stop {
				return
			}
			wg.Add(1)
			poolC <- true
			go func() {
				defer func() {
					wg.Done()
					<-poolC
				}()
				v, err := f()
				if err != nil {
					stop = true
					wp.err = err
					return
				}
				resultC <- v
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(resultC)
	}()
	return resultC
}
