package http

import (
	"strconv"

	"go-common/app/service/main/location/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// pgcZone get ip infos by gids.
func pgcZone(c *bm.Context) {
	var (
		params  = c.Request.Form
		err     error
		zoneID  int64
		zoneIDs = []int64{}
		ip      *model.InfoComplete
	)
	zoneStr := params.Get("zone_id")
	ipStr := params.Get("ip")
	if zoneStr == "" && ipStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if zoneID, err = strconv.ParseInt(zoneStr, 10, 64); err != nil || zoneID == 0 {
		if ip, err = svr.InfoComplete(c, ipStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		} else {
			zoneIDs = ip.ZoneID
		}
	}
	if len(zoneIDs) == 0 {
		if zoneID != 0 {
			zoneIDs = append(zoneIDs, 0)
		}
		zoneIDs = append(zoneIDs, zoneID)
	}
	c.JSON(svr.PgcZone(c, zoneIDs))
}

// auth get auth by aid & ip & cip & mid.
func auth(c *bm.Context) {
	var (
		params        = c.Request.Form
		err           error
		ipaddr, cdnip string
		mid, aid      int64
	)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, ok := c.Get("mid"); !ok {
		mid = 0
	} else {
		mid = vmid.(int64)
	}
	ipaddr = params.Get("ip")
	cdnip = params.Get("cdnip")
	if ipaddr == "" && cdnip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Auth(c, aid, mid, ipaddr, cdnip))
}

// archive2 get auth by aid & ip & cip & mid.
func archive2(c *bm.Context) {
	var (
		params        = c.Request.Form
		err           error
		ipaddr, cdnip string
		mid, aid      int64
	)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, ok := c.Get("mid"); !ok {
		mid = 0
	} else {
		mid = vmid.(int64)
	}
	ipaddr = params.Get("ip")
	cdnip = params.Get("cdnip")
	if ipaddr == "" && cdnip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Archive2(c, aid, mid, ipaddr, cdnip))
}

// authGID get auth by gids & ip & cid & mid.
func authGID(c *bm.Context) {
	var (
		mid, gid      int64
		ipaddr, cdnip string
		params        = c.Request.Form
		err           error
	)
	gidStr := params.Get("gid")
	if gid, err = strconv.ParseInt(gidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, ok := c.Get("mid"); !ok {
		mid = 0
	} else {
		mid = vmid.(int64)
	}
	ipaddr = params.Get("ip")
	cdnip = params.Get("cdnip")
	if ipaddr == "" && cdnip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.AuthGID(c, gid, mid, ipaddr, cdnip), nil)
}

// authGIDs get auth by gids & ip & cid & mid.
func authGIDs(c *bm.Context) {
	var (
		gids          []int64
		mid           int64
		ipaddr, cdnip string
		params        = c.Request.Form
		err           error
	)
	gidsStr := params.Get("gids")
	if gids, err = xstr.SplitInts(gidsStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, ok := c.Get("mid"); !ok {
		mid = 0
	} else {
		mid = vmid.(int64)
	}
	ipaddr = params.Get("ip")
	cdnip = params.Get("cdnip")
	if ipaddr == "" && cdnip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.AuthGIDs(c, gids, mid, ipaddr, cdnip), nil)
}
