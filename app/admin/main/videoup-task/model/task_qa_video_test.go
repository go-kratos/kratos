package model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var qv = &QATaskVideo{
	VideoDetail: VideoDetail{
		Attribute: 1048691,
		UPGroups:  []int64{1, 2},
	},
}

func TestQATaskVideo_GetAttributeList(t *testing.T) {
	Convey("GetAttributeqv.AttributeList", t, func() {
		qv.GetAttributeList()
		t.Logf("attributeList(%+v)", qv.AttributeList)
		So(qv.AttributeList["norank"], ShouldEqual, 1)
		So(qv.AttributeList["nosearch"], ShouldEqual, 1)
		So(qv.AttributeList["nodynamic"], ShouldEqual, 1)
		So(qv.AttributeList["norecommend"], ShouldEqual, 1)
		So(qv.AttributeList["oversea_block"], ShouldEqual, 1)
		So(qv.AttributeList["push_blog"], ShouldEqual, 1)
	})
}
