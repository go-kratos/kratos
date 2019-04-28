package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// videoAudit up firstRound info.
func videoAudit(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var vp = &archive.VideoParam{}
	if err = json.Unmarshal(bs, &vp); err != nil {
		log.Error("http firstRound() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if vp.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attrs := make(map[uint]int32, 6)
	attrs[archive.AttrBitNoRank] = vp.Attrs.NoRank
	attrs[archive.AttrBitNoDynamic] = vp.Attrs.NoDynamic
	attrs[archive.AttrBitNoSearch] = vp.Attrs.NoSearch
	attrs[archive.AttrBitNoRecommend] = vp.Attrs.NoRecommend
	attrs[archive.AttrBitOverseaLock] = vp.Attrs.OverseaLock
	attrs[archive.AttrBitPushBlog] = vp.Attrs.PushBlog
	if vp.TagID > 0 && vp.Status >= 0 {
		attrs[archive.AttrBitParentMode] = 1
	}
	c.JSON(nil, vdaSvc.VideoAudit(c, vp, attrs))
}

func batchVideo(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var vps = []*archive.VideoParam{}
	if err = json.Unmarshal(bs, &vps); err != nil {
		log.Error("http firstRound() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(vps) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := vdaSvc.CheckVideo(vps); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.BatchVideo(c, vps, archive.ActionVideoSubmit))
}

func upVideo(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var vp = &archive.VideoParam{}
	if err = json.Unmarshal(bs, &vp); err != nil {
		log.Error("http firstRound() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if vp.Aid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.UpVideo(c, vp))
}

func upWebLink(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var vp = &archive.VideoParam{}
	if err = json.Unmarshal(bs, &vp); err != nil {
		log.Error("http firstRound() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if vp.ID == 0 || vp.WebLink == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.UpWebLink(c, vp))
}

func delVideo(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var vp = &archive.VideoParam{}
	if err = json.Unmarshal(bs, &vp); err != nil {
		log.Error("http firstRound() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// TODO check data.
	if vp.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.DelVideo(c, vp))
}

func changeIndex(c *bm.Context) {
	req := c.Request
	// read
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	// params
	var lo = &archive.IndexParam{}
	if err = json.Unmarshal(bs, &lo); err != nil {
		log.Error("http changeIndex() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lo.Aid == 0 {
		log.Error("aid==%d", lo.Aid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.ChangeIndex(c, lo))
}
