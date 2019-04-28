package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_activityNotStart = iota
	_activityStart
	_activityEnd
)

// AddActivity add creative activity
func (s *Service) AddActivity(c context.Context, ac *model.CActivity, creator string) (err error) {
	ac.Creator = creator
	_, err = s.dao.GetActivityByName(c, ac.Name)
	if err == sql.ErrNoRows {
		err = s.insertActivity(c, ac, "", true)
		return
	}
	if err != nil {
		log.Error("s.dao.GetActivityByName error(%v)", err)
		return
	}
	err = fmt.Errorf("activity has exist")
	return
}

func actQueyStmt(name, sort string) string {
	query := ""
	if name != "" {
		query = fmt.Sprintf("WHERE name = '%s'", name)
	}
	if sort != "" {
		query += " ORDER BY "
		query += sort
	}
	return query
}

// ListActivity list activity by query
func (s *Service) ListActivity(c context.Context, name string, from, limit int, sort string) (acs []*model.CActivity, total int, err error) {
	query := actQueyStmt(name, sort)
	total, err = s.dao.ActivityCount(c, query)
	if err != nil {
		log.Error("s.dao.ActivityCount error(%v)", err)
		return
	}
	acs, err = s.dao.GetActivities(c, query, from, limit)
	if err != nil {
		log.Error("s.dao.GetActivities error(%v)", err)
		return
	}
	if len(acs) == 0 {
		return
	}
	ids := make([]int64, len(acs))
	for i := 0; i < len(acs); i++ {
		ids[i] = acs[i].ID
	}
	bonus, err := s.getActivityBonus(c, ids)
	if err != nil {
		log.Error("s.getActivityBonus error(%v)", err)
		return
	}

	now := xtime.Time(time.Now().Unix())
	for _, ac := range acs {
		if now < ac.SignUpStart {
			ac.State = _activityNotStart
		} else if now >= ac.SignUpStart && now <= ac.BonusTime {
			ac.State = _activityStart
		} else if now > ac.BonusTime {
			ac.State = _activityEnd
		}
		ac.BonusMoney = bonus[ac.ID]
		ac.Enrolment, err = s.dao.UpActivityStateCount(c, ac.ID, []int64{1, 2, 3})
		if err != nil {
			log.Error("s.dao.UpActivityStateCount error(%v)", err)
			return
		}
		ac.WinNum, err = s.dao.UpActivityStateCount(c, ac.ID, []int64{2, 3})
		if err != nil {
			log.Error("s.dao.UpActivityStateCount error(%v)", err)
			return
		}
	}
	return
}

func checkSignUp(oldAc, newAc *model.CActivity) bool {
	if oldAc.SignedStart != newAc.SignedStart ||
		oldAc.SignedEnd != newAc.SignedEnd ||
		oldAc.SignUp != newAc.SignUp ||
		oldAc.SignUpStart != newAc.SignUpStart {
		return false
	}
	return true
}

func checkWin(oldAc, newAc *model.CActivity) bool {
	if oldAc.Object != newAc.Object ||
		oldAc.UploadStart != newAc.UploadStart ||
		oldAc.UploadEnd != newAc.UploadEnd ||
		oldAc.WinType != newAc.WinType ||
		oldAc.RequireItems != newAc.RequireItems ||
		oldAc.RequireValue != newAc.RequireValue ||
		oldAc.StatisticsStart != newAc.StatisticsStart {
		return false
	}
	return true
}

func checkBonus(oldAc, newAc *model.CActivity) bool {
	if oldAc.BonusType != newAc.BonusType ||
		oldAc.BonusTime != newAc.BonusTime {
		return false
	}
	return true
}

func checkProgress(oldAc, newAc *model.CActivity) bool {
	if oldAc.ProgressFrequency != newAc.ProgressFrequency ||
		oldAc.UpdatePage != newAc.UpdatePage ||
		oldAc.ProgressStart != newAc.ProgressStart ||
		oldAc.ProgressSync != newAc.ProgressSync {
		return false
	}
	return true
}

func checkOpenBonus(oldAc, newAc *model.CActivity) bool {
	if oldAc.BonusQuery != newAc.BonusQuery ||
		oldAc.BonusQuerStart != newAc.BonusQuerStart {
		return false
	}
	return true
}

