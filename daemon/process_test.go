package daemon

import (
	"github.com/ihaiker/gokit/logs"
	"testing"
	"time"
)

func init() {
	logs.SetDebugMode(true)
}

func TestForegroundProcess(t *testing.T) {
	program := NewProgram()
	program.Name = "pingbaidu"
	program.Start = &Command{
		Command: "ping",
		Args: []string{
			"baidu.com",
		},
	}
	program.Logger = "/tmp/sudis.log"
	process, err := NewProcess(program)
	if err != nil {
		t.Fatal(err)
	}
	process.statusListener = func(pro *Process, oldState, newState FSMState) {
		logger.Debugf("from %s to %s", oldState, newState)
	}

	if err := process.startCommand(nil); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 10)

	if err := process.stopCommand(); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second)
}

func TestDaemonProcess(t *testing.T) {
	program := NewProgram()
	program.Name = "nginx"
	program.IgnoreAlreadyStarted = true
	program.Start = &Command{
		Command: "nginx",
		CheckHealth: &CheckHealth{
			CheckAddress: "http://127.0.0.1",
			CheckMode:    HTTP,
			SecretToken:  "",
			CheckTtl:     3,
		},
	}
	program.Stop = &Command{
		Command: "nginx",
		Args:    []string{"-s", "quit"},
	}

	process, err := NewProcess(program)
	if err != nil {
		t.Fatal(err)
	}

	if err := process.startCommand(nil); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Hour)

	if err := process.stopCommand(); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second)
}
