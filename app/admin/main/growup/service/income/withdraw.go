package income

import (
	"context"
	"fmt"
	"sort"
	"time"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/app/admin/main/growup/service"

	"go-common/library/log"
	"go-common/library/xstr"
)

// UpWithdraw get up_account infos
func (s *Service) UpWithdraw(c context.Context, mids []int64, isDeleted, from, limit int) (data []*model.UpWithdrawRes, total int64, err error) {
	query := ""
	if len(mids) != 0 {
		query += fmt.Sprintf("mid in (%s)", xstr.JoinInts(mids))
	}
	total, err = s.UpAccountCount(c, query, isDeleted)
	if err != nil {
		log.Error("s.UpAccountCount error(%v)", err)
		return
	}
	ups, err := s.ListUpAccount(c, query, isDeleted, from, limit)
	if err != nil {
		log.Error("s.ListUpAccount error(%v)", err)
		return
	}

	if len(mids) == 0 {
		for _, up := range ups {
			mids = append(mids, up.MID)
		}
	}
	query = ""
	if len(mids) != 0 {
		query = fmt.Sprintf("mid in (%s)", xstr.JoinInts(mids))
	}
	upWithdraw, err := s.GetUpWithdraw(c, query)
	if err != nil {
		log.Error("s.GetUpWithdraw error(%v)", err)
		return
	}
	withdrawMap := make(map[string]int64)
	for _, w := range upWithdraw {
		key := fmt.Sprintf("%d+%s", w.MID, w.DateVersion)
		withdrawMap[key] = w.WithdrawIncome
	}

	upInfo, err := s.dao.ListUpInfo(c, mids)
	if err != nil {
		log.Error("s.dao.ListUpInfo error(%v)", err)
		return
	}

	sort.Slice(ups, func(i, j int) bool {
		return ups[i].LastWithdrawTime > ups[j].LastWithdrawTime
	})

	data = make([]*model.UpWithdrawRes, 0)
	for _, up := range ups {
		key := fmt.Sprintf("%d+%s", up.MID, up.WithdrawDateVersion)
		data = append(data, &model.UpWithdrawRes{
			MID:                up.MID,
			Nickname:           upInfo[up.MID],
			WithdrawIncome:     fmt.Sprintf("%0.2f", fromYuanToFen(up.TotalWithdrawIncome)),
			UnWithdrawIncome:   fmt.Sprintf("%0.2f", fromYuanToFen(up.TotalUnwithdrawIncome)),
			LastWithdrawIncome: fmt.Sprintf("%0.2f", fromYuanToFen(withdrawMap[key])),
			WithdrawDate:       up.LastWithdrawTime.Time().Format(_layout),
			MTime:              up.MTime,
		})
	}
	return
}

// UpWithdrawExport export up withdraw
func (s *Service) UpWithdrawExport(c context.Context, mids []int64, isDeleted, from, limit int) (res []byte, err error) {
	upWithdraw, _, err := s.UpWithdraw(c, mids, isDeleted, from, limit)
	if err != nil {
		log.Error("s.UpWithdraw error(%v)", err)
		return
	}

	records := formatUpWithdraw(upWithdraw, isDeleted)
	res, err = service.FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)")
	}
	return
}

// GetUpWithdraw get up_withdraw
func (s *Service) GetUpWithdraw(c context.Context, query string) (upW []*model.UpIncomeWithdraw, err error) {
	upW = make([]*model.UpIncomeWithdraw, 0)
	var id int64
	limit := 2000
	for {
		var w []*model.UpIncomeWithdraw
		w, err = s.ListUpWithdraw(c, id, query, limit)
		if err != nil {
			return
		}
		upW = append(upW, w...)
		if len(w) < limit {
			break
		}
		id = w[len(w)-1].ID
	}
	return
}

