package income

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	upModel "go-common/app/admin/main/growup/model"
	model "go-common/app/admin/main/growup/model/income"
	"go-common/app/admin/main/growup/service"

	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"golang.org/x/sync/errgroup"
)

// BreachList list
func (s *Service) BreachList(c context.Context, mids, aids []int64, typ int, fromTime, toTime int64, reason string, from, limit int) (breachs []*model.AvBreach, total int, err error) {
	query := formatBreachQuery(mids, aids, typ, fromTime, toTime, reason)
	total, err = s.dao.BreachCount(c, query)
	if err != nil {
		log.Error("s.dao.GetBreachCount error(%v)", err)
		return
	}
	query = fmt.Sprintf("%s LIMIT %d,%d", query, from, limit)
	breachs, err = s.dao.ListArchiveBreach(c, query)
	if err != nil {
		log.Error("s.dao.ListArchiveBreach error(%v)", err)
	}

	mids = make([]int64, 0, len(breachs))
	for _, b := range breachs {
		mids = append(mids, b.MID)
	}
	nickname, err := s.dao.ListUpInfo(c, mids)
	for _, b := range breachs {
		b.Nickname = nickname[b.MID]
	}
	return
}

// BreachStatis statis
func (s *Service) BreachStatis(c context.Context, mids, aids []int64, typ, groupType int, fromTime, toTime int64, reason string) (date interface{}, err error) {
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	query := formatBreachQuery(mids, aids, typ, from.Unix(), to.Unix(), reason)
	breachs, err := s.dao.ListArchiveBreach(c, query)
	if err != nil {
		log.Error("s.dao.ListArchiveBreach error(%v)", err)
		return
	}
	date = breachStatis(breachs, from, to, groupType)
	return
}

func breachStatis(breachs []*model.AvBreach, from, to time.Time, groupType int) interface{} {
	dateIncome := make(map[string]int64)
	dateUps := make(map[string]map[int64]struct{})
	for _, breach := range breachs {
		date := formatDateByGroup(breach.CDate.Time(), groupType)
		if _, ok := dateIncome[date]; ok {
			dateIncome[date] += breach.Money
		} else {
			dateIncome[date] = breach.Money
		}
		if _, ok := dateUps[date]; !ok {
			dateUps[date] = make(map[int64]struct{})
		}
		dateUps[date][breach.MID] = struct{}{}
	}

	income, counts, xAxis := []string{}, []int{}, []string{}
	// get result by date
	to = to.AddDate(0, 0, 1)
	for from.Before(to) {
		dateStr := formatDateByGroup(from, groupType)
		xAxis = append(xAxis, dateStr)
		if val, ok := dateIncome[dateStr]; ok {
			income = append(income, fmt.Sprintf("%.2f", float64(val)/float64(100)))
			counts = append(counts, len(dateUps[dateStr]))
		} else {
			income = append(income, "0")
			counts = append(counts, 0)
		}
		from = addDayByGroup(groupType, from)
	}

	return map[string]interface{}{
		"counts":  counts,
		"incomes": income,
		"xaxis":   xAxis,
	}
}

// ExportBreach export
func (s *Service) ExportBreach(c context.Context, mids, aids []int64, typ int, fromTime, toTime int64, reason string, from, limit int) (res []byte, err error) {
	breachs, _, err := s.BreachList(c, mids, aids, typ, fromTime, toTime, reason, from, limit)
	if err != nil {
		log.Error("s.BreachList error(%v)", err)
		return
	}
	records := formatBreach(breachs)
	res, err = service.FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)")
	}
	return
}

