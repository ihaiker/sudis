package dao

import (
	"errors"
	"github.com/ihaiker/sudis/daemon"
	"strings"
	"xorm.io/xorm"
)

type Tags []string

func (t *Tags) FromDB(db []byte) error {
	if db != nil && string(db) != "" {
		*t = strings.Split(string(db), ",")
		for i, s := range *t {
			(*t)[i] = s[1 : len(s)-1]
		}
	} else {
		*t = []string{}
	}
	return nil
}
func (t Tags) ToDB() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	nt := make([]string, len(t))
	for i, s := range t {
		nt[i] = "(" + s + ")"
	}
	return []byte(strings.Join(nt, ",")), nil
}

func (t *Tags) Add(tag string) {
	*t = append(*t, tag)
}

func (t *Tags) Remove(tag string) {
	var idx int = -1
	for i, s := range *t {
		if tag == s {
			idx = i
		}
	}
	if idx != -1 {
		*t = append((*t)[:idx], (*t)[idx+1:]...)
	}
}

type Program struct {
	Name        string          `json:"name" yaml:"name" toml:"name" xorm:"name"`
	Description string          `json:"description" yaml:"description" toml:"description"`
	Node        string          `json:"node" yaml:"node" toml:"node" xorm:"node"`
	Tags        Tags            `json:"tags" yaml:"tags" toml:"tags" xorm:"tags"`
	Status      daemon.FSMState `json:"status" yaml:"status" toml:"status" xorm:"status"`
	Time        string          `json:"time" yaml:"time" toml:"time" xorm:"time"`
	Sort        uint64          `json:"sort" yaml:"sort" toml:"sort" xorm:"sort"`
}

type programDao struct {
}

func (self *programDao) Ready() {
	_, err := engine.Cols("status").Update(&Program{Status: "ready"})
	logger.Debug(err)
}

func (self *programDao) List(name, node, tag string, status string, page, limit int) (pageInfo *Page, err error) {
	sn := func() *xorm.Session {
		s := engine.Desc("sort").Asc("time")
		if name != "" {
			s = s.And("(name like ? or description like ?)", "%"+name+"%", "%"+name+"%")
		}
		if node != "" {
			s = s.And("node = ?", node)
		}
		if tag != "" {
			s = s.And(" tags like ?", "%("+tag+")%")
		}
		if status != "" {
			s = s.And("status = ?", status)
		}
		return s
	}

	programs := make([]*Program, 0)
	pageInfo = &Page{
		Total: 0,
		Page:  page,
		Limit: limit,
	}
	if pageInfo.Total, err = sn().Count(new(Program)); err != nil {
		return
	}
	err = sn().Limit(limit, (page-1)*limit).Find(&programs)
	pageInfo.Data = programs
	return
}

func (self *programDao) Get(name, node string) (pro *Program, has bool, err error) {
	pro = new(Program)
	has, err = engine.Where("name = ? and node = ?", name, node).Get(pro)
	return
}

func (self *programDao) Add(program *Program) error {
	if _, has, err := self.Get(program.Node, program.Node); err != nil {
		return err
	} else if has {
		return errors.New("exists")
	} else {
		_, err = engine.InsertOne(program)
		return err
	}
}

func (self *programDao) UpdateStatus(node, name, description string, status daemon.FSMState) error {
	_, err := engine.Cols("status", "description").
		Update(&Program{Status: status, Description: description}, &Program{Node: node, Name: name})
	return err
}

func (self *programDao) ModifyTag(name, node, tag string, add bool) error {
	if pro, has, err := self.Get(name, node); err != nil {
		return err
	} else if !has {
		return ErrNotExist
	} else {
		if add {
			pro.Tags.Add(tag)
		} else {
			pro.Tags.Remove(tag)
		}
		_, err = engine.Update(&Program{Tags: pro.Tags}, &Program{Name: name, Node: node})
		return err
	}
}

func (self *programDao) Lost(key string) {
	_, _ = engine.Cols("status").Update(&Program{Status: daemon.Stoped}, &Program{Node: key})
}

func (self *programDao) Remove(node string, name string) error {
	_, err := engine.Cols("name", "node").Delete(&Program{Node: node, Name: name})
	return err
}

var ProgramDao = new(programDao)
