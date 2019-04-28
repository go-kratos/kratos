package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

	"go-common/app/service/live/live-dm/dao"
	"go-common/app/service/live/live-dm/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	uuid "github.com/satori/go.uuid"
)

//History 存档信息
type History struct {
	Text       string        `json:"text"`
	UID        int64         `json:"uid"`
	NickName   string        `json:"nickname"`
	UnameColor string        `json:"uname_color"`
	TimeLine   string        `json:"timeline"`
	Isadmin    int32         `json:"isadmin"`
	Vip        int           `json:"vip"`
	SVip       int           `json:"svip"`
	Medal      []interface{} `json:"medal"`
	Title      []interface{} `json:"title"` //内容格式待定
	UserLevel  []interface{} `json:"user_level"`
	Rank       int32         `json:"rank"`
	Teamid     int64         `json:"teamid"`
	RND        string        `json:"rnd"`
	UserTitle  string        `json:"user_title"`
	GuardLevel int           `json:"guard_level"`
	Bubble     int64         `json:"bubble"`
}

//DatatoString 数据处理转换为json
func (h *History) DatatoString(s *SendMsg) string {
	h.Text = s.SendMsgReq.Msg
	h.NickName = s.UserBindInfo.Uname
	h.UID = s.SendMsgReq.Uid
	h.UnameColor = s.UserInfo.UnameColor
	h.TimeLine = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	if s.UserInfo.RoomAdmin {
		h.Isadmin = 1
	} else {
		h.Isadmin = 0
	}
	h.Vip = s.UserInfo.Vip
	h.SVip = s.UserInfo.Svip
	h.Rank = s.UserBindInfo.URank
	h.Teamid = 0
	h.RND = s.SendMsgReq.Rnd
	h.UserTitle = s.TitleConf.Title
	h.GuardLevel = s.UserInfo.PrivilegeType
	h.Bubble = s.UserInfo.Bubble

	md := make([]interface{}, 0, 6)
	if s.UserInfo.MedalInfo.MedalName != "" {
		md = append(md, s.UserInfo.MedalInfo.MedalLevel)
		md = append(md, s.UserInfo.MedalInfo.MedalName)
		md = append(md, s.UserInfo.MedalInfo.RUName) //TODO 获取主播get_weared_medal
		md = append(md, s.UserInfo.MedalInfo.RoomID)
		md = append(md, s.UserInfo.MedalInfo.MColor)
		md = append(md, s.UserInfo.MedalInfo.SpecialMedal)
	}
	h.Medal = md

	ul := make([]interface{}, 0, 4)
	ul = append(ul, s.UserInfo.UserLever)
	ul = append(ul, 0)
	ul = append(ul, s.UserInfo.ULevelColor)
	if s.UserInfo.ULevelRank > 50000 {
		ul = append(ul, ">50000")
	} else {
		ul = append(ul, s.UserInfo.ULevelRank)
	}
	h.UserLevel = ul

	tl := make([]interface{}, 0, 2)
	tl = append(tl, s.TitleConf.OldTitle)
	tl = append(tl, s.TitleConf.Title)
	h.Title = tl

	msg, _ := json.Marshal(h)
	return string(msg)
}

//BroadCastMsg 广播信息
type BroadCastMsg struct {
	CMD  string        `json:"cmd"`
	Info []interface{} `json:"info"`
}

