package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_TagSub(t *testing.T) {
	convey.Convey("test tag sub", t, func(ctx convey.C) {
		mid := int64(2089809)
		tid := int64(600)
		err := d.TagSub(context.Background(), mid, tid)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_TagCancelSub(t *testing.T) {
	convey.Convey("test cancel tag cancel sub", t, func(ctx convey.C) {
		mid := int64(2089809)
		tid := int64(600)
		err := d.TagCancelSub(context.Background(), mid, tid)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_TagSubList(t *testing.T) {
	convey.Convey("test tag sub list", t, func(ctx convey.C) {
		vmid := int64(88889018)
		pn := 1
		ps := 15
		data, count, err := d.TagSubList(context.Background(), vmid, pn, ps)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v,%d", data, count)
	})
}
