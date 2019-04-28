package archive

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PorderConfig(t *testing.T) {
	Convey("PorderConfig", t, WithDao(func(d *Dao) {
		p, err := d.PorderConfig(context.TODO())
		So(err, ShouldBeNil)
		So(p, ShouldNotBeNil)
	}))
}

func Test_Porder(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		p, err := d.Porder(context.Background(), 10098814)
		So(err, ShouldBeNil)
		So(p, ShouldNotBeNil)
	}))
}

func Test_TxUpPorder(t *testing.T) {
	var c = context.Background()
	Convey("TxUpPorder", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpPorder(tx, 0, &archive.ArcParam{})
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}
