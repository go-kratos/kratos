package upcrmservice

import (
	"context"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceformatLog(t *testing.T) {
	convey.Convey("formatLog", t, func(ctx convey.C) {
		var (
			log upcrmmodel.SimpleCreditLogWithContent
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := formatLog(log)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceCreditLogQueryUp(t *testing.T) {
	convey.Convey("CreditLogQueryUp", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &upcrmmodel.CreditLogQueryArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.CreditLogQueryUp(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
