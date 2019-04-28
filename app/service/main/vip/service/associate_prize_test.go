package service

import (
	"fmt"
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceBilibiliPrizeGrant(t *testing.T) {
	Convey(" TestServiceBilibiliPrizeGrant ", t, func() {
		res, err := s.BilibiliPrizeGrant(c, &model.ArgBilibiliPrizeGrant{})
		fmt.Println("res", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})
}

func TestServiceThirdPrizeGrant(t *testing.T) {
	Convey(" TestServiceThirdPrizeGrant ", t, func() {
		err := s.ThirdPrizeGrant(c, &model.ArgThirdPrizeGrant{})
		So(err, ShouldBeNil)

	})
}
