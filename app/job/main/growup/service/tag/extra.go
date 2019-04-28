package tag

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/tag"

	"go-common/library/log"
)

// TagExtraIncome get tag extra income
func (s *Service) TagExtraIncome(c context.Context, tagIDs []int64, tagName string, ratio int64, start, end string) (err error) {
	startDate, err := time.ParseInLocation("2006-01-02", start, time.Local)
	if err != nil {
		return
	}
	endDate, err := time.ParseInLocation("2006-01-02", end, time.Local)
	if err != nil {
		return
	}
	var ups map[int64]int64
	if ratio == 0 {
		ups, err = s.getExtraByUpIncome(c, tagIDs, startDate, endDate)
	} else {
		ups, err = s.getExtraByTag(c, tagIDs, ratio)
	}
	if err != nil {
		log.Error("getExtra income error(%v)", err)
		return
	}
	err = s.insertUpTagYear(c, ups, tagName)
	if err != nil {
		log.Error("s.insertUpTagYear error(%v)", err)
	}
	return
}

func (s *Service) getExtraByUpIncome(c context.Context, tagIDs []int64, startDate, endDate time.Time) (ups map[int64]int64, err error) {
	ups = make(map[int64]int64)
	upTags, err := s.getExtraByTag(c, tagIDs, 0)
	if err != nil {
		return
	}
	upIncomes := make([]*model.UpIncome, 0)
	for !startDate.After(endDate) {
		query := fmt.Sprintf("date = '%s'", startDate.Format("2006-01-02"))
		var upIncome []*model.UpIncome
		upIncome, err = s.getUpIncome(c, query)
		if err != nil {
			log.Error("s.GetUpIncome error(%v)", err)
			return
		}
		upIncomes = append(upIncomes, upIncome...)
		startDate = startDate.AddDate(0, 0, 1)
	}
	for _, up := range upIncomes {
		if _, ok := upTags[up.MID]; ok {
			ups[up.MID] += up.AvIncome - up.AvBaseIncome
		}
	}
	return
}

func (s *Service) getExtraByTag(c context.Context, tagIDs []int64, ratio int64) (ups map[int64]int64, err error) {
	ups = make(map[int64]int64)
	if len(tagIDs) == 0 {
		return
	}
	for _, tagID := range tagIDs {
		var upTags []*model.UpTagIncome
		upTags, err = s.getTagUpById(c, tagID)
		if err != nil {
			return
		}
		for _, up := range upTags {
			ups[up.MID] = ratio
		}
	}
	return
}

func (s *Service) getTagUpById(c context.Context, tagID int64) (upTags []*model.UpTagIncome, err error) {
	upTags = make([]*model.UpTagIncome, 0)
	var id int64
	limit := 2000
	for {
		var upTag []*model.UpTagIncome
		upTag, err = s.dao.GetTagUpByID(c, tagID, id, limit)
		if err != nil {
			log.Error("s.dao.GetTagUpByID error(%v)", err)
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

func (s *Service) insertUpTagYear(c context.Context, ups map[int64]int64, tagName string) (err error) {
	up := make(map[int64]int64)
	count := 0
	for mid, income := range ups {
		up[mid] = income
		count++
		if count >= 2000 {
			_, err = s.dao.InsertUpTagYear(c, assembleUpTagYear(up), tagName)
			if err != nil {
				return
			}
			up = make(map[int64]int64)
			count = 0
		}
	}
	if len(up) > 0 {
		_, err = s.dao.InsertUpTagYear(c, assembleUpTagYear(up), tagName)
	}
	return
}

func assembleUpTagYear(ups map[int64]int64) (vals string) {
	var buf bytes.Buffer
	for mid, income := range ups {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(income, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}
