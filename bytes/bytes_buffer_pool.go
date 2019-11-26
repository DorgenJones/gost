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

package gxbytes

import (
	"sync"
	"sync/atomic"
)

const shardSize = 4

var index uint64 = 1

var (
	shardingBufferPool [shardSize]*ObjectPool
)

func init() {
	for i := 0; i < shardSize; i++ {
		shardingBufferPool[i] = NewObjectPool(func() PoolObject {
			return new(ByteBuffer)
		})
	}
}

func GetByteBuffer(size int) Buffer {
	i := atomic.AddUint64(&index, 1) % shardSize
	return shardingBufferPool[i].Get(size).(Buffer)
}

func PutByteBuffer(buf Buffer) {
	i := atomic.AddUint64(&index, 1) % shardSize
	shardingBufferPool[i].Put(buf)
}

// Pool object
type PoolObject interface {
	Name() string
	Free()
	Init(param interface{})
}

type New func() PoolObject

// Pool is bytes.buffer Pool
type ObjectPool struct {
	New  New
	pool sync.Pool
}

func NewObjectPool(n New) *ObjectPool {
	return &ObjectPool{New: n}
}

// take returns *bytes.buffer from Pool
func (p *ObjectPool) Get(param interface{}) PoolObject {
	v := p.pool.Get()
	if v == nil {
		v = p.New()
	}
	result := v.(PoolObject)
	result.Init(param)
	return result
}

// give returns *byes.buffer to Pool
func (p *ObjectPool) Put(o PoolObject) {
	o.Free()
	p.pool.Put(o)
}