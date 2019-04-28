package v1

import (
	"context"
	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"testing"
)

//group=qa01 DEPLOY_ENV=uat go test -run TestSaveHistory
func TestSaveHistory(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:      1877309,
			Roomid:   5392,
			Ip:       "115.239.211.112",
			Msg:      "66666",
			Rnd:      "1000",
			Fontsize: 15,
			Mode:     0,
			Platform: "ios",
			Msgtype:  0,
		},
		Dmservice: &DMService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		},
		UserInfo: &dao.UserInfo{
			ULevelColor: 9868950,
			ULevelRank:  2,
			MedalInfo:   &dao.FansMedalInfo{},
			RoomAdmin:   false,
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

	saveHistory(context.TODO(), s)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestIncrDMNum
func TestIncrDMNum(t *testing.T) {
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
			ULevelColor: 9868950,
			ULevelRank:  2,
			MedalInfo:   &dao.FansMedalInfo{},
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
	incrDMNum(context.TODO(), s)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestSendBroadCast
func TestSendBroadCast(t *testing.T) {
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
			ULevelColor: 9868950,
			ULevelRank:  2,
			MedalInfo:   &dao.FansMedalInfo{},
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
	sendBroadCast(context.TODO(), s)
}
