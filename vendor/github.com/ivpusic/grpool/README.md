# grpool
[![Build Status](https://travis-ci.org/ivpusic/grpool.svg?branch=master)](https://travis-ci.org/ivpusic/grpool)

Lightweight Goroutine pool

Clients can submit jobs. Dispatcher takes job, and sends it to first available worker.
When worker is done with processing job, will be returned back to worker pool.

Number of workers and Job queue size is configurable.

## Docs
https://godoc.org/github.com/ivpusic/grpool

## Installation
```
go get github.com/ivpusic/grpool
```

## Simple example
```Go
package main

import (
  "fmt"
  "runtime"
  "time"

  "github.com/ivpusic/grpool"
)

func main() {
  // number of workers, and size of job queue
  pool := grpool.NewPool(100, 50)

  // release resources used by pool
  defer pool.Release()

  // submit one or more jobs to pool
  for i := 0; i < 10; i++ {
    count := i

    pool.JobQueue <- func() {
      fmt.Printf("I am worker! Number %d\n", count)
    }
  }

  // dummy wait until jobs are finished
  time.Sleep(1 * time.Second)
}
```

## Example with waiting jobs to finish
```Go
package main

import (
  "fmt"
  "runtime"

  "github.com/ivpusic/grpool"
)

func main() {
  // number of workers, and size of job queue
  pool := grpool.NewPool(100, 50)
  defer pool.Release()

  // how many jobs we should wait
  pool.WaitCount(10)

  // submit one or more jobs to pool
  for i := 0; i < 10; i++ {
    count := i

    pool.JobQueue <- func() {
      // say that job is done, so we can know how many jobs are finished
      defer pool.JobDone()

      fmt.Printf("hello %d\n", count)
    }
  }

  // wait until we call JobDone for all jobs
  pool.WaitAll()
}
```

## License
*MIT*
