package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/live/xroom-feed/internal/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_sqlConf      = "select id, name, type, rules, module_type, position, percent, priority from ap_rec_pool_conf where is_del = 0 and (status = 2 or (status = 1 and start_time < '%s' and end_time > '%s'))"
	_sqlWhiteList = "select room_id from ap_rec_white_list where rec_id = %d and is_del = 0"
)

func (d *Dao) GetConfFromDb() (ruleConf []*model.RecPoolConf, err error) {
	var rows *sql.Rows
	ruleConf = make([]*model.RecPoolConf, 0)
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	rows, err = d.liveAppDb.Query(context.TODO(), fmt.Sprintf(_sqlConf, timeNow, timeNow))
	if err != nil {
		log.Error("getConfFromDb_error:%+v", err)
		return
	}
	for rows.Next() {
		r := new(model.RecPoolConf)
		if err = rows.Scan(&r.Id, &r.Name, &r.ConfType, &r.Rules, &r.ModuleType, &r.Position, &r.Percent, &r.Priority); err != nil {
			log.Error("getConfFromDb_parseError:%+v", err)
			continue
		}

		ruleConf = append(ruleConf, r)
	}

	return
}

func (d *Dao) GetWhiteList(ctx context.Context, id int) (whiteStr string, err error) {
	//白名单获取
	var whiteRows *sql.Rows
	whiteRows, err = d.liveAppDb.Query(ctx, fmt.Sprintf(_sqlWhiteList, id))
	if err != nil {
		log.Error("getWhiteListFromDb_error:%+v", err)
		return
	}
	whiteList := make([]string, 0)
	for whiteRows.Next() {
		r := new(model.RecWhiteList)
		if err = whiteRows.Scan(&r.RoomId); err != nil {
			log.Error("getWhiteListFromDb_parseError:%+v", err)
			continue
		}
		if r.RoomId == 0 {
			continue
		}
		whiteList = append(whiteList, strconv.Itoa(r.RoomId))

	}

	if len(whiteList) <= 0 {
		return
	}
	whiteStr = strings.Join(whiteList, ",")
	return
}
