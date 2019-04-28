package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getWhiteList    = "select prize_id, item_id from capsule_prize_extra WHERE prize_id in (%v) and type = ? and status = 1 order by id"
	_updateWhiteList = "UPDATE capsule_prize_extra SET status = ? WHERE prize_id = ?"
	_addWhiteUser    = "INSERT INTO capsule_prize_extra (prize_id,type,item_id,status) VALUES "
	_whiteType       = 1
)

// GetWhiteUserMap 获取白名单用户
func (d *Dao) GetWhiteUserMap(ctx context.Context, prizeIds []int64) (userMap map[int64][]int64, err error) {
	var rows *sql.Rows
	stringPrizeIds := make([]string, len(prizeIds))
	for ix, prizeId := range prizeIds {
		stringPrizeIds[ix] = strconv.FormatInt(prizeId, 10)
	}
	sqlStr := fmt.Sprintf(_getWhiteList, strings.Join(stringPrizeIds, ","))
	if rows, err = d.db.Query(ctx, sqlStr, _whiteType); err != nil {
		log.Error("[dao.pool_prize | GetPoolPrize]query(%s) error(%v)", _getWhiteList, err)
		return
	}
	defer rows.Close()

	userMap = make(map[int64][]int64)
	for rows.Next() {
		var uid, prizeId int64
		if err = rows.Scan(&prizeId, &uid); err != nil {
			log.Error("[dao.pool_prize | GetWhiteUser] scan error, err %v", err)
			return
		}
		if _, ok := userMap[prizeId]; !ok {
			userMap[prizeId] = make([]int64, 0)
		}
		userMap[prizeId] = append(userMap[prizeId], uid)
	}
	return
}

// AddWhiteUser 添加白名单用户
func (d *Dao) AddWhiteUser(ctx context.Context, prizeId int64, uids []int64) (status bool, err error) {
	_, err = d.db.Exec(ctx, _updateWhiteList, 0, prizeId)
	if err != nil {
		log.Error("[dao.pool_prize | AddWhiteUids] update(%s) error(%v)", _updateWhiteList, err)
		return
	}
	rowFields := make([][4]interface{}, 0)
	for _, uid := range uids {
		rowFields = append(rowFields, [4]interface{}{prizeId, _whiteType, uid, 1})
	}
	sqlStr := _addWhiteUser + strings.Repeat("(?,?,?,?),", len(rowFields))
	slen := len(sqlStr) - 1
	sqlStr = sqlStr[0:slen]
	values := make([]interface{}, len(rowFields)*4)
	for i, field := range rowFields {
		for j, v := range field {
			ix := i*4 + j
			values[ix] = v
		}
	}
	res, err := d.db.Exec(ctx, sqlStr, values...)
	if err != nil {
		log.Error("[dao.prize_extra | AddWhiteUser] add(%s) error (%v)", sqlStr, err)
		return
	}
	rows, _ := res.RowsAffected()
	status = rows > 0
	return
}
