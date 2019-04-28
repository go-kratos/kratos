package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	v12 "go-common/app/service/live/rc/api/liverpc/v1"
	"go-common/library/cache/redis"

	"go-common/app/job/live/xlottery/internal/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// NormalCoinId 普通扭蛋id
const NormalCoinId = 1

// ColorfulCoinId 梦幻扭蛋id
const ColorfulCoinId = 2

// WeekCoinId 周星扭蛋id
const WeekCoinId = 3

// LplCoinId 周星扭蛋id
const LplCoinId = 4

// BlessCoinId 祈福抽奖券
const BlessCoinId = 5

// IsBottomPool 是否为保底奖池
const IsBottomPool = 1

// CapsuleActionTrans 转换扭蛋币
const CapsuleActionTrans = "trans"

const (
	// NormalCoinString 普通扭蛋字符串标识，数据库和redis
	NormalCoinString = "normal"

	// ColorfulCoinString 梦幻扭蛋字符串标识，数据库和redis
	ColorfulCoinString = "colorful"

	// WeekCoinString 周星扭蛋字符串标识，数据库和redis
	WeekCoinString = "week"

	// LplCoinString 周星扭蛋字符串标识，数据库和redis
	LplCoinString = "lpl"

	// BlessCoinString 祈福扭蛋字符串标识，数据库和redis
	BlessCoinString = "bless"
)
const (
	_capsuleNotice        = "capsule:notice:%s:%d"
	_userCapsuleCoinRedis = "hash:capsule:user:%d"
	_userInfoRedis        = "capsule:user:info:%d:%d"
	_capsuleConfRand      = "capsule:rand"
	_lplSendGiftRedis     = "capsule:lpl:send:gift:%s:%d"
)
const (
	_ = iota
	// CapsulePrizeGift1Type 辣条
	CapsulePrizeGift1Type // 辣条
	// CapsulePrizeTitleType 头衔
	CapsulePrizeTitleType
	// CapsulePrizeStuff1Type 经验原石
	CapsulePrizeStuff1Type
	// CapsulePrizeStuff2Type 经验曜石
	CapsulePrizeStuff2Type
	// CapsulePrizeStuff3Type 贤者之石
	CapsulePrizeStuff3Type
	// CapsulePrizeSmallTvType 小电视
	CapsulePrizeSmallTvType
	// CapsulePrizeGuard3Type 舰长体验
	CapsulePrizeGuard3Type
	// CapsulePrizeGuard2Type 提督体验
	CapsulePrizeGuard2Type
	// CapsulePrizeGuard1Type 总督体验
	CapsulePrizeGuard1Type
)

const (
	_ = iota
	// CapsuleGiftTypeAll gift_type 为全部道具
	CapsuleGiftTypeAll
	// CapsuleGiftTypeGold gift_type 为金瓜子道具
	CapsuleGiftTypeGold
	// CapsuleGiftTypeSelected gift_type 为指定道具
	CapsuleGiftTypeSelected
)
const (
	_getActiveColorPool       = "SELECT id, coin_id, title, rule, start_time, end_time, status FROM capsule_pool WHERE status = 1 AND coin_id = ?"
	_getCapsuleMaxId          = "SELECT id FROM capsule_%s order by id desc limit 1"
	_transCapsule             = "UPDATE capsule_%s SET normal_score = normal_score + ?,colorful_score = 0 WHERE uid = ?"
	_getChangeNum             = "SELECT change_num FROM capsule_coin WHERE id = ?"
	_getUserInfoById          = "SELECT uid,normal_score,colorful_score FROM capsule_%s WHERE id >= ? AND id < ? and colorful_score > 0"
	_getUserInfoByUids        = "SELECT uid,normal_score,colorful_score FROM capsule_%s WHERE uid in(%s)"
	_getOnCoin                = "SELECT id, title, gift_type, change_num, start_time, end_time, status FROM capsule_coin WHERE status = 1"
	_getCoinConfigMap         = "SELECT coin_id, type, area_v2_parent_id, area_v2_id, gift_id FROM capsule_coin_config WHERE coin_id in (%v) AND status = 1"
	_getPoolMap               = "SELECT id, coin_id, title, rule, start_time, end_time, status, is_bottom FROM capsule_pool WHERE status = 1 and coin_id in (%v)"
	_getPoolPrizeMap          = "SELECT id, pool_id, type, num, object_id,expire, web_url, mobile_url, description, jump_url, pro_type, chance, loop_num, limit_num, weight FROM capsule_pool_prize WHERE pool_id in (%s) and status = 1 order by ctime"
	_reportCapsuleChangeMysql = "insert into capsule_log_%s(uid, type, score, action, platform, pre_normal_score, pre_colorful_score, cur_normal_score, cur_colorful_score) values(?,?,?,?,?,?,?,?,?)"
	_userCapsuleCoinMysql     = "select normal_score, colorful_score from capsule_%d where uid = ?"
	_addCapsuleCoinMysql      = "insert into capsule_%d(uid, normal_score, colorful_score) values(?, ?, ?)"
	_userInfoMysql            = "select score from capsule_info_%d where uid = ? and type = ?"
	_addInfoMysql             = "insert into capsule_info_%d(uid,type,score) values(?, ?, ?)"
	_transUserCapsuleMysql    = "update capsule_%d set colorful_score = 0, normal_score = normal_score + ? where uid = ?"
	_updateUserCapsuleMysql   = "update capsule_%d set %s_score = %s_score + ? where uid = ?"
	_updateUserInfoMysql      = "update capsule_info_%d set score = score + ? where uid = ? and type = ?"
	_getExtraDataMysql        = "select item_value, item_extra from capsule_extra_data where uid = ? and type = ?"
	_addExtraDataMysql        = "insert into capsule_extra_data(uid,type,item_value,item_extra) values(?,?,?,?)"
	_getExtraDataByTimeMysql  = "select id, uid, type, item_value, item_extra from capsule_extra_data where mtime >= ? and mtime < ?"
	_updateExtraValueMysql    = "update capsule_extra_data set item_value = ? where id = ?"
	_updateExtraMtimeMysql    = "update capsule_extra_data set mtime = ? where id = ?"
	_updateExtraMysql         = "update capsule_extra_data set item_value = ?, item_extra = ? where id = ?"
	_getExtraDataMaxIdMysql   = "select id from capsule_extra_data order by id desc limit 1"
	_getExtraDataByIdMysql    = "select id, uid, type, item_value, item_extra from capsule_extra_data where id >= ? and id < ?"
)

