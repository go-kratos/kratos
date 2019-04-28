package model

import (
	"fmt"
	"go-common/library/ecode"
	"time"
)

/*
const STREAM_OP_RESULT_SUCCESS = 0              //交易成功
const STREAM_OP_RESULT_IN_PROGRESS = 1          //交易进行中
const STREAM_OP_RESULT_FAILED = 2               //交易失败
const STREAM_OP_RESULT_ROLLBACK_SUCC = 2        //回滚成功
const STREAM_OP_RESULT_ROLLBACK_IN_PROGRESS = 2 //交易回滚中
const STREAM_OP_RESULT_ROLLBACK_FAILED = 2      //回滚成功

*/

const STREAM_OP_RESULT_SUB_SUCC = 1    //扣款成功
const STREAM_OP_RESULT_ADD_SUCC = 2    //加款成功
const STREAM_OP_RESULT_SUB_FAILED = -1 //扣款失败
const STREAM_OP_RESULT_ADD_FAILED = -2 //加款失败

const STREAM_OP_REASON_EXECUTE_UNKNOWN = -6   //系统内部逻辑错误: 结果未知，当作失败
const STREAM_OP_REASON_POST_QUERY_FAILED = -5 //系统内部错误，后置查询失败
const STREAM_OP_REASON_QUERY_FAILED = -4      //系统内部错误，查询失败
const STREAM_OP_REASON_LOCK_ERROR = -3        //系统内部错误，获取锁异常
const STREAM_OP_REASON_EXECUTE_FAILED = -2    //系统内部逻辑错误: 执行失败，连接异常 OR SQL错误 OR 修改要求的条件不满足（有别的进程越过用户锁，对货币数进行了变更）
const STREAM_OP_REASON_PRE_QUERY_FAILED = -1  //前置查询失败
const STREAM_OP_REASON_NOT_ENOUGH_COIN = 1
const STREAM_OP_REASON_LOCK_FAILED = 2

const (
	LIVE_PLATFORM_IOS     = "ios"
	LIVE_PLATFORM_PC      = "pc"
	LIVE_PLATFORM_ANDROID = "android"
	LIVE_PLATFORM_H5      = "h5"
	COIN_TYPE_IOS_GOLD    = "iap_gold"
	COIN_TYPE_GOLD        = "gold"
	COIN_TYPE_SILVER      = "silver"
	COIN_TYPE_METAL       = "metal" //主站硬币（主站提供硬币数查询、硬币扣除接口供调用）
)

type RechargeOrPayForm struct {
	Uid           int64  `form:"uid" validate:"required"`
	CoinType      string `form:"coin_type" validate:"required"`
	CoinNum       int64  `form:"coin_num" validate:"required"`
	ExtendTid     string `form:"extend_tid" validate:"required"`
	Timestamp     int64  `form:"timestamp" validate:"required"`
	TransactionId string `form:"transaction_id" validate:"required"`
}

type ExchangeForm struct {
	Uid           int64  `form:"uid" validate:"required"`
	ExtendTid     string `form:"extend_tid" validate:"required"`
	Timestamp     int64  `form:"timestamp" validate:"required"`
	TransactionId string `form:"transaction_id" validate:"required"`
	SrcCoinType   string `form:"src_coin_type" validate:"required"`
	SrcCoinNum    int64  `form:"src_coin_num" validate:"required"`
	DestCoinType  string `form:"dest_coin_type" validate:"required"`
	DestCoinNum   int64  `form:"dest_coin_num" validate:"required"`
}

type RecordCoinStreamForm struct {
	Uid  int64  `form:"uid" validate:"required"`
	Data string `form:"data" validate:"required"`
}

type ServiceType int32

const (
	PAYTYPE            ServiceType = 0
	RECHARGETYPE       ServiceType = 1
	EXCHANGETYPE       ServiceType = 2
	ROLLBACKTYPE       ServiceType = 3
	SysCoinTypeIosGold int32       = 2
	SysCoinTypeGold    int32       = 1
	SysCoinTypeSilver  int32       = 0
	SysCoinTypeMetal   int32       = 3
)

func IsValidServiceType(serviceType int32) bool {
	st := ServiceType(serviceType)
	return st == PAYTYPE ||
		st == RECHARGETYPE ||
		st == EXCHANGETYPE ||
		st == ROLLBACKTYPE
}

