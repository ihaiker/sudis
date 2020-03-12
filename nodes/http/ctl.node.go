package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/libs/config"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris"
)

type NodeController struct {
	clusterManger *cluster.DaemonManager
}

func (self *NodeController) queryNodeList() []*dao.Node {
	return self.clusterManger.QueryNode()
}

func (self *NodeController) modifyNodeTag(json *dao.JSON) int {

	key := json.String("key")
	errors.True(key != "", ErrNodeIsEmpty)

	if key == config.Config.Key {
		return iris.StatusNoContent
	}
	tag := json.String("tag")

	self.clusterManger.ModifyNodeTag(key, tag)

	return iris.StatusNoContent
}
