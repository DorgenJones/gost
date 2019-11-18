/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gxsync

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	defaultTaskQNumber = 10
	defaultTaskQLen    = 128
)

/////////////////////////////////////////
// Task Pool
/////////////////////////////////////////

// Task pool: manage Task ts
type TaskPool struct {
	TaskPoolOptions

	idx    uint32 // round robin index
	qArray []chan Task
	wg     sync.WaitGroup

	once sync.Once
	done chan struct{}
}

// build a Task pool
func NewTaskPool(opts ...TaskPoolOption) *TaskPool {
	var tOpts TaskPoolOptions
	for _, opt := range opts {
		opt(&tOpts)
	}

	tOpts.validate()

	p := &TaskPool{
		TaskPoolOptions: tOpts,
		qArray:          make([]chan Task, tOpts.tQNumber),
		done:            make(chan struct{}),
	}

	for i := 0; i < p.tQNumber; i++ {
		p.qArray[i] = make(chan Task, p.tQLen)
	}
	p.start()

	return p
}

// start Task pool
func (p *TaskPool) start() {
	for i := 0; i < p.tQPoolSize; i++ {
		p.wg.Add(1)
		workerID := i
		q := p.qArray[workerID%p.tQNumber]
		go p.run(int(workerID), q)
	}
}

// worker
func (p *TaskPool) run(id int, q chan Task) error {
	defer p.wg.Done()

	var (
		ok bool
		t  Task
	)

	for {
		select {
		case <-p.done:
			if 0 < len(q) {
				return fmt.Errorf("Task worker %d exit now while its Task buffer length %d is greater than 0",
					id, len(q))
			}

			return nil

		case t, ok = <-q:
			if ok {
				t()
			}
		}
	}
}

// add Task
func (p *TaskPool) AddTask(t Task) {
	index := atomic.AddUint32(&p.idx, 1)
	p.AddShardTask(int(index), t)
}

// add sharding Task
func (p *TaskPool) AddShardTask(index int, t Task) {
	id := index % p.tQNumber

	select {
	case <-p.done:
		return
	case p.qArray[id] <- t:
	}
}

// stop all tasks
func (p *TaskPool) stop() {
	select {
	case <-p.done:
		return
	default:
		p.once.Do(func() {
			close(p.done)
		})
	}
}

// check whether the session has been closed.
func (p *TaskPool) IsClosed() bool {
	select {
	case <-p.done:
		return true

	default:
		return false
	}
}

func (p *TaskPool) Close() {
	p.stop()
	p.wg.Wait()
	for i := range p.qArray {
		close(p.qArray[i])
	}
}