var (
	capsuleConf CapsuleConf

	// ReprotConfig map
	ReprotConfig map[int64]string

	// CoinIdIntMap map
	CoinIdIntMap map[int64]string
	// PrizeNameMap map
	PrizeNameMap map[int64]string
)
var (
	// UnLockGetWrong flag
	UnLockGetWrong = "UnLockGetWrong"
	// ErrUnLockGet error
	ErrUnLockGet = errors.New(UnLockGetWrong)
)

// CapsuleConf 扭蛋全局配置
type CapsuleConf struct {
	CoinConfMap map[int64]*CapsuleCoinConf
	CacheTime   int64
	ChangeFlag  int64
	RwLock      sync.RWMutex
}

// CapsuleCoinConf 扭蛋币配置
type CapsuleCoinConf struct {
	Id          int64
	Title       string
	GiftType    int64
	ChangeNum   int64
	StartTime   int64
	EndTime     int64
	Status      int64
	GiftMap     map[int64]struct{}
	AreaMap     map[int64]struct{}
	PoolConf    *CapsulePoolConf
	AllPoolConf []*CapsulePoolConf
}

// CapsulePoolConf  奖池配置
type CapsulePoolConf struct {
	Id                 int64
	CoinId             int64
	Title              string
	Rule               string
	StartTime, EndTime int64
	Status             int64
	IsBottom           int64
	PoolPrize          []*CapsulePoolPrize
}

// CapsulePoolPrize 奖池奖品
type CapsulePoolPrize struct {
	Id, PoolId, Type, Num, ObjectId, Expire           int64
	Name, WebImage, MobileImage, Description, JumpUrl string
	ProType                                           int64
	Chance                                            int64
	LoopNum, LimitNum, Weight                         int64
}

func init() {
	CoinIdIntMap = make(map[int64]string)
	CoinIdIntMap[NormalCoinId] = NormalCoinString
	CoinIdIntMap[ColorfulCoinId] = ColorfulCoinString
	CoinIdIntMap[WeekCoinId] = WeekCoinString
	CoinIdIntMap[LplCoinId] = LplCoinString
	CoinIdIntMap[BlessCoinId] = BlessCoinString
	PrizeNameMap = make(map[int64]string)
	PrizeNameMap[CapsulePrizeGift1Type] = "辣条"
	PrizeNameMap[CapsulePrizeTitleType] = "头衔"
	PrizeNameMap[CapsulePrizeStuff1Type] = "经验原石"
	PrizeNameMap[CapsulePrizeStuff2Type] = "经验曜石"
	PrizeNameMap[CapsulePrizeStuff3Type] = "贤者之石"
	PrizeNameMap[CapsulePrizeSmallTvType] = "小电视抱枕"
	PrizeNameMap[CapsulePrizeGuard3Type] = "舰长体验"
	PrizeNameMap[CapsulePrizeGuard2Type] = "提督体验"
	ReprotConfig = make(map[int64]string)
	ReprotConfig[0] = "未知"
	ReprotConfig[1] = "增加普通扭蛋"
	ReprotConfig[2] = "增加梦幻扭蛋"
	ReprotConfig[3] = "减少普通扭蛋"
	ReprotConfig[4] = "减少梦幻扭蛋"
	ReprotConfig[5] = "梦幻转化普通"
}

// GetActiveColorPool 获取奖池
func (d *Dao) GetActiveColorPool(ctx context.Context) (pool []*model.Pool, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getActiveColorPool, ColorfulCoinId); err != nil {
		log.Error("[dao.capsule | GetActiveColorPool]query(%s) error(%v)", _getActiveColorPool, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Pool{}
		if err = rows.Scan(&p.Id, &p.CoinId, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.Status); err != nil {
			log.Error("[dao.capsule | GetActiveColorPool] scan error, err %v", err)
			return
		}
		pool = append(pool, p)
	}
	return
}

// GetUserInfoById 获取扭蛋币信息
func (d *Dao) GetUserInfoById(ctx context.Context, table string, start int64, end int64) (infos [][4]int64, err error) {
	var rows *sql.Rows
	sqlStr := fmt.Sprintf(_getUserInfoById, table)
	if rows, err = d.db.Query(ctx, sqlStr, start, end); err != nil {
		log.Error("[dao.capsule | GetUserInfos]query(%s) error(%v)", _getUserInfoById, err)
		return
	}
	defer rows.Close()
	infos = make([][4]int64, 0)
	for rows.Next() {
		p := &model.UserInfo{}
		if err = rows.Scan(&p.Uid, &p.NormalScore, &p.ColorfulScore); err != nil {
			log.Error("[dao.capsule | GetUserInfos] scan error, err %v", err)
			return
		}
		infos = append(infos, [4]int64{p.Uid, p.NormalScore, p.ColorfulScore, 0})
	}
	return
}

