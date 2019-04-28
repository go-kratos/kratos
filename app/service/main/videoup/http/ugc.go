package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// addArchive add archive.
func addArchive(c *bm.Context) {
	req := c.Request
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal(bs, ap); err != nil {
		log.Error("http addArchive() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || len(ap.Videos) == 0 {
		log.Error("http addArchive() func param is empty (%d)|(%d)|(%d)|(%s)|(%s)|(%s)|(%v)", ap.Mid, ap.TypeID, ap.Copyright, ap.Title, ap.Tag, ap.Desc, ap.Videos)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdpSvc.AllowType(c, ap.TypeID); !ok {
		log.Error("typeId (%d) not exist", ap.TypeID)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := archive.InCopyrights(ap.Copyright); !ok {
		log.Error("Copyright (%d) not in (1,2)", ap.Copyright)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := vdpSvc.AddByUGC(c, ap)
	if err != nil {
		log.Error("vdpSvc.AddArchive() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": aid,
	}, nil)
}

// editArchive edit archive.
func editArchive(c *bm.Context) {
	req := c.Request
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var ap = &archive.ArcParam{}
	if err = json.Unmarshal(bs, ap); err != nil {
		log.Error("http editArchive() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 || ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" {
		log.Error("http editArchive() func param is empty (%v)", ap)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdpSvc.AllowType(c, ap.TypeID); !ok {
		log.Error("typeId (%d) not exist", ap.TypeID)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := archive.InCopyrights(ap.Copyright); !ok {
		log.Error("Copyright (%d) not in (1,2)", ap.Copyright)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vdpSvc.EditByUGC(c, ap); err != nil {
		log.Error("vdpSvc.EditByUGC() error(%v)", err)
		c.JSON(nil, err)
		return
	}

	c.JSON(nil, nil)
}

// editMissionByUGC edit archive mission by UGC.
func editMissionByUGC(c *bm.Context) {
	v := &archive.ArcMissionParam{}
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, vdpSvc.EditMissionByUGC(c, v))
}

// queryDynamic edit archive mission by UGC.
func queryDynamic(c *bm.Context) {
	v := &archive.ArcDynamicParam{}
	if err := c.Bind(v); err != nil {
		return
	}
	can, err := vdpSvc.QueryDynamicSetting(c, v)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var repostBanned = int8(0)
	if !can {
		repostBanned = 1
	}
	c.JSON(map[string]int8{
		"repost_banned": repostBanned,
	}, err)
}
