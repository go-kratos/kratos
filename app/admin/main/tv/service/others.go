package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

// Playurl new playurl function, get url from API
func (s *Service) Playurl(cid int) (playurl string, err error) {
	if playurl, err = s.dao.Playurl(ctx, cid); err != nil {
		log.Error("Playurl API Error(%d) (%v)", cid, err)
		return
	}
	if playurl, err = s.hostChange(playurl); err != nil {
		log.Error("hostChange Error(%s)-(%v)", playurl, err)
		return
	}
	log.Info("NewPlayURL cid = %d, playurl = %s", cid, playurl)
	return
}

// hostChange can change the url from playurl api to tvshenhe's host
func (s *Service) hostChange(playurl string) (replacedURL string, err error) {
	u, err := url.Parse(playurl)
	if err != nil {
		log.Error("hostChange ParseURL error (%v)", err)
		return
	}
	log.Info("[hostChange] for URL: %s, Original Host: %s, Now we change it to: %s", playurl, u.Host, s.c.Cfg.Playpath)
	u.Host = s.c.Cfg.Playpath // replace the host
	u.RawQuery = ""           // remove useless query
	replacedURL = u.String()
	return
}

// Upload can upload a file object: store the info in Redis, and transfer the file to Bfs
func (s *Service) Upload(c context.Context, fileName string, fileType string, timing int64, body []byte) (location string, err error) {
	if location, err = s.dao.Upload(c, fileName, fileType, timing, body); err != nil {
		log.Error("s.upload.Upload() error(%v)", err)
	}
	return
}

// unshelveReqT treats the unshelve request to db ( for update ) and dbSel (for select )
func (s *Service) unshelveReqT(req *model.ReqUnshelve) (db, dbSel *gorm.DB, err error) {
	if length := len(req.IDs); length == 0 || length > s.c.Cfg.AuditConsult.UnshelveNb {
		err = ecode.RequestErr
		return
	}
	switch req.Type {
	case 1: // sid
		db = s.DB.Model(&model.TVEpSeason{}).Where("is_deleted = 0").
			Where(fmt.Sprintf("id IN (%s)", xstr.JoinInts(req.IDs)))
		dbSel = db.Select("id")
	case 2: // epid
		db = s.DB.Model(&model.Content{}).Where("is_deleted = 0").
			Where(fmt.Sprintf("epid IN (%s)", xstr.JoinInts(req.IDs)))
		dbSel = db.Select("epid")
	case 3: // aid
		db = s.DB.Model(&model.Archive{}).Where("deleted = 0").
			Where(fmt.Sprintf("aid IN (%s)", xstr.JoinInts(req.IDs)))
		dbSel = db.Select("aid")
	case 4: // cid
		db = s.DB.Model(&model.Video{}).Where("deleted = 0").
			Where(fmt.Sprintf("cid IN (%s)", xstr.JoinInts(req.IDs)))
		dbSel = db.Select("cid")
	default:
		err = ecode.RequestErr
	}
	return
}

// Unshelve is to soft delete the media data
func (s *Service) Unshelve(c context.Context, req *model.ReqUnshelve, username string) (resp *model.RespUnshelve, err error) {
	var (
		rows         *sql.Rows
		existMap     = make(map[int64]int, len(req.IDs))
		db, dbSelect *gorm.DB
		updField     = make(map[string]int, 1)
	)
	log.Warn("Unshelve Req Type %d, IDs %v, Username %s", req.Type, req.IDs, username) // record user's action
	resp = &model.RespUnshelve{
		SuccIDs: make([]int64, 0),
		FailIDs: make([]int64, 0),
	}
	if db, dbSelect, err = s.unshelveReqT(req); err != nil {
		log.Error("unshelve ReqT Err %v", err)
		return
	}
	if rows, err = dbSelect.Rows(); err != nil {
		log.Error("db rows Ids %v, Err %v", req.IDs, err)
		return
	}
	defer rows.Close()
	for rows.Next() { // pick existing ids
		var sid int64
		if err = rows.Scan(&sid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		resp.SuccIDs = append(resp.SuccIDs, sid)
		existMap[sid] = 1
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err %v", err)
		return
	}
	for _, v := range req.IDs { // treat to have the non-existing ids
		if _, ok := existMap[v]; !ok {
			resp.FailIDs = append(resp.FailIDs, v)
		}
	}
	if len(resp.SuccIDs) == 0 { // there isn't any to update ids
		return
	}
	switch req.Type {
	case 1, 2: // sid, epid
		updField["is_deleted"] = 1
	case 3, 4: // aid, cid
		updField["deleted"] = 1
	}
	if err = db.Update(updField).Error; err != nil {
		log.Error("update Ids %v, err %v", req.IDs, err)
	}
	return
}

// ChlSplash gets channel's splash data
func (s *Service) ChlSplash(c context.Context, req *model.ReqChannel) (res *model.ChannelPager, err error) {
	var (
		db    = s.DB.Model(&model.Channel{}).Where("deleted!=?", _isDeleted)
		items []*model.ChannelFmt
		count int64
	)
	if req.Desc != "" {
		db = db.Where("`desc` LIKE ?", "%"+req.Desc+"%")
	}
	if req.Title != "" {
		db = db.Where("title = ?", req.Title)
	}
	db.Count(&count)
	if req.Order == model.OrderDesc {
		db = db.Order("mtime DESC")
	} else {
		db = db.Order("mtime ASC")
	}
	if err = db.Offset((req.Page - 1) * _pagesize).Limit(_pagesize).Find(&items).Error; err != nil {
		log.Error("chlList Error (%v)", err)
		return
	}
	for _, v := range items {
		v.MtimeFormat = s.TimeFormat(v.Mtime)
		v.Mtime = 0
	}
	res = &model.ChannelPager{
		TotalCount: count,
		Pn:         req.Page,
		Ps:         _pagesize,
		Items:      items,
	}
	return
}

func atoi(str string) (res int) {
	res, _ = strconv.Atoi(str)
	return
}
