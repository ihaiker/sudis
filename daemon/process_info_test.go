package daemon

import (
	"github.com/ihaiker/gokit/logs"
	"testing"
	"time"
)

func TestProcInfo(t *testing.T) {
	p, err := NewProcessInfo(77681)
	if err != nil {
		t.Fatal(err)
	}

	for {
		if pi, err := p.ProcInfo(); err != nil {
			t.Fatal(err)
		} else {
			logs.Debug("cpu: ", pi.PCpu, ", rss: ", pi.Rss)
		}
		time.Sleep(time.Second)
	}

}
