package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	v12 "go-common/app/service/live/rc/api/liverpc/v1"

	"github.com/pkg/errors"

	"go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_userCapsuleCoinMysql     = "select normal_score, colorful_score from capsule_%d where uid = ?"
	_addCapsuleCoinMysql      = "insert into capsule_%d(uid, normal_score, colorful_score) values(?, ?, ?)"
	_userInfoMysql            = "select score from capsule_info_%d where uid = ? and type = ?"
	_addInfoMysql             = "insert into capsule_info_%d(uid,type,score) values(?, ?, ?)"
	_updateUserCapsuleMysql   = "update capsule_%d set %s_score = %s_score + ? where uid = ?"
	_updateUserInfoMysql      = "update capsule_info_%d set score = score + ? where uid = ? and type = ?"
	_transUserCapsuleMysql    = "update capsule_%d set colorful_score = 0, normal_score = normal_score + ? where uid = ?"
	_reportCapsuleChangeMysql = "insert into capsule_log_%s(uid, type, score, action, platform, pre_normal_score, pre_colorful_score, cur_normal_score, cur_colorful_score) values(?,?,?,?,?,?,?,?,?)"
	_userCapsuleCoinRedis     = "hash:capsule:user:%d"
	_userInfoRedis            = "capsule:user:info:%d:%d"
	_openHistoryRedis         = "list:capsule:%s:list"
	_openLockRedis            = "capsule-pay-%d"
	_capsuleConfRand          = "capsule:rand"
	_historyOpenCount         = "hash:capsule:count"
	_historyGiftCount         = "capsule:gift:count:%d:%s"
	_capsuleNotice            = "capsule:notice:%s:%d"
	_whiteUserPrizeRedis      = "capsule:white:user:%d"
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
	// CapsulePrizeScoreAdd 积分加成卡
	CapsulePrizeScoreAdd
	// CapsulePrizeSmallStar 小星星
	CapsulePrizeSmallStar
	// CapsulePrizeWeekScore 抽奖券
	CapsulePrizeWeekScore
	// CapsulePrizeDanmuColor 弹幕颜色
	CapsulePrizeDanmuColor
	// CapsulePrizeLplScore lpl抽奖券
	CapsulePrizeLplScore
	// CapsulePrizeLplProduct1 lpl实物奖励1
	CapsulePrizeLplProduct1
	// CapsulePrizeLplProduct2 lpl实物奖励2
	CapsulePrizeLplProduct2
	// CapsulePrizeLplProduct3 lpl实物奖励3
	CapsulePrizeLplProduct3
)

const (
	// CapsulePrizeProduct1 .
	CapsulePrizeProduct1 = 100000 // 实物奖励
	// CapsulePrizeProduct2 .
	CapsulePrizeProduct2 = 100001 // 实物奖励
	// CapsulePrizeProduct3 .
	CapsulePrizeProduct3 = 100002
	// CapsulePrizeProduct4 .
	CapsulePrizeProduct4 = 100003 // 实物奖励
	// CapsulePrizeProduct5 .
	CapsulePrizeProduct5 = 100004 // 实物奖励
	// CapsulePrizeProduct6 .
	CapsulePrizeProduct6 = 100005
)

const (
	// CapsulePrizeCoupon1 .
	CapsulePrizeCoupon1 = 200000
	// CapsulePrizeCoupon2 .
	CapsulePrizeCoupon2 = 200001
	// CapsulePrizeCoupon3 .
	CapsulePrizeCoupon3 = 200002
)

const (
	// CapsulePrizeExpire1Day 过期时间24小时
	CapsulePrizeExpire1Day = 1
	// CapsulePrizeExpire3Day 过期时间72小时
	CapsulePrizeExpire3Day = 10
	// CapsulePrizeExpire1Week 过期时间1周
	CapsulePrizeExpire1Week = 20
	// CapsulePrizeExpire3Month 过期时间3个月
	CapsulePrizeExpire3Month = 30
	// CapsulePrizeExpireForever 过期时间永久
	CapsulePrizeExpireForever = 100
)
const (
	_ = iota
	// ProTypeNormal 概率
	ProTypeNormal
	// ProTypeFixDay 每天固定数量
	ProTypeFixDay
	// ProTypeFixWeek 每周固定数量
	ProTypeFixWeek
	// ProTypeWhite 白名单
	ProTypeWhite
)
const (
	// CapsuleGiftTypeAll gift_type 为全部道具
	CapsuleGiftTypeAll = 1
)
const (
	// NormalCoinId 普通扭蛋id
	NormalCoinId = 1
	// ColorfulCoinId 梦幻扭蛋id
	ColorfulCoinId = 2
	// WeekCoinId 梦幻扭蛋id
	WeekCoinId = 3
	// LplCoinId 梦幻扭蛋id
	LplCoinId = 4
	// BlessCoinId 祈福券
	BlessCoinId = 5
	// OpenHistoryNum 开奖历史
	OpenHistoryNum = 30
	// NormalCoinString 普通扭蛋字符串标识，数据库和redis
	NormalCoinString = "normal"
	// ColorfulCoinString 梦幻扭蛋字符串标识，数据库和redis
	ColorfulCoinString = "colorful"
	// WeekCoinString 周星扭蛋字符串标识，数据库和redis
	WeekCoinString = "week"
	// LplCoinString lpl扭蛋字符串标识，数据库和redis
	LplCoinString = "lpl"
	// BlessCoinString 新年祈愿扭蛋字符串标识，数据库和redis
	BlessCoinString = "bless"
	// GetCapsuleDetailFromRoom 接口来源
	GetCapsuleDetailFromRoom = "room"
	// GetCapsuleDetailFromWeb 接口来源
	GetCapsuleDetailFromWeb = "web"
	// GetCapsuleDetailFromH5 接口来源
	GetCapsuleDetailFromH5 = "h5"
)

