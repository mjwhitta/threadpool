package threadpool

import "sync"

// Task is a function pointer to be passed to Queue().
type Task func(threadId uint, data map[string]interface{})

// ThreadPool is a struct containing all relevant metadata to maintain
// a pool of threads.
type ThreadPool struct {
	pool chan uint
	wg   *sync.WaitGroup
}

// New will return a pointer to a new ThreadPool instance of the
// specified size.
func New(size uint) *ThreadPool {
	var i uint
	var tp = &ThreadPool{
		pool: make(chan uint, size),
		wg:   &sync.WaitGroup{},
	}

	for i = 0; i < size; i++ {
		tp.pool <- i + 1
	}

	return tp
}

// Queue will add a task to the ThreadPool.
func (tp *ThreadPool) Queue(task Task, scope map[string]interface{}) {
	// Notify that task is queued
	tp.wg.Add(1)

	go func(threadId uint, data map[string]interface{}) {
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
