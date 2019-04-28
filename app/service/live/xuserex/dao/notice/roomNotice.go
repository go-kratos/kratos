package notice

import (
	"context"
	"go-common/app/service/live/xuser/api/grpc/v1"
	v1pb "go-common/app/service/live/xuserex/api/grpc/v1"
	"go-common/app/service/live/xuserex/conf"
	"go-common/library/cache/memcache"

	"bytes"
	"encoding/json"
	"fmt"
	"go-common/app/service/live/resource/sdk"
	"go-common/app/service/live/xuserex/model/roomNotice"
	"go-common/library/database/hbase.v2"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
	"strconv"
	"time"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	mc *memcache.Pool
	// acc rpc
	client *bm.Client
	hbase  *hbase.Client
	xuser  *v1.Client
}

type guardGuideConf struct {
	Open      int64 `json:"open"`
	Threshold int64 `json:"threshold"`
}

const (
	//  缓存过期时间
	keyShouldNoticeExpire = 3600
	keyNoticePre          = "kn_v1_%d_%d_%s"
	// HBaseMonthlyConsumeTable .
	HBaseMonthlyConsumeTable = "livemonthconsume"
	// BuyGuardGuideTitanKey .
	BuyGuardGuideTitanKey = "buy_guard_guide"
	// TaskFinishKey .
	TaskFinishKey = "task_finish_%s"
)

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		mc:     memcache.NewPool(c.Memcache),
		client: bm.NewClient(c.BMClient),
		hbase:  hbase.NewClient(c.HBase),
	}

	conn, err := v1.NewClient(c.Warden)
	if err != nil {
		panic(err)
	}
	dao.xuser = conn
	return
}

func keyShouldNotice(UID int64, targetID int64, date string) string {
	return fmt.Sprintf(keyNoticePre, UID, targetID, date)
}

