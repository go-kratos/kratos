package http

import (
	"go-common/app/service/live/userexp/conf"
	"go-common/app/service/live/userexp/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"strconv"
)

func level(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(expSvr.Level(c, uid))
}

func multiGetLevel(c *bm.Context) {
	uidsStr := c.Request.Form.Get("uids")
	// check params
	uids, err := xstr.SplitInts(uidsStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	levels, err := expSvr.MultiGetLevel(c, uids)
	if err != nil {
		log.Error("[http.exp|multiGetLevel] expSvr.MultiGetLevel(%v) error(%v)", uids, err)
		c.JSON(nil, err)
		return
	}
	levelInfo := make(map[string]*model.Level, len(levels))
	for _, v := range levels {
		levelInfo[strconv.FormatInt(v.Uid, 10)] = v
	}
	c.JSON(levelInfo, nil)
}

func addUexp(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	uexpStr := c.Request.Form.Get("uexp")
	var exp *model.Exp
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uexp, err := strconv.ParseInt(uexpStr, 10, 64)
	if err != nil || uexp <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	ric := infocArg(c)
	err = expSvr.AddUexp(c, uid, uexp, ric)
	if err != nil {
		log.Error("[http.exp|addUexp] expSvr.AddUexp1(%u) error(%v)", uid, err)
	}
	LogSwitch := conf.Conf.Switch
	if LogSwitch != nil {
		var QueryConfig = conf.Conf.Switch.QueryExp
		if QueryConfig == 1 {
			exp, err = expSvr.Exp(c, uid)
			log.Info("addUexpUpdate uid:%d,Uexp:%d,Uexp:%d,delta:%d", exp.Uid, exp.Uexp, exp.Rexp, uexp)
		}
	}
	err = expSvr.AddUExpLog(c, uid, uexp, exp.Uexp, exp.Rexp, ric)

	c.JSON(nil, err)
}

func addRexp(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	rexpStr := c.Request.Form.Get("rexp")
	var exp *model.Exp
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rexp, err := strconv.ParseInt(rexpStr, 10, 64)
	if err != nil || rexp <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	ric := infocArg(c)
	err = expSvr.AddRexp(c, uid, rexp, ric)
	if err != nil {
		log.Error("[http.exp|addRexp] expSvr.AddRexp(%u) error(%v)", uid, err)
	}
	LogSwitch := conf.Conf.Switch
	if LogSwitch != nil {
		var QueryConfig = conf.Conf.Switch.QueryExp
		if QueryConfig == 1 {
			exp, err = expSvr.Exp(c, uid)
			log.Info("addRexpUpdate   uid:%d,Uexp:%d,Rexp:%d,delta:%d", exp.Uid, exp.Uexp, exp.Rexp, rexp)
		}
	}
	err = expSvr.AddRExpLog(c, uid, rexp, exp.Uexp, exp.Rexp, ric)

	c.JSON(nil, err)
}
