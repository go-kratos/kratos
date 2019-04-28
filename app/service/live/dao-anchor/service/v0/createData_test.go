package v0

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/service/live/dao-anchor/api/grpc/v0"

	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/library/net/metadata"

	. "github.com/smartystreets/goconvey/convey"
)

var s *CreateDataService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = NewCreateDataService(conf.Conf)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestCreateCache
func TestCreateCache(t *testing.T) {
	Convey("TestCreateCache", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{})
		replay, err := s.CreateLiveCacheList(ctx, &v0.CreateLiveCacheListReq{Content: dao.DANMU_NUM, RoomIds: []int64{1, 1003, 1005}})
		fmt.Println(replay, err)
		So(err, ShouldBeNil)
	})
}

func TestCreateDBData(t *testing.T) {
	Convey("TestCreateDBData", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{})
		reply, err := s.CreateDBData(ctx, &v0.CreateDBDataReq{Content: dao.DANMU_NUM, RoomIds: []int64{1003, 1005}})
		fmt.Println(reply, err)
		So(err, ShouldBeNil)
	})

}
