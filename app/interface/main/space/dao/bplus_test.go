package dao

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_GroupsCount(t *testing.T) {
	convey.Convey("test group count", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.groupsCountURL).Reply(200).JSON(`{"code":0,"data":{"num":1}}`)
		mid := int64(28272030)
		vmid := int64(28272030)
		data, err := d.GroupsCount(context.Background(), mid, vmid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)
		convey.Printf("%d", data)
	})
}

func TestDao_DynamicCnt(t *testing.T) {
	convey.Convey("test dynamic cnt", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.dynamicCntURL).Reply(200).JSON(`{"code":0,"msg":"","message":"","data":{"items":[{"uid":2089809,"num":345}],"_gt_":0}}`)
		vmid := int64(2089809)
		data, err := d.DynamicCnt(context.Background(), vmid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%d", data)
	})
}

func TestDao_DynamicList(t *testing.T) {
	convey.Convey("test dynamic list", t, func(ctx convey.C) {
		mid := int64(29313802)
		vmid := int64(34709144)
		data, err := d.DynamicList(context.Background(), mid, vmid, 0, 16, 1)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%d", data)
	})
}

func TestDao_Dynamic(t *testing.T) {
	convey.Convey("test dynamic item", t, func(ctx convey.C) {
		mid := int64(27515256)
		dyID := int64(118606711587078278)
		data, err := d.Dynamic(context.Background(), mid, dyID, 16)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}
