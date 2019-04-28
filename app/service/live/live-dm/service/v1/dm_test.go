package v1

import (
	"context"
	"flag"
	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"path/filepath"
	"testing"
)

func init() {
	dir, _ := filepath.Abs("../../cmd/test.toml")
	flag.Set("conf", dir)
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao.InitAPI()
	dao.InitGrpc(conf.Conf)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestDMService_SendMsg
func TestDMService_SendMsg(t *testing.T) {

}

//group=qa01 DEPLOY_ENV=uat go test -run TestGetHistory
func TestGetHistory(t *testing.T) {
	req := &v1pb.HistoryReq{
		Roomid: 460828,
	}

	s := &DMService{
		conf: conf.Conf,
		dao:  dao.New(conf.Conf),
	}
	s.GetHistory(context.TODO(), req)
}