const (
	//IsBottomPool 是保底奖池
	IsBottomPool = 1
	//CapsuleActionTrans 转换扭蛋
	CapsuleActionTrans = "trans"
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
	WhiteUserMap                                      map[int64]struct{}
}

// HistoryOpenInfo 开奖历史
type HistoryOpenInfo struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Date string `json:"date"`
	Num  int64  `json:"num"`
}

var (
	// CoinIdIntMap map
	CoinIdIntMap map[int64]string
	// CoinIdStringMap map
	CoinIdStringMap map[string]int64
	// ReprotConfig map
	ReprotConfig map[int64]string
	// PrizeNameMap map
	PrizeNameMap map[int64]string
	// PrizeExpireMap map
	PrizeExpireMap map[int64]string
	// UnLockGetWrong flag
	UnLockGetWrong = "UnLockGetWrong"
	// ErrUnLockGet error
	ErrUnLockGet  = errors.New(UnLockGetWrong)
	capsuleConf   CapsuleConf
	whitePrizeMap sync.Map
)

func init() {
	CoinIdIntMap = make(map[int64]string)
	CoinIdIntMap[NormalCoinId] = NormalCoinString
	CoinIdIntMap[ColorfulCoinId] = ColorfulCoinString
	CoinIdIntMap[WeekCoinId] = WeekCoinString
	CoinIdIntMap[LplCoinId] = LplCoinString
	CoinIdIntMap[BlessCoinId] = BlessCoinString
	CoinIdStringMap = make(map[string]int64)
	CoinIdStringMap[NormalCoinString] = NormalCoinId
	CoinIdStringMap[ColorfulCoinString] = ColorfulCoinId
	CoinIdStringMap[WeekCoinString] = WeekCoinId
	CoinIdStringMap[LplCoinString] = LplCoinId
	CoinIdStringMap[BlessCoinString] = BlessCoinId
	ReprotConfig = make(map[int64]string)
	ReprotConfig[0] = "未知"
	ReprotConfig[1] = "增加普通扭蛋"
	ReprotConfig[2] = "增加梦幻扭蛋"
	ReprotConfig[3] = "减少普通扭蛋"
	ReprotConfig[4] = "减少梦幻扭蛋"
	ReprotConfig[5] = "梦幻转化普通"
	PrizeNameMap = make(map[int64]string)
	PrizeNameMap[CapsulePrizeGift1Type] = "辣条"
	PrizeNameMap[CapsulePrizeTitleType] = "头衔"
	PrizeNameMap[CapsulePrizeStuff1Type] = "经验原石"
	PrizeNameMap[CapsulePrizeStuff2Type] = "经验曜石"
	PrizeNameMap[CapsulePrizeStuff3Type] = "贤者之石"
	PrizeNameMap[CapsulePrizeSmallTvType] = "小电视抱枕"
	PrizeNameMap[CapsulePrizeGuard3Type] = "舰长体验"
	PrizeNameMap[CapsulePrizeGuard2Type] = "提督体验"
	PrizeNameMap[CapsulePrizeScoreAdd] = "积分加成卡"
	PrizeNameMap[CapsulePrizeSmallStar] = "小星星"
	PrizeNameMap[CapsulePrizeWeekScore] = "抽奖券"
	PrizeNameMap[CapsulePrizeDanmuColor] = "金色弹幕"
	PrizeNameMap[CapsulePrizeLplScore] = "LPL抽奖券"
	PrizeNameMap[CapsulePrizeLplProduct1] = "2019拜年祭小电视猪"
	PrizeNameMap[CapsulePrizeLplProduct2] = "小电视毛绒公仔"
	PrizeNameMap[CapsulePrizeLplProduct3] = "机械之心桌垫"
	PrizeNameMap[CapsulePrizeProduct1] = "迎新礼盒+新年台历"
	PrizeNameMap[CapsulePrizeProduct2] = "新年锦鲤围巾"
	PrizeNameMap[CapsulePrizeProduct3] = "拜年祭挂件+拜年祭立牌+年夜饭挂画"
	PrizeNameMap[CapsulePrizeProduct4] = "拜年祭耳罩"
	PrizeNameMap[CapsulePrizeProduct5] = "小电视猪挂件+小电视猪公仔"
	PrizeNameMap[CapsulePrizeProduct6] = "拜年祭徽章"
	PrizeNameMap[CapsulePrizeCoupon1] = "会员购20元优惠券"
	PrizeNameMap[CapsulePrizeCoupon2] = "会员购40元优惠券"
	PrizeNameMap[CapsulePrizeCoupon3] = "会员购60元优惠券"
	PrizeExpireMap = make(map[int64]string)
	PrizeExpireMap[CapsulePrizeExpire1Day] = "1天"
	PrizeExpireMap[CapsulePrizeExpire3Day] = "3天"
	PrizeExpireMap[CapsulePrizeExpire1Week] = "1周"
	PrizeExpireMap[CapsulePrizeExpire3Month] = "3个月"
	PrizeExpireMap[CapsulePrizeExpireForever] = "永久"
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
			log.Error("[dap.capsule | GetCapsuleConf] CapsuleCoinConf is empty")
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

func getCapsuleTable(Uid int64) int64 {
	return Uid % 10
}

func userKey(uid int64) string {
	return fmt.Sprintf(_userCapsuleCoinRedis, uid)
}

func userInfoKey(uid int64, coinId int64) string {
	return fmt.Sprintf(_userInfoRedis, uid, coinId)
}

func openHistoryKey(coinType string) string {
	return fmt.Sprintf(_openHistoryRedis, coinType)
}

func openTotalCountKey() string {
	return _historyOpenCount
}

func openGiftCountKey(giftId int64, day string) string {
	return fmt.Sprintf(_historyGiftCount, giftId, day)
}

func whiteUserPrizeKey(uid int64) string {
	return fmt.Sprintf(_whiteUserPrizeRedis, uid)
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
		coinMap, err = d.GetUserCapsuleInfo(c, uid)
		return
	}
	var (
		isEmpty bool
		uKey    string
		score   int64
	)
	uKey = userInfoKey(uid, coinId)

	conn := d.redis.Get(c)
	defer conn.Close()
	score, err = redis.Int64(conn.Do("GET", uKey))
	if err != nil {
		if err == redis.ErrNil {
			isEmpty = true
			err = nil
		} else {
			log.Error("[dao.redis_lottery|GetUserInfo] getUserInfoCache conn.HMGET(%s) error(%v)", uKey, err)
			return
		}

	}
	coinMap = make(map[int64]int64)
	if isEmpty {
		sqlStr := fmt.Sprintf(_userInfoMysql, getCapsuleTable(uid))
		row := d.db.QueryRow(c, sqlStr, uid, CoinIdIntMap[coinId])
		err = row.Scan(&score)
		if err != nil && err != sql.ErrNoRows {
			log.Error("[dao.redis_lottery|GetUserInfo] getUserInfoFromDB uid(%d) error(%v)", uid, err)
			return
		}
		if err == sql.ErrNoRows {
			sqlStr := fmt.Sprintf(_addInfoMysql, getCapsuleTable(uid))
			_, err = d.db.Exec(c, sqlStr, uid, CoinIdIntMap[coinId], 0)
			if err != nil {
				log.Error("[dao.redis_lottery|GetUserInfo] init sql(%s) uid(%d) error(%v)", sqlStr, uid, err)
				return nil, err
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

// GetOpenHistory 获取扭蛋币历史
func (d *Dao) GetOpenHistory(c context.Context, coinType int64) (ret []*HistoryOpenInfo, err error) {
	hkey := openHistoryKey(CoinIdIntMap[coinType])
	conn := d.redis.Get(c)
	defer conn.Close()
	jsons, err := redis.Strings(conn.Do("LRANGE", hkey, 0, OpenHistoryNum-1))
	if err != nil {
		return
	}
	length := len(jsons)
	if length == 0 {
		return
	}
	ret = make([]*HistoryOpenInfo, length)
	for ix, jsonStr := range jsons {
		var openInfo HistoryOpenInfo
		json.Unmarshal([]byte(jsonStr), &openInfo)
		ret[ix] = &openInfo
	}
	return
}

// GetCoin 获取扭蛋数量
func (d *Dao) GetCoin(score int64, coinConf *CapsuleCoinConf) (coinNum int64) {
	if coinConf == nil || coinConf.ChangeNum == 0 {
		return 0
	}
	coinNum = score / coinConf.ChangeNum
	return
}

// GetProgress 获取扭蛋币进度
func (d *Dao) GetProgress(score int64, coinConf *CapsuleCoinConf) (process *v1.Progress) {
	process = &v1.Progress{}
	if coinConf == nil || coinConf.ChangeNum == 0 {
		return
	}
	process.Max = coinConf.ChangeNum
	process.Now = score % coinConf.ChangeNum
	return
}

// UpdateScore 更新扭蛋币积分
func (d *Dao) UpdateScore(ctx context.Context, uid, coinId, score int64, action, platform string, coinConf *CapsuleCoinConf) (affect int64, err error) {
	var (
		sqlStr, uKey, iKey string
	)
	if action == CapsuleActionTrans {
		sqlStr = fmt.Sprintf(_transUserCapsuleMysql, getCapsuleTable(uid))
	} else {
		sqlStr = fmt.Sprintf(_updateUserCapsuleMysql, getCapsuleTable(uid), CoinIdIntMap[coinId], CoinIdIntMap[coinId])
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	uKey = userKey(uid)
	iKey = userInfoKey(uid, coinId)
	log.Info("trace UpdateScore start")
	affect, err = d.execSqlWithBindParams(ctx, &sqlStr, score, uid)
	log.Info("trace UpdateScore end")
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
	return
}

// UpdateCapsule 更新扭蛋币积分
func (d *Dao) UpdateCapsule(ctx context.Context, uid, coinId, score int64, action, platform string, coinConf *CapsuleCoinConf) (affect int64, err error) {
	var (
		sqlStr, uKey, iKey string
	)
	sqlStr = fmt.Sprintf(_updateUserInfoMysql, getCapsuleTable(uid))
	conn := d.redis.Get(ctx)
	defer conn.Close()
	uKey = userKey(uid)
	iKey = userInfoKey(uid, coinId)
	log.Info("trace UpdateCapsule start")
	affect, err = d.execSqlWithBindParams(ctx, &sqlStr, score, uid, CoinIdIntMap[coinId])
	log.Info("trace UpdateCapsule end")
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
	return
}

// ReportCapsuleChange 扭蛋流水
func (d *Dao) ReportCapsuleChange(ctx context.Context, coinId, uid, score int64, action, platform string, pInfo, cInfo map[int64]int64, coinConf *CapsuleCoinConf) bool {
	if _, ok := pInfo[coinId]; !ok {
		return false
	}
	if _, ok := cInfo[coinId]; !ok {
		return false
	}
	chnageType := coinId
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

// PayCoin 支付扭蛋币
func (d *Dao) PayCoin(ctx context.Context, uid int64, coinConf *CapsuleCoinConf, openCount int64, action, platform string) (status int64, pInfo map[int64]int64, err error) {
	lockKey := fmt.Sprintf(_openLockRedis, uid)
	isGetLock, lockString, err := d.Lock(ctx, lockKey, 10000, 0, 0)
	if err != nil || !isGetLock {
		return
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	userData, err := d.GetUserInfo(ctx, uid, coinConf.Id)
	if err != nil {
		d.UnLock(ctx, lockKey, lockString)
		return
	}
	var score int64
	if _, ok := userData[coinConf.Id]; ok {
		score = userData[coinConf.Id]
	}
	coinNum := d.GetCoin(score, coinConf)
	if coinNum < openCount {
		d.UnLock(ctx, lockKey, lockString)
		return 1, userData, nil

	}
	value := openCount * coinConf.ChangeNum
	_, err = d.UpdateScore(ctx, uid, coinConf.Id, -value, action, platform, coinConf)
	if err != nil {
		d.UnLock(ctx, lockKey, lockString)
		return
	}
	d.UnLock(ctx, lockKey, lockString)
	return 0, userData, nil
}

// PayCapsule 支付扭蛋币
func (d *Dao) PayCapsule(ctx context.Context, uid int64, coinConf *CapsuleCoinConf, openCount int64, action, platform string) (status int64, pInfo map[int64]int64, err error) {
	lockKey := fmt.Sprintf(_openLockRedis, uid)
	isGetLock, lockString, err := d.Lock(ctx, lockKey, 10000, 0, 0)
	if err != nil {
		return
	}
	if !isGetLock {
		return 1, nil, nil
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	userData, err := d.GetUserInfo(ctx, uid, coinConf.Id)
	if err != nil {
		d.UnLock(ctx, lockKey, lockString)
		return
	}
	var score int64
	if _, ok := userData[coinConf.Id]; ok {
		score = userData[coinConf.Id]
	}
	coinNum := d.GetCoin(score, coinConf)
	if coinNum < openCount {
		d.UnLock(ctx, lockKey, lockString)
		return 2, userData, nil

	}
	value := openCount * coinConf.ChangeNum
	_, err = d.UpdateCapsule(ctx, uid, coinConf.Id, -value, action, platform, coinConf)
	if err != nil {
		d.UnLock(ctx, lockKey, lockString)
		return
	}
	d.UnLock(ctx, lockKey, lockString)
	return 0, userData, nil
}

// IsPoolOpen 判断扭蛋池是否开启
func (d *Dao) IsPoolOpen(coinConf *CapsuleCoinConf, coinId int64) bool {
	if coinConf == nil {
		return false
	}
	if coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
		return false
	}
	now := time.Now().Unix()
	if coinConf.PoolConf.StartTime < now && coinConf.PoolConf.EndTime > now {
		return true
	}
	return false
}

// GetGift 获取扭蛋奖池奖品
func (d *Dao) GetGift(ctx context.Context, coinId int64) (gift []*CapsulePoolPrize, err error) {
	coinConfMap, err := d.GetCapsuleConf(ctx)
	if err != nil || len(coinConfMap) == 0 {
		return
	}
	if _, ok := coinConfMap[coinId]; !ok {
		return
	}
	conf := coinConfMap[coinId]
	if conf.PoolConf == nil || len(conf.PoolConf.PoolPrize) == 0 {
		return
	}
	gift = conf.PoolConf.PoolPrize
	return
}

// IncrOpenCount 增加开奖次数
func (d *Dao) IncrOpenCount(ctx context.Context, coinId int64) (cnt int64) {
	hkey := openTotalCountKey()
	conn := d.redis.Get(ctx)
	defer conn.Close()
	cnt, err := redis.Int64(conn.Do("HINCRBY", hkey, CoinIdIntMap[coinId], 1))
	if err != nil {
		return
	}
	if cnt > 0 {
		return cnt
	}
	return 0
}

// GetOpenCount 获取开奖次数
func (d *Dao) GetOpenCount(ctx context.Context, coinId int64) (cnt int64) {
	hkey := openTotalCountKey()
	conn := d.redis.Get(ctx)
	defer conn.Close()
	cnt, err := redis.Int64(conn.Do("HGET", hkey, CoinIdIntMap[coinId]))
	if err != nil {
		return 0
	}
	if cnt > 0 {
		return cnt
	}
	return 0
}

func getWhiteGift(coinConf *CapsuleCoinConf) (fixPrize []*CapsulePoolPrize) {
	if coinConf == nil || coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
		return
	}
	fLen := 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeWhite {
			fLen++
		}
	}
	if fLen <= 0 {
		return
	}
	fixPrize = make([]*CapsulePoolPrize, fLen)
	fLen = 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeWhite {
			fixPrize[fLen] = prize
			fLen++
		}
	}
	return
}

func getFixGift(coinConf *CapsuleCoinConf) (fixPrize []*CapsulePoolPrize) {
	if coinConf == nil || coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
		return
	}
	fLen := 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeFixDay || prize.ProType == ProTypeFixWeek {
			fLen++
		}
	}
	if fLen <= 0 {
		return
	}
	fixPrize = make([]*CapsulePoolPrize, fLen)
	fLen = 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeFixDay || prize.ProType == ProTypeFixWeek {
			fixPrize[fLen] = prize
			fLen++
		}
	}
	return
}

func getRandomGift(coinConf *CapsuleCoinConf) (randomPrize []*CapsulePoolPrize) {
	if coinConf == nil || coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
		return
	}
	rLen := 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeNormal {
			rLen++
		}
	}
	if rLen <= 0 {
		return
	}
	randomPrize = make([]*CapsulePoolPrize, rLen)
	rLen = 0
	for _, prize := range coinConf.PoolConf.PoolPrize {
		if prize.ProType == ProTypeNormal {
			randomPrize[rLen] = prize
			rLen++
		}
	}
	return
}

