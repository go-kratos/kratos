package http

import (
	"fmt"
	"go-common/app/admin/main/apm/model/app"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"strings"

	bm "go-common/library/net/http/blademaster"
)

func appList(c *bm.Context) {
	var err error
	v := new(struct {
		AppID string `form:"app_id"`
		Pn    int    `form:"pn" default:"1"`
		Ps    int    `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		aa    []*app.App
		count int
		lk    = "%" + v.AppID + "%"
	)
	if v.AppID != "" {
		err = apmSvc.DB.Where("app_id LIKE ?", lk).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	} else {
		err = apmSvc.DB.Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	}
	if err != nil {
		log.Error("apmSvc.AppList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.AppID != "" {
		err = apmSvc.DB.Where("app_id LIKE ?", lk).Model(&app.App{}).Count(&count).Error
	} else {
		err = apmSvc.DB.Model(&app.App{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.AppList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: aa,
		Total: count,
	}
	c.JSON(data, nil)
}

func appAdd(c *bm.Context) {
	var err error
	username := name(c)
	v := new(struct {
		AppTreeID int64  `form:"app_tree_id" validate:"required"`
		AppID     string `form:"app_id" validate:"required"`
		Limit     int64  `form:"limit"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = apmSvc.AppAdd(c, username, v.AppTreeID, v.AppID, v.Limit); err != nil {
		log.Error("apmSvc.appAdd error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func appEdit(c *bm.Context) {
	var err error
	v := new(struct {
		ID    int `form:"id" validate:"required"`
		Limit int `form:"limit" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	a := &app.App{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(a).Error; err != nil {
		log.Error("apmSvc.appEdit find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	ups := map[string]interface{}{
		"limit": v.Limit,
	}
	if err = apmSvc.DB.Model(&app.App{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
		log.Error("apmSvc.appEdit updates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  ups,
		"Old":     a,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, int64(v.ID), "apmSvc.appEdit", sqlLog)
	c.JSON(nil, err)
}

func appDelete(c *bm.Context) {
	v := new(struct {
		ID int `form:"id" validate:"required"`
	})
	username := name(c)
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	a := &app.App{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(a).Error; err != nil {
		log.Error("apmSvc.appDelete find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Delete(a).Error; err != nil {
		log.Error("apmSvc.appDelete delete(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	apmSvc.SendLog(*c, username, 0, 3, int64(v.ID), "apmSvc.appDelete", a)
	c.JSON(nil, err)
}

func appAuthList(c *bm.Context) {
	v := new(struct {
		AppID     string `form:"app_id"`
		ServiceID string `form:"service_id"`
		Pn        int    `form:"pn" default:"1"`
		Ps        int    `form:"ps" default:"20"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		aa    []*app.Auth
		count int
	)
	if v.AppID != "" && v.ServiceID != "" {
		err = apmSvc.DB.Where("app_id=? and service_id=?", v.AppID, v.ServiceID).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	} else if v.AppID != "" {
		err = apmSvc.DB.Where("app_id=?", v.AppID).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	} else if v.ServiceID != "" {
		err = apmSvc.DB.Where("service_id=?", v.ServiceID).Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	} else {
		err = apmSvc.DB.Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	}
	if err != nil {
		log.Error("apmSvc.appAuthList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.AppID != "" {
		err = apmSvc.DB.Where("app_id=?", v.AppID).Model(&app.Auth{}).Count(&count).Error
	} else {
		err = apmSvc.DB.Model(&app.Auth{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.appAuthList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: aa,
		Total: count,
	}
	c.JSON(data, nil)
}

func appAuthAdd(c *bm.Context) {
	v := new(struct {
		AppTreeID     int64  `form:"app_tree_id" validate:"required"`
		AppID         string `form:"app_id" validate:"required"`
		ServiceTreeID int64  `form:"service_tree_id" validate:"required"`
		ServiceID     string `form:"service_id" validate:"required"`
		RPCMethod     string `form:"rpc_method"`
		HTTPMethod    string `form:"http_method"`
		Quota         int64  `form:"quota"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ServiceTreeID == v.AppTreeID {
		log.Error("apmSvc.appAuthAdd service_tree_id=app_tree_id error(%v)", v.ServiceTreeID)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cnt := 0
	if err = apmSvc.DB.Model(&app.Auth{}).Where("service_tree_id=? and app_tree_id=?", v.ServiceTreeID, v.AppTreeID).Count(&cnt).Error; err != nil {
		log.Error("apmSvc.appAuthAdd count error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cnt > 0 {
		log.Error("apmSvc.appAuthAdd count (%v)", cnt)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	a := &app.Auth{
		AppTreeID:     v.AppTreeID,
		AppID:         v.AppID,
		ServiceTreeID: v.ServiceTreeID,
		ServiceID:     v.ServiceID,
		RPCMethod:     v.RPCMethod,
		HTTPMethod:    v.HTTPMethod,
		Quota:         v.Quota,
	}
	if err = apmSvc.DB.Create(a).Error; err != nil {
		log.Error("apmSvc.appAuthAdd create error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 1, a.ID, "apmSvc.appAuthAdd", a)
	c.JSON(nil, err)
}

func appAuthEdit(c *bm.Context) {
	v := new(struct {
		ID         int    `form:"id" validate:"required"`
		RPCMethod  string `form:"rpc_method"`
		HTTPMethod string `form:"http_method"`
		Quota      int    `form:"quota" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	a := &app.Auth{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(a).Error; err != nil {
		log.Error("apmSvc.appAuthEdit find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	ups := map[string]interface{}{
		"rpc_method":  v.RPCMethod,
		"http_method": v.HTTPMethod,
		"quota":       v.Quota,
	}
	if err = apmSvc.DB.Model(&app.Auth{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
		log.Error("apmSvc.appAuthEdit updates error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  ups,
		"Old":     a,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, int64(v.ID), "apmSvc.appAuthEdit", sqlLog)
	c.JSON(nil, err)
}

func appAuthDelete(c *bm.Context) {
	v := new(struct {
		ID int `form:"id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	a := &app.Auth{}
	if err = apmSvc.DB.Where("id = ?", v.ID).Find(a).Error; err != nil {
		log.Error("apmSvc.appAuthDelete find(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DB.Delete(a).Error; err != nil {
		log.Error("apmSvc.appAuthDelete delete(%d) error(%v)", v.ID, err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 3, int64(v.ID), "apmSvc.appAuthDelete", a)
	c.JSON(nil, err)
}

func appTree(c *bm.Context) {
	v := new(struct {
		Bu   string `form:"bu"`
		Team string `form:"team"`
		Name string `form:"name"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	trees, err := apmSvc.Trees(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("apmSvc.appTree error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var data []string
	temp := make(map[string]string)
	if v.Bu == "" && v.Team == "" {
		for _, val := range trees {
			nameArr := strings.Split(val.Path, ".")
			newName := nameArr[1]
			if f := temp[newName]; f == "" {
				data = append(data, newName)
				temp[newName] = newName
			}
		}
	} else if v.Bu != "" && v.Team != "" && v.Name != "" {
		for _, val := range trees {
			nameArr := strings.Split(val.Path, ".")
			if v.Bu == nameArr[1] && v.Team == nameArr[2] && v.Name == nameArr[3] {
				data = append(data, strconv.Itoa(val.TreeID))
			}
		}
	} else if v.Bu != "" && v.Team != "" {
		for _, val := range trees {
			nameArr := strings.Split(val.Path, ".")
			if v.Bu == nameArr[1] && v.Team == nameArr[2] {
				newName := nameArr[1] + "." + nameArr[2] + "." + nameArr[3]
				if f := temp[newName]; f == "" {
					data = append(data, nameArr[3])
					temp[newName] = newName
				}
			}
		}
	} else if v.Bu != "" {
		for _, val := range trees {
			nameArr := strings.Split(val.Path, ".")
			if v.Bu == nameArr[1] {
				newName := v.Bu + "." + nameArr[2]
				if f := temp[newName]; f == "" {
					data = append(data, nameArr[2])
					temp[newName] = newName
				}
			}
		}
	}
	c.JSON(data, nil)
}

func appCallerSearch(c *bm.Context) {
	v := new(struct {
		AppID string `form:"app_id"`
		Pn    int    `form:"pn" default:"1"`
		Ps    int    `form:"ps" default:"20"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		aa    []*app.Auth
		count int
		lk    = "%" + v.AppID + "%"
	)
	if v.AppID != "" {
		err = apmSvc.DB.Where("app_id LIKE ?", lk).Group("app_tree_id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	} else {
		err = apmSvc.DB.Group("app_tree_id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&aa).Error
	}
	if err != nil {
		log.Error("apmSvc.appCallerSearch error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.AppID != "" {
		err = apmSvc.DB.Select("count(distinct(app_tree_id))").Where("app_id LIKE ?", lk).Model(&app.Auth{}).Count(&count).Error
	} else {
		err = apmSvc.DB.Select("count(distinct(app_tree_id))").Model(&app.Auth{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.appCallerSearch count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	type result struct {
		AppTreeID int64  `gorm:"column:app_tree_id" json:"app_tree_id"`
		AppID     string `gorm:"column:app_id" json:"app_id"`
		Services  []*app.Auth
	}
	var results []*result
	if count > 0 {
		var in []int64
		auth := []*app.Auth{}
		for _, val := range aa {
			in = append(in, val.AppTreeID)
			fmt.Printf("apptreeid=%v", val.AppTreeID)
		}
		apmSvc.DB.Where("app_tree_id in (?)", in).Find(&auth)
		for _, vv := range aa {
			rs := new(result)
			rs.AppTreeID = vv.AppTreeID
			rs.AppID = vv.AppID
			for _, vvv := range auth {
				if vv.AppTreeID == vvv.AppTreeID {
					rs.Services = append(rs.Services, vvv)
				}
			}
			results = append(results, rs)
		}
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: results,
		Total: count,
	}
	c.JSON(data, nil)
}
