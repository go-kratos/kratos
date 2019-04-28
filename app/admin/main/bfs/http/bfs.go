package http

import (
	"go-common/app/admin/main/bfs/model"
	bm "go-common/library/net/http/blademaster"
)

func clusters(c *bm.Context) {
	clusters := srv.Clusters(c)
	c.JSON(clusters, nil)
}

func bfsTotal(c *bm.Context) {
	arg := &model.ArgCluster{}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Total(c, arg))
}

func rackMeta(c *bm.Context) {
	arg := &model.ArgCluster{}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Racks(c, arg))
}

func groupMeta(c *bm.Context) {
	arg := &model.ArgCluster{}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Groups(c, arg))
}

func volumeMeta(c *bm.Context) {
	arg := &model.ArgCluster{}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Volumes(c, arg))
}

func addVolume(c *bm.Context) {
	arg := new(model.ArgAddVolume)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.AddVolume(c, arg))
}

func addFreeVolume(c *bm.Context) {
	arg := new(model.ArgAddFreeVolume)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.AddFreeVolume(c, arg))
}

func compact(c *bm.Context) {
	arg := new(model.ArgCompact)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.Compact(c, arg))
}

func setGroupStatus(c *bm.Context) {
	arg := new(model.ArgGroupStatus)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, srv.SetGroupStatus(c, arg))
}
