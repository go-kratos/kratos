package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	v1pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getCoinConfig    = "SELECT coin_id, type, area_v2_parent_id, area_v2_id, gift_id FROM capsule_coin_config WHERE coin_id = ? AND status = 1 "
	_getCoinConfigMap = "SELECT coin_id, type, area_v2_parent_id, area_v2_id, gift_id FROM capsule_coin_config WHERE coin_id in (%v) AND status = 1"
	_updateCoinConfig = "UPDATE capsule_coin_config SET status = ? WHERE coin_id = ?"
	_addCoinConfig    = "INSERT INTO capsule_coin_config (type,gift_id,area_v2_parent_id,area_v2_id,coin_id,status) VALUES "
)

//GetCoinConfig 获取扭蛋币配置
func (d *Dao) GetCoinConfig(ctx context.Context, coinId int64) (configs []*model.CoinConfig, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getCoinConfig, coinId); err != nil {
		log.Error("[dao.coin_config | GetCoinConfig]query(%s) error(%v)", _getCoinConfig, err)
		return
	}
	defer rows.Close()

	configs = make([]*model.CoinConfig, 0)
	for rows.Next() {
		d := &model.CoinConfig{}
		if err = rows.Scan(&d.CoinId, &d.Type, &d.AreaV2ParentId, &d.AreaV2Id, &d.GiftId); err != nil {
			log.Error("[dao.coin_config | GetCoinConfig] scan error, err %v", err)
			return
		}
		configs = append(configs, d)
	}
	return
}

//GetCoinConfigMap 批量获取扭蛋币
func (d *Dao) GetCoinConfigMap(ctx context.Context, coinIds []int64) (configMap map[int64][]*model.CoinConfig, err error) {
	var rows *sql.Rows
	stringCoinIds := make([]string, 0)
	for _, coinId := range coinIds {
		stringCoinIds = append(stringCoinIds, strconv.FormatInt(coinId, 10))
	}
	coinString := strings.Join(stringCoinIds, ",")
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getCoinConfigMap, coinString)); err != nil {
		log.Error("[dao.coin_config | GetCoinConfigMap]query(%s) error(%v)", _getCoinConfigMap, err)
		return
	}
	defer rows.Close()

	configMap = make(map[int64][]*model.CoinConfig)
	for rows.Next() {
		d := &model.CoinConfig{}
		if err = rows.Scan(&d.CoinId, &d.Type, &d.AreaV2ParentId, &d.AreaV2Id, &d.GiftId); err != nil {
			log.Error("[dao.coin_config | GetCoinConfigMap] scan error, err %v", err)
			return
		}
		if _, ok := configMap[d.CoinId]; !ok {
			configMap[d.CoinId] = make([]*model.CoinConfig, 0)
		}
		configMap[d.CoinId] = append(configMap[d.CoinId], d)
	}
	return
}

//UpdateCoinConfig 更新扭蛋币配置
func (d *Dao) UpdateCoinConfig(ctx context.Context, coinId int64, areaIds []*v1pb.UpdateCoinConfigReq_AreaIds, giftIds []int64) (status bool, err error) {
	_, err = d.db.Exec(ctx, _updateCoinConfig, 2, coinId)
	if err != nil {
		log.Error("[dao.coin_config | UpdateCoinConfig] query(%s) error(%v)", _updateCoinConfig, err)
		return
	}
	rowFields := make([][6]interface{}, 0)
	for _, giftId := range giftIds {
		rowFields = append(rowFields, [6]interface{}{2, giftId, 0, 0, coinId, 1})
	}
	if len(areaIds) != 0 {
		for _, areaId := range areaIds {
			if areaId.IsAll == 1 {
				rowFields = append(rowFields, [6]interface{}{1, 0, areaId.ParentId, areaId.ParentId, coinId, 1})
			} else {
				for _, areaV2Id := range areaId.List {
					rowFields = append(rowFields, [6]interface{}{1, 0, areaId.ParentId, areaV2Id, coinId, 1})
				}

			}
		}
	}
	sqlStr := _addCoinConfig + strings.Repeat("(?,?,?,?,?,?),", len(rowFields))
	slen := len(sqlStr) - 1
	sqlStr = sqlStr[0:slen]
	values := make([]interface{}, len(rowFields)*6)
	for i, field := range rowFields {
		for j, v := range field {
			ix := i*6 + j
			values[ix] = v
		}
	}
	res, err := d.db.Exec(ctx, sqlStr, values...)
	if err != nil {
		log.Error("[dao.coin_config | UpdateCoinConfig] insert(%s) error(%v)", sqlStr, err)
		return
	}
	rows, _ := res.RowsAffected()
	status = rows > 0
	return
}