func formatBreachQuery(mids, aids []int64, typ int, fromTime, toTime int64, reason string) (query string) {
	query = fmt.Sprintf("cdate >= '%s' AND cdate <= '%s'", time.Unix(fromTime, 0).Format(_layout), time.Unix(toTime, 0).Format(_layout))
	if typ != 4 {
		query = fmt.Sprintf("%s AND ctype = %d", query, typ)
	}
	if len(mids) > 0 {
		query = fmt.Sprintf("%s AND mid IN (%s)", query, xstr.JoinInts(mids))
	}
	if len(aids) > 0 {
		query = fmt.Sprintf("%s AND av_id IN (%s)", query, xstr.JoinInts(aids))
	}
	if reason != "" {
		query = fmt.Sprintf("%s AND reason = '%s'", query, reason)
	}
	return
}

// ArchiveBreach breach archive batch
func (s *Service) ArchiveBreach(c context.Context, typ int, aids []int64, mid int64, reason string, operator string) (err error) {
	count, err := s.dao.BreachCount(c, fmt.Sprintf("av_id in (%s) AND cdate = '%s'", xstr.JoinInts(aids), time.Now().Format(_layout)))
	if err != nil {
		log.Error("s.dao.AvBreachCount error(%v)", err)
		return
	}
	if count > 0 {
		err = fmt.Errorf("有稿件已被扣除")
		return
	}
	return s.avBreach(c, typ, aids, mid, reason, operator)
}

