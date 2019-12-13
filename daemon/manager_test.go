package daemon

import (
	"github.com/ihaiker/gokit/logs"
	"testing"
	"time"
)

func TestDaemonManager_Start(t *testing.T) {
	dm := NewDaemonManager("../conf/programs")
	dm.SetStatusListener(func(process *Process, oldStatus, newStatus FSMState) {
		logs.Info("program:", process.Program.Name, ", from:", oldStatus, ", to:", newStatus)
	})
	if err := dm.Start(); err != nil {
		t.Fatal(err)
	}

	if err := dm.StartProgram("nginx", nil); err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 30)

	if err := dm.StopProgram("nginx"); err != nil {
		t.Fatal(err)
	}

	dm.Stop()
}
