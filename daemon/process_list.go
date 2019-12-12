package daemon

type ProcessList []*Process

func (pl ProcessList) Len() int {
	return len(pl)
}

func (pl ProcessList) Less(i, j int) bool {
	return pl[j].Program.Id > pl[i].Program.Id
}

func (pl ProcessList) Swap(i, j int) {
	tmp := (pl)[i]
	(pl)[i] = (pl)[j]
	(pl)[j] = tmp
}

func (pl ProcessList) Names() []string {
	names := make([]string, len(pl))
	for i := 0; i < len(pl); i++ {
		names[i] = pl[i].Program.Name
	}
	return names
}

func (pl ProcessList) GetName(name string) (*Process, error) {
	for _, p := range pl {
		if p.Program.Name == name {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (pl ProcessList) GetIdx(idx int) (*Process, error) {
	if len(pl) > idx {
		return pl[idx], nil
	} else {
		return nil, ErrNotFound
	}
}

func (pl *ProcessList) Remove(name string) error {
	removeIdx := -1
	for idx, p := range *pl {
		if p.Program.Name == name {
			removeIdx = idx
			break
		}
	}
	if removeIdx == -1 {
		return ErrNotFound
	}
	*pl = append((*pl)[:removeIdx], (*pl)[removeIdx+1:]...)
	return nil
}
