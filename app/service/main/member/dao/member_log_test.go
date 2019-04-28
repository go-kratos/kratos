package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddMoralLogReport(t *testing.T) {
	var (
		c  = context.Background()
		ul = &model.UserLog{
			Mid: 4780461,
		}
	)
	convey.Convey("AddMoralLogReport", t, func(ctx convey.C) {
		d.AddMoralLogReport(c, ul)
		t.Logf("log:%+v", ul)
	})
}

func TestDaoaddLogReport(t *testing.T) {
	var (
		c        = context.Background()
		business = int(0)
		action   = ""
		ul       = &model.UserLog{
			Mid: 4780461,
		}
	)
	convey.Convey("addLogReport", t, func(ctx convey.C) {
		d.addLogReport(c, business, action, ul)
		t.Logf("log:%+v", ul)
	})
}
