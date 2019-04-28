package v1

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strconv"
	"time"

	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/dao/vip"
	"go-common/app/service/live/xuser/model"
	"go-common/library/cache"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/go-sql-driver/mysql"
	pkgerr "github.com/pkg/errors"
)

const (
	// mysql-driver duplicated entry error code
	_mySQLErrCodeDuplicateEntry = 1062
)

var (
	// get vip info error
	errGetVipInfo = errors.New("get vip info error")
	// buy vip params error
	errBuyParamsInvalid = errors.New("buy vip params invalid")
	// buy vip order_id already exists
	errOrderIDExists = errors.New("orderID already exists")
	// buy vip check order_id server error
	errLockOrderID = errors.New("check orderID failed")
	// buy vip server error
	errBuyVipFailed = errors.New("buy vip failed")
)

// VipService vip service
type VipService struct {
	c                *conf.Config
	dao              *vip.Dao
	liveVipChangePub *databus.Databus
	infoRunCache     *cache.Cache
	addRunCache      *cache.Cache
}

// New new vip service
func New(c *conf.Config) *VipService {
	return &VipService{
		c:                c,
		dao:              vip.New(c),
		liveVipChangePub: databus.New(c.LiveVipChangePub),
		infoRunCache:     cache.New(1, 10240),
		addRunCache:      cache.New(1, 1024),
	}
}

// Info get vip info by uid , grpc wrapper
func (s *VipService) Info(ctx context.Context, req *v1pb.UidReq) (reply *v1pb.InfoReply, err error) {
	vipInfo, err := s.VipInfo(ctx, req.Uid)
	if err != nil {
		log.Error("[service.v1.vip|Info] VipInfo error(%v), uid(%d)", err, req.Uid)
		err = errGetVipInfo
		return
	}
	// 处理vipInfo为nil的情况
	if vipInfo == nil {
		log.Error("[service.v1.vip|Info] VipInfo nil, uid(%d)", req.Uid)
		vipInfo = &model.VipInfo{
			VipTime:  model.TimeEmpty,
			SvipTime: model.TimeEmpty,
		}
	}
	reply = &v1pb.InfoReply{
		Info: &v1pb.Info{
			Vip:      vipInfo.Vip,
			VipTime:  vipInfo.VipTime,
			Svip:     vipInfo.Svip,
			SvipTime: vipInfo.SvipTime,
		},
	}
	return
}

// VipInfo vip info service
func (s *VipService) VipInfo(ctx context.Context, uid int64) (info *model.VipInfo, err error) {
	// cache first
	info, err = s.dao.GetVipFromCache(ctx, uid)
	if err != nil {
		// if cache error return, don't request to db
		log.Error("[service.v1.vip|VipInfo] VipInfo get from cache error(%v), uid(%d)", err, uid)
		return
	}
	if info != nil {
		return
	}

	// db then
	info, err = s.dao.GetVipFromDB(ctx, uid)
	if err != nil {
		// no row in db, set empty cache
		if err == sql.ErrNoRows {
			info.VipTime = model.TimeEmpty
			info.SvipTime = model.TimeEmpty
			log.Info("[service.v1.vip|VipInfo] uid(%d) no row in db", uid)
			goto AsyncSetCache
		}
		log.Error("[service.v1.vip|VipInfo] VipInfo get from db error(%v), uid(%d)", err, uid)
		return
	}

AsyncSetCache:
	// cache set
	s.asyncSetVipCache(ctx, uid, info)
	return
}

// asyncSetVipCache 异步设置vip缓存，通过library/cache队列控制goroutine个数
func (s *VipService) asyncSetVipCache(ctx context.Context, uid int64, info *model.VipInfo) error {
	c := metadata.WithContext(ctx)
	f := func(c context.Context) {
		if err := s.dao.SetVipCache(c, uid, info); err != nil {
			log.Error("[service.v1.vip|asyncSetVipCache] async set vip cache error(%v), uid(%d), info(%v)",
				err, uid, info)
		}
	}
	if runErr := s.infoRunCache.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error("[service.v1.vip|asyncSetVipCache] run cache is full, do it sync.uid(%d), info(%v)", uid, info)
		f(c)
	}
	return nil
}

// Buy buy vip, grpc wrapper
func (s *VipService) Buy(ctx context.Context, req *v1pb.BuyReq) (reply *v1pb.BuyReply, err error) {
	var status int
	buyVip := &model.VipBuy{
		Uid:      req.Uid,
		OrderID:  req.OrderId,
		GoodID:   req.GoodId,
		GoodNum:  req.GoodNum,
		Platform: req.Platform,
		Source:   req.Source,
	}
	status, err = s.AddVip(ctx, buyVip)
	if err != nil {
		log.Error("[service.v1.vip|Buy] buy vip error(%+v), params(%+v)", err, buyVip)
	}
	reply = &v1pb.BuyReply{
		Status: status,
	}
	return
}

