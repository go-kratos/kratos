package http

import (
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

//查询某个资源所命中的所有流量套餐
func hitFlows(c *bm.Context) {
	pm := new(struct {
		OID int64 `form:"aid" validate:"required"`
	})
	if err := c.Bind(pm); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(vdaSvc.HitFlowGroups(c, pm.OID, []int8{archive.PoolArcForbid}))
}
