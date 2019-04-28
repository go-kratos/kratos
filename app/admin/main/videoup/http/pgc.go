package http

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// passByPGC update archive state pass by pgc
func passByPGC(c *bm.Context) {
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
	var ap struct {
		Aid         int64  `json:"aid"`
		Gid         int64  `json:"gid"`
		IsJump      int32  `json:"is_jump"`
		AllowBp     int32  `json:"allow_bp"`
		IsBangumi   int32  `json:"is_bangumi"`
		IsMovie     int32  `json:"is_movie"`
		BadgePay    int32  `json:"is_pay"`
		IsPGC       int32  `json:"is_pgc"`
		RedirectURL string `json:"redirect_url"`
	}
	if err = json.Unmarshal(bs, &ap); err != nil {
		log.Error("http upArchiveArr() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		log.Error("aid==0")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attrs := make(map[uint]int32, 7)
	attrs[archive.AttrBitJumpURL] = ap.IsJump
	attrs[archive.AttrBitAllowBp] = ap.AllowBp
	attrs[archive.AttrBitIsBangumi] = ap.IsBangumi
	attrs[archive.AttrBitIsMovie] = ap.IsMovie
	attrs[archive.AttrBitBadgepay] = ap.BadgePay
	attrs[archive.AttrBitIsPGC] = ap.IsPGC
	attrs[archive.AttrBitLimitArea] = 0
	if ap.Gid > 1 {
		attrs[archive.AttrBitLimitArea] = 1
	}
	c.JSON(nil, vdaSvc.PassByPGC(c, ap.Aid, ap.Gid, attrs, ap.RedirectURL, time.Now()))
}

// modifyByPGC update archive attr by pgc
func modifyByPGC(c *bm.Context) {
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
	var ap struct {
		Aid         int64  `json:"aid"`
		Gid         int64  `json:"gid"`
		IsJump      int32  `json:"is_jump"`
		AllowBp     int32  `json:"allow_bp"`
		IsBangumi   int32  `json:"is_bangumi"`
		IsMovie     int32  `json:"is_movie"`
		BadgePay    int32  `json:"is_pay"`
		IsPGC       int32  `json:"is_pgc"`
		RedirectURL string `json:"redirect_url"`
	}
	if err = json.Unmarshal(bs, &ap); err != nil {
		log.Error("http modArchiveArr() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		log.Error("aid==0")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	attrs := make(map[uint]int32, 7)
	attrs[archive.AttrBitJumpURL] = ap.IsJump
	attrs[archive.AttrBitAllowBp] = ap.AllowBp
	attrs[archive.AttrBitIsBangumi] = ap.IsBangumi
	attrs[archive.AttrBitIsMovie] = ap.IsMovie
	attrs[archive.AttrBitBadgepay] = ap.BadgePay
	attrs[archive.AttrBitIsPGC] = ap.IsPGC
	attrs[archive.AttrBitLimitArea] = 0
	if ap.Gid > 1 {
		attrs[archive.AttrBitLimitArea] = 1
	}

	c.JSON(nil, vdaSvc.ModifyByPGC(c, ap.Aid, ap.Gid, attrs, ap.RedirectURL))
}

// lockByPGC update archive state to lockbid by pgc
func lockByPGC(c *bm.Context) {
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
	var ap struct {
		Aid int64 `json:"aid"`
	}
	if err = json.Unmarshal(bs, &ap); err != nil {
		log.Error("http struct aid json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		log.Error("aid==0")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.LockByPGC(c, ap.Aid))
}
