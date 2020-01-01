package http

import (
	"github.com/ihaiker/sudis/master/dao"
	"github.com/ihaiker/sudis/master/server"
	"github.com/kataras/iris"
)

type NodeController struct {
	api server.Api
}

func (self *NodeController) queryNodeList(ctx iris.Context) []*dao.Node {
	nodes, err := dao.NodeDao.List()
	AssertErr(err)
	return nodes
}

func (self *NodeController) modifyNodeTag(ctx iris.Context) int {
	json := &dao.JSON{}
	AssertErr(ctx.ReadJSON(json))

	key := json.String("key")
	Assert(key != "", ErrNodeIsEmpty)

	tag := json.String("tag")

	err := dao.NodeDao.Modify(key, tag)
	AssertErr(err)
	return iris.StatusNoContent
}

//fix 强制同步问题该如何处理呢？
func (self *NodeController) forceReload(ctx iris.Context) int {
	//json := &JSON{}
	//AssertErr(ctx.ReadJSON(json))
	//ip := json.String("ip")
	//eventbus.Send(eventbus.NewNode(ip))
	return iris.StatusNoContent
}
