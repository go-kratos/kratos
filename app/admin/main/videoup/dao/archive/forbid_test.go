package archive

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Forbid(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		af, err := d.Forbid(context.Background(), 10098814)
		So(err, ShouldBeNil)
		So(af, ShouldNotBeNil)
	}))
}

func Test_TxUpForbid(t *testing.T) {
	var c = context.Background()
	Convey("TxUpForbid", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpForbid(tx, &archive.ForbidAttr{})
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpFlowID(t *testing.T) {
	var c = context.Background()
	Convey("TxUpFlowID", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpFlowID(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}
