package http

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

func simpleArchive(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	modeStr := params.Get("mode")
	// check params
	mode, _ := strconv.Atoi(modeStr)
	if mode <= 0 || mode > 2 {
		mode = 0 // 0 novideo 1 open video 2 all video
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	av, err := vdpSvc.SimpleArchive(c, aid, mode)
	if err != nil {
		log.Error(" vdpSvc.SimpleArchive(%d) error(%v)", aid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(av, nil)
}

func simpleVideos(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	vs, err := vdpSvc.SimpleVideos(c, aid)
	if err != nil {
		log.Error(" vdpSvc.SimpleVideos(%d) error(%v)", aid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(vs, nil)
}

// viewArchive get archive info.
func viewArchive(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	withPOIStr := params.Get("need_poi")
	withVoteStr := params.Get("need_vote")
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	av, err := vdpSvc.Archive(c, aid)
	if err != nil {
		log.Error(" vdpSvc.Archive(%d) error(%v)", aid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	//详情页展示LBS
	var poi *archive.PoiObj
	eg, errCtx := errgroup.WithContext(c)
	if withPOIStr == "1" {
		eg.Go(func() (err error) {
			poi, _ = vdpSvc.ArchivePOI(errCtx, aid)
			return nil
		})
	}
	var vote *archive.Vote
	if withVoteStr == "1" {
		eg.Go(func() (err error) {
			vote, _ = vdpSvc.ArchiveVote(errCtx, aid)
			return nil
		})
	}
	eg.Wait()
	if poi != nil {
		av.Archive.POI = poi
	}
	if vote != nil {
		av.Archive.Vote = vote
	}
	//详情页展示staff
	var staffs []*archive.StaffApply
	if staffs, err = vdpSvc.ApplysByAID(c, aid); err != nil || staffs == nil || len(staffs) == 0 {
		log.Error(" vdpSvc.ApplysByAID(%d) error(%v)", aid, err)
	} else {
		av.Archive.Staffs = staffs
	}
	c.JSON(av, nil)
}

// viewArchives get archive info.
func viewArchives(c *bm.Context) {
	params := c.Request.Form
	aidsStr := params.Get("aids")
	// check params
	aids, err := xstr.SplitInts(aidsStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aids) > 50 {
		log.Error("viewArchives aids(%s) too long, appkey(%s)", aidsStr, params.Get("appkey"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	avm, err := vdpSvc.Archives(c, aids)
	if err != nil {
		log.Error(" vdpSvc.Archive(%d) error(%v)", aids, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(avm, nil)
}

// archivesByMid get archive list by mid.
func archivesByMid(c *bm.Context) {
	params := c.Request.Form
	// check params
	midStr := params.Get("mid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	groupStr := params.Get("group")
	mid, _ := strconv.ParseInt(midStr, 10, 64)
	if mid <= 0 {
		log.Error("http.archivesByMid  mid(%d) <=0 ", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, _ := strconv.Atoi(pnStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 100 {
		ps = 10
	}
	group, _ := strconv.Atoi(groupStr)
	gp := int8(group)
	if gp < 0 || gp > 2 {
		gp = 0
	}
	uav, err := vdpSvc.UpArchives(c, mid, pn, ps, gp)
	if err != nil {
		log.Error(" vdpSvc.Archive(%d,%d,%d,%d) error(%v)", mid, pn, ps, gp, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(uav, nil)
}

// upArchiveTag add archive.
func upArchiveTag(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	tag := params.Get("tag")
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		return
	}
	c.JSON(nil, vdpSvc.UpTag(c, aid, tag))
}

// delArchive del archive.
func delArchive(c *bm.Context) {
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
		log.Error("http delArchive() json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ap.Aid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdpSvc.DelArchive(c, ap.Aid, ap.Mid))
}

// arcHistory get archive edit history.
func arcHistory(c *bm.Context) {
	params := c.Request.Form
	hidStr := params.Get("hid")
	// check params
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(hid(%s)) error(%v)", hidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(vdpSvc.ArcHistory(c, hid), nil)
}

// arcHistorys get archive edit history.
func arcHistorys(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(aid(%s)) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(vdpSvc.ArcHistorys(c, aid), nil)
}

// types get all types info
func types(c *bm.Context) {
	c.JSON(vdpSvc.Types(c), nil)
}

// videoBycid get video bid cid
func videoBycid(c *bm.Context) {
	params := c.Request.Form
	cidStr := params.Get("cid")
	// check params
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(cid(%s)) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	v, err := vdpSvc.VideoBycid(c, cid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(v, nil)
}

// archivesByCids get archives by cids
func archivesByCids(c *bm.Context) {
	params := c.Request.Form
	cidsStr := params.Get("cids")
	appkey := params.Get("appkey")
	// check params
	if appkey != config.DmVerifyKey {
		log.Warn("appkey(%s) invalid", appkey)
		c.JSON(nil, ecode.AppKeyInvalid)
		return
	}
	cids, err := xstr.SplitInts(cidsStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(cids) > 100 {
		log.Error("cids(%d) number gt 100", len(cids))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(vdpSvc.ArchivesByCids(c, cids), nil)
}

// flows get flow list
func flows(c *bm.Context) {
	c.JSON(vdpSvc.Flows(c), nil)
}

// flows get specialUps list
func specialUps(c *bm.Context) {
	params := c.Request.Form
	groupStr := params.Get("group_id")
	var err error
	// check params
	//default all groups
	groupID := int64(0)
	if groupStr != "" {
		groupID, err = strconv.ParseInt(groupStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(groupStr(%s)) error(%v)", groupStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	c.JSON(vdpSvc.UpsByGroup(c, groupID), nil)
}

// arcReasonTag .
func arcReasonTag(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	var (
		aid, tagID int64
		err        error
	)

	aid, err = strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(aidStr(%s)) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tagID, err = vdpSvc.ArcTag(c, aid)
	if err != nil {
		log.Error("vdpSvc.ArcTag error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(map[string]interface{}{
		"tag_id": tagID,
	}, nil)
}

//flowForbid
func flowForbid(c *bm.Context) {
	c.JSON(vdpSvc.UpsForbid(c), nil)
}

func appFeedAids(c *bm.Context) {
	aids, err := vdpSvc.AppFeedAids(c)
	if err != nil {
		log.Error("vdpSvc.AppFeedAids() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(aids, nil)
}

func descFormats(c *bm.Context) {
	dfs, err := vdpSvc.DescFormats(c)
	if err != nil {
		log.Error("vdpSvc.DescFormats() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(dfs, nil)
}

// videoJam video check traffic jam time
func videoJam(c *bm.Context) {
	level, _ := vdpSvc.VideoJamLevel(c)
	c.JSON(map[string]interface{}{
		"level": level,
	}, nil)
}

// archiveAddit get archive addit
func archiveAddit(c *bm.Context) {
	var (
		err   error
		addit *archive.Addit
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if addit, err = vdpSvc.ArchiveAddit(c, aid); err != nil {
		log.Error("vdpSvc.archiveAddit(%d) error(%v)", aid, err)
		c.JSON(nil, err)
		return
	}
	if addit == nil {
		err = ecode.NothingFound
		c.JSON(nil, err)
		return
	}
	c.JSON(addit, nil)
}

func rejectedArchives(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	stateStr := params.Get("state")
	startStr := params.Get("start")
	state, _ := strconv.Atoi(stateStr)
	if state >= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	start, err := getTimeFromSecStr(startStr)
	if err != nil || start == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := strconv.ParseInt(midStr, 10, 64)
	pn, _ := strconv.Atoi(pnStr)
	ps, _ := strconv.Atoi(psStr)
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 20
	}
	arcs, count, err := vdpSvc.RejectedArchives(c, mid, int32(state), int32(pn), int32(ps), start)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"page": map[string]int{
			"num":   pn,
			"size":  ps,
			"total": int(count),
		},
		"archives": arcs,
	}
	c.JSON(data, nil)
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
