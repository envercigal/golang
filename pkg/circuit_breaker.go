package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("Error opening circuit breaker")
var ErrHalfOpen = errors.New("Error half opening circuit breaker")

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type Breaker struct {
	mu           sync.Mutex
	state        State
	failureCount int
	lastFailure  time.Time
	maxFailures  int
	resetTimeout time.Duration
}

func New(maxFailures int, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		state:        Closed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
	}
}

func (b *Breaker) Execute(fn func() error) error {
	if !b.allowRequest() {
		return ErrOpen
	}

	err := fn()

	if err != nil {
		b.recordFailure()
		return err
	}

	b.recordSuccess()
	return nil
}

func (b *Breaker) allowRequest() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Open:
		if time.Since(b.lastFailure) > b.resetTimeout {
			b.state = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	default:
		return true
	}
}

func (b *Breaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount++
	b.lastFailure = time.Now()
	if b.failureCount >= b.maxFailures {
		b.state = Open
	}
}

func (b *Breaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount = 0
	b.state = Closed
}
