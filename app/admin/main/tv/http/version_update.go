package http

import (
	"fmt"
	"net/url"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_isVerDeleted = 2
)

func verUpdateList(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*model.VersionUpdate
		limit []*model.VersionUpdateLimit
		rets  []*model.VersionUpdateDetail
		count int64
		vid   = atoi(req.Get("vid"))
		page  = atoi(req.Get("page"))
		size  = 20
	)
	if page == 0 {
		page = 1
	}
	db := tvSrv.DBShow.Where("vid=?", vid).Where("state!=?", _isVerDeleted)
	db.Model(&model.VersionUpdate{}).Count(&count)
	if err = db.Model(&model.VersionUpdate{}).Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	for _, val := range items {
		if err = tvSrv.DBShow.Model(&model.VersionUpdateLimit{}).Where("up_id=?", val.ID).Find(&limit).Error; err != nil {
			log.Error("%v\n", err)
			c.JSON(nil, err)
			return
		}
		rets = append(rets, &model.VersionUpdateDetail{VersionUpdate: val, VerLimit: limit})
	}

	version := model.Version{}
	if err = tvSrv.DBShow.Where("id=?", vid).Where("state!=?", _isDeleted).First(&version).Error; err != nil {
		log.Error("version info not exists-(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	lists := map[string]interface{}{
		"verUpdate": rets,
		"version":   version,
	}
	pager := &model.VersionUpdatePager{
		TotalCount: count,
		Pn:         page,
		Ps:         size,
		Items:      lists,
	}
	c.JSON(pager, nil)
}

func saveVerUpdate(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	exist := model.VersionUpdate{}
	if err = tvSrv.DBShow.Where("id=?", id).Where("state!=?", _isVerDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	alert, simple := validateVerUpdatePostData(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	if err = tvSrv.DBShow.Model(&model.VersionUpdate{}).Where("id = ?", id).Update(simple).Error; err != nil {
		log.Error("tvSrv.saveVerUpdate error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.DBShow.Model(&model.VersionUpdate{}).Where("id = ?", id).Update(map[string]int8{"is_push": simple.IsPush, "is_force": simple.IsForce}).Error; err != nil {
		log.Error("tvSrv.saveVerUpdate error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.DBShow.Model(&model.VersionUpdateLimit{}).Where("up_id=?", id).Delete(&model.VersionUpdateLimit{}).Error; err != nil {
		log.Error("tvSrv.DeleteVerUpdateLimit error(%v)\n", err)
		return
	}
	log.Info("saveVerUpdate exist.ID = %d", exist.ID)
	if exist.ID > 0 {
		addVerUpdateLimit(req, exist.ID)
	}
	c.JSON(nil, nil)
}

func addVerUpdate(c *bm.Context) {
	var (
		req = c.Request.PostForm
		err error
	)
	alert, simple := validateVerUpdatePostData(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	db := tvSrv.DBShow.Create(simple)
	if err = db.Error; err != nil {
		log.Error("tvSrv.addVerUpdate error(%v)", err)
		c.JSON(nil, err)
		return
	}
	insertID := (db.Value.(*model.VersionUpdate)).ID
	if insertID > 0 {
		addVerUpdateLimit(req, insertID)
	}
	c.JSON(nil, nil)
}

func addVerUpdateLimit(req url.Values, upid int64) (err error) {
	var (
		condi = req["condi[]"]
		value = req["value[]"]
	)
	if len(condi) > 0 {
		for key, val := range condi {
			li := &model.VersionUpdateLimit{UPID: int32(upid)}
			if key < len(condi) {
				li.Condi = val
			}
			if key < len(value) {
				li.Value = atoi(value[key])
			}
			if err = tvSrv.DBShow.Create(li).Error; err != nil {
				log.Error("tvSrv.addVerUpdateLimit error(%v)", err)
				break
			}
		}
	}
	return
}

func verUpdateEnable(c *bm.Context) {
	var (
		req = c.Request.PostForm

		id    = parseInt(req.Get("id"))
		state = atoi(req.Get("state"))
		err   error
	)
	exist := model.VersionUpdate{}
	if err = tvSrv.DBShow.Where("id=?", id).Where("state!=?", _isVerDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err := tvSrv.DBShow.Model(&model.VersionUpdate{}).Where("id = ?", id).Update(map[string]int{"state": state}).Error; err != nil {
		log.Error("tvSrv.verUpdateEnable error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func fullPackageImport(c *bm.Context) {
	var (
		req = c.Request.Form

		vid   = atoi(req.Get("vid"))
		build = atoi(req.Get("build"))
		err   error
	)
	result, err := tvSrv.FullImport(c, build)
	if err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	for _, val := range result {
		if !createApk2Version(val, vid) {
			renderErrMsg(c, ecode.RequestErr.Code(), fmt.Sprintf("fullPackageImport fail(%v)", val))
			return
		}
	}
	c.JSON(nil, nil)
}

func createApk2Version(val *model.APKInfo, vid int) (b bool) {
	var (
		err    error
		simple *model.VersionUpdate
		limit  *model.VersionUpdateLimit
	)
	b = false
	simple = new(model.VersionUpdate)
	simple.VID = vid
	simple.PolicyName = "指定版本导入更新"
	simple.IsForce = 0
	simple.IsPush = 0
	simple.Channel = "bili"
	simple.URL = val.CDNAddr
	simple.Size = val.Size
	simple.Md5 = val.SignMd5
	simple.Sdkint = 0
	simple.Model = ""
	simple.Policy = 0
	simple.Coverage = 100
	db := tvSrv.DBShow.Create(simple)
	if err = db.Error; err != nil {
		log.Error("tvSrv.createApk2Version error(%v)", err)
		return
	}
	insertID := (db.Value.(*model.VersionUpdate)).ID
	limit = &model.VersionUpdateLimit{UPID: int32(insertID), Condi: "", Value: 0}
	if err = tvSrv.DBShow.Create(limit).Error; err != nil {
		log.Error("tvSrv.createAPK2Version error(%v)", err)
		return
	}
	return true
}

func validateVerUpdatePostData(c *bm.Context) (alert string, simple *model.VersionUpdate) {
	var (
		req        = c.Request.PostForm
		vid        = atoi(req.Get("vid"))
		isForce    = atoi(req.Get("is_force"))
		isPush     = atoi(req.Get("is_push"))
		channel    = req.Get("channel")
		url        = req.Get("url")
		size       = atoi(req.Get("size"))
		md5        = req.Get("md5")
		sdkint     = atoi(req.Get("sdkint"))
		mod        = req.Get("model")
		policy     = atoi(req.Get("policy"))
		coverage   = atoi(req.Get("coverage"))
		policyName = req.Get("policy_name")
	)
	alert = string("")
	simple = new(model.VersionUpdate)
	simple.VID = vid
	simple.PolicyName = policyName
	simple.IsForce = int8(isForce)
	simple.IsPush = int8(isPush)
	simple.Channel = channel
	simple.URL = url
	simple.Size = size
	simple.Md5 = md5
	simple.Sdkint = sdkint
	simple.Model = mod
	simple.Policy = int8(policy)
	simple.Coverage = int32(coverage)
	if simple.Channel == "" {
		alert = "渠道不能为空"
		return
	}
	if simple.URL == "" {
		alert = "安装包地址不能为空"
		return
	}
	if simple.Size == 0 {
		alert = "文件大小不能为0"
		return
	}
	if simple.Md5 == "" {
		alert = "文件hash值不能为空"
		return
	}
	return
}
