package http

import (
	"strings"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

func seasonList(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*model.SeaRepoDB
		count int64
		order = atoi(req.Get("order"))
		page  = atoi(req.Get("page"))
		size  = 20
	)
	if page == 0 {
		page = 1
	}
	db := seasonWhere(c)
	db.Model(&model.SeaRepoDB{}).Count(&count)
	if order == 1 {
		db = db.Order("mtime ASC")
	} else {
		db = db.Order("mtime DESC")
	}
	if err = db.Model(&model.SeaRepoDB{}).Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &model.SeasonRepoPager{
		TotalCount: count,
		Pn:         page,
		Ps:         size,
	}
	for _, v := range items {
		pager.Items = append(pager.Items, v.ToList())
	}
	c.JSON(pager, nil)
}

func seasonInfo(c *bm.Context) {
	var (
		req = c.Request.Form

		sid = parseInt(req.Get("id"))
		err error
	)
	exist := model.TVEpSeason{}
	if err = tvSrv.DB.Where("id=?", sid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(exist, nil)
}

func saveSeason(c *bm.Context) {
	var (
		req = c.Request.PostForm

		sid = parseInt(req.Get("id"))
		err error
	)
	exist := model.TVEpSeason{}
	if err = tvSrv.DB.Where("id=?", sid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	title := req.Get("title")
	desc := req.Get("desc")
	staff := req.Get("staff")
	cover := req.Get("cover")
	if title == "" {
		renderErrMsg(c, ecode.RequestErr.Code(), "标题不能为空")
		return
	}
	if desc == "" {
		renderErrMsg(c, ecode.RequestErr.Code(), "简介不能为空")
		return
	}
	if staff == "" {
		renderErrMsg(c, ecode.RequestErr.Code(), "staff不能为空")
		return
	}
	if cover == "" {
		renderErrMsg(c, ecode.RequestErr.Code(), "封面不能为空")
		return
	}
	if err := tvSrv.DB.Model(&model.TVEpSeason{}).Where("id = ?", sid).Update(map[string]string{"title": title, "desc": desc, "staff": staff, "cover": cover}).Error; err != nil {
		log.Error("tvSrv.saveSeason error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func seasonOptions(id string, valid int) (ret bool) {
	var (
		sid   = parseInt(id)
		exist = model.TVEpSeason{}
	)
	ret = false
	if err := tvSrv.DB.Where("id=?", sid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		log.Error("tvSrv.seasonOptions error(%v)", err)
		return
	}
	if err := tvSrv.DB.Model(&model.TVEpSeason{}).Where("id=?", sid).Update(map[string]int{"valid": valid}).Error; err != nil {
		log.Error("tvSrv.seasonOptions error(%v)", err)
		return
	}
	return true
}

func seasonOnline(c *bm.Context) {
	var (
		req = c.Request.PostForm

		ids = req.Get("ids")
	)
	idList := strings.Split(ids, ",")
	if len(idList) == 0 {
		renderErrMsg(c, ecode.RequestErr.Code(), _errIDNotFound)
		return
	}
	for _, val := range idList {
		if !seasonOptions(val, 1) {
			renderErrMsg(c, ecode.RequestErr.Code(), "Online("+val+") fail")
			return
		}
	}
	c.JSON(nil, nil)
}

func seasonHidden(c *bm.Context) {
	var (
		req = c.Request.PostForm

		ids = req.Get("ids")
	)
	idList := strings.Split(ids, ",")
	if len(idList) == 0 {
		renderErrMsg(c, ecode.RequestErr.Code(), _errIDNotFound)
		return
	}
	for _, val := range idList {
		if !seasonOptions(val, 0) {
			renderErrMsg(c, ecode.RequestErr.Code(), "Hidden("+val+") fail")
			return
		}
	}
	c.JSON(nil, nil)
}

func seasonWhere(c *bm.Context) *gorm.DB {
	var (
		req      = c.Request.Form
		sid      = atoi(req.Get("sid"))
		cat      = atoi(req.Get("category"))
		validStr = req.Get("valid")
		title    = req.Get("title")
	)
	db := tvSrv.DB.Select("*").
		Where("`check`=?", 1).
		Where("is_deleted=?", 0)
	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if sid != 0 {
		db = db.Where("id=?", sid)
	}
	if cat != 0 {
		db = db.Where("category=?", cat)
	}
	if validStr == "" {
		return db
	}
	if valid := atoi(validStr); valid == 0 {
		db = db.Where("valid=?", 0)
	} else if valid == 1 {
		db = db.Where("valid=?", 1)
	}
	return db
}
