package conf

import (
	"github.com/ihaiker/gokit/files"
	"github.com/pelletier/go-toml"
	"testing"
)

func TestConfig(t *testing.T) {
	bs, err := toml.Marshal(*Config)
	if err != nil {
		t.Fatal(err)
	}
	f := files.New("./sudis.toml")
	t.Log(f.GetPath())
	w, err := f.GetWriter(false)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	i, err := w.Write(bs)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(i)
}
