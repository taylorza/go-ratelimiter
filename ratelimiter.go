package ratelimiter

import (
	"math"
	"sync"
	"time"
)

const defaultRefreshPeriod = 10

// Limiter used to limit the rate at which work is done
type Limiter struct {
	tpp    int
	tokens int
	m      sync.Mutex
	c      *sync.Cond
}

// New creates a rate limiter that can be used to throttle work to target rate per second
func New(rate int) *Limiter {
	l := &Limiter{}

	l.tpp = int(math.Round(float64(rate)/1000) * defaultRefreshPeriod)
	l.tokens = l.tpp
	l.c = sync.NewCond(&l.m)

	go func() {
		t := time.NewTicker(defaultRefreshPeriod * time.Millisecond)
		for range t.C {
			l.m.Lock()
			notify := l.tokens == 0
			l.tokens += l.tpp
			if l.tokens > l.tpp {
				l.tokens = l.tpp
			}
			l.m.Unlock()
			if notify {
				l.c.Broadcast()
			}
		}
	}()

	return l
}

// Throttle blocks if the current rate of work exceeds the limiter
func (l *Limiter) Throttle() {
	l.m.Lock()
	for l.tokens == 0 {
		l.c.Wait()
	}
	l.tokens--
	l.m.Unlock()
}
