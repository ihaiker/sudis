package daemon

import (
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	err := Async(func() error {
		time.Sleep(time.Second)
		return nil
	})

	select {
	case e := <-err:
		t.Log(e)
	}
}
