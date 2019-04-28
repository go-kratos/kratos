package anchorReward

import (
	"context"
	"fmt"
	"go-common/app/service/live/xrewardcenter/conf"
	AnchorTaskModel "go-common/app/service/live/xrewardcenter/model/anchorTask"
	model "go-common/app/service/live/xrewardcenter/model/anchorTask"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"math"
	"time"

	"bytes"
	"encoding/json"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

// OrderExist .
const (
	//  缓存过期时间
	rewardConfExpire = 3600
	rewardConfPrefix = "rconf_v1_%d"
)

// Dao dao
type Dao struct {
	c                   *conf.Config
	mc                  *memcache.Pool
	redis               *redis.Pool
	orm                 *gorm.DB
	db                  *xsql.DB
	keyRewardConfExpire int32
	client              *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                   c,
		mc:                  memcache.NewPool(c.Memcache),
		redis:               redis.NewPool(c.Redis),
		db:                  xsql.NewMySQL(c.MySQL),
		orm:                 orm.NewMySQL(c.ORM),
		keyRewardConfExpire: rewardConfExpire,
		client:              bm.NewClient(c.HTTPClient),
	}
	dao.initORM()
	return
}

func keyRewardConf(id int64) string {
	return fmt.Sprintf(rewardConfPrefix, id)
}

