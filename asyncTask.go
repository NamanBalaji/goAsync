package goasync

import (
	"context"
	"fmt"
	"time"
)

type (
	// for handling errors inside if tasks
	ErrorHandler func(error)

	//The actual task
	Task interface {
		Process(ctx context.Context) error
	}

	// external interface
	Keeper interface {
		AddTask(ctx context.Context, task Task) error
		UnprocessedTaskSize() int
	}

	//internal interface
	keeper struct {
		errorHandler ErrorHandler
		queueSize    int
		workerSize   int
		timeout      time.Duration
		dispatcher   *dispatcher
	}
)

// NewAsyncTask creates an object implemented Keeper interface
func NewAsyncTask(opts ...Option) (*keeper, error) {
	k := &keeper{}
	o := []Option{
		WithErrorHandlerOption(func(err error) {
			fmt.Printf("%+v \n", err)
		}),
		WithQueueSizeOption(1000),
		WithWorkerSizeOption(5),
		WithTimeoutOption(time.Duration(60 * time.Second)),
	}
	o = append(o, opts...)
	for _, opt := range o {
		opt.apply(k)
	}

	d, err := k.newDispatcher()
	if err != nil {
		return nil, err
	}
	k.dispatcher = d
	k.dispatcher.start()
	return k, nil
}

// AddTask adds a task in asynchronously
func (k *keeper) AddTask(ctx context.Context, task Task) error {
	if ctx == nil || task == nil {
		return fmt.Errorf("[err] AddTask empty params")
	}

	// check context timeout
	select {
	case k.dispatcher.taskQueue <- task:
	case <-ctx.Done():
		return fmt.Errorf("[err] AddTask timeout")
	}
	return nil
}

// UnProcessedTaskSize returns unprocessed task size.
func (k *keeper) UnProcessedTaskSize() int {
	return len(k.dispatcher.taskQueue)
}

// newDispatcher creates a dispatcher object
func (k *keeper) newDispatcher() (*dispatcher, error) {
	workerPool := make(chan chan Task, k.workerSize)
	var workers []*worker

	for i := 0; i < k.workerSize; i++ {
		worker := &worker{
			id:           i,
			workerPool:   workerPool,
			taskChannel:  make(chan Task),
			quit:         make(chan bool),
			errorHandler: k.errorHandler,
			timeout:      k.timeout,
		}
		workers = append(workers, worker)
	}

	return &dispatcher{
		taskQueue:    make(chan Task, k.queueSize),
		workerPool:   workerPool,
		workers:      workers,
		quit:         make(chan bool),
		errorHandler: k.errorHandler,
	}, nil
}
