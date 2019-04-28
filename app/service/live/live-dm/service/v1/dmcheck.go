package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ipipdotnet/ipdb-go"

	avService "go-common/app/service/live/av/api/liverpc/v1"
	bannedService "go-common/app/service/live/banned_service/api/liverpc/v1"
	"go-common/app/service/live/live-dm/dao"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

//检查弹幕合法性
func checkLegitimacy(ctx context.Context, sdm *SendMsg) error {
	g, ctx := errgroup.WithContext(ctx)
	//防止竞态条件，透传mid以及userip
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Mid] = sdm.SendMsgReq.Uid
		md[metadata.RemoteIP] = sdm.SendMsgReq.Ip
	}
	//全局禁言检查
	g.Go(func() error { return globalRestriction(ctx, sdm) })
	//付费直播检查
	g.Go(func() error { return paidRoomInspection(ctx, sdm) })
	//弹幕长度检查
	g.Go(func() error { return messageLenCheck(ctx, sdm) })
	//房间配置禁言检查
	g.Go(func() error { return roomInfoCheck(ctx, sdm) })
	// 房间禁言名单
	g.Go(func() error { return roomBluckList(ctx, sdm) })
	// 主播过滤用户
	g.Go(func() error { return anchorFilterUser(ctx, sdm) })
	// 主播过滤内容
	g.Go(func() error { return anchorFilterMsg(ctx, sdm) })
	// 用户绑定信息判断是否禁言
	g.Go(func() error { return authenticationAuthority(ctx, sdm) })
	// 用户区域检查
	if sdm.LimitConf.AreaLimit {
		g.Go(func() error { return userAreaLive(sdm) })
	}
	//弹幕内容检查
	g.Go(func() error { return messageCheck(ctx, sdm) })
	// 调用真实接口判断，弹幕是否可以广播
	g.Go(func() error { return canBroadCastMsg(ctx, sdm) })
	//2019 拜年祭白名单
	roomid := strconv.FormatInt(sdm.SendMsgReq.Roomid, 10)
	if sdm.LimitConf.DMwhitelist && sdm.Dmservice.conf.BNJRoomList[roomid] {
		g.Go(func() error { return dmWhiteList(ctx, sdm) })
	}

	return g.Wait()
}

//DMWhiteList 拜年祭白名单
func dmWhiteList(ctx context.Context, sdm *SendMsg) error {
	if sdm.LimitConf.DMwhitelistID == "ID-0" {
		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"非拜年祭白名单用户", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "非拜年祭白名单用户")
		return ecode.Error(0, "w")
	}

	key := fmt.Sprintf("%s_%d", sdm.LimitConf.DMwhitelistID, sdm.SendMsgReq.Uid)
	if isWhite := sdm.Dmservice.dao.IsWhietListUID(ctx, key); isWhite {
		return nil
	}

	log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
		"非拜年祭白名单用户", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
		sdm.SendMsgReq.Msg)
	lancer(sdm, "非拜年祭白名单用户")
	return ecode.Error(0, "w")
}

//globalRestriction 全局禁言检查
func globalRestriction(ctx context.Context, sdm *SendMsg) error {
	if sdm.LimitConf.AllUserLimit {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"开启全局禁言", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)

		lancer(sdm, "开启全局禁言")
		return ecode.Error(0, "系统正在维护")
	}

	if sdm.UserInfo.UserLever <= sdm.LimitConf.LevelLimit && sdm.LimitConf.LevelLimitStatus {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s, limitLevel:%d, userLever:%d]",
			"全局等级禁言", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg, sdm.LimitConf.LevelLimit, sdm.UserInfo.UserLever)

		lancer(sdm, "开启全局等级禁言")
		return ecode.Error(0, "系统正在维护")
	}

	//全站禁言
	req := &bannedService.SiteBlockMngIsBlockUserReq{
		Tuid: sdm.SendMsgReq.Uid,
	}

	resp, err := dao.BannedServiceClient.V1SiteBlockMng.IsBlockUser(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: banned_service V1SiteBlockMng IsBlockUser err : %+v", err)
		}
		return nil
	}

	if resp.Code != 0 {
		log.Error("DM: banned_service V1SiteBlockMng IsBlockUser err code: %d", resp.Code)
		return nil
	}
	if resp.Data.IsBlock {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"直播黑名单用户", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "直播黑名单用户")
		return ecode.Error(0, "你被禁言啦")
	}
	return nil
}

