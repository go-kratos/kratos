package http

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

type playerInfoc struct {
	ip        string
	ctime     string
	api       string
	buvid     string
	mid       string
	client    string
	itemID    string
	displayID string
	errCode   string
	from      string
	build     string
	trackID   string
}

const (
	_api    = "player"
	_client = "web"
)

var (
	emptyByte    = []byte{}
	_headerBuvid = "Buvid"
	_buvid       = "buvid3"
)

func player(c *bm.Context) {
	var (
		aid, cid, mid int64
		sid           model.Sid
		sidCookie     *http.Cookie
		err           error
		ip            = metadata.String(c, metadata.RemoteIP)
	)
	request := c.Request
	params := request.Form
	idStr := params.Get("id")
	aidStr := params.Get("aid")
	cidStr := params.Get("cid")
	refer := request.Referer()
	// response header
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("X-Remote-IP", ip)
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if strings.Index(idStr, "cid:") == 0 {
		cidStr = idStr[4:]
	}
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// get aid
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sidCookie, err = request.Cookie("sid")
	if err != nil || sidCookie.Value == "" {
		http.SetCookie(c.Writer, &http.Cookie{Name: "sid", Value: string(sid.Create()), Path: "/", Domain: ".bilibili.com"})
	}
	if sidCookie != nil {
		if sidCookie.Value != "" {
			sidStr := sidCookie.Value
			if !model.Sid(sidStr).Valid() {
				http.SetCookie(c.Writer, &http.Cookie{Name: "sid", Value: string(sid.Create()), Path: "/", Domain: ".bilibili.com"})
			}
		}
	}
	var (
		ips      = strings.Split(c.Request.Header.Get("X-Forwarded-For"), ",")
		cdnIP    string
		playByte []byte
	)
	if len(ips) >= 2 {
		cdnIP = ips[1]
	}
	buvid := request.Header.Get(_headerBuvid)
	if buvid == "" {
		cookie, _ := request.Cookie(_buvid)
		if cookie != nil {
			buvid = cookie.Value
		}
	}
	now := time.Now()
	if playInfoc != nil {
		showInfoc(ip, now, buvid, aid, mid)
	}
	if playByte, err = playSvr.Player(c, mid, aid, cid, cdnIP, refer, now); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Bytes(http.StatusOK, "text/plain; charset=utf-8", playByte)
}

func carousel(c *bm.Context) {
	// response header
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	type msg struct {
		Items []*model.Item `xml:"item"`
	}
	var (
		items []*model.Item
		err   error
	)
	if items, err = playSvr.Carousel(c); err != nil {
		c.Bytes(http.StatusOK, "text/xml; charset=utf-8", emptyByte)
		return
	}
	result := msg{Items: items}
	output, err := xml.MarshalIndent(result, "  ", "    ")
	if err != nil {
		output = emptyByte
	}
	c.Bytes(http.StatusOK, "text/xml; charset=utf-8", output)
}

func policy(c *bm.Context) {
	var (
		id  int64
		mid int64
		err error
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	params := c.Request.Form
	idStr := params.Get("id")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil || id < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(playSvr.Policy(c, id, mid))
}

func view(c *bm.Context) {
	var (
		aid int64
		err error
	)
	aidStr := c.Request.Form.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(playSvr.View(c, aid))
}

func matPage(c *bm.Context) {
	c.JSON(playSvr.Matsuri(c, time.Now()), nil)
}

func showInfoc(ip string, now time.Time, buvid string, aid, mid int64) {
	if playInfoc == nil {
		return
	}
	info := &playerInfoc{
		ip:        ip,
		ctime:     strconv.FormatInt(now.Unix(), 10),
		api:       _api,
		buvid:     buvid,
		mid:       strconv.FormatInt(mid, 10),
		client:    _client,
		itemID:    strconv.FormatInt(aid, 10),
		displayID: "",
		errCode:   "",
		from:      "",
		build:     "",
		trackID:   "",
	}
	err := playInfoc.Info(info.ip, info.ctime, info.api, info.buvid, info.mid, info.client, info.itemID, info.displayID, info.errCode, info.from, info.build, info.trackID)
	log.Info("playInfoc.Info(%s,%s,%s,%s,%s) error(%v)", info.ip, info.ctime, info.buvid, info.mid, info.itemID, err)
}
