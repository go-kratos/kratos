package data

import (
	"context"
	"testing"
	"time"

	"go-common/library/database/hbase.v2"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDatagetDataWithBackup(t *testing.T) {
	convey.Convey("getDataWithBackup", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			client        = &hbase.Client{}
			tableNameFunc func(retryCount int) string
			maxRetry      = int(0)
			key           = ""
			options       func(hrpc.Call) error
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.getDataWithBackup(c, client, tableNameFunc, maxRetry, key, options)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldBeNil)
			})
		})
	})
}

func TestDatagetTableName(t *testing.T) {
	convey.Convey("getTableName", t, func(ctx convey.C) {
		var (
			tablePrefix = ""
			date        = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := getTableName(tablePrefix, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatagenerateTableNameFunc(t *testing.T) {
	convey.Convey("generateTableNameFunc", t, func(ctx convey.C) {
		var (
			tablePrefix = ""
			date        = time.Now()
			dayDiff     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := generateTableNameFunc(tablePrefix, date, dayDiff)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
