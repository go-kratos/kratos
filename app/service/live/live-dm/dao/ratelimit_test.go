package dao

import (
	"context"
	"flag"
	"go-common/app/service/live/live-dm/conf"
	"path/filepath"
	"testing"
)

func init() {
	dir, _ := filepath.Abs("../cmd/test.toml")
	flag.Set("conf", dir)
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	InitAPI()
	InitGrpc(conf.Conf)
}

//group=qa01 DEPLOY_ENV=uat go test  -run TestLimitPerSec
func TestLimitPerSec(t *testing.T) {
	l := LimitCheckInfo{
		UID:     111,
		RoomID:  222,
		Msg:     "6666",
		Dao:     New(conf.Conf),
		MsgType: 0,
		Conf: &LimitConf{
			DmNum:     20,
			DMPercent: 25,
		},
	}
	if err := l.LimitPerSec(context.TODO()); err != nil {
		t.Error("每秒限制错误:", err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test  -run TestLimitSameMsg
func TestLimitSameMsg(t *testing.T) {
	l := LimitCheckInfo{
		UID:     111,
		RoomID:  222,
		Msg:     "6666",
		Dao:     New(conf.Conf),
		MsgType: 0,
		Conf: &LimitConf{
			DmNum:     20,
			DMPercent: 25,
		},
	}
	if err := l.LimitSameMsg(context.TODO()); err != nil {
		t.Error("5秒相同发言错误:", err)
	}
}

//group=qa01 DEPLOY_ENV=uat go test  -run TestLimitRoomPerSecond
func TestLimitRoomPerSecond(t *testing.T) {
	l := LimitCheckInfo{
		UID:     111,
		RoomID:  222,
		Msg:     "6666",
		Dao:     New(conf.Conf),
		MsgType: 0,
		Conf: &LimitConf{
			DmNum:     20,
			DMPercent: 25,
		},
	}
	if err := l.LimitRoomPerSecond(context.TODO()); err != nil {
		t.Error("每秒20弹幕错误:", err)
	}
}