var (
	validPlatformMap   = map[string]string{LIVE_PLATFORM_ANDROID: LIVE_PLATFORM_ANDROID, LIVE_PLATFORM_H5: LIVE_PLATFORM_H5, LIVE_PLATFORM_PC: LIVE_PLATFORM_PC, LIVE_PLATFORM_IOS: LIVE_PLATFORM_IOS}
	validCoinTypeMap   = map[string]int32{COIN_TYPE_IOS_GOLD: SysCoinTypeIosGold, COIN_TYPE_GOLD: SysCoinTypeGold, COIN_TYPE_SILVER: SysCoinTypeSilver, COIN_TYPE_METAL: SysCoinTypeMetal}
	validPlatformNoMap = map[string]int32{LIVE_PLATFORM_PC: 1, LIVE_PLATFORM_ANDROID: 2, LIVE_PLATFORM_IOS: 3, LIVE_PLATFORM_H5: 4}
)

func IsValidCoinType(coinType string) bool {
	_, ok := validCoinTypeMap[coinType]
	return ok
}

func GetCoinTypeNumber(coinType string) int32 {
	n := validCoinTypeMap[coinType]
	return n
}

func IsValidPlatform(platform string) bool {
	_, ok := validPlatformMap[platform]
	return ok
}

func IsPlatformIOS(platform string) bool {
	return platform == LIVE_PLATFORM_IOS
}

func IsLocalCoin(coinTypeNo int32) bool {
	return coinTypeNo != SysCoinTypeMetal
}
func GetSysCoinType(coinType string, platform string) string {
	if IsPlatformIOS(platform) && coinType == COIN_TYPE_GOLD {
		coinType = COIN_TYPE_IOS_GOLD
	}
	return coinType
}

func GetSysCoinTypeByNo(coinTypeNo int32) string {
	switch coinTypeNo {
	case SysCoinTypeGold:
		return COIN_TYPE_GOLD
	case SysCoinTypeIosGold:
		return COIN_TYPE_IOS_GOLD
	case SysCoinTypeSilver:
		return COIN_TYPE_SILVER
	case SysCoinTypeMetal:
		return COIN_TYPE_METAL
	default:
		return "not_define"
	}
}

func GetRechargeCnt(coinTypeNo int32) string {
	var rechargeCntField string
	if coinTypeNo == SysCoinTypeSilver {
		rechargeCntField = ""
	} else if coinTypeNo == SysCoinTypeIosGold {
		rechargeCntField = "gold_recharge_cnt"
	} else if coinTypeNo == SysCoinTypeGold {
		rechargeCntField = "gold_recharge_cnt"
	}
	return rechargeCntField
}

func GetPayCnt(coinTypeNo int32) string {
	var cntField string
	if coinTypeNo == SysCoinTypeSilver {
		cntField = "silver_pay_cnt"
	} else if coinTypeNo == SysCoinTypeIosGold {
		cntField = "gold_pay_cnt"
	} else if coinTypeNo == SysCoinTypeGold {
		cntField = "gold_pay_cnt"
	}
	return cntField
}

func GetWalletFormatTime(opTime int64) string {
	tm := time.Unix(opTime, 0)
	date := tm.Format("2006-01-02 15:04:05")
	return date
}

func NewCoinStream(uid int64, tid string, extendTid string, coinType int32, coinNum int64, opType int32, opTime int64, bizCode string, area int64, source string, bizSource string, metadata string) *CoinStreamRecord {
	return &CoinStreamRecord{
		Uid:           uid,
		TransactionId: tid,
		ExtendTid:     extendTid,
		CoinType:      coinType,
		DeltaCoinNum:  coinNum,
		OpType:        opType,
		OpTime:        opTime,
		BizCode:       bizCode,
		Area:          area,
		Source:        source,
		BizSource:     bizSource,
		MetaData:      metadata,
	}
}

func NewExchangeSteam(uid int64, tid string, srcCoinType int32, srcCoinNum int32, destCoinType int32, destCoinNum int32, opTime int64, status int32) *CoinExchangeRecord {
	return &CoinExchangeRecord{
		Uid:           uid,
		TransactionId: tid,
		SrcType:       srcCoinType,
		SrcNum:        srcCoinNum,
		DestType:      destCoinType,
		DestNum:       destCoinNum,
		ExchangeTime:  opTime,
		Status:        status,
	}
}

func (m *CoinStreamRecord) SetOpReason(r int32) {
	m.OpReason = r
}

