# ThreadPool

<a href="https://www.buymeacoffee.com/mjwhitta">🍪 Buy me a cookie</a>

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/mjwhitta/threadpool)](https://goreportcard.com/report/gitlab.com/mjwhitta/threadpool)

## What is this?

This Go module is a very simple threadpool implementation.

## How to install

Open a terminal and run the following:

```
$ go get --ldflags="-s -w" --trimpath -u \
    gitlab.com/mjwhitta/threadpool
```

## Usage

```
package main

import (
    "fmt"
    "time"

    tp "gitlab.com/mjwhitta/threadpool"
)

func main() {
    var e error
    var pool *tp.ThreadPool

    // Create a threadpool with 10 worker threads
    if pool, e = tp.New(10); e != nil {
        panic(e)
    }
    defer pool.Close()

    // Queue 32 tasks
    for i := 0; i < 32; i++ {
        pool.Queue(
            func(tid int, data tp.ThreadData) {
                fmt.Printf("%d - %d\n", tid, data["counter"].(int))
                time.Sleep(time.Second)
            },
            tp.ThreadData{"counter": i},
        )
    }

    // Wait for tasks to complete
    pool.Wait()
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/threadpool)
