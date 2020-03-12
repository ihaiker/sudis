package daemon

import (
	"testing"
)

func TestProcInfo(t *testing.T) {
	cup, rss, err := GetProcessInfo(362)
	t.Log(cup)
	t.Log(rss)
	t.Log(err)
}
