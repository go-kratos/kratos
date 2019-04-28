package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_insertTagSQL = "INSERT INTO tag_info(tag,category_id,business_id,start_at,end_at,creator,ratio,activity_id,icon,dimension,adjust_type,upload_start_time,upload_end_time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"

	// insert tag_up_info
	_inTagUpInfoSQL = "INSERT INTO tag_up_info(tag_id,mid,is_deleted) VALUES (?,?,?) ON DUPLICATE KEY UPDATE is_deleted=VALUES(is_deleted)"

	// update
	_updateTagInfoSQL     = "UPDATE tag_info SET tag=?,category_id=?,business_id=?,start_at=?,end_at=?,ratio=?,activity_id=?,icon=?,dimension=?,adjust_type=?,upload_start_time=?,upload_end_time=? WHERE id = ?"
	_updatetagStateSQL    = "UPDATE tag_info SET is_deleted = ? WHERE id = ?"
	_updateTagComSQL      = "UPDATE tag_info SET is_common = ? WHERE id = ?"
	_updateTagActivitySQL = "UPDATE tag_info set activity_id = ? WHERE id = ?"

	// count(*)
	_tagCntSQL   = "SELECT count(*) FROM tag_info %s"
	_tagUpMIDSQL = "SELECT mid FROM tag_up_info WHERE tag_id = ? AND is_deleted = ?"

	// select
	_tagInfoSQL       = "SELECT id, tag, category_id, business_id, start_at, end_at, ctime, creator, ratio, adjust_type, is_common, activity_id, icon, is_deleted FROM tag_info WHERE id = ?"
	_tagInfoByNameSQL = "SELECT id FROM tag_info WHERE tag = ? AND dimension = ? AND category_id = ? AND business_id = ?"
	_tagInfosSQL      = "SELECT id, tag, dimension, category_id, business_id, total_income, start_at, end_at, ctime, creator, ratio, adjust_type, is_common, is_deleted, activity_id, icon, upload_start_time, upload_end_time, ups FROM tag_info %s LIMIT ?,?"
	_getNicknameSQL   = "SELECT nick_name FROM up_category_info WHERE mid = ? AND is_deleted = 0"
)

