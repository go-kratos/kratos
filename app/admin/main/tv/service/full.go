package service

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

// FullImport .
func (s *Service) FullImport(c context.Context, build int) (result []*model.APKInfo, err error) {
	result, err = s.dao.FullImport(c, build)
	return
}

func (s *Service) loadSnsproc() {
	for {
		time.Sleep(time.Duration(s.c.Cfg.LoadSnFre))
		s.loadSns(context.Background())
	}
}

// loadSns loads all not deleted season Info
func (s *Service) loadSns(c context.Context) (err error) {
	var (
		rows    *sql.Rows
		sns     = make(map[int64]*model.TVEpSeason)
		sidCats = make(map[int][]int64)
	)
	if rows, err = s.DB.Model(&model.TVEpSeason{}).Where("is_deleted = 0").Select("id, title, category, state, valid, `check`").Rows(); err != nil {
		log.Error("rows Err %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cont := &model.TVEpSeason{}
		if err = rows.Scan(&cont.ID, &cont.Title, &cont.Category, &cont.State, &cont.Valid, &cont.Check); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		sns[cont.ID] = cont
		if dataSet, ok := sidCats[cont.Category]; !ok {
			sidCats[cont.Category] = []int64{cont.ID}
		} else {
			sidCats[cont.Category] = append(dataSet, cont.ID)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err %v", err)
		return
	}
	if len(sns) > 0 {
		s.snsInfo = sns
		s.snsCats = sidCats
	}
	return
}
