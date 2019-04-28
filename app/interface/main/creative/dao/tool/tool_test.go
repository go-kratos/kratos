package tool

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestToolSign(t *testing.T) {
	var (
		params url.Values
	)
	convey.Convey("Sign", t, func(ctx convey.C) {
		query, err := Sign(params)
		ctx.Convey("Then err should be nil.query should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(query, convey.ShouldNotBeNil)
		})
	})
}

func TestToolDeDuplicationSlice(t *testing.T) {
	var (
		a = []int64{1, 2}
	)
	convey.Convey("DeDuplicationSlice", t, func(ctx convey.C) {
		b := DeDuplicationSlice(a)
		ctx.Convey("Then b should not be nil.", func(ctx convey.C) {
			ctx.So(b, convey.ShouldNotBeNil)
		})
	})
}

func TestToolContainAll(t *testing.T) {
	var (
		a = []int64{1}
		b = []int64{2}
	)
	convey.Convey("ContainAll", t, func(ctx convey.C) {
		p1 := ContainAll(a, b)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestToolContainAtLeastOne(t *testing.T) {
	var (
		a = []int64{}
		b = []int64{}
	)
	convey.Convey("ContainAtLeastOne", t, func(ctx convey.C) {
		p1 := ContainAtLeastOne(a, b)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestToolElementInSlice(t *testing.T) {
	var (
		a = int64(1)
		b = []int64{1, 2}
	)
	convey.Convey("ElementInSlice", t, func(ctx convey.C) {
		p1 := ElementInSlice(a, b)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestToolRandomSliceKeys(t *testing.T) {
	var (
		start = 0
		end   = 100
		count = 100
		seed  = int64(100)
	)
	convey.Convey("RandomSliceKeys", t, func(ctx convey.C) {
		p1 := RandomSliceKeys(start, end, count, seed)
		fmt.Println(p1)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(len(p1), convey.ShouldBeGreaterThan, 0)
		})
	})
}
