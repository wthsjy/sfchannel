package sfchannel

import (
	"errors"
	"sync"
	"sync/atomic"
)

var SafeChannelErrClosed = errors.New("channel is closed")

type SafeChannel[T any] struct {
	ch     chan T
	closed atomic.Bool
	mux    sync.Mutex
}

func NewChannel[T any](capacity uint) *SafeChannel[T] {
	var sc = &SafeChannel[T]{
		ch: make(chan T, capacity),
	}
	return sc
}

// ProduceOrDiscard  if the buf is full ,then discard value
func (sc *SafeChannel[T]) ProduceOrDiscard(value T) (success bool, err error) {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	if sc.closed.Load() {
		return success, SafeChannelErrClosed
	}

	defer func() {
		if r := recover(); r != nil {
			success = false
			err = SafeChannelErrClosed
		}
	}()

	select {
	case sc.ch <- value:
		success = true
	default:
		success = false
	}
	return success, err
}

// Produce produce
func (sc *SafeChannel[T]) Produce(value T) (err error) {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	if sc.closed.Load() {
		return SafeChannelErrClosed
	}

	defer func() {
		if r := recover(); r != nil {
			err = SafeChannelErrClosed
		}
	}()

	sc.ch <- value
	return err
}

func (sc *SafeChannel[T]) Consume() <-chan T {
	return sc.ch
}

func (sc *SafeChannel[T]) Close() {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	if sc.closed.Load() {
		return
	}
	sc.closed.Store(true)
	close(sc.ch)
}
