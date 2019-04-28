package http

import (
	"encoding/json"
	"strconv"
	"strings"

	"go-common/app/interface/main/dm/model"
	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

var (
	_userValidTime = 86400
	_glbValidTime  = 86400
)

func userRules(c *bm.Context) {
	mid, _ := c.Get("mid")
	rs, err := dmSvc.UserRules(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"rule":  rs,
		"ver":   123,
		"valid": _userValidTime,
	}
	c.JSON(data, err)
}

func globalRuleEmpty(c *bm.Context) {
	var (
		param      = c.Request.Form
		curVersion = dmSvc.GlobalRuleVersion()
	)
	_, err := strconv.ParseUint(param.Get("ver"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := map[string]interface{}{
		"rule":  make([]*dm2Mdl.GlobalFilter, 0),
		"ver":   curVersion,
		"valid": _glbValidTime,
	}
	c.JSON(data, err)
}

// func globalRules(c context.Context) {
// 	var (
// 		res        = c.Result()
// 		param      = c.Request().Form
// 		curVersion = dmSvc.GlobalRuleVersion()
// 	)
// 	version, err := strconv.ParseUint(param.Get("ver"), 10, 64)
// 	if err != nil {
// 		res["code"] = ecode.RequestErr
// 		return
// 	}
// 	if version == curVersion {
// 		res["code"] = ecode.NotModified
// 		return
// 	}
// 	rs, err := dmSvc.GlobalRules(c)
// 	if err != nil {
// 		res["code"] = err
// 		return
// 	}
// 	data := map[string]interface{}{
// 		"rule":  rs,
// 		"ver":   curVersion,
// 		"valid": _glbValidTime,
// 	}
// 	res["data"] = data
// }

func addUserRule(c *bm.Context) {
	params := c.Request.Form
	typeStr := params.Get("type")
	filter := params.Get("filter")
	comment := params.Get("comment")
	mid, _ := c.Get("mid")
	tp, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if filter == "" {
		c.JSON(nil, ecode.DMFilterIsEmpty)
		return
	}
	v, err := dmSvc.AddUserRule(c, int8(tp), mid.(int64), []string{filter}, comment)
	if err != nil || len(v) <= 0 {
		c.JSON(nil, err)
		return
	}
	c.JSON(v[0], err)
}

func multiAddUserRule(c *bm.Context) {
	params := c.Request.Form
	typeStr := params.Get("type")
	filters := params.Get("filters")
	comment := params.Get("comment")
	mid, _ := c.Get("mid")
	tp, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if filters == "" {
		c.JSON(nil, ecode.DMFilterIsEmpty)
		return
	}
	values := strings.Split(filters, ",")
	v, err := dmSvc.AddUserRule(c, int8(tp), mid.(int64), values, comment)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(v, err)
}

func addGlobalRule(c *bm.Context) {
	params := c.Request.Form
	filter := params.Get("filter")
	typeStr := params.Get("type")
	tp, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if filter == "" {
		c.JSON(nil, ecode.DMFilterIsEmpty)
		return
	}
	v, err := dmSvc.AddGlobalRule(c, int8(tp), filter)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(v, err)
}

func delUserRules(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	idsStr := params.Get("ids")
	if len(idsStr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ids, err := xstr.SplitInts(idsStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, err = dmSvc.DelUserRules(c, mid.(int64), ids)
	c.JSON(nil, err)
}

func delGlobalRules(c *bm.Context) {
	var params = c.Request.Form
	idsStr := params.Get("ids")
	if len(idsStr) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ids, err := xstr.SplitInts(idsStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, err = dmSvc.DelGlobalRules(c, ids)
	c.JSON(nil, err)
}

// filterList
func filterList(c *bm.Context) {
	var (
		cid int64
		err error
		p   = c.Request.Form
	)
	cidStr := p.Get("cid")
	if len(cidStr) > 0 {
		if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	mid, _ := c.Get("mid")
	data, err := dmSvc.FilterList(c, mid.(int64), cid)
	if err != nil {
		log.Error("dmSvc.FilterList(mid:%v cid:%d) error(%v)", mid, cid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func editFilter(c *bm.Context) {
	var (
		cid   int64
		fList = make([]*model.IndexFilter, 0)
		p     = c.Request.Form
		err   error
	)
	mid, _ := c.Get("mid")
	cidStr := p.Get("cid")
	if cidStr != "" {
		if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	fStr := p.Get("filters")
	if err = json.Unmarshal([]byte(fStr), &fList); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, f := range fList {
		if f.Regex == dm2Mdl.FilterTypeRev ||
			f.Regex == dm2Mdl.FilterTypeTop ||
			f.Regex == dm2Mdl.FilterTypeBottom {
			f.Filter = dm2Mdl.FilterContent
		}
		if err = dmSvc.EditFilter(c, cid, mid.(int64), f.Filter, f.Regex, f.Activate); err != nil {
			log.Error("dmSvc.EditFilter(mid:%v cid:%d filter:%s type:%d) error(%v)", mid, cid, f.Filter, f.Regex, err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
}
