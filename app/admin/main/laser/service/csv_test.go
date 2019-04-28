package service

import (
	"go-common/app/admin/main/laser/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceFormatCSV(t *testing.T) {
	convey.Convey("FormatCSV", t, func(convCtx convey.C) {
		var (
			records = [][]string{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			FormatCSV(records)
			convCtx.Convey("Then err should be nil.data should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceformatAuditCargo(t *testing.T) {
	convey.Convey("formatAuditCargo", t, func(convCtx convey.C) {
		var (
			wrappers  = []*model.CargoViewWrapper{}
			lineWidth = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			formatAuditCargo(wrappers, lineWidth)
			convCtx.Convey("Then data should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceformatVideoAuditStat(t *testing.T) {
	convey.Convey("formatVideoAuditStat", t, func(convCtx convey.C) {
		var (
			statViewExts = []*model.StatViewExt{}
			lineWidth    = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			formatVideoAuditStat(statViewExts, lineWidth)
			convCtx.Convey("Then data should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}
