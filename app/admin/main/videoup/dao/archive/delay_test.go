package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Delay(t *testing.T) {
	var (
		err error
	)
	Convey("PopMsgCache", t, WithDao(func(d *Dao) {
		_, err = d.Delay(context.Background(), 10098814, 1)
		So(err, ShouldBeNil)
	}))
}

func Test_TxUpDelay(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	Convey("TxUpDelay", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err = d.TxUpDelay(tx, 0, 0, 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpDelState(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	Convey("TxUpDelState", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err = d.TxUpDelState(tx, 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpDelayDtime(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	Convey("TxUpDelayDtime", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err = d.TxUpDelayDtime(tx, 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxDelDelay(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	Convey("TxDelDelay", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err = d.TxDelDelay(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}
