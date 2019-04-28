package http

import (
	"context"

	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/net/http/blademaster"
)

func mcnState(c *blademaster.Context) {
	oarg := new(mcnmodel.GetStateReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnGetState(context, arg.(*mcnmodel.GetStateReq))
		},
		"mcnState",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnExist(c *blademaster.Context) {
	oarg := new(mcnmodel.GetStateReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnExist(context, arg.(*mcnmodel.GetStateReq))
		},
		"mcnExist",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnBaseInfo(c *blademaster.Context) {
	oarg := new(mcnmodel.GetStateReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnBaseInfo(context, arg.(*mcnmodel.GetStateReq))
		},
		"mcnBaseInfo",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnApply(c *blademaster.Context) {
	oarg := new(mcnmodel.McnApplyReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnApply(context, arg.(*mcnmodel.McnApplyReq))
		},
		"mcnApply",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnBindUpApply(c *blademaster.Context) {
	oarg := new(mcnmodel.McnBindUpApplyReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnBindUpApply(context, arg.(*mcnmodel.McnBindUpApplyReq))
		},
		"mcnBindUpApply",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnUpConfirm(c *blademaster.Context) {
	oarg := new(mcnmodel.McnUpConfirmReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnUpConfirm(context, arg.(*mcnmodel.McnUpConfirmReq))
		},
		"mcnBindUpApply",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnUpGetBind(c *blademaster.Context) {
	oarg := new(mcnmodel.McnUpGetBindReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnUpGetBind(context, arg.(*mcnmodel.McnUpGetBindReq))
		},
		"mcnUpGetBind",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetDataSummary(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetDataSummaryReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnDataSummary(context, arg.(*mcnmodel.McnGetDataSummaryReq))
		},
		"mcnGetDataSummary",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetDataUpList(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetUpListReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnDataUpList(context, arg.(*mcnmodel.McnGetUpListReq))
		},
		"mcnGetDataUpList",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetAccountInfo(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetAccountReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetUpAccountInfo(context, arg.(*mcnmodel.McnGetAccountReq))
		},
		"mcnGetAccountInfo",
		nil,
		nil,
	)(c)
}

func mcnGetOldInfo(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnOldInfoReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnGetOldInfo(context, arg.(*mcnmodel.McnGetMcnOldInfoReq))
		},
		"McnGetOldInfo",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetRankUpFans(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetRankReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnGetRankUpFans(context, arg.(*mcnmodel.McnGetRankReq))
		},
		"McnGetRankUpFans",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetRankArchiveLikesOuter(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetRankReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnGetRankArchiveLikes(context, arg.(*mcnmodel.McnGetRankReq))
		},
		"mcnGetRankArchiveLikesOuter",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetRecommendPool(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetRecommendPoolReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetRecommendPool(context, arg.(*mcnmodel.McnGetRecommendPoolReq))
		},
		"GetRecommendPool",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetRecommendPoolTidList(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetRecommendPoolTidListReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetRecommendPoolTidList(context, arg.(*mcnmodel.McnGetRecommendPoolTidListReq))
		},
		"GetRecommendPoolTidList",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnGetChangePermit(c *blademaster.Context) {
	oarg := new(mcnmodel.McnChangePermitReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnChangePermit(context, arg.(*mcnmodel.McnChangePermitReq))
		},
		"McnChangePermit",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnPermitApplyGetBind(c *blademaster.Context) {
	oarg := new(mcnmodel.McnUpGetBindReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnPermitApplyGetBind(context, arg.(*mcnmodel.McnUpGetBindReq))
		},
		"McnPermitApplyGetBind",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnUpPermitApplyConfirm(c *blademaster.Context) {
	oarg := new(mcnmodel.McnUpConfirmReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnUpPermitApplyConfirm(context, arg.(*mcnmodel.McnUpConfirmReq))
		},
		"McnUpPermitApplyConfirm",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

func mcnPublicationPriceChange(c *blademaster.Context) {
	oarg := new(mcnmodel.McnPublicationPriceChangeReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnPublicationPriceChange(context, arg.(*mcnmodel.McnPublicationPriceChangeReq))
		},
		"McnPublicationPriceChange",
		[]preBindFuncType{getCookieMid(c, oarg)},
		[]preHandleFuncType{cheatReq},
	)(c)
}

// -------------- command ------------------
func cmdReloadRank(c *blademaster.Context) {
	oarg := new(mcnmodel.CmdReloadRank)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.CmdReloadRankCache(context, arg.(*mcnmodel.CmdReloadRank))
		},
		"CmdReloadRank",
		nil,
		nil,
	)(c)
}
