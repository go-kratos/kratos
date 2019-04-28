package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// addByPGC add archive info by PGC.
func addByPGC(c *bm.Context) {
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
		log.Error("http addByPGC() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || ap.Desc == "" || len(ap.Videos) == 0 {
		log.Error("addByPGC func param is empty (%v)", ap)
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
	ap.UpFrom = archive.UpFromPGC
	aid, err := vdpSvc.AddByPGC(c, ap)
	if err != nil {
		log.Error("vdpSvc.AddByPGC() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": aid,
	}, nil)
}

// saddByPGC add secret archive info by PGC.
func saddByPGC(c *bm.Context) {
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
		log.Error("http addByPGC() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || ap.Desc == "" || len(ap.Videos) == 0 {
		log.Error("addByPGC func param is empty (%v)", ap)
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
	ap.UpFrom = archive.UpFromSecretPGC
	aid, err := vdpSvc.AddByPGC(c, ap)
	if err != nil {
		log.Error("vdpSvc.AddByPGC() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": aid,
	}, nil)
}

// editByPGC edit archive info by PGC.
func editByPGC(c *bm.Context) {
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
		log.Error("http editByPGC() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 || ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || ap.Desc == "" {
		log.Error("editByPGC func param is empty (%v)", ap)
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
	ap.UpFrom = archive.UpFromPGC
	c.JSON(nil, vdpSvc.EditByPGC(c, ap))
}

// caddByPGC add coopera archive info by PGC.
func caddByPGC(c *bm.Context) {
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
		log.Error("http addByPGC() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || ap.Desc == "" || len(ap.Videos) == 0 {
		log.Error("addByPGC func param is empty (%v)", ap)
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
	ap.UpFrom = archive.UpFromCoopera
	aid, err := vdpSvc.AddByPGC(c, ap)
	if err != nil {
		log.Error("vdpSvc.AddByPGC() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int64{
		"aid": aid,
	}, nil)
}

// ceditByPGC edit coopera archive info by PGC.
func ceditByPGC(c *bm.Context) {
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
		log.Error("http editByPGC() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 || ap.Mid == 0 || ap.TypeID == 0 || ap.Copyright == 0 || ap.Title == "" || ap.Tag == "" || ap.Desc == "" {
		log.Error("editByPGC func param is empty (%v)", ap)
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
	ap.UpFrom = archive.UpFromCoopera
	c.JSON(nil, vdpSvc.EditByPGC(c, ap))
}
