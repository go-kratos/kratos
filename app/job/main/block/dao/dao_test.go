package dao

import (
	"context"
	"flag"
	"testing"

	"go-common/app/job/main/block/conf"
	"go-common/app/job/main/block/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/block-service-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao = New()
	defer dao.Close()
	m.Run()
}

func TestDB(t *testing.T) {
	Convey("db", t, func() {
		tx, err := dao.BeginTX(ctx)
		So(err, ShouldBeNil)

		var (
			mid int64 = 46333
		)
		err = dao.UpsertAddBlockCount(ctx, mid)
		So(err, ShouldBeNil)

		err = dao.TxUpsertUser(ctx, tx, mid, model.BlockStatusFalse)
		So(err, ShouldBeNil)

		var (
			history = &model.DBHistory{
				MID:     mid,
				Source:  model.BlockSourceRemove,
				Comment: "ut test",
				Action:  model.BlockActionAdminRemove,
			}
		)
		err = dao.TxInsertHistory(ctx, tx, history)
		So(err, ShouldBeNil)

		err = tx.Rollback()
		So(err, ShouldBeNil)
	})
}

func TestTool(t *testing.T) {
	Convey("tool", t, func() {
		var (
			mids = []int64{1, 2, 3, 46333, 35858}
		)
		str := midsToParam(mids)
		So(str, ShouldEqual, "1,2,3,46333,35858")
	})
}