// AddVip add vip service logic
func (s *VipService) AddVip(ctx context.Context, req *model.VipBuy) (status int, err error) {
	log.Info("[service.v1.vip|BuyVip] buy vip params(%v)", req)
	var (
		recordID int64
		info     *model.VipInfo
		record   *model.VipRecord
	)
	status = model.BuyStatusRetry

	defer func() {
		log.Info("[service.v1.vip|BuyVip] buy vip status(%d), error(%v), params(%v)", status, err, req)
		// order id exists, do nothing
		if err == errOrderIDExists {
			err = nil
			return
		}
		if status == model.BuyStatusSuccess && err == nil {
			s.asyncAfterAddVipSuccess(ctx, recordID, req.Uid, record, info)
		}
	}()

	// check params
	if !checkBuyParams(req) {
		log.Error("[service.v1.vip|BuyVip] check params(%+v) error", req)
		err = errBuyParamsInvalid
		return
	}

	// check order id
	recordID, err = s.lockVipRecord(ctx, req)
	if err != nil {
		if err == errOrderIDExists {
			log.Info("[service.v1.vip|BuyVip] orderID exists buy success, params(%v)", req)
			goto SUC
		}
		log.Error("[service.v1.vip|BuyVip] lockVipRecord error(%v), params(%v), record id(%d)", err, req, recordID)
		err = errLockOrderID
		return
	}

	// get vip info from db
	info, err = s.dao.GetVipFromDB(ctx, req.Uid)
	if err != nil {
		log.Error("[service.v1.vip|BuyVip] VipInfo error(%v), params(%v)", err, req)
		err = errGetVipInfo
		return
	}

	// add vip month or year
	if req.GoodID == model.Vip {
		record, err = s.addMonth(ctx, req, info)
	} else {
		record, err = s.addYear(ctx, req, info)
	}
	if err != nil {
		log.Error("[service.v1.vip|BuyVip] buy error(%v), params(%v)", err, req)
		err = errBuyVipFailed
		return
	}

SUC:
	status = model.BuyStatusSuccess
	return
}

// lockVipRecord lock user vip record by uid order_id, otherwise create one
func (s *VipService) lockVipRecord(ctx context.Context, req *model.VipBuy) (recordID int64, err error) {
	recordID, err = s.dao.CreateVipRecord(ctx, req)
	switch nErr := pkgerr.Cause(err).(type) {
	case *mysql.MySQLError:
		if nErr.Number == _mySQLErrCodeDuplicateEntry {
			err = errOrderIDExists
		}
	}
	return
}

// addMonth add vip time
func (s *VipService) addMonth(ctx context.Context, req *model.VipBuy, info *model.VipInfo) (record *model.VipRecord, err error) {
	var (
		newVipTime, newSvipTime xtime.Time
		oldVipTime, oldSvipTime xtime.Time
		monthDuration           = xtime.Time(req.GoodNum * 30 * 86400)
		currentTime             = xtime.Time(time.Now().Unix())
	)
	// old vip time
	if info.VipTime != model.TimeEmpty {
		t, _ := time.Parse(model.TimeNano, info.VipTime)
		oldVipTime = xtime.Time(t.Unix())
	}
	// new vip time
	if info.Vip == 1 && oldVipTime > currentTime {
		newVipTime = oldVipTime + monthDuration
	} else {
		newVipTime = currentTime + monthDuration
	}

	// 自动转年费
	// old svip time
	if info.SvipTime != model.TimeEmpty {
		t, _ := time.Parse(model.TimeNano, info.SvipTime)
		oldSvipTime = xtime.Time(t.Unix())
	}
	if oldSvipTime < currentTime {
		oldSvipTime = currentTime
	}
	vipType := model.Vip
	yearNum := math.Floor(float64(newVipTime-oldSvipTime) / float64(359*86400))
	yearDuration := xtime.Time(yearNum * 365 * 86400)
	if yearNum > 0 {
		if info.Svip == 1 && oldSvipTime > currentTime {
			newSvipTime = oldSvipTime + yearDuration
		} else {
			newSvipTime = currentTime + yearDuration
		}
		newVipTime = newVipTime + xtime.Time(yearNum*5*86400) // 额外补五天
		vipType = model.Svip
	}

	// add vip,svip time
	_, err = s.dao.AddVip(ctx, req.Uid, newVipTime, newSvipTime)
	if err != nil {
		log.Error("[service.v1.vip|addMonth] dao.UpdateVipTime error(%v), req(%v), newVip(%d), newsVip(%d)",
			err, req, newVipTime, newSvipTime)
		return
	}

	// build record struct
	ot := model.TimeEmpty
	if oldVipTime > 0 {
		ot = oldVipTime.Time().Format(model.TimeNano)
	}
	record = &model.VipRecord{
		Uid:           req.Uid,
		Opcode:        model.OpcodeAdd,
		BuyType:       req.GoodID,
		BuyNum:        req.GoodNum,
		VipType:       vipType,
		BeforeVipTime: ot,
		AfterVipTime:  newVipTime.Time().Format(model.TimeNano),
		Platform:      req.Source,
	}
	return
}

