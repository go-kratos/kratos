package tag

import (
	"bytes"
	"context"
	"strconv"

	model "go-common/app/job/main/growup/model/tag"
)

func (s *Service) delArchiveRatio(c context.Context, ctype int) (err error) {
	var limit int64 = 2000
	for {
		var rows int64
		rows, err = s.dao.DelArchiveRatio(c, ctype, limit)
		if err != nil {
			return
		}
		if rows < limit {
			break
		}
	}
	return
}

// delUpChargeRatio del up_charge_ratio
func (s *Service) delUpChargeRatio(c context.Context, ctype int) (err error) {
	var limit int64 = 2000
	for {
		var rows int64
		rows, err = s.dao.DelUpRatio(c, ctype, limit)
		if err != nil {
			return
		}
		if rows < limit {
			break
		}
	}
	return
}

// insertAvRatio insert av_charge_ratio
func (s *Service) insertAvRatio(c context.Context, avs map[int64]*model.AvTagRatio, ctype int) (err error) {
	var buf bytes.Buffer
	var count int64
	for _, av := range avs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(av.TagID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(av.AvID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(av.Ratio))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(av.AdjustType))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ctype))
		buf.WriteString("),")
		count++
		if count%2000 == 0 {
			buf.Truncate(buf.Len() - 1)
			_, err = s.dao.InsertAvRatio(c, buf.String())
			if err != nil {
				return
			}
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		_, err = s.dao.InsertAvRatio(c, buf.String())
		if err != nil {
			return
		}
		buf.Reset()
	}
	return
}

// insertUpRatio insert up_charge_ratio
func (s *Service) insertUpRatio(c context.Context, ups map[int64]*model.AvTagRatio, ctype int) (err error) {
	var buf bytes.Buffer
	var count int64
	for _, up := range ups {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(up.TagID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(up.MID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(up.Ratio))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(up.AdjustType))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(ctype))
		buf.WriteString("),")
		count++
		if count%2000 == 0 {
			buf.Truncate(buf.Len() - 1)
			_, err = s.dao.InsertUpRatio(c, buf.String())
			if err != nil {
				return
			}
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		_, err = s.dao.InsertUpRatio(c, buf.String())
		if err != nil {
			return
		}
		buf.Reset()
	}
	return
}
