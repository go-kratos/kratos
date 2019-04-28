package http

import (
	"go-common/app/admin/main/app/model/aids"
	bm "go-common/library/net/http/blademaster"
)

// saveAids save
func saveAids(c *bm.Context) {
	v := &aids.Param{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, aidsSvc.Save(c, v.Aids))
}
