package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_openAwardsSQL   = "SELECT award_id,award_name,cycle_start FROM special_award WHERE open_status=2"
	_joinedAwardsSQL = "SELECT award_id,award_name,cycle_start,cycle_end,open_status FROM special_award WHERE award_id IN (%s)"

	_awardScheSQL     = "SELECT award_name,cycle_start,cycle_end,announce_date FROM special_award WHERE award_id=?"
	_awardResourceSQL = "SELECT resource_type,content,resource_index FROM special_award_resource WHERE award_id=? AND is_deleted=0"
	_awardWinnerSQL   = "SELECT mid FROM special_award_winner WHERE award_id=? AND is_deleted=0"

	_winnerIDsSQL      = "SELECT award_id,prize_id FROM special_award_winner WHERE mid=? AND is_deleted=0"
	_winnerDivisionSQL = "SELECT award_id,division_name FROM special_award_winner WHERE mid=? AND is_deleted=0"

	_awardBonusSQL = "SELECT bonus FROM special_award_prize WHERE award_id=? AND prize_id=?"

	_inSpecialAwardRecordSQL = "INSERT INTO special_award_record(mid,award_id) VALUES(?,?)"
	_joinedCountSQL          = "SELECT count(*) FROM special_award_record WHERE mid=? AND award_id=?"

	_specialAwardsSQL   = "SELECT award_id, award_name, cycle_start, cycle_end, announce_date FROM special_award WHERE display_status = 2 AND is_deleted = 0"
	_awardDivisionSQL   = "SELECT division_name FROM special_award_division WHERE award_id = ? AND is_deleted = 0"
	_awardWinRecordSQL  = "SELECT award_id FROM special_award_winner WHERE mid = ? AND is_deleted = 0"
	_awardJoinRecordSQL = "SELECT award_id FROM special_award_record WHERE mid = ? AND is_deleted = 0"
)

// PastAwards get award basic info
func (d *Dao) PastAwards(c context.Context) (as []*model.SimpleSpecialAward, err error) {
	rows, err := d.db.Query(c, _openAwardsSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	as = make([]*model.SimpleSpecialAward, 0)
	for rows.Next() {
		a := &model.SimpleSpecialAward{}
		err = rows.Scan(&a.AwardID, &a.AwardName, &a.CycleStart)
		if err != nil {
			return
		}
		as = append(as, a)
	}
	err = rows.Err()
	return
}

// JoinedSpecialAwards get joined awards
func (d *Dao) JoinedSpecialAwards(c context.Context, awardIDs []int64) (sas []*model.SpecialAward, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_joinedAwardsSQL, xstr.JoinInts(awardIDs)))
	if err != nil {
		return
	}
	defer rows.Close()
	sas = make([]*model.SpecialAward, 0)
	for rows.Next() {
		sa := &model.SpecialAward{}
		err = rows.Scan(&sa.AwardID, &sa.AwardName, &sa.CycleStart, &sa.CycleEnd, &sa.OpenStatus)
		if err != nil {
			log.Error("row scan error(%v)", err)
			return
		}
		sas = append(sas, sa)
	}
	err = rows.Err()
	return
}

