package ugc

import (
	"fmt"
	"go-common/library/database/sql"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_DeletedArc(t *testing.T) {
	Convey("TestDao_DeletedArc", t, WithDao(func(d *Dao) {
		res, err := d.DeletedArc(ctx)
		if err == sql.ErrNoRows {
			fmt.Println("No to delete data")
			return
		}
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
}

func TestDao_PpDelArc(t *testing.T) {
	Convey("TestDao_PpDelArc", t, WithDao(func(d *Dao) {
		err := d.PpDelArc(ctx, 333)
		So(err, ShouldBeNil)
	}))
}

func TestDao_FinishDelArc(t *testing.T) {
	Convey("TestDao_FinishDelArc", t, WithDao(func(d *Dao) {
		err := d.FinishDelArc(ctx, 333)
		So(err, ShouldBeNil)
	}))
}
