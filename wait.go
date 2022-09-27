package srv

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	PollTime time.Duration = 100 * time.Microsecond
)

type Pool struct {
	maxConn int32
	numConn int32
	sync.RWMutex
}

func NewPool(size int32) *Pool {
	return &Pool{
		maxConn: size,
		numConn: 0,
	}
}

func (p *Pool) Add() {
	p.Lock()
	defer p.Unlock()
	// Block until we have room in the pool
	for atomic.LoadInt32(&p.numConn) >= p.maxConn {
		time.Sleep(PollTime)
	}
	atomic.AddInt32(&p.numConn, 1)
}

func (p *Pool) Remove() {
	atomic.AddInt32(&p.numConn, -1)
}