//DatatoString 数据处理转换为json
func (b *BroadCastMsg) DatatoString(s *SendMsg) string {
	b.CMD = "DANMU_MSG"
	b.Info = make([]interface{}, 0, 10)

	//弹幕配置
	dc := make([]interface{}, 0, 11)
	dc = append(dc, 0)
	dc = append(dc, s.SendMsgReq.Mode)
	dc = append(dc, s.SendMsgReq.Fontsize)
	dc = append(dc, s.DMconf.Color)
	dc = append(dc, time.Now().Unix())
	var rand int64
	rand, _ = strconv.ParseInt(s.SendMsgReq.Rnd, 10, 64)
	dc = append(dc, rand)
	dc = append(dc, 0)
	dc = append(dc, fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(strconv.FormatInt(s.SendMsgReq.Uid, 10)))))
	dc = append(dc, 0)
	dc = append(dc, s.SendMsgReq.Msgtype)
	dc = append(dc, s.UserInfo.Bubble)

	//用户信息
	userInfo := make([]interface{}, 0, 8)
	userInfo = append(userInfo, s.SendMsgReq.Uid)
	userInfo = append(userInfo, s.UserBindInfo.Uname)
	var admin int
	if s.UserInfo.RoomAdmin {
		admin = 1
	} else {
		admin = 0
	}
	userInfo = append(userInfo, admin)
	userInfo = append(userInfo, s.UserInfo.Vip)
	userInfo = append(userInfo, s.UserInfo.Svip)
	userInfo = append(userInfo, s.UserBindInfo.URank)
	userInfo = append(userInfo, s.UserBindInfo.MobileVerify)
	userInfo = append(userInfo, s.UserInfo.UnameColor)

	//勋章配置
	medalInfo := make([]interface{}, 0, 6)
	if s.UserInfo.MedalInfo.MedalName != "" {
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.MedalLevel)
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.MedalName)
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.RUName) // 主播名称
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.RoomID)
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.MColor)
		medalInfo = append(medalInfo, s.UserInfo.MedalInfo.SpecialMedal)
	}

	//用户等级信息
	ul := make([]interface{}, 0, 4)
	ul = append(ul, s.UserInfo.UserLever)
	ul = append(ul, 0)
	ul = append(ul, s.UserInfo.ULevelColor)
	if s.UserInfo.ULevelRank > 50000 {
		ul = append(ul, ">50000")
	} else {
		ul = append(ul, s.UserInfo.ULevelRank)
	}

	//头衔
	tl := make([]interface{}, 0, 2)
	tl = append(tl, s.TitleConf.OldTitle)
	tl = append(tl, s.TitleConf.Title)

	//组合
	b.Info = append(b.Info, dc)
	b.Info = append(b.Info, s.SendMsgReq.Msg)
	b.Info = append(b.Info, userInfo)
	b.Info = append(b.Info, medalInfo)
	b.Info = append(b.Info, ul)
	b.Info = append(b.Info, tl)
	b.Info = append(b.Info, 0)
	b.Info = append(b.Info, s.UserInfo.PrivilegeType)

	b.Info = append(b.Info, nil)

	msg, _ := json.Marshal(b)
	return string(msg)
}

func send(ctx context.Context, sdm *SendMsg) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		saveHistory(ctx, sdm)
		return nil
	})
	g.Go(func() error {
		incrDMNum(ctx, sdm)
		return nil
	})
	g.Go(func() error {
		return sendBroadCast(ctx, sdm)
	})
	//databus投递
	err := sdm.Dmservice.dao.Databus.Do(context.Background(), func(c context.Context) {
		sendDataBus(context.Background(), sdm)
	})
	if err != nil {
		log.Error("DM: send databus save err: %+v", err)
	}
	return g.Wait()
}

func saveHistory(ctx context.Context, sdm *SendMsg) {
	hm := &History{}
	jhm := hm.DatatoString(sdm)
	sdm.Dmservice.dao.SaveHistory(ctx, jhm, sdm.UserInfo.RoomAdmin, sdm.SendMsgReq.Roomid)
}

func incrDMNum(ctx context.Context, sdm *SendMsg) {
	dao.IncrDMNum(ctx, sdm.SendMsgReq.Roomid, sdm.SendMsgReq.Msgtype)
}

func sendBroadCast(ctx context.Context, sdm *SendMsg) error {
	bm := &BroadCastMsg{}
	jbm := bm.DatatoString(sdm)
	if err := dao.SendBroadCastGrpc(ctx, jbm, sdm.SendMsgReq.Roomid); err != nil {
		// if err := dao.SendBroadCast(ctx, jbm, sdm.SendMsgReq.Roomid); err != nil {
		lancer(sdm, "弹幕投递消息失败")
		return nil
		// }
		// return nil
	}
	return nil
}

func sendDataBus(ctx context.Context, sdm *SendMsg) {
	info := &model.BNDatabus{
		Roomid:    sdm.SendMsgReq.Roomid,
		UID:       sdm.SendMsgReq.Uid,
		Uname:     sdm.UserBindInfo.Uname,
		UserLever: sdm.UserInfo.UserLever,
		Color:     sdm.DMconf.Color,
		Msg:       sdm.SendMsgReq.Msg,
		Time:      time.Now().Unix(),
		MsgType:   sdm.SendMsgReq.Msgtype,
		MsgID:     uuid.NewV4().String(),
	}
	dao.SendBNDatabus(ctx, sdm.SendMsgReq.Uid, info)
}