// GetUserInfoByUids 获取扭蛋币信息
func (d *Dao) GetUserInfoByUids(ctx context.Context, table string, uids []int64) (infos map[int64][2]int64, err error) {
	var rows *sql.Rows
	uidStrArray := make([]string, len(uids))
	for ix, uid := range uids {
		uidStrArray[ix] = strconv.FormatInt(uid, 10)
	}
	uidStr := strings.Join(uidStrArray, ",")
	sqlStr := fmt.Sprintf(_getUserInfoByUids, table, uidStr)
	if rows, err = d.db.Query(ctx, sqlStr); err != nil {
		log.Error("[dao.capsule | GetUserInfoByUids]query(%s) error(%v)", _getUserInfoByUids, err)
		return
	}
	defer rows.Close()

	infos = make(map[int64][2]int64)
	for rows.Next() {
		p := &model.UserInfo{}
		if err = rows.Scan(&p.Uid, &p.NormalScore, &p.ColorfulScore); err != nil {
			log.Error("[dao.capsule | GetUserInfoByUids] scan error, err %v", err)
			return
		}
		infos[p.Uid] = [2]int64{p.NormalScore, p.ColorfulScore}
	}
	return
}

// GetTransNum 获取扭蛋币数量
func (d *Dao) GetTransNum(ctx context.Context, coinId int64) (changeNum int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getChangeNum, coinId); err != nil {
		log.Error("[dao.capsule | GetTransNum]query(%s) error(%v)", _getChangeNum, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Coin{}
		if err = rows.Scan(&p.ChangeNum); err != nil {
			log.Error("[dao.capsule | GetTransNum]scan error, err %v", err)
			return
		}
		changeNum = p.ChangeNum
	}
	return
}

// TransCapsule 转换扭蛋币
func (d *Dao) TransCapsule(ctx context.Context, table string, colorChangeNum int64, normalChangeNum int64) (err error) {
	var maxId int64
	row := d.db.QueryRow(ctx, fmt.Sprintf(_getCapsuleMaxId, table))
	if err = row.Scan(&maxId); err != nil {
		log.Error("[dao.capsule | TransCapsule] query(%s),err(%v)", _getCapsuleMaxId, err)
		return
	}
	log.Info("[dao.capsule | TransCapsule] table: %s, max: %d", table, maxId)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	var i int64
	for i = 0; i <= maxId; i = i + 1000 {
		var userInfos [][4]int64
		userInfos, err = d.GetUserInfoById(ctx, table, i, i+1000)
		if err != nil {
			log.Error("[dao.capsule | TransCapsule] GetUserInfos error %v", err)
			return err
		}
		ulen := len(userInfos)
		if ulen < 1 {
			continue
		}
		uids := make([]int64, ulen)
		for ix, userInfo := range userInfos {
			uid := userInfo[0]
			changeScore := userInfo[2] - userInfo[2]%normalChangeNum
			_, err = d.db.Exec(ctx, fmt.Sprintf(_transCapsule, table), changeScore, uid)
			if err != nil {
				log.Error("[dao.capsule | TransCapsule]query(%s) error(%v)", _transCapsule, err)
				return
			}
			userInfos[ix][3] = changeScore
			uids[ix] = uid
			uKey := userKey(uid)
			_, e := conn.Do("DEL", uKey)
			if e != nil {
				log.Error("redis_lottery|delete score error,%v，  uid %v", e, uid)
			}

		}
		var userMap map[int64][2]int64
		userMap, err = d.GetUserInfoByUids(ctx, table, uids)
		if err != nil {
			log.Error("[dao.capsule | TransCapsule] GetUserInfoByUids error %v", err)
			return err
		}
		if userMap == nil {
			continue
		}
		for _, userInfo := range userInfos {
			changeScore := userInfo[3]
			uid := userInfo[0]
			if _, ok := userMap[uid]; !ok {
				continue
			}
			report.User(&report.UserInfo{
				Business: 101,
				Type:     int(1),
				Oid:      uid,
				Action:   "capsule_change",
				Ctime:    time.Now(),
				Index: []interface{}{
					"梦幻转化普通",
					changeScore,
					userMap[uid][0],
					userMap[uid][1],
					"trans",
				},
			})
			date := time.Now().Format("200601")
			sqlStr := fmt.Sprintf(_reportCapsuleChangeMysql, date)
			_, err := d.execSqlWithBindParams(ctx, &sqlStr, uid, 1, changeScore, "trans", "", userInfo[1], userInfo[2], userMap[uid][0], userMap[uid][1])
			if err != nil {
				log.Error("[dao.capsule | TransCapsule] AddCapsuleLog error %v", err)
				continue
			}
		}
	}
	return
}

func userKey(uid int64) string {
	return fmt.Sprintf(_userCapsuleCoinRedis, uid)
}

func userInfoKey(uid int64, coinId int64) string {
	return fmt.Sprintf(_userInfoRedis, uid, coinId)
}

// GetCapsuleChangeFlag 获取扭蛋配置变化标记
func (d *Dao) GetCapsuleChangeFlag(ctx context.Context) (randNum int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	randNum, err = redis.Int64(conn.Do("GET", _capsuleConfRand))
	if err != nil {
		log.Error("[dao.redis_lottery|GetCapsuleChangeFlag] conn.GET(%s) error(%v)", _capsuleConfRand, err)
		return
	}
	return
}

