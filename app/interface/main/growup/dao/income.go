package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	// select
	_avIncomeByMIDSQL      = "SELECT av_id, income, total_income, date FROM av_income WHERE mid = ? AND date >= ? AND date <= ?"
	_avIncomeByAvIDSQL     = "SELECT income, date FROM av_income WHERE av_id = ? AND date <= ?"
	_blacklistByAvIDSQL    = "SELECT av_id FROM av_black_list WHERE av_id in (%s) AND ctype = ? AND is_delete = 0"
	_activityInfoByAvIDSQL = "SELECT archive_id, tag_id FROM activity_info WHERE archive_id in (%s)"
	_tagInfoByTagIDSQL     = "SELECT id, ratio, icon FROM tag_info WHERE id in (%s) and is_deleted = 0"
	_upIncomeTableSQL      = "SELECT mid,income,av_income,column_income,bgm_income,total_income,base_income,av_base_income,column_base_income,bgm_base_income,date FROM %s WHERE mid = ? AND date >= ? AND date <= ? AND is_deleted = 0"
	_upAccountSQL          = "SELECT mid, total_income, total_unwithdraw_income, withdraw_date_version, version FROM up_account WHERE mid = ? AND is_deleted = 0"
	_upIncomeSQL           = "SELECT mid,base_income,income,date FROM up_income WHERE mid=? AND date>=? AND date <=? ORDER BY date"
	_firstUpIncomeSQL      = "SELECT date FROM up_income WHERE mid = ? ORDER BY date LIMIT 1"
	_upIncomeCountSQL      = "SELECT count(*) FROM up_income WHERE date = ?"
	_upDailyCharge         = "SELECT inc_charge FROM up_daily_charge WHERE mid=? AND date>=?"
)

