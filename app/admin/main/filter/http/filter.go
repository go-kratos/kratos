package http

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_baseArea    = "common"
	_paddingArea = "padding"
	_normal      = 0
	_deleted     = 1
	_expired     = 2
)

func filterRuleByID(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt() err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		filterInfo    *model.FilterInfo
		managerFilter = &model.FilterForGet{}
	)
	filterInfo, err = svc.AdminRuleByID(c, id)
	if filterInfo == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// 兼容前端格式
	managerFilter.LoadFromFilter(filterInfo)
	c.JSON(managerFilter, err)
}

func filterAdd(c *bm.Context) {
	var (
		err     error
		mode    int64
		level   = &model.AreaLevel{}
		adid    int64
		stime   int64
		etime   int64
		areas   []string
		tps     []int64
		source  int64
		keyType int64
	)
	params := c.Request.Form
	ruleStr := params.Get("rule")
	areaStr := params.Get("area")
	modeStr := params.Get("mode")
	levelJSON := params.Get("level") // level json
	commentStr := params.Get("comment")
	adidStr := params.Get("adid")
	nameStr := params.Get("name")
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	tpStr := params.Get("tpid")
	sourceStr := params.Get("source")
	typeStr := params.Get("key_type")
	if modeStr != "" {
		if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
			log.Error("strconv.ParseInt() err(%v)", err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	// Parse area level
	if err = json.Unmarshal([]byte(levelJSON), level); err != nil {
		log.Error("json.Unmarshal(%s) err(%+v)")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tpStr != "" {
		if tps, err = xstr.SplitInts(tpStr); err != nil {
			log.Error("xstr.SplintInts(tpStr %s) err(%v)", tpStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		tps = []int64{0}
	}
	if areaStr == "" {
		areas = append(areas, _paddingArea)
	} else {
		areas = strings.Split(areaStr, ",")
	}
	if ruleStr == "" {
		log.Error("ruleStr == nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var rules []string
	if mode == 1 {
		rules = strings.Split(ruleStr, "|")
	} else {
		rules = []string{ruleStr}
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adidStr, adid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if stimeStr != "" {
		stime, err = strconv.ParseInt(stimeStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", stimeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		stime = time.Now().Unix()
	}
	if etimeStr != "" {
		etime, err = strconv.ParseInt(etimeStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", etimeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		etime = stime + 3600*24*365*10 // 10年
	}
	if len([]rune(commentStr)) > 50 || len([]rune(commentStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if sourceStr != "" {
		source, err = strconv.ParseInt(sourceStr, 10, 64)
		if err != nil {
			log.Error("source strconv.ParseInt(%s) error(%v)", source, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if typeStr != "" {
		keyType, err = strconv.ParseInt(typeStr, 10, 64)
		if err != nil {
			log.Error("source strconv.ParseInt(%s) error(%v)", source, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(nil, svc.AdminAdd(c, areas, rules, level, commentStr, nameStr, int8(mode), tps, adid, stime, etime, int8(source), int8(keyType)))
}

func filterEdit(c *bm.Context) {
	var (
		err      error
		id, mode int64
		tps      []int64
		areas    []string
		source   int64
		keyType  int64
		level    = &model.AreaLevel{}
	)

	params := c.Request.Form
	ruleStr := params.Get("rule")
	areaStr := params.Get("area")
	modeStr := params.Get("mode")
	levelJSON := params.Get("level") // level json
	commentStr := params.Get("comment")
	adidStr := params.Get("adid")
	reasonStr := params.Get("reason")
	idStr := params.Get("id")
	nameStr := params.Get("name")
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	tpStr := params.Get("tpid")
	sourceStr := params.Get("source")
	typeStr := params.Get("key_type")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if modeStr != "" {
		if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
			log.Error("strconv.ParseInt() err(%v)", err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	// Parse area level
	if err = json.Unmarshal([]byte(levelJSON), level); err != nil {
		log.Error("json.Unmarshal(%s) err(%+v)")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tpStr != "" {
		if tps, err = xstr.SplitInts(tpStr); err != nil {
			log.Error("xstr.SplintInts(tpStr %s) err(%v)", tpStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		tps = []int64{0}
	}
	if areaStr == "" {
		areas = append(areas, _paddingArea)
	} else {
		areas = strings.Split(areaStr, ",")
	}
	if ruleStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	stime, err := strconv.ParseInt(stimeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	etime, err := strconv.ParseInt(etimeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(commentStr)) > 50 || len([]rune(commentStr)) < 1 || len([]rune(reasonStr)) > 50 || len([]rune(reasonStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if sourceStr != "" {
		source, err = strconv.ParseInt(sourceStr, 10, 64)
		if err != nil {
			log.Error("source strconv.ParseInt(%s) error(%v)", source, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if typeStr != "" {
		keyType, err = strconv.ParseInt(typeStr, 10, 64)
		if err != nil {
			log.Error("source strconv.ParseInt(%s) error(%v)", source, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(nil, svc.AdminEdit(c, areas, ruleStr, commentStr, reasonStr, nameStr, int8(mode), level, tps, id, adid, stime, etime, int8(source), int8(keyType)))
}

func filterDel(c *bm.Context) {

	params := c.Request.Form
	fidStr := params.Get("fid")
	adidStr := params.Get("adid")
	reasonStr := params.Get("reason")
	nameStr := params.Get("name")
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt() err(%v)", err)
		return
	}
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(reasonStr)) > 50 || len(reasonStr) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AdminDel(c, fid, adid, reasonStr, nameStr))
}

func filterList(c *bm.Context) {
	var (
		err    error
		pn, ps int64
	)

	params := c.Request.Form
	areaStr := params.Get("area")
	psStr := params.Get("ps")
	pnStr := params.Get("pn")

	if areaStr == "" {
		areaStr = _baseArea
	}
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 0 || ps > 200 {
		ps = 200
	}
	rules, total, err := svc.AdminSearch(c, "", areaStr, "", "", 0, model.FilterStateNormal, 0, pn, ps)
	searchRules := make([]*model.FilterForSearch, 0)
	for _, rule := range rules {
		var sf = &model.FilterForSearch{}
		sf.LoadFromFilter(rule)
		searchRules = append(searchRules, sf)
	}
	var data = map[string]interface{}{
		"rules": searchRules,
		"total": total,
	}
	c.JSON(data, err)
}

func filterSearch(c *bm.Context) {
	var (
		err      error
		params   = c.Request.Form
		msg      = params.Get("msg")
		area     = params.Get("area")
		source   = params.Get("source")
		levelStr = params.Get("level")
		ftype    = params.Get("filter_type")
		stageStr = params.Get("stage")
		pnStr    = params.Get("pn")
		psStr    = params.Get("ps")
	)
	var (
		pn      int64
		ps      int64
		state   int
		deleted int
		level   int
	)
	if area == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(msg) > 100 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 0 || ps > 200 {
		ps = 200
	}

	if levelStr == "" {
		level = 0
	} else {
		if level, err = strconv.Atoi(levelStr); err != nil {
			log.Error("%+v", errors.WithStack(err))
			return
		}
	}
	switch stageStr {
	case "", "0":
		state = _normal
		deleted = _normal
	case "1":
		state = _deleted
		deleted = _deleted
	case "2":
		state = _expired
		deleted = _normal
	}
	rules, total, err := svc.AdminSearch(c, msg, area, source, ftype, level, state, deleted, pn, ps)
	searchRules := make([]*model.FilterForSearch, 0)
	for _, rule := range rules {
		var sf = &model.FilterForSearch{}
		sf.LoadFromFilter(rule)
		searchRules = append(searchRules, sf)
	}
	var data = map[string]interface{}{
		"rules": searchRules,
		"total": total,
	}
	c.JSON(data, err)
}

func filterLog(c *bm.Context) {

	params := c.Request.Form
	idStr := params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.AdminLog(c, id))
}

func filterOrigin(c *bm.Context) {
	var (
		err error
		id  int64
	)

	params := c.Request.Form
	idStr := params.Get("id")
	areaStr := params.Get("area")
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.AdminOrigin(c, id, areaStr))
}

func filterOrigins(c *bm.Context) {
	var (
		err error
		ids []int64
	)

	params := c.Request.Form
	idsStr := params.Get("ids")
	areaStr := params.Get("area")
	if ids, err = xstr.SplitInts(idsStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.AdminOrigins(c, ids, areaStr))
}
