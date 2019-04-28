package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"go-common/library/ecode"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

const _pagesize = 20

//WaterMarkist water mark list
func (s *Service) WaterMarkist(c *bm.Context, param *model.WaterMarkListParam) (pager *model.WaterMarkListPager, err error) {
	var (
		order    string
		total    int
		markList []*model.WaterMarkList
	)
	selectStr := []string{
		"tv_content.id",
		"tv_content.epid",
		"tv_content.season_id",
		"tv_content.title AS content_title",
		"tv_content.mark_time",
		"tv_ep_season.category",
		"tv_ep_season.title AS season_title",
	}
	db := s.DB.Model(&model.WaterMarkList{})
	db = db.Select(selectStr).
		Where("tv_content.is_deleted = ?", 0).
		Where("tv_content.mark = ?", model.WatermarkWhite).
		Joins("LEFT JOIN tv_ep_season ON tv_ep_season.id = tv_content.season_id")
	if param.Category != "" {
		db = db.Where("tv_ep_season.category = ?", param.Category)
	}
	if param.EpID != "" {
		db = db.Where("tv_content.epid = ?", param.EpID)
	}
	if param.SeasonID != "" {
		db = db.Where("tv_content.season_id = ?", param.SeasonID)
	}
	if param.Order == model.OrderDesc {
		order = "tv_content.mtime DESC"
	} else {
		order = "tv_content.mtime ASC"
	}
	if err = db.Order(order).Offset((param.Pn - 1) * param.Ps).Limit(param.Ps).Find(&markList).Error; err != nil {
		return
	}
	for i := range markList {
		attr := markList[i]
		attr.Category = s.pgcCatToName(atoi(attr.Category))
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	pager = &model.WaterMarkListPager{
		Items: markList,
		Page: &model.Page{
			Num:   param.Pn,
			Size:  param.Ps,
			Total: int(total),
		},
	}
	return
}

//AddEpID add water mark by ep id
func (s *Service) AddEpID(c *bm.Context, ids []int64) (res *model.AddEpIDResp, err error) {
	var (
		notExist bool
	)
	res = &model.AddEpIDResp{
		Succ:     []int64{},
		NotExist: []int64{},
		Invalids: []int64{},
	}
	for _, v := range ids {
		var mark model.WaterMarkOne
		if notExist, err = s.valueWithEpID(v, &mark); err != nil {
			return
		}
		if notExist {
			res.NotExist = append(res.NotExist, v)
			continue
		}
		if mark.Mark == model.WatermarkWhite {
			res.Invalids = append(res.Invalids, v)
			continue
		}
		if err = s.updateWithEpID(v); err != nil {
			return
		}
		res.Succ = append(res.Succ, v)
	}
	return
}

//valueWithEpID get value with epid
func (s *Service) valueWithEpID(id int64, m *model.WaterMarkOne) (exist bool, err error) {
	if err = s.DB.Where("is_deleted = 0").Where("epid = ?", id).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return
	}
	return false, nil
}

//updateWithEpID update with epid
func (s *Service) updateWithEpID(id int64) (err error) {
	up := map[string]interface{}{
		"mark":      model.WatermarkWhite,
		"mark_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err = s.DB.Model(&model.WaterMarkOne{}).Where("epid=?", id).Update(up).Error; err != nil {
		return
	}
	return
}

//AddSeasonID add water mark with season id
func (s *Service) AddSeasonID(c *bm.Context, ids []int64) (res *model.AddEpIDResp, err error) {
	var (
		notExist bool
	)
	res = &model.AddEpIDResp{
		Succ:     []int64{},
		NotExist: []int64{},
		Invalids: []int64{},
	}
	for _, v := range ids {
		var mark model.WaterMarkOne
		if notExist, err = s.valueWithSeasonID(v, &mark); err != nil {
			return
		}
		if notExist {
			res.NotExist = append(res.NotExist, v)
			continue
		}
		if err = s.updateWithSeasonID(v); err != nil {
			return
		}
		res.Succ = append(res.Succ, v)
	}
	return
}

//valueWithSeasonID get value with season id
func (s *Service) valueWithSeasonID(id int64, m *model.WaterMarkOne) (exist bool, err error) {
	w := map[string]interface{}{
		"is_deleted": 0,
		"season_id":  id,
		//"mark":       model.WatermarkDefault,
	}
	if err = s.DB.Where(w).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return
	}
	return false, nil
}

