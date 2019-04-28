package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/app/interface/main/creative/model/archive"
	whmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

func appH5ArcTags(c *bm.Context) {
	params := c.Request.Form
	tidStr := params.Get("typeid")
	title := params.Get("title")
	filename := params.Get("filename")
	desc := params.Get("desc")
	cover := params.Get("cover")
	midStr, ok := c.Get("mid")
	mid := midStr.(int64)
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	tags, _ := dataSvc.TagsWithChecked(c, mid, uint16(tid), title, filename, desc, cover, archive.TagPredictFromAPP)
	c.JSON(tags, nil)
}

func appH5ArcTagInfo(c *bm.Context) {
	params := c.Request.Form
	tagNameStr := params.Get("tag_name")
	midStr, ok := c.Get("mid")
	mid := midStr.(int64)
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if len(tagNameStr) == 0 {
		log.Error("tagNameStr len zero (%s)", tagNameStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	code, msg := arcSvc.TagCheck(c, mid, tagNameStr)
	c.JSON(map[string]interface{}{
		"code": code,
		"msg":  msg,
	}, nil)
}

func appH5Pre(c *bm.Context) {
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid := midStr.(int64)
	c.JSON(map[string]interface{}{
		"activities": arcSvc.Activities(c),
		"fav":        arcSvc.Fav(c, mid),
	}, nil)
}

func appH5MissionByType(c *bm.Context) {
	params := c.Request.Form
	tidStr := params.Get("tid")
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	tid, _ := strconv.Atoi(tidStr)
	actWithTP, _ := arcSvc.MissionOnlineByTid(c, int16(tid), 1)
	c.JSON(actWithTP, nil)
}

func toInt(s string) (i int, err error) {
	if s == "" {
		return 0, nil
	}
	i, err = strconv.Atoi(s)
	if err != nil {
		log.Error("strconv.Atoi s(%s) error(%v)", s, err)
		err = ecode.RequestErr
	}
	return
}

func toInt64(s string) (i int64, err error) {
	if s == "" {
		return 0, nil
	}
	i, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Error("strconv.Atoi s(%s) error(%v)", s, err)
		err = ecode.RequestErr

	}
	return
}

func h5ViewPlay(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	busStr := params.Get("business")

	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid == 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}

	aid, err := toInt64(aidStr)
	if err != nil || aid <= 0 {
		c.JSON(nil, err)
		return
	}

	bus, err := toInt(busStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	if aid == 0 || bus == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	pl, err := acaSvc.ViewPlay(c, mid, aid, int8(bus))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pl, nil)
}

