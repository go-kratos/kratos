package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_activityByNameSQL       = "SELECT id FROM creative_activity WHERE name = ?"
	_activityCountSQL        = "SELECT count(*) FROM creative_activity %s"
	_activitisSQL            = "SELECT id,name,signed_start,signed_end,sign_up,sign_up_start,sign_up_end,object,upload_start,upload_end,win_type,require_items,require_value,statistics_start,statistics_end,bonus_type,bonus_time,progress_frequency,update_page,progress_start,progress_end,progress_sync,bonus_query,bonus_query_start,bonus_query_end,background,win_desc,unwin_desc,details FROM creative_activity %s LIMIT ?,?"
	_activityBonusSQL        = "SELECT activity_id,ranking,bonus_money FROM activity_bonus WHERE activity_id IN (%s)"
	_upActivityStateCountSQL = "SELECT count(*) FROM up_activity WHERE activity_id = ? AND state IN (%s)"
	_upActivitySQL           = "SELECT mid,nickname,sign_up_time,bonus,ranking,state FROM up_activity WHERE activity_id = ? AND state != 0 AND is_deleted = 0 LIMIT ?,?"
	_upActivitySuccessSQL    = "SELECT mid,nickname,aids,aid_num,bonus,success_time,state FROM up_activity WHERE activity_id = ? AND state IN (2,3) %s AND is_deleted = 0 LIMIT ?,?"

	_txInsertActivitySQL  = "INSERT INTO creative_activity(name,creator,signed_start,signed_end,sign_up,sign_up_start,sign_up_end,object,upload_start,upload_end,win_type,require_items,require_value,statistics_start,statistics_end,bonus_type,bonus_time,progress_frequency,update_page,progress_start,progress_end,progress_sync,bonus_query,bonus_query_start,bonus_query_end,background,win_desc,unwin_desc,details) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) %s"
	_txInsertBonusRankSQL = "INSERT INTO activity_bonus(activity_id,ranking,bonus_money) VALUES %s ON DUPLICATE KEY UPDATE activity_id=VALUES(activity_id),ranking=VALUES(ranking),bonus_money=VALUES(bonus_money)"

	_updateUpActivityStateSQL = "UPDATE up_activity SET state = ? WHERE activity_id = ? AND mid IN (%s) AND state = ?"
)

// GetActivityByName get activity by name
func (d *Dao) GetActivityByName(c context.Context, name string) (id int64, err error) {
	err = d.rddb.QueryRow(c, _activityByNameSQL, name).Scan(&id)
	return
}

// ActivityCount get activity count by query
func (d *Dao) ActivityCount(c context.Context, query string) (count int, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_activityCountSQL, query)).Scan(&count)
	return
}

// GetActivities get activity by query
func (d *Dao) GetActivities(c context.Context, query string, from, limit int) (acs []*model.CActivity, err error) {
	acs = make([]*model.CActivity, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_activitisSQL, query), from, limit)
	if err != nil {
		log.Error("GetActivities d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ac := &model.CActivity{}
		err = rows.Scan(&ac.ID, &ac.Name, &ac.SignedStart, &ac.SignedEnd, &ac.SignUp, &ac.SignUpStart, &ac.SignUpEnd, &ac.Object, &ac.UploadStart, &ac.UploadEnd, &ac.WinType, &ac.RequireItems, &ac.RequireValue, &ac.StatisticsStart, &ac.StatisticsEnd, &ac.BonusType, &ac.BonusTime, &ac.ProgressFrequency, &ac.UpdatePage, &ac.ProgressStart, &ac.ProgressEnd, &ac.ProgressSync, &ac.BonusQuery, &ac.BonusQuerStart, &ac.BonusQueryEnd, &ac.Background, &ac.WinDesc, &ac.UnwinDesc, &ac.Details)
		if err != nil {
			log.Error("GetActivities row.Scan error(%v)", err)
			return
		}
		acs = append(acs, ac)
	}
	err = rows.Err()
	return
}

