package charge

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
)

const (
	_layout = "2006-01-02"

	_avDailyChargeSQL   = "SELECT id,av_id,mid,tag_id,is_original,upload_time,danmaku_count,comment_count,collect_count,coin_count,share_count,elec_pay_count,total_play_count,web_play_count,app_play_count,h5_play_count,lv_unknown,lv_0,lv_1,lv_2,lv_3,lv_4,lv_5,lv_6,v_score,inc_charge,total_charge,date,is_deleted,ctime,mtime FROM av_daily_charge_%s WHERE id > ? AND date = ? AND inc_charge > 0 ORDER BY id LIMIT ?"
	_avWeeklyChargeSQL  = "SELECT id,av_id,mid,tag_id,is_original,upload_time,danmaku_count,comment_count,collect_count,coin_count,share_count,elec_pay_count,total_play_count,web_play_count,app_play_count,h5_play_count,lv_unknown,lv_0,lv_1,lv_2,lv_3,lv_4,lv_5,lv_6,v_score,inc_charge,total_charge,date,is_deleted,ctime,mtime FROM av_weekly_charge WHERE id > ? AND date = ? ORDER BY id LIMIT ?"
	_avMonthlyChargeSQL = "SELECT id,av_id,mid,tag_id,is_original,upload_time,danmaku_count,comment_count,collect_count,coin_count,share_count,elec_pay_count,total_play_count,web_play_count,app_play_count,h5_play_count,lv_unknown,lv_0,lv_1,lv_2,lv_3,lv_4,lv_5,lv_6,v_score,inc_charge,total_charge,date,is_deleted,ctime,mtime FROM av_monthly_charge WHERE id > ? AND date = ? ORDER BY id LIMIT ?"

	_inAvChargeTableSQL = "INSERT INTO %s(av_id,mid,tag_id,is_original,danmaku_count,comment_count,collect_count,coin_count,share_count,elec_pay_count,total_play_count,web_play_count,app_play_count,h5_play_count,lv_unknown,lv_0,lv_1,lv_2,lv_3,lv_4,lv_5,lv_6,v_score,inc_charge,total_charge,date,upload_time) VALUES %s ON DUPLICATE KEY UPDATE danmaku_count=VALUES(danmaku_count),comment_count=VALUES(comment_count),collect_count=VALUES(collect_count),coin_count=VALUES(coin_count),share_count=VALUES(share_count),elec_pay_count=VALUES(elec_pay_count),total_play_count=VALUES(total_play_count),web_play_count=VALUES(web_play_count),app_play_count=VALUES(app_play_count),h5_play_count=VALUES(h5_play_count),lv_unknown=VALUES(lv_unknown),lv_0=VALUES(lv_0),lv_1=VALUES(lv_1),lv_2=VALUES(lv_2),lv_3=VALUES(lv_3),lv_4=VALUES(lv_4),lv_5=VALUES(lv_5),lv_6=VALUES(lv_6),v_score=VALUES(v_score),inc_charge=VALUES(inc_charge),total_charge=VALUES(total_charge)"
)

// IAvCharge av_charge_weekly monthly
type IAvCharge func(c context.Context, date time.Time, id int64, limit int) (data []*model.AvCharge, err error)

// AvDailyCharge av_daily_charge
func (d *Dao) AvDailyCharge(c context.Context, date time.Time, id int64, limit int) (data []*model.AvCharge, err error) {
	month := date.Month()
	monthStr := strconv.Itoa(int(month))
	if month < 10 {
		monthStr = "0" + monthStr
	}
	return d.GetAvCharge(c, fmt.Sprintf(_avDailyChargeSQL, monthStr), date.Format(_layout), id, limit)
}

// AvWeeklyCharge av_weekly_charge
func (d *Dao) AvWeeklyCharge(c context.Context, date time.Time, id int64, limit int) (data []*model.AvCharge, err error) {
	return d.GetAvCharge(c, _avWeeklyChargeSQL, date.Format(_layout), id, limit)
}

// AvMonthlyCharge av_monthly_charge
func (d *Dao) AvMonthlyCharge(c context.Context, date time.Time, id int64, limit int) (data []*model.AvCharge, err error) {
	return d.GetAvCharge(c, _avMonthlyChargeSQL, date.Format(_layout), id, limit)
}

// GetAvCharge get av_charge
func (d *Dao) GetAvCharge(c context.Context, sql string, date string, id int64, limit int) (data []*model.AvCharge, err error) {
	rows, err := d.db.Query(c, sql, id, date, limit)
	if err != nil {
		log.Error("d.Query av_charge(%s) error(%v)", sql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		adc := &model.AvCharge{}
		err = rows.Scan(&adc.ID, &adc.AvID, &adc.MID, &adc.TagID, &adc.IsOriginal, &adc.UploadTime, &adc.DanmakuCount, &adc.CommentCount, &adc.CollectCount, &adc.CoinCount, &adc.ShareCount, &adc.ElecPayCount, &adc.TotalPlayCount, &adc.WebPlayCount, &adc.AppPlayCount, &adc.H5PlayCount, &adc.LvUnknown, &adc.Lv0, &adc.Lv1, &adc.Lv2, &adc.Lv3, &adc.Lv4, &adc.Lv5, &adc.Lv6, &adc.VScore, &adc.IncCharge, &adc.TotalCharge, &adc.Date, &adc.IsDeleted, &adc.CTime, &adc.MTime)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		data = append(data, adc)
	}
	return
}

// InsertAvChargeTable add av charge batch
func (d *Dao) InsertAvChargeTable(c context.Context, vals, table string) (rows int64, err error) {
	if vals == "" {
		return
	}

	res, err := d.db.Exec(c, fmt.Sprintf(_inAvChargeTableSQL, table, vals))
	if err != nil {
		log.Error("InsertAvChargeTable(%s) tx.Exec error(%v)", table, err)
		return
	}
	return res.RowsAffected()
}