func (d *Dao) checkWhiteLimit(ctx context.Context, uid int64, prize *CapsulePoolPrize) bool {
	if prize == nil || prize.ProType != ProTypeWhite || len(prize.WhiteUserMap) == 0 {
		return false
	}
	if _, ok := prize.WhiteUserMap[uid]; !ok {
		return false
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	dtime := time.Now()
	day := dtime.Format("2006-01-02")
	uKey := whiteUserPrizeKey(uid)
	lastTime, err := redis.Int64(conn.Do("GET", uKey))
	if err != nil {
		if err == redis.ErrNil {
			// 回源数据库
			prizeLog, err := d.GetUserPrizeLog(ctx, prize.Id, uid)
			if err != nil {
				return false
			}
			if prizeLog != nil {
				lastTime = prizeLog.Timestamp
			} else {
				lastTime = 0
			}
			conn.Do("SET", uKey, lastTime, "EX", 30*86400)
		} else {
			return false
		}
	}
	if dtime.Unix()-lastTime < 7*86400 {
		return false
	}
	gKey := openGiftCountKey(prize.Id, day)
	isGetLock, lockString, errLock := d.Lock(ctx, gKey, 1000000, 0, 0)
	if errLock != nil || !isGetLock {
		return false
	}
	mKey := day + strconv.FormatInt(prize.Id, 10)
	_, ok := whitePrizeMap.Load(mKey)
	if ok {
		d.UnLock(ctx, gKey, lockString)
		return false
	}
	// 回源数据库
	prizeLog, errDb := d.GetPrizeDayLog(ctx, prize.Id, day)
	if errDb != nil {
		d.UnLock(ctx, gKey, lockString)
		return false
	}
	if prizeLog != nil {
		whitePrizeMap.Store(mKey, prizeLog.Uid)
		d.UnLock(ctx, gKey, lockString)
		return false
	}
	stutus, errAdd := d.AddPrizeData(ctx, prize.Id, uid, day, dtime.Unix())
	if errAdd != nil || !stutus {
		d.UnLock(ctx, gKey, lockString)
		return false
	}
	whitePrizeMap.Store(mKey, uid)
	conn.Do("DEL", uKey)
	d.UnLock(ctx, gKey, lockString)
	return true

}

func (d *Dao) checkPrizeLimit(ctx context.Context, prize *CapsulePoolPrize) bool {
	if prize == nil {
		return false
	}
	if prize.ProType != ProTypeFixDay && prize.ProType != ProTypeFixWeek {
		return false
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()
	var status bool
	switch prize.ProType {
	case ProTypeFixDay:
		day := time.Now().Format("2006-01-02")
		gKey := openGiftCountKey(prize.Id, day)
		cnt, err := redis.Int64(conn.Do("INCRBY", gKey, 1))
		if err != nil {
			status = false
			break
		}
		status = cnt <= prize.LimitNum
	case ProTypeFixWeek:
		wDay := time.Now().Weekday()
		if wDay == 0 {
			wDay = 7
		}
		diff := time.Duration(wDay)
		day := time.Now().Add(time.Second * 86400 * (diff - 1)).Format("2006-01-02")
		gKey := openGiftCountKey(prize.Id, day)
		cnt, err := redis.Int64(conn.Do("INCRBY", gKey, 1))
		if err != nil {
			status = false
			break
		}
		status = cnt <= prize.LimitNum
	default:
		status = false
	}
	return status
}

func (d *Dao) getWhiteAward(ctx context.Context, uid int64, coinConf *CapsuleCoinConf) (prize *CapsulePoolPrize) {
	whitePrize := getWhiteGift(coinConf)
	if len(whitePrize) == 0 {
		return
	}
	openCount := d.GetOpenCount(ctx, coinConf.Id)
	if openCount == 0 {
		return
	}
	for _, wprize := range whitePrize {
		if !d.checkWhiteLimit(ctx, uid, wprize) {
			continue
		}
		return &CapsulePoolPrize{
			Id:          wprize.Id,
			PoolId:      wprize.PoolId,
			Type:        wprize.Type,
			Num:         wprize.Num,
			Name:        wprize.Name,
			WebImage:    wprize.WebImage,
			MobileImage: wprize.MobileImage,
			Description: wprize.Description,
			ProType:     wprize.ProType,
			JumpUrl:     wprize.JumpUrl,
			ObjectId:    wprize.ObjectId,
			Expire:      wprize.Expire,
			Weight:      wprize.Weight,
		}
	}
	return
}

func (d *Dao) getFixAward(ctx context.Context, uid, openCount int64, coinConf *CapsuleCoinConf) (prize *CapsulePoolPrize) {
	fixPrize := getFixGift(coinConf)
	if len(fixPrize) == 0 {
		return
	}
	if openCount == 0 {
		return
	}
	var loop int64
	for _, fprize := range fixPrize {
		loop = fprize.LoopNum
		if openCount%loop != 0 {
			continue
		}
		if d.checkPrizeLimit(ctx, fprize) {
			return &CapsulePoolPrize{
				Id:          fprize.Id,
				PoolId:      fprize.PoolId,
				Type:        fprize.Type,
				Num:         fprize.Num,
				Name:        fprize.Name,
				WebImage:    fprize.WebImage,
				MobileImage: fprize.MobileImage,
				Description: fprize.Description,
				ProType:     fprize.ProType,
				JumpUrl:     fprize.JumpUrl,
				ObjectId:    fprize.ObjectId,
				Expire:      fprize.Expire,
				Weight:      fprize.Weight,
			}
		}
	}
	return
}

func (d *Dao) getRandomAward(ctx context.Context, uid int64, coinConf *CapsuleCoinConf) (prize *CapsulePoolPrize) {
	randomPrize := getRandomGift(coinConf)
	if len(randomPrize) == 0 {
		return
	}
	var start, random, total int64
	for _, prize := range randomPrize {
		if prize == nil {
			rbyte, _ := json.Marshal(randomPrize)
			log.Error("[dao.capsule | getRandomAward] randomPrize error : %s", string(rbyte))
			continue
		}
		total += prize.Chance
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random = r.Int63n(total - 1)
	for _, prize := range randomPrize {
		if random >= start && random < start+prize.Chance {
			return &CapsulePoolPrize{
				Id:          prize.Id,
				PoolId:      prize.PoolId,
				Type:        prize.Type,
				Num:         prize.Num,
				Name:        prize.Name,
				WebImage:    prize.WebImage,
				MobileImage: prize.MobileImage,
				Description: prize.Description,
				ProType:     prize.ProType,
				JumpUrl:     prize.JumpUrl,
				ObjectId:    prize.ObjectId,
				Expire:      prize.Expire,
				Weight:      prize.Weight,
			}
		}
		start += prize.Chance
	}
	return
}

// OpenCapsule 开启扭蛋
func (d *Dao) OpenCapsule(ctx context.Context, uid int64, coinConf *CapsuleCoinConf, iTime, openCount int64, isGetFixAward bool, entryMap map[int64]bool) (award *CapsulePoolPrize) {
	if iTime == 0 {
		award = d.getWhiteAward(ctx, uid, coinConf)
		if award != nil {
			return award
		}
	}
	if isGetFixAward {
		award = d.getFixAward(ctx, uid, openCount, coinConf)
		if award != nil {
			if _, ok := entryMap[award.Id]; !ok {
				return
			}
		}
	}
	award = d.getRandomAward(ctx, uid, coinConf)
	if award != nil {
		return
	}
	return
}

// LogAward 记录抽奖奖励
func (d *Dao) LogAward(ctx context.Context, uid int64, coinId int64, awards []*CapsulePoolPrize) {
	if len(awards) == 0 {
		return
	}
	hkey := openHistoryKey(CoinIdIntMap[coinId])
	day := time.Now().Format("2006-01-02")
	logs := make([]interface{}, len(awards)+1)
	logs[0] = hkey
	ll := 1
	for _, award := range awards {
		if award.Type == CapsulePrizeGift1Type {
			continue
		}
		info := HistoryOpenInfo{Uid: uid, Name: award.Name, Num: award.Num, Date: day}
		b, err := json.Marshal(info)
		if err == nil {
			logs[ll] = string(b)
			ll++
		}
	}
	if ll == 1 {
		return
	}
	logs = logs[0:ll]
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("LPUSH", logs...)
}

// AddNotice 增加标记
func (d *Dao) AddNotice(ctx context.Context, uid, coinId, coinNum int64) {
	nKey := fmt.Sprintf(_capsuleNotice, CoinIdIntMap[coinId], uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("SET", nKey, coinNum, 30*86400)

}

// ClearNotice 清除标记
func (d *Dao) ClearNotice(ctx context.Context, uid, coinId int64) {
	nKey := fmt.Sprintf(_capsuleNotice, CoinIdIntMap[coinId], uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("DEL", nKey)
}

// ClearNoticeBoth 清除标记
func (d *Dao) ClearNoticeBoth(ctx context.Context, uid int64) {
	keys := make([]interface{}, len(CoinIdIntMap))
	var i = 0
	for _, coinType := range CoinIdIntMap {
		keys[i] = fmt.Sprintf(_capsuleNotice, coinType, uid)
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("DEL", keys...)
}

// GetChangeNum 获取扭蛋变化数量
func (d *Dao) GetChangeNum(ctx context.Context, uid, coinId int64) int64 {
	nKey := fmt.Sprintf(_capsuleNotice, CoinIdIntMap[coinId], uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	change, err := redis.Int64(conn.Do("GET", nKey))
	if err != nil {
		return 0
	}
	return change
}

// SetCapsuleChangeFlag 设置扭蛋配置变化标记
func (d *Dao) SetCapsuleChangeFlag(ctx context.Context) (status string, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	status, err = redis.String(conn.Do("SET", _capsuleConfRand, time.Now().Unix()))
	if err != nil {
		log.Error("[dao.capsule | SetCapsuleChangeFlag] redis set error : %v", err)
		return
	}
	return
}

// GetCapsuleChangeFlag 获取扭蛋配置变化标记
func (d *Dao) GetCapsuleChangeFlag(ctx context.Context) (changeFlag int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	changeFlag, err = redis.Int64(conn.Do("GET", _capsuleConfRand))
	if err != nil {
		return
	}
	return
}

// GetCapsuleChangeInfo 获取扭蛋配置信息
func (d *Dao) GetCapsuleChangeInfo(ctx context.Context) (int64, int64) {
	capsuleConf.RwLock.RLock()
	CacheTime := capsuleConf.CacheTime
	ChangeFlag := capsuleConf.ChangeFlag
	capsuleConf.RwLock.RUnlock()
	return CacheTime, ChangeFlag
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
	prizeIds := make([]int64, 0)
	idMap := make(map[int64]struct{})
	for _, prizeList := range poolPrizeMap {
		for _, prize := range prizeList {
			if prize.ObjectId != 0 && prize.Type == CapsulePrizeTitleType {
				if _, ok := idMap[prize.ObjectId]; !ok {
					ids = append(ids, prize.ObjectId)
					idMap[prize.ObjectId] = struct{}{}
				}
			}
			prizeIds = append(prizeIds, prize.Id)
		}
	}
	prizeWhiteMap, err1 := d.GetWhiteUserMap(ctx, prizeIds)
	if err1 != nil {
		log.Error("[dao.capsule | RelaodCapsuleConfig]  GetWhiteUserMap error")
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
						prize.WhiteUserMap = make(map[int64]struct{})
						if prize.ProType == ProTypeWhite {
							if _, ok := prizeWhiteMap[prize.Id]; ok {
								if len(prizeWhiteMap[prize.Id]) > 0 {
									for _, wuid := range prizeWhiteMap[prize.Id] {
										prize.WhiteUserMap[wuid] = struct{}{}
									}
								}
							}
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

// GetBottomPrize 获取保底奖品
func (d *Dao) GetBottomPrize(ctx context.Context, coinConf *CapsuleCoinConf) (bottomPrize *CapsulePoolPrize) {
	if coinConf == nil || coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
		return nil
	}

	for _, prize := range coinConf.PoolConf.PoolPrize {
		if bottomPrize == nil || bottomPrize.Weight > prize.Weight {
			bottomPrize = &CapsulePoolPrize{
				Id:          prize.Id,
				PoolId:      prize.PoolId,
				Type:        prize.Type,
				Name:        prize.Name,
				WebImage:    prize.WebImage,
				MobileImage: prize.MobileImage,
				Description: prize.Description,
				ProType:     prize.ProType,
				JumpUrl:     prize.JumpUrl,
				ObjectId:    prize.ObjectId,
				Expire:      prize.Expire,
				Weight:      prize.Weight,
				Num:         prize.Num,
			}
		}
	}
	return bottomPrize
}

//GetExpireTime 获取过期时间
func (d *Dao) GetExpireTime(expire int64) time.Time {
	var td time.Time
	if expire == CapsulePrizeExpire1Day {
		year, month, day := time.Now().Date()
		td = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Add(86400 * time.Second).Add(86400 * time.Second)
	} else if expire == CapsulePrizeExpire3Day {
		year, month, day := time.Now().Date()
		td = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Add(86400 * time.Second).Add(3 * 86400 * time.Second)
	} else if expire == CapsulePrizeExpire1Week {
		year, month, day := time.Now().Date()
		td = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Add(86400 * time.Second).Add(6 * 86400 * time.Second)
	} else if expire == CapsulePrizeExpire3Month {
		year, month, day := time.Now().Date()
		td = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Add(86400 * time.Second).Add(90 * 86400 * time.Second)
	} else if expire == CapsulePrizeExpireForever {
		td = time.Unix(0, 0)
	} else {
		td = time.Unix(0, 0)
	}
	return td
}

// IsAwardEntry 是否是实物奖励
func (d *Dao) IsAwardEntry(awardType int64) bool {
	if awardType == CapsulePrizeSmallTvType || awardType == CapsulePrizeLplProduct1 || awardType == CapsulePrizeLplProduct2 || awardType == CapsulePrizeLplProduct3 {
		return true
	}
	if awardType >= CapsulePrizeProduct1 && awardType < CapsulePrizeCoupon1 {
		return true
	}
	return false
}

// IsAwardCoupon 是否是会员券
func (d *Dao) IsAwardCoupon(awardType int64) bool {
	return awardType >= CapsulePrizeCoupon1
}
