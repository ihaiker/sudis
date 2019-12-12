package daemon

import (
	"github.com/ihaiker/gokit/logs"
	"os"
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
	process := NewProcess(program)
	process.statusListener = func(pro *Process, oldState, newState FSMState) {
		logger.Debugf("from %s to %s", oldState, newState)
	}
	process.stdout = os.Stdout
	process.stderr = os.Stderr

	if err := process.startCommand(func(process *Process) {}); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 60)

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
	process := NewProcess(program)
	process.stdout = os.Stdout
	process.stderr = os.Stderr

	if err := process.startCommand(nil); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Hour)

	if err := process.stopCommand(); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second)
}
