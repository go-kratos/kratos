package service

import (
	"testing"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSeatStock(t *testing.T) {
	Convey("SeatStock", t, func() {
		req := &item.SeatStockRequest{
			Screen: 1633,
			Area:   239,
			SeatInfo: []*item.SeatPrice{
				{
					X:     1,
					Y:     1,
					Price: 13424,
				},
				{
					X:     2,
					Y:     2,
					Price: 13424,
				},
				{
					X:     3,
					Y:     3,
					Price: 13424,
				},
			},
		}
		res, err := s.SeatStock(ctx, req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
