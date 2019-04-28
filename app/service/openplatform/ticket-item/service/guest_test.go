package service

import (
	"testing"

	"go-common/app/service/openplatform/ticket-item/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGuestSearch(t *testing.T) {
	Convey("GuestSearch", t, func() {
		param := &model.GuestSearchParam{
			Keyword: "1",
			Ps:      10,
			Pn:      1,
		}
		res, err := s.GuestSearch(ctx, param)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
