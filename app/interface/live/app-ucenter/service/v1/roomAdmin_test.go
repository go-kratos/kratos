package v1

import (
	"context"
	"flag"
	"go-common/app/interface/live/app-ucenter/api/http/v1"
	"go-common/app/interface/live/app-ucenter/conf"
	"go-common/library/net/metadata"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/smartystreets/goconvey/convey"
)

var roomadmin *RoomAdminService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	roomadmin = NewRoomAdminService(conf.Conf)
}

func TestV1ShowEntry(t *testing.T) {
	convey.Convey("ShowEntry", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(10000),
		})
		res, err := roomadmin.ShowEntry(ctx, &v1.ShowEntryReq{})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

func TestV1SearchForAdmin(t *testing.T) {
	convey.Convey("SearchForAdmin", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(10000),
		})
		res, err := roomadmin.SearchForAdmin(ctx, &v1.RoomAdminSearchForAdminReq{
			KeyWord: "test",
		})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

func TestV1IsAny(t *testing.T) {
	convey.Convey("IsAny", t, func() {
		ctx := metadata.NewContext(context.TODO(), metadata.MD{
			metadata.Mid: int64(10000),
		})
		res, err := roomadmin.IsAny(ctx, &v1.ShowEntryReq{})
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
