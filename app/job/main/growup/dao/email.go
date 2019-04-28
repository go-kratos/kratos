package dao

import (
	"context"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_totalIncome = "SELECT id, av_id, mid, income, total_income, is_deleted FROM av_income WHERE date = ? AND id > ? ORDER BY id LIMIT ?"
	_getAV       = "SELECT av_id, income, total_income FROM av_income WHERE date = ?"
	_tag         = "SELECT id, tag, category_id, is_common FROM tag_info WHERE is_deleted = 0 AND start_at <= ? AND end_at >= ?"
	_ncMID       = "SELECT mid FROM tag_up_info WHERE is_deleted = 0 AND tag_id = ?"
	_commonAV    = "SELECT av_id, income, total_income, is_deleted FROM av_income WHERE tag_id = ? AND date = ?"
	_midAV       = "SELECT av_id, income, total_income, is_deleted FROM av_income WHERE mid = ? AND date = ? AND tag_id = ?"
)

// TotalIncome get date totalincome,upcount, avcount
func (d *Dao) TotalIncome(c context.Context, date time.Time, from, limit int64) (infos []*model.IncomeInfo, err error) {
	rows, err := d.db.Query(c, _totalIncome, date, from, limit)
	if err != nil {
		log.Error("growup-job dao.TotalIncome error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.IncomeInfo{}
		err = rows.Scan(&info.ID, &info.AVID, &info.MID, &info.Income, &info.TotalIncome, &info.IsDeleted)
		if err != nil {
			log.Error("growup-job dao.TotalIncome rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// GetAV get av income.
func (d *Dao) GetAV(c context.Context, date time.Time) (infos []*model.IncomeInfo, err error) {
	rows, err := d.db.Query(c, _getAV, date)
	if err != nil {
		log.Error("growup-job dao.GetAV error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.IncomeInfo{}
		err = rows.Scan(&info.AVID, &info.Income, &info.TotalIncome)
		if err != nil {
			log.Error("growup-job dao.GetAV rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// GetTag get tag.
func (d *Dao) GetTag(c context.Context, date time.Time) (infos []*model.TagInfo, err error) {
	infos = make([]*model.TagInfo, 0)
	rows, err := d.db.Query(c, _tag, date, date)
	if err != nil {
		log.Error("growup-job dao.GetTag error(%v),sql(%s)", err, _tag)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.TagInfo{}
		err = rows.Scan(&info.ID, &info.Tag, &info.Category, &info.IsCommon)
		if err != nil {
			log.Error("growup-job dao.GetTag rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// GetMID get tag mid.
func (d *Dao) GetMID(c context.Context, TagID int64) (infos []*model.MIDInfo, err error) {
	infos = make([]*model.MIDInfo, 0)
	rows, err := d.db.Query(c, _ncMID, TagID)
	if err != nil {
		log.Error("growup-job dao.GetMID error(%v),sql(%s)", err, _ncMID)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.MIDInfo{}
		err = rows.Scan(&info.MID)
		if err != nil {
			log.Error("growup-job dao.GetMID rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// TagToAV common tag av.
func (d *Dao) TagToAV(c context.Context, category int, date time.Time) (infos []*model.TagInfo, err error) {
	infos = make([]*model.TagInfo, 0)
	rows, err := d.db.Query(c, _commonAV, category, date)
	if err != nil {
		log.Error("growup-job dao.TagToAV error(%v),sql(%s)", err, _commonAV)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.TagInfo{}
		err = rows.Scan(&info.AVID, &info.Income, &info.TotalIncome, &info.IsDeleted)
		if err != nil {
			log.Error("growup-job dao.TagToAV rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// MIDToAV no common tag av.
func (d *Dao) MIDToAV(c context.Context, mid int64, category int, date time.Time) (infos []*model.TagInfo, err error) {
	infos = make([]*model.TagInfo, 0)
	rows, err := d.db.Query(c, _midAV, mid, date, category)
	if err != nil {
		log.Error("growup-job dao.MIDToAV error(%v),sql(%s)", err, _midAV)
		return
	}
	defer rows.Close()
	for rows.Next() {
		info := &model.TagInfo{}
		err = rows.Scan(&info.AVID, &info.Income, &info.TotalIncome, &info.IsDeleted)
		if err != nil {
			log.Error("growup-job dao.MIDToAV rows scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}
