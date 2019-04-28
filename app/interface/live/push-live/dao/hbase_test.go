package dao

import (
	"context"
	"go-common/app/interface/live/push-live/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_fans1(t *testing.T) {
	initd()
	Convey("Parse Json To Struct", t, func() {
		fan := int64(27515316)
		f, fsp1, err := d.Fans(context.TODO(), fan, model.RelationAttention)
		t.Logf("the included(%v) includedSP(%v) err(%v)", f, fsp1, err)
		f, fsp2, err := d.Fans(context.TODO(), fan, model.RelationSpecial)
		t.Logf("the included(%v) includedSP(%v) err(%v)", f, fsp2, err)
		f, fsp3, err := d.Fans(context.TODO(), fan, model.RelationAll)
		t.Logf("the included(%v) includedSP(%v) err(%v)", f, fsp3, err)

		So(len(fsp1)+len(fsp2), ShouldEqual, len(fsp3))

	})
}

func Test_fans2(t *testing.T) {
	initd()
	Convey("Parse Json To Struct", t, func() {
		upper := int64(27515316)
		fans := make(map[int64]bool)
		fans[1232032] = true
		fans[21231134] = true
		fans[27515398] = true
		fans[27515275] = true

		f1, f2, err := d.SeparateFans(context.TODO(), upper, fans)

		t.Logf("the included(%v) includedSP(%v) err(%v)", f1, f2, err)

		So(0, ShouldEqual, 0)
	})
}
