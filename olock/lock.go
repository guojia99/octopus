/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午2:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package olock

type lockType int

const (
	OptimisticL lockType = iota // 乐观锁 sync.RWMutex
	ReentrantL                  // 悲观锁 sync.Mutex
)