// GetCapsuleConf 获取扭蛋币配置
func (d *Dao) GetCapsuleConf(ctx context.Context) (conf map[int64]*CapsuleCoinConf, err error) {
	capsuleConf.RwLock.RLock()
	tmpConf := capsuleConf.CoinConfMap
	capsuleConf.RwLock.RUnlock()
	if len(tmpConf) == 0 {
		redisChangeFlag, _ := d.GetCapsuleChangeFlag(ctx)
		tmpConf, err = d.RelaodCapsuleConfig(ctx, redisChangeFlag)
		if err != nil || tmpConf == nil || len(tmpConf) == 0 {
			log.Error("[dao.capsule | GetCapsuleConf] CapsuleCoinConf is empty")
			return nil, err
		}
	}
	now := time.Now().Unix()
	conf = make(map[int64]*CapsuleCoinConf)
	for coinId := range CoinIdIntMap {
		if _, ok := tmpConf[coinId]; ok {
			coinConf := tmpConf[coinId]
			if coinConf.AllPoolConf != nil && len(coinConf.AllPoolConf) > 0 {
				for _, poolConf := range coinConf.AllPoolConf {
					if poolConf.StartTime < now && poolConf.EndTime > now {
						if _, ok := conf[coinId]; !ok {
							conf[coinId] = &CapsuleCoinConf{
								Id:          coinConf.Id,
								Title:       coinConf.Title,
								GiftType:    coinConf.GiftType,
								ChangeNum:   coinConf.ChangeNum,
								StartTime:   coinConf.StartTime,
								EndTime:     coinConf.EndTime,
								Status:      coinConf.Status,
								GiftMap:     coinConf.GiftMap,
								AreaMap:     coinConf.AreaMap,
								PoolConf:    poolConf,
								AllPoolConf: coinConf.AllPoolConf,
							}
						} else {
							if poolConf.IsBottom != IsBottomPool {
								conf[coinId] = &CapsuleCoinConf{
									Id:          coinConf.Id,
									Title:       coinConf.Title,
									GiftType:    coinConf.GiftType,
									ChangeNum:   coinConf.ChangeNum,
									StartTime:   coinConf.StartTime,
									EndTime:     coinConf.EndTime,
									Status:      coinConf.Status,
									GiftMap:     coinConf.GiftMap,
									AreaMap:     coinConf.AreaMap,
									PoolConf:    poolConf,
									AllPoolConf: coinConf.AllPoolConf,
								}
							}
						}
					}

				}
			}
		}
	}
	return conf, nil
}

//GetCoinMap 批量获取扭蛋币
func (d *Dao) GetCoinMap(ctx context.Context) (coinMap map[int64]*model.Coin, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getOnCoin); err != nil {
		log.Error("[dao.coin | GetCoinMap] query(%s) error(%v)", _getOnCoin, err)
		return
	}
	defer rows.Close()

	coinMap = make(map[int64]*model.Coin)
	for rows.Next() {
		p := &model.Coin{}
		if err = rows.Scan(&p.Id, &p.Title, &p.GiftType, &p.ChangeNum, &p.StartTime, &p.EndTime, &p.Status); err != nil {
			log.Error("[dao.coin | GetCoinMap] scan error, err %v", err)
			return
		}
		coinMap[p.Id] = p
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
		configMap[d.CoinId] = append(configMap[d.CoinId], d)
	}
	return
}

//GetPoolMap 批量奖池信息
func (d *Dao) GetPoolMap(ctx context.Context, coinIds []int64) (poolMap map[int64][]*model.Pool, err error) {
	var rows *sql.Rows
	stringCoinIds := make([]string, 0)
	for _, coinId := range coinIds {
		stringCoinIds = append(stringCoinIds, strconv.FormatInt(coinId, 10))
	}
	coinString := strings.Join(stringCoinIds, ",")
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getPoolMap, coinString)); err != nil {
		log.Error("[dao.pool | GetPoolMap]query(%s) error(%v)", _getPoolMap, err)
		return
	}
	defer rows.Close()

	poolMap = make(map[int64][]*model.Pool)
	for rows.Next() {
		d := &model.Pool{}
		if err = rows.Scan(&d.Id, &d.CoinId, &d.Title, &d.Description, &d.StartTime, &d.EndTime, &d.Status, &d.IsBottom); err != nil {
			log.Error("[dao.pool |GetPoolMap] scan error, err %v", err)
			return
		}
		if _, ok := poolMap[d.CoinId]; !ok {
			poolMap[d.CoinId] = make([]*model.Pool, 0)
		}
		poolMap[d.CoinId] = append(poolMap[d.CoinId], d)
	}
	return
}

// GetPoolPrizeMap 批量奖池奖品
func (d *Dao) GetPoolPrizeMap(ctx context.Context, poolIds []int64) (poolPrizeMap map[int64][]*model.PoolPrize, err error) {
	var rows *sql.Rows
	stringPoolIds := make([]string, 0)
	for _, poolId := range poolIds {
		stringPoolIds = append(stringPoolIds, strconv.FormatInt(poolId, 10))
	}
	poolString := strings.Join(stringPoolIds, ",")
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getPoolPrizeMap, poolString)); err != nil {
		log.Error("[dao.pool_prize | GetPoolPrizeMap] query(%s) error(%v)", _getPoolPrizeMap, err)
		return
	}
	defer rows.Close()

	poolPrizeMap = make(map[int64][]*model.PoolPrize)
	for rows.Next() {
		d := &model.PoolPrize{}
		if err = rows.Scan(&d.Id, &d.PoolId, &d.Type, &d.Num, &d.ObjectId, &d.Expire, &d.WebUrl, &d.MobileUrl, &d.Description, &d.JumpUrl, &d.ProType, &d.Chance, &d.LoopNum, &d.LimitNum, &d.Weight); err != nil {
			log.Error("[dao.pool_prize | GetPoolPrizeMap] scan error, err %v", err)
			return
		}
		if _, ok := PrizeNameMap[d.Type]; !ok {
			continue
		}
		poolPrizeMap[d.PoolId] = append(poolPrizeMap[d.PoolId], d)
	}
	return
}

