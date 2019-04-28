package dao

import (
	"context"

	activityService "go-common/app/service/live/activity/api/liverpc/v1"
	rankdbService "go-common/app/service/live/rankdb/api/liverpc/v1"
	rcService "go-common/app/service/live/rc/api/liverpc/v1"
	roomService "go-common/app/service/live/room/api/liverpc/v2"
	userextService "go-common/app/service/live/userext/api/liverpc/v1"
	acctountService "go-common/app/service/main/account/api"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//CommentTitle 头衔信息
type CommentTitle struct {
	OldTitle string `json:"oldtitle"`
	Title    string `json:"title"`
}

//GetUnameColor 获取用户的昵称颜色
func (u *UserInfo) GetUnameColor(ctx context.Context, uid int64, rid int64) error {
	req := &userextService.ColorGetUnameColorReq{
		Uid:    uid,
		RoomId: rid,
	}
	resq, err := UserExtServiceClient.V1Color.GetUnameColor(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: UserExt  GetUnameColor err: %v", err)
		}
		u.lock.Lock()
		u.UnameColor = ""
		u.lock.Unlock()
		return nil
	}
	if resq.Code != 0 {
		log.Error("DM:  UserExt GetUnameColor errror code: %d", resq.Code)
		u.lock.Lock()
		u.UnameColor = ""
		u.lock.Unlock()
		return nil
	}
	u.lock.Lock()
	u.UnameColor = resq.Data.UnameColor
	u.lock.Unlock()
	return nil
}

//GetSpeicalMedal 获取特殊勋章
func (f *FansMedalInfo) GetSpeicalMedal(ctx context.Context, uid int64, rid int64) error {
	//更新special
	req := &activityService.UnionFansGetSpecialMedalReq{
		Uid:  uid,
		Ruid: rid,
	}
	resq, err := ActivityServiceClient.V1UnionFans.GetSpecialMedal(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: GetSpecialMedal err: %v", err)
		}
		f.lock.Lock()
		f.SpecialMedal = ""
		f.lock.Unlock()
		return nil
	}
	if resq.Code != 0 {
		log.Error("DM: GetSpecialMedal err code: %d", resq.Code)
		f.lock.Lock()
		f.SpecialMedal = ""
		f.lock.Unlock()
		return nil
	}
	f.lock.Lock()
	f.SpecialMedal = resq.Data.SpecialMedal
	f.lock.Unlock()
	return nil
}

//GetUserLevelRank 获取用户等级RANK
func (u *UserInfo) GetUserLevelRank(ctx context.Context, uid int64) error {
	defer u.lock.Unlock()
	req := &rankdbService.UserRankGetUserRankReq{
		Uid:  uid,
		Type: "user_level",
	}
	resq, err := RankdbServiceClient.V1UserRank.GetUserRank(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: GetUserRank err: %v", err)
		}
		u.lock.Lock()
		u.ULevelRank = 1000000
		return nil
	}
	if resq.Code != 0 {
		log.Error("DM: GetUserRank error code: %d", resq.Code)
		u.lock.Lock()
		u.ULevelRank = 1000000
		return nil
	}
	u.lock.Lock()
	u.ULevelRank = resq.Data.Rank
	return nil
}

//GetCommentTitle 获取头衔
func (c *CommentTitle) GetCommentTitle(ctx context.Context) error {
	req := &rcService.UserTitleGetCommentTitleReq{}
	resq, err := RcServiceClient.V1UserTitle.GetCommentTitle(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: GetCommentTitle err: %v", err)
		}
		c.OldTitle = ""
		c.Title = ""
		return nil
	}
	if resq.Code != 0 {
		log.Error("DM: GetCommentTitle error code: %d", resq.Code)
		c.OldTitle = ""
		c.Title = ""
		return nil
	}
	if len(resq.Data) == 0 {
		c.OldTitle = ""
		c.Title = ""
		return nil
	}

	c.OldTitle = resq.Data[0]
	c.Title = resq.Data[1]
	return nil
}

//GetMedalanchorName 获取勋章对应主播的昵称
func (f *FansMedalInfo) GetMedalanchorName(ctx context.Context, uid int64) error {
	if uid == 0 {
		f.RUName = ""
		return nil
	}

	req := &acctountService.MidReq{
		Mid: uid,
	}
	resp, err := ac.Profile3(ctx, req)
	if err != nil {
		log.Error("DM: acctountService Profile3 err: %v", err)
		f.lock.Lock()
		f.RUName = ""
		f.lock.Unlock()
		return nil
	}
	f.lock.Lock()
	f.RUName = resp.Profile.Name
	f.lock.Unlock()
	return nil
}

//GetMedalRoomid 获取勋章对应主播的房间
func (f *FansMedalInfo) GetMedalRoomid(ctx context.Context, uid int64) error {
	if uid == 0 {
		f.lock.Lock()
		f.RoomID = 0
		f.lock.Unlock()
		return nil
	}
	req := &roomService.RoomRoomIdByUidReq{
		Uid: uid,
	}
	resp, err := RoomServiceClient.V2Room.RoomIdByUid(ctx, req)
	if err != nil {
		log.Error("DM: room  RoomIdByUid err: %v", err)
		f.lock.Lock()
		f.RoomID = 0
		f.lock.Unlock()
		return nil
	}
	if resp.Code == 404 {
		f.lock.Lock()
		f.RoomID = 0
		f.lock.Unlock()
		return nil
	}

	if resp.Code != 0 {
		log.Error("DM: room  RoomIdByUid err code: %d", resp.Code)
		f.lock.Lock()
		f.RoomID = 0
		f.lock.Unlock()
		return nil
	}
	f.lock.Lock()
	f.RoomID = resp.Data.RoomId
	f.lock.Unlock()
	return nil
}

//GetUserBubble 判断用户是否有气泡
func (u *UserInfo) GetUserBubble(ctx context.Context, uid int64, roomid int64, bubble int64, guardLevel int) error {
	req := &userextService.BubbleGetBubbleReq{
		Uid:        uid,
		RoomId:     roomid,
		BubbleId:   bubble,
		GuardLevel: int64(guardLevel),
	}
	defer u.lock.Unlock()
	resp, err := UserExtServiceClient.V1Bubble.GetBubble(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: userext Bubble check err: %v", err)
		}
		u.lock.Lock()
		u.Bubble = bubble
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: userext Bubble check err code: %d", resp.Code)
		u.lock.Lock()
		u.Bubble = bubble
		return nil
	}
	if resp.Data == nil {
		log.Error("DM: userext Bubble check err not data")
		u.lock.Lock()
		u.Bubble = bubble
		return nil
	}
	u.lock.Lock()
	u.Bubble = resp.Data.Bubble
	return nil
}
