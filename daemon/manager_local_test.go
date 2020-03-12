package daemon_test

import (
	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/daemon"
	"testing"
	"time"
)

func TestDaemonManager_Start(t *testing.T) {
	dm := daemon.NewDaemonManager("/tmp")

	dm.SetStatusListener(func(process *daemon.Process, oldStatus, newStatus daemon.FSMState) {
		logs.Info("program:", process.Program.Name, ", from:", oldStatus, ", to:", newStatus)
	})
	if err := dm.Start(); err != nil {
		t.Fatal(err)
	}

	if err := dm.StartProgram("ping", nil); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 30)

	if err := dm.StopProgram("ping"); err != nil {
		t.Fatal(err)
	}

	dm.Stop()
}
