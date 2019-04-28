package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PubReport(t *testing.T) {
	Convey("Test_PubReport", t, WithDao(func(d *Dao) {
		err := d.PubReport(context.Background(), &model.Report{})
		So(err, ShouldBeNil)
	}))
}

func Test_PubCallback(t *testing.T) {
	Convey("Test_PubCallback", t, WithDao(func(d *Dao) {
		err := d.PubCallback(context.Background(), []*model.Callback{&model.Callback{}})
		So(err, ShouldBeNil)
	}))
}
