package model

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_attribute(t *testing.T) {
	list := AttributeList(1048691)
	convey.Convey("属性列表", t, func() {
		convey.So(list["norank"], convey.ShouldEqual, 1)
		convey.So(list["nosearch"], convey.ShouldEqual, 1)
		convey.So(list["nodynamic"], convey.ShouldEqual, 1)
		convey.So(list["norecommend"], convey.ShouldEqual, 1)
		convey.So(list["oversea_block"], convey.ShouldEqual, 1)
		convey.So(list["push_blog"], convey.ShouldEqual, 1)
	})
}
