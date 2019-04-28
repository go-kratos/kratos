package service

import (
	"context"

	"go-common/app/job/main/search/model"
	"go-common/library/log"
)

// all all data to es
func (s *Service) all(c context.Context, appid string, writeEntityIndex bool) {
	var stat = new(model.Stat)
	app := s.base.D.AppPool[appid]
	app.InitIndex(c)
	app.InitOffset(c)
	//app.Offset(c)
	app.Sleep(c)
	for {
		start := 0
		length, err := app.AllMessages(c)
		if err != nil {
			log.Error("AllMessages error(%v)", err)
			app.Sleep(c)
			continue
		}
		for {
			end := start + _bulkSize
			diff := length - start
			if diff > _bulkSize {
				if err := app.BulkIndex(c, start, end, writeEntityIndex); err != nil {
					log.Error("es:BulkIndex error(%v)", err)
					app.Sleep(c)
					continue
				}
				start = end
			} else if diff > 0 && diff <= _bulkSize {
				if err := app.BulkIndex(c, start, length, writeEntityIndex); err != nil {
					log.Error("BulkIndex error(%v)", err)
					app.Sleep(c)
					continue
				}
				if err := app.Commit(c); err != nil {
					log.Error("UpdateOffsetID error(%v)", err)
					app.Sleep(c)
					continue
				}
				app.Sleep(c)
				break
			} else {
				app.Sleep(c)
				break
			}
		}
		stat.Counts += length
		s.updateStat(appid, stat)
		if length < app.Size(c) {
			switch appid {
			case "pgc_media", "esports", "esports_contests", "academy_archive", "esports_fav_all", "activity_all":
				app.SetRecover(c, 0, "", 0)
				app.Sleep(c)
				continue
			}
			break
		}
		app.Sleep(c)
	}
	log.Info("appid:%s, all data to es successful!!!", appid)
}
