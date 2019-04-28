package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/live/xlottery/internal/model"
	"go-common/library/ecode"
	"go-common/library/queue/databus/report"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron"

	"go-common/app/job/live/xlottery/internal/conf"
	"go-common/app/job/live/xlottery/internal/dao"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c                    *conf.Config
	dao                  *dao.Dao
	cron                 *cron.Cron
	giftPaySub           *databus.Databus
	giftFreeSub          *databus.Databus
	capsuleSub           *databus.Databus
	ExpireCountFrequency string
	CouponRetryFrequency string
	httpClient           *bm.Client
	wg                   *sync.WaitGroup
}

const _sendGiftKey = "lottery:gift:msgid:%s"

const _addCapsuleKey = "lottery:gift:msgid:%s"

type info struct {
	MsgContent string `json:"msg_content"`
}

type msgContent struct {
	Body *body `json:"body"`
}
type body struct {
	GiftId    int64     `json:"giftid"`
	RoomId    int64     `json:"roomid"`
	Num       int64     `json:"num"`
	Uid       int64     `json:"uid"`
	Ruid      int64     `json:"ruid"`
	TotalCoin int64     `json:"totalCoin"`
	CoinType  string    `json:"coinType"`
	Tid       string    `json:"tid"`
	Platform  string    `json:"platform"`
	RoomInfo  *roomInfo `json:"roomInfo"`
}
type roomInfo struct {
	AreaV2Id       int64 `json:"area_v2_id"`
	AreaV2ParentId int64 `json:"area_v2_parent_id"`
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                    c,
		dao:                  dao.New(c),
		cron:                 cron.New(),
		giftPaySub:           databus.New(c.GiftPaySub),
		giftFreeSub:          databus.New(c.GiftFreeSub),
		capsuleSub:           databus.New(c.AddCapsuleSub),
		wg:                   new(sync.WaitGroup),
		ExpireCountFrequency: c.Cfg.ExpireCountFrequency,
		CouponRetryFrequency: c.Cfg.CouponRetryFrequency,
		httpClient:           bm.NewClient(c.HTTPClient),
	}
	report.InitUser(conf.Conf.UserReport)
	dao.InitAPI()
	s.addCrontab()
	s.cron.Start()
	s.tickerReloadCapsuleConf(context.TODO())
	log.Info("[service.lottery| 11start")
	var i int64
	for i = 0; i < c.Cfg.ConsumerProcNum; i++ {
		s.wg.Add(1)
		go s.giftConsumeProc()
	}
	s.wg.Add(1)
	go s.capsuleConsumeProc()
	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.subClose()
	s.wg.Wait()
	s.dao.Close()
}

// subClose Close all sub channels
func (s *Service) subClose() {
	s.giftPaySub.Close()
	s.giftFreeSub.Close()
	s.capsuleSub.Close()
}

func (s *Service) addCrontab() (err error) {
	//spew.Dump(s.ExpireCountFrequency)
	err = s.cron.AddFunc(s.ExpireCountFrequency, s.TransCapsule)
	if err != nil {
		log.Error("cron job transCapsule error(%v)", err)
	}
	err = s.cron.AddFunc(s.CouponRetryFrequency, s.CouponRetry)
	if err != nil {
		log.Error("cron job couponRetry error(%v)", err)
	}
	return
}

// CouponRetry 抽奖券重试
func (s *Service) CouponRetry() {
	var ctx = context.Background()
	if s.c.CouponConf == nil || s.c.CouponConf.Url == "" || len(s.c.CouponConf.Coupon) == 0 {
		log.Error("[service.capsule | sendAward] couponConf is empty")
		return
	}
	nowTime := time.Now()
	log.Info("[service.service | couponRetry]couponRetry %s", nowTime.Format("2006-01-02 15:04:05"))
	extraData, _ := s.dao.GetCouponData(ctx)
	if len(extraData) == 0 {
		return
	}
	for _, extra := range extraData {
		s.dao.UpdateExtraMtimeById(ctx, extra.Id, nowTime.Format("2006-01-02 15:04:05"))
		awardType := extra.ItemExtra
		if _, ok := s.c.CouponConf.Coupon[awardType]; !ok {
			log.Error("[service.capsule | sendAward] couponConf.coupon is empty %s", awardType)
			continue
		}
		uid := extra.Uid
		var res struct {
			Code int    `json:"code"`
			Msg  string `json:"message"`
		}
		endPoint := s.c.CouponConf.Url
		postJson := make(map[string]interface{})
		postJson["mid"] = uid
		postJson["couponId"] = s.c.CouponConf.Coupon[awardType]
		bytesData, err := json.Marshal(postJson)
		if err != nil {
			log.Error("[service.capsule | sendAward] json.Marshal(%v) error(%v)", postJson, err)
			continue
		}
		req, err := http.NewRequest("POST", endPoint, bytes.NewReader(bytesData))
		if err != nil {
			log.Error("[service.capsule | sendAward] http.NewRequest(%v) url(%v) error(%v)", postJson, endPoint, err)
			continue
		}
		req.Header.Add("Content-Type", "application/json;charset=UTF-8")
		log.Info("coupon vip mid(%d) couponID(%s)", uid, s.c.CouponConf.Coupon[awardType])
		if err = s.httpClient.Do(ctx, req, &res); err != nil {
			log.Error("[service.capsule | sendAward] s.client.Do error(%v)", err)
			continue
		}
		if res.Code != 0 && res.Code != 83110005 {
			err = ecode.Int(res.Code)
			log.Error("coupon vip url(%v) res code(%d)", endPoint, res.Code)
			continue
		}
		log.Info("[service.capsule | sendAward] s.client.Do endpoint (%v) req (%v)", endPoint, postJson)
		s.dao.UpdateExtraValueById(ctx, extra.Id, 1)
	}

}

