package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_layout = "2006-01-02"
)

// InsertTagIncome insert up_tag_income.
func (s *Service) InsertTagIncome(c context.Context, date time.Time) (err error) {
	infos, err := s.getTagAvInfo(c, date)
	if err != nil {
		log.Error("s.InsertTagIncome getTagAVInfo error(%v)", err)
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.InsertTagIncome dao.BeginTran error(%v)", err)
		return
	}
	if err = s.insertTagIncome(c, tx, infos); err != nil {
		tx.Rollback()
		log.Error("s.InsertTagIncome insertTagIncome error(%v)", err)
		return
	}
	if err = s.updateTagInfo(tx, infos); err != nil {
		tx.Rollback()
		log.Error("s.InsertTagIncome updateTagInfo error(%v)", err)
		return
	}
	if err = s.updateTagUpInfo(tx, infos); err != nil {
		tx.Rollback()
		log.Error("s.InsertTagIncome updateTagUpInfo error(%v)", err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("s.InsertTagIncome tx.Commit error")
		return
	}
	return
}

func (s *Service) getTagAvInfo(c context.Context, date time.Time) (infos []*model.TagAvIncome, err error) {
	var (
		from, limit int64
		av, avs     []*model.ActivityAVInfo
	)
	from, limit = 0, 3000
	for {
		av, err = s.dao.GetAvTagRatio(c, from, limit)
		if err != nil {
			log.Error("s.getTagAvInfo dao.GetAvTagRatio error(%v)", err)
			return
		}
		avs = append(avs, av...)
		if int64(len(av)) < limit {
			break
		}
		from = av[len(av)-1].MID
	}
	for _, a := range avs {
		var income *model.TagAvIncome
		income, err = s.dao.GetAvIncomeInfo(c, a.AVID, date)
		if err != nil {
			log.Error("s.GetAvIncomes dao.GetAvIncomeInfo error(%v)", err)
			return
		}
		if income == nil {
			continue
		}
		income.TagID = a.TagID
		infos = append(infos, income)
	}
	return
}

func (s *Service) getTagAVLatestTotalIncome(c context.Context, avID, tagID int64) (totalIncome int, err error) {
	infos, err := s.dao.GetTagAvTotalIncome(c, tagID, avID)
	if err != nil {
		log.Error("s.getTagAVLatestTotalIncome dao.GetTagAvTotalIncome error(%v)", err)
		return
	}
	for _, info := range infos {
		if int(info.TotalIncome) > totalIncome {
			totalIncome = int(info.TotalIncome)
		}
	}
	return
}

func (s *Service) insertTagIncome(c context.Context, tx *sql.Tx, infos []*model.TagAvIncome) (err error) {
	var buf bytes.Buffer
	var cnt, totalIncome int
	var rows, totalRows int64
	for _, info := range infos {
		totalIncome, err = s.getTagAVLatestTotalIncome(c, info.TagID, info.AVID)
		if err != nil {
			log.Error("s.insertTagIncome dao.GetTagAvTotalIncome error(%v)", err)
			return
		}
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(info.TagID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(info.MID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(info.AVID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(info.Income))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(totalIncome + info.Income))
		buf.WriteString(",")
		buf.WriteString("'")
		buf.WriteString(strconv.Itoa(info.Date.Year()))
		buf.WriteString("-")
		if int(info.Date.Month()) < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(strconv.Itoa(int(info.Date.Month())))
		buf.WriteString("-")
		if info.Date.Day() < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(strconv.Itoa(info.Date.Day()))
		buf.WriteString("'")
		buf.WriteString("),")
		cnt++
		if cnt%1000 == 0 {
			buf.Truncate(buf.Len() - 1)
			rows, err = s.dao.TxInsertTagIncome(tx, buf.String())
			if err != nil {
				log.Error("s.InsertTagIncome dao.TxInsertTagIncome error(%v)", err)
				return
			}
			totalRows += rows
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		rows, err = s.dao.TxInsertTagIncome(tx, buf.String())
		if err != nil {
			log.Error("s.InsertTagIncome dao.TxInsertTagIncome error(%v)", err)
			return
		}
		totalRows += rows
	}
	log.Info("s.InsertTagIncome insert up_tag_income (%d) rows", totalRows)
	return
}

// updateTagInfo update tag_info total_income.
func (s *Service) updateTagInfo(tx *sql.Tx, infos []*model.TagAvIncome) (err error) {
	tim := make(map[int64]int) // key-value: tag_id-total av income.
	for _, info := range infos {
		tim[info.TagID] += info.Income
	}
	for k, v := range tim {
		query := "total_income = total_income + "
		query += strconv.Itoa(v)
		_, err = s.dao.TxUpdateTagInfo(tx, k, query)
		if err != nil {
			log.Error("s.updateTagInfo dao.UpdateTagInfo error(%v)", err)
			return
		}
	}
	return
}

// updateTagUpInfo update tag_up_info totalIncome.
func (s *Service) updateTagUpInfo(tx *sql.Tx, infos []*model.TagAvIncome) (err error) {
	utm := make(map[int64]*model.TagAvIncome) // key-value: mid-totalIncome
	for _, info := range infos {
		_, ok := utm[info.MID]
		if !ok {
			a := &model.TagAvIncome{TagID: info.TagID, MID: info.MID, TotalIncome: info.Income}
			utm[info.MID] = a
		} else {
			utm[info.MID].TotalIncome += info.Income
		}
	}
	cnt := 0
	var buf bytes.Buffer
	for _, v := range utm {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.TagID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.MID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(v.TotalIncome))
		buf.WriteString("),")
		cnt++
		if cnt%2000 == 0 {
			buf.Truncate(buf.Len() - 1)
			_, err = s.dao.TxUpdateTagUpInfo(tx, buf.String())
			if err != nil {
				log.Error("s.updateTagUpInfo dao.UpdateTagUpInfo error(%v)", err)
				return
			}
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		_, err = s.dao.TxUpdateTagUpInfo(tx, buf.String())
		if err != nil {
			log.Error("s.updateTagUpInfo dao.UpdateTagUpInfo error(%v)", err)
			return
		}
	}
	return
}

// GetAvIncomeStatis get av monthly income
func (s *Service) GetAvIncomeStatis(c context.Context, date string) error {
	d, _ := time.Parse(_layout, date)
	endTime := d.AddDate(0, 1, 0).Format(_layout)
	avs, err := s.GetAvIncome(c)
	if err != nil {
		log.Error("s.GetAvIncome error(%v)", err)
		return err

	}
	log.Info("GetAvIncomeStatis get %d avs", len(avs))

	avsMap := make(map[int64]*model.AvIncome)
	avIncomeStatis(avsMap, avs, date, endTime)

	data := make([]*model.AvIncome, 0)
	for _, av := range avsMap {
		data = append(data, av)
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Income > data[j].Income
	})
	log.Info("GetAvIncomeStatis calculate success: %d", len(data))
	return s.batchSend(data)
}

func (s *Service) batchSend(data []*model.AvIncome) error {
	fileNo, start, batchSize := 0, 0, 50000
	for {
		if start+batchSize >= len(data) {
			batchSize = len(data) - start
		}
		if batchSize <= 0 {
			break
		}

		records := formatAvIncome(data[start : start+batchSize])
		filename := fmt.Sprintf("av_statis0%d", fileNo)
		err := WriteCSV(records, filename)
		if err != nil {
			return err
		}

		err = s.email.SendMailAttach(filename, "稿件月收入", []string{"shaozhenyu@bilibili.com"})
		if err != nil {
			return err
		}
		fileNo++
		start += batchSize
	}
	return nil
}

func avIncomeStatis(avsMap map[int64]*model.AvIncome, avs []*model.AvIncome, fromTime, toTime string) {
	for _, av := range avs {
		d := av.Date.Time().Format(_layout)
		if d < fromTime || d >= toTime {
			continue
		}
		if _, ok := avsMap[av.AvID]; ok {
			avsMap[av.AvID].Income += av.Income
			if avsMap[av.AvID].Date < av.Date {
				avsMap[av.AvID].TotalIncome = av.TotalIncome
				avsMap[av.AvID].Date = av.Date
			}
		} else {
			avsMap[av.AvID] = av
		}
	}
}

// GetAvIncome get av_income
func (s *Service) GetAvIncome(c context.Context) (avs []*model.AvIncome, err error) {
	limit := 2000
	var id int64
	for {
		av, err := s.dao.ListAvIncome(c, id, limit)
		if err != nil {
			return avs, err
		}
		avs = append(avs, av...)
		if len(av) < limit {
			break
		}
		id = av[len(av)-1].ID
	}
	return
}

// GetUpIncomeStatis get up statis
func (s *Service) GetUpIncomeStatis(c context.Context, date string, hasWithdraw int) (err error) {
	var upAccount []*model.UpAccount
	if hasWithdraw == 1 {
		upAccount, err = s.getUpIncomeStatisAfterWithdraw(c, date)
		if err != nil {
			log.Error("s.getUpIncomeStatisAfterWithdraw error(%v)", err)
			return err
		}
	} else {
		upAccount, err = s.getUpIncomeStatisBeforeWithdraw(c, date)
		if err != nil {
			log.Error("s.getUpIncomeStatisBeforeWithdraw error(%v)", err)
			return err
		}
	}

	upa := make(map[int64]*model.UpAccount)
	mids := make([]int64, len(upAccount))
	for i := 0; i < len(upAccount); i++ {
		mids[i] = upAccount[i].MID
		upa[upAccount[i].MID] = upAccount[i]
	}

	upNick, err := s.GetUpNickname(c, mids)
	if err != nil {
		log.Error("s.GetUpNickname error(%v)", err)
		return
	}

	upIncome, err := s.GetUpIncome(c, "up_income_monthly", date)
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}
	upIncomeStatis(upa, upIncome, upNick)

	data := []*model.UpAccount{}
	for _, up := range upa {
		data = append(data, up)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].MonthIncome > data[j].MonthIncome
	})

	d, _ := time.Parse(_layout, date)
	records := formatUpAccount(data, int(d.Month()))
	err = WriteCSV(records, "up_statis.csv")
	if err != nil {
		log.Error("WriteCSV error(%v)", err)
		return
	}

	return s.email.SendMailAttach("up_statis.csv", "up主月结算", []string{"shaozhenyu@bilibili.com"})

}

