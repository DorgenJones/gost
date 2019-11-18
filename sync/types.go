/**
*
* Copyright (c) 2018 MaoYan
* All rights reserved
* Author: dujiang02
* Date: 2019-11-17
 */
package gxsync

import "fmt"

type Task func()


/////////////////////////////////////////
// Task Pool Options
/////////////////////////////////////////

type TaskPoolOptions struct {
	tQLen      int // Task queue length
	tQNumber   int // Task queue number
	tQPoolSize int // Task pool size
}

func (o *TaskPoolOptions) validate() {
	if o.tQPoolSize < 1 {
		panic(fmt.Sprintf("illegal pool size %d", o.tQPoolSize))
	}

	if o.tQLen < 1 {
		o.tQLen = defaultTaskQLen
	}

	if o.tQNumber < 1 {
		o.tQNumber = defaultTaskQNumber
	}

	if o.tQNumber > o.tQPoolSize {
		o.tQNumber = o.tQPoolSize
	}
}

type TaskPoolOption func(*TaskPoolOptions)

// @size is the Task queue pool size
func WithTaskPoolTaskPoolSize(size int) TaskPoolOption {
	return func(o *TaskPoolOptions) {
		o.tQPoolSize = size
	}
}

// @length is the Task queue length
func WithTaskPoolTaskQueueLength(length int) TaskPoolOption {
	return func(o *TaskPoolOptions) {
		o.tQLen = length
	}
}

// @number is the Task queue number
func WithTaskPoolTaskQueueNumber(number int) TaskPoolOption {
	return func(o *TaskPoolOptions) {
		o.tQNumber = number
	}
}

type IShardingTaskPool interface {
	AddTask(t Task)
	AddShardTask(index int, t Task)
	IsClosed()  bool
	Close()
}
