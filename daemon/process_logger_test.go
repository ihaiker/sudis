package daemon

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestProcessLogger(t *testing.T) {
	logger, _ := NewLogger("")

	logger.Reg("1", func(id, line string) {
		fmt.Println(id, " = ", line)
	}, 10)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		if i == 50 {
			logger.UnReg("1")
		} else if i == 20 {
			logger.Reg("2", func(id, line string) {
				fmt.Println(id, " = ", line)
			}, 10)
		}
		if _, err := logger.Write([]byte(strconv.Itoa(i) + "----\n")); err != nil {
			t.Fatal(err)
		}
	}
}
