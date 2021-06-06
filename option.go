package goasync

import "time"

type Option interface {
	apply(k *keeper)
}

type OptionFunc func(k *keeper)

func (f OptionFunc) apply(k *keeper) { f(k) }

// WithErrorHandlerOption returns a function which sets handler for the error.
func WithErrorHandlerOption(f ErrorHandler) OptionFunc {
	return func(k *keeper) {
		k.errorHandler = f
	}
}

// WithQueueSizeOption returns a function which sets queue size.
func WithQueueSizeOption(size int) OptionFunc {
	return func(k *keeper) {
		k.queueSize = size
	}
}

// WithWorkerSizeOption returns a function which sets worker size.
func WithWorkerSizeOption(size int) OptionFunc {
	return func(k *keeper) {
		k.workerSize = size
	}
}

// WithTimeoutOption returns a function which set timeout.
func WithTimeoutOption(timeout time.Duration) OptionFunc {
	return func(k *keeper) {
		k.timeout = timeout
	}
}