// TransCapsule 转换扭蛋币
func (s *Service) TransCapsule() {
	var ctx = context.Background()
	pools, err := s.dao.GetActiveColorPool(ctx)
	if err != nil {
		log.Error("[service.service | TransCapsule]CronJob TransCapsule GetActiveColorPool error(%v)", err)
		return
	}
	nowTime := time.Now().Add(-(60 * time.Second)).Format("2006-01-02 15:04")
	log.Info("[service.service | TransCapsule]TranCapsule %s", nowTime)
	flag := 0
	coinId := int64(0)
	for _, pool := range pools {
		if pool.EndTime == 0 {
			continue
		} else {
			endTimeUnix := time.Unix(pool.EndTime, 0)
			endTime := endTimeUnix.Format("2006-01-02 15:04")
			if endTime == nowTime {
				flag = 1
				coinId = pool.CoinId
			}
		}
	}
	if flag == 1 {
		colorChangeNum, err := s.dao.GetTransNum(ctx, coinId)
		if err != nil || colorChangeNum == 0 {
			log.Error("[service.service | TransCapsule] GetTransNum colorChangeNum: %d, err: %v", colorChangeNum, err)
			return
		}
		normalChangeNum, err := s.dao.GetTransNum(ctx, dao.NormalCoinId)
		if err != nil || normalChangeNum == 0 {
			log.Error("[service.service | TransCapsule] GetTransNum normalChangeNum: %d, err: %v", normalChangeNum, err)
			return
		}
		for i := int64(0); i < 10; i++ {
			err := s.dao.TransCapsule(ctx, strconv.FormatInt(i, 10), colorChangeNum, normalChangeNum)
			if err != nil {
				log.Error("[service.service | TransCapsule]TranCapsule error %v", err)
				return
			}
			log.Info("[service.service | TransCapsule]TranCapsule %s", strconv.FormatInt(i, 10))
		}
	}
}

