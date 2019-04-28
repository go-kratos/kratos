package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/live/app-blink/conf"
	"go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

//FILTER_MPOST_URL 主站屏蔽词url
const FILTER_MPOST_URL = "http://api.bilibili.co/x/internal/filter/mpost"

//CheckMsgIsLegal ...
//area值:"live_danmu";"live_biaoti";"live_tag";"live_jianjie",;"live_xuyuanping"
func (d *Dao) CheckMsgIsLegal(c context.Context, msg string, area string, roomId int64) (resp bool, err error) {
	if msg == "" || area == "" {
		log.Error("checkMsgIsLegal_params_error:msg=%s;area=%s", msg, area)
		return
	}
	uid, _ := metadata.Value(c, metadata.Mid).(int64)
	filterReply, err := d.FilterSite(c, msg, area)
	if err != nil {
		err = ecode.CallMainFilterError
		return
	}
	areaMap := map[string]int{"live_biaoti": 2, "live_tag": 3, "live_jianjie": 4, "live_xuyuanping": 5}
	historyData := &v1.RoomMngSaveHistoryReq_List{}
	historyDatas := &v1.RoomMngSaveHistoryReq{}
	resp = false
	for _, filterInfo := range filterReply {
		level := filterInfo.Level
		hit := filterInfo.Hit
		msg := filterInfo.Msg
		if level >= 20 {
			resp = true
		}
		if roomId > 0 && level == 15 && hit != nil && areaMap[area] > 0 {
			historyData.Area = int64(areaMap[area])
			historyData.Uid = uid
			historyData.Roomid = roomId
			historyData.Content = msg
			historyData.Status = 0
			historyData.Oname = ""
		}
		historyDatas.List = append(historyDatas.List, historyData)
	}
	//保存命中审核的数据
	d.RoomApi.V1RoomMng.SaveHistory(c, historyDatas)
	return

}

type mpostData struct {
	Level int                    `json:"level"`
	Limit int                    `json:"limit"`
	Hit   map[string]interface{} `json:"hit"`
	Msg   string                 `json:"msg"`
}

//FilterSite 主站12月初grpc接口上线，可以替换
func (d *Dao) FilterSite(c context.Context, msg string, area string) (resp []*mpostData, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("appkey", conf.APPKey)
	params.Set("msg", msg)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("area", area)
	var res struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data []*mpostData `json:"data"`
	}
	if d.HttpCli.Post(c, FILTER_MPOST_URL, ip, params, &res); err != nil {
		log.Error("filterSite:%v;params:%v;code:%d;msg:%s", err, params, res.Code, res.Msg)
		return
	}
	if res.Code != 0 {
		log.Error("filterSite:%v;code:%d;msg:%s", err, res.Code, res.Msg)
	}
	//接入sla
	resp = res.Data
	return
}
