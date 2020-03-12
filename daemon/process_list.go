package daemon

import "github.com/ihaiker/sudis/libs/errors"

type ProcessGroup []*Process

func (pl ProcessGroup) Len() int {
	return len(pl)
}

func (pl ProcessGroup) Less(i, j int) bool {
	return pl[j].Program.Id > pl[i].Program.Id
}

func (pl ProcessGroup) Swap(i, j int) {
	tmp := (pl)[i]
	(pl)[i] = (pl)[j]
	(pl)[j] = tmp
}

func (pl ProcessGroup) Names() []string {
	names := make([]string, len(pl))
	for i := 0; i < len(pl); i++ {
		names[i] = pl[i].Program.Name
	}
	return names
}

func (pl ProcessGroup) Get(name string) (*Process, error) {
	for _, p := range pl {
		if p.Program.Name == name {
			return p, nil
		}
	}
	return nil, errors.ErrProgramNotFound
}

func (pl *ProcessGroup) Remove(name string) error {
	removeIdx := -1
	for idx, p := range *pl {
		if p.Program.Name == name {
			removeIdx = idx
			break
		}
	}
	if removeIdx == -1 {
		return errors.ErrProgramNotFound
	}
	*pl = append((*pl)[:removeIdx], (*pl)[removeIdx+1:]...)
	return nil
}
