package service

import (
	"context"

	"go-common/app/job/main/figure/model"
)

const (
	_reportDel = "report_del"
)

// DanmakuReport .
func (s *Service) DanmakuReport(c context.Context, d *model.DMAction) (err error) {
	if err = s.figureDao.DanmakuReport(c, d.Data.OwnerUID, model.ACColumnPublishDanmakuDeleted, -1); err != nil {
		return
	}
	if err = s.figureDao.DanmakuReport(c, d.Data.ReportUID, model.ACColumnDanmakuReoprtPassed, 1); err != nil {
		return
	}
	return
}