func GetMelonseedResp(platform string, melonseed *Melonseed) *MelonseedResp {
	gold := getPlatformGold(melonseed.Gold, melonseed.IapGold, platform)
	return &MelonseedResp{
		Silver: fmt.Sprintf("%d", melonseed.Silver),
		Gold:   fmt.Sprintf("%d", gold),
	}
}

func GetMelonseedWithMetalResp(platform string, melonseed *Melonseed, metal float64) *MelonseedWithMetalResp {
	gold := getPlatformGold(melonseed.Gold, melonseed.IapGold, platform)
	return &MelonseedWithMetalResp{
		Silver: fmt.Sprintf("%d", melonseed.Silver),
		Gold:   fmt.Sprintf("%d", gold),
		Metal:  fmt.Sprintf("%.2f", metal),
	}
}

func GetDetailResp(platform string, detail *Detail) *DetailResp {
	gold := getPlatformGold(detail.Gold, detail.IapGold, platform)
	return &DetailResp{
		Silver:          fmt.Sprintf("%d", detail.Silver),
		Gold:            fmt.Sprintf("%d", gold),
		GoldPayCnt:      fmt.Sprintf("%d", detail.GoldPayCnt),
		GoldRechargeCnt: fmt.Sprintf("%d", detail.GoldRechargeCnt),
		SilverPayCnt:    fmt.Sprintf("%d", detail.SilverPayCnt),
		CostBase:        detail.CostBase,
	}
}

func GetDetailWithMetalResp(platform string, detail *Detail, metal float64) *DetailWithMetalResp {
	gold := getPlatformGold(detail.Gold, detail.IapGold, platform)
	return &DetailWithMetalResp{
		Silver:          fmt.Sprintf("%d", detail.Silver),
		Gold:            fmt.Sprintf("%d", gold),
		GoldPayCnt:      fmt.Sprintf("%d", detail.GoldPayCnt),
		GoldRechargeCnt: fmt.Sprintf("%d", detail.GoldRechargeCnt),
		SilverPayCnt:    fmt.Sprintf("%d", detail.SilverPayCnt),
		Metal:           fmt.Sprintf("%.2f", metal),
		CostBase:        detail.CostBase,
	}
}

func GetTidResp(tid string) *TidResp {
	return &TidResp{TransactionId: tid}
}

func getPlatformGold(normalGold int64, iapGold int64, platform string) int64 {
	gold := normalGold
	if IsPlatformIOS(platform) {
		gold = iapGold
	}
	return gold
}

func IncrMelonseedCoin(userCoins *Melonseed, num int64, coinTypeNo int32) {
	switch coinTypeNo {
	case SysCoinTypeIosGold:
		userCoins.IapGold += num
	case SysCoinTypeGold:
		userCoins.Gold += num
	case SysCoinTypeSilver:
		userCoins.Silver += num
	default:
	}
}

func GetCoinByMelonseed(coinTypeNo int32, userCoin *Melonseed) int64 {
	switch coinTypeNo {
	case SysCoinTypeIosGold:
		return userCoin.IapGold
	case SysCoinTypeGold:
		return userCoin.Gold
	case SysCoinTypeSilver:
		return userCoin.Silver
	default:
		return 0
	}
}

func GetCoinByDetailWithSnapShot(coinTypeNo int32, userCoin *DetailWithSnapShot) int64 {
	switch coinTypeNo {
	case SysCoinTypeIosGold:
		return userCoin.IapGold
	case SysCoinTypeGold:
		return userCoin.Gold
	case SysCoinTypeSilver:
		return userCoin.Silver
	default:
		return 0
	}
}

func CompareCoin(origin interface{}, num int64) bool {
	switch origin.(type) {
	case int64:
		return origin.(int64) >= num
	case float64:
		return int64(origin.(float64)) >= num
	default:
		return false
	}
}

// 得到数据库适配的货币数据，由于数据库的org_coin_num delta_coin_num都是整型，但是硬币的类型是浮点数，所以做一下适配
func GetDbFitCoin(v interface{}) int64 {
	switch v.(type) {
	case int64:
		return v.(int64)
	case float64:
		return int64(v.(float64))
	default:
		return 0
	}
}

func SubCoin(v1 interface{}, v2 interface{}) int64 {
	switch v1.(type) {
	case int64:
		return v1.(int64) - v2.(int64)
	case float64:
		return int64(v1.(float64) - v2.(float64))
	default:
		return 0
	}
}

func AddMoreParam2CoinStream(stream *CoinStreamRecord, bp *BasicParam, platform string) {
	platformNo := GetPlatformNo(platform)
	stream.Platform = platformNo
	stream.Reserved1 = bp.Reason
	stream.Version = bp.Version
}

