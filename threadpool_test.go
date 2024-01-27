package threadpool_test

import (
	"testing"

	tp "github.com/mjwhitta/threadpool"
	assert "github.com/stretchr/testify/require"
)

func TestEmptyPool(t *testing.T) {
	var e error

	_, e = tp.New(0)
	assert.NotNil(t, e)
}

func TestPool(t *testing.T) {
	var collected map[int]struct{} = map[int]struct{}{}
	var collector chan int
	var e error
	var pool *tp.ThreadPool
	var psz int = 10

	// Create chan
	collector = make(chan int, psz)

	// Create ppool
	pool, e = tp.New(psz)
	assert.Nil(t, e)
	assert.NotNil(t, pool)

	// Queue 10 tasks
	for i := 0; i < psz; i++ {
		pool.Queue(
			func(tid int, data tp.ThreadData) {
				collector <- data["int"].(int)

				// Queue an additional task to test for deadlock
				pool.Queue(func(tid int, data tp.ThreadData) {}, nil)
			},
			tp.ThreadData{"int": i},
		)
	}

	// Wait for tasks to complete
	pool.Wait()

	// Close chan
	close(collector)

	// Collect data
	for i := range collector {
		collected[i] = struct{}{}
	}

	// Validate they are all unique values
	assert.Equal(t, psz, len(collected))
}
