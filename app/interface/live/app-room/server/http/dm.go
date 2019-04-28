package http

import (
	apiv1 "go-common/app/interface/live/app-room/api/http/v1"
	v1index "go-common/app/interface/live/app-room/service/v1/dm"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func sendMsgSendMsg(ctx *bm.Context) {
	p := new(apiv1.SendDMReq)
	if err := ctx.Bind(p); err != nil {
		return
	}
	_, err := dmservice.SendMsg(ctx, p)
	res := map[string]interface{}{}
	if e, ok := err.(*ecode.Status); ok {
		res["msg"] = e.Message()
		res["message"] = e.Message()
		//验证返回
		if e.Code() == 1990000 {
			res["data"] = map[string]string{
				"verify_url": "https://live.bilibili.com/p/html/live-app-captcha/index.html?is_live_half_webview=1&hybrid_half_ui=1,5,290,332,0,0,30,0;2,5,290,332,0,0,30,0;3,5,290,332,0,0,30,0;4,5,290,332,0,0,30,0;5,5,290,332,0,0,30,0;6,5,290,332,0,0,30,0;7,5,290,332,0,0,30,0;8,5,290,332,0,0,30,0",
			}
		}
		ctx.JSONMap(res, err)
		return
	}
	res["msg"] = ""
	res["message"] = ""
	res["data"] = []string{}
	ctx.JSONMap(res, err)
}

func getHistory(ctx *bm.Context) {
	p := new(apiv1.HistoryReq)
	if err := ctx.Bind(p); err != nil {
		return
	}

	resp, err := dmservice.GetHistory(ctx, p)

	res := map[string]interface{}{}
	res["msg"] = ""
	res["message"] = ""

	empty := make(map[string][]string)
	empty["room"] = make([]string, 0)
	empty["admin"] = make([]string, 0)
	if err != nil {
		res["data"] = empty
		ctx.JSONMap(res, err)
		return
	}
	res["data"] = v1index.HistoryData(resp)
	ctx.JSONMap(res, err)
}