func (s *Service) getUpIncomeStatisAfterWithdraw(c context.Context, date string) (upAccount []*model.UpAccount, err error) {
	d, _ := time.Parse(_layout, date)
	withdrawDateStr := d.Format("2006-01")
	ctime := d.AddDate(0, 1, 1).Format(_layout)
	upAccount, err = s.GetUpAccount(c, withdrawDateStr, ctime)
	if err != nil {
		log.Error("s.GetUpAccount error(%v)", err)
		return
	}

	upWithdraw, err := s.GetUpWithdraw(c, withdrawDateStr)
	if err != nil {
		log.Error("s.GetUpWithdraw error(%v)", err)
		return
	}
	for _, up := range upAccount {
		up.TotalUnwithdrawIncome = upWithdraw[up.MID]
	}
	return
}

func (s *Service) getUpIncomeStatisBeforeWithdraw(c context.Context, date string) (upAccount []*model.UpAccount, err error) {
	d, _ := time.Parse(_layout, date)
	withdrawDateStr := d.AddDate(0, -1, 0).Format("2006-01")
	ctime := d.AddDate(0, 1, 1).Format(_layout)
	return s.GetUpAccount(c, withdrawDateStr, ctime)
}

func upIncomeStatis(upa map[int64]*model.UpAccount, upIncome []*model.UpIncome, upNick map[int64]string) {
	for _, income := range upIncome {
		if _, ok := upa[income.MID]; ok {
			upa[income.MID].AvCount = income.AvCount
			upa[income.MID].MonthIncome = income.Income
			upa[income.MID].Nickname = upNick[income.MID]
		}
	}
}