// expCanalConsumeproc consumer archive
func (s *Service) giftConsumeProc() {
	defer func() {
		log.Warn("giftConsumeProc exited.")
		s.wg.Done()
	}()
	var (
		payMsgs  = s.giftPaySub.Messages()
		freeMsgs = s.giftFreeSub.Messages()
	)
	log.Info("[service.lottery|giftConsumeProc")
	for {
		select {
		case msg, ok := <-payMsgs:
			if !ok {
				log.Warn("[service.lottery|giftConsumeProc] giftPaySub has been closed.")
				return
			}

			var value *info
			var subValue *msgContent
			err := json.Unmarshal([]byte(msg.Value), &value)
			if err != nil {
				log.Error("[service.lottery|giftConsumeProc] giftPaySub json decode error:%v", err)
				continue
			}
			err = json.Unmarshal([]byte(value.MsgContent), &subValue)
			if err != nil {
				log.Error("[service.lottery|giftConsumeProc] giftPaySub json decode error:%v", err)
				continue
			}
			areaV2Id := subValue.Body.RoomInfo.AreaV2Id
			areaV2ParentId := subValue.Body.RoomInfo.AreaV2ParentId
			giftId := subValue.Body.GiftId
			roomId := subValue.Body.RoomId
			num := subValue.Body.Num
			uid := subValue.Body.Uid
			ruid := subValue.Body.Ruid
			totalCoin := subValue.Body.TotalCoin
			coinType := subValue.Body.CoinType
			platform := subValue.Body.Platform
			key := fmt.Sprintf(_sendGiftKey, subValue.Body.Tid)
			isGetLock, _, err := s.dao.Lock(context.Background(), key, 86400000, 0, 0)
			if err != nil || !isGetLock {
				log.Error("[service.lottery|giftConsumeProc Lock Error msgKey(%s) uid(%d) ruid(%d) roomId(%d) giftId(%d) num(%d) totalCoin(%d) coinType(%s) tid(%s) key(%s) offset(%d) err(%v)", msg.Key, uid, ruid, roomId, giftId, num, totalCoin, coinType, subValue.Body.Tid, msg.Key, msg.Offset, err)
				continue
			}
			msg.Commit()
			log.Info("[service.lottery|giftConsumeProc] pay-msgKey(%s) uid(%d) ruid(%d) roomId(%d) giftId(%d) num(%d) totalCoin(%d) coinType(%s) tid(%s) key(%s) offset(%d)", msg.Key, uid, ruid, roomId, giftId, num, totalCoin, coinType, subValue.Body.Tid, msg.Key, msg.Offset)
			s.sendGift(context.Background(), uid, giftId, num, totalCoin, coinType, areaV2ParentId, areaV2Id, platform)
		case msg, ok := <-freeMsgs:
			if !ok {
				log.Warn("[service.lottery|giftConsumeProc] giftFreeSub has been closed.")
				return
			}
			var value *info
			var subValue *msgContent
			err := json.Unmarshal([]byte(msg.Value), &value)
			if err != nil {
				log.Error("[service.lottery|giftConsumeProc] giftFreeSub message:%s json decode error:%v", msg.Value, err)
				continue
			}
			err = json.Unmarshal([]byte(value.MsgContent), &subValue)
			if err != nil {
				log.Error("[service.lottery|giftConsumeProc] giftFreeSub message:%s json decode error:%v", msg.Value, err)
				continue
			}
			areaV2Id := subValue.Body.RoomInfo.AreaV2Id
			areaV2ParentId := subValue.Body.RoomInfo.AreaV2ParentId
			giftId := subValue.Body.GiftId
			roomId := subValue.Body.RoomId
			num := subValue.Body.Num
			uid := subValue.Body.Uid
			ruid := subValue.Body.Ruid
			totalCoin := subValue.Body.TotalCoin
			coinType := subValue.Body.CoinType
			platform := subValue.Body.Platform
			key := fmt.Sprintf(_sendGiftKey, subValue.Body.Tid)
			isGetLock, _, err := s.dao.Lock(context.Background(), key, 86400000, 0, 0)
			if err != nil || !isGetLock {
				log.Error("[service.lottery|giftConsumeProc Lock Error msgKey(%s) uid(%d) ruid(%d) roomId(%d) giftId(%d) num(%d) totalCoin(%d) coinType(%s) tid(%s) key(%s) offset(%d) err(%v)", msg.Key, uid, ruid, roomId, giftId, num, totalCoin, coinType, subValue.Body.Tid, msg.Key, msg.Offset, err)
				continue
			}
			msg.Commit()
			log.Info("[service.lottery|giftConsumeProc] pay-msgKey(%s) uid(%d) ruid(%d) roomId(%d) giftId(%d) num(%d) totalCoin(%d) coinType(%s) tid(%s) key(%s) offset(%d)", msg.Key, uid, ruid, roomId, giftId, num, totalCoin, coinType, subValue.Body.Tid, msg.Key, msg.Offset)
			s.sendGift(context.Background(), uid, giftId, num, totalCoin, coinType, areaV2ParentId, areaV2Id, platform)
		default:
			time.Sleep(time.Second * 3)
			continue
		}
	}
}

func (s *Service) capsuleConsumeProc() {
	defer func() {
		log.Warn("capsuleConsumeProc exited.")
		s.wg.Done()
	}()
	var (
		capsuleMsgs = s.capsuleSub.Messages()
	)
	log.Info("[service.lottery|capsuleConsumeProc")
	for {
		select {
		case msg, ok := <-capsuleMsgs:
			if !ok {
				log.Warn("[service.lottery|capsuleConsumeProc] giftPaySub has been closed.")
				return
			}
			var msgContent *info
			var value *model.AddCapsule
			err := json.Unmarshal([]byte(msg.Value), &msgContent)
			if err != nil {
				log.Error("[service.lottery|capsuleConsumeProc] json decode error:%v", err)
				continue
			}
			err = json.Unmarshal([]byte(msgContent.MsgContent), &value)
			if err != nil {
				log.Error("[service.lottery|giftConsumeProc] giftFreeSub message:%s json decode error:%v", msg.Value, err)
				continue
			}

			uid := value.Uid
			cType := value.Type
			coinId := value.CoinId
			num := value.Num
			key := fmt.Sprintf(_addCapsuleKey, value.MsgId)
			isGetLock, _, err := s.dao.Lock(context.Background(), key, 86400000, 0, 0)
			if err != nil || !isGetLock {
				log.Error("[service.lottery|capsuleConsumeProc Lock Error msgKey(%s) uid(%d) num(%d) type(%s) coinId(%d) tid(%s) offset(%d) err(%v)", msg.Key, uid, num, cType, coinId, value.MsgId, msg.Offset, err)
				continue
			}
			msg.Commit()
			log.Info("[service.lottery|capsuleConsumeProc] msgKey(%s) uid(%d) num(%d) type(%s) coinId(%s)  tid(%s) offset(%d)", msg.Key, uid, num, cType, coinId, value.MsgId, msg.Offset)
			s.addCapsule(context.Background(), uid, coinId, num)
		default:
			time.Sleep(time.Second * 3)
			continue
		}
	}
}