// TxGetActivityByName get activity by name
func (d *Dao) TxGetActivityByName(tx *sql.Tx, name string) (id int64, err error) {
	err = tx.QueryRow(_activityByNameSQL, name).Scan(&id)
	return
}

// TxInsertActivity tx insert into creative_activity
func (d *Dao) TxInsertActivity(tx *sql.Tx, ac *model.CActivity, update string) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_txInsertActivitySQL, update), ac.Name, ac.Creator, ac.SignedStart, ac.SignedEnd, ac.SignUp, ac.SignUpStart, ac.SignUpEnd, ac.Object, ac.UploadStart, ac.UploadEnd, ac.WinType, ac.RequireItems, ac.RequireValue, ac.StatisticsStart, ac.StatisticsEnd, ac.BonusType, ac.BonusTime, ac.ProgressFrequency, ac.UpdatePage, ac.ProgressStart, ac.ProgressEnd, ac.ProgressSync, ac.BonusQuery, ac.BonusQuerStart, ac.BonusQueryEnd, ac.Background, ac.WinDesc, ac.UnwinDesc, ac.Details)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetActivityBonus get activity_bonus by id
func (d *Dao) GetActivityBonus(c context.Context, ids []int64) (brs []*model.BonusRank, err error) {
	brs = make([]*model.BonusRank, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_activityBonusSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("GetActivityBonus d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		br := &model.BonusRank{}
		err = rows.Scan(&br.ID, &br.Rank, &br.Money)
		if err != nil {
			log.Error("GetActivityBonus row.Scan error(%v)", err)
			return
		}
		brs = append(brs, br)
	}
	err = rows.Err()
	return
}

// TxInsertActivityBonusBatch tx insert into activity_bonus
func (d *Dao) TxInsertActivityBonusBatch(tx *sql.Tx, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_txInsertBonusRankSQL, vals))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpActivityStateCount get up activity state count
func (d *Dao) UpActivityStateCount(c context.Context, id int64, states []int64) (count int, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_upActivityStateCountSQL, xstr.JoinInts(states)), id).Scan(&count)
	return
}

// ListUpActivity list up_activity where state != 0
func (d *Dao) ListUpActivity(c context.Context, id int64, from, limit int) (ups []*model.UpActivity, err error) {
	ups = make([]*model.UpActivity, 0)
	rows, err := d.rddb.Query(c, _upActivitySQL, id, from, limit)
	if err != nil {
		log.Error("ListUpActivity d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpActivity{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.SignUpTime, &up.Bonus, &up.Rank, &up.State)
		if err != nil {
			log.Error("ListUpActivity row.Scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// ListUpActivitySuccess list up_activity where state != 0
func (d *Dao) ListUpActivitySuccess(c context.Context, id, mid int64, from, limit int) (ups []*model.UpActivity, err error) {
	ups = make([]*model.UpActivity, 0)
	query := ""
	if mid != 0 {
		query = fmt.Sprintf("AND mid = %d", mid)
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upActivitySuccessSQL, query), id, from, limit)
	if err != nil {
		log.Error("ListUpActivitySuccess d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpActivity{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.AIDs, &up.AIDNum, &up.Bonus, &up.SuccessTime, &up.State)
		if err != nil {
			log.Error("ListUpActivitySuccess row.Scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// UpdateUpActivityState update up state
func (d *Dao) UpdateUpActivityState(c context.Context, activityID int64, mids []int64, oldState, newState int) (rows int64, err error) {
	if oldState == newState {
		return
	}
	res, err := d.rddb.Exec(c, fmt.Sprintf(_updateUpActivityStateSQL, xstr.JoinInts(mids)), newState, activityID, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxUpdateUpActivityState tx update up state
func (d *Dao) TxUpdateUpActivityState(tx *sql.Tx, activityID int64, mids []int64, oldState, newState int) (rows int64, err error) {
	if oldState == newState {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_updateUpActivityStateSQL, xstr.JoinInts(mids)), newState, activityID, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
