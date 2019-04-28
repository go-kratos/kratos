package service

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/siddontang/go-mysql/mysql"
)

//VideoList is used for getting PGC video from DB
func (s *Service) VideoList(c *bm.Context, param *model.VideoListParam) (pager *model.VideoListPager, err error) {
	var (
		order  string
		total  int
		videos []*model.VideoListQuery
	)
	selectStr := []string{
		"ugc_video.id",
		"ugc_video.aid",
		"ugc_video.cid",
		"ugc_video.eptitle",
		"ugc_video.valid",
		"ugc_video.mtime",
		"ugc_video.index_order",
		"ugc_archive.title",
		"ugc_archive.typeid",
	}
	//只筛选出未删除 且审核过的视频
	w := map[string]interface{}{"deleted": 0, "result": 1}
	db := s.DB.Model(&model.VideoListQuery{}).Where(w)
	db = db.Select(selectStr).
		Where("ugc_archive.deleted = ?", 0).
		Where("ugc_archive.result = ?", 1).
		Joins("LEFT JOIN ugc_archive ON ugc_archive.aid = ugc_video.aid")
	if param.CID != "" {
		db = db.Where("ugc_video.aid = ?", param.CID)
	}
	if param.VID != "" {
		db = db.Where("`cid` = ?", param.VID)
	}
	if param.Typeid > 0 {
		db = db.Where("typeid = ?", param.Typeid)
	}
	if param.Pid > 0 {
		db = db.Where("typeid in (?)", s.arcPTids[param.Pid])
	}
	if param.Valid != "" {
		db = db.Where("ugc_video.valid = ?", param.Valid)
	}
	if param.Order == 2 {
		order = "ugc_video.mtime DESC"
	} else {
		order = "ugc_video.mtime ASC"
	}
	if err = db.Order(order).Offset((param.Pn - 1) * param.Ps).Limit(param.Ps).Find(&videos).Error; err != nil {
		return
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	//get parent id
	for k, v := range videos {
		videos[k].PTypeID = s.GetArchivePid(v.TypeID)
	}
	pager = &model.VideoListPager{
		Items: videos,
		Page: &model.Page{
			Num:   param.Pn,
			Size:  param.Ps,
			Total: int(total),
		},
	}
	return
}

//VideoOnline is used for online PGC video
func (s *Service) VideoOnline(ids []int64) (err error) {
	w := map[string]interface{}{"deleted": 0, "result": 1}
	tx := s.DB.Begin()
	for _, v := range ids {
		video := model.Video{}
		if errDB := tx.Model(&model.Video{}).Where(w).Where("id=?", v).First(&video).Error; errDB != nil {
			err = fmt.Errorf("找不到id为%v的数据", v)
			tx.Rollback()
			return
		}
		if errDB := tx.Model(&model.Video{}).Where("id=?", v).
			Update("valid", 1).Error; errDB != nil {
			err = errDB
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

//VideoHidden is used for hidden UGC video
func (s *Service) VideoHidden(ids []int64) (err error) {
	w := map[string]interface{}{"deleted": 0, "result": 1}
	tx := s.DB.Begin()
	for _, v := range ids {
		video := model.Video{}
		if errDB := tx.Model(&model.Video{}).Where(w).Where("id=?", v).First(&video).Error; errDB != nil {
			err = fmt.Errorf("找不到id为%v的数据", v)
			tx.Rollback()
			return
		}
		if errDB := tx.Model(&model.Video{}).Where("id=?", v).
			Update("valid", 0).Error; errDB != nil {
			err = errDB
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

//VideoUpdate is used for hidden UGC video
func (s *Service) VideoUpdate(id int, title string) (err error) {
	w := map[string]interface{}{"id": id}
	if err = s.DB.Model(&model.Video{}).Where(w).Update("eptitle", title).Error; err != nil {
		return
	}
	return
}

func (s *Service) getArc(aid int64) (res *model.Archive, err error) {
	var data = model.Archive{}
	if err = s.DB.Where("aid = ?", aid).Where("deleted =0").First(&data).Error; err != nil {
		log.Error("getArc Aid %d, Err %v", aid, err)
		return
	}
	res = &data
	return
}

func (s *Service) loadAbnCids() {
	var (
		rows    *sql.Rows
		err     error
		abnCids []*model.AbnorCids
		cfg     = s.c.Cfg.Abnormal
	)
	// select cid from ugc_video where cid > 12780000 and deleted =0 and mark = 1 and submit = 1 and ctime < '2018-10-17 18:00:00' and transcoded = 0
	if rows, err = s.DB.Model(&model.Video{}).Where("cid > ?", cfg.CriticalCid).
		Where("deleted = 0").Where("mark = 1").Where("submit = 1").
		Where(fmt.Sprintf("ctime < DATE_SUB(NOW(), INTERVAL %d HOUR)", cfg.AbnormHours)).
		Where("transcoded = 0").Select("cid,aid,ctime,eptitle").Rows(); err != nil {
		log.Error("loadAbnCids Err %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			video = &model.AbnorVideo{}
			arc   *model.Archive
		)
		if err = rows.Scan(&video.CID, &video.AID, &video.CTime, &video.VideoTitle); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if arc, err = s.getArc(video.AID); err != nil {
			log.Error("getArc Aid %d, Err %v", video.AID, err)
			continue
		}
		abnCids = append(abnCids, video.ToCids(arc))
	}
	if err = rows.Err(); err != nil {
		log.Error("loadAbnCids rows Err %v", err)
		return
	}
	s.abnCids = abnCids
	log.Info("loadAbnCids Update Memory, Length %d", len(abnCids))
}

func (s *Service) loadAbnCidsproc() {
	for {
		time.Sleep(time.Duration(s.c.Cfg.Abnormal.ReloadFre))
		s.loadAbnCids()
	}
}

// AbnDebug returns the memory abnormal cids
func (s *Service) AbnDebug() (data []*model.AbnorCids) {
	return s.abnCids
}

// AbnormExport exports the abnormal cids in CSV format
func (s *Service) AbnormExport() (data *bytes.Buffer, fileName string) {
	var cfg = s.c.Cfg.Abnormal
	data = &bytes.Buffer{}
	csvWriter := csv.NewWriter(data)
	fileName = fmt.Sprintf("attachment;filename=\"Abnormal_%dh_%s.csv\"", cfg.AbnormHours, time.Now().Format(mysql.TimeFormat))
	csvWriter.Write(cfg.ExportTitles)
	for _, v := range s.abnCids {
		csvWriter.Write(v.Export())
	}
	csvWriter.Flush()
	return
}
