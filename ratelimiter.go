package ratelimiter

import (
	"runtime"
	"sync"
	"time"
)

const defaultRefreshPeriod = 10

// Limiter used to limit the rate at which work is done
type Limiter struct {
	tpp    float64
	tokens float64
	m      sync.Mutex
	c      *sync.Cond
	t      *time.Ticker
	done   chan bool
}

// New creates a rate limiter that can be used to throttle work to target rate per second
func New(rate uint) *Limiter {
	if rate == 0 {
		panic("rate must be greater than 0")
	}
	l := new(Limiter)

	l.done = make(chan bool)
	l.tpp = float64(rate) / 1000 * defaultRefreshPeriod
	l.tokens = l.tpp
	l.c = sync.NewCond(&l.m)
	l.t = time.NewTicker(defaultRefreshPeriod * time.Millisecond)

	go l.tokenReplenisher()

	runtime.SetFinalizer(l, finalizer)

	return l
}

// Throttle blocks if the current rate of work exceeds the limiter
func (l *Limiter) Throttle() {
	l.m.Lock()
	for l.tokens <= 0 {
		l.c.Wait()
	}
	l.tokens--
	l.m.Unlock()
}

func (l *Limiter) tokenReplenisher() {
	for {
		select {
		case <-l.done:
			return
		case <-l.t.C:
			l.m.Lock()
			notify := l.tokens <= 0
			l.tokens += l.tpp
			if l.tokens > l.tpp {
				l.tokens = l.tpp
			}
			l.m.Unlock()
			if notify {
				l.c.Broadcast()
			}
		}
	}
}

func finalizer(l *Limiter) {
	l.done <- true
	l.t.Stop()
}