func h5AddPlay(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	busStr := params.Get("business")
	watchStr := params.Get("watch")

	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid == 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}

	aid, err := toInt64(aidStr)
	if err != nil || aid <= 0 {
		c.JSON(nil, err)
		return
	}

	bus, err := toInt(busStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	watch, err := toInt(watchStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	if aid == 0 || bus == 0 || watch == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	id, err := acaSvc.PlayAdd(c, mid, aid, int8(bus), int8(watch))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func h5DelPlay(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	busStr := params.Get("business")

	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok || mid == 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}

	aid, err := toInt64(aidStr)
	if err != nil || aid <= 0 {
		c.JSON(nil, err)
		return
	}

	bus, err := toInt(busStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	if aid == 0 || bus == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	id, err := acaSvc.PlayDel(c, mid, aid, int8(bus))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func h5PlayList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")

	// check user
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, ok := midStr.(int64)
	if !ok || mid == 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}

	pn, err := toInt(pnStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	ps, err := toInt(psStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pn <= 0 {
		pn = 1
	}
	if ps > 20 || ps <= 0 {
		ps = 20
	}

	arcs, err := acaSvc.PlayList(c, mid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(arcs, nil)
}

func h5ThemeDir(c *bm.Context) {
	occs, err := acaSvc.Occupations(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(occs, nil)
}

func h5NewbCourse(c *bm.Context) {
	nc, err := acaSvc.NewbCourse(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nc, nil)
}

func h5Tags(c *bm.Context) {
	c.JSON(acaSvc.Tags(c), nil)
}

func h5Archive(c *bm.Context) {
	params := c.Request.Form
	tidsStr := params.Get("tids")
	bsStr := params.Get("business")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	keyword := params.Get("keyword")
	order := params.Get("order")
	drStr := params.Get("duration")
	ip := metadata.String(c, metadata.RemoteIP)

	var (
		tids []int64
		err  error
	)
	// check params
	if tidsStr != "" {
		if tids, err = xstr.SplitInts(tidsStr); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tidsStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	bs, err := toInt(bsStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	dr, err := toInt(drStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	pn, err := toInt(pageStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	ps, err := toInt(psStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pn <= 0 {
		pn = 1
	}
	if ps > 20 || ps <= 0 {
		ps = 20
	}

	aca := &academy.EsParam{
		Tid:      tids,
		Business: bs,
		Pn:       pn,
		Ps:       ps,
		Keyword:  keyword,
		Order:    order,
		IP:       ip,
		Duration: dr,
	}

	var arcs *academy.ArchiveList
	arcs, err = acaSvc.Archives(c, aca)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arcs, nil)
}

func h5Feature(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)

	// check params
	aca := &academy.EsParam{
		Pn:      1,
		Ps:      50,
		Keyword: "",
		Order:   "",
		IP:      ip,
	}

	aca2 := &academy.EsParam{
		Pn:      1,
		Ps:      50,
		Keyword: "",
		Order:   "",
		IP:      ip,
	}

	var (
		g               = &errgroup.Group{}
		offArcs, chArcs *academy.ArchiveList
	)

	g.Go(func() error {
		aca.Tid = []int64{acaSvc.OfficialID} //官方课程
		offArcs, _ = acaSvc.ArchivesWithES(c, aca)
		return nil
	})

	g.Go(func() error {
		aca2.Tid = []int64{acaSvc.EditorChoiceID} //编辑精选
		chArcs, _ = acaSvc.ArchivesWithES(c, aca2)
		return nil
	})
	g.Wait()

	c.JSON(map[string]interface{}{
		"official_course": offArcs,
		"editor_choice":   chArcs,
	}, nil)
}

func weeklyHonor(c *bm.Context) {
	midStr, _ := c.Get("mid")
	var mid int64
	uid, ok := midStr.(int64)
	if ok {
		mid = uid
	}
	arg := new(struct {
		UID   int64  `form:"uid"`
		Token string `form:"token"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	if mid == 0 && arg.UID == 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	rec, err := honorSvc.WeeklyHonor(c, mid, arg.UID, arg.Token)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rec, nil)
}

func weeklyHonorSubSwitch(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	stateStr := params.Get("state")
	st, err := strconv.Atoi(stateStr)
	state := uint8(st)
	if err != nil || (state != whmdl.HonorSub && state != whmdl.HonorUnSub) {
		c.JSON(nil, ecode.ReqParamErr)
	}
	err = honorSvc.ChangeSubState(c, mid, state)
	c.JSON(nil, err)
}

func h5RecommendV2(c *bm.Context) {
	midStr, _ := c.Get("mid")

	var mid int64
	uid, ok := midStr.(int64)
	if ok {
		mid = uid
	}

	rec, err := acaSvc.RecommendV2(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rec, nil)
}

func h5ThemeCousreV2(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pidStr := params.Get("pid")
	skidStr := params.Get("skid")
	sidStr := params.Get("sid")

	pn, err := toInt(pnStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	ps, err := toInt(psStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pn <= 0 {
		pn = 1
	}
	if ps > 20 || ps <= 0 {
		ps = 20
	}

	var pids, skids, sids []int64

	if pidStr != "" {
		if pids, err = xstr.SplitInts(pidStr); err != nil {
			log.Error("strconv.ParseInt pidStr(%s) error(%v)", pidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if skidStr != "" {
		if skids, err = xstr.SplitInts(skidStr); err != nil {
			log.Error("strconv.ParseInt skidStr(%s) error(%v)", skidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if sidStr != "" {
		if sids, err = xstr.SplitInts(sidStr); err != nil {
			log.Error("strconv.ParseInt sidStr(%s) error(%v)", sidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	arcs, err := acaSvc.ProfessionSkill(c, pids, skids, sids, pn, ps, false)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(arcs, nil)
}

func h5Keywords(c *bm.Context) {
	c.JSON(acaSvc.Keywords(c), nil)
}
