package service

import (
	"testing"
	"time"

	"go-common/app/admin/main/coupon/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestBatchInfo
func TestInitCodes(t *testing.T) {
	Convey("TestInitCodes ", t, func() {
		err := s.InitCodes(c, "allowance_batch100")
		time.Sleep(2 * time.Second)
		So(err, ShouldBeNil)
	})
}

func TestCodeBlock(t *testing.T) {
	Convey("TestCodeBlock ", t, func() {
		err := s.CodeBlock(c, &model.ArgCouponCode{ID: 103})
		So(err, ShouldBeNil)
	})
}

func TestCodeUnBlock(t *testing.T) {
	Convey("TestCodeUnBlock ", t, func() {
		err := s.CodeUnBlock(c, &model.ArgCouponCode{ID: 103})
		So(err, ShouldBeNil)
	})
}
