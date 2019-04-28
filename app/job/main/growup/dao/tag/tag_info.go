package tag

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// select
	_allTagInfoSQL    = "SELECT id, start_at, end_at FROM tag_info WHERE is_deleted = 0"
	_tagInfoByDateSQL = "SELECT id,tag,category_id,business_id,adjust_type,ratio,is_common,activity_id,upload_start_time,upload_end_time FROM tag_info WHERE dimension = ? AND business_id = ? AND start_at <= ? AND end_at >= ? AND is_deleted = 0"

	// update
	_updateTagInfoIncomeSQL = "UPDATE tag_info SET total_income = total_income + %d WHERE id = ?"
	_upTagUpsSQL            = "UPDATE tag_info SET ups = ? WHERE id = ?"
)

// AllTagInfo get all tagInfo
func (d *Dao) AllTagInfo(c context.Context) (tagInfos []*model.TagInfo, err error) {
	rows, err := d.db.Query(c, _allTagInfoSQL)
	if err != nil {
		log.Error("AllTagInfo Query(%v), error(%v)", _allTagInfoSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.TagInfo{}
		err = rows.Scan(&t.ID, &t.StartAt, &t.EndAt)
		if err != nil {
			log.Error("AllTagInfo rows scan error(%v)", err)
			return
		}
		tagInfos = append(tagInfos, t)
	}
	return
}

// GetTagInfoByDate get tag_info by date
func (d *Dao) GetTagInfoByDate(c context.Context, dimension int, ctype int, startAt string, endAt string) (tagInfos []*model.TagInfo, err error) {
	rows, err := d.db.Query(c, _tagInfoByDateSQL, dimension, ctype, startAt, endAt)
	if err != nil {
		log.Error("GetTagInfoByDate Query(%v), error(%v)", _tagInfoByDateSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.TagInfo{}
		err = rows.Scan(&t.ID, &t.TagName, &t.CategoryID, &t.BusinessID, &t.AdjustType, &t.Ratio, &t.IsCommon, &t.ActivityID, &t.UploadStartTime, &t.UploadEndTime)
		if err != nil {
			log.Error("GetTagInfoByDate rows scan error(%v)", err)
			return
		}
		tagInfos = append(tagInfos, t)
	}
	return
}

// TxUpdateTagInfoIncome update tag_info income
func (d *Dao) TxUpdateTagInfoIncome(tx *sql.Tx, id, income int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateTagInfoIncomeSQL, income), id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateTagUps update tag info effect ups count
func (d *Dao) UpdateTagUps(c context.Context, tagID int64, ups int) (rows int64, err error) {
	res, err := d.db.Exec(c, _upTagUpsSQL, ups, tagID)
	if err != nil {
		log.Error("dao.UpdateTagUps exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
