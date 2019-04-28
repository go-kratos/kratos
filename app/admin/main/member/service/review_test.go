package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_Reviews(t *testing.T) {
	convey.Convey("Reviews", t, func() {
		o, logs, err := s.Reviews(context.Background(), &model.ArgReviewList{Property: []int8{1}, IsDesc: true, STime: 10000, Pn: 1, Ps: 10})
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(logs, convey.ShouldNotBeNil)
	})
}
