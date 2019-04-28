package service

import (
	"context"
	"testing"

	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.TODO()
)

func Test_Frozen(t *testing.T) {
	Convey("test user frozen", t, func() {
		la := []uint32{1, 2, 3, 4, 5}
		for _, i := range la {
			s.dao.AddLogginIP(ctx, 7593666, i)
		}
		s.Frozen(ctx, &model.LoginLog{Mid: 7593666, IP: 234566})
	})
}

func Test_UnFrozen(t *testing.T) {
	Test_Frozen(t)
	Convey("test user unfrozen", t, func() {
		s.unFrozenJob()
	})
}
