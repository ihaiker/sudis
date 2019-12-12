package eventbus

type Event struct {
	Name  string
	Value interface{}
}

func NewNode(nodeIp, key string) *Event {
	return &Event{
		Name:  "NewNode",
		Value: []string{nodeIp, key},
	}
}

func LostNode(nodeIp, key string) *Event {
	return &Event{
		Name:  "LostNode",
		Value: []string{nodeIp, key},
	}
}

func Shutdown() *Event {
	return &Event{
		Name: "Shutdown",
	}
}

type ProgramStatusEvent struct {
	Ip        string
	Key       string
	Name      string
	OldStatus string
	NewStatus string
}

//node, key, name,oldStatus,newStatus
func ProgramStatus(event *ProgramStatusEvent) *Event {
	return &Event{
		Name:  "ProgramStatus",
		Value: event,
	}
}
