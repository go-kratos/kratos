package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_ThemeInfoByMid(t *testing.T) {
	convey.Convey("test theme info by mid", t, func(ctx convey.C) {
		mid := int64(545840)
		data, err := d.ThemeInfoByMid(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		for _, v := range data {
			convey.Printf("%+v", v)
		}
	})
}

func TestDao_RawTheme(t *testing.T) {
	convey.Convey("test theme", t, func(ctx convey.C) {
		mid := int64(545840)
		data, err := d.RawTheme(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		for _, v := range data.List {
			convey.Printf("%+v", v)
		}
	})
}

func TestDao_ThemeActive(t *testing.T) {
	convey.Convey("test theme active", t, func(ctx convey.C) {
		mid := int64(27515256)
		themeID := int64(11)
		err := d.ThemeActive(context.Background(), mid, themeID)
		convey.So(err, convey.ShouldBeNil)
	})
}
