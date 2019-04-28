package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// AddTagInfo add tag info service
func (s *Service) AddTagInfo(c context.Context, t *model.TagInfo, creator string) (err error) {
	t.Creator = creator
	if t.Dimension == 0 && (t.Category == 0 || t.Business == 0) {
		err = fmt.Errorf("维度为稿件时必须填写业务和分区")
		return
	}
	// bgm
	if t.Business == 3 {
		t.Category = 0
	}

	if t.Icon != "" && t.ActivityID == 0 {
		err = fmt.Errorf("请填写活动id")
		return
	}

	_, err = s.dao.GetTagInfoByName(c, t.Tag, t.Dimension, t.Category, t.Business)
	if err == sql.ErrNoRows {
		_, err = s.dao.InsertTag(c, t)
		if err != nil {
			log.Error("s.dao.InsertTag error(%v)", err)
		}
		return
	}
	if err != nil {
		log.Error("s.dao.GetTagInfoByName error(%v)", err)
		return
	}
	// tag has exist, can not add
	return ecode.GrowupTagAddForbit
}

// UpdateTagInfo update tag
func (s *Service) UpdateTagInfo(c context.Context, t *model.TagInfo, creator string) (err error) {
	t.Creator = creator
	if t.Dimension == 0 && (t.Category == 0 || t.Business == 0) {
		err = fmt.Errorf("维度为稿件时必须填写业务和分区")
		return
	}
	// bgm
	if t.Business == 3 {
		t.Category = 0
	}

	if t.Icon != "" && t.ActivityID == 0 {
		err = fmt.Errorf("请填写活动id")
		return
	}
	old, err := s.dao.GetTagInfo(c, int(t.ID))
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	if old.StartTime.Time().Unix() <= today && old.EndTime.Time().Unix() >= today {
		err = fmt.Errorf("标签生效中，不能修改")
		return
	}
	_, err = s.dao.UpdateTagInfo(c, t)
	if err != nil {
		log.Error("s.dao.UpdateTagInfo error(%v)", err)
	}
	return
}

// ModeTagState mode tag state
func (s *Service) ModeTagState(c context.Context, tagID int, isDeleted int) (err error) {
	rows, err := s.dao.UpdateTagState(c, tagID, isDeleted)
	if err != nil {
		log.Error("s.dao.UpdateTagState error(%v)", err)
		return
	}
	if rows == 0 {
		return fmt.Errorf("Modification has not taken effect")
	}
	return
}

// AddTagUps add tag ups service
func (s *Service) AddTagUps(c context.Context, tagID int, mids []int64, isCommon int) (err error) {
	info, err := s.dao.GetTagInfo(c, tagID)
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}
	_, err = s.dao.UpdateTagCom(c, tagID, isCommon)
	if err != nil {
		log.Error("s.dao.UpdateTagCom error(%v)", err)
		return
	}

	if isCommon == 0 && len(mids) > 0 {
		if err = s.InsertTagUpInfos(c, info.ID, mids); err != nil {
			log.Error("s.InsertTagUpInfos error(%v)", err)
		}
	}
	return
}

// InsertTagUpInfos insert tag_info
func (s *Service) InsertTagUpInfos(c context.Context, tagID int64, mids []int64) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	for _, mid := range mids {
		_, err = s.dao.TxInsertTagUpInfo(tx, tagID, mid, 0)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.TxInsertTagUpInfo error(%v)", err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
	}
	return
}

// ReleaseUp release tag ups service
func (s *Service) ReleaseUp(c context.Context, tagID int, mid int64) (err error) {
	info, err := s.dao.GetTagInfo(c, tagID)
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}
	if info.IsCommon == 1 {
		return ecode.GrowupTagForbit
	}
	_, err = s.dao.InsertTagUpInfo(c, info.ID, mid, 1)
	if err != nil {
		log.Error("s.dao.InsertTagUpInfo error(%v)", err)
	}
	return
}

