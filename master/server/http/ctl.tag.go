package http

import (
	"github.com/ihaiker/sudis/master/dao"
	"github.com/kataras/iris"
)

type TagsController struct {
}

func (self *TagsController) queryTag(ctx iris.Context) []*dao.Tag {
	tags, err := dao.TagDao.List()
	AssertErr(err)
	return tags
}

func (self *TagsController) addOrModify(ctx iris.Context) int {
	json := &dao.JSON{}
	AssertErr(ctx.ReadJSON(json))

	name := json.String("name")
	Assert(name != "", ErrNameEmpty)

	class := json.String("class")
	Assert(class != "", ErrClassEmpty)

	err := dao.TagDao.AddOrUpdate(name, class)
	AssertErr(err)
	return iris.StatusNoContent
}

func (self *TagsController) removeTag(ctx iris.Context) int {
	name := ctx.Params().GetString("name")
	AssertErr(dao.TagDao.Remove(name))
	return iris.StatusNoContent
}
