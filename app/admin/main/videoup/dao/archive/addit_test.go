package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddit(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		aid := int64(10098814)
		a, err := d.Addit(context.Background(), aid)
		So(err, ShouldBeNil)
		So(a, ShouldNotBeNil)
	}))
}

func TestUpAdditRedirect(t *testing.T) {
	Convey("UpAdditRedirect", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpAdditRedirect(tx, 0, "")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestTxUpAddit(t *testing.T) {
	Convey("TxUpAddit", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpAddit(tx, 0, 0, "", "", "")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestAdditBatch(t *testing.T) {
	Convey("AdditBatch", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		_, err := d.AdditBatch(c, []int64{1, 2, 3})
		So(err, ShouldBeNil)
	}))
}
