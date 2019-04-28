package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_activitisSQL      = "SELECT id,name,signed_start,signed_end,sign_up,sign_up_start,sign_up_end,object,upload_start,upload_end,win_type,require_items,require_value,statistics_start,statistics_end,bonus_type,bonus_time,progress_frequency,update_page,progress_start,progress_end,progress_sync,bonus_query,bonus_query_start,bonus_query_end FROM creative_activity"
	_upActivitySQL     = "SELECT mid,activity_id,state,success_time FROM up_activity WHERE activity_id = ? AND state > 0 AND is_deleted = 0"
	_activityBonuseSQL = "SELECT bonus_money,ranking FROM activity_bonus WHERE activity_id = ?"
	_upAvInfoSQL       = "SELECT id,av_id,mid,upload_time FROM up_av_info WHERE id > ? ORDER BY id LIMIT ?"
	_archiveInfoSQL    = "SELECT id,av_id,state,likes,share,play,reply,danmu FROM activity_archive_info WHERE id > ? AND activity_id = ? ORDER BY id LIMIT ?"

	_inUpActivitySQL  = "INSERT INTO up_activity(mid,nickname,activity_id,aids,aid_num,ranking,bonus,state,success_time) VALUES %s ON DUPLICATE KEY UPDATE aids=VALUES(aids),aid_num=VALUES(aid_num),ranking=VALUES(ranking),bonus=VALUES(bonus),state=VALUES(state)"
	_updateUpStateSQL = "UPDATE up_activity SET state = ?, bonus = 0, ranking = 0, aids = '' WHERE activity_id = ? AND state = ?"
)

// GetCActivities get activity by query
func (d *Dao) GetCActivities(c context.Context) (acs []*model.CActivity, err error) {
	acs = make([]*model.CActivity, 0)
	rows, err := d.db.Query(c, _activitisSQL)
	if err != nil {
		log.Error("GetCActivities d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ac := &model.CActivity{}
		err = rows.Scan(&ac.ID, &ac.Name, &ac.SignedStart, &ac.SignedEnd, &ac.SignUp, &ac.SignUpStart, &ac.SignUpEnd, &ac.Object, &ac.UploadStart, &ac.UploadEnd, &ac.WinType, &ac.RequireItems, &ac.RequireValue, &ac.StatisticsStart, &ac.StatisticsEnd, &ac.BonusType, &ac.BonusTime, &ac.ProgressFrequency, &ac.UpdatePage, &ac.ProgressStart, &ac.ProgressEnd, &ac.ProgressSync, &ac.BonusQuery, &ac.BonusQuerStart, &ac.BonusQueryEnd)
		if err != nil {
			log.Error("GetCActivities row.Scan error(%v)", err)
			return
		}
		acs = append(acs, ac)
	}
	err = rows.Err()
	return
}

// ListUpActivity get up from up_activity
func (d *Dao) ListUpActivity(c context.Context, id int64) (ups []*model.UpActivity, err error) {
	ups = make([]*model.UpActivity, 0)
	rows, err := d.db.Query(c, _upActivitySQL, id)
	if err != nil {
		log.Error("ListUpActivity d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpActivity{}
		err = rows.Scan(&up.MID, &up.ActivityID, &up.State, &up.SuccessTime)
		if err != nil {
			log.Error("ListUpActivity rows.Scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// GetActivityBonus get activity_bonus by activity_id
func (d *Dao) GetActivityBonus(c context.Context, id int64) (actBonus map[int64]int64, err error) {
	actBonus = make(map[int64]int64)
	rows, err := d.db.Query(c, _activityBonuseSQL, id)
	if err != nil {
		log.Error("GetActivityBonus d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rank, money int64
		err = rows.Scan(&money, &rank)
		if err != nil {
			log.Error("GetActivityBonus rows.Scan error(%v)", err)
			return
		}
		actBonus[rank] = money
	}
	err = rows.Err()
	return
}

// GetAvUploadByMID get up_av_info by mid
func (d *Dao) GetAvUploadByMID(c context.Context, id int64, limit int) (avs []*model.AvUpload, err error) {
	avs = make([]*model.AvUpload, 0)
	rows, err := d.db.Query(c, _upAvInfoSQL, id, limit)
	if err != nil {
		log.Error("GetAvUploadByMID d.dbQuery error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.AvUpload{}
		err = rows.Scan(&av.ID, &av.AvID, &av.MID, &av.UploadTime)
		if err != nil {
			log.Error("GetAvUploadByMID rows.Scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}
	return
}

// GetArchiveInfo get activity archive info
func (d *Dao) GetArchiveInfo(c context.Context, activityID, id int64, limit int) (avs []*model.ArchiveStat, err error) {
	avs = make([]*model.ArchiveStat, 0)
	rows, err := d.db.Query(c, _archiveInfoSQL, id, activityID, limit)
	if err != nil {
		log.Error("GetArchiveInfo d.dbQuery error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.ArchiveStat{}
		err = rows.Scan(&av.ID, &av.AvID, &av.State, &av.Like, &av.Share, &av.Play, &av.Reply, &av.Dm)
		if err != nil {
			log.Error("GetArchiveInfo rows.Scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}
	return
}

// UpdateUpActivityState update up_activity state
func (d *Dao) UpdateUpActivityState(c context.Context, id int64, oldState, newState int) (count int64, err error) {
	res, err := d.db.Exec(c, _updateUpStateSQL, newState, id, oldState)
	if err != nil {
		log.Error("UpdateUpActivityState tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUpActivityBatch insert up_activity
func (d *Dao) InsertUpActivityBatch(c context.Context, vals string) (count int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpActivitySQL, vals))
	if err != nil {
		log.Error("InsertUpActivityBatch d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
