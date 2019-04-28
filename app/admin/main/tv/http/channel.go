package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

func chlInfo(c *bm.Context) {
	var (
		req = c.Request.Form

		vid = parseInt(req.Get("id"))
		err error
	)
	exist := model.ChannelFmt{}
	if err = tvSrv.DB.Where("id=?", vid).Where("deleted!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	exist.MtimeFormat = tvSrv.TimeFormat(exist.Mtime)
	exist.Mtime = 0
	c.JSON(exist, nil)
}

func chlList(c *bm.Context) {
	param := new(model.ReqChannel)
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(tvSrv.ChlSplash(c, param))
}

func chlEdit(c *bm.Context) {
	var (
		req     = c.Request.PostForm
		vid     = parseInt(req.Get("id"))
		allowed bool
		err     error
	)
	exist := model.Channel{}
	if err = tvSrv.DB.Where("id=?", vid).Where("deleted!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	alert, simple := validateChl(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	if allowed, _ = nameExist(simple.Title, int(vid)); !allowed {
		renderErrMsg(c, ecode.RequestErr.Code(), "Title exists")
		return
	}
	if err = tvSrv.DB.Model(&model.Channel{}).Where("id=?", vid).Update(simple).Error; err != nil {
		log.Error("tvSrv.saveChannel error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func chlAdd(c *bm.Context) {
	var (
		err     error
		allowed bool
	)
	alert, simple := validateChl(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	if allowed, _ = nameExist(simple.Title, 0); !allowed {
		renderErrMsg(c, ecode.RequestErr.Code(), _errTitleExist)
		return
	}
	if err = tvSrv.DB.Create(simple).Error; err != nil {
		log.Error("tvSrv.addChannel error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func chlDel(c *bm.Context) {
	var (
		req = c.Request.PostForm

		vid = parseInt(req.Get("id"))
		err error
	)
	exist := model.Channel{}
	if err = tvSrv.DB.Where("id=?", vid).Where("deleted!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.DB.Model(&model.Channel{}).Where("id=?", vid).Update(map[string]int{"deleted": _isDeleted}).Error; err != nil {
		log.Error("tvSrv.chlDel error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// avoid two same names in DB - conflict
func nameExist(name string, myID int) (allowed bool, err error) {
	var (
		exist = model.Channel{}
		db    = tvSrv.DB.Where("title = ?", name).Where("deleted!=?", _isDeleted)
	)
	if myID != 0 {
		db = db.Where("id != ?", myID)
	}
	if err = db.First(&exist).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("tvSrv.nameExist error(%v)", err)
			return
		}
		return true, nil
	}
	return false, nil
}

// validate Channel params
func validateChl(c *bm.Context) (alert string, simple *model.Channel) {
	var (
		req    = c.Request.PostForm
		title  = req.Get("title")
		desc   = req.Get("desc")
		splash = req.Get("splash")
	)
	if title == "" {
		alert = "Channel Title can't be empty"
		return
	}
	if splash == "" {
		alert = "Splash can't be empty"
		return
	}
	return "", &model.Channel{
		Title:  title,
		Desc:   desc,
		Splash: splash,
	}
}