//
//PaidRoomInspection 付费直播检查
func paidRoomInspection(ctx context.Context, sdm *SendMsg) error {
	req := &avService.PayLiveLiveValidateReq{
		RoomId:   sdm.SendMsgReq.Roomid,
		Platform: "grpc",
	}
	// if md, ok := metadata.FromContext(ctx); ok {
	// 	md[metadata.Mid] = sdm.SendMsgReq.Uid
	// 	md[metadata.RemoteIP] = sdm.SendMsgReq.Ip
	// }
	resp, err := dao.AvServiceClient.V1PayLive.LiveValidate(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: V1PayLive err: %v", err)
		}
		return nil
	}

	if resp.Code == 5001 {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"付费直播间", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)

		lancer(sdm, "付费直播间")
		return ecode.Error(-102, "非常抱歉，本场直播需要购票，即可参与互动")
	} else if resp.Code == 5003 {
		return nil
	}

	if resp.Code != 0 {
		log.Error("DM: V1OayLive error code: %d", resp.Code)

		lancer(sdm, "付费直播间")
		return ecode.Error(-102, "非常抱歉，本场直播需要购票，即可参与互动")
	}

	if resp.Data.Permission != 1 {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"付费直播间", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)

		lancer(sdm, "付费直播间")
		return ecode.Error(-102, "非常抱歉，本场直播需要购票，即可参与互动")
	}
	return nil
}

//MessageCheck 弹幕内容检查
func messageCheck(ctx context.Context, sdm *SendMsg) error {
	//命中等级大于等于20
	if sdm.UserScore.MsgLevel >= 20 {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"弹幕内容非法", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)

		lancer(sdm, "内容非法")
		return ecode.Error(0, "内容非法")
	}
	return nil
}

//MessageLenCheck 弹幕长度检查
func messageLenCheck(ctx context.Context, sdm *SendMsg) error {
	ml := len([]rune(sdm.SendMsgReq.Msg))
	// 小于默认长度 允许发送
	if ml < sdm.LimitConf.MsgLength {
		return nil
	}
	// 弹幕长度
	if ml <= int(sdm.DMconf.Length) {
		return nil
	}

	log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s, msgLength:%d, limitLength:%d]",
		"超出限制长度", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
		sdm.SendMsgReq.Msg, ml, sdm.DMconf.Length)
	lancer(sdm, "超出限制长度")
	return ecode.Error(-500, "超出限制长度")
}

//RoomInfoCheck 房间配置禁言检查
func roomInfoCheck(ctx context.Context, sdm *SendMsg) error {
	if sdm.SendMsgReq.Uid == sdm.RoomConf.UID ||
		sdm.UserInfo.RoomAdmin {
		return nil
	}
	//查询房间禁言
	req := &bannedService.SilentGetRoomSilentReq{
		RoomId: sdm.SendMsgReq.Roomid,
	}
	resp, err := dao.BannedServiceClient.V1Silent.GetRoomSilent(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: banned_service GetRoomSilent err: %v", err)
		}
		return nil
	}

	if resp.Code != 0 {
		log.Error("DM: banned_service GetRoomSilent err code: %d", resp.Code)
		return nil
	}

	switch resp.Data.Type {
	case "level":
		// 等级限制
		if resp.Data.Level < sdm.UserInfo.UserLever {
			return nil
		}
		// 老爷免疫
		if sdm.UserInfo.Vip != 0 || sdm.UserInfo.Svip != 0 {
			return nil
		}
		// 守护免疫
		if sdm.UserInfo.PrivilegeType != 0 {
			return nil
		}

		var errm string
		if resp.Data.Second == -1 {
			errm = "本场直播结束后解除~"
		} else {
			errm = strconv.FormatInt(resp.Data.Second/60, 10) + "分钟后解除～"
		}
		em := fmt.Sprintf("主播对UL %d以下开启了禁言，%s", resp.Data.Level, errm)

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			em, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "播主开启等级禁言")

		return ecode.Error(-403, em)
	case "medal":

		var errm string
		if resp.Data.Second == -1 {
			errm = "本场直播结束后解除~"
		} else {
			errm = strconv.FormatInt(resp.Data.Second/60, 10) + "分钟后解除～"
		}
		em := fmt.Sprintf("主播对粉丝勋章%d以下开启了禁言，%s", resp.Data.Level, errm)

		//佩戴勋章不一致
		if sdm.RoomConf.UID != sdm.UserInfo.MedalInfo.RUID {

			log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
				em, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
				sdm.SendMsgReq.Msg)
			lancer(sdm, "播主开启粉丝勋章限制")

			return ecode.Error(-403, em)
		}
		// 勋章等级免疫
		if resp.Data.Level < sdm.UserInfo.MedalInfo.MedalLevel {
			return nil
		}
		// 守护免疫(等级限制)
		pt := sdm.UserInfo.PrivilegeType
		if pt >= 1 && pt <= 2 {
			return nil
		}

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			em, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "播主开启粉丝勋章限制")

		return ecode.Error(-403, em)
	case "member":
		//守护免疫（总督免疫）
		if sdm.UserInfo.PrivilegeType == 1 {
			return nil
		}

		var errm string
		if resp.Data.Second == -1 {
			errm = "本场直播结束后解除~"
		} else {
			errm = strconv.FormatInt(resp.Data.Second/60, 10) + "分钟后解除～"
		}
		em := fmt.Sprintf("主播对全体用户开启了禁言，%s", errm)

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			em, sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "播开启房间全体禁言")

		return ecode.Error(-403, em)
	}
	return nil
}

