package http

import (
	"strconv"

	"encoding/json"
	"fmt"
	"go-common/app/admin/main/creative/model/app"
	"go-common/app/admin/main/creative/model/logcli"
	Mamdl "go-common/app/admin/main/creative/model/material"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/xstr"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//素材库
func searchMaterialDb(c *bm.Context) {
	var (
		req           = c.Request.Form
		err           error
		items, dItems []*Mamdl.Material
		count         int64
		typeID        = req.Get("type")
		Platform      = atoi(req.Get("platform"))
		categoryID    = req.Get("category_id")
		id            = atoi(req.Get("id"))
		name          = req.Get("name")
		page          = atoi(req.Get("page"))
		size          = 20
	)
	if page == 0 {
		page = 1
	}
	db := svc.DB.Model(&Mamdl.Material{}).Where("material.state!=?", Mamdl.StateDelete).
		Joins("left join material_with_category on material_with_category.material_id=material.id and material_with_category.state !=?", Mamdl.StateOff).
		Joins("left join material_category on material_with_category.category_id=material_category.id and material_category.state !=?", Mamdl.StateOff).
		Select("material.*,material_category.name as category_name,material_with_category.category_id,material_with_category.index as category_index")
	if id != 0 {
		//id查询时需要提供type参数否则可能导致搜索为空
		db = db.Where("material.id=?", id)
	}
	if Platform != -1 {
		db = db.Where("material.platform=?", Platform)
	}
	//类型filter
	if categoryID != "" {
		db = db.Where("material_with_category.category_id=?", atoi(categoryID))
	}
	if typeID != "" {
		db = db.Where("material.type=?", atoi(typeID))
	}
	if name != "" {
		db = db.Where("material.name=?", name)
	}
	db.Count(&count)
	if categoryID != "" {
		db = db.Order("material_with_category.index")
	} else {
		db = db.Order("material.rank")
	}
	if err = db.Offset((page - 1) * size).Limit(size).Find(&dItems).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}

	if len(dItems) > 0 {
		items = make([]*Mamdl.Material, 0, len(dItems))
		for _, v := range dItems {
			i := &Mamdl.Material{}
			if v == nil {
				continue
			}
			i.ID = v.ID
			i.UID = v.UID
			i.Name = v.Name
			i.Extra = v.Extra
			i.Rank = v.Rank
			i.Type = v.Type
			i.Platform = v.Platform
			i.Build = v.Build
			i.State = v.State
			i.CategoryID = v.CategoryID
			i.CategoryIndex = v.CategoryIndex
			i.CategoryName = v.CategoryName
			if i.CategoryName == "" {
				i.CategoryID = 0
				i.CategoryIndex = 0
			}
			i.CTime = v.CTime
			i.MTime = v.MTime
			items = append(items, i)
		}
	} else {
		items = []*Mamdl.Material{}
	}
	pager := &Mamdl.Result{
		Items: items,
		Pager: &Mamdl.Pager{Num: page, Size: size, Total: count},
	}
	c.JSON(pager, nil)
}

func infoMaterial(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &Mamdl.Material{}
	if err = svc.DB.Where("id=?", id).First(&m).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]*Mamdl.Material{
		"data": m,
	}, nil)
}

func getUIDName(c *bm.Context) (uid int64, uname string) {
	unamei, ok := c.Get("username")
	if ok {
		uname = unamei.(string)
	}
	uidi, ok := c.Get("uid")
	if ok {
		uid = uidi.(int64)
	}
	return
}