// GetUpDailyCharge get up daily charge
func (d *Dao) GetUpDailyCharge(c context.Context, mid int64, begin string) (incs []int, err error) {
	rows, err := d.db.Query(c, _upDailyCharge, mid, begin)
	if err != nil {
		log.Error("GetUpDailyCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var incCharge int
		err = rows.Scan(&incCharge)
		if err != nil {
			log.Error("rows Scan error(%v)", err)
			return
		}
		incs = append(incs, incCharge)
	}
	return
}

// ListAvIncome list av_income by mid
func (d *Dao) ListAvIncome(c context.Context, mid int64, startTime, endTime string) (avs []*model.ArchiveIncome, err error) {
	avs = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _avIncomeByMIDSQL, mid, startTime, endTime)
	if err != nil {
		log.Error("ListAvIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.ArchiveIncome{}
		err = rows.Scan(&av.ArchiveID, &av.Income, &av.TotalIncome, &av.Date)
		if err != nil {
			log.Error("ListAvIncome rows scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}

	err = rows.Err()
	return
}

// ListAvIncomeByID list av_income by av_id
func (d *Dao) ListAvIncomeByID(c context.Context, avID int64, endTime string) (avs []*model.ArchiveIncome, err error) {
	avs = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _avIncomeByAvIDSQL, avID, endTime)
	if err != nil {
		log.Error("ListAvIncomeByID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.ArchiveIncome{}
		err = rows.Scan(&av.Income, &av.Date)
		if err != nil {
			log.Error("ListAvIncomeByID rows scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}
	err = rows.Err()
	return
}

// ListAvBlackList list av_blakc_list by av_id
func (d *Dao) ListAvBlackList(c context.Context, avIds []int64, typ int) (avb map[int64]struct{}, err error) {
	avb = make(map[int64]struct{})
	rows, err := d.db.Query(c, fmt.Sprintf(_blacklistByAvIDSQL, xstr.JoinInts(avIds)), typ)
	if err != nil {
		log.Error("ListBlackList d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avID int64
		err = rows.Scan(&avID)
		if err != nil {
			log.Error("ListBlackList rows scan error(%v)", err)
			return
		}
		avb[avID] = struct{}{}
	}
	err = rows.Err()
	return
}

// ListActiveInfo list active_info by avid
func (d *Dao) ListActiveInfo(c context.Context, avIds []int64) (acM map[int64]int64, err error) {
	acM = make(map[int64]int64)
	rows, err := d.db.Query(c, fmt.Sprintf(_activityInfoByAvIDSQL, xstr.JoinInts(avIds)))
	if err != nil {
		log.Error("ListActiveInfo d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avID, tagID int64
		err = rows.Scan(&avID, &tagID)
		if err != nil {
			log.Error("ListActiveInfo rows scan error(%v)", err)
			return
		}
		acM[avID] = tagID
	}
	err = rows.Err()
	return
}

// ListTagInfo list tag_info by avid
func (d *Dao) ListTagInfo(c context.Context, tagIds []int64) (tagM map[int64]*model.TagInfo, err error) {
	tagM = make(map[int64]*model.TagInfo)
	rows, err := d.db.Query(c, fmt.Sprintf(_tagInfoByTagIDSQL, xstr.JoinInts(tagIds)))
	if err != nil {
		log.Error("ListTagInfo d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tagInfo := model.TagInfo{}
		err = rows.Scan(&tagInfo.ID, &tagInfo.Radio, &tagInfo.Icon)
		if err != nil {
			log.Error("ListTagInfo rows scan error(%v)", err)
			return
		}
		if val, ok := tagM[tagInfo.ID]; !ok {
			tagM[tagInfo.ID] = &tagInfo
		} else {
			if val.Radio < tagInfo.Radio {
				tagM[tagInfo.ID] = &tagInfo
			}
		}
	}
	err = rows.Err()
	return
}

// ListUpIncome list up_income by mid
func (d *Dao) ListUpIncome(c context.Context, mid int64, table, startTime, endTime string) (ups []*model.UpIncome, err error) {
	ups = make([]*model.UpIncome, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeTableSQL, table), mid, startTime, endTime)
	if err != nil {
		log.Error("ListUpIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpIncome{}
		err = rows.Scan(&up.MID, &up.Income, &up.AvIncome, &up.ColumnIncome, &up.BgmIncome, &up.TotalIncome, &up.BaseIncome, &up.AvBaseIncome, &up.ColumnBaseIncome, &up.BgmBaseIncome, &up.Date)
		if err != nil {
			log.Error("ListUpIncome rows scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}

	err = rows.Err()
	return
}

// ListUpAccount  list up_account by mid
func (d *Dao) ListUpAccount(c context.Context, mid int64) (up *model.UpAccount, err error) {
	up = &model.UpAccount{}
	row := d.db.QueryRow(c, _upAccountSQL, mid)
	if err = row.Scan(&up.MID, &up.TotalIncome, &up.TotalUnwithdrawIncome, &up.WithdrawDateVersion, &up.Version); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// GetUpIncome list up income by date
func (d *Dao) GetUpIncome(c context.Context, mid int64, begin string, end string) (result []*model.UpIncomeStat, err error) {
	rows, err := d.db.Query(c, _upIncomeSQL, mid, begin, end)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		stat := &model.UpIncomeStat{}
		err = rows.Scan(&stat.MID, &stat.BaseIncome, &stat.Income, &stat.Date)
		if err != nil {
			return
		}
		stat.ExtraIncome = stat.Income - stat.BaseIncome
		result = append(result, stat)
	}
	return
}

// GetFirstUpIncome get first up income
func (d *Dao) GetFirstUpIncome(c context.Context, mid int64) (date time.Time, err error) {
	err = d.db.QueryRow(c, _firstUpIncomeSQL, mid).Scan(&date)
	if err == sql.ErrNoRows {
		err = nil
		date = time.Time(0)
	}
	return
}

// GetUpIncomeCount get up income count by date
func (d *Dao) GetUpIncomeCount(c context.Context, date string) (count int, err error) {
	err = d.db.QueryRow(c, _upIncomeCountSQL, date).Scan(&count)
	return
}
