package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
)

var (
	srv *ResourceService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	srv = NewResourceService(conf.Conf)
}

// go test  -test.v -test.run TestAdd
func TestAdd(t *testing.T) {
	Convey("TestAdd", t, func() {
		res, err := srv.Add(context.TODO(), &pb.AddReq{
			Platform:     "mng",
			Title:        "splash first",
			JumpPath:     "www.sina.com",
			JumpPathType: 1,
			JumpTime:     3,
			Type:         "splash",
			Device:       "[{\"platform\":\"pc\",\"build\":3,\"limit\":1}]",
			StartTime:    "2018-01-01 00:00:00",
			EndTime:      "2018-11-10 00:00:00",
			ImageUrl:     "www.baidu.com",
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
func TestEdit(t *testing.T) {
	Convey("TestEdit", t, func() {
		res, err := srv.Edit(context.TODO(), &pb.EditReq{
			Platform:     "mng",
			Id:           23,
			Title:        "splash second",
			JumpPath:     "www.sina.com",
			JumpPathType: 1,
			JumpTime:     6,
			StartTime:    "2018-1-1 00:00:00",
			EndTime:      "2018-1-9 00:00:00",
			ImageUrl:     "www.sina.com",
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
func TestGetList(t *testing.T) {
	Convey("TestGetList", t, func() {
		res, err := srv.GetList(context.TODO(), &pb.GetListReq{
			Platform: "mng",
			Type:     "splash",
			Page:     1,
			PageSize: 10,
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
func TestOffline(t *testing.T) {
	Convey("TestOffline", t, func() {
		res, err := srv.Offline(context.TODO(), &pb.OfflineReq{
			Platform: "mng",
			Id:       23,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
