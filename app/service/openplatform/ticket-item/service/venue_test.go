package service

import (
	"testing"

	"go-common/app/service/openplatform/ticket-item/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVenueSearch(t *testing.T) {
	Convey("VenueSearch", t, func() {
		param := &model.VenueSearchParam{
			ID: 1,
			PageParam: model.PageParam{
				Pn: 1,
				Ps: 10,
			},
		}
		res, err := s.VenueSearch(ctx, param)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