// SendGift 送礼增加扭蛋积分
func (s *Service) sendGift(ctx context.Context, uid, giftId, num, totalCoin int64, coinType string, areaV2ParentId, areaV2Id int64, platform string) {
	if totalCoin <= 0 {
		return
	}
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || len(coinConfMap) == 0 {
		return
	}
	var addCoinId = int64(dao.NormalCoinId)
	var coinIds = []int64{dao.BlessCoinId, dao.LplCoinId, dao.WeekCoinId, dao.ColorfulCoinId, dao.NormalCoinId}
	for _, coinId := range coinIds {
		if _, ok := coinConfMap[coinId]; ok {
			if coinConfMap[coinId].AreaMap != nil {
				_, v2ID := coinConfMap[coinId].AreaMap[areaV2Id]
				_, v2ParentID := coinConfMap[coinId].AreaMap[areaV2ParentId]
				if v2ID || v2ParentID {
					if coinConfMap[coinId].GiftType == dao.CapsuleGiftTypeAll {
						addCoinId = coinId
					} else if coinConfMap[coinId].GiftType == dao.CapsuleGiftTypeGold {
						if coinType == "gold" {
							addCoinId = coinId
						}
					} else if coinConfMap[coinId].GiftType == dao.CapsuleGiftTypeSelected {
						if coinConfMap[coinId].GiftMap != nil {
							if _, ok := coinConfMap[coinId].GiftMap[giftId]; ok {
								addCoinId = coinId
							}
						}
					}
				}
			}
		}
		if addCoinId != dao.NormalCoinId {
			break
		}
	}
	// 首次赠送
	if addCoinId == dao.LplCoinId {
		if s.dao.CheckLplFirstGift(ctx, uid, giftId) {
			totalCoin = totalCoin + coinConfMap[addCoinId].ChangeNum
		}
	}
	if addCoinId <= dao.ColorfulCoinId {
		_, err = s.dao.UpdateScore(ctx, uid, addCoinId, totalCoin, "sendGift", platform, nil, coinConfMap[addCoinId])
	} else {
		_, err = s.dao.UpdateCapsule(ctx, uid, addCoinId, totalCoin, "sendGift", platform, coinConfMap[addCoinId])
	}
	if err != nil {
		log.Error("[service.lottery|sendGift] UpdateScore type:%d error:%v", addCoinId, err)
		return
	}
}

func (s *Service) addCapsule(ctx context.Context, uid, coinId, num int64) {
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || len(coinConfMap) == 0 {
		return
	}
	addCoinId := coinId
	if _, ok := coinConfMap[addCoinId]; !ok {
		return
	}
	totalCoin := coinConfMap[addCoinId].ChangeNum * num
	if addCoinId <= dao.ColorfulCoinId {
		_, err = s.dao.UpdateScore(ctx, uid, addCoinId, totalCoin, "sendGift", "", nil, coinConfMap[addCoinId])
	} else {
		_, err = s.dao.UpdateCapsule(ctx, uid, addCoinId, totalCoin, "sendGift", "", coinConfMap[addCoinId])
	}
	if err != nil {
		log.Error("[service.lottery|addCapsule] UpdateScore type:%d error:%v", addCoinId, err)
		return
	}
}

//定时重置Capusule
func (s *Service) tickerReloadCapsuleConf(ctx context.Context) {
	changeFlag, _ := s.dao.GetCapsuleChangeFlag(ctx)
	s.dao.RelaodCapsuleConfig(ctx, changeFlag)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			redisChangeFlag, _ := s.dao.GetCapsuleChangeFlag(ctx)
			capsuleCacheTime, capsuleChangeFlag := s.dao.GetCapsuleChangeInfo(ctx)
			if redisChangeFlag != capsuleChangeFlag || time.Now().Unix()-capsuleCacheTime > 60 {
				s.dao.RelaodCapsuleConfig(ctx, redisChangeFlag)
			}
		}
	}()
}
