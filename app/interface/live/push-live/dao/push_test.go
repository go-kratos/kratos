package dao

import (
	"go-common/app/interface/live/push-live/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_getSign(t *testing.T) {
	initd()
	Convey("should return correct sign string by given params and secret", t, func() {
		params := map[string]string{
			"aa": "abc",
			"bb": "xyz",
			"cc": "opq",
		}
		secret := "abc"
		sign := d.getSign(params, secret)

		So(sign, ShouldEqual, "4571d284b198823bbf62f34cf38c9307")
	})
}

func TestService_GetPushTemplate(t *testing.T) {
	initd()
	Convey("should return correct template by different type", t, func() {
		name := "test"
		t1 := d.GetPushTemplate(model.AttentionGroup, name)
		t2 := d.GetPushTemplate(model.SpecialGroup, name)
		t3 := d.GetPushTemplate("test group", name)

		So(t1, ShouldEqual, "你关注的【test】正在直播~")
		So(t2, ShouldEqual, "你特别关注的【test】正在直播~")
		// default type template
		So(t3, ShouldEqual, name)
	})
}
