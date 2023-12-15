# ThreadPool

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/threadpool?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/threadpool)
![License](https://img.shields.io/github/license/mjwhitta/threadpool?style=for-the-badge)

## What is this?

This Go module is a very simple thread pool implementation. It may be
a little non-traditional. Rather than a pool of identical worker
threads waiting for data, this module sets the maximum number of
active threads. Additionally, each active thread can do a completely
different task, if desired.

**Note:** Versions prior to v1.5.0 caused deadlock if tasks created
sub-tasks.

## How to install

Open a terminal and run the following:

```
$ go get --ldflags "-s -w" --trimpath -u \
    github.com/mjwhitta/threadpool
```

## Usage

```
package main

import (
    "fmt"
    "time"

    tp "github.com/mjwhitta/threadpool"
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

- [Source](https://github.com/mjwhitta/threadpool)
