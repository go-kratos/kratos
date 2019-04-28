package tag

import (
	"bytes"
	"context"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// TagUps cal tag effect up count
func (s *Service) TagUps(c context.Context, date time.Time) (err error) {
	tags, err := s.dao.AllTagInfo(c)
	if err != nil {
		log.Error("s.dao.AllTagInfo error(%v)", err)
		return
	}
	if len(tags) == 0 {
		return
	}
	d := xtime.Time(date.Unix())
	startTime := d
	for _, t := range tags {
		if t.EndAt > d && t.StartAt < startTime {
			startTime = t.StartAt
		}
	}
	uptags, err := s.upTagIncomeByDate(c, startTime.Time().Format("2006-01-02"))
	if err != nil {
		log.Error("s.upTagIncomeByDate error(%v)", err)
		return
	}
	tagMap := make(map[int64]map[int64]struct{})
	for _, u := range uptags {
		if _, ok := tagMap[u.TagID]; !ok {
			tagMap[u.TagID] = make(map[int64]struct{})
		}
		tagMap[u.TagID][u.MID] = struct{}{}
	}

	// update
	for id, tag := range tagMap {
		_, err = s.dao.UpdateTagUps(c, id, len(tag))
		if err != nil {
			log.Error("s.dao.UpdateTagUps error(%v)", err)
			return
		}
	}
	return
}

func (s *Service) upTagIncomeByDate(c context.Context, date string) (upTags []*model.UpTagIncome, err error) {
	upTags = make([]*model.UpTagIncome, 0)
	var id int64
	limit := 2000
	for {
		var upTag []*model.UpTagIncome
		upTag, err = s.dao.UpTagIncomeByDate(c, date, id, limit)
		if err != nil {
			log.Error("s.dao.UpTagIncomeByDate error(%v)", err)
			return
		}
		upTags = append(upTags, upTag...)
		if len(upTag) < limit {
			break
		}
		id = upTag[len(upTag)-1].ID
	}
	return
}

// GetUpTagIncomeMap get up_tag_income map[mid]income
func (s *Service) GetUpTagIncomeMap(c context.Context) (avs map[int64]int64, err error) {
	avs = make(map[int64]int64)
	var id int64
	count, limit := 0, 2000
	for {
		id, count, err = s.dao.GetUpTagIncomeMap(c, id, limit, avs)
		if err != nil {
			return
		}
		if count < limit {
			break
		}
	}
	return
}

// TxInsertUpTagIncome insert up_tag_income
func (s *Service) TxInsertUpTagIncome(tx *sql.Tx, avs []*model.AvTagRatio) (err error) {
	if len(avs) == 0 {
		return
	}
	start, offset := 0, 2000
	if len(avs) < offset {
		offset = len(avs)
	}
	for start+offset <= len(avs) {
		_, err = s.InsertUpTagIncomeBatch(tx, avs[start:start+offset])
		if err != nil {
			tx.Rollback()
			return
		}
		start += offset
		if start < len(avs) && start+offset > len(avs) {
			offset = len(avs) - start
		}
	}
	return
}

// InsertUpTagIncomeBatch insert up_tag_income batch
func (s *Service) InsertUpTagIncomeBatch(tx *sql.Tx, avs []*model.AvTagRatio) (rows int64, err error) {
	var buf bytes.Buffer
	for _, a := range avs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(a.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.BaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + a.Date + "\"")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	_, err = s.dao.TxInsertUpTagIncome(tx, values)
	return
}
