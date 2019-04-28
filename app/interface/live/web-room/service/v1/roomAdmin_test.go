package v1

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/conf"
	"go-common/library/net/metadata"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var roomAdmin *RoomAdminService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	roomAdmin = NewRoomAdminService(conf.Conf)
}

func TestV1NewRoomAdminService(t *testing.T) {
	convey.Convey("NewRoomAdminService", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(25156756),
		})

		res, err := roomAdmin.GetByRoom(ctx, &v1.RoomAdminGetByRoomReq{
			Page:     0,
			Roomid:   1008,
			PageSize: 0,
		})
		fmt.Println("------", res, err, "-------")
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}
