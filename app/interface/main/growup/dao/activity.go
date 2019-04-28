package dao

import (
	"context"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	// "go-common/library/database/sql"
)

const (
	_activitySQL   = "SELECT id,name,signed_start,signed_end,sign_up_start,sign_up_end,sign_up,win_type,progress_start,progress_end,progress_sync,update_page,bonus_query,bonus_query_start,bonus_query_end,background,win_desc,unwin_desc,details FROM creative_activity WHERE id = ?"
	_upActivitySQL = "SELECT mid,nickname,ranking,state FROM up_activity WHERE activity_id = ? AND is_deleted = 0"

	_inUpActivitySQL = "INSERT INTO up_activity(mid,nickname,activity_id,state,sign_up_time) VALUES(?,?,?,?,?)"
)

// GetActivity get activity by query
func (d *Dao) GetActivity(c context.Context, id int64) (ac *model.CActivity, err error) {
	ac = &model.CActivity{}
	err = d.db.QueryRow(c, _activitySQL, id).Scan(&ac.ID, &ac.Name, &ac.SignedStart, &ac.SignedEnd, &ac.SignUpStart, &ac.SignUpEnd, &ac.SignUp, &ac.WinType, &ac.ProgressStart, &ac.ProgressEnd, &ac.ProgressSync, &ac.UpdatePage, &ac.BonusQuery, &ac.BonusQueryStart, &ac.BonusQueryEnd, &ac.Background, &ac.WinDesc, &ac.UnwinDesc, &ac.Details)
	return
}

// ListUpActivity get up from up_activity
func (d *Dao) ListUpActivity(c context.Context, id int64) (ups []*model.UpBonus, err error) {
	ups = make([]*model.UpBonus, 0)
	rows, err := d.db.Query(c, _upActivitySQL, id)
	if err != nil {
		log.Error("ListUpActivity d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpBonus{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.Rank, &up.State)
		if err != nil {
			log.Error("ListUpActivity rows.Scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// SignUpActivity up sign up activity
func (d *Dao) SignUpActivity(c context.Context, up *model.UpBonus) (rows int64, err error) {
	res, err := d.db.Exec(c, _inUpActivitySQL, up.MID, up.Nickname, up.ActivityID, up.State, up.SignUpTime)
	if err != nil {
		log.Error("SignUpActivity.Exec(%s) error(%v)", _inUpActivitySQL, err)
		return
	}
	return res.RowsAffected()
}
