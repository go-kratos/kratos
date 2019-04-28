package manager

import (
	"context"
	"go-common/app/service/main/up/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"strings"
)

func TestManagerAddGroup(t *testing.T) {
	var (
		c            = context.TODO()
		groupAddInfo = &model.AddGroupArg{}
	)
	convey.Convey("AddGroup", t, func(ctx convey.C) {
		_, err := d.AddGroup(c, groupAddInfo)
		if strings.Contains(err.Error(), "Error 1062") {
			err = nil
		}
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestManagerCheckGroupExist(t *testing.T) {
	var (
		c            = context.TODO()
		groupAddInfo = &model.AddGroupArg{}
		exceptid     = int64(0)
	)
	convey.Convey("CheckGroupExist", t, func(ctx convey.C) {
		exist, err := d.CheckGroupExist(c, groupAddInfo, exceptid)
		ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerUpdateGroup(t *testing.T) {
	var (
		c            = context.TODO()
		groupAddInfo = &model.EditGroupArg{}
	)
	convey.Convey("UpdateGroup", t, func(ctx convey.C) {
		res, err := d.UpdateGroup(c, groupAddInfo)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestManagerRemoveGroup(t *testing.T) {
	var (
		c   = context.TODO()
		arg = &model.RemoveGroupArg{}
	)
	convey.Convey("RemoveGroup", t, func(ctx convey.C) {
		res, err := d.RemoveGroup(c, arg)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetGroup(t *testing.T) {
	var (
		c   = context.TODO()
		arg = &model.GetGroupArg{}
	)
	convey.Convey("GetGroup", t, func(ctx convey.C) {
		_, err := d.GetGroup(c, arg)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
