package daemon

type FSMEvent string

const (
	StartProgram      = FSMEvent("start-program")       //程序启动
	StopProgram       = FSMEvent("stop-program")        //程序已停止
	ProgramNotSurvive = FSMEvent("program-not-survive") //程序不存活
	ProgramSurvive    = FSMEvent("program-survive")     //程序恢复存活
)

type FSMState string

const (
	Ready     = FSMState("ready")
	Starting  = FSMState("starting")
	Running   = FSMState("running")
	Stopped   = FSMState("stopped")
	Fatal     = FSMState("fatal")
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

//如果oldStatus为空，则是新程序添加
//如果newStatus为空，则为删除管理
type FSMStatusListener func(process *Process, oldStatus, newStatus FSMState)

func StdoutStatusListener(process *Process, oldStatus, newStatus FSMState) {
	logger.Debug("program(%s)  from %s to %s", process.Program.Name, oldStatus, newStatus)
}
