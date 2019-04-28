package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	dmv1pb "go-common/app/interface/live/open-interface/api/http/v1"
	"go-common/app/interface/live/open-interface/internal/dao"
	broadcasrtService "go-common/app/service/live/broadcast-proxy/api/v1"
	titansSdk "go-common/app/service/live/resource/sdk"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

func checksign(ts string, sign string, group string) bool {
	sts := time.Now().Unix()
	cts, _ := strconv.ParseInt(ts, 10, 64)
	dv := sts - cts
	if math.Abs(float64(dv)) > 5 {
		log.Info("[dm] ts  err: %+d", dv)
		return false
	}

	userSecrets, terr := titansSdk.Get("dmUser")
	if terr != nil {
		log.Error("[dm] get titan conf err: %+v", terr)
		return false
	}
	juserSecrets := make(map[string]string)
	if jerr := json.Unmarshal([]byte(userSecrets), &juserSecrets); jerr != nil {
		log.Error("[dm] JSON decode titan conf err: %+v", terr)
		return false
	}

	secret, ok := juserSecrets[group]
	if !ok {
		log.Info("[dm] unknow  group: %+v", group)
		return false
	}

	newSign := fmt.Sprintf("%x", md5.Sum([]byte(group+secret+ts)))
	if newSign != sign {
		log.Info("[dm] check sign err sign: %+s service sign: %s", sign, newSign)
		return false
	}
	return true
}

//Sendmsg 发送弹幕消息
func (s *Service) Sendmsg(ctx context.Context, req *dmv1pb.SendMsgReq) (resp *dmv1pb.SendMsgResp, err error) {
	if ok := checksign(req.GetTs(), req.GetSign(), req.GetGroup()); !ok {
		return nil, ecode.Error(-403, "sign or ts error")
	}

	var dmString = "{\"cmd\":\"DANMU_MSG\",\"info\":[[0,1,25,16777215,%d,0,0,\"\",0,0,0],\"%s\",[0,\"\",0,0,0,0,0,\"\"],[],[0,0,0,\"\"],[\"\",\"\"],0,0]}"

	breq := &broadcasrtService.RoomMessageRequest{
		RoomId:  int32(req.GetRoomID()),
		Message: fmt.Sprintf(dmString, time.Now().Unix(), req.GetMsg()),
	}
	_, berr := dao.BcastClient.DanmakuClient.RoomMessage(ctx, breq)
	if berr != nil {
		log.Error("[dm] SendBroadCastGrpc err: %+v", berr)
		return nil, ecode.Error(-400, "send msg err")
	}
	return
}

//GetConf 获取弹幕配置
func (s *Service) GetConf(ctx context.Context, req *dmv1pb.GetConfReq) (resp *dmv1pb.GetConfResp, err error) {
	if ok := checksign(req.GetTs(), req.GetSign(), req.GetGroup()); !ok {
		return nil, ecode.Error(-403, "sign or ts error")
	}

	resp = &dmv1pb.GetConfResp{
		WSPort:     []int64{2244},
		WSSPort:    []int64{443},
		TCPPort:    []int64{2243, 80},
		DomianList: []string{},
	}

	breq := &broadcasrtService.DispatchRequest{
		UserIp: metadata.String(ctx, metadata.RemoteIP),
	}
	bresp, berr := dao.BcastClient.Dispatch(ctx, breq)
	if berr != nil {
		resp.IPList = []string{"broadcastlv.chat.bilibili.com"}
		resp.DomianList = []string{"broadcastlv.chat.bilibili.com"}
		log.Error("[dm] get IPList by BcastClient Dispatch err:%+v", berr)
	} else {
		resp.IPList = bresp.Ip
		resp.DomianList = bresp.Host
	}
	return resp, nil
}
