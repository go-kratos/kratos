package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_RecheckByAid(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.RecheckByAid(context.Background(), 1, 10098814)
		So(err, ShouldBeNil)
	}))
}

func Test_RecheckIDByAID(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, _, err := d.RecheckIDByAID(context.Background(), 1, []int64{10098814})
		So(err, ShouldBeNil)
	}))
}

func Test_RecheckStateMap(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.RecheckStateMap(context.Background(), 1, []int64{10098814})
		So(err, ShouldBeNil)
	}))
}

func Test_TxUpRecheckState(t *testing.T) {
	var c = context.Background()
	Convey("TxUpRecheckState", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		err := d.TxUpRecheckState(tx, 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}
