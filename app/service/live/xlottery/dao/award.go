package dao

import (
	"context"
	"time"

	"go-common/library/log"
)

const (
	_addAwardRecordMysql = "insert into ap_award_record(uid, gift_name, gift_type, gift_num, source, source_id, create_time, expire_time, status, user_extra_field) values(?,?,?,?,?,?,?,?,?,?)"
)

// AddAward .
func (d *Dao) AddAward(ctx context.Context, uid int64, expireTime string, giftType int64, giftName string, giftNum int64, source string, sourceId int64) bool {
	sql := _addAwardRecordMysql
	affect, err := d.execSqlWithBindParams(ctx, &sql, uid, giftName, giftType, giftNum, source, sourceId, time.Now().Format("2006-01-02 15:04:05"), expireTime, 0, "")
	if err != nil {
		log.Error("[dao.mysql_lottery | AddAward] uid(%d) gift_name(%s) gift_type(%s) source(%s) error(%v)", uid, giftName, giftType, source, err)
	}
	return affect > 0
}
