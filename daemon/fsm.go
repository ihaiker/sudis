package daemon

type (
	FSMState string

	FSMStatusEvent struct {
		Process    *Process `json:"process"`
		FromStatus FSMState `json:"fromStatus"`
		ToStatus   FSMState `json:"toStatus"`
	}

	FSMStatusListener func(event FSMStatusEvent)
)

const (
	Ready     = FSMState("ready")
	Starting  = FSMState("starting")
	Running   = FSMState("running")
	Fail      = FSMState("fail")
	RetryWait = FSMState("retry")
	Stopping  = FSMState("stopping")
	Stoped    = FSMState("stoped")
)

func (f FSMState) String() string {
	return string(f)
}

func (f FSMState) IsRunning() bool {
	return f == Running || f == Starting
}

func StdoutStatusListener(event FSMStatusEvent) {
	logger.Debugf("program(%s)  from %s to %s", event.Process.Name, event.FromStatus, event.ToStatus)
}
