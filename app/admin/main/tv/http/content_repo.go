package http

import (
	"strings"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

func contList(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*model.ContentRepo
		count int64
		order = atoi(req.Get("order"))
		page  = atoi(req.Get("page"))
		size  = 20
	)
	if page == 0 {
		page = 1
	}
	db := contWhere(c)
	db.Model(&model.ContentRepo{}).Count(&count)
	if order == 1 {
		db = db.Order("tv_content.mtime ASC")
	} else {
		db = db.Order("tv_content.mtime DESC")
	}
	if err = db.Model(&model.ContentRepo{}).Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	for _, v := range items {
		v.MtimeFormat = tvSrv.TimeFormat(v.Mtime)
		v.Mtime = 0
	}
	pager := &model.ContentRepoPager{
		TotalCount: count,
		Pn:         page,
		Ps:         size,
		Items:      items,
	}
	c.JSON(pager, nil)
}

func contInfo(c *bm.Context) {
	var (
		req  = c.Request.Form
		epid = parseInt(req.Get("id"))
		err  error
	)
	exist := model.Content{}
	if err = tvSrv.DB.Where("epid=?", epid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(exist, nil)
}

func saveCont(c *bm.Context) {
	var (
		req  = c.Request.PostForm
		epid = atoi(req.Get("id"))
		err  error
	)
	exist := model.Content{}
	if err = tvSrv.DB.Where("epid=?", epid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	title := req.Get("title")
	cover := req.Get("cover")
	if cover == "" {
		renderErrMsg(c, ecode.RequestErr.Code(), "封面不能为空")
		return
	}
	if err := tvSrv.DB.Model(&model.Content{}).Where("epid = ?", epid).Update(map[string]string{"title": title, "cover": cover}).Error; err != nil {
		log.Error("tvSrv.saveCont error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err := tvSrv.DB.Model(&model.TVEpContent{}).Where("id = ?", epid).Update(map[string]string{"long_title": title, "cover": cover}).Error; err != nil {
		log.Error("tvSrv.saveCont error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func preview(c *bm.Context) {
	var (
		req  = c.Request.Form
		err  error
		epid = atoi(req.Get("id"))
	)
	exist := model.Content{}
	if err = tvSrv.DB.Where("epid=?", epid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	url, err := tvSrv.Playurl(exist.CID)
	if err != nil {
		log.Error("tvSrv.Playurl error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(url, nil)
}

func contOptions(id string, valid int) (ret bool) {
	var (
		epid  = atoi(id)
		exist = model.Content{}
	)
	ret = false
	if err := tvSrv.DB.Where("epid=?", epid).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		log.Error("tvSrv.contOptions error(%v)", err)
		return
	}
	if err := tvSrv.DB.Model(&model.Content{}).Where("epid=?", epid).Update(map[string]int{"valid": valid}).Error; err != nil {
		log.Error("tvSrv.contOptions error(%v)", err)
		return
	}
	return true
}

func contOnline(c *bm.Context) {
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
		if !contOptions(val, 1) {
			renderErrMsg(c, ecode.RequestErr.Code(), "Online("+val+") fail")
			return
		}
	}
	c.JSON(nil, nil)
}

func contHidden(c *bm.Context) {
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
		if !contOptions(val, 0) {
			renderErrMsg(c, ecode.RequestErr.Code(), "Hide ("+val+") fail")
			return
		}
	}
	c.JSON(nil, nil)
}

func contWhere(c *bm.Context) *gorm.DB {
	var (
		req      = c.Request.Form
		sid      = atoi(req.Get("sid"))
		cat      = atoi(req.Get("category"))
		epid     = atoi(req.Get("epid"))
		validStr = req.Get("valid")
	)
	db := tvSrv.DB.
		Joins("LEFT OUTER JOIN tv_ep_season ON tv_content.season_id=tv_ep_season.id").
		Select("tv_content.*, tv_ep_season.category, tv_ep_season.title AS season_title").
		Where("tv_content.state=?", 3).
		Where("tv_content.is_deleted=?", 0).
		Where("tv_ep_season.check=?", 1).
		Where("tv_ep_season.is_deleted=?", 0)
	if sid != 0 {
		db = db.Where("tv_content.season_id=?", sid)
	}
	if epid != 0 {
		db = db.Where("tv_content.epid=?", epid)
	}
	if cat != 0 {
		db = db.Where("tv_ep_season.category=?", cat)
	}
	if validStr == "" {
		return db
	}
	if valid := atoi(validStr); valid == 0 {
		db = db.Where("tv_content.valid=?", 0)
	} else if valid == 1 {
		db = db.Where("tv_content.valid=?", 1)
	}
	return db
}
