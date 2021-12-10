package threadpool

import (
	"sync"

	"gitlab.com/mjwhitta/errors"
)

// Task is a function pointer to be passed to Queue().
type Task func(threadId int, data ThreadData)

// ThreadData is simple a map[string]interface{}.
type ThreadData map[string]interface{}

// ThreadPool is a struct containing all relevant metadata to maintain
// a pool of threads.
type ThreadPool struct {
	pool chan int
	wg   *sync.WaitGroup
}

// New will return a pointer to a new ThreadPool instance of the
// specified size.
func New(size int) (*ThreadPool, error) {
	var tp *ThreadPool
	var wg = &sync.WaitGroup{}

	if size <= 0 {
		return nil, errors.New("pool size must be greater than 0")
	}

	// Initialize ThreadPool
	tp = &ThreadPool{
		pool: make(chan int, size),
		wg:   &sync.WaitGroup{},
	}

	// Fill pool with workers
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(threadId int) {
			tp.pool <- threadId
			wg.Done()
		}(i + 1)
	}
	wg.Wait()

	return tp, nil
}

// Close will close the ThreadPool's chan preventing it from being
// used further.
func (tp *ThreadPool) Close() {
	close(tp.pool)
}

// Queue will add a task to the ThreadPool.
func (tp *ThreadPool) Queue(task Task, scope ThreadData) {
	// Notify that task is queued
	tp.wg.Add(1)

	go func(threadId int, data ThreadData) {
		// Run task with ready thread
		task(threadId, data)

		// Put thread back in pool
		tp.pool <- threadId

		// Notify when finished
		tp.wg.Done()
	}(<-tp.pool, scope) // Grab the next ready thread
}

// Wait will block until the ThreadPool has finished it's tasks.
func (tp *ThreadPool) Wait() {
	tp.wg.Wait()
}
