package goasync

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type mockTask struct {
	Data int
}

func (m *mockTask) Process(ctx context.Context) error {
	return nil
}

func TestNewAsyncTask(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputs []Option
		wants  []interface{}
	}{
		"empty": {inputs: []Option{},
			wants: []interface{}{1000, 5, time.Duration(60 * time.Second)}},
		"check": {inputs: []Option{WithQueueSizeOption(10),
			WithWorkerSizeOption(20),
			WithTimeoutOption(10 * time.Second)},
			wants: []interface{}{10, 20, time.Duration(10 * time.Second)}},
	}

	for _, test := range tests {
		k, err := NewAsyncTask(test.inputs...)
		assert.NoError(err)
		assert.True(cmp.Equal(k.queueSize, test.wants[0]))
		assert.True(cmp.Equal(k.workerSize, test.wants[1]))
		assert.True(cmp.Equal(k.timeout, test.wants[2]))
	}
}

func TestKeeper_AddTask(t *testing.T) {
	assert := assert.New(t)

	k, err := NewAsyncTask(WithQueueSizeOption(10), WithWorkerSizeOption(1), WithTimeoutOption(5*time.Second))
	assert.NoError(err)

	ctx := context.Background()

	// fail
	err = k.AddTask(ctx, nil)
	assert.Error(err)

	// success
	for i := 0; i < 100; i++ {
		err := k.AddTask(ctx, &mockTask{Data: i})
		assert.NoError(err)
	}
	time.Sleep(1 * time.Second)

	// stop
	k.dispatcher.stop()

	for i := 0; i < 10; i++ {
		err := k.AddTask(ctx, &mockTask{Data: i})
		assert.NoError(err)
	}
	assert.Equal(10, len(k.dispatcher.taskQueue))

	timectx, _ := context.WithTimeout(ctx, time.Second*2)
	err = k.AddTask(timectx, &mockTask{Data: 10})
	assert.Error(err)

	// start
	k.dispatcher.start()
	time.Sleep(1 * time.Second)
	assert.Equal(0, len(k.dispatcher.taskQueue))
}

func TestKeeper_UnProcessedTaskSize(t *testing.T) {
	assert := assert.New(t)

	empty, err := NewAsyncTask(WithQueueSizeOption(10), WithWorkerSizeOption(1), WithTimeoutOption(5*time.Second))
	assert.NoError(err)

	exist, err := NewAsyncTask(WithQueueSizeOption(200), WithWorkerSizeOption(1), WithTimeoutOption(5*time.Second))
	assert.NoError(err)
	for i := 0; i < 10; i++ {
		exist.AddTask(context.Background(), &mockTask{})
	}

	tests := map[string]struct {
		dispatcher *keeper
		count      int
	}{
		"empty": {empty, 0},
		"exist": {exist, 10},
	}

	for _, t := range tests {
		assert.Equal(t.count, t.dispatcher.UnProcessedTaskSize())
	}
}
