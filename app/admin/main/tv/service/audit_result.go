package service

import (
	"net/url"

	"context"
	"fmt"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

// order const
const (
	newOrder = 1
)

// EpResult gives the result of ep audit
func (s *Service) EpResult(req url.Values, page int, order int) (pager *model.EPResultPager, err error) {
	var (
		count int64
		size  = s.c.Cfg.AuditRSize
		items []*model.EPResDB

		db = s.DB.Model(&model.Content{}).
			Where("tv_content.is_deleted=?", 0).
			Joins("LEFT JOIN tv_ep_season ON tv_content.season_id=tv_ep_season.id").
			Select("tv_content.*,tv_ep_season.title as stitle,tv_ep_season.category")
	)
	// order treatment
	if order == newOrder {
		db = db.Order("tv_content.inject_time DESC")
	} else {
		db = db.Order("tv_content.inject_time ASC")
	}
	// category treatment
	if category := req.Get("category"); category != "" {
		db = db.Where("tv_ep_season.category=?", category)
	}
	// audit status treatment
	if state := req.Get("state"); state != "" {
		switch state {
		case "1": // passed
			db = db.Where("tv_content.`state` = ?", 3)
		case "2": // reject
			db = db.Where("tv_content.`state` = ?", 4)
		default: // waiting result
			db = db.Where("tv_content.`state` NOT IN (3,4)")
		}
	}
	// season_id treatment
	if sid := req.Get("season_id"); sid != "" {
		db = db.Where("tv_content.season_id=?", sid)
	}
	// epid treatment
	if epid := req.Get("epid"); epid != "" {
		db = db.Where("tv_content.epid=?", epid)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Error("Count Err %v", err)
		return
	}
	pager = &model.EPResultPager{
		Page: &model.Page{
			Num:   page,
			Size:  size,
			Total: int(count),
		},
	}
	if err = db.Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return
	}
	// use time in string format to replace the time in number format
	for _, v := range items {
		pager.Items = append(pager.Items, v.ToItem())
	}
	return
}

//TimeFormat is used for format time
func (s *Service) TimeFormat(time time.Time) (format string) {
	if time < 0 {
		return ""
	}
	return time.Time().Format("2006-01-02 15:04:05")
}

// SeasonResult gives the result of ep audit
func (s *Service) SeasonResult(req url.Values, page int, order int) (pager *model.SeasonResultPager, err error) {
	var (
		count   int64
		size    = s.c.Cfg.AuditRSize
		dbTerms []*model.SeasonResDB
		db      = s.DB.Model(&model.TVEpSeason{}).Where("is_deleted=?", 0)
	)
	// order treatment
	if order == newOrder {
		db = db.Order("inject_time DESC")
	} else {
		db = db.Order("inject_time ASC")
	}
	// category treatment
	if category := req.Get("category"); category != "" {
		db = db.Where("tv_ep_season.category=?", category)
	}
	// audit status treatment
	if state := req.Get("check"); state != "" {
		switch state {
		case "1": // passed
			db = db.Where("tv_ep_season.`check` = ?", 1)
		case "2": // reject
			db = db.Where("tv_ep_season.`check` = ?", 0)
		default: // waiting result
			db = db.Where("tv_ep_season.`check` NOT IN (0,1)")
		}
	}
	// season_id treatment
	if sid := req.Get("season_id"); sid != "" {
		db = db.Where("id=?", sid)
	}
	// title treatment
	if title := req.Get("title"); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if err = db.Count(&count).Error; err != nil {
		log.Error("db Count Err %v", err)
		return
	}
	pager = &model.SeasonResultPager{
		Page: &model.Page{
			Num:   page,
			Size:  size,
			Total: int(count),
		},
	}
	if err = db.Offset((page - 1) * size).Limit(size).Find(&dbTerms).Error; err != nil {
		return
	}
	for _, v := range dbTerms {
		pager.Items = append(pager.Items, v.ToItem())
	}
	return
}

func (s *Service) typeidsTreat(secondCat int32, firstCat int32) (typeids []int32) {
	if secondCat != 0 { // typeid logic
		typeids = []int32{secondCat}
	} else if firstCat != 0 {
		if secondCats, ok := s.arcPTids[firstCat]; ok && len(secondCats) > 0 {
			typeids = secondCats
		}
	}
	return
}

// ArcResult picks archive result data
func (s *Service) ArcResult(c context.Context, req *model.ReqArcCons) (data *model.ArcResPager, err error) {
	if data, err = s.arcByES(req, s.typeidsTreat(req.SecondCat, req.FirstCat)); err != nil {
		log.Error("arcByEs Err %v", err)
		return
	}
	return
}

func (s *Service) arcByES(req *model.ReqArcCons, typeids []int32) (data *model.ArcResPager, err error) {
	var (
		esRes   *model.EsUgcResult
		aids    []int64
		arcs    []*model.Archive
		arcsMap = make(map[int64]*model.ArcRes)
		reqES   = new(model.ReqArcES)
	)
	reqES.FromAuditConsult(req, typeids)
	if esRes, err = s.dao.ArcES(ctx, reqES); err != nil {
		log.Error("UgcConsult Err %v", err)
		return
	}
	data = &model.ArcResPager{
		Page: esRes.Page,
	}
	if len(esRes.Result) == 0 {
		return
	}
	for _, v := range esRes.Result {
		aids = append(aids, v.AID)
	}
	if err = s.DB.Where(fmt.Sprintf("aid IN (%s)", xstr.JoinInts(aids))).Find(&arcs).Error; err != nil {
		log.Error("arcByES DB Aids %v Err %v", aids, err)
		return
	}
	if len(arcs) == 0 {
		return
	}
	for _, v := range arcs {
		arcsMap[v.AID] = v.ConsultRes(s.ArcTypes)
	}
	for _, v := range aids {
		if arc, ok := arcsMap[v]; ok {
			data.Items = append(data.Items, arc)
		}
	}
	return
}

// VideoResult picks video audit consult result
func (s *Service) VideoResult(c context.Context, req *model.ReqVideoCons) (data *model.VideoResPager, err error) {
	var (
		videos []*model.Video
	)
	data = &model.VideoResPager{
		Page: &model.Page{
			Num:  req.Pn,
			Size: _pagesize,
		},
	}
	db := s.DB.Model(&model.Video{}).Where("aid = ?", req.AVID).Where("deleted = 0")
	if req.Status != "" {
		db = db.Where("result = ?", req.Status)
	}
	if req.Title != "" {
		db = db.Where("eptitle LIKE ?", "%"+req.Title+"%")
	}
	if req.CID != 0 {
		db = db.Where("cid = ?", req.CID)
	}
	if err = db.Count(&data.Page.Total).Error; err != nil {
		log.Error("VideoResult Count Err %v", err)
		return
	}
	if req.Order != 1 {
		db = db.Order("index_order ASC")
	} else {
		db = db.Order("index_order DESC")
	}
	if err = db.Offset((req.Pn - 1) * _pagesize).Limit(_pagesize).Find(&videos).Error; err != nil {
		log.Error("arcByDB Err %v", err)
	}
	if len(videos) == 0 {
		return
	}
	for _, v := range videos {
		data.Items = append(data.Items, v.ConsultRes())
	}
	return
}
