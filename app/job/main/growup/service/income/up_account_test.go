package income

import (
	"context"
	"testing"

	model "go-common/app/job/main/growup/model/income"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpAccount1(t *testing.T) {
	Convey("UpAccount", t, func() {
		_, err := s.income.upAccountSvr.UpAccount(context.Background(), 10)
		So(err, ShouldBeNil)
	})
}

func Test_BatchInsertUpAccount(t *testing.T) {
	Convey("BatchInsertUpAccount", t, func() {
		err := s.income.upAccountSvr.BatchInsertUpAccount(context.Background(), map[int64]*model.UpAccount{})
		So(err, ShouldBeNil)
	})
}

func Test_UpdateUpAccount(t *testing.T) {
	Convey("UpdateUpAccount", t, func() {
		err := s.income.upAccountSvr.UpdateUpAccount(context.Background(), map[int64]*model.UpAccount{})
		So(err, ShouldBeNil)
	})
}
