package http

import (
	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func updateNodes(c *bm.Context) {
	var (
		umnr = &model.UpdateMachineNodeRequest{}
		err  error
	)
	if err = c.BindWith(umnr, binding.JSON); err != nil {
		return
	}
	if err = umnr.VerifyNodes(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.UpdateMachineNode(c, umnr))
}

func queryNodes(c *bm.Context) {
	v := new(struct {
		MachineID int64 `form:"machine_id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QueryMachineNodes(v.MachineID))
}
