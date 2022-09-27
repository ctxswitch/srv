package srv

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	PollTime time.Duration = 100 * time.Microsecond
)

type Limiter struct {
	maxConn int32
	numConn int32
	sync.RWMutex
}

func NewLimiter(size int32) *Limiter {
	return &Limiter{
		maxConn: size,
		numConn: 0,
	}
}

func (l *Limiter) Add() {
	l.Lock()
	defer l.Unlock()
	// Block until we have room in the Limiter
	for atomic.LoadInt32(&l.numConn) >= l.maxConn {
		time.Sleep(PollTime)
	}
	atomic.AddInt32(&l.numConn, 1)
}

func (l *Limiter) Remove() {
	atomic.AddInt32(&l.numConn, -1)
}