func syncMaterial(c *bm.Context) {
	var (
		err    error
		m      *Mamdl.Material
		action string
	)
	mp := &Mamdl.Param{}
	if err = c.BindWith(mp, binding.Form); err != nil {
		log.Error("syncMaterial  bind error trace(%+v)", errors.Wrap(err, "sync bind  error"))
		httpCode(c, fmt.Sprintf("bind error(%v)", err), ecode.RequestErr)
		return
	}
	m, err = transferDB(c, mp)
	if err != nil {
		log.Error("transferDB  find error(%+v)", err)
		httpCode(c, fmt.Sprintf("transferDB validate error(%v)", err), ecode.RequestErr)
		return
	}
	uid, uname := getUIDName(c)
	m.UID = uid
	if m.ID == 0 {
		action = "add"
		if err = svc.DB.Create(m).Error; err != nil {
			log.Error("syncMaterial  Create error(%+v)", err)
			httpCode(c, fmt.Sprintf("syncMaterial  Create error(%v)", err), ecode.RequestErr)
			return
		}
	} else {
		action = "edit"
		exist := &Mamdl.Material{}
		if err = svc.DB.Where("id=?", mp.ID).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("syncMaterial  find error(%+v)", err)
			httpCode(c, fmt.Sprintf("syncMaterial  find error(%v)", err), ecode.RequestErr)
			return
		}
		if exist.ID > 0 {
			m.ID = exist.ID
			if err = svc.DB.Model(&Mamdl.Material{}).Where("id=?", mp.ID).Update(m).Update(map[string]int8{"type": m.Type}).Error; err != nil {
				log.Error("syncMaterial update error(%+v)", err)
				httpCode(c, fmt.Sprintf("syncMaterial update error(%v)", err), ecode.RequestErr)
				return
			}
		}
	}
	logMaterial, _ := json.Marshal(mp)
	svc.SendMusicLog(c, logcli.LogClientArchiveMaterialType, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: action, Name: string(logMaterial)})
	//不做事务允许分类绑定异常
	if mp.Type == Mamdl.TypeFilter || mp.Type == Mamdl.TypeCreativeSticks {
		_, err = svc.BindWithCategory(c, m.ID, mp.CategoryID, mp.CategoryIndex)
		if err != nil {
			httpCode(c, fmt.Sprintf("BindWithCategory update error(%v)", err), ecode.RequestErr)
			return
		}
	}

	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func transferDB(c *bm.Context, form *Mamdl.Param) (material *Mamdl.Material, err error) {
	var extraBytes []byte
	material = &Mamdl.Material{}
	if form == nil {
		return material, fmt.Errorf("params is wrong")
	}
	if form.Type != Mamdl.TypeCreativeSticks && (form.Name == "") {
		return material, fmt.Errorf("name is need")
	}
	if form.Type == Mamdl.TypeSubTitle && (form.Max == 0) {
		return material, fmt.Errorf("max is need")
	}
	if form.Type != Mamdl.TypeSticksIcon && (form.Rank == 0) {
		return material, fmt.Errorf("rank is need")
	}
	//不需要cover
	if form.Type != Mamdl.TypeHotWord && form.Type != Mamdl.TypeSticksIcon && form.Type != Mamdl.TypeCooperate && (form.Cover == "" || form.DownloadURL == "") {
		return material, fmt.Errorf("cover and download_url is need")
	}
	if form.Cover != "" && (!strings.HasPrefix(form.Cover, "http") || !strings.Contains(form.Cover, "/bfs/")) {
		return material, fmt.Errorf("cover is wrong")
	}
	if form.Type == Mamdl.TypeCooperate && form.DownloadURL != "" && !strings.Contains(form.DownloadURL, "acgvideo.com/") {
		//acg
		return material, fmt.Errorf("download_url is wrong")
	}
	if form.Type != Mamdl.TypeCooperate && form.DownloadURL != "" && (!strings.HasSuffix(form.DownloadURL, ".zip") || !strings.HasPrefix(form.DownloadURL, "http") || !strings.Contains(form.DownloadURL, "/bfs/creative/")) {
		//bfs creative
		return material, fmt.Errorf("download_url is wrong")
	}
	if !Mamdl.InMaterialType(form.Type) || form.Type == Mamdl.TypeBGM {
		//Mamdl.TypeBGM 为迁入bgm预留
		return material, fmt.Errorf("type is wrong")
	}
	if !checkBuild(form.Build) {
		return material, fmt.Errorf("build is wrong")
	}
	switch form.Type {
	case Mamdl.TypeSubTitle:
		//字幕库
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"max":          form.Max,
			"cover":        form.Cover,
		})
	case Mamdl.TypeFont, Mamdl.TypeHotWord, Mamdl.TypeSticksIcon:
		//字体库,贴纸，热词,贴纸Icon
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
		})
	case Mamdl.TypeSticks:
		//贴纸
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
			"sub_type":     form.SubType,
			"white_list":   form.WhilteList,
			"tip":          form.Tip,
		})
	case Mamdl.TypeFilter:
		//滤镜库
		if len(form.ExtraURL) > 0 && (!strings.HasSuffix(form.ExtraURL, ".zip") || !strings.HasPrefix(form.ExtraURL, "http") || !strings.Contains(form.ExtraURL, "/bfs/creative/")) {
			return material, fmt.Errorf("extra_url is wrong")
		}
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
			"extra_url":    form.ExtraURL,
			"extra_field":  form.ExtraField,
			"filter_type":  form.FilterType,
		})
	case Mamdl.TypeCooperate:
		//合拍库
		if form.MaterialAID*form.MaterialCID*form.DemoAID == 0 || form.SubType == 0 {
			return material, fmt.Errorf("参数错误")
		}
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
			"sub_type":     form.SubType,
			"style":        form.Style,
			"material_aid": form.MaterialAID,
			"material_cid": form.MaterialCID,
			"demo_aid":     form.DemoAID,
			"demo_cid":     form.DemoCID,
			"mission_id":   form.MissionID,
		})
	case Mamdl.TypeTheme:
		//主题库
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
		})
	default:
		extraBytes, err = json.Marshal(map[string]interface{}{
			"download_url": form.DownloadURL,
			"cover":        form.Cover,
		})
	}
	if err != nil {
		return material, err
	}
	material.ID = form.ID
	material.Name = form.Name
	material.Extra = string(extraBytes)
	material.Type = form.Type
	material.Rank = form.Rank
	material.Platform = form.Platform
	material.Build = form.Build
	if form.CategoryID > 0 {
		if form.CategoryIndex < 1 {
			return material, fmt.Errorf("category_index is wrong")
		}
		cate, _ := svc.CategoryByID(c, form.CategoryID)
		if cate == nil || cate.State == Mamdl.StateOff {
			return material, fmt.Errorf("category_id is wrong")
		}
	}
	return

}