//RoomBluckList 房间禁言名单
func roomBluckList(ctx context.Context, sdm *SendMsg) error {
	req := &bannedService.SilentMngIsBlockUserReq{
		Uid:    sdm.SendMsgReq.Uid,
		Roomid: sdm.SendMsgReq.Roomid,
		Type:   1,
	}

	resp, err := dao.BannedServiceClient.V1SilentMng.IsBlockUser(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: BannedService IsBlockUser err: %v", err)
		}
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: BannedService IsBlockUser err code: %d", resp.Code)
		return nil
	}
	if resp.Data.IsBlockUser {
		if sdm.UserInfo.PrivilegeType == 1 {
			return nil
		}
		tm := time.Unix(resp.Data.BlockEndTime, 0)
		ms := "你在本房间被禁言至 " + tm.Format("2006-01-02 15:04:05")

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"主播禁言名单", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "播主禁言名单")

		return ecode.Error(1003, ms)
	}
	return nil
}

//AnchorFilterUser 主播过滤用户
func anchorFilterUser(ctx context.Context, sdm *SendMsg) error {
	if sdm.RoomConf.RoomShield != 1 {
		return nil
	}
	req := &bannedService.ShieldMngIsShieldUserReq{
		Uid:       sdm.RoomConf.UID,
		ShieldUid: sdm.SendMsgReq.Uid,
	}
	resp, err := dao.BannedServiceClient.V1ShieldMng.IsShieldUser(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: BannedService IsShieldUser err: %v", err)
		}
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: BannedService IsShieldUser err code: %d", resp.Code)
		return nil
	}
	if resp.Data.IsShieldUser {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"主播过滤用户", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "被播主过滤的用户")

		return ecode.Error(0, "u")
	}
	return nil
}

//AnchorFilterMsg 主播过滤内容
func anchorFilterMsg(ctx context.Context, sdm *SendMsg) error {
	if sdm.RoomConf.RoomShield != 1 {
		return nil
	}

	req := &bannedService.ShieldIsShieldContentReq{
		Uid:     sdm.RoomConf.UID,
		Content: sdm.SendMsgReq.Msg,
	}
	resp, err := dao.BannedServiceClient.V1Shield.IsShieldContent(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: BannedService IsShieldContent err: %v", err)
		}
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: BannedService IsShieldContent err code: %d", resp.Code)
		return nil
	}
	if resp.Data.IsShieldContent {

		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
			"主播过滤内容", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg)
		lancer(sdm, "被播主过滤的弹幕内容")

		return ecode.Error(0, "k")
	}
	return nil
}