// GetAwardSchedule get special award by award id
func (d *Dao) GetAwardSchedule(c context.Context, awardID int64) (award *model.SpecialAward, err error) {
	award = &model.SpecialAward{}
	row := d.db.QueryRow(c, _awardScheSQL, awardID)
	if err = row.Scan(&award.AwardName, &award.CycleStart, &award.CycleEnd, &award.AnnounceDate); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// GetResources get special award resource map[type]map[index]content
func (d *Dao) GetResources(c context.Context, awardID int64) (res map[int]map[int]string, err error) {
	rows, err := d.db.Query(c, _awardResourceSQL, awardID)
	if err != nil {
		log.Error("GetResources d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int]map[int]string)
	for rows.Next() {
		var rtype, index int
		var content string
		err = rows.Scan(&rtype, &content, &index)
		if err != nil {
			log.Error("GetDivisions rows.Scan error(%v)", err)
			return
		}
		if cs, ok := res[rtype]; ok {
			cs[index] = content
		} else {
			m := make(map[int]string)
			m[index] = content
			res[rtype] = m
		}
	}
	err = rows.Err()
	return
}

// GetWinners get winner mids by award_id
func (d *Dao) GetWinners(c context.Context, awardID int64) (mids []int64, err error) {
	rows, err := d.db.Query(c, _awardWinnerSQL, awardID)
	if err != nil {
		log.Error("GetWinners d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	mids = make([]int64, 0)
	for rows.Next() {
		var mid int64
		err = rows.Scan(&mid)
		if err != nil {
			log.Error("GetWinners rows.Scan error(%v)", err)
			return
		}
		mids = append(mids, mid)
	}
	err = rows.Err()
	return
}

// AwardIDsByWinner get winner's record, am map[award_id]prize_id
func (d *Dao) AwardIDsByWinner(c context.Context, mid int64) (am map[int64]int64, err error) {
	rows, err := d.db.Query(c, _winnerIDsSQL, mid)
	if err != nil {
		return
	}
	defer rows.Close()
	am = make(map[int64]int64)
	for rows.Next() {
		var awardID, prizeID int64
		err = rows.Scan(&awardID, &prizeID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		am[awardID] = prizeID
	}
	err = rows.Err()
	return
}

// DivisionName division name by mid, names map[award_id]division_name
func (d *Dao) DivisionName(c context.Context, mid int64) (names map[int64]string, err error) {
	rows, err := d.db.Query(c, _winnerDivisionSQL, mid)
	if err != nil {
		return
	}
	defer rows.Close()
	names = make(map[int64]string)
	for rows.Next() {
		var awardID int64
		var name string
		err = rows.Scan(&awardID, &name)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		names[awardID] = name
	}
	err = rows.Err()
	return
}

// JoinedCount get signed up count by mid, award_id
func (d *Dao) JoinedCount(c context.Context, mid, awardID int64) (count int64, err error) {
	row := d.db.QueryRow(c, _joinedCountSQL, mid, awardID)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// AwardBonus get bonus by award_id and prize_id
func (d *Dao) AwardBonus(c context.Context, awardID, prizeID int64) (bonus int64, err error) {
	row := d.db.QueryRow(c, _awardBonusSQL, awardID, prizeID)
	if err = row.Scan(&bonus); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// AddToAwardRecord add to award record
func (d *Dao) AddToAwardRecord(c context.Context, mid, awardID int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _inSpecialAwardRecordSQL, mid, awardID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetSpecialAwards get all display sepcial award
func (d *Dao) GetSpecialAwards(c context.Context) (awards []*model.SpecialAward, err error) {
	awards = make([]*model.SpecialAward, 0)
	rows, err := d.db.Query(c, _specialAwardsSQL)
	if err != nil {
		log.Error("GetSpecialAwards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sa := &model.SpecialAward{}
		err = rows.Scan(&sa.AwardID, &sa.AwardName, &sa.CycleStart, &sa.CycleEnd, &sa.AnnounceDate)
		if err != nil {
			log.Error("GetSpecialAwards rows.Scan error(%v)", err)
			return
		}
		awards = append(awards, sa)
	}
	err = rows.Err()
	return
}

// GetSpecialAwardDivision get special award division name
func (d *Dao) GetSpecialAwardDivision(c context.Context, awardID int64) (divisions []string, err error) {
	divisions = make([]string, 0)
	rows, err := d.db.Query(c, _awardDivisionSQL, awardID)
	if err != nil {
		log.Error("GetSpecialAwardDivision d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Error("GetSpecialAwardDivision rows.Scan error(%v)", err)
			return
		}
		divisions = append(divisions, name)
	}
	err = rows.Err()
	return
}

// GetAwardWinRecord get award win record
func (d *Dao) GetAwardWinRecord(c context.Context, mid int64) (awardIDs map[int64]bool, err error) {
	awardIDs = make(map[int64]bool)
	rows, err := d.db.Query(c, _awardWinRecordSQL, mid)
	if err != nil {
		log.Error("GetAwardWinRecord d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Error("GetAwardWinRecord rows.Scan error(%v)", err)
			return
		}
		awardIDs[id] = true
	}
	err = rows.Err()
	return
}

// GetAwardJoinRecord get award join record
func (d *Dao) GetAwardJoinRecord(c context.Context, mid int64) (awardIDs map[int64]bool, err error) {
	awardIDs = make(map[int64]bool)
	rows, err := d.db.Query(c, _awardJoinRecordSQL, mid)
	if err != nil {
		log.Error("GetAwardJoinRecord d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Error("GetAwardJoinRecord rows.Scan error(%v)", err)
			return
		}
		awardIDs[id] = true
	}
	err = rows.Err()
	return
}
