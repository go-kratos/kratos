package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/web/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Coins(t *testing.T) {
	Convey("test coin Coins", t, WithService(func(s *Service) {
		var (
			mid int64 = 27515256
			aid int64 = 37515257
		)
		res, err := s.Coins(context.Background(), mid, aid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_AddCoin(t *testing.T) {
	Convey("test coin AddCoin", t, WithService(func(s *Service) {
		var (
			mid        int64 = 27515256
			aid        int64 = 37515257
			upID       int64 = 37515257
			multiply   int64 = 1
			avtype     int64 = 1
			business         = model.CoinArcBusiness
			selectLike       = 1
		)
		like, err := s.AddCoin(context.Background(), aid, mid, upID, multiply, avtype, business, "", "", "", time.Now(), selectLike)
		So(err, ShouldBeNil)
		So(like, ShouldNotBeNil)
	}))
}

func TestService_CoinExp(t *testing.T) {
	Convey("test coin CoinExp", t, WithService(func(s *Service) {
		var (
			mid int64 = 27515256
		)
		res, err := s.CoinExp(context.Background(), mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
