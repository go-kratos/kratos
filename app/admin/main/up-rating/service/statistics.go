package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/up-rating/dao/global"
	"go-common/app/admin/main/up-rating/model"

	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_layout  = "2006-01-02"
	_segment = 10
)

// StatisGraph get statistics graph data
func (s *Service) StatisGraph(c context.Context, ctype int64, tagID string, Compare int) (data interface{}, err error) {
	totalScore, err := s.getTypeScore(c, ctype)
	if err != nil {
		log.Error("s.getTypeScore error(%v)", err)
		return
	}
	nowStatis, now, err := s.GetLastRatingStatis(c, ctype, tagID)
	if err != nil {
		log.Error("s.GetRatingStatis error(%v)", err)
		return
	}
	last := now.AddDate(0, -1*Compare, 0)
	lastStatis, err := s.GetRatingStatis(c, ctype, tagID, getStartMonthlyDate(last))
	if err != nil {
		log.Error("s.GetRatingStatis error(%v)", err)
		return
	}
	data = map[string]interface{}{
		"xAxis":      getSections(int(totalScore)),
		"this_month": statisSections(nowStatis, int(totalScore)),
		"compare":    statisSections(lastStatis, int(totalScore)),
	}
	return
}

// StatisList get statistics list
func (s *Service) StatisList(c context.Context, ctype int64, tagID string, Compare int) (list []*model.RatingStatis, err error) {
	list = make([]*model.RatingStatis, 0)
	totalScore, err := s.getTypeScore(c, ctype)
	if err != nil {
		log.Error("s.getTypeScore error(%v)", err)
		return
	}
	nowStatis, now, err := s.GetLastRatingStatis(c, ctype, tagID)
	if err != nil {
		log.Error("s.GetRatingStatis error(%v)", err)
		return
	}
	last := now.AddDate(0, -1*Compare, 0)
	lastStatis, err := s.GetRatingStatis(c, ctype, tagID, getStartMonthlyDate(last))
	if err != nil {
		log.Error("s.GetRatingStatis error(%v)", err)
		return
	}
	list = statisProportion(nowStatis, lastStatis, getSections(int(totalScore)), int(totalScore))
	return
}

// StatisExport export statis
func (s *Service) StatisExport(c context.Context, ctype int64, tagID string, Compare int) (res []byte, err error) {
	list, err := s.StatisList(c, ctype, tagID, Compare)
	if err != nil {
		log.Error("s.StatisList error(%v)", err)
		return
	}
	res, err = formatCSV(formatStatis(list, ctype))
	if err != nil {
		log.Error("StatisExport formatCSV error(%v)", err)
	}
	return
}

func statisProportion(now, last []*model.RatingStatis, sections []string, totalScore int) (list []*model.RatingStatis) {
	var totalUps, lastTotalUps int64
	list = make([]*model.RatingStatis, _segment)
	for i := 0; i < _segment; i++ {
		list[i] = &model.RatingStatis{
			Tips: sections[i],
		}
	}
	offset := totalScore / _segment
	for _, s := range now {
		totalUps += s.Ups
		idx := int(s.Section) * 10 / offset
		if idx >= _segment {
			idx--
		}

		list[idx].Ups += s.Ups
		list[idx].Score += s.Score
		list[idx].CreativityScore += s.CreativityScore
		list[idx].InfluenceScore += s.InfluenceScore
		list[idx].CreditScore += s.CreditScore
		list[idx].Fans += s.Fans
		list[idx].Avs += s.Avs
		list[idx].Coin += s.Coin
		list[idx].Play += s.Play
	}
	for _, s := range last {
		lastTotalUps += s.Ups
		idx := int(s.Section) * 10 / offset
		if idx >= _segment {
			idx--
		}
		list[idx].Compare += s.Ups
	}
	for i := 0; i < len(list); i++ {
		if totalUps > 0 {
			list[i].Proportion = fmt.Sprintf("%.02f", float64(list[i].Ups*100)/float64(totalUps))
		}
		if lastTotalUps > 0 {
			list[i].ComparePropor = fmt.Sprintf("%.02f", float64(list[i].Compare*100)/float64(lastTotalUps))
		}
		if list[i].Ups > 0 {
			list[i].Score /= list[i].Ups
			list[i].CreativityScore /= list[i].Ups
			list[i].InfluenceScore /= list[i].Ups
			list[i].CreditScore /= list[i].Ups
			list[i].Fans /= list[i].Ups
			list[i].Avs /= list[i].Ups
			list[i].Coin /= list[i].Ups
			list[i].Play /= list[i].Ups
		}
	}
	return
}

