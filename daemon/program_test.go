package daemon

import (
	"encoding/json"
	"github.com/ihaiker/gokit/files"
	"testing"
	"time"
)

func TestNewProgram(t *testing.T) {
	p := NewProgram()
	p.Name = "pingbaidu"
	p.Start = &Command{
		Command: "ping",
		Args: []string{
			"baidu.com",
		},
	}
	p.AddTime = time.Now()
	p.UpdateTime = time.Now()

	if bs, err := json.MarshalIndent(p, "\t", "\n"); err != nil {
		t.Fatal(err)
	} else {
		w, err := files.New("../conf/programs/" + p.Name + ".json").GetWriter(false)
		if err != nil {
			t.Fatal(err)
		}
		defer w.Close()
		if _, err = w.Write(bs); err != nil {
			t.Fatal(err)
		}
	}
}
