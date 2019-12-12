package daemon

import (
	"errors"
	"time"
)

var ErrAsyncTimeout = errors.New("Timeout")

func Async(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

func AsyncTimeout(timeout time.Duration, f func() error) chan error {
	ch := make(chan error)
	go func() {
		select {
		case err := <-Async(f):
			ch <- err
		case <-time.After(timeout):
			ch <- ErrAsyncTimeout
		}
	}()
	return ch
}