// GetUpAccount get up_account
func (s *Service) GetUpAccount(c context.Context, date, ctime string) (ups []*model.UpAccount, err error) {
	offset, limit := 0, 2000
	for {
		up, err := s.dao.ListUpAccount(c, date, ctime, offset, limit)
		if err != nil {
			return ups, err
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
		offset += limit
	}
	return
}

// GetUpIncome get up_income
func (s *Service) GetUpIncome(c context.Context, table, date string) (ups []*model.UpIncome, err error) {
	ups = make([]*model.UpIncome, 0)
	var id int64
	limit := 2000
	for {
		var up []*model.UpIncome
		up, err = s.dao.ListUpIncome(c, table, date, id, limit)
		if err != nil {
			return
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
		id = up[len(up)-1].ID
	}
	return
}

// GetUpWithdraw get up_income_withdraw
func (s *Service) GetUpWithdraw(c context.Context, date string) (ups map[int64]int64, err error) {
	ups = make(map[int64]int64)
	offset, limit := 0, 2000
	for {
		up, err := s.dao.ListUpWithdraw(c, date, offset, limit)
		if err != nil {
			return ups, err
		}
		for mid, income := range up {
			ups[mid] = income
		}
		if len(up) < limit {
			break
		}
		offset += limit
	}
	return
}

// GetUpNickname get up nickname
func (s *Service) GetUpNickname(c context.Context, mids []int64) (upNick map[int64]string, err error) {
	upNick = make(map[int64]string)
	offset, limit := 0, 2000
	for {
		if offset+limit > len(mids) {
			limit = len(mids) - offset
		}
		if limit <= 0 {
			break
		}
		err = s.dao.ListUpNickname(c, mids[offset:offset+limit], upNick)
		if err != nil {
			log.Error("s.dao.ListUpNickname error(%v)", err)
			return
		}
		offset += limit
	}
	return
}