func (d *Dao) initORM() {
	d.orm.LogMode(true)
	d.orm.SingularTable(true)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -sync=true -nullcache=&model.AnchorRewardConf{ID:-1} -check_null_code=$.ID==-1
	RewardConf(c context.Context, id int64) (*model.AnchorRewardConf, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// 获取奖励配置
	// mc: -key=keyRewardConf
	CacheRewardConf(c context.Context, id int64) (*model.AnchorRewardConf, error)

	// 保存奖励配置
	// mc: -key=keyRewardConf -expire=d.keyRewardConfExpire -encode=json|gzip
	AddCacheRewardConf(c context.Context, id int64, value *model.AnchorRewardConf) error
}

// AddReward add Reward to a user.
func (d *Dao) AddReward(c context.Context, iRewardID int64, uid int64, iSource int64, iRoomid int64, iLifespan int64) (err error) {
	//aReward, _ := getRewardConfByLid(iRewardID)

	m, _ := time.ParseDuration(fmt.Sprintf("+%dh", iLifespan))

	arg := &AnchorTaskModel.AnchorReward{
		Uid:         uid,
		RewardId:    iRewardID,
		Roomid:      iRoomid,
		Source:      iSource,
		AchieveTime: xtime.Time(time.Now().Unix()),
		ExpireTime:  xtime.Time(time.Now().Add(m).Unix()),
		Status:      model.RewardUnUsed,
	}

	//spew.Dump
	// (arg)
	if err := d.orm.Create(arg).Error; err != nil {
		log.Error("addReward(%v) error(%v)", arg, err)
		return err
	}

	if err := d.SetNewReward(c, uid, int64(1)); err != nil {
		log.Error("addRewardMc(%v) error(%v)", uid, err)
	}

	if err := d.SetHasReward(c, uid, int64(1)); err != nil {
		log.Error("SetHasReward(%v) error(%v)", uid, err)
	}

	log.Info("addReward (%v) succ", arg)

	return
}

//GetByUidPage get reward by uid and page.
func (d *Dao) GetByUidPage(c context.Context, uid int64, page int64, pageSize int64, status []int64) (pager *model.AnchorRewardPager, list []*model.AnchorRewardObject, err error) {
	err = nil
	pager = &model.AnchorRewardPager{}
	list = []*model.AnchorRewardObject{}

	var (
		Items []*AnchorTaskModel.AnchorReward
		count int64
	)
	iOffSet := (page - 1) * pageSize

	db := d.orm.Where("status in (?)", status).Where("uid=?", uid)
	db.Model(&AnchorTaskModel.AnchorReward{}).Count(&count)

	if err = db.Model(&AnchorTaskModel.AnchorReward{}).Limit(pageSize).Offset(iOffSet).Order("mtime DESC, id").Find(&Items).Error; err != nil {
		log.Error("get ap_anchor_task_reward_list uid(%v) error(%v)", uid, err)
		return
	}

	for _, v := range Items {
		aReward, err := d.RewardConf(c, v.RewardId)
		if err != nil {
			log.Error("RewardConf(%v) error(%v)", v.RewardId, err)
			return pager, list, err
		}

		if aReward == nil {
			continue
		}

		aListItem := &model.AnchorRewardObject{
			Id:          v.Id,
			RewardType:  aReward.RewardType,
			Status:      v.Status,
			RewardId:    v.RewardId,
			Name:        aReward.Name,
			Icon:        aReward.Icon,
			AchieveTime: v.AchieveTime.Time().Format("2006-01-02 15:04:05"),
			ExpireTime:  v.ExpireTime.Time().Format("2006-01-02 15:04:05"),
			UseTime:     v.UseTime.Time().Format("2006-01-02 15:04:05"),
			Source:      v.Source,
			RewardIntro: aReward.RewardIntro,
		}

		list = append(list, aListItem)
	}
	pager = &model.AnchorRewardPager{
		Page:       page,
		PageSize:   pageSize,
		TotalPage:  int64(math.Ceil(float64(count) / float64(pageSize))),
		TotalCount: count,
	}

	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if d.orm != nil {
		d.orm.DB().PingContext(c)
	}

	err = d.pingMC(c)
	return
}

// GetById get reward by id.
func (d *Dao) GetById(id int64) (reward *model.AnchorReward, err error) {
	rewards := []*model.AnchorReward{}
	if err := d.orm.Model(&model.AnchorReward{}).Find(&rewards, "id=?", id).Error; err != nil {
		log.Error("getRewardById (%v) error(%v)", id, err)
		return reward, err
	}
	if len(rewards) != 0 {
		reward = rewards[0]
	}

	return
}

// UseReward use reward by id.
func (d *Dao) UseReward(id int64, usePlat string) (rst bool, err error) {
	if err := d.orm.
		Model(&model.AnchorReward{}).
		Where("id=?", id).
		Update(map[string]interface{}{"status": model.RewardUsed, "use_plat": usePlat, "use_time": xtime.Time(time.Now().Unix())}).Error; err != nil {
		log.Error("useReward (%v) error(%v)", id, err)
		return rst, err
	}
	rst = true
	return
}

// HasNewReward .
func (d *Dao) HasNewReward(c context.Context, uid int64) (rst int64, err error) {
	rst, _ = d.GetNewReward(c, uid)
	return
}

func (d *Dao) findByUid(uid int64, limitOne bool) (reward *model.AnchorReward, err error) {
	rewards := []*model.AnchorReward{}
	db := d.orm.Where("uid=?", uid)

	if limitOne {
		db = db.Limit(1)
	}

	if err := db.Model(&model.AnchorReward{}).Find(&rewards).Error; err != nil {
		log.Error("getRewardById (%v) error(%v)", uid, err)
		return reward, err
	}
	if len(rewards) != 0 {
		reward = rewards[0]
	}

	return
}

// HasReward returns if a user have reward.
func (d *Dao) HasReward(c context.Context, uid int64) (r int64, err error) {
	rst, err := d.GetHasReward(c, uid)
	if err != nil {
		if err == memcache.ErrNotFound {
			reward, err2 := d.findByUid(uid, true)
			if err2 != nil {
				return rst, err2
			}
			if reward != nil {
				rst = int64(1)
				d.SetHasReward(c, uid, rst)
			} else {
				rst = int64(0)
				d.SetHasReward(c, uid, rst)
			}
			return rst, err
		}
		log.Error("HasReward(%v) error(%v)", uid, err)
		return rst, err
	}
	return rst, err
}

// CheckOrderID check orderid is valid.
func (d *Dao) CheckOrderID(c context.Context, id string) (exist int64, err error) {
	exist = 0
	if exist, err = d.GetOrder(c, id); err != nil {
		//spew.Dump(exist, err)
		if err == memcache.ErrNotFound {
			err = nil
		}
		return exist, err
	}
	return exist, err
}

// SaveOrderID save order id.
func (d *Dao) SaveOrderID(c context.Context, id string) error {
	err := d.SaveOrder(c, id)
	return err
}

// SetExpire .
func (d *Dao) SetExpire(now time.Time) (err error) {
	var (
		db = d.orm
	)

	setMap := map[string]interface{}{
		"status": model.RewardExpired,
	}

	if err = db.Model(model.AnchorReward{}).
		Where("status=? AND reward_id = ? AND expire_time <= ?", model.RewardUnUsed, 1, now.Format("2006-01-02 15:04:05")).
		Update(setMap).
		Error; err != nil {
		log.Error("SetExpire (%v) error(%v)", setMap, err)
		return err
	}
	return
}

// CountExpire .
func (d *Dao) CountExpire(interval int64, now time.Time) (err error) {
	var (
		c      = context.TODO()
		db     = d.orm
		result = &[]model.AnchorReward{}
	)

	dur, _ := time.ParseDuration("-" + strconv.FormatInt(interval, 10) + "s")
	begin := now.Add(dur)

	//spew.Dump(begin.Format("2006-01-02 15:04:05"))
	//spew.Dump(now.Format("2006-01-02 15:04:05"))
	sqlTemp :=
		"SELECT * FROM ap_anchor_task_reward_list WHERE status = ? AND expire_time > ? AND expire_time <= ? AND reward_id= ?"

	db.Raw(sqlTemp,
		model.RewardExpired,
		begin.Format("2006-01-02 15:04:05"),
		now.Format("2006-01-02 15:04:05"),
		1).Scan(&result)

	for _, v := range *result {
		d.AddExpireCountCache(c, fmt.Sprintf(model.CountExpireUserKey, v.Uid), model.ExpireCountTime)
	}
	return
}

// SendBroadcastV2 .
func (d *Dao) SendBroadcastV2(c context.Context, uid int64, roomid int64, rewardId int64) (err error) {
	log.Info("send reward broadcast begin:%d", roomid)

	var endPoint string = fmt.Sprintf("http://live-dm.bilibili.co/dm/1/push?cid=%d&ensure=1", roomid)

	postJson := make(map[string]interface{})
	postJson["cmd"] = "new_anchor_reward"
	postJson["uid"] = uid
	postJson["roomid"] = roomid
	postJson["reward_id"] = rewardId

	bytesData, err := json.Marshal(postJson)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", postJson, err)
		return
	}

	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(bytesData))

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Error("http.NewRequest(%v) url(%v) error(%v)", postJson, endPoint, err)
		return
	}

	var v interface{}

	if err = d.client.Do(c, req, v); err != nil {
		log.Error("s.client.Do error(%v) res (%v)", err, v)
		return
	}
	log.Info("s.client.Do endpoint (%v) req (%v) res (%v)", endPoint, postJson, v)

	return
}

// SendBroadcast .
func (d *Dao) SendBroadcast(uid int64, roomid int64, rewardId int64) (err error) {
	log.Info("send reward broadcast begin:%d", roomid)

	var endPoint = fmt.Sprintf("http://live-dm.bilibili.co/dm/1/push?cid=%d&ensure=1", roomid)

	postJson := make(map[string]interface{})
	postJson["cmd"] = "new_anchor_reward"
	postJson["uid"] = uid
	postJson["roomid"] = roomid
	postJson["reward_id"] = rewardId

	bytesData, err := json.Marshal(postJson)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", postJson, err)
		return
	}

	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(bytesData))

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Error("http.NewRequest(%v) url(%v) error(%v)", postJson, endPoint, err)
		return
	}

	client := http.Client{
		Timeout: time.Second,
	}

	// use httpClient to send request
	response, err := client.Do(req)

	if err != nil {
		log.Error("sending request to API endpoint(%v) error(%v)", req, err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("parse resp body(%v) error(%v)", body, err)
	}

	log.Info("send reward broadcast end:%d", roomid)

	return
}
