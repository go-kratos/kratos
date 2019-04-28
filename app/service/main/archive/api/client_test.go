package api

import (
	"context"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

var client ArchiveClient

func init() {
	var err error
	client, err = NewClient(nil)
	if err != nil {
		panic(err)
	}
}

func TestTypes(t *testing.T) {
	convey.Convey("Types", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.Types(c, &NoArgRequest{})
			ctx.So(err, convey.ShouldBeNil)
			for k, v := range reply.Types {
				ctx.Printf("key:%d id:%d name:%s pid:%d\n", k, v.ID, v.Name, v.Pid)
			}
		})
	})
}

func TestArc(t *testing.T) {
	convey.Convey("TestArc", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.Arc(c, &ArcRequest{Aid: 10100696})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("%+v\n", reply.Arc)
		})
		ctx.Convey("When error", func(ctx convey.C) {
			reply, err := client.Arc(context.TODO(), &ArcRequest{Aid: 99999999999})
			ctx.So(err, convey.ShouldEqual, ecode.NothingFound)
			ctx.So(reply, convey.ShouldBeNil)
		})
	})
}

func TestArcs(t *testing.T) {
	convey.Convey("TestArcs", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.Arcs(c, &ArcsRequest{Aids: []int64{10100696}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("%+v\n", reply.Arcs)
		})
		ctx.Convey("When empty", func(ctx convey.C) {
			reply, err := client.Arcs(c, &ArcsRequest{Aids: []int64{99999999999}})
			// 批量接口 err=nil arcs=nil
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(reply.Arcs, convey.ShouldBeNil)
		})
	})
}

func TestView(t *testing.T) {
	convey.Convey("TestView", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.View(c, &ViewRequest{Aid: 10100696})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("arc:%+v\n", reply.Arc)
			ctx.Printf("pages:%+v\n", reply.Pages)
		})
		ctx.Convey("When empty", func(ctx convey.C) {
			reply, err := client.View(c, &ViewRequest{Aid: 99999999999})
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(reply, convey.ShouldBeNil)
		})
	})
}

func TestViews(t *testing.T) {
	convey.Convey("TestViews", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.Views(c, &ViewsRequest{Aids: []int64{10100696}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("%+v\n", reply.Views)
		})
		ctx.Convey("When empty", func(ctx convey.C) {
			arcs, err := client.Views(c, &ViewsRequest{Aids: []int64{99999999999}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(arcs.Views, convey.ShouldBeNil)
		})
	})
}
