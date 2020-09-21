package http

import (
	"github.com/ihaiker/gokit/errors"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
)

type TagsController struct {
}

func (self *TagsController) queryTag(ctx iris.Context) []*dao.Tag {
	tags, err := dao.TagDao.List()
	errors.Assert(err)
	return tags
}

func (self *TagsController) addOrModify(ctx iris.Context) int {
	json := &dao.JSON{}
	errors.Assert(ctx.ReadJSON(json))

	name := json.String("name")
	errors.True(name != "", ErrNameEmpty)

	class := json.String("class")
	errors.True(class != "", ErrClassEmpty)

	err := dao.TagDao.AddOrUpdate(name, class)
	errors.Assert(err)
	return iris.StatusNoContent
}

func (self *TagsController) removeTag(ctx iris.Context) int {
	name := ctx.Params().GetString("name")
	errors.Assert(dao.TagDao.Remove(name))
	return iris.StatusNoContent
}
