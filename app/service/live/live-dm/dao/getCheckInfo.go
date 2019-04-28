package dao

import (
	"context"
	"strconv"
	"sync"

	fansMedalService "go-common/app/service/live/fans_medal/api/liverpc/v2"
	"go-common/app/service/live/live-dm/conf"
	liveUserService "go-common/app/service/live/live_user/api/liverpc/v1"
	roomService "go-common/app/service/live/room/api/liverpc/v2"
	userService "go-common/app/service/live/user/api/liverpc/v3"
	userextService "go-common/app/service/live/userext/api/liverpc/v1"
	xuserService "go-common/app/service/live/xuser/api/grpc/v1"
	acctountService "go-common/app/service/main/account/api"
	filterService "go-common/app/service/main/filter/api/grpc/v1"
	spyService "go-common/app/service/main/spy/api"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//UserInfo 用户等级,经验等信息
type UserInfo struct {
	Vip           int
	Svip          int
	UserLever     int64
	UserScore     int64
	ULevelRank    int64
	ULevelColor   int64
	UnameColor    string
	RoomAdmin     bool
	PrivilegeType int
	Bubble        int64
	MedalInfo     *FansMedalInfo
	lock          sync.Mutex
}

//UserScore 用户真实分已经弹幕ai分
type UserScore struct {
	UserScore int64
	MsgAI     int64
	MsgLevel  int64
	lock      sync.Mutex
}

//DMConf 弹幕配置
type DMConf struct {
	Mode   int64
	Color  int64
	Length int64
}

//FansMedalInfo 粉丝勋章信息
type FansMedalInfo struct {
	MedalID      int64
	RUID         int64
	RUName       string
	MedalLevel   int64
	MedalName    string
	AnchorName   string
	RoomID       int64
	MColor       int64
	SpecialMedal string
	lock         sync.Mutex
}

//RoomConf 播主房间配置信息
type RoomConf struct {
	UID        int64
	RoomShield int64
	RoomID     int64
	Anchor     string
}

//UserBindInfo 用户手机,实名等信息
type UserBindInfo struct {
	Identification int32
	MobileVerify   int32
	Uname          string
	URank          int32
	// MobileVirtual  int64 无虚拟号段信息
}

//Get 获取用户的弹幕配置信息
func (d *DMConf) Get(ctx context.Context, uid int64, roomid int64, c *conf.Config) error {
	req := &userextService.DanmuConfGetAllReq{
		Uid:    uid,
		Roomid: roomid,
	}
	resp, err := UserExtServiceClient.V1DanmuConf.GetAll(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: userextService GetAll err: %v", err)
		}
		d.Color = 16777215
		d.Mode = 1
		d.Length = 20
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: userextService GetAll err code: %d", resp.Code)
		d.Color = 16777215
		d.Mode = 1
		d.Length = 20
		return nil
	}
	d.Color = resp.Data.Color
	d.Mode = resp.Data.Mode
	d.Length = resp.Data.Length
	su := strconv.FormatInt(uid, 10)
	if c.DmRules.Nixiang[su] {
		d.Mode = 6
	}
	if c.DmRules.Color[su] != 0 {
		d.Color = c.DmRules.Color[su]
	}
	return nil
}

// Get  获取用户等级,经验等信息
func (u *UserInfo) Get(ctx context.Context, uid int64) error {
	req := &userService.UserGetUserLevelInfoReq{
		Uid: uid,
	}
	resp, err := userClient.V3User.GetUserLevelInfo(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: user err: %v", err)
		}
		u.lock.Lock()
		u.UserLever = 0
		u.UserScore = 0
		u.ULevelColor = 16752445
		u.lock.Unlock()
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: user GetUserLevelInfo err code: %d", resp.Code)
		u.lock.Lock()
		u.UserLever = 0
		u.UserScore = 0
		u.ULevelColor = 16752445
		u.lock.Unlock()
		return nil
	}
	u.lock.Lock()
	u.UserLever = resp.Data.Level
	u.UserScore = resp.Data.Exp
	u.ULevelColor = resp.Data.Color
	u.lock.Unlock()
	return nil
}

//GetVipInfo 获取用户的老爷等级
func (u *UserInfo) GetVipInfo(ctx context.Context, uid int64) error {
	req := &xuserService.UidReq{
		Uid: uid,
	}
	resp, err := vipCli.Info(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: xuserService get vip Info err: %v", err)
		}
		u.lock.Lock()
		u.Vip = 0
		u.Svip = 0
		u.lock.Unlock()
		return nil
	}
	u.lock.Lock()
	u.Vip = resp.Info.Vip
	u.Svip = resp.Info.Svip
	u.lock.Unlock()
	return nil
}

//GetPrivilegeType 获取大航海信息
func (u *UserInfo) GetPrivilegeType(ctx context.Context, uid int64, ruid int64) error {
	req := &liveUserService.GuardGetByUidTargetIdReq{
		Uid:        uid,
		TargetId:   ruid,
		IsLimitOne: 0,
	}
	resp, err := LiveUserServiceClient.V1Guard.GetByUidTargetId(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: liveUserService GetByUidTargetId err: %v", err)
		}
		u.lock.Lock()
		u.PrivilegeType = 0
		u.lock.Unlock()
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: liveUserService GetByUidTargetId err code: %d", resp.Code)
		u.lock.Lock()
		u.PrivilegeType = 0
		u.lock.Unlock()
		return nil
	}
	var PrivilegeType int64
	for _, val := range resp.Data {
		i := val.PrivilegeType
		if PrivilegeType == 0 {
			PrivilegeType = i
		}
		if PrivilegeType != 0 && PrivilegeType > i {
			PrivilegeType = i
		}
	}
	u.lock.Lock()
	u.PrivilegeType = int(PrivilegeType)
	u.lock.Unlock()
	return nil
}

