package http

import (
	"strconv"

	"go-common/app/interface/main/playlist/conf"
	"go-common/app/interface/main/playlist/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func videoList(c *bm.Context) {
	var (
		pid    int64
		pn, ps int
		err    error
		list   *model.ArcList
	)
	params := c.Request.Form
	pidStr := params.Get("pid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxPlArcsPs {
		ps = conf.Conf.Rule.MaxPlArcsPs
	}
	if list, err = plSvc.Videos(c, pid, pn, ps); err != nil {
		c.JSON(nil, switchCode(err, favmdl.TypePlayVideo))
		return
	}
	c.JSON(list, nil)
}

func toView(c *bm.Context) {
	var (
		pid, mid int64
		err      error
		list     *model.ToView
	)
	params := c.Request.Form
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	pidStr := params.Get("pid")
	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, err = plSvc.ToView(c, mid, pid); err != nil {
		c.JSON(nil, switchCode(err, favmdl.TypePlayVideo))
		return
	}
	c.JSON(list, nil)
}

func check(c *bm.Context) {
	var (
		err      error
		mid, pid int64
		aids     []int64
		videos   model.Videos
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aids")
	if aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aids, err = xstr.SplitInts(aidStr); err != nil || len(aids) == 0 || len(aids) > conf.Conf.Rule.MaxArcChangeLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidMap := make(map[int64]int64, len(aids))
	for _, aid := range aids {
		aidMap[aid] = aid
	}
	if len(aidMap) < len(aids) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pidStr := params.Get("pid")
	if pidStr != "" {
		if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if videos, err = plSvc.CheckVideo(c, mid, pid, aids); err != nil {
		c.JSON(nil, switchCode(err, favmdl.TypePlayVideo))
		return
	}
	c.JSON(videos, nil)
}

func addVideo(c *bm.Context) {
	var (
		err      error
		mid, pid int64
		aids     []int64
		videos   model.Videos
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aids")
	if aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aids, err = xstr.SplitInts(aidStr); err != nil || len(aids) == 0 || len(aids) > conf.Conf.Rule.MaxArcChangeLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidMap := make(map[int64]int64, len(aids))
	for _, aid := range aids {
		aidMap[aid] = aid
	}
	if len(aidMap) < len(aids) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pidStr := params.Get("pid")
	if pidStr != "" {
		if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if videos, err = plSvc.AddVideo(c, mid, pid, aids); err != nil {
		c.JSON(nil, switchCode(err, favmdl.TypePlayVideo))
		return
	}
	c.JSON(videos, nil)
}

func delVideo(c *bm.Context) {
	var (
		err      error
		mid, pid int64
		aids     []int64
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aids")
	if aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aids, err = xstr.SplitInts(aidStr); err != nil || len(aids) == 0 || len(aids) > conf.Conf.Rule.MaxArcChangeLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidMap := make(map[int64]int64, len(aids))
	for _, aid := range aids {
		aidMap[aid] = aid
	}
	if len(aidMap) < len(aids) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pidStr := params.Get("pid")
	if pidStr != "" {
		if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(nil, switchCode(plSvc.DelVideo(c, mid, pid, aids), favmdl.TypePlayVideo))
}

func sortVideo(c *bm.Context) {
	var (
		mid, pid, aid, sort int64
		err                 error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pidStr := params.Get("pid")
	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sortStr := params.Get("sort")
	if sort, err = strconv.ParseInt(sortStr, 10, 64); err != nil || sort <= 0 || sort > int64(conf.Conf.Rule.MaxVideoCnt) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, switchCode(plSvc.SortVideo(c, mid, pid, aid, sort), favmdl.TypePlayVideo))
}

func editVideoDesc(c *bm.Context) {
	var (
		err           error
		mid, pid, aid int64
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	pidStr := params.Get("pid")
	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	desc := params.Get("desc")
	if len([]rune(desc)) > conf.Conf.Rule.MaxVideoDescLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, plSvc.EditVideoDesc(c, mid, pid, aid, desc))
}

func searchVideo(c *bm.Context) {
	var (
		err           error
		pn, ps, count int
		list          []*model.SearchArc
	)
	params := c.Request.Form
	pnStr := params.Get("pn")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	psStr := params.Get("ps")
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxSearchArcPs {
		ps = conf.Conf.Rule.MaxSearchArcPs
	}
	query := params.Get("keyword")
	if query == "" || len([]rune(query)) > conf.Conf.Rule.MaxSearchLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, count, err = plSvc.SearchVideos(c, pn, ps, query); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   pn,
		"size":  ps,
		"count": count,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}
