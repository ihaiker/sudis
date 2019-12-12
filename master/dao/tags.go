package dao

type Tag struct {
	Name  string `json:"name" yaml:"name" toml:"name" xorm:"varchar(64) notnull pk 'name'"`
	Class string `json:"class" yaml:"class" toml:"class" xorm:"class"`
}

type tagDao struct {
}

func (self *tagDao) List() (tags []*Tag, err error) {
	tags = make([]*Tag, 0)
	err = engine.Find(&tags)
	return
}

func (self *tagDao) AddOrUpdate(name, class string) error {
	tag := new(Tag)
	if has, err := engine.Where("name = ?", name).Get(tag); err != nil {
		return err
	} else if has {
		tag.Class = class
		_, err = engine.Update(tag, &Tag{Name: name})
		return err
	} else {
		tag.Name = name
		tag.Class = class
		_, err = engine.InsertOne(tag)
		return err
	}
}

func (self *tagDao) Remove(name string) error {
	_, err := engine.Delete(&Tag{Name: name})
	return err
}

var TagDao = new(tagDao)