// RelaodCapsuleConfig 重新加载扭蛋配置
func (d *Dao) RelaodCapsuleConfig(ctx context.Context, changeFlag int64) (conf map[int64]*CapsuleCoinConf, err error) {
	coinMap, err := d.GetCoinMap(ctx)
	if err != nil || len(coinMap) == 0 {
		log.Error("[dao.capsule | RelaodCapsuleConfig] coinMap is empty")
		return
	}
	coinIds := make([]int64, len(coinMap))
	ix := 0
	for _, coinInfo := range coinMap {
		coinIds[ix] = coinInfo.Id
		ix++
	}
	coinConfigMap, err := d.GetCoinConfigMap(ctx, coinIds)
	if err != nil || len(coinConfigMap) == 0 {
		log.Error("[dao.capsule | RelaodCapsuleConfig] CoinConfigMap is empty")
		return
	}
	poolMap, err := d.GetPoolMap(ctx, coinIds)
	if err != nil || len(poolMap) == 0 {
		log.Error("[dao.capsule | RelaodCapsuleConfig] PoolMap is empty")
		return
	}
	poolIds := make([]int64, 0)
	for _, pools := range poolMap {
		for _, pool := range pools {
			poolIds = append(poolIds, pool.Id)
		}

	}
	poolPrizeMap, err := d.GetPoolPrizeMap(ctx, poolIds)
	if err != nil || len(poolPrizeMap) == 0 {
		log.Error("[dao.capsule | RelaodCapsuleConfig] PoolPrizeMap is empty")
		return
	}
	coinConfMap := make(map[int64]*CapsuleCoinConf)
	ids := make([]int64, 0)
	idMap := make(map[int64]struct{})
	for _, prizeList := range poolPrizeMap {
		for _, prize := range prizeList {
			if prize.ObjectId != 0 && prize.Type == CapsulePrizeTitleType {
				if _, ok := idMap[prize.ObjectId]; !ok {
					ids = append(ids, prize.ObjectId)
					idMap[prize.ObjectId] = struct{}{}
				}
			}
		}
	}
	titleMap := make(map[int64]string)
	if len(ids) != 0 {
		TitleData, err1 := RcApi.V1UserTitle.GetTitleByIds(ctx, &v12.UserTitleGetTitleByIdsReq{Ids: ids})
		if err1 != nil {
			log.Error("[dao.capsule | RelaodCapsuleConfig]  GetTitleByIds error")
		}
		if TitleData != nil && TitleData.Data != nil {
			titleMap = TitleData.Data
		}
	}
	for coinId, coinConf := range coinMap {
		conf := &CapsuleCoinConf{}
		conf.Status = coinConf.Status
		conf.GiftType = coinConf.GiftType
		conf.Title = coinConf.Title
		conf.EndTime = coinConf.EndTime
		conf.StartTime = coinConf.StartTime
		conf.Id = coinConf.Id
		conf.ChangeNum = coinConf.ChangeNum
		if _, ok := coinConfigMap[coinId]; ok {
			coinConfig := coinConfigMap[coinId]
			gifts := make(map[int64]struct{})
			areas := make(map[int64]struct{})
			for _, config := range coinConfig {
				if config.GiftId > 0 {
					gifts[config.GiftId] = struct{}{}
				}
				if config.AreaV2ParentId > 0 {
					areas[config.AreaV2Id] = struct{}{}
				}
			}
			conf.AreaMap = areas
			conf.GiftMap = gifts
		}
		if _, ok := poolMap[coinId]; ok {
			for _, poolConf := range poolMap[coinId] {
				pool := &CapsulePoolConf{}
				pool.Id = poolConf.Id
				pool.StartTime = poolConf.StartTime
				pool.EndTime = poolConf.EndTime
				pool.Title = poolConf.Title
				pool.Status = poolConf.Status
				pool.Rule = poolConf.Description
				pool.CoinId = poolConf.CoinId
				pool.IsBottom = poolConf.IsBottom
				if _, ok := poolPrizeMap[pool.Id]; ok {
					prizeConfigs := poolPrizeMap[pool.Id]
					pool.PoolPrize = make([]*CapsulePoolPrize, len(poolPrizeMap[pool.Id]))
					for ix, prizeConfig := range prizeConfigs {
						name := PrizeNameMap[prizeConfig.Type]
						if prizeConfig.Type == CapsulePrizeTitleType && titleMap != nil {
							if _, ok := titleMap[prizeConfig.ObjectId]; ok {
								name = titleMap[prizeConfig.ObjectId]
							}
						}
						prize := &CapsulePoolPrize{
							Id:          prizeConfig.Id,
							PoolId:      prizeConfig.PoolId,
							Type:        prizeConfig.Type,
							Num:         prizeConfig.Num,
							ObjectId:    prizeConfig.ObjectId,
							Expire:      prizeConfig.Expire,
							Name:        name,
							WebImage:    prizeConfig.WebUrl,
							MobileImage: prizeConfig.MobileUrl,
							Description: prizeConfig.Description,
							JumpUrl:     prizeConfig.JumpUrl,
							ProType:     prizeConfig.ProType,
							Chance:      prizeConfig.Chance,
							LoopNum:     prizeConfig.LoopNum,
							LimitNum:    prizeConfig.LimitNum,
							Weight:      prizeConfig.Weight,
						}
						pool.PoolPrize[ix] = prize
					}
				}
				if conf.AllPoolConf == nil {
					conf.AllPoolConf = make([]*CapsulePoolConf, 0)
				}
				conf.AllPoolConf = append(conf.AllPoolConf, pool)
			}

		}
		coinConfMap[coinId] = conf
	}
	cacheTime := time.Now().Unix()
	capsuleConf.RwLock.Lock()
	capsuleConf.CacheTime = cacheTime
	capsuleConf.ChangeFlag = changeFlag
	capsuleConf.CoinConfMap = coinConfMap
	capsuleConf.RwLock.Unlock()
	log.Info("[dao.capsule | RelaodCapsuleConfig] reload conf")
	return coinConfMap, nil
}