// UpdateActivity update creative activity
func (s *Service) UpdateActivity(c context.Context, newAc *model.CActivity) (err error) {
	acs, _, err := s.ListActivity(c, newAc.Name, 0, 1, "")
	if err != nil {
		log.Error("s.ListActivity error(%v)", err)
		return
	}
	if len(acs) == 0 {
		err = fmt.Errorf("activity(%s) not exist", newAc.Name)
		return
	}
	old := acs[0]
	// 报名标准
	signUpStr := "signed_start=VALUES(signed_start),signed_end=VALUES(signed_end),sign_up=VALUES(sign_up),sign_up_start=VALUES(sign_up_start)"
	// 中奖标准
	winStr := "object=VALUES(object),upload_start=VALUES(upload_start),upload_end=VALUES(upload_end),win_type=VALUES(win_type),require_items=VALUES(require_items),require_value=VALUES(require_value),statistics_start=VALUES(statistics_start),statistics_end=VALUES(statistics_end)"
	// 奖金设置
	bonusStr := "bonus_type=VALUES(bonus_type),bonus_time=VALUES(bonus_time)"
	// 进展同步
	progressStr := "progress_frequency=VALUES(progress_frequency),update_page=VALUES(update_page),progress_start=VALUES(progress_start),progress_end=VALUES(progress_end),progress_sync=VALUES(progress_sync)"
	// 开奖查询
	openBonusStr := "bonus_query=VALUES(bonus_query),bonus_query_start=VALUES(bonus_query_start),bonus_query_end=VALUES(bonus_query_end)"
	// others
	otherStr := "background=VALUES(background),win_desc=VALUES(win_desc),unwin_desc=VALUES(unwin_desc),details=VALUES(details)"

	var (
		update      = ""
		updateBonus = false
		now         = xtime.Time(time.Now().Unix())
	)
	switch {
	case now < old.SignUpStart:
		// 报名未时间开始
		update = fmt.Sprintf("%s,%s,%s,%s,%s,%s,sign_up_end=VALUES(sign_up_end)", signUpStr, winStr, bonusStr, progressStr, openBonusStr, otherStr)
		updateBonus = true
	case now >= old.SignUpStart && now <= old.SignUpEnd && now < old.ProgressStart:
		// 报名已开始未结束,进展同步未开始
		if !checkSignUp(old, newAc) || !checkWin(old, newAc) || !checkBonus(old, newAc) {
			err = fmt.Errorf("报名已开始,无法修改报名、中奖、奖金相关内容,请检查修改项")
			return
		}
		update = fmt.Sprintf("sign_up_end=VALUES(sign_up_end),%s,%s,%s", progressStr, openBonusStr, otherStr)
	case now > old.SignUpEnd && now >= old.ProgressStart && now <= old.ProgressEnd:
		// 报名已结束,进展同步开始未结束
		if !checkSignUp(old, newAc) || !checkWin(old, newAc) || !checkBonus(old, newAc) || !checkProgress(old, newAc) {
			err = fmt.Errorf("报名已结束,进展同步开始未结束,无法修改报名、中奖、奖金、进展相关内容,请检查修改项")
			return
		}
		update = fmt.Sprintf("progress_end=VALUES(progress_end),%s,%s", openBonusStr, otherStr)
	case now > old.ProgressEnd && now < old.BonusQueryEnd:
		// 进展同步已结束,开奖查询未结束
		if !checkSignUp(old, newAc) || !checkWin(old, newAc) || !checkBonus(old, newAc) || !checkProgress(old, newAc) || !checkOpenBonus(old, newAc) || old.ProgressEnd != newAc.ProgressEnd {
			err = fmt.Errorf("进展同步已结束,开奖查询未结束,无法修改报名、中奖、奖金、进展、开奖相关内容,请检查修改项")
			return
		}
		update = fmt.Sprintf("bonus_query_end=VALUES(bonus_query_end),%s", otherStr)
	default:
		err = fmt.Errorf("不符合任何修改时间段,没有任何修改")
		return
	}
	update = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", update)
	err = s.insertActivity(c, newAc, update, updateBonus)
	return
}

func (s *Service) getActivityBonus(c context.Context, ids []int64) (bm map[int64][]int64, err error) {
	bm = make(map[int64][]int64)
	bonus, err := s.dao.GetActivityBonus(c, ids)
	if err != nil {
		return
	}
	sort.Slice(bonus, func(i, j int) bool {
		return bonus[i].Rank < bonus[j].Rank
	})
	for i := 0; i < len(bonus); i++ {
		id := bonus[i].ID
		if _, ok := bm[id]; !ok {
			bm[id] = make([]int64, 0)
		}
		bm[id] = append(bm[id], bonus[i].Money)
	}
	return
}

