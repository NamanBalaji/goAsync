# goAsync

This package provides workers that help process tasks and job asynchronously .

## Installation
```
go get github.com/NamanBalaji/goAsync
```
## Usage
```go
package main

import (
	"log"
	"github.com/NamanBalaji/goAsync"
	"context"
)
type mockTask struct {
	Data int
} 

func main() {
	ctx := context.Background()
	# create worker 
    worker, err := goAsync.NewAsyncTask()
    if err != nil {
        log.fatal(err)
    }
    # process task asynchronously
    err = worker.AddTask(ctx, &mockTask{Data: 100})
    if err != nil {
        log.fatal(err)
    }
}
```

### Changing the Parameters

When creating a new `async task` different parameters can be passed to the `NewAsyncTask()` function. The parameters are:

* `WithQueueSizeOption(size)` — Allows you set the queue size.
* `WithWorkerSizeOption(size)` — Allows you set the worker size.
* `WithTimeoutOption(time)` — Allows you set the timeout value. The value should be of type `time`
* `WithErrorHandlerOption(func)` — Allows you to set a custom handler for errors. The custom error function should be of type `ErrorHandler`. 


```go
k, err := NewAsyncTask(
    WithQueueSizeOption(10), 
    WithWorkerSizeOption(1), WithTimeoutOption(5*time.Hour),
    )
```
### Types and functions
- `NewAsyncTask()` returns a type `keeper`.
        - The `keeper` type has methods like `AddTask`, `UnProcessedTaskSize`
- `dispatcher` is an internal type that controls workers.

### Run locally
Download and install `go`
Clone this repo and make changes.
For testing run `go test -v`