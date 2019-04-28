package http

import (
	"strconv"

	"go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// stat get stat.
func stat(c *bm.Context) {
	var (
		err    error
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(relationSvc.Stat(c, mid))
}

// stat get user's follower list.
func stats(c *bm.Context) {
	var (
		err     error
		params  = c.Request.Form
		midsStr = params.Get("mids")
	)
	mids, err := xstr.SplitInts(midsStr)
	if err != nil || len(mids) > 20 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(relationSvc.Stats(c, mids))
}

// setStat set stat.
func setStat(c *bm.Context) {
	var (
		err         error
		mid         int64
		f, w, b, fr int64
		st          *model.Stat
		params      = c.Request.Form
		midStr      = params.Get("mid")
		fStr        = params.Get("following")
		wStr        = params.Get("whisper")
		bStr        = params.Get("black")
		frStr       = params.Get("follower")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fStr == "" {
		f = -1
	} else {
		if f, err = strconv.ParseInt(fStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if wStr == "" {
		w = -1
	} else {
		if w, err = strconv.ParseInt(wStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if bStr == "" {
		b = -1
	} else {
		if b, err = strconv.ParseInt(bStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if frStr == "" {
		fr = -1
	} else {
		if fr, err = strconv.ParseInt(frStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	st = &model.Stat{Following: f, Whisper: w, Black: b, Follower: fr}
	if st.Following == -1 && st.Whisper == -1 && st.Black == -1 && st.Follower == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.SetStat(c, mid, st))
}

// delStatCache del stat cache.
func delStatCache(c *bm.Context) {
	var (
		err    error
		mid    int64
		params = c.Request.Form
		midStr = params.Get("mid")
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, relationSvc.DelStatCache(c, mid))
}
