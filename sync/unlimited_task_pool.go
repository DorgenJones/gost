/**
*
* Copyright (c) 2018 MaoYan
* All rights reserved
* Author: dujiang02
* Date: 2019-11-17
 */
package gxsync

import "sync"

type UnlimitedTaskPool struct {
	TaskPoolOptions

	idx    uint32 // round robin index
	qArray []chan task
	wg     sync.WaitGroup

	once sync.Once
	done chan struct{}
}