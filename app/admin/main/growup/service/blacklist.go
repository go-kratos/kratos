package service

import (
	"fmt"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

// QueryBlacklist query blacklist
func (s *Service) QueryBlacklist(fromTime, toTime int64, ctype, reason int, mid int64, nickname string, aid int64, from, limit int, sort string) (total int, blacklist []*model.Blacklist, err error) {
	blacklistQuery := buildBlacklistQuery(fromTime, toTime, ctype, reason, mid, nickname, aid)
	blacklist, total, err = s.dao.ListBlacklist(blacklistQuery, from, limit, sort)
	if err != nil {
		log.Error("s.dao.ListBlacklist error(%v)", err)
		return
	}
	if total == 0 {
		return
	}
	avIDs := make([]int64, len(blacklist))
	for i := 0; i < len(blacklist); i++ {
		avIDs[i] = blacklist[i].AvID
	}

	// get av total income
	avIncomeMap, err := s.getAvIncomeStatis(avIDs)
	if err != nil {
		log.Error("s.getAvIncomeStatis error(%v)", err)
		return
	}

	for i := 0; i < len(blacklist); i++ {
		blacklist[i].Income = avIncomeMap[blacklist[i].AvID]
	}

	return
}

func buildBlacklistQuery(fromTime, toTime int64, ctype, reason int, mid int64, nickname string, aid int64) (query string) {
	query += fmt.Sprintf("ctime >= '%s' AND ctime <= '%s'", time.Unix(fromTime, 0).Format("2006-01-02"), time.Unix(toTime, 0).Format("2006-01-02"))
	query += " AND "

	if aid <= 0 && mid <= 0 && nickname == "" {
		query += "has_signed = 1"
		query += " AND "
	}

	if mid > 0 {
		query += fmt.Sprintf("mid = %d", mid)
		query += " AND "
	}
	if nickname != "" {
		query += fmt.Sprintf("nickname = \"%s\"", nickname)
		query += " AND "
	}
	if aid > 0 {
		query += fmt.Sprintf("av_id = %d", aid)
		query += " AND "
	}
	if ctype != 4 {
		query += fmt.Sprintf("ctype = %d", ctype)
		query += " AND "
	}
	if reason > 0 {
		query += fmt.Sprintf("reason = %d", reason)
		query += " AND "
	}
	query += "is_delete = 0"
	return
}

func (s *Service) getAvIncomeStatis(avIDs []int64) (avIncomeMap map[int64]int64, err error) {
	avIncomeMap = make(map[int64]int64)
	query := fmt.Sprintf("av_id IN (%s) AND is_deleted = 0", xstr.JoinInts(avIDs))
	avIncomes, err := s.dao.GetAvIncomeStatis(query)
	if err != nil {
		log.Error("s.dao.GetAvIncomeStatis error(%v)", err)
		return
	}

	for _, avIncome := range avIncomes {
		avIncomeMap[avIncome.AvID] = avIncome.TotalIncome
	}
	return
}

// ExportBlacklist blacklist export csv
func (s *Service) ExportBlacklist(fromTime, toTime int64, ctype, reason int, mid int64, nickname string, aid int64, from, limit int, sort string) (res []byte, err error) {
	_, blacklist, err := s.QueryBlacklist(fromTime, toTime, ctype, reason, mid, nickname, aid, from, limit, sort)
	if err != nil {
		log.Error("s.QueryBlacklist error(%v)", err)
		return
	}

	records := formatBlacklist(blacklist)
	res, err = FormatCSV(records)
	if err != nil {
		log.Error("FormatCSV error(%v)", err)
	}
	return
}

// RecoverBlacklist recover av from blacklist
func (s *Service) RecoverBlacklist(aID int64, ctype int) (err error) {
	update := map[string]interface{}{
		"is_delete": 1,
	}
	err = s.dao.UpdateBlacklist(aID, ctype, update)
	if err != nil {
		log.Error("s.dao.UpdateBlacklist error(%v)", err)
		return
	}
	return
}
