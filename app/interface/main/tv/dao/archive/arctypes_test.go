package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveloadTypes(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("loadTypes", t, func(ctx convey.C) {
		d.loadTypes(c)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestArchiveGetPTypeName(t *testing.T) {
	var (
		typeID = int32(3)
	)
	convey.Convey("GetPTypeName", t, func(ctx convey.C) {
		firstName, secondName := d.GetPTypeName(typeID)
		ctx.Convey("Then firstName,secondName should not be nil.", func(ctx convey.C) {
			ctx.So(secondName, convey.ShouldNotBeNil)
			ctx.So(firstName, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveTargetTypes(t *testing.T) {
	convey.Convey("TargetTypes", t, func(ctx convey.C) {
		tids, err := d.TargetTypes()
		ctx.Convey("Then err should be nil.tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveFirstTypes(t *testing.T) {
	convey.Convey("FirstTypes", t, func(ctx convey.C) {
		typeMap, err := d.FirstTypes()
		ctx.Convey("Then err should be nil.typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(typeMap, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveTypeInfo(t *testing.T) {
	var (
		typeid = int32(3)
	)
	convey.Convey("TypeInfo", t, func(ctx convey.C) {
		p1, err := d.TypeInfo(typeid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveTypeChildren(t *testing.T) {
	var (
		typeid = int32(3)
	)
	convey.Convey("TypeChildren", t, func(ctx convey.C) {
		children, err := d.TypeChildren(typeid)
		ctx.Convey("Then err should be nil.children should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(children, convey.ShouldNotBeNil)
		})
	})
}