func tagQueryStmt(startTime int64, endTime int64, categories []int64, business []int64, tag string, effect int, sort string) (query string) {
	if effect == 1 {
		startTime = time.Now().Unix()
		endTime = time.Now().Unix()
		query += "is_deleted = 0"
		query += " AND "
	}
	if startTime != 0 && endTime != 0 {
		st := time.Unix(startTime, 0).Format("2006-01-02")
		et := time.Unix(endTime, 0).Format("2006-01-02")
		query += "("
		query += "(start_at <= '" + st
		query += "' AND "
		query += "end_at >= '" + et
		query += "') OR "
		query += "(start_at BETWEEN '" + st + "' AND '" + et
		query += "' OR "
		query += "end_at BETWEEN '" + st + "' AND '" + et
		query += "') OR "
		query += "(start_at >= '" + st
		query += "' AND "
		query += "end_at <= '" + et
		query += "'))"
		query += " AND "
	}
	if tag != "" {
		query += fmt.Sprintf("tag = \"%s\"", tag)
		query += " AND "
	}
	query = strings.TrimSuffix(query, " AND ")

	queryArchive := ""
	if len(categories) != 0 {
		queryArchive += fmt.Sprintf("category_id IN (%s)", xstr.JoinInts(categories))
		queryArchive += " AND "
	}
	queryUp := "dimension = 1"
	if len(business) != 0 {
		queryArchive += fmt.Sprintf("business_id IN (%s)", xstr.JoinInts(business))
		queryUp += fmt.Sprintf(" AND business_id IN (%s)", xstr.JoinInts(business))
	}
	if queryArchive != "" {
		queryArchive = strings.TrimSuffix(queryArchive, " AND ")
		query = fmt.Sprintf("%s AND ((%s AND dimension = 0) OR (%s))", query, queryArchive, queryUp)
	}

	query = strings.TrimPrefix(query, " AND ")
	if query != "" {
		query = fmt.Sprintf("WHERE %s", query)
	}

	if strings.HasPrefix(sort, "-") {
		sort = strings.TrimPrefix(sort, "-")
		sort += " DESC"
	}
	if sort != "" {
		query += " ORDER BY "
		query += sort
	}
	return
}

// QueryTagInfo query tag info service
func (s *Service) QueryTagInfo(c context.Context, start int64, end int64, categories []int64, business []int64, tag string, effect, from int, limit int, sort string) (total int, tagInfos []*model.TagInfo, err error) {
	if len(categories) == 0 && len(business) == 0 {
		return
	}
	for _, b := range business {
		if b == 3 {
			categories = append(categories, 0)
			break
		}
	}
	query := tagQueryStmt(start, end, categories, business, tag, effect, sort)
	total, err = s.dao.TagsCount(c, query)
	if err != nil {
		log.Error("s.dao.TagsCount error(%v)", err)
		return
	}
	if total == 0 {
		tagInfos = []*model.TagInfo{}
		return
	}

	tagInfos, err = s.dao.GetTagInfos(c, query, from, limit)
	if err != nil {
		log.Error("dao.GetTagInfos error(%v)", err)
		return
	}
	for _, tag := range tagInfos {
		tag.RetRatio = float32(tag.Ratio) / float32(100)
	}
	return
}

func (s *Service) getUpTagIncomeByDate(c context.Context, startTime, endTime time.Time, tagID int64, query string) (avs []*model.UpTagIncome, err error) {
	endTime = endTime.AddDate(0, 0, 1)
	avs = make([]*model.UpTagIncome, 0)
	for startTime.Before(endTime) {
		var av []*model.UpTagIncome
		date := startTime.Format("2006-01-02")
		av, err = s.dao.GetUpTagIncome(c, date, tagID, query)
		if err != nil {
			return
		}
		avs = append(avs, av...)
		startTime = startTime.AddDate(0, 0, 1)
	}
	return
}