func (s *Service) insertActivity(c context.Context, ac *model.CActivity, updateVal string, updateBonus bool) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	// insert activity
	if _, err = s.dao.TxInsertActivity(tx, ac, updateVal); err != nil {
		log.Error("s.dao.TxInsertActivity error(%v)", err)
		return
	}
	id, err := s.dao.TxGetActivityByName(tx, ac.Name)
	if err != nil {
		log.Error("s.dao.GetActivityByName error(%v)", err)
		return
	}

	// is update bonus
	if updateBonus && len(ac.BonusMoney) > 0 {
		bonus := make([]*model.BonusRank, 0)
		if ac.WinType == 1 {
			bonus = append(bonus, &model.BonusRank{ID: id, Rank: 0, Money: ac.BonusMoney[0]})
		} else if ac.WinType == 2 {
			for i := 0; i < len(ac.BonusMoney); i++ {
				bonus = append(bonus, &model.BonusRank{ID: id, Rank: i + 1, Money: ac.BonusMoney[i]})
			}
		}
		// insert bonus money
		if err = s.txInsertActivityBonus(tx, bonus); err != nil {
			log.Error("s.TxInsertBonusRank error(%v)", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

func (s *Service) txInsertActivityBonus(tx *sql.Tx, bonus []*model.BonusRank) (err error) {
	var buf bytes.Buffer
	for _, row := range bonus {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.ID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.Rank))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Money, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals := buf.String()
	buf.Reset()
	_, err = s.dao.TxInsertActivityBonusBatch(tx, vals)
	return
}

// ListActivitySignUp list activity who sign up
func (s *Service) ListActivitySignUp(c context.Context, activityID int64, from, limit int) (ups []*model.UpActivity, total int, err error) {
	total, err = s.dao.UpActivityStateCount(c, activityID, []int64{1, 2, 3})
	if err != nil {
		log.Error("s.dao.UpActivityStateCount error(%v)", err)
		return
	}
	ups, err = s.dao.ListUpActivity(c, activityID, from, limit)
	if err != nil {
		log.Error("s.dao.ListUpActivity error(%v)", err)
	}
	return
}

// ListActivityWinners list activity winners
func (s *Service) ListActivityWinners(c context.Context, activityID, mid int64, from, limit int) (ups []*model.UpActivity, total int, err error) {
	total, err = s.dao.UpActivityStateCount(c, activityID, []int64{2, 3})
	if err != nil {
		log.Error("s.dao.UpActivityStateCount error(%v)", err)
		return
	}
	ups, err = s.dao.ListUpActivitySuccess(c, activityID, mid, from, limit)
	if err != nil {
		log.Error("s.dao.ListUpActivity error(%v)", err)
		return
	}
	if mid != 0 {
		total = len(ups)
	}
	return
}

// ActivityAward activity award
func (s *Service) ActivityAward(c context.Context, activityID int64, activityName string, date, statisticsEnd xtime.Time, creator string) (err error) {
	if xtime.Time(time.Now().Unix()) <= statisticsEnd {
		err = fmt.Errorf("统计阶段未结束，不能发奖")
		return
	}
	ups, err := s.listUpActivity(c, activityID)
	if err != nil {
		log.Error("s.listUpActivity error(%v)", err)
		return
	}
	if len(ups) == 0 {
		return
	}
	rankMID := make(map[int][]int64)
	rankMoney := make(map[int]int64)
	for _, up := range ups {
		if up.State != 2 {
			continue
		}
		rank := up.Rank
		rankMoney[rank] = up.Bonus
		if _, ok := rankMID[rank]; !ok {
			rankMID[rank] = make([]int64, 0)
		}
		rankMID[rank] = append(rankMID[rank], up.MID)
	}

	// insert to tag
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	for rank, money := range rankMoney {
		mids, ok := rankMID[rank]
		if !ok || len(mids) == 0 {
			continue
		}
		tagName := fmt.Sprintf("act-%s-%d", activityName, rank)
		err = s.addActivityUpTag(tx, money, creator, tagName, mids, date)
		if err != nil {
			log.Error("s.addActivityUpTag error(%v)", err)
			return
		}
		// update mids state
		if _, err = s.dao.TxUpdateUpActivityState(tx, activityID, mids, 2, 3); err != nil {
			tx.Rollback()
			log.Error("s.dao.TxUpdateUpActivityState error(%v)", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
	}
	return
}

func (s *Service) addActivityUpTag(tx *sql.Tx, money int64, creator, tagName string, mids []int64, date xtime.Time) (err error) {
	tag := &model.TagInfo{
		Tag:             tagName,
		Creator:         creator,
		Dimension:       1,
		StartTime:       date,
		EndTime:         date,
		AdjustType:      1,
		Ratio:           int(money),
		UploadStartTime: date,
		UploadEndTime:   date,
	}
	if _, err = s.dao.TxInsertTag(tx, tag); err != nil {
		tx.Rollback()
		log.Error("s.dao.TxInsertTag error(%v)", err)
		return
	}
	tagID, err := s.dao.TxGetTagInfoByName(tx, tagName, 1, 0, 0)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.TxGetTagInfoByName error(%v)", err)
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
	return
}

func (s *Service) listUpActivity(c context.Context, activityID int64) (ups []*model.UpActivity, err error) {
	from, limit := 0, 2000
	ups = make([]*model.UpActivity, 0)
	for {
		var up []*model.UpActivity
		up, err = s.dao.ListUpActivity(c, activityID, from, limit)
		if err != nil {
			return
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
		from += limit
	}
	return
}
