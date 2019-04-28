package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/interface/live/app-blink/api/http/v1"
	"go-common/library/net/metadata"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/interface/live/app-blink/conf"
)

var room *RoomService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	room = NewRoomService(conf.Conf)
}

// group=qa01 DEPLOY_ENV=uat go test -run TestGetRoomInfo
func TestGetRoomInfo(t *testing.T) {
	Convey("TestNewRoomService", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(16299525),
		})
		res, err := room.GetInfo(ctx, &v1.GetRoomInfoReq{
			Platform: "ios",
		})
		fmt.Println(1111, res, err, 22222)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

func TestCreate(t *testing.T) {
	Convey("TestNewRoomService", t, func() {

		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(16299525),
		})
		res, err := room.Create(ctx, &v1.CreateReq{
			Platform: "ios",
		})
		fmt.Println(1111, res, err, 22222)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
