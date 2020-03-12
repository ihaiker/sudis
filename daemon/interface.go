package daemon

type Manager interface {
	Start() error
	Stop() error

	GetProcess(name string) (*Process, error)

	AddProgram(program *Program) error
	RemoveProgram(name string, skip bool) error
	ModifyProgram(program *Program) error

	ModifyTag(name string, add bool, tag string) error

	ListProgramNames() ([]string, error)
	ListProcess() ([]*Process, error)

	SetStatusListener(lis FSMStatusListener)

	StartProgram(name string, determinedResult chan *Process) error
	StopProgram(name string) error

	MustGetProcess(name string) *Process
	MustAddProgram(program *Program)
	MustRemoveProgram(name string, skip bool)
	MustModifyProgram(program *Program)
	MustModifyTag(name string, add bool, tag string)
	MustListProgramNames() []string
	MustListProcess() []*Process

	MustStartProgram(name string, determinedResult chan *Process)
	MustStopProgram(name string)

	SubscribeLogger(uid string, name string, tail TailLogger, firstLine int) error
	UnSubscribeLogger(uid, name string) error
}