// ListUps list ups
func (s *Service) ListUps(c context.Context, tagID int, mid int64, from int, limit int) (total int, data []*model.UpIncomeInfo, err error) {
	tagInfo, err := s.dao.GetTagInfo(c, tagID)
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}
	query := ""
	if mid != 0 {
		query = fmt.Sprintf("mid = %d", mid)
	}

	avs, err := s.getUpTagIncomeByDate(c, tagInfo.StartTime.Time(), tagInfo.EndTime.Time(), tagInfo.ID, query)
	if err != nil {
		log.Error("s.getUpTagIncomeByDate error(%v)", err)
		return
	}

	releaseUps, err := s.dao.GetTagUpInfoMID(c, tagInfo.ID, 1)
	if err != nil {
		log.Error("s.dao.GetTagUpInfoUps error(%v)", err)
		return
	}

	upIncomes := statisTagUps(c, avs, releaseUps)
	total = len(upIncomes)
	if from > len(upIncomes) {
		return
	}
	if from+limit > len(upIncomes) {
		limit = len(upIncomes) - from
	}
	data = upIncomes[from : from+limit]
	for _, d := range data {
		d.Nickname, err = s.dao.GetNickname(c, d.MID)
		if err != nil {
			log.Error("s.dao.GetNickname error(%v)", err)
			return
		}
	}
	return
}

func statisTagUps(c context.Context, avs []*model.UpTagIncome, releaseUps map[int64]int) (ups []*model.UpIncomeInfo) {
	upsMap := make(map[int64]*model.UpIncomeInfo)
	for _, av := range avs {
		if av.IsDeleted == 1 {
			continue
		}
		if _, ok := upsMap[av.MID]; ok {
			upsMap[av.MID].BaseIncome += av.BaseIncome
			upsMap[av.MID].TotalIncome += av.Income
			upsMap[av.MID].AdjustIncome += av.Income - av.BaseIncome
		} else {
			upsMap[av.MID] = &model.UpIncomeInfo{
				MID:          av.MID,
				CreateTime:   xtime.Time(av.Date.Unix()),
				BaseIncome:   av.BaseIncome,
				TotalIncome:  av.Income,
				AdjustIncome: av.Income - av.BaseIncome,
				IsDeleted:    releaseUps[av.MID],
			}
		}
	}
	ups = make([]*model.UpIncomeInfo, 0)
	for _, up := range upsMap {
		ups = append(ups, up)
	}
	sort.Slice(ups, func(i, j int) bool {
		return ups[i].TotalIncome > ups[j].TotalIncome
	})
	return
}

// ListAvs list avs
func (s *Service) ListAvs(c context.Context, tagID, from, limit int, avID int64) (total int, data []*model.AvIncomeInfo, err error) {
	data = make([]*model.AvIncomeInfo, 0)
	tagInfo, err := s.dao.GetTagInfo(c, tagID)
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}
	if tagInfo.Dimension == 1 {
		return
	}
	query := ""
	if avID != 0 {
		query = fmt.Sprintf("av_id = %d", avID)
	}

	avs, err := s.getUpTagIncomeByDate(c, tagInfo.StartTime.Time(), tagInfo.EndTime.Time(), tagInfo.ID, query)
	if err != nil {
		log.Error("s.getUpTagIncomeByDate error(%v)", err)
		return
	}

	avIncomes := statisTagAvs(c, avs)
	total = len(avIncomes)
	if from > len(avIncomes) {
		return
	}
	if from+limit > len(avIncomes) {
		limit = len(avIncomes) - from
	}
	data = avIncomes[from : from+limit]
	for _, d := range data {
		d.Category = tagInfo.Category
		d.Nickname, err = s.dao.GetNickname(c, d.MID)
		if err != nil {
			log.Error("s.dao.GetNickname error(%v)", err)
			return
		}
	}
	return
}

