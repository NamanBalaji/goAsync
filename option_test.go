package goasync

import (
	"log"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorHandlerOption(t *testing.T) {
	assert := assert.New(t)

	f := func(err error) {
		log.Println(err)
	}
	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{f}, wants: []interface{}{f}},
	}

	for _, t := range tests {
		k := &keeper{}
		opt := WithErrorHandlerOption(t.inputs[0].(func(error)))
		opt.apply(k)
		assert.NotNil(k.errorHandler)
	}
}

func TestWithWorkerSizeOption(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{10}, wants: []interface{}{10}},
	}

	for _, t := range tests {
		k := &keeper{}
		opt := WithWorkerSizeOption(t.inputs[0].(int))
		opt.apply(k)
		assert.True(cmp.Equal(t.wants[0], k.workerSize))
	}
}

func TestWithQueueSizeOption(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{10}, wants: []interface{}{10}},
	}

	for _, t := range tests {
		k := &keeper{}
		opt := WithQueueSizeOption(t.inputs[0].(int))
		opt.apply(k)
		assert.True(cmp.Equal(t.wants[0], k.queueSize))
	}
}

func TestWithTimeoutOption(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{time.Duration(10 * time.Second)}, wants: []interface{}{time.Duration(10 * time.Second)}},
	}

	for _, t := range tests {
		k := &keeper{}
		opt := WithTimeoutOption(t.inputs[0].(time.Duration))
		opt.apply(k)
		assert.True(cmp.Equal(t.wants[0], k.timeout))
	}
}