func getSections(score int) (sections []string) {
	sections = make([]string, _segment)
	offset := score / _segment
	for i := 0; i < len(sections); i++ {
		sections[i] = fmt.Sprintf("%d-%d", i*offset, (i+1)*offset)
	}
	return sections
}

func statisSections(statis []*model.RatingStatis, totalScore int) (ups []int64) {
	ups = make([]int64, _segment)
	offset := totalScore / _segment
	for _, s := range statis {
		idx := int(s.Section) * 10 / offset
		if idx >= len(ups) {
			idx--
		}
		ups[idx] += s.Ups
	}
	return
}

// GetLastRatingStatis get last rating statis
func (s *Service) GetLastRatingStatis(c context.Context, ctype int64, tagID string) (statis []*model.RatingStatis, date time.Time, err error) {
	statis = make([]*model.RatingStatis, 0)
	times := 0
	date = getStartMonthlyDate(time.Now()).AddDate(0, -1, 0)
	for times < 2 {
		statis, err = s.GetRatingStatis(c, ctype, tagID, date)
		if err != nil {
			return
		}
		if len(statis) > 0 {
			break
		}
		times++
		date = date.AddDate(0, -1, 0)
	}
	if times == 2 && len(statis) == 0 {
		err = fmt.Errorf("get last statis error")
		return
	}
	return
}

// GetRatingStatis get rating statis
func (s *Service) GetRatingStatis(c context.Context, ctype int64, tagID string, date time.Time) (statis []*model.RatingStatis, err error) {
	return s.dao.GetRatingStatis(c, ctype, date.Format(_layout), tagID)
}

func getStartMonthlyDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
}

// GetTrendAsc get trend asc
func (s *Service) GetTrendAsc(c context.Context, ctype string, tags []int64, date time.Time, frl, fru int, mid int64, offset, limit int) (total int, ts []*model.Trend, err error) {
	q := query(ctype, tags, frl, mid)
	dsr := getStartMonthlyDate(date).Format(_layout)

	total, err = s.dao.AscTrendCount(c, dsr, q)
	if err != nil {
		return
	}

	q += fmt.Sprintf(" ORDER BY %s_diff DESC,id LIMIT %d,%d", ctype, offset, limit)
	ts, err = s.dao.GetTrendAsc(c, ctype, dsr, q)
	if err != nil {
		return
	}

	err = fillUpNames(c, ts)
	return
}

//GetTrendDesc get trend desc
func (s *Service) GetTrendDesc(c context.Context, ctype string, tags []int64, date time.Time, frl, fru int, mid int64, offset, limit int) (total int, ts []*model.Trend, err error) {
	q := query(ctype, tags, frl, mid)
	dsr := getStartMonthlyDate(date).Format(_layout)

	total, err = s.dao.DescTrendCount(c, dsr, q)
	if err != nil {
		return
	}

	q += fmt.Sprintf(" ORDER BY %s_diff,id LIMIT %d,%d", ctype, offset, limit)
	ts, err = s.dao.GetTrendDesc(c, ctype, dsr, q)
	if err != nil {
		return
	}

	err = fillUpNames(c, ts)
	return
}

func fillUpNames(c context.Context, ts []*model.Trend) (err error) {
	var mids []int64
	for _, trend := range ts {
		mids = append(mids, trend.MID)
	}
	if len(mids) == 0 {
		return
	}
	ns, err := global.Names(c, mids)
	if err != nil {
		return
	}
	for _, trend := range ts {
		trend.Nickname = ns[trend.MID]
	}
	return
}

func getSection(ctype string, frl int) (section int, typ int) {
	if ctype == "magnetic" {
		typ = 0
		interval := 600 / 10
		section = frl / interval
	}

	if ctype == "creativity" {
		typ = 1
		interval := 200 / 10
		section = frl / interval
	}

	if ctype == "influence" {
		typ = 2
		interval := 200 / 10
		section = frl / interval
	}

	if ctype == "credit" {
		typ = 3
		interval := 200 / 10
		section = frl / interval
	}
	return
}

func query(ctype string, tags []int64, frl int, mid int64) (q string) {
	if len(tags) > 0 {
		q += fmt.Sprintf(" AND tag_id IN (%s)", xstr.JoinInts(tags))
	}
	section, typ := getSection(ctype, frl)
	q += fmt.Sprintf(" AND section=%d", section)
	q += fmt.Sprintf(" AND ctype=%d", typ)

	if mid > 0 {
		q += fmt.Sprintf(" AND mid=%d", mid)
	}
	return
}
