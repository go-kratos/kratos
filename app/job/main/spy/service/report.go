package service

import (
	"context"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/model"
	"go-common/library/log"
)

// AddReport add daill report.
func (s *Service) AddReport(c context.Context) {
	var (
		scount      int64
		pcount      int64
		dateVersion string
		err         error
	)
	year, month, day := time.Now().Date()
	stoday := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	etoday := time.Date(year, month, day, 23, 59, 59, 999, time.Local)
	syesday := stoday.AddDate(0, 0, -1)
	eyesday := etoday.AddDate(0, 0, -1)
	dateVersion = syesday.Format("20060102")
	if pcount, err = s.dao.PunishmentCount(c, syesday, eyesday); err != nil {
		log.Error("s.dao.PunishmentCount(%s, %s), err(%v)", syesday, eyesday, err)
		return
	}
	s.dao.AddReport(c, &model.Report{
		Name:        model.BlockCount,
		DateVersion: dateVersion,
		Val:         pcount,
		Ctime:       time.Now(),
	})
	for i := int64(0); i < conf.Conf.Property.HistoryShard; i++ {
		var count int64
		if count, err = s.dao.SecurityLoginCount(c, i, "导入二次验证,恢复行为得分", syesday, eyesday); err != nil {
			log.Error("s.dao.SecurityLoginCount(%s, %s), err(%v)", syesday, eyesday, err)
			return
		}
		scount = scount + count
		time.Sleep(s.blockTick)
	}
	s.dao.AddReport(c, &model.Report{
		Name:        model.SecurityLoginCount,
		DateVersion: dateVersion,
		Val:         scount,
		Ctime:       time.Now(),
	})
}