//CanBroadCastMsg 调用真实接口判断，弹幕是否可以广播
func canBroadCastMsg(ctx context.Context, sdm *SendMsg) error {
	req := &bannedService.AdminSilentGetShieldRuleReq{
		Roomid: sdm.SendMsgReq.Roomid,
	}
	resp, err := dao.BannedServiceClient.V1AdminSilent.GetShieldRule(ctx, req)
	if err != nil {
		if errors.Cause(err) != context.Canceled {
			log.Error("DM: BannedService GetShieldRule err: %v", err)
		}
		return nil
	}
	if resp.Code != 0 {
		log.Error("DM: BannedService GetShieldRule err code: %d", resp.Code)
		return nil
	}

	if sdm.UserScore.UserScore >= resp.Data.RealScore &&
		sdm.UserScore.MsgAI <= resp.Data.AiScore {
		return nil
	}

	if sdm.UserScore.UserScore < resp.Data.RealScore {
		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s, UserScore: %d, RealScore: %d]",
			"用户真实分拦截", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg, sdm.UserScore.UserScore, resp.Data.RealScore)
		lancer(sdm, "用户真实分拦截")
	}
	if sdm.UserScore.MsgAI > resp.Data.AiScore {
		log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s MsgAI: %d,  AiScore: %d]",
			"弹幕AI分拦截", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
			sdm.SendMsgReq.Msg, sdm.UserScore.MsgAI, resp.Data.AiScore)
		lancer(sdm, "弹幕AI分拦截")
	}

	return ecode.Error(0, "fire")
}

//AuthenticationAuthority 用户绑定信息判断是否禁言
func authenticationAuthority(ctx context.Context, sdm *SendMsg) error {
	if sdm.LimitConf.RealName {
		if 1 != sdm.UserBindInfo.Identification {
			log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
				"未实名认证", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
				sdm.SendMsgReq.Msg)
			lancer(sdm, "未实名认证")

			return ecode.Error(0, "实名认证才可以发言")
		}
	}

	mvf := sdm.UserBindInfo.MobileVerify
	if sdm.LimitConf.PhoneLimit {
		if mvf != 1 {
			log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]",
				"未实绑定手机", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
				sdm.SendMsgReq.Msg)
			lancer(sdm, "未绑定手机")

			return ecode.Error(1001, "根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作")
		}
	}
	return nil
}

//UserArea 用户区域判断
// func userArea(ctx context.Context, sdm *SendMsg) error {
// 	var ip = strings.Split(sdm.SendMsgReq.Ip, ":")
// 	req := &locationService.InfoReq{
// 		Addr: ip[0],
// 	}
// 	resp, err := dao.LcClient.Info(ctx, req)
// 	if err != nil {
// 		log.Error("DM: 主站区域查询接口失败 ERR: %v ip:%s", err, sdm.SendMsgReq.Ip)
// 		return nil
// 	}
// 	if resp.Country == "中国" || resp.Country == "局域网" {
// 		return nil
// 	}
// 	log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s,  用户IP:%s 用户区域:%s]",
// 		"用户所在区域无法发言", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
// 		sdm.SendMsgReq.Msg, sdm.SendMsgReq.Ip, resp.Country)

// 	return &errPb.Error{
// 		ErrCode:    0,
// 		ErrMessage: "你所在的地区暂无法发言",
// 	}
// }

//userAreaLive  用户区域判断
func userAreaLive(sdm *SendMsg) error {
	var cityInfo *ipdb.CityInfo
	var err error

	for _, v := range sdm.SendMsgReq.Ip {
		if string(v) == "." {
			cityInfo, err = dao.IP4db.FindInfo(sdm.SendMsgReq.Ip, "CN")
			break
		} else if string(v) == ":" {
			cityInfo, err = dao.IP6db.FindInfo(sdm.SendMsgReq.Ip, "CN")
			break
		}
	}

	if cityInfo == nil {
		log.Error("DM: ip解析失败: %s", sdm.SendMsgReq.Ip)
		return nil
	}
	if err != nil {
		log.Error("IPdb errr:%+v", err)
		return err
	}
	if cityInfo.CountryName == "中国" || cityInfo.CountryName == "局域网" {
		return nil
	}

	log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s,  用户IP:%s 用户区域:%s]",
		"用户所在区域无法发言", sdm.SendMsgReq.Uid, sdm.SendMsgReq.Roomid,
		sdm.SendMsgReq.Msg, sdm.SendMsgReq.Ip, cityInfo.CountryName)
	lancer(sdm, "用户区域限制")

	return ecode.Error(0, "你所在的地区暂无法发言")
}
