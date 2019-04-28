package service

import (
	"context"
	"go-common/app/service/main/vip/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceTips
func TestServiceTips(t *testing.T) {
	Convey("TestServiceTips", t, func() {
		var (
			version     int64 = 6000
			platformStr       = "ios"
		)
		arg := &model.ArgTips{
			Version:  version,
			Platform: platformStr,
			Position: int8(1),
		}
		r, err := s.Tips(context.TODO(), arg)
		t.Logf("-------------len(%d)", len(r))
		for _, v := range r {
			t.Logf("data(+%v)", v)
		}
		So(err, ShouldBeNil)
	})
}