func getCapsuleTable(Uid int64) int64 {
	return Uid % 10
}

// UpdateScore 更新扭蛋币积分
func (d *Dao) UpdateScore(ctx context.Context, uid, coinId, score int64, action, platform string, userInfo map[int64]int64, coinConf *CapsuleCoinConf) (affect int64, err error) {
	var (
		sqlStr, uKey, iKey string
	)
	if _, ok := userInfo[coinId]; !ok {
		userInfo, _ = d.GetUserCapsuleInfo(ctx, uid)
	}
	if action == CapsuleActionTrans {
		sqlStr = fmt.Sprintf(_transUserCapsuleMysql, getCapsuleTable(uid))
	} else {
		sqlStr = fmt.Sprintf(_updateUserCapsuleMysql, getCapsuleTable(uid), CoinIdIntMap[coinId], CoinIdIntMap[coinId])
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	uKey = userKey(uid)
	iKey = userInfoKey(uid, coinId)
	affect, err = d.execSqlWithBindParams(ctx, &sqlStr, score, uid)
	if err != nil {
		log.Error("[dao.mysql_lottery|updateScore] uid(%d) type(%d) score(%d) error(%v)", uid, coinId, score, err)
		_, e := conn.Do("DEL", uKey, iKey)
		if e != nil {
			log.Error("[dao.redis_lottery|updateScore] conn.DEL(%s, %s) error(%v)", uKey, iKey, e)
		}
		return
	}
	_, e := conn.Do("DEL", uKey, iKey)
	if e != nil {
		log.Error("[dao.redis_lottery|updateScore] conn.DEL(%s, %s) error(%v)", uKey, iKey, e)
	}
	d.ReportCapsuleChange(ctx, coinId, uid, score, action, platform, userInfo, coinConf)
	return
}

// UpdateCapsule 更新扭蛋币积分
func (d *Dao) UpdateCapsule(ctx context.Context, uid, coinId, score int64, action, platform string, coinConf *CapsuleCoinConf) (affect int64, err error) {
	var (
		sqlStr, uKey, iKey string
	)
	userInfo, _ := d.GetUserInfo(ctx, uid, coinId)
	sqlStr = fmt.Sprintf(_updateUserInfoMysql, getCapsuleTable(uid))
	conn := d.redis.Get(ctx)
	defer conn.Close()
	uKey = userKey(uid)
	iKey = userInfoKey(uid, coinId)
	affect, err = d.execSqlWithBindParams(ctx, &sqlStr, score, uid, CoinIdIntMap[coinId])
	if err != nil {
		log.Error("[dao.mysql_lottery|UpdateCapsule] uid(%d) type(%d) score(%d) error(%v)", uid, coinId, score, err)
		_, e := conn.Do("DEL", uKey, iKey)
		if e != nil {
			log.Error("[dao.redis_lottery|UpdateCapsule] conn.DEL(%s, %s) error(%v)", uKey, iKey, e)
		}
		return
	}
	_, e := conn.Do("DEL", uKey, iKey)
	if e != nil {
		log.Error("[dao.redis_lottery|UpdateCapsule] conn.DEL(%s, %s) error(%v)", uKey, iKey, e)
	}
	d.ReportCapsuleChange(ctx, coinId, uid, score, action, platform, userInfo, coinConf)
	return
}

// GetUserCapsuleInfo 获取扭蛋币积分
func (d *Dao) GetUserCapsuleInfo(c context.Context, uid int64) (coinMap map[int64]int64, err error) {
	var (
		isEmpty       bool
		uKey          string
		normalScore   int64
		colorfulScore int64
	)
	uKey = userKey(uid)

	conn := d.redis.Get(c)
	defer conn.Close()
	uInfo, err := redis.Int64Map(conn.Do("HGETALL", uKey))
	if err != nil {
		if err == redis.ErrNil {
			isEmpty = true
			err = nil
		} else {
			log.Error("[dao.redis_lottery|setUserCapsuleInfoCache] getUserCapsuleInfoCache conn.HMGET(%s) error(%v)", uKey, err)
			return
		}

	} else if len(uInfo) == 0 {
		isEmpty = true
	}
	coinMap = make(map[int64]int64)
	if isEmpty {
		sqlStr := fmt.Sprintf(_userCapsuleCoinMysql, getCapsuleTable(uid))
		row := d.db.QueryRow(c, sqlStr, uid)
		err = row.Scan(&normalScore, &colorfulScore)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows {
			sqlStr := fmt.Sprintf(_addCapsuleCoinMysql, getCapsuleTable(uid))
			_, err = d.db.Exec(c, sqlStr, uid, 0, 0)
			if err != nil {
				log.Error("[dao.redis_lottery|GetUserCapsuleInfo] init sql(%s) uid(%d) error(%v)", sqlStr, uid, err)
				return
			}
			normalScore, colorfulScore = 0, 0
		}
		_, err = conn.Do("HMSET", uKey, CoinIdIntMap[NormalCoinId], normalScore, CoinIdIntMap[ColorfulCoinId], colorfulScore)
		if err != nil {
			log.Error("[dao.redis_lottery|GetUserCapsuleInfo] setUserCapsuleInfoCache conn.HMSET(%s) error(%v)", uKey, err)
		}
		coinMap[NormalCoinId] = normalScore
		coinMap[ColorfulCoinId] = colorfulScore
		err = nil
	} else {
		if _, ok := uInfo[CoinIdIntMap[NormalCoinId]]; ok {
			coinMap[NormalCoinId] = uInfo[CoinIdIntMap[NormalCoinId]]
		}
		if _, ok := uInfo[CoinIdIntMap[ColorfulCoinId]]; ok {
			coinMap[ColorfulCoinId] = uInfo[CoinIdIntMap[ColorfulCoinId]]
		}
	}
	return
}

// GetUserInfo 获取扭蛋币详情
func (d *Dao) GetUserInfo(c context.Context, uid, coinId int64) (coinMap map[int64]int64, err error) {
	if coinId <= ColorfulCoinId {
		return d.GetUserCapsuleInfo(c, uid)
	}
	var (
		isEmpty bool
		uKey    string
	)
	uKey = userInfoKey(uid, coinId)

	conn := d.redis.Get(c)
	defer conn.Close()
	score, err := redis.Int64(conn.Do("GET", uKey))
	if err != nil {
		if err == redis.ErrNil {
			isEmpty = true
			err = nil
		} else {
			log.Error("[dao.redis_lottery|setUserCapsuleInfoCache] getUserInfoCache conn.HMGET(%s) error(%v)", uKey, err)
			return
		}

	}
	coinMap = make(map[int64]int64)
	if isEmpty {
		sqlStr := fmt.Sprintf(_userInfoMysql, getCapsuleTable(uid))
		row := d.db.QueryRow(c, sqlStr, uid, CoinIdIntMap[coinId])
		err = row.Scan(&score)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows {
			sqlStr := fmt.Sprintf(_addInfoMysql, getCapsuleTable(uid))
			_, err = d.db.Exec(c, sqlStr, uid, CoinIdIntMap[coinId], 0)
			if err != nil {
				return
			}
			score = 0
		}
		_, err = conn.Do("SET", uKey, score)
		if err != nil {
			log.Error("[dao.redis_lottery|GetUserInfo] setUserCapsuleInfoCache conn.HMSET(%s) error(%v)", uKey, err)
		}
		coinMap[coinId] = score
		err = nil
	} else {
		coinMap[coinId] = score
	}
	return
}

// ReportCapsuleChange 通知扭蛋积分变换
func (d *Dao) ReportCapsuleChange(ctx context.Context, coinId, uid, score int64, action, platform string, pInfo map[int64]int64, coinConf *CapsuleCoinConf) bool {
	if _, ok := pInfo[coinId]; !ok {
		return false
	}
	chnageType := coinId
	cInfo, err := d.GetUserInfo(ctx, uid, coinId)
	if err != nil {
		return false
	}
	if _, ok := cInfo[coinId]; !ok {
		return false
	}
	change := d.GetCoin(cInfo[coinId], coinConf) - d.GetCoin(pInfo[coinId], coinConf)
	if change > 0 {
		d.AddNotice(ctx, uid, coinId, change)
	}
	var normalPreScore, colorPreScore, normalNowScore, colorNowScore int64
	if coinId <= ColorfulCoinId {
		if action == CapsuleActionTrans {
			chnageType = 5
		} else {
			if score < 0 {
				chnageType += 2
				score = -score
			}
		}
		normalPreScore, colorPreScore, normalNowScore, colorNowScore = pInfo[NormalCoinId], pInfo[ColorfulCoinId], cInfo[NormalCoinId], cInfo[ColorfulCoinId]
	} else {
		chnageType = coinId * 10
		if score < 0 {
			chnageType += 1
		}
		normalPreScore, colorPreScore, normalNowScore, colorNowScore = 0, pInfo[coinId], 0, cInfo[coinId]
	}

	date := time.Now().Format("200601")
	sqlStr := fmt.Sprintf(_reportCapsuleChangeMysql, date)
	affect, _ := d.execSqlWithBindParams(ctx, &sqlStr, uid, chnageType, score, action, platform, normalPreScore, colorPreScore, normalNowScore, colorNowScore)
	var rcontent string
	if _, ok := ReprotConfig[chnageType]; ok {
		rcontent = ReprotConfig[chnageType]
	}
	if rcontent == "" {
		if chnageType%10 == 0 {
			rcontent = "减少" + coinConf.Title
		} else if chnageType%10 == 1 {
			rcontent = "增加" + coinConf.Title
		} else if chnageType%10 == 2 {
			rcontent = "转化" + coinConf.Title
		}
	}
	report.User(&report.UserInfo{
		Platform: platform,
		Business: 101, // 101 102 103 104 105 106
		Type:     int(chnageType),
		Oid:      uid,
		Action:   "capsule_change",
		Ctime:    time.Now(),
		Index: []interface{}{
			rcontent,
			score,
			normalNowScore,
			colorNowScore,
			action,
		},
	})
	return affect > 0
}

// GetCoin 获取扭蛋数量
func (d *Dao) GetCoin(score int64, coinConf *CapsuleCoinConf) (coinNum int64) {
	if coinConf == nil || coinConf.ChangeNum == 0 {
		return 0
	}
	coinNum = score / coinConf.ChangeNum
	return
}

// AddNotice 增加标记
func (d *Dao) AddNotice(ctx context.Context, uid, coinId, coinNum int64) {
	nKey := fmt.Sprintf(_capsuleNotice, CoinIdIntMap[coinId], uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err := conn.Do("SET", nKey, coinNum, "EX", 30*86400)
	if err != nil {
		log.Error("[dao.redis_lottery|AddNotice] conn.SET(%s) error(%v)", nKey, err)
	}
}

// GetCapsuleChangeInfo 获取扭蛋配置信息
func (d *Dao) GetCapsuleChangeInfo(ctx context.Context) (int64, int64) {
	capsuleConf.RwLock.RLock()
	CacheTime := capsuleConf.CacheTime
	ChangeFlag := capsuleConf.ChangeFlag
	capsuleConf.RwLock.RUnlock()
	return CacheTime, ChangeFlag
}

// CheckLplFirstGift 检测是否首次送礼
func (d *Dao) CheckLplFirstGift(ctx context.Context, uid, giftId int64) bool {
	var (
		value int64
		extra string
		day   string
		cType string
	)
	day = time.Now().Format("2006-01-02")
	cType = "lpl" + day
	nKey := fmt.Sprintf(_lplSendGiftRedis, day, uid)

	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err := redis.Int64(conn.Do("GET", nKey))
	log.Info("[dao.redis_lottery|CheckLplFirstGift] conn.GET(%s) error(%v)", nKey, err)
	if err == redis.ErrNil {
		row := d.db.QueryRow(ctx, _getExtraDataMysql, uid, cType)
		err = row.Scan(&value, &extra)
		if err == sql.ErrNoRows {
			_, err = d.db.Exec(ctx, _addExtraDataMysql, uid, cType, time.Now().Unix(), strconv.FormatInt(giftId, 10))
			if err != nil {
				log.Error("[dao.redis_lottery|CheckLplFirstGift] conn.addExtraData(%s) error(%v)", nKey, err)
				return false
			}
			return true
		}
		if err != nil {
			log.Error("[dao.redis_lottery|CheckLplFirstGift] conn.getExtraData(%s) error(%v)", nKey, err)
			return false
		}
		_, err = conn.Do("SET", nKey, 1, "EX", 86400)
		if err != nil {
			log.Error("[dao.redis_lottery|CheckLplFirstGift] conn.SET(%s) error(%v)", nKey, err)
		}
		return false
	}
	if err != nil {
		log.Error("[dao.redis_lottery|CheckLplFirstGift] conn.GET(%s) error(%v)", nKey, err)
	}
	return false
}

// GetExtraDataByTime 获取数据
func (d *Dao) GetExtraDataByTime(ctx context.Context, startTime, endTime string) (extraData []*model.ExtraData, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getExtraDataByTimeMysql, startTime, endTime); err != nil {
		log.Error("[dao.extra | GetExtraDataByType] query(%s) error (%v)", _getExtraDataByTimeMysql, err)
		return
	}
	log.Info("[dao.extra | GetExtraDataByType] start(%d) end(%s)", startTime, endTime)
	defer rows.Close()
	extraData = make([]*model.ExtraData, 0)
	for rows.Next() {
		p := &model.ExtraData{}
		if err = rows.Scan(&p.Id, &p.Uid, &p.Type, &p.ItemValue, &p.ItemExtra); err != nil {
			log.Error("[dao.extra | GetExtraDataByType] scan error, err %v", err)
			return
		}
		extraData = append(extraData, p)
	}
	return
}

// UpdateExtraValueById 更新数据
func (d *Dao) UpdateExtraValueById(ctx context.Context, id int64, itemValue int64) (status bool, err error) {
	res, err := d.db.Exec(ctx, _updateExtraValueMysql, itemValue, id)
	if err != nil {
		log.Error("[dao.extra | UpdateExtraValue] update(%s) error(%v)", _updateExtraValueMysql, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.extra | UpdateExtraValue]  err %v", err)
		return false, err
	}
	return rows > 0, nil
}

// UpdateExtraMtimeById 更新时间数据
func (d *Dao) UpdateExtraMtimeById(ctx context.Context, id int64, mtime string) (status bool, err error) {
	res, err := d.db.Exec(ctx, _updateExtraMtimeMysql, mtime, id)
	if err != nil {
		log.Error("[dao.extra | UpdateExtraMtimeById] update(%s) error(%v)", _updateExtraMtimeMysql, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.extra | UpdateExtraMtimeById]  err %v", err)
		return false, err
	}
	return rows > 0, nil
}

// UpdateExtraById 更新数据
func (d *Dao) UpdateExtraById(ctx context.Context, id int64, itemValue int64, itemExtra string) (status bool, err error) {
	res, err := d.db.Exec(ctx, _updateExtraMysql, itemValue, itemExtra, id)
	if err != nil {
		log.Error("[dao.extra | UpdateExtraById] update(%s) error(%v)", _updateExtraMysql, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.extra | UpdateExtraById]  err %v", err)
		return false, err
	}
	return rows > 0, nil
}

// GetExtraDataByIds 获取数据
func (d *Dao) GetExtraDataByIds(ctx context.Context, start, end int64) (extraData []*model.ExtraData, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getExtraDataByIdMysql, start, end); err != nil {
		log.Error("[dao.capsule | GetExtraDataByIds]query(%s) error(%v)", _getUserInfoById, err)
		return
	}
	defer rows.Close()
	extraData = make([]*model.ExtraData, 0)
	for rows.Next() {
		p := &model.ExtraData{}
		if err = rows.Scan(&p.Id, &p.Uid, &p.Type, &p.ItemValue, &p.ItemExtra); err != nil {
			log.Error("[dao.capsule | GetExtraDataByIds] scan error, err %v", err)
			return
		}
		extraData = append(extraData, p)
	}
	return
}

// GetCouponData 添加数据
func (d *Dao) GetCouponData(ctx context.Context) (extraData []*model.ExtraData, err error) {
	var i, maxId int64
	row := d.db.QueryRow(ctx, _getExtraDataMaxIdMysql)
	if err = row.Scan(&maxId); err != nil {
		log.Error("[dao.capsule | GetCouponData] query(%s),err(%v)", _getExtraDataMaxIdMysql, err)
		return
	}
	extraData = make([]*model.ExtraData, 0)
	for i = 0; i < maxId; i = i + 10000 {
		var curExtraData []*model.ExtraData
		curExtraData, err = d.GetExtraDataByIds(ctx, i, i+10000)
		if err != nil {
			return
		}
		for _, extra := range curExtraData {
			if extra.ItemValue != 0 {
				continue
			}
			if !strings.HasPrefix(extra.Type, "CouponRetry") {
				continue
			}
			extraData = append(extraData, extra)
		}

	}
	return
}
