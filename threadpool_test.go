//nolint:godoclint // These are tests
package threadpool_test

import (
	"testing"

	tp "github.com/mjwhitta/threadpool"
	assert "github.com/stretchr/testify/require"
)

func TestEmptyPool(t *testing.T) {
	var e error

	_, e = tp.New(0)
	assert.Error(t, e)
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
	assert.NoError(t, e)
	assert.NotNil(t, pool)

	// Queue 10 tasks
	for i := range psz {
		pool.Queue(
			func(tid int, data tp.ThreadData) {
				if i, ok := data["int"].(int); ok {
					collector <- i
				}

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
	assert.Len(t, collected, psz)

	// Close pool
	pool.Close()
}
