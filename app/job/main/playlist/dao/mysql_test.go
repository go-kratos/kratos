package dao

import (
	"context"
	"testing"
	"time"

	plmdl "go-common/app/interface/main/playlist/model"
	"go-common/app/job/main/playlist/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Stat(t *testing.T) {
	Convey("test stat", t, WithDao(func(d *Dao) {
		pid := int64(41)
		res, err := d.Stat(context.TODO(), pid)
		So(err, ShouldBeNil)
		Printf("%+v", res)
	}))
}

func TestDao_Update(t *testing.T) {
	Convey("test update", t, WithDao(func(d *Dao) {
		var (
			pid  int64 = 1
			aid  int64 = 1
			view int64 = 2
		)
		arg := &model.StatM{Type: plmdl.PlDBusType, ID: pid, Aid: aid, Count: &view, Timestamp: xtime.Time(time.Now().Unix()), IP: ""}
		res, err := d.Update(context.TODO(), arg, model.ViewCountType)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
