package http

import (
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/sudis/libs/config"
	. "github.com/ihaiker/sudis/libs/errors"
	"github.com/ihaiker/sudis/nodes/cluster"
	"github.com/ihaiker/sudis/nodes/dao"
	"github.com/kataras/iris/v12"
)

type NodeController struct {
	clusterManger *cluster.DaemonManager
}

func (self *NodeController) queryNodeList() []*dao.Node {
	return self.clusterManger.QueryNode()
}

func (self *NodeController) addOrModifyNodeToken(json *dao.JSON) int {
	key := json.String("key")
	errors.True(key != "", ErrNodeIsEmpty)

	token := json.String("token")
	errors.True(token != "", ErrNodeIsEmpty)

	errors.Assert(self.clusterManger.ModifyNodeToken(key, token))
	return iris.StatusNoContent
}

func (self *NodeController) modifyNodeTag(json *dao.JSON) int {

	key := json.String("key")
	errors.True(key != "", ErrNodeIsEmpty)

	if key == config.Config.Key {
		return iris.StatusNoContent
	}
	tag := json.String("tag")

	errors.Assert(self.clusterManger.ModifyNodeTag(key, tag))

	return iris.StatusNoContent
}

func (self *NodeController) removeNode(key string) int {
	errors.Assert(self.clusterManger.RemoveJoin(key), "删除节点")
	return iris.StatusNoContent
}