func stateMaterial(c *bm.Context) {
	var (
		req   = c.Request.PostForm
		ids   []int64
		state = parseInt(req.Get("state"))
		err   error
	)
	idStr := req.Get("id")
	if ids, err = xstr.SplitInts(idStr); err != nil {
		log.Error("stateMaterial strconv.ParseInt(%s) error(%v)", idStr, err)
		httpCode(c, fmt.Sprintf("stateMaterial strconv.ParseInt error(%v)", err), ecode.RequestErr)
		return
	}
	if state > 2 {
		httpCode(c, "state参数可选值为 0 上架，1下架，2 删除", ecode.RequestErr)
		return
	}
	//删除时保护上架状态的
	if state == 2 {
		if err = svc.DB.Model(Mamdl.Material{}).Where("id IN (?)", ids).Where("state !=?", 0).Update(map[string]int64{"state": state}).Error; err != nil {
			log.Error("svc.stateMaterial error(%v)", err)
			httpCode(c, fmt.Sprintf("stateMaterial update state error(%v)", err), ecode.RequestErr)
			return
		}
	} else {
		if err = svc.DB.Model(Mamdl.Material{}).Where("id IN (?)", ids).Update(map[string]int64{"state": state}).Error; err != nil {
			log.Error("svc.stateMaterial error(%v)", err)
			httpCode(c, fmt.Sprintf("stateMaterial update state error(%v)", err), ecode.RequestErr)
			return
		}
	}

	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func parseInt(value string) int64 {
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		intval = 0
	}
	return intval
}

func atoi(value string) (intval int) {
	intval, err := strconv.Atoi(value)
	if err != nil {
		intval = 0
	}
	return intval
}

func httpCode(c *bm.Context, message string, err error) {
	c.JSON(map[string]interface{}{
		"message": message,
	}, err)
}

func checkBuild(build string) bool {
	var err error
	if len(build) > 0 {
		type buildItem struct {
			Build      int64 `json:"build"`
			Conditions int8  `json:"conditions"`
		}
		//比较版本号符号类型,0-等于,1-小于,2-大于,3-不等于,4-小于等于,5-大于等于
		//[{"conditions":0,"build":5290000},{"conditions":0,"build":5290000}]
		var buildExp []*buildItem
		if err = json.Unmarshal([]byte(build), &buildExp); err != nil {
			return false
		}
	}
	return true
}

func checkWhite(white string) bool {
	var err error
	if len(white) > 0 {
		var whiteExp []*app.WhiteExp
		if err = json.Unmarshal([]byte(white), &whiteExp); err != nil {
			return false
		}
	}
	return true
}