//IsRoomAdmin 判断用户是否是房管
func (u *UserInfo) IsRoomAdmin(ctx context.Context, uid int64, roomid int64) error {
	defer u.lock.Unlock()
	if roomid == 59125 || roomid == 5440 {
		u.lock.Lock()
		u.RoomAdmin = false
		return nil
	}

	req := &xuserService.RoomAdminIsAdminShortReq{
		Uid:    uid,
		Roomid: roomid,
	}
	resp, err := isAdmin.IsAdminShort(ctx, req)
	if err != nil {
		log.Error("DM:  xuser IsAdminShort err: %+v", err)
		u.lock.Lock()
		u.RoomAdmin = false
		return nil
	}
	if resp.Result == 1 {
		u.lock.Lock()
		u.RoomAdmin = true
		return nil
	}

	u.lock.Lock()
	u.RoomAdmin = false
	return nil
}

//GetFansMedal 获取主播的粉丝勋章信息
func (f *FansMedalInfo) GetFansMedal(ctx context.Context, uid int64) error {
	req := &fansMedalService.HighQpsLiveWearedReq{
		Uid: uid,
	}
	resq, err := FansMedalServiceClient.V2HighQps.LiveWeared(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: fansMedalService LiveWeared err: %v", err)
		}
		f.RUID = 0
		f.MedalLevel = 1
		f.MedalName = ""
		f.MColor = 0
		return nil
	}
	if resq.Code != 0 {
		log.Error("DM: fansMedalService LiveWeared err code: %d", resq.Code)
		f.RUID = 0
		f.MedalLevel = 1
		f.MedalName = ""
		f.MColor = 0
		return nil
	}
	f.RUID = resq.Data.TargetId
	f.MedalLevel = resq.Data.Level
	f.MedalName = resq.Data.MedalName
	f.MColor = resq.Data.MedalColor
	return nil
}

//Get 获取主播房间配置信息s
func (r *RoomConf) Get(ctx context.Context, roomID int64) error {
	req := &roomService.RoomGetByIdsReq{
		Ids: []int64{roomID},
	}
	resp, err := RoomServiceClient.V2Room.GetByIds(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: roomService GetByIds err: %v", err)
		}
		r.RoomID = roomID
		r.UID = 0
		r.RoomShield = 0
		r.Anchor = ""
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: roomService GetByIds err code: %d", resp.Code)
		r.RoomID = roomID
		r.UID = 0
		r.RoomShield = 0
		r.Anchor = ""
		return nil
	}
	if _, ok := resp.Data[roomID]; !ok {
		log.Error("DM: roomService GetByIds error roomid:%d", roomID)
		r.RoomID = roomID
		r.UID = 0
		r.RoomShield = 0
		r.Anchor = ""
		return nil
	}
	r.RoomID = resp.Data[roomID].Roomid
	r.UID = resp.Data[roomID].Uid
	r.RoomShield = resp.Data[roomID].RoomShield
	r.Anchor = resp.Data[roomID].Uname
	return nil
}

//Get 获取用户绑定信息,实名认知,绑定手机等信息
func (u *UserBindInfo) Get(ctx context.Context, uid int64) error {
	req := &acctountService.MidReq{
		Mid: uid,
	}
	resp, err := ac.Profile3(ctx, req)
	if err != nil {
		log.Error("DM: acctountService Profile3 err: %v", err)
		u.Identification = 0
		u.MobileVerify = 0
		u.Uname = ""
		u.URank = 0
		return nil
	}
	u.Identification = resp.Profile.Identification
	u.MobileVerify = resp.Profile.TelStatus
	u.Uname = resp.Profile.Name
	u.URank = resp.Profile.Rank
	return nil
}

//GetUserScore 用户真实分
func (u *UserScore) GetUserScore(ctx context.Context, uid int64) error {
	req := &spyService.InfoReq{
		Mid: uid,
	}
	resp, err := SpyClient.Info(ctx, req)
	if err != nil {
		log.Error("DM: 获取用户真实分错误 err: %v", err)
		u.lock.Lock()
		u.UserScore = 0
		u.lock.Unlock()
		return nil
	}
	u.lock.Lock()
	u.UserScore = int64(resp.Ui.Score)
	u.lock.Unlock()
	return nil
}

//GetMsgScore 获取弹幕AI分
func (u *UserScore) GetMsgScore(ctx context.Context, msg string) error {
	req := &filterService.FilterReq{
		Area:    "live_danmu",
		Message: msg,
	}
	resp, err := FilterClient.Filter(ctx, req)
	if err != nil {
		log.Error("DM: main filter err: %v", err)
		u.lock.Lock()
		u.MsgLevel = 0
		u.MsgAI = 0
		u.lock.Unlock()
		return nil
	}

	u.lock.Lock()
	u.MsgLevel = int64(resp.Level)
	if resp.Ai == nil {
		log.Error("DM: main filter err: miss ai scores")
		u.MsgAI = 0
		u.lock.Unlock()
		return nil
	}
	if len(resp.Ai.Scores) == 0 {
		u.MsgAI = 0
		u.lock.Unlock()
		return nil
	}
	u.MsgAI = int64(resp.Ai.Scores[0] * 10)
	u.lock.Unlock()
	return nil
}
