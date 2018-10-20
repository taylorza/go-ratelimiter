package ratelimiter

import (
	"sync"
	"time"
)

const defaultRefreshPeriod = 10

// Limiter used to limit the rate at which work is done
type Limiter struct {
	tpp     float64
	tokens  float64
	m       sync.Mutex
	c       *sync.Cond
	t       *time.Ticker
	done    chan bool
	started bool
}

// New creates a rate limiter that can be used to throttle work to target rate per second. The returned ratelimiter is started and ready to throttle.
func New(rate uint) *Limiter {
	if rate == 0 {
		panic("rate must be greater than 0")
	}
	l := &Limiter{}
	l.done = make(chan bool)
	l.c = sync.NewCond(&l.m)
	l.t = time.NewTicker(defaultRefreshPeriod * time.Millisecond)
	l.SetRate(rate)
	l.Start()
	return l
}

// Start a stopped ratelimiter if the rate limiter is already started the operation does nothing.
func (l *Limiter) Start() {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.started {
		l.tokens = l.tpp
		l.started = true
		go l.tokenReplenisher()
	}
}

// Stop the rate limiter and releases the internal resources.
func (l *Limiter) Stop() {
	l.m.Lock()
	defer l.m.Unlock()
	if l.started {
		l.done <- true
		l.t.Stop()
		l.started = false
	}
}

// SetRate updates the rate of the ratelimiter on the fly.
func (l *Limiter) SetRate(rate uint) {
	l.m.Lock()
	defer l.m.Unlock()
	l.tpp = float64(rate) / 1000 * defaultRefreshPeriod
}

// Throttle blocks if the current rate of work exceeds the limiter.
func (l *Limiter) Throttle() {
	l.m.Lock()
	if l.started {
		for l.tokens <= 0 {
			l.c.Wait()
		}
		l.tokens--
	}
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