// IsNotice returns whether should pop a purchase notice.
func (dao *Dao) IsNotice(c context.Context, UID int64, targetID int64) (*v1pb.RoomNoticeBuyGuardResp, error) {
	term := dao.GetTermBegin()
	begin := term.Unix()
	end := dao.GetTermEnd()

	resp := &v1pb.RoomNoticeBuyGuardResp{
		Begin:   begin,
		End:     end.Unix(),
		Now:     time.Now().Unix(),
		Title:   "感谢支持主播",
		Content: "成为船员为主播保驾护航吧～",
		Button:  "开通大航海",
	}

	shouldNotice, err := dao.getShouldNotice(c, UID, targetID, term)
	if err != nil {
		log.Error("dao getShouldNotice uid(%v)roomid(%v)term(%v) error(%v)", UID, targetID, term.Format("2006-01-02"), err)
		err = nil
		return resp, err
	}
	resp.ShouldNotice = int64(shouldNotice)
	return resp, nil
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -sync=true -nullcache=&roomNotice.MonthConsume{Amount:-1} -check_null_code=$.Amount==-1
	MonthConsume(c context.Context, UID int64, targetID int64, date string) (*roomNotice.MonthConsume, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// 获取某个月消费
	// mc: -key=keyShouldNotice
	CacheMonthConsume(c context.Context, UID int64, targetID int64, date string) (*roomNotice.MonthConsume, error)

	// 保存获取某个月消费
	// mc: -key=keyShouldNotice -expire=d.keyShouldNoticeExpire -encode=json|gzip
	AddCacheMonthConsume(c context.Context, UID int64, targetID int64, date string, value *roomNotice.MonthConsume) error
}

func (dao *Dao) getThreshold() (threshold *guardGuideConf, err error) {
	threshold = &guardGuideConf{}

	guideConf, err := titansSdk.Get(BuyGuardGuideTitanKey)
	log.Info("getThreshold_key(%v) conf(%+v)", BuyGuardGuideTitanKey, guideConf)
	if err != nil {
		log.Error("getThreshold(%v) error(%v)", BuyGuardGuideTitanKey, err)
		return
	}

	if "" == guideConf {
		return
	}

	if err = json.Unmarshal([]byte(guideConf), threshold); err != nil {
		log.Error("json Unmarshal guideconf(%+v) error(%v)", guideConf, err)
		return
	}
	log.Info("getThreshold_unmarshal_succ key(%v) conf (%v) Threshold(%+v)", BuyGuardGuideTitanKey, guideConf, threshold)
	return
}

func (dao *Dao) getShouldNotice(ctx context.Context, UID int64, targetID int64, term time.Time) (shouldNotice int, err error) {
	shouldNotice = 0
	taskFinish, err := dao.GetTaskFinish(ctx, term)
	if err != nil {
		return shouldNotice, err
	}
	if !taskFinish {
		log.Info("task_not_finish")
		return shouldNotice, err
	}

	// 获取配置的收入门槛
	threshold, err := dao.getThreshold()
	if err != nil {
		log.Error("get_threshold_error(%v)", err)
		return
	}

	if nil == threshold {
		log.Error("get_threshold_nil")
		return
	}

	if 0 == threshold.Open {
		log.Info("guard_guide not Open (%+v)", threshold)
		return
	}

	monthConsume, err := dao.MonthConsume(ctx, UID, targetID, dao.termToString(term))
	log.Info("get_monthConsume(%+v) Threshold (%+v)", monthConsume, threshold)

	if err != nil {
		log.Error("get_monthConsum_err uid(%d) targetid(%v) term (%v) error(%v)", UID, targetID, term, err)
		return
	}

	if nil == monthConsume {
		return
	}

	if int64(monthConsume.Amount) >= threshold.Threshold*1000 { // coin to rmb
		isGuard, err := dao.isGuard(ctx, UID, targetID)
		log.Info("show guard guide uid(%v) target (%v) guard (%v) Threshold (%v)", UID, targetID, isGuard, threshold)
		if err != nil {
			log.Error("get gaurd UID(%v) targetid (%v) error(%v)", UID, targetID, err)
			return shouldNotice, err
		}
		if !isGuard {
			shouldNotice = 1
		}
	}
	return
}

func hbaseRowKey(UID int64, targetID int64, date string) string {
	return fmt.Sprintf("%s_%d_%d", date, UID, targetID)
}

//RawMonthConsume get month consume from hbase
func (dao *Dao) RawMonthConsume(ctx context.Context, UID int64, targetID int64, date string) (res *roomNotice.MonthConsume, err error) {
	var (
		tableName = HBaseMonthlyConsumeTable
		key       = hbaseRowKey(UID, targetID, date)
	)
	result, err := dao.hbase.GetStr(ctx, tableName, key)
	log.Info("RawMonthConsume_getstr tableName (%v) key (%v) res (%v)", tableName, key, result)

	if err != nil {
		log.Error("dao.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, UID, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil {
		return
	}
	res = &roomNotice.MonthConsume{}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		v, _ := strconv.ParseInt(string(c.Value[:]), 10, 64)
		if !bytes.Equal(c.Family, []byte("info")) {
			log.Error("family_type_err(%v) error", c.Family)

			continue
		}
		switch {
		case bytes.Equal(c.Qualifier, []byte("uid")):
			res.Uid = v
		case bytes.Equal(c.Qualifier, []byte("ruid")):
			res.Ruid = v
		case bytes.Equal(c.Qualifier, []byte("amount")):
			res.Amount = v
		case bytes.Equal(c.Qualifier, []byte("time")):
			res.Date = v
		}
	}
	log.Info("RawMonthConsume_succ uid (%v) target (%v) date (%v) res (%+v)", UID, targetID, date, res)
	return
}

// GetTermBegin return first day of last month
func (dao *Dao) GetTermBegin() time.Time {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth.AddDate(0, -1, 0)
}

// GetTermEnd returns last second of last month
func (dao *Dao) GetTermEnd() time.Time {

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	second, _ := time.ParseDuration("-1s")
	return thisMonth.Add(second)
}

// IsValidTerm returns whether a term is valid.
func (dao *Dao) IsValidTerm(term time.Time) bool {

	return true
}

// TOOD
func (dao *Dao) isGuard(c context.Context, UID int64, targetID int64) (isGuard bool, err error) {
	ret, err := dao.xuser.GetByUIDTargetID(c, &v1.GetByUidTargetIdReq{
		Uid:      UID,
		TargetId: targetID,
	})

	log.Info("dao.xuser.GetByUIDTargetID uid (%v) target (%v) res (%v)", UID, targetID, ret)
	if err != nil {
		log.Error("get guard uid (%v) target (%v) error(%v)", UID, targetID, err)
		return
	}

	if nil == ret || nil == ret.Data || 0 == len(ret.Data) {
		log.Info("not_guard uid (%v) target (%v) res (%v)", UID, targetID, ret)
		return
	}
	isGuard = true
	return
}

// GetTaskFinish .
func (dao *Dao) GetTaskFinish(c context.Context, term time.Time) (isOn bool, err error) {
	conn := dao.mc.Get(c)
	defer conn.Close()
	key := dao.keyTaskFinish(term)
	reply, err := conn.Get(key)
	log.Info("GetTaskFinish key (%v) term (%v)", key, term)
	if err != nil {
		if err == memcache.ErrNotFound {
			log.Info("GetTaskFinish_not_found key (%v) term (%v)", key, term)
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:GetTaskFinish")
		log.Error("GetTaskFinish_fail key(%v) error(%v)", key, err)
		return
	}

	res := &roomNotice.TaskFinish{}
	err = conn.Scan(reply, &res)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:GetTaskFinish")
		log.Error("GetTaskFinish_fail_scan key(%v) error(%v)", key, err)
		return
	}
	log.Info("GetTaskFinish_succ key (%v) term (%v) res(%+v)", key, term, res)
	if res == nil {
		return
	}
	if 1 == res.Finish {
		isOn = true
	}
	return
}

// SetTaskFinish .
func (dao *Dao) SetTaskFinish(c context.Context, term time.Time, isFinish int64) (err error) {
	conn := dao.mc.Get(c)
	defer conn.Close()
	key := dao.keyTaskFinish(term)
	value := &roomNotice.TaskFinish{
		Finish: isFinish,
	}
	log.Info("SetTaskFinish key (%v) term (%v) value (%+v)", key, term, value)
	item := &memcache.Item{
		Key:        key,
		Object:     value,
		Flags:      memcache.FlagJSON,
		Expiration: 0,
	}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:SetTaskFinish")
		log.Error("SetTaskFinish_fail key(%v) value (%+v) error(%v)", key, value, err)
		return
	}
	log.Info("SetTaskFinish_succ key (%v) term (%v) value (%+v)", key, term, value)
	return
}

func (dao *Dao) termToString(term time.Time) string {
	return term.Format("20060102")
}

func (dao *Dao) keyTaskFinish(term time.Time) (key string) {
	return fmt.Sprintf(TaskFinishKey, dao.termToString(term))
}