// addYear add svip time
func (s *VipService) addYear(ctx context.Context, req *model.VipBuy, info *model.VipInfo) (record *model.VipRecord, err error) {
	var (
		newVipTime, newSvipTime xtime.Time
		oldVipTime, oldSvipTime xtime.Time
		yearDuration            = xtime.Time(req.GoodNum * 365 * 86400)
		currentTime             = xtime.Time(time.Now().Unix())
	)
	// old vip time
	if info.VipTime != model.TimeEmpty {
		t, _ := time.Parse(model.TimeNano, info.VipTime)
		oldVipTime = xtime.Time(t.Unix())
	}
	// old svip time
	if info.SvipTime != model.TimeEmpty {
		t, _ := time.Parse(model.TimeNano, info.SvipTime)
		oldSvipTime = xtime.Time(t.Unix())
	}
	// new svip time
	if info.Svip == 1 && oldSvipTime > currentTime {
		newSvipTime = oldSvipTime + yearDuration
	} else {
		newSvipTime = currentTime + yearDuration
	}
	// new vip time
	if info.Vip == 1 && oldVipTime > currentTime {
		newVipTime = oldVipTime + yearDuration
	} else {
		newVipTime = currentTime + yearDuration
	}

	// add vip,svip time
	_, err = s.dao.AddVip(ctx, req.Uid, newVipTime, newSvipTime)
	if err != nil {
		log.Error("[service.v1.vip|addYear] dao.UpdateVipTime error(%v), req(%v), newVip(%v), newsVip(%v)",
			err, req, newVipTime, newSvipTime)
		return
	}

	// build record struct
	ot := model.TimeEmpty
	if oldSvipTime > 0 {
		ot = oldSvipTime.Time().Format(model.TimeNano)
	}
	record = &model.VipRecord{
		Uid:           req.Uid,
		Opcode:        model.OpcodeAdd,
		BuyType:       req.GoodID,
		BuyNum:        req.GoodNum,
		VipType:       model.Svip,
		BeforeVipTime: ot,
		AfterVipTime:  newSvipTime.Time().Format(model.TimeNano),
		Platform:      req.Source,
	}
	return
}

// asyncAfterAddVipSuccess 异步处理增加姥爷成功后的后续逻辑
func (s *VipService) asyncAfterAddVipSuccess(ctx context.Context, recordID, uid int64, record *model.VipRecord, info *model.VipInfo) error {
	c := metadata.WithContext(ctx)
	f := func(c context.Context) {
		// 同步清除缓存，失败异步重试一次
		if clearErr := s.dao.ClearCache(ctx, uid); clearErr != nil {
			go s.dao.ClearCache(c, uid)
		}
		// 异步处理购买成功后续逻辑
		go s.dao.UpdateVipRecord(c, recordID, uid, info)
		go s.dao.CreateApVipRecord(c, record)
		go s.notifyVipChange(c, record)
	}

	if runErr := s.addRunCache.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error("[service.v1.vip|asyncAfterAddVipSuccess] run cache is full, do it sync.record id(%d), uid(%d), record(%v), info(%v)",
			record, uid, record, info)
		f(c)
	}

	return nil
}

// notifyVipChange 异步通知vip购买消息
func (s *VipService) notifyVipChange(ctx context.Context, record *model.VipRecord) (err error) {
	if err = s.liveVipChangePub.Send(ctx, strconv.FormatInt(record.Uid, 10), record); err != nil {
		log.Error("[service.v1.vip|notifyVipChange] send error(%v), record(%v)", err, record)
	}
	return
}

// Close close vip service
func (s *VipService) Close() {
	s.dao.Close()
}

func checkBuyParams(req *model.VipBuy) bool {
	if req.Uid <= 0 || !checkGoodID(req.GoodID) || req.GoodNum <= 0 || req.OrderID == "" {
		return false
	}
	return true
}

func checkGoodID(goodID int) bool {
	return goodID == model.Vip || goodID == model.Svip
}
