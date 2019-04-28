package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/search/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_getOffsetSQL = "SELECT offset_incr_id,offset_incr_time,review_incr_id,review_icnr_time FROM digger_offset WHERE project=? AND table_name=?"
	//TODO 有问题，字段少了
	_updateOffsetSQL = "UPDATE digger_offset SET offset_incr_id=?,offset_incr_time=?,mtime=? WHERE project=? AND table_name=?"
	//_initOffsetSQL   = "INSERT INTO digger_offset(project,table_name,offset_incr_time,offset_recover_id,offset_recover_time) VALUES(?,?,?,?,?,?) " +
	//	"ON DUPLICATE KEY UPDATE offset_recover_id=?, offset_recover_time=?"
)

// Offset get offset
func (d *Dao) Offset(c context.Context, appid, tableName string) (res *model.Offset, err error) {
	res = new(model.Offset)
	row := d.SearchDB.QueryRow(c, _getOffsetSQL, appid, tableName)
	if err = row.Scan(&res.OffID, &res.OffTime, &res.ReviewID, &res.ReviewTime); err != nil {
		log.Error("OffsetID row.Scan error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
			res.OffID = 1
			res.OffTime = xtime.Time(time.Now().Unix())
			return
		}
		log.Error("offset row.Scan error(%v)", err)
	}
	return
}

// updateOffset update offset
func (d *Dao) updateOffset(c context.Context, offset *model.LoopOffset, appid, tableName string) (err error) {
	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	if _, err = d.SearchDB.Exec(c, _updateOffsetSQL, offset.OffsetID, offset.OffsetTime, nowFormat, appid, tableName); err != nil {
		log.Error("updateOffset Exec() error(%v)", err)
	}
	return
}

// bulkInitOffset .
func (d *Dao) bulkInitOffset(c context.Context, offset *model.LoopOffset, attrs *model.Attrs, arr []string) (err error) {
	var (
		values          = []string{}
		nowFormat       = time.Now().Format("2006-01-02 15:04:05")
		insertOffsetSQL = "INSERT INTO digger_offset(project,table_name,table_suffix,offset_incr_time,offset_recover_id,offset_recover_time) VALUES"
	)
	if len(arr) == 0 {
		for i := attrs.Table.TableFrom; i <= attrs.Table.TableTo; i++ {
			if attrs.Table.TableTo == 0 {
				arr = append(arr, attrs.Table.TablePrefix)
			} else {
				arr = append(arr, fmt.Sprintf("%s%0"+attrs.Table.TableZero+"d", attrs.Table.TablePrefix, i))
			}
		}
	}
	for _, v := range arr {
		// TODO why???
		// table := attrs.Table.TablePrefix + v
		// value := "('" + attrs.AppID + "','" + table + "','" + v + "','" + nowFormat + "'," + strconv.FormatInt(offset.RecoverID, 10) + ",'" + offset.RecoverTime + "')"
		value := "('" + attrs.AppID + "','" + v + "','" + attrs.Table.TablePrefix + "','" + nowFormat + "'," + strconv.FormatInt(offset.RecoverID, 10) + ",'" + offset.RecoverTime + "')"
		values = append(values, value)
	}
	valueStr := strings.Join(values, ",") + " ON DUPLICATE KEY UPDATE offset_recover_id=VALUES(offset_recover_id),offset_recover_time=VALUES(offset_recover_time)"
	bulkInserSQL := insertOffsetSQL + valueStr
	_, err = d.SearchDB.Exec(c, bulkInserSQL)
	return
}