// ListUpWithdraw list up_withdraw
func (s *Service) ListUpWithdraw(c context.Context, id int64, query string, limit int) (upWithdraw []*model.UpIncomeWithdraw, err error) {
	if query != "" {
		query += " AND"
	}
	return s.dao.ListUpWithdraw(c, id, query, limit)
}

// UpWithdrawStatis up_withdraw statis
func (s *Service) UpWithdrawStatis(c context.Context, from, to int64, isDeleted int) (data interface{}, err error) {
	now := time.Now()
	var fromTime, toTime time.Time
	if from == 0 || to == 0 {
		fromTime = time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, time.Local)
		toTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	} else {
		fromTime = time.Unix(from, 0)
		toTime = time.Unix(to, 0)
	}
	query := fmt.Sprintf("date_version >= '%s' and date_version <= '%s'", fromTime.Format(_layoutMonth), toTime.Format(_layoutMonth))
	upWithdraw, err := s.GetUpWithdraw(c, query)
	if err != nil {
		log.Error("s.GetUpWithdraw error(%v)", err)
		return
	}

	deletedUp, err := s.GetUpAccount(c, "", 1)
	if err != nil {
		log.Error("s.GetUpAccount error(%v)", err)
		return
	}
	deletedUpMap := make(map[int64]struct{})
	for _, up := range deletedUp {
		deletedUpMap[up.MID] = struct{}{}
	}

	data = upWithdrawStatis(upWithdraw, deletedUpMap, isDeleted, fromTime, toTime)
	return
}

func upWithdrawStatis(upWithdraw []*model.UpIncomeWithdraw, deletedUp map[int64]struct{}, isDeleted int, startTime, endTime time.Time) interface{} {
	dateCount := make(map[string]int)
	dateIncome := make(map[string]int64)
	for _, up := range upWithdraw {
		_, ok := deletedUp[up.MID]
		if (ok && isDeleted == 1) || (!ok && isDeleted == 0) {
			dateCount[up.DateVersion]++
			dateIncome[up.DateVersion] += up.WithdrawIncome
		}
	}

	incomes, counts, xAxis := []string{}, []int{}, []string{}
	endTime = endTime.AddDate(0, 0, 1)
	for startTime.Before(endTime) {
		key := startTime.Format(_layoutMonth)
		incomes = append(incomes, fmt.Sprintf("%0.2f", fromYuanToFen(dateIncome[key])))
		counts = append(counts, dateCount[key])
		xAxis = append(xAxis, key)
		startTime = startTime.AddDate(0, 1, 0)
	}

	return map[string]interface{}{
		"incomes": incomes,
		"counts":  counts,
		"xaxis":   xAxis,
	}
}

// UpWithdrawDetail get up withdraw by mid
func (s *Service) UpWithdrawDetail(c context.Context, mid int64) (upWithdraw []*model.UpIncomeWithdraw, err error) {
	query := fmt.Sprintf("mid = %d", mid)
	upWithdraw, err = s.GetUpWithdraw(c, query)
	if err != nil {
		log.Error("s.GetUpWithdraw error(%v)", err)
		return
	}
	if len(upWithdraw) == 0 {
		return
	}

	upInfo, err := s.dao.GetUpInfoNickname(c, []int64{upWithdraw[0].MID})
	if err != nil {
		log.Error("s.dao.GetUpInfoNickname error(%v)", err)
		return
	}
	for _, up := range upWithdraw {
		up.Nickname = upInfo[up.MID]
		up.Income = fmt.Sprintf("%.2f", fromYuanToFen(up.WithdrawIncome))
	}
	return
}

// UpWithdrawDetailExport export up withdraw detail
func (s *Service) UpWithdrawDetailExport(c context.Context, mid int64) (res []byte, err error) {
	upWithdraw, err := s.UpWithdrawDetail(c, mid)
	if err != nil {
		log.Error("s.UpWithdrawDetail error(%v)", err)
		return
	}

	records := formatUpIncomeWithdraw(upWithdraw)
	res, err = service.FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)")
	}
	return
}
