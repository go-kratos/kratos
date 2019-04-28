package client

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMask(t *testing.T) {
	var (
		cid int64 = 632
	)
	Convey("test mask", t, func() {
		arg := &model.ArgMask{Cid: cid, Plat: 0}
		res, err := svr.Mask(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("===============%+v", res)
	})
}
