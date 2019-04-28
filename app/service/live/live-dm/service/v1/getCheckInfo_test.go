package v1

import (
	"context"
	"encoding/json"
	"fmt"
	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"testing"
)

//group=qa01 DEPLOY_ENV=uat go test -race  -run TestGetCheckMsg
func TestGetCheckMsg(t *testing.T) {
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:      111,
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
			MedalInfo: &dao.FansMedalInfo{},
		},
		RoomConf:     &dao.RoomConf{},
		UserBindInfo: &dao.UserBindInfo{},
		TitleConf:    &dao.CommentTitle{},
		DMconf:       &dao.DMConf{},
		UserScore:    &dao.UserScore{},
	}
	if err := getCheckMsg(context.TODO(), s); err != nil {
		t.Error("获取检查配置失败:", err)
	}
	j, _ := json.Marshal(s)
	fmt.Print("####", string(j))
}

//group=qa01 DEPLOY_ENV=uat go test -race -run TestGetRoomConf
func TestGetRoomConf(t *testing.T) {
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		UserBindInfo: &dao.UserBindInfo{
			Identification: 0,
			MobileVerify:   0,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    1877309,
			Roomid: 5392,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID: 158807,
		},
	}
	if err := getRoomConf(context.TODO(), s); err != nil {
		t.Error("房间配置获取失败")
	}
	fmt.Printf("######%+v", s)
	fmt.Printf("######%+v", s.UserInfo)
	fmt.Printf("######%+v", s.RoomConf)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestGetPrivilegeType
func TestGetPrivilegeType(t *testing.T) {
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		UserBindInfo: &dao.UserBindInfo{
			Identification: 0,
			MobileVerify:   0,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    1877309,
			Roomid: 5392,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID: 158807,
		},
	}
	if err := s.UserInfo.GetPrivilegeType(context.TODO(), 1877309, 999); err != nil {
		t.Error("大航海获取失败")
	}
	fmt.Printf("######%d", s.UserInfo.PrivilegeType)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestIsRoomAdim111
func TestIsRoomAdim111(t *testing.T) {
	var s = &SendMsg{
		UserInfo: &dao.UserInfo{
			UserLever: 1,
		},
		UserBindInfo: &dao.UserBindInfo{
			Identification: 0,
			MobileVerify:   0,
		},
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:    1877309,
			Roomid: 5392,
			Ip:     "115.239.211.112",
		},
		Dmservice: &DMService{
			conf: conf.Conf,
		},
		RoomConf: &dao.RoomConf{
			UID: 158807,
		},
	}
	s.UserInfo.IsRoomAdmin(context.TODO(), 110000654, 460874)
	if !s.UserInfo.RoomAdmin {
		t.Error("房管判断失败")
	}

	s.UserInfo.IsRoomAdmin(context.TODO(), 110000655, 460874)
	if s.UserInfo.RoomAdmin {
		t.Error("房管判断失败")
	}
}
