package http

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func webArchVideos(c *bm.Context) {
	params := c.Request.Form
	ck := c.Request.Header.Get("Cookie")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	ap, err := arcSvc.SimpleArchiveVideos(c, mid, aid, "", ck, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ap, nil)
}

func webRecommandCover(c *bm.Context) {
	params := c.Request.Form
	fnsStr := params.Get("fns")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	cvs, err := dataSvc.Covers(c, mid, strings.Split(fnsStr, ","), metadata.String(c, metadata.RemoteIP))
	if err != nil {
		log.Error(" arcSvc.CoverList fnsStr(%s), mid(%d), ip(%s) error(%v)", fnsStr, mid, metadata.String(c, metadata.RemoteIP), err)
		c.JSON(nil, err)
		return
	}
	c.JSON(cvs, nil)
}

func webMissionProtocol(c *bm.Context) {
	params := c.Request.Form
	missionIDStr := params.Get("msid")
	missionID, err := strconv.ParseInt(missionIDStr, 10, 64)
	if err != nil || missionID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	missionProtocol, err := arcSvc.MissionProtocol(c, missionID)
	if err != nil {
		log.Error(" arcSvc.MissionProtocol missionID(%s), mid(%d), ip(%s) error(%v)", missionIDStr, mid, metadata.String(c, metadata.RemoteIP), err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"protocol": missionProtocol.Protocol,
	}, nil)
}

func webUgcPayProtocol(c *bm.Context) {
	params := c.Request.Form
	protocolID := params.Get("protocol_id")
	if len(protocolID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	payProtocol, err := paySvc.Protocol(c, protocolID)
	if err != nil {
		log.Error("paySvc.Protocol protocolID(%s) error(%v)", protocolID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"protocol": payProtocol,
	}, nil)
}

func voteAcsByTime(c *bm.Context) {
	params := c.Request.Form
	startStr := params.Get("start")
	endStr := params.Get("end")
	start, err := getTimeFromSecStr(startStr)
	if err != nil || start == nil || start.After(time.Now()) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	end, err := getTimeFromSecStr(endStr)
	if err != nil || start == nil || !isInHalfYear(start, end) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.BIZsByTime(c, start, end, archive.VoteType))
}

func getTimeFromSecStr(secStr string) (t *time.Time, err error) {
	sec, err := strconv.ParseInt(secStr, 10, 64)
	if err != nil || sec <= 0 {
		return
	}
	ti := time.Unix(sec, 0)
	t = &ti
	return
}

func isInHalfYear(start, end *time.Time) bool {
	return end.Unix() <= start.Add(6*31*24*time.Hour).Unix()
}
