package threadpool_test

import (
	"testing"

	tp "gitlab.com/mjwhitta/threadpool"
)

func TestEmptyPool(t *testing.T) {
	var e error
	var err string = "threadpool: pool size must be greater than 0"
	if _, e = tp.New(0); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", err)
	} else if e.Error() != err {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), err)
	}
}

func TestPool(t *testing.T) {
	var collected = map[int]struct{}{}
	var collector chan int
	var e error
	var pool *tp.ThreadPool
	var psz int = 10

	// Create chan
	collector = make(chan int, psz)

	// Create ppool
	if pool, e = tp.New(psz); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	// Queue 10 tasks
	for i := 0; i < psz; i++ {
		pool.Queue(
			func(tid int, data tp.ThreadData) {
				collector <- data["int"].(int)
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
	if len(collected) != psz {
		t.Errorf("\ngot: %d\nwant: %d", len(collected), psz)
	}
}