func assembleAvBreach(aids []int64, mid int64, ctype int, reason string, breach map[int64]int64, upload map[int64]string) (vals string) {
	var buf bytes.Buffer
	for _, aid := range aids {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(aid, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + time.Now().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(breach[aid], 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(ctype))
		buf.WriteByte(',')
		buf.WriteString("\"" + reason + "\"")
		buf.WriteByte(',')
		buf.WriteString("'" + upload[aid] + "'")
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

// AvBreach av breach from av_income
func (s *Service) avBreach(c context.Context, ctype int, aids []int64, mid int64, reason string, operator string) (err error) {
	archives, withdrawMonth, err := s.GetArchiveByUpAccount(c, ctype, aids, mid)
	if err != nil {
		log.Error("s.GetArchiveByUpAccount error(%v)", err)
		return
	}
	if len(archives) == 0 {
		return
	}
	preMonthBreach, thisMonthBreach, avBreach, avUpload := getBreachMoney(archives, withdrawMonth)

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	var eg errgroup.Group
	// insert av_breach
	eg.Go(func() (err error) {
		if _, err = s.dao.TxInsertAvBreach(tx, assembleAvBreach(aids, mid, ctype, reason, avBreach, avUpload)); err != nil {
			log.Error("s.TxInsertAvBreach error(%v)", err)
			tx.Rollback()
		}
		return
	})

	// update av breach pre state = 2
	eg.Go(func() (err error) {
		if _, err = s.dao.TxUpdateBreachPre(tx, aids, time.Now().Format(_layout)); err != nil {
			log.Error("s.TxUpdateBreachPre error(%v)", err)
			tx.Rollback()
		}
		return
	})

	// save av_black_list
	eg.Go(func() (err error) {
		if err = s.TxInsertAvBlacklist(c, tx, ctype, aids, mid, _avBreach, len(aids)); err != nil {
			log.Error("s.InsertAvBlacklist error(%v)", err)
		}
		return
	})
	// update up_account
	eg.Go(func() (err error) {
		if err = s.TxUpAccountBreach(c, tx, mid, preMonthBreach, thisMonthBreach); err != nil {
			log.Error("s.UpdateUpAccount error(%v)", err)
		}
		return
	})
	// update up credit score
	eg.Go(func() (err error) {
		if err = s.UpdateUpCredit(c, tx, ctype, mid, aids, operator); err != nil {
			log.Error("s.UpdateUpCredit error(%v)", err)
		}
		return
	})

	eg.Go(func() (err error) {
		if _, err = s.upDao.TxUpdateAvSpyState(tx, 1, aids); err != nil {
			tx.Rollback()
			log.Error("s.upDao.TxUpdateAvSpyState error(%v)", err)
		}
		return
	})

	if err = eg.Wait(); err != nil {
		log.Error("run eg.Wait error(%v)", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
		return
	}

	var business string
	switch ctype {
	case _video:
		business = "avid"
	case _column:
		business = "cv"
	case _bgm:
		business = "au"
	}

	for _, aid := range aids {
		err = s.msg.Send(c, "1_14_5",
			fmt.Sprintf("您的稿件 %s %d 违反创作激励计划规则。", business, aid),
			fmt.Sprintf("您的稿件 %s %d 因为%s原因被取消参加创作激励计划资格，已获得收入将被扣除。如有疑问，请联系客服。", business, aid, reason),
			[]int64{mid},
			time.Now().Unix())
		if err != nil {
			log.Error("s.msg.Send error(%v)", err)
			return
		}
	}

	return
}

// UpdateUpCredit update up_info credit score
func (s *Service) UpdateUpCredit(c context.Context, tx *sql.Tx, ctype int, mid int64, aids []int64, operator string) (err error) {
	score, err := s.upDao.CreditScore(c, mid)
	if err != nil {
		return
	}
	// insert credit_score_record
	creditRecord := &upModel.CreditRecord{
		MID:       mid,
		OperateAt: xtime.Time(time.Now().Unix()),
		Operator:  operator,
		Reason:    9,
		Deducted:  3 * len(aids),
		Remaining: score - 3*len(aids),
	}
	r1, err := s.upDao.TxInsertCreditRecord(tx, creditRecord)
	if err != nil {
		tx.Rollback()
		return
	}

	r2, err := s.upDao.TxUpdateCreditScore(tx, mid, score-3*len(aids))
	if err != nil {
		tx.Rollback()
		return
	}
	if r1 != r2 {
		tx.Rollback()
		return
	}
	return
}

// GetArchiveByUpAccount get archive income by withdraw date
func (s *Service) GetArchiveByUpAccount(c context.Context, typ int, aids []int64, mid int64) (archives []*model.ArchiveIncome, withdrawMonth time.Month, err error) {
	archives = make([]*model.ArchiveIncome, 0)
	upAccount, err := s.dao.GetUpAccount(c, mid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("s.dao.GetUpAccount error(%v)", err)
		return
	}

	withdrawDateStr := upAccount.WithdrawDateVersion + "-01"
	withdrawDate, err := time.Parse(_layout, withdrawDateStr)
	if err != nil {
		log.Error("time.Parse error(%v)", err)
		return
	}
	withdrawMonth = withdrawDate.AddDate(0, 1, 0).Month()
	from, now := withdrawDate.AddDate(0, 1, 0).Format(_layout), time.Now().Format(_layout)
	query := ""
	switch typ {
	case _video:
		query = fmt.Sprintf("av_id in (%s)", xstr.JoinInts(aids))
	case _column:
		query = fmt.Sprintf("aid in (%s)", xstr.JoinInts(aids))
	case _bgm:
		query = fmt.Sprintf("sid in (%s)", xstr.JoinInts(aids))
	}
	archives, err = s.GetArchiveIncome(c, typ, query, from, now)
	if err != nil {
		log.Error("s.GetArchiveIncome error(%v)", err)
	}
	return
}

func getBreachMoney(archives []*model.ArchiveIncome, withdrawMonth time.Month) (int64, int64, map[int64]int64, map[int64]string) {
	var preMonthBreach, thisMonthBreach int64 = 0, 0
	avBreach := make(map[int64]int64)
	avUpload := make(map[int64]string)
	for _, arch := range archives {
		if arch.Date.Time().Month() == withdrawMonth {
			preMonthBreach += arch.Income
		} else {
			thisMonthBreach += arch.Income
		}
		avBreach[arch.AvID] += arch.Income
		avUpload[arch.AvID] = arch.UploadTime.Time().Format(_layout)
	}
	return preMonthBreach, thisMonthBreach, avBreach, avUpload
}
