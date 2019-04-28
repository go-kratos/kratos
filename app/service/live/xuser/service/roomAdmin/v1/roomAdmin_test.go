package v1

import (
	"flag"
	"github.com/smartystreets/goconvey/convey"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"testing"
)

var (
	RoomAdmin *RoomAdminService
)

func init() {
	flag.Set("conf", "../../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	RoomAdmin = NewRoomAdminService(conf.Conf)
}

// go test  -test.v -test.run TestRoomAdminService_IsAny
func TestRoomAdminService_IsAny(t *testing.T) {
	convey.Convey("TestRoomAdminService_IsAny", t, func() {

		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := RoomAdmin.IsAny(ctx, &v1.RoomAdminShowEntryReq{
			Uid: 10000,
		})
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestRoomAdminService_IsAdminShort
func TestRoomAdminService_IsAdminShort(t *testing.T) {
	convey.Convey("TestRoomAdminService_IsAdminShort", t, func() {

		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := RoomAdmin.IsAdminShort(ctx, &v1.RoomAdminIsAdminShortReq{
			Uid:    10000,
			Roomid: 1,
		})
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}
