package archive

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpGtimeByID(t *testing.T) {
	Convey("UpGtimeByID", t, WithDao(func(d *Dao) {
		_, err := d.UpGtimeByID(context.Background(), 0, "")
		So(err, ShouldBeNil)
	}))
}

func Test_TxUpTaskByID(t *testing.T) {
	var c = context.Background()
	Convey("TxUpTaskByID", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpTaskByID(tx, 0, map[string]interface{}{"id": 0})
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxReleaseByID(t *testing.T) {
	var c = context.Background()
	Convey("TxReleaseByID", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxReleaseByID(tx, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_MulReleaseMtime(t *testing.T) {
	var c = context.Background()
	Convey("MulReleaseMtime", t, WithDao(func(d *Dao) {
		_, err := d.MulReleaseMtime(c, []int64{1, 2}, time.Now())
		So(err, ShouldBeNil)
	}))
}

func Test_GetTimeOutTask(t *testing.T) {
	var c = context.Background()
	Convey("GetTimeOutTask", t, WithDao(func(d *Dao) {
		_, err := d.GetTimeOutTask(c)
		So(err, ShouldBeNil)
	}))
}

func Test_GetRelTask(t *testing.T) {
	var c = context.Background()
	Convey("GetRelTask", t, WithDao(func(d *Dao) {
		_, _, err := d.GetRelTask(c, 0)
		So(err, ShouldBeNil)
	}))
}

func Test_TxReleaseSpecial(t *testing.T) {
	var c = context.Background()
	Convey("TxReleaseSpecial", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxReleaseSpecial(tx, time.Now(), 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}
