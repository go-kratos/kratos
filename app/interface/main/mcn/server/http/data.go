package http

import (
	"context"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/net/http/blademaster"
)

func mcnGetMcnGetIndexInc(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetIndexIncReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnGetIndexInc(context, arg.(*mcnmodel.McnGetIndexIncReq))
		},
		"GetMcnGetIndexInc",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetMcnGetIndexSource(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetIndexSourceReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnGetIndexSource(context, arg.(*mcnmodel.McnGetIndexSourceReq))
		},
		"GetMcnGetIndexSource",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetPlaySource(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetPlaySourceReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetPlaySource(context, arg.(*mcnmodel.McnGetPlaySourceReq))
		},
		"GetPlaySource",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetMcnFans(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFans(context, arg.(*mcnmodel.McnGetMcnFansReq))
		},
		"GetMcnFans",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetMcnFansInc(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansIncReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansInc(context, arg.(*mcnmodel.McnGetMcnFansIncReq))
		},
		"GetMcnFansInc",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetMcnFansDec(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansDecReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansDec(context, arg.(*mcnmodel.McnGetMcnFansDecReq))
		},
		"GetMcnFansDec",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetMcnFansAttentionWay(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansAttentionWayReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansAttentionWay(context, arg.(*mcnmodel.McnGetMcnFansAttentionWayReq))
		},
		"GetMcnFansAttentionWay",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetBaseFansAttrReq(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetBaseFansAttrReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetBaseFansAttrReq(context, arg.(*mcnmodel.McnGetBaseFansAttrReq))
		},
		"GetBaseFansAttrReq",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetFansArea(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansAreaReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansArea(context, arg.(*mcnmodel.McnGetFansAreaReq))
		},
		"GetFansArea",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetFansType(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansTypeReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansType(context, arg.(*mcnmodel.McnGetFansTypeReq))
		},
		"GetFansType",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetFansTag(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansTagReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansTag(context, arg.(*mcnmodel.McnGetFansTagReq))
		},
		"GetFansTag",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}