type CoinStreamFieldInject interface {
	GetExtendTid() string
	GetTimestamp() int64
	GetTransactionId() string
	GetBizCode() string
	GetArea() int64
	GetBizSource() string
	GetSource() string
	GetReason() int64
	GetVersion() int64
	GetMetaData() string
	GetPlatform() string
	GetUid() int64
}

func InjectFieldToCoinStream(stream *CoinStreamRecord, inject CoinStreamFieldInject) {
	stream.ExtendTid = inject.GetExtendTid()
	stream.TransactionId = inject.GetTransactionId()
	stream.OpTime = inject.GetTimestamp()

	stream.BizCode = inject.GetBizCode()
	stream.Area = inject.GetArea()
	stream.BizSource = inject.GetBizSource()
	stream.MetaData = inject.GetMetaData()
	stream.Source = inject.GetSource()

	stream.Reserved1 = inject.GetReason()
	stream.Version = inject.GetVersion()
	platformNo := GetPlatformNo(inject.GetPlatform())
	stream.Platform = platformNo

	stream.Uid = inject.GetUid()
}

func GetPlatformNo(platform string) int32 {
	platformNo, ok := validPlatformNoMap[platform]
	if !ok {
		platformNo = 0
	}
	return platformNo
}

var (
	validRecordCoinStreamItemType = map[string]bool{"recharge": true, "pay": true}
)

func (m *RecordCoinStreamItem) IsValidType() bool {
	_, ok := validRecordCoinStreamItemType[m.Type]
	return ok
}

func (m *RecordCoinStreamItem) IsPayType() bool {
	return m.Type == "pay"
}

func (m *RecordCoinStreamItem) IsRechargeType() bool {
	return m.Type == "recharge"
}

func (m *RecordCoinStreamItem) GetOpType() int32 {
	if m.IsPayType() {
		return int32(PAYTYPE)
	} else {
		return int32(RECHARGETYPE)
	}
}

func (m *RecordCoinStreamItem) GetOpResult() int32 {
	if m.IsPayType() {
		return STREAM_OP_RESULT_SUB_SUCC
	} else {
		return STREAM_OP_RESULT_ADD_SUCC
	}
}

func (m *RecordCoinStreamItem) IsValid() (valid bool) {
	valid = false
	if m.OrgCoinNum < 0 {
		return
	}

	if !m.IsValidType() {
		return
	}

	if !IsValidCoinType(m.CoinType) {
		return
	}
	if m.IsPayType() && m.CoinNum >= 0 {
		return
	}

	if m.IsRechargeType() && m.CoinNum <= 0 {
		return
	}

	valid = true
	return
}

func GetMelonByDetailWithSnapShot(wallet *DetailWithSnapShot, platform string) (melon *MelonseedResp) {
	gold := wallet.Gold
	if platform == LIVE_PLATFORM_IOS {
		gold = wallet.IapGold
	}
	return &MelonseedResp{
		Silver: fmt.Sprintf("%d", wallet.Silver),
		Gold:   fmt.Sprintf("%d", gold),
	}
}

func ModifyCoinInDetailWithSnapShot(wallet *DetailWithSnapShot, sysCoinTypeNo int32, coinNum int64) {
	switch sysCoinTypeNo {
	case SysCoinTypeGold:
		wallet.Gold += coinNum
	case SysCoinTypeIosGold:
		wallet.IapGold += coinNum
	case SysCoinTypeSilver:
		wallet.Silver += coinNum
	}
}

// 根据锁的错误设置数据库的reason
func SetReasonByLockErr(lockErr error, coinStream *CoinStreamRecord) {
	if lockErr == ecode.TargetBlocked {
		coinStream.OpReason = STREAM_OP_REASON_LOCK_FAILED
	} else {
		coinStream.OpReason = STREAM_OP_REASON_LOCK_ERROR
	}
}

func NeedSnapshot(wallet *DetailWithSnapShot, now time.Time) bool {
	lastTime, _ := time.Parse("2006-01-02 15:04:05", wallet.SnapShotTime)
	return now.After(lastTime)
}

func GetTodayTime(now time.Time) time.Time {
	timeStr := now.Format("2006-01-02") + " 00:00:00"
	today, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return today
}

func TodayNeedSnapShot(wallet *DetailWithSnapShot) bool {
	now := GetTodayTime(time.Now())
	return NeedSnapshot(wallet, now)
}