//updateWithSeasonID update with seasonID
func (s *Service) updateWithSeasonID(id int64) (err error) {
	up := map[string]interface{}{
		"mark":      model.WatermarkWhite,
		"mark_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err = s.DB.Model(&model.WaterMarkOne{}).Where("season_id=?", id).
		Update(up).Error; err != nil {
		return
	}
	return
}

//DeleteWatermark delete watermark
func (s *Service) DeleteWatermark(ids []int64) (err error) {
	if len(ids) > 50 {
		err = fmt.Errorf("更新数量最多为50条")
		return
	}
	up := map[string]interface{}{
		"mark": model.WatermarkDefault,
	}
	if err = s.DB.Model(&model.WaterMarkOne{}).Where("id in (?)", ids).
		Update(up).Error; err != nil {
		return
	}
	return
}

// titleMatch picks the title match seasons
func (s *Service) titleMatch(title string) (match []int64) {
	if len(s.snsInfo) == 0 {
		return
	}
	for _, v := range s.snsInfo {
		if strings.Contains(v.Title, title) {
			match = append(match, v.ID)
		}
	}
	return
}

// TransList picks the transcode list
func (s *Service) TransList(ctx context.Context, req *model.TransReq) (data *model.TransPager, err error) {
	var (
		db    *gorm.DB
		cntEp int
		cntSn = new(model.SnCount)
	)
	data = &model.TransPager{
		Page: &model.Page{
			Num: req.Pn, Size: _pagesize, Total: 0,
		},
		Items: make([]*model.TransReply, 0),
	}
	if db, err = s.transReqT(req); err != nil {
		return
	}
	if err = db.Count(&cntEp).Error; err != nil { //
		log.Error("countEp Err %v", err)
		return
	}
	if cntEp == 0 { // if no result ,just return
		data.CountSn = 0
		return
	}
	if (req.Pn-1)*_pagesize >= cntEp { // if page is much bigger than existing pages, just return error
		err = ecode.TvDangbeiPageNotExist
		return
	}
	if err = db.Table("tv_content").Select("COUNT(DISTINCT(season_id)) AS count").First(cntSn).Error; err != nil { // count season
		log.Error("countSn Err %v", err)
		return
	}
	db = postReqT(req, db) // order by & limit
	if data.Items, err = s.transDB(db); err != nil {
		log.Error("transDB Err %v", err)
		return
	}
	data.Page = &model.Page{
		Num:   req.Pn,
		Size:  _pagesize,
		Total: cntEp,
	}
	data.CountSn = cntSn.Count
	return
}

func (s *Service) transDB(db *gorm.DB) (items []*model.TransReply, err error) {
	var rows *sql.Rows
	rows, err = db.Select("epid, title, transcoded, apply_time, mark_time, season_id").Rows()
	if err != nil {
		log.Error("rows Err %v", err)
		return
	}
	for rows.Next() {
		var (
			cont         = &model.TransReply{}
			aTime, mTime xtime.Time
		)
		if err = rows.Scan(&cont.EpID, &cont.Etitle, &cont.Transcoded, &aTime, &mTime, &cont.SeasonID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if mTime > xtime.Time(0) { // resolve negative value issue
			cont.MarkTime = mTime.Time().Format("2006-01-02 15:04:05")
		}
		cont.ApplyTime = aTime.Time().Format("2006-01-02 15:04:05")
		if sn, ok := s.snsInfo[cont.SeasonID]; ok {
			cont.Stitle = sn.Title
			cont.Category = s.pgcCatToName(sn.Category)
		}
		items = append(items, cont)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows Err %v", err)
		return
	}
	return
}

func (s *Service) transReqT(req *model.TransReq) (db *gorm.DB, err error) {
	var (
		inSids, catSids []int64
		ok              bool
	)
	db = s.DB.Model(model.Content{}).Where("apply_time != '0000-00-00 00:00:00'").Where("is_deleted = 0")
	if req.Status != "" { // status
		switch req.Status {
		case "0":
			db = db.Where("transcoded = 0")
		case "1":
			db = db.Where("transcoded = 1")
		case "2":
			db = db.Where("transcoded = 2")
		}
	}
	if req.EpID != 0 {
		db = db.Where("epid = ?", req.EpID)
	}
	if req.SeasonID != 0 {
		db = db.Where("season_id = ?", req.SeasonID)
	}
	if req.Title != "" {
		if inSids = s.titleMatch(req.Title); len(inSids) == 0 {
			log.Warn("titleMatch %s, Empty", req.Title)
			err = ecode.NothingFound
			return
		}
	}
	if req.Category != 0 {
		if catSids, ok = s.snsCats[req.Category]; !ok || len(catSids) == 0 {
			log.Warn("snsCats %d, Empty", req.Category)
			err = ecode.NothingFound
			return
		}
		if len(inSids) > 0 { // title match have sids
			if inSids = intersect(catSids, inSids); len(inSids) == 0 {
				err = ecode.NothingFound
				return
			}
		} else {
			inSids = catSids
		}
	}
	if len(inSids) > 0 {
		db = db.Where(fmt.Sprintf("season_id IN (%s)", xstr.JoinInts(inSids)))
	}
	return
}

func postReqT(req *model.TransReq, db *gorm.DB) (newDb *gorm.DB) {
	if req.Order == 2 { // order
		newDb = db.Order("apply_time ASC")
	} else {
		newDb = db.Order("apply_time DESC")
	}
	newDb = newDb.Offset((req.Pn - 1) * _pagesize).Limit(_pagesize)
	return
}

func intersect(big, small []int64) (out []int64) {
	var amap = make(map[int64]int, len(big))
	for _, v := range big {
		amap[v] = 1
	}
	for _, v := range small {
		if _, ok := amap[v]; ok {
			out = append(out, v)
		}
	}
	return
}
