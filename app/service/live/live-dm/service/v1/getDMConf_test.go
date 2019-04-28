package v1

import (
	"context"
	"encoding/json"
	"fmt"
	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"go-common/library/net/metadata"
	"testing"
)

//group=qa01 DEPLOY_ENV=uat go test -race -run TestGetDMconfig
func TestGetDMconfig(t *testing.T) {
	ctx := metadata.NewContext(context.TODO(), metadata.MD{})
	s := &SendMsg{
		SendMsgReq: &v1pb.SendMsgReq{
			Uid:      520,
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
	if err := getDMconfig(ctx, s); err != nil {
		t.Error("获取弹幕配置失败", err)
	}
	j, _ := json.Marshal(s)
	fmt.Print("####", string(j))
}
