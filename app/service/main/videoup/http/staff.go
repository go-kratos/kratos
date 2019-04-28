package http

import (
	"encoding/json"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

func addApply(c *bm.Context) {
	//1.同意 拒绝 忽略
	//2.申请解除
	v := new(archive.ApplyParam)
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("addApply data(%v)", v)
	c.JSON(vdpSvc.DoApply(c, v, "申请单"))
}

//批量修改
func batchApplys(c *bm.Context) {
	var (
		req = c.Request
		bs  []byte
		err error
		aps archive.StaffBatchParam
	)
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	if err = json.Unmarshal(bs, &aps); err != nil {
		log.Error("http batchApplys() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aps.AID == 0 {
		log.Error("http batchApplys() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	//允许为空 就是删除
	if ok := vdpSvc.CheckStaff(aps.Staffs); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vdpSvc.HandleArchiveApplys(c, aps.AID, aps.Staffs, "admin_edit", true); err != nil {
		log.Error("vdaSvc.batchApplys() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
func viewApply(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("viewApply data(%v)", v)
	c.JSON(vdpSvc.Apply(c, v.ID))
}

func checkMid(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("checkMid data(%v)", v)
	c.JSON(vdpSvc.MidCount(c, v.ID))
}

func applys(c *bm.Context) {
	v := new(struct {
		IDS []int64 `form:"ids,split" validate:"required" `
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.IDS) > 200 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("applys data(%v)", v)
	c.JSON(vdpSvc.Applys(c, v.IDS))
}

func filterApplys(c *bm.Context) {
	v := new(struct {
		ADS []int64 `form:"aids,split" validate:"required" `
		MID int64   `form:"mid" validate:"required" `
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if len(v.ADS) > 200 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("filterApplys data(%v)", v)
	c.JSON(vdpSvc.FilterApplys(c, v.ADS, v.MID))
}

func archiveApplys(c *bm.Context) {
	v := new(struct {
		AID int64 `form:"aid" validate:"required" `
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("archiveApplys data(%v)", v)
	c.JSON(vdpSvc.ApplysByAID(c, v.AID))
}

func staffs(c *bm.Context) {
	v := new(struct {
		AID int64 `form:"aid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	log.Info("staffs data(%v)", v)
	c.JSON(vdpSvc.Staffs(c, v.AID))
}
