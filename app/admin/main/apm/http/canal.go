package http

import (
	"context"
	"strings"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/canal"
	"go-common/app/admin/main/apm/model/user"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_upsertMasterInfo = "INSERT INTO master_info (addr,remark,leader,cluster) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE remark=?,leader=?,cluster=?"
)

//canalList get canalinfo list
func canalList(c *bm.Context) {
	var err error
	v := new(canal.ListReq)
	if err = c.Bind(v); err != nil {
		return
	}
	data, _ := apmSvc.ProcessCanalList(c, v)

	c.JSON(data, nil)

}

//canalList get canalinfo add
func canalAdd(c *bm.Context) {

	v := new(canal.Canal)
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	err = apmSvc.ApplyAdd(c, v, username)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(nil, err)
}

//canalList get canalinfo edit
func canalEdit(c *bm.Context) {
	v := new(canal.EditReq)
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	if err = apmSvc.ApplyEdit(c, v, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

//canalDelete 根据addr查询对应id进行软删除
func canalDelete(c *bm.Context) {
	var err error
	v := new(canal.ScanReq)
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	if err = apmSvc.ApplyDelete(c, v, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

//canalScanByAddrFromConfig 根据Addr查询对应有效的配置
func canalScanByAddrFromConfig(c *bm.Context) {
	var confData *canal.Results
	cookie := c.Request.Header.Get("Cookie")
	v := new(canal.ScanReq)
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	if confData, err = apmSvc.GetScanInfo(c, v, username, cookie); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(confData, nil)
}

func canalApplyList(c *bm.Context) {
	v := new(canal.ListReq)
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	data, err := apmSvc.ProcessApplyList(c, v)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

//canalApplyDetailToConfig is
func canalApplyDetailToConfig(c *bm.Context) {
	var (
		username string
		err      error
		v        = new(canal.ConfigReq)
		cookie   = c.Request.Header.Get("Cookie")
	)
	if u, err := c.Request.Cookie("username"); err == nil {
		username = u.Value
	} else {
		username = "third"
	}
	if err = c.Bind(v); err != nil {
		return
	}
	//judge legal params
	f := strings.Contains(v.Addr, ":")
	if !f {
		log.Error("canalApplyAdd addr not standard error(%v)", err)
		c.JSON(nil, ecode.CanalAddrFmtErr)
		return
	}
	if v.User != "" && v.Password != "" {
		if err = apmSvc.CheckMaster(c, v); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if err = apmSvc.ProcessCanalInfo(c, v, username); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.ProcessConfigInfo(c, v, cookie, username); err != nil {
		c.JSON(nil, err)
		return
	}
	go apmSvc.SendWechatMessage(context.Background(), v.Addr, canal.TypeMap[canal.TypeApply], "", username, v.Mark, conf.Conf.Canal.Reviewer)

	c.JSON(nil, err)
}

//canalApplyConfigEdit is
func canalApplyConfigEdit(c *bm.Context) {
	cookie := c.Request.Header.Get("Cookie")
	v := new(canal.ConfigReq)
	if err := c.Bind(v); err != nil {
		return
	}
	ap := &canal.Apply{}
	err := apmSvc.DBCanal.Model(&canal.Apply{}).Select("`operator`").Where("addr=?", v.Addr).Find(ap).Error
	if err != nil {
		log.Error("no apply error", err)
		err = ecode.CanalAddrNotFound
		c.JSON(nil, err)
		return
	}
	username := name(c)
	err = apmSvc.Permit(c, username, user.CanalView)
	err0 := apmSvc.Permit(c, username, user.CanalEdit)

	if err != nil || (username != ap.Operator && err0 != nil) {
		log.Error("permit(%v, %s,%s)", username, err)
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if v.User != "" && v.Password != "" {
		if err = apmSvc.CheckMaster(c, v); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if err0 != nil {
		v = &canal.ConfigReq{
			Addr:          v.Addr,
			MonitorPeriod: v.MonitorPeriod,
			Databases:     v.Databases,
			Project:       v.Project,
			Leader:        v.Leader,
			Mark:          v.Mark,
		}
	}
	//judge legal params
	f := strings.Contains(v.Addr, ":")
	if !f {
		log.Error("canalApplyAdd addr not standard error(%v)", err)
		c.JSON(nil, ecode.CanalAddrFmtErr)
		return
	}
	if err = apmSvc.ProcessCanalInfo(c, v, username); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.ProcessConfigInfo(c, v, cookie, username); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(nil, err)
}

//canalAddrAll get canal all addr info
func canalAddrAll(c *bm.Context) {
	var (
		list  []string
		items []*canal.Canal
	)
	err := apmSvc.DBCanal.Model(&canal.Canal{}).Where("is_delete= 0").Select("`addr`").Scan(&items).Error
	if err != nil {
		log.Error("canalAddrAll get addr error(%v)", err)
		c.JSON(nil, err)
		return
	}

	for _, v := range items {
		list = append(list, v.Addr)
	}

	c.JSON(list, nil)
}

// canal 审核
func canalApplyApprovalProcess(c *bm.Context) {
	var (
		err error
	)
	cookie := c.Request.Header.Get("Cookie")
	res := map[string]interface{}{}
	v := new(struct {
		ID    int  `form:"id" validate:"required"`
		State int8 `form:"state" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	apply := &canal.Apply{}
	if err = apmSvc.DBCanal.Where("id = ?", v.ID).First(apply).Error; err != nil {
		log.Error("canalApplyApprovalProcess id error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if !(apply.State == 1 || apply.State == 2) {
		log.Error("canalApplyApprovalProcess apply.state error(%v)", apply.State)
		res["message"] = "只有申请中和打回才可审核"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if !(v.State == 2 || v.State == 3 || v.State == 4) {
		log.Error("canalApplyApprovalProcess v.state error(%v)", v.State)
		res["message"] = "state值范围2,3,4"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	ups := map[string]interface{}{
		"state": v.State,
	}
	if v.State == 3 && apply.State != 3 {
		if err = apmSvc.UpdateProcessTag(c, apply.ConfID, cookie); err != nil {
			log.Error("apmSvc.UpdateProcessTag error(%v)", apply.ID)
			c.JSON(nil, err)
			return
		}
		if err = apmSvc.DBCanal.Exec(_upsertMasterInfo, apply.Addr, apply.Remark, apply.Leader, apply.Cluster, apply.Remark, apply.Leader, apply.Cluster).Error; err != nil {
			log.Error("canalProcess update  master_info error(%v)", err)
			c.JSONMap(nil, err)
			return
		}
	}
	// 更新apply
	if err = apmSvc.DBCanal.Model(apply).Where("id = ?", apply.ID).Update(ups).Error; err != nil {
		log.Error("canalApplyApprovalProcess update error(%v)", apply.ID)
		res["message"] = "修改状态失败"
		c.JSONMap(res, err)
		return
	}
	go apmSvc.SendWechatMessage(context.Background(), apply.Addr, canal.TypeMap[canal.TypeReview], canal.TypeMap[v.State], username, "", []string{apply.Operator})

	c.JSON(nil, err)
}
