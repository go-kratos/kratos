package http

import (
	"encoding/json"
	"errors"
	"strings"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func checkParams(c *bm.Context) (sp *model.QueryParams, err error) {
	sp = &model.QueryParams{
		QueryBody: &model.QueryBody{},
	}
	if err = c.Bind(sp); err != nil {
		return
	}
	decoder := json.NewDecoder(strings.NewReader(sp.QueryBodyStr))
	decoder.UseNumber()
	if err = decoder.Decode(&sp.QueryBody); err != nil {
		log.Error("s.http.upsert(%v) json error(%v)", err, sp.QueryBodyStr)
	}
	if sp.QueryBody == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// eq
	if sp.QueryBody.Where != nil {

		for k, v := range sp.QueryBody.Where.EQ {
			if integer, ok := v.(json.Number); ok {
				if sp.QueryBody.Where.EQ[k], err = integer.Int64(); err != nil {
					log.Error("sp.QueryBody.Where.EQ.numberToInt64(%v)(%v)", integer, err)
				}
			}
		}
		// or
		for k, v := range sp.QueryBody.Where.Or {
			if integer, ok := v.(json.Number); ok {
				if sp.QueryBody.Where.Or[k], err = integer.Int64(); err != nil {
					log.Error("sp.QueryBody.Where.EQ.numberToInt64(%v)(%v)", integer, err)
				}
			}
		}
		// in
		for k, v := range sp.QueryBody.Where.In {
			vs := make([]interface{}, 0)
			for _, v2 := range v {
				if integer, ok := v2.(json.Number); ok {
					vs = append(vs, integer)
				} else {
					vs = append(vs, v2)
				}
			}
			sp.QueryBody.Where.In[k] = vs
		}
	}

	if sp.QueryBody.From == "" {
		err = errors.New("query from is empty")
		return
	}
	// deny *
	if b, ok := model.PermConf["star"][sp.Business]; !(ok && b == "true") {
		if strings.Contains(sp.QueryBody.From, "*") {
			err = errors.New("query from * is denied")
			return
		}
	}
	if err = svr.CheckQueryConf(c, sp); err != nil {
		log.Error("svr.QueryConf(%+v) error(%v)", sp, err)
		return
	}
	if sp.QueryBody.Pn == 0 {
		sp.QueryBody.Pn = 1
	}
	if sp.QueryBody.Ps == 0 {
		sp.QueryBody.Ps = 10
	}
	return
}

func querySearch(c *bm.Context) {
	var (
		err error
		sp  *model.QueryParams
		res *model.QueryResult
	)
	res = &model.QueryResult{}
	if sp, err = checkParams(c); err != nil {
		log.Error("checkParams(%v) error(%v)", sp, err)
		c.JSON(res, err)
		return
	}
	//query mode
	switch sp.AppIDConf.QueryMode {
	case model.QueryModeBasic: //完全通用查询，支持nested方式
		res, _, err = svr.QueryBasic(c, sp)
	case model.QueryModeExtra: //半自定义，在查询体基础上加额外条件
		res, _, err = svr.QueryExtra(c, sp)
	case 4: //完整自定义查询
		// todo
	default:
		res, _, err = svr.QueryBasic(c, sp)
	}
	if err != nil {
		log.Error("srv.QueryBasic(%v) error(%v)", sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

// queryDebug 1. queryBody is right or not 2. return dsl body 3. return dsl profile
func queryDebug(c *bm.Context) {
	var (
		err      error
		sp       *model.QueryParams
		resDebug = &model.QueryDebugResult{}
	)
	if sp, err = checkParams(c); err != nil {
		log.Error("checkParams(%v) error(%v)", c, err)
		resDebug.AddErrMsg(err.Error())
		resDebug.QueryBody = sp.QueryBodyStr
		c.JSON(resDebug, ecode.RequestErr)
		return
	}
	// 默认高级debug
	if sp.DebugLevel == 0 {
		sp.DebugLevel = 2
	}
	switch sp.AppIDConf.QueryMode {
	case 0, model.QueryModeBasic: //完全通用查询
		if _, resDebug, err = svr.QueryBasic(c, sp); err != nil {
			log.Error("QueryDebug(%v) error(%v)", sp, err)
			resDebug.AddErrMsg(err.Error())
			c.JSON(resDebug, ecode.RequestErr)
			return
		}
	case model.QueryModeExtra: //半自定义，根据查询体基础再做额外处理
		if _, resDebug, err = svr.QueryExtra(c, sp); err != nil {
			log.Error("QueryDebug(%v) error(%v)", sp, err)
			resDebug.AddErrMsg(err.Error())
			c.JSON(resDebug, ecode.RequestErr)
			return
		}
	case 3: //完全自定义
		//todo
	default:

	}
	resDebug.QueryBody = sp.QueryBodyStr
	c.JSON(resDebug, err)
}
