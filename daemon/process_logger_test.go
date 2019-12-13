package daemon

import (
	"strconv"
	"sync"
	"testing"
)

func TestProcessLogger(t *testing.T) {

	logger, err := NewLogger("")
	if err != nil {
		t.Fatal(err)
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)

	for i := 0; i < 100000; i++ {
		if _, err := logger.Write([]byte(strconv.Itoa(i) + "----\n")); err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()

}
