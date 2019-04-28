package ugc

import (
	"go-common/app/service/main/archive/api"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_FinishUpper(t *testing.T) {
	Convey("TestDao_FinishUpper", t, WithDao(func(d *Dao) {
		err := d.FinishUpper(ctx, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_Import(t *testing.T) {
	Convey("TestDao_Import", t, WithDao(func(d *Dao) {
		res, err := d.Import(ctx)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_PpUpper(t *testing.T) {
	Convey("TestDao_PpUpper", t, WithDao(func(d *Dao) {
		err := d.PpUpper(ctx, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_BeginTran(t *testing.T) {
	Convey("TestDao_BeginTran", t, WithDao(func(d *Dao) {
		tx, err := d.BeginTran(ctx)
		So(err, ShouldBeNil)
		tx.Rollback()
	}))
}

func TestDao_FilterExist(t *testing.T) {
	Convey("TestDao_FilterExist", t, WithDao(func(d *Dao) {
		pMap := make(map[int64]*api.Arc)
		pAids := []int64{10106351, 10106309, 10106308, 10106307, 10106306, 10106284, 10105807, 10101484, 10100856, 10100855, 10100486, 10100146, 10100145, 10100144,
			10100044, 10099332, 10099167, 10099150, 10099149, 10098970}
		err := d.FilterExist(ctx, &pMap, pAids)
		So(err, ShouldBeNil)
	}))
}
