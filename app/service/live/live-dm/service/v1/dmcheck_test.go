package v1

import (
	"context"
	"fmt"
	"testing"

	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
)

//group=qa01 DEPLOY_ENV=uat go test -race -run  TestCheckLegitimacy
func TestCheckLegitimacy(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:      1877309,
			Roomid:   5392,
			Ip:       "115.239.211.112",
			Msg:      "6666",
			Rnd:      "1000",
			Fontsize: 15,
			Mode:     1,
			Platform: "ios",
			Msgtype:  0,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		UserInfo: &dao.UserInfo{
			UserLever: 100,
			MedalInfo: &dao.FansMedalInfo{},
		},
		RoomConf: &dao.RoomConf{
			UID:    158807,
			RoomID: 5392,
			Anchor: "沢奇Sawaki",
		},
		UserBindInfo: &dao.UserBindInfo{
			Uname: "Bili_111",
			URank: 1000,
		},
		TitleConf: &dao.CommentTitle{},
		DMconf: &dao.DMConf{
			Mode:   6,
			Color:  5555,
			Length: 20,
		},
		UserScore: &dao.UserScore{
			UserScore: 90,
		},
	}
	if err := checkLegitimacy(context.TODO(), s); err != nil {
		fmt.Print("####\n", err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestGlobalRestriction
func TestGlobalRestriction(t *testing.T) {
	s := &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
	}
	if err := globalRestriction(context.TODO(), s); err != nil {
		if err != ecode.DMUserLevel {
			t.Error("全局禁言检查错误:", err)
		}
	}
	s.Dmservice.conf.DmRules.LevelLimit = 10
	s.UserInfo.UserLever = 5
	if err := globalRestriction(context.TODO(), s); err != nil {
		if err != ecode.DMUserLevel {
			t.Error("全局等级禁言检查错误:", err)
		}
	}
	s.Dmservice.conf.DmRules.AllUserLimit = true
	if err := globalRestriction(context.TODO(), s); err != nil {
		if err != ecode.DMallUser {
			t.Error("全体禁言检查错误:", err)
		}
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestPaidRoomInspection
func TestPaidRoomInspection(t *testing.T) {
	ctx := metadata.NewContext(context.TODO(), metadata.MD{})
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    111,
			Roomid: 460460,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
	}
	if err := paidRoomInspection(ctx, s); err != nil {
		if err != ecode.PayLive {
			t.Error("付费检查失败")
		}
	}
	s.SendMsgReq.Uid = 123
	if err := paidRoomInspection(ctx, s); err != nil {
		t.Error("付费检查失败")
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_MessageCheck
func TestDmcheck_MessageCheck(t *testing.T) {
	s := &SendMsg{
		UserScore: &dao.UserScore{
			MsgLevel: 30,
		},
		SendMsgReq: &v1pb.SendMsgReq{},
	}
	if err := messageCheck(context.TODO(), s); err != nil {
		if err != ecode.FilterLimit {
			t.Error("弹幕内容检查失败 Level 30")
		}
	}
	s.UserScore.MsgLevel = 10
	if err := messageCheck(context.TODO(), s); err != nil {
		t.Error("弹幕内容检查失败 Level 10")
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_MessageLenCheck
func TestDmcheck_MessageLenCheck(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg: "6666",
		},
		DMconf: &dao.DMConf{Length: 30},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
	}
	if err := messageLenCheck(context.TODO(), s); err != nil {
		t.Error("弹幕长度检查失败 msg:6666")
	}
	s.SendMsgReq.Msg = "666666666666666666666"
	if err := messageLenCheck(context.TODO(), s); err != nil {
		t.Error("弹幕长度检查失败 msglength:21")
	}

	s.SendMsgReq.Msg = "一一一一一一一一一一一一一一一一一一一一一"
	s.DMconf.Length = 20
	if err := messageLenCheck(context.TODO(), s); err != nil {
		if err != ecode.MsgLengthLimit {
			t.Error("弹幕长度检查失败 msglength:21个一")
		}
	}

	s.SendMsgReq.Msg = "一一一一一一一一一一一一一一一一一一一一一"
	s.DMconf.Length = 30
	if err := messageLenCheck(context.TODO(), s); err != nil {
		t.Error("弹幕长度检查失败 msglength:21个一 长度限制30")
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_RoomInfoCheck
func TestDmcheck_RoomInfoCheck(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg:    "6666",
			Roomid: 460758,
			Uid:    222,
		},
		DMconf: &dao.DMConf{
			Length: 30,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID: 111,
		},
		UserInfo: &dao.UserInfo{
			RoomAdmin: false,
			UserLever: 10,
		},
	}
	if err := roomInfoCheck(context.TODO(), s); err != nil {
		if err == ecode.RoomAllLimit && err == ecode.RoomMedalLimit &&
			err == ecode.RoomLeverLimit {
			t.Error("房间禁言检查错误")
		}
	}
	s.UserInfo.PrivilegeType = 1

	if err := roomInfoCheck(context.TODO(), s); err != nil {
		t.Error("总督禁言检查错误")
	}

	s.UserInfo.PrivilegeType = 2

	if err := roomInfoCheck(context.TODO(), s); err != nil {
		if err == ecode.RoomAllLimit && err == ecode.RoomMedalLimit &&
			err == ecode.RoomLeverLimit {
			t.Error("房间禁言检查错误")
		}
	}
	s.UserInfo.PrivilegeType = 0
	s.UserInfo.UserLever = 6
	if err := roomInfoCheck(context.TODO(), s); err != nil {
		fmt.Print("11111", err)
		if err == ecode.RoomAllLimit && err == ecode.RoomMedalLimit &&
			err == ecode.RoomLeverLimit {
			t.Error("房间禁言检查错误")
		}
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_RoomBluckList
func TestDmcheck_RoomBluckList(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg:    "666",
			Roomid: 460758,
			Uid:    111,
		},
		DMconf: &dao.DMConf{
			Length: 30,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID:        110000298,
			RoomShield: 1,
		},
		UserInfo: &dao.UserInfo{
			RoomAdmin: false,
			UserLever: 10,
		},
	}
	if err := roomBluckList(context.TODO(), s); err != nil {
		t.Error("房间黑名单检查错误", err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_AnchorFilterMsg
func TestDmcheck_AnchorFilterMsg(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg:    "666",
			Roomid: 460758,
			Uid:    222,
		},
		DMconf: &dao.DMConf{
			Length: 30,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID:        110000298,
			RoomShield: 1,
		},
		UserInfo: &dao.UserInfo{
			RoomAdmin: false,
			UserLever: 10,
		},
	}
	if err := anchorFilterMsg(context.TODO(), s); err != nil {
		if err != ecode.ShieldContent {
			t.Error("主播屏蔽过滤词失败")
		}
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_AnchorFilterUser
func TestDmcheck_AnchorFilterUser(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg:    "666",
			Roomid: 460758,
			Uid:    222,
		},
		DMconf: &dao.DMConf{
			Length: 30,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID:        110000298,
			RoomShield: 1,
		},
		UserInfo: &dao.UserInfo{
			RoomAdmin: false,
			UserLever: 10,
		},
	}
	if err := anchorFilterUser(context.TODO(), s); err != nil {
		if err != ecode.ShieldContent {
			t.Error("主播屏蔽过滤用户失败")
		}
	}

}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_CanBroadCastMsg
func TestDmcheck_CanBroadCastMsg(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Msg:    "666",
			Roomid: 460758,
			Uid:    222,
		},
		DMconf: &dao.DMConf{
			Length: 30,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID:        110000298,
			RoomShield: 1,
		},
		UserInfo: &dao.UserInfo{
			RoomAdmin: false,
			UserLever: 10,
		},
		UserScore: &dao.UserScore{
			UserScore: 10,
			MsgAI:     0,
		},
	}
	if err := canBroadCastMsg(context.TODO(), s); err != nil {
		t.Error("真实分服务错误")
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_AuthenticationAuthority
func TestDmcheck_AuthenticationAuthority(t *testing.T) {
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		UserBindInfo: &dao.UserBindInfo{
			Identification: 0,
			MobileVerify:   0,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    111,
			Roomid: 460460,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
	}
	if err := authenticationAuthority(context.TODO(), s); err != nil {
		if err != ecode.RealName {
			t.Error("实名认证检查失败")
		}
	}
	s.UserBindInfo.Identification = 1
	if err := authenticationAuthority(context.TODO(), s); err != nil {
		if err != ecode.PhoneBind {
			t.Error("手机绑定检查失败")
		}
	}
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDmcheck_UserArea
// func TestDmcheck_UserArea(t *testing.T) {
// 	var s = &SendMsg{
// 		UserInfo: &dao.UserInfo{
// 			UserLever: 1,
// 		},
// 		SendMsgReq: &v1pb.SendMsgReq{
// 			Uid:    111,
// 			Roomid: 460460,
// 			Ip:     "115.239.211.112",
// 		},
// 		Dmservice: &DMService{
// 			conf: conf.Conf,
// 		},
// 	}
// 	if err := userArea(context.TODO(), s); err != nil {
// 		t.Error("区域判断出错")
// 	}
// 	s.SendMsgReq.Ip = "198.11.179.19"
// 	if err := userArea(context.TODO(), s); err != nil {
// 		if err != ecode.CountryLimit {
// 			t.Error("区域判断出错 美国")
// 		}
// 	}
// }

//group=fat1 DEPLOY_ENV=uat go test -run TestDmcheck_userAreaLive
func TestDmcheck_userAreaLive(t *testing.T) {
	dao.InitIPdb()
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    111,
			Roomid: 460460,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
	}
	err := userAreaLive(s)
	if err != nil {
		t.Error("#####", err)
	}
}