// UpdateTagState mode tag is_deleted.
func (d *Dao) UpdateTagState(c context.Context, tagID int, isDeleted int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updatetagStateSQL, isDeleted, tagID)
	if err != nil {
		log.Error("dao.UpdateTagState error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxInsertTagUpInfo insert tag mid into tag_up_info
func (d *Dao) TxInsertTagUpInfo(tx *sql.Tx, tagID, mid int64, isDeleted int) (rows int64, err error) {
	res, err := tx.Exec(_inTagUpInfoSQL, tagID, mid, isDeleted)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertTagUpInfo insert tag mid into tag_up_info
func (d *Dao) InsertTagUpInfo(c context.Context, tagID, mid int64, isDeleted int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inTagUpInfoSQL, tagID, mid, isDeleted)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateTagCom update tag is_common
func (d *Dao) UpdateTagCom(c context.Context, tagID int, isCommon int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updateTagComSQL, isCommon, tagID)
	if err != nil {
		log.Error("d.rddb.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertTag insert tag_info.
func (d *Dao) InsertTag(c context.Context, tag *model.TagInfo) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _insertTagSQL, tag.Tag, tag.Category, tag.Business, tag.StartTime, tag.EndTime, tag.Creator, tag.Ratio, tag.ActivityID, tag.Icon, tag.Dimension, tag.AdjustType, tag.UploadStartTime, tag.UploadEndTime)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxInsertTag insert tag_info.
func (d *Dao) TxInsertTag(tx *sql.Tx, tag *model.TagInfo) (rows int64, err error) {
	res, err := tx.Exec(_insertTagSQL, tag.Tag, tag.Category, tag.Business, tag.StartTime, tag.EndTime, tag.Creator, tag.Ratio, tag.ActivityID, tag.Icon, tag.Dimension, tag.AdjustType, tag.UploadStartTime, tag.UploadEndTime)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateTagInfo update tag_info.
func (d *Dao) UpdateTagInfo(c context.Context, t *model.TagInfo) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updateTagInfoSQL, t.Tag, t.Category, t.Business, t.StartTime, t.EndTime, t.Ratio, t.ActivityID, t.Icon, t.Dimension, t.AdjustType, t.UploadStartTime, t.UploadEndTime, t.ID)
	if err != nil {
		log.Error("growup-admin dao.UpdateTag error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetTagInfo get tag info by tag id.
func (d *Dao) GetTagInfo(c context.Context, tagID int) (info *model.TagInfo, err error) {
	info = new(model.TagInfo)
	row := d.rddb.QueryRow(c, _tagInfoSQL, tagID)
	err = row.Scan(&info.ID, &info.Tag, &info.Category, &info.Business, &info.StartTime, &info.EndTime, &info.CreateTime, &info.Creator, &info.Ratio, &info.AdjustType, &info.IsCommon, &info.ActivityID, &info.Icon, &info.IsDeleted)
	return
}

// GetTagInfoByName get tag info by tag name.
func (d *Dao) GetTagInfoByName(c context.Context, tag string, dimension, category, business int) (id int64, err error) {
	row := d.rddb.QueryRow(c, _tagInfoByNameSQL, tag, dimension, category, business)
	err = row.Scan(&id)
	return
}

// TxGetTagInfoByName get tag info by tag name.
func (d *Dao) TxGetTagInfoByName(tx *sql.Tx, tag string, dimension, category, business int) (id int64, err error) {
	row := tx.QueryRow(_tagInfoByNameSQL, tag, dimension, category, business)
	err = row.Scan(&id)
	return
}

// TagsCount get tag count
func (d *Dao) TagsCount(c context.Context, query string) (count int, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_tagCntSQL, query)).Scan(&count)
	return
}

// GetTagInfos get tag_infos by query
func (d *Dao) GetTagInfos(c context.Context, query string, from, limit int) (tagInfos []*model.TagInfo, err error) {
	tagInfos = make([]*model.TagInfo, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_tagInfosSQL, query), from, limit)
	if err != nil {
		log.Error("d.rddb.Query query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.TagInfo{}
		if err = rows.Scan(&a.ID, &a.Tag, &a.Dimension, &a.Category, &a.Business, &a.TotalIncome, &a.StartTime, &a.EndTime, &a.CreateTime, &a.Creator, &a.Ratio, &a.AdjustType, &a.IsCommon, &a.IsDeleted, &a.ActivityID, &a.Icon, &a.UploadStartTime, &a.UploadEndTime, &a.UpCount); err != nil {
			log.Error("row.Scan error (%v)", err)
			return
		}
		tagInfos = append(tagInfos, a)
	}
	return
}

// GetNickname get nickname
func (d *Dao) GetNickname(c context.Context, mid int64) (nickname string, err error) {
	err = d.rddb.QueryRow(c, _getNicknameSQL, mid).Scan(&nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.rddb.QueryRow.Scan error(%v)", err)
	}
	return
}

// GetTagUpInfoMID get all tag_up_info mids
func (d *Dao) GetTagUpInfoMID(c context.Context, tagID int64, isDeleted int) (mids map[int64]int, err error) {
	mids = make(map[int64]int)
	rows, err := d.rddb.Query(c, _tagUpMIDSQL, tagID, isDeleted)
	if err != nil {
		log.Error("d.rddb.Query query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("row.Scan error (%v)", err)
			return
		}
		mids[mid] = 1
	}
	return
}

// UpdateTagActivity update tag activity_id.
func (d *Dao) UpdateTagActivity(c context.Context, tagID, activityID int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updateTagActivitySQL, activityID, tagID)
	if err != nil {
		log.Error("dao.UpdateTagActivity error(%v)", err)
		return
	}
	return res.RowsAffected()
}