func statisTagAvs(c context.Context, avs []*model.UpTagIncome) (avIncome []*model.AvIncomeInfo) {
	avsMap := make(map[int64]*model.AvIncomeInfo)
	for _, av := range avs {
		if av.IsDeleted == 1 {
			continue
		}

		if _, ok := avsMap[av.AvID]; ok {
			avsMap[av.AvID].BaseIncome += av.BaseIncome
			avsMap[av.AvID].AdjustIncome += av.Income - av.BaseIncome
			avsMap[av.AvID].TotalIncome += av.Income
		} else {
			avsMap[av.AvID] = &model.AvIncomeInfo{
				AVID:         av.AvID,
				MID:          av.MID,
				CreateTime:   xtime.Time(av.Date.Unix()),
				BaseIncome:   av.BaseIncome,
				TotalIncome:  av.Income,
				AdjustIncome: av.Income - av.BaseIncome,
			}

		}
	}
	avIncome = make([]*model.AvIncomeInfo, 0)
	for _, av := range avsMap {
		avIncome = append(avIncome, av)
	}
	sort.Slice(avIncome, func(i, j int) bool {
		return avIncome[i].TotalIncome > avIncome[j].TotalIncome
	})
	return
}

// TagDetails query tag details.
func (s *Service) TagDetails(c context.Context, tagID, from, limit int) (total int, data []*model.Details, err error) {
	tagInfo, err := s.dao.GetTagInfo(c, tagID)
	if err != nil {
		log.Error("s.dao.GetTagInfo error(%v)", err)
		return
	}

	avs, err := s.getUpTagIncomeByDate(c, tagInfo.StartTime.Time(), tagInfo.EndTime.Time(), tagInfo.ID, "")
	if err != nil {
		log.Error("s.getUpTagIncomeByDate error(%v)", err)
		return
	}

	tags := statisTags(c, avs)
	total = len(tags)
	if from > len(tags) {
		return
	}
	if from+limit > len(tags) {
		limit = len(tags) - from
	}
	data = tags[from : from+limit]
	return
}

func statisTags(c context.Context, avs []*model.UpTagIncome) (tags []*model.Details) {
	var (
		tagsMap                 = make(map[string]*model.Details)
		dateMID                 = make(map[string]map[int64]struct{})
		avMap                   = make(map[int64]struct{})
		upMap                   = make(map[int64]struct{})
		totalIncome, baseIncome int
	)
	for _, av := range avs {
		key := av.Date.Format("2006-01-02")
		if av.IsDeleted == 1 {
			continue
		}
		if _, ok := tagsMap[key]; ok {
			tagsMap[key].Income += av.Income
			tagsMap[key].BaseIncome += av.BaseIncome
			tagsMap[key].AdjustIncome += av.Income - av.BaseIncome
			tagsMap[key].AvCnt++
			if _, ok := dateMID[key][av.MID]; !ok {
				tagsMap[key].UpCnt++
			}
		} else {
			tagsMap[key] = &model.Details{
				Date:   key,
				Income: av.Income,
				UpCnt:  1,
				AvCnt:  1,
			}
			dateMID[key] = make(map[int64]struct{})
		}
		dateMID[key][av.MID] = struct{}{}

		// all
		if av.AvID != 0 {
			avMap[av.AvID] = struct{}{}
		}
		upMap[av.MID] = struct{}{}
		totalIncome += av.Income
		baseIncome += av.BaseIncome
	}

	tags = make([]*model.Details, 1)
	tags[0] = &model.Details{
		Date:         "累计",
		UpCnt:        len(upMap),
		AvCnt:        len(avMap),
		Income:       totalIncome,
		BaseIncome:   baseIncome,
		AdjustIncome: totalIncome - baseIncome,
	}
	for _, t := range tagsMap {
		tags = append(tags, t)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date > tags[j].Date
	})
	return
}

// UpdateTagActivity update tag activity_id
func (s *Service) UpdateTagActivity(c context.Context, tagID, activityID int64) (err error) {
	rows, err := s.dao.UpdateTagActivity(c, tagID, activityID)
	if err == nil && rows != 1 {
		err = fmt.Errorf("UpdateActivity effect rows(%d) error ", rows)
	}
	return
}
