package http

import (
	"encoding/json"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func staffs(c *bm.Context) {
	v := new(struct {
		AID int64 `form:"aid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("staffs data(%v)", v)
	c.JSON(vdaSvc.Staffs(c, v.AID))
}

//batchStaff .
func batchStaff(c *bm.Context) {
	v := new(struct {
		AID    int64  `form:"aid" validate:"required"`
		DelAll int8   `form:"del_all"`
		Staffs string `form:"staffs"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("batchStaff data(%v)", v)
	var err error
	var aps = &archive.StaffBatchParam{}
	if err = json.Unmarshal([]byte(v.Staffs), &aps.Staffs); err != nil {
		log.Error("http batchStaff Staffs json.Unmarshal(%s) error(%v)", v.Staffs, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.AID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aps.Staffs) == 0 && v.DelAll != 1 {
		log.Info("batchStaff del_all data is wrong(%v) staffs(%+v)", v, aps)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aps.Staffs == nil {
		aps.Staffs = []*archive.StaffParam{}
	}
	if ok := vdaSvc.CheckStaff(aps.Staffs); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aps.AID = v.AID
	c.JSON(nil, vdaSvc.StaffApplyBatchSubmit(c, aps))
}
