package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	coinmdl "go-common/app/service/main/coin/model"
	pointmdl "go-common/app/service/main/point/model"
	"go-common/app/service/main/usersuit/model"
	vipmdl "go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_gtbig   int64 = 31
	_entryID int64 = 30
	_dayTime int64 = 86400
	_expires int64 = 31 * 86400
	_coin    int64 = 0
	_point   int64 = 2
	_bcoin   int64 = 1
	_subject       = "头像挂件"
	_reason        = "购买挂件"
	_ctype         = "硬币"
	_ptype         = "积分"
	_btype         = "B币"
)

var (
	_emptyPendant = &model.PendantEquip{Pid: 0}
)

type orderParam struct {
	Mid, Pid, Expires, Price, Tp int64
	Cost, IP, OrderID            string
	Pendant                      *model.Pendant
}

// GroupInfo get group info for web.
func (s *Service) GroupInfo(c context.Context) (res []*model.PendantGroupInfo, err error) {
	return s.groupInfo, nil
}

// AllGroupInfo get all group information
func (s *Service) AllGroupInfo(c context.Context) (res []*model.PendantGroupInfo, err error) {
	var (
		pendant []*model.Pendant
		temp    []*model.PendantGroupInfo
	)
	if temp, err = s.pendantDao.PendantGroupInfo(c); err != nil || temp == nil {
		err = ecode.PendantNotFound
		return
	}
	res = make([]*model.PendantGroupInfo, 0)
	for i, x := range temp {
		if pendant, err = s.pendantInfoByGid(c, x.ID); err != nil {
			return
		}
		if len(pendant) != 0 {
			pgi := temp[i]
			pgi.SubPendant = pendant
			pgi.Number = int64(len(pendant))
			res = append(res, pgi)
		}
	}
	return
}

// GroupInfoByID get group info by gid
func (s *Service) GroupInfoByID(c context.Context, gid int64) (res *model.PendantGroupInfo, err error) {
	var ok bool
	if res, ok = s.groupMap[gid]; ok {
		return
	}
	return
}

// PendantAll get all pendant info
func (s *Service) PendantAll(c context.Context) (res []*model.Pendant, err error) {
	if res, err = s.pendantDao.PendantList(c); err != nil || res == nil {
		err = ecode.PendantNotFound
		return
	}
	for _, v := range res {
		if err = s.assemblePendant(c, v); err != nil {
			return
		}
	}
	return
}

// PendantInfo get pendant info by id
func (s *Service) PendantInfo(c context.Context, pid int64) (res *model.Pendant, err error) {
	var ok bool
	if res, ok = s.pendantMap[pid]; !ok {
		err = ecode.PendantNotFound
	}
	return
}

// PendantPoint return pendant which has point pay type
func (s *Service) PendantPoint(c context.Context, mid int64) (res []*model.Pendant, err error) {
	var (
		pkg []*model.PendantPackage
		tmp = make(map[int64]int64)
	)
	if pkg, err = s.PackageInfo(c, mid); err != nil {
		return
	}
	for _, v := range pkg {
		tmp[v.Pid] = v.Expires
	}
	all := s.pendantInfo
	res = make([]*model.Pendant, 0)
	for _, v := range all {
		if v.Status == model.PendantStatusOFF {
			continue
		}
		if _, ok := tmp[v.ID]; ok {
			v.Expires = tmp[v.ID]
		}
		if v.Point > 0 && v.Gid != _gtbig {
			res = append(res, v)
		}
	}
	l := len(res)
	for i := 0; i < l; i++ {
		now := res[i]
		var nowGRank = &model.PendantGroupInfo{}
		if nowGRank, err = s.GroupInfoByID(c, now.Gid); err != nil || nowGRank == nil {
			log.Error("PendantPoint nowGRank(%+v) err(%+v)", nowGRank, err)
			return
		}
		for j := i + 1; j < l-1; j++ {
			next := res[j]
			if now.Rank == next.Rank {
				var nextGRank = &model.PendantGroupInfo{}
				if nextGRank, err = s.GroupInfoByID(c, next.Gid); err != nil || nextGRank == nil {
					log.Error("PendantPoint nextGRank(%+v) err(%+v)", nextGRank, err)
					return
				}
				if nowGRank.Rank > nextGRank.Rank {
					tp := res[i]
					res[i] = res[j]
					res[j] = tp
				}
			}
		}
	}
	return
}

// OrderHistory get order info.
func (s *Service) OrderHistory(c context.Context, arg *model.ArgOrderHistory) (res []*model.PendantOrderInfo, count map[string]int64, err error) {
	var (
		p  *model.Pendant
		ct int64
	)
	if res, ct, err = s.pendantDao.OrderInfo(c, arg); err != nil {
		return
	}
	for i, r := range res {
		res[i].TimeLength = r.TimeLength / (_dayTime)
		if p, err = s.PendantInfo(c, res[i].Pid); err != nil {
			return
		}
		res[i].Name = p.Name
		res[i].Image = p.Image
	}
	count = make(map[string]int64, 3)
	count["page_current"] = arg.Page
	count["page_size"] = 20
	count["result_count"] = ct
	return
}

// PackageInfo get package by mid
func (s *Service) PackageInfo(c context.Context, mid int64) (res []*model.PendantPackage, err error) {
	defer func() {
		if err == nil && res != nil && len(res) != 0 {
			if err = s.pendantDao.DelRedPointCache(c, mid); err != nil {
				log.Error("s.pendantDao.DelRedPointCache(%d) error(%+v)", mid, err)
				err = nil
			}
		}
	}()
	var (
		cache = false
		t     = time.Now().Unix()
	)
	if res, err = s.pendantDao.PKGCache(c, mid); err != nil {
		return
	}
	for _, v := range res {
		if v.Expires < t {
			cache = true
			res = nil
			break
		}
	}
	if !cache && res != nil {
		return
	}
	cache = true
	if res, err = s.pendantDao.PackageByMid(c, mid); err != nil {
		return
	}
	if len(res) == 0 {
		s.pendantDao.DelPKGCache(c, mid)
		return
	}
	for k, v := range res {
		if res[k].Pendant, err = s.PendantInfo(c, v.Pid); err != nil {
			return
		}
	}
	if cache && len(res) > 0 {
		s.addCache(func() {
			s.pendantDao.AddPKGCache(context.Background(), mid, res)
		})
	}
	return
}

// Equipment return equipped pendant.
func (s *Service) Equipment(c context.Context, mid int64) (res *model.PendantEquip, err error) {
	defer func() {
		if err == nil && res != nil && res.Pid != 0 {
			if res.Pendant, err = s.PendantInfo(c, res.Pid); err != nil {
				log.Error("can not found user(%d) equip(%d) error(%+v)", mid, res.Pid, err)
				err = nil
			}
		}
	}()
	var (
		cache = true
		noRow = false
		t     = time.Now().Unix()
	)
	res = &model.PendantEquip{}
	if res, err = s.pendantDao.EquipCache(c, mid); err != nil {
		log.Error("s.pendantDao.EquipCache error(%v)", err)
		err = nil
		cache = false
	}
	if res != nil {
		if res.Pid == 0 {
			res = nil
			return
		}
		if t < res.Expires {
			return
		}
		log.Info("user(%d) equip expire.", mid)
		s.pendantDao.DelEquipCache(c, mid)
	}
	if res, noRow, err = s.pendantDao.EquipByMid(c, mid, t); err != nil {
		log.Error("s.pendantDao.EquipByMid(%d), err(%+v)", mid, err)
		return
	}
	if noRow {
		res = _emptyPendant
	}
	if cache {
		cres := &model.PendantEquip{
			Mid:     res.Mid,
			Pid:     res.Pid,
			Expires: res.Expires,
			Pendant: nil,
		}
		s.addCache(func() {
			s.pendantDao.AddEquipCache(context.Background(), mid, cres)
		})
	}
	return
}

// Equipments obtain equipemt from cache or db .
func (s *Service) Equipments(c context.Context, mids []int64) (res map[int64]*model.PendantEquip, err error) {
	defer func() {
		// err == nil && res != nil && res.Pid != 0
		if err != nil || res == nil {
			return
		}
		for _, e := range res {
			if e.Pid == 0 {
				continue
			}
			if e.Pendant, err = s.PendantInfo(c, e.Pid); err != nil {
				log.Error("can not found user(%d) equip(%d) error(%+v)", e.Mid, e.Pid, err)
				err = nil
			}
		}
	}()
	var (
		cache     = true
		missedMID []int64
		t         = time.Now().Unix()
	)
	if len(mids) > 100 {
		err = ecode.RequestErr
		return
	}
	if res, missedMID, err = s.pendantDao.EquipsCache(c, mids); err != nil {
		log.Error("s.pendantDao.EquipCache error(%v)", err)
		err = nil
		cache = false
	}
	for _, v := range res {
		if t >= v.Expires {
			delete(res, v.Mid)
		}
	}
	if len(missedMID) == 0 {
		return
	}
	equipMap, err := s.pendantDao.EquipByMids(c, missedMID, t)
	if err != nil {
		return
	}
	if len(equipMap) == 0 {
		return
	}
	var cleanMids []int64
	for _, v := range equipMap {
		if t >= v.Expires {
			cleanMids = append(cleanMids, v.Mid)
			delete(equipMap, v.Mid)
			continue
		}
		cv := &model.PendantEquip{
			Mid:     v.Mid,
			Pid:     v.Pid,
			Expires: v.Expires,
			Pendant: nil,
		}
		res[v.Mid] = cv
	}
	if len(cleanMids) > 0 {
		s.pendantDao.DelEquipsCache(c, cleanMids)
	}
	if cache {
		s.addCache(func() {
			s.pendantDao.AddEquipsCache(context.Background(), equipMap)
		})
	}
	return
}

// PackageByID get package by mid and pid
func (s *Service) PackageByID(c context.Context, mid, pid int64) (res *model.PendantPackage, err error) {
	return s.pendantDao.PackageByID(c, mid, pid)
}

// OrderPendant order pendant
func (s *Service) OrderPendant(c context.Context, mid, pid, expires, tp int64) (res *model.PayInfo, err error) {
	var pendant *model.Pendant

	if pendant, err = s.PendantInfo(c, pid); err != nil || pendant == nil {
		err = ecode.PendantNotFound
		return
	}

	if pendant.Status == model.PendantStatusOFF {
		err = ecode.PendantNotFound
		return
	}
	// if pendant.Gid == _gtbig {
	// 	return s.vipOrder(c, mid, pid, pendant)
	// }
	// common order
	orderInfo := &orderParam{Mid: mid, Pid: pid, Expires: expires, Tp: tp, IP: metadata.String(c, metadata.RemoteIP), Pendant: pendant}
	switch tp {
	case 0:
		return s.coinOrder(c, mid, pid, expires, orderInfo, pendant)
	case 2:
		return s.pointOrder(c, mid, pid, expires, orderInfo, pendant)
	case 1:
		return s.bcoinOrder(c, mid, pid, expires, orderInfo, pendant)
	default:
		err = ecode.RequestErr
		return
	}
}

// BatchGrantPendantByMid batch grant pendant.
func (s *Service) BatchGrantPendantByMid(c context.Context, pid, expire int64, mids []int64) (err error) {
	var pendant *model.Pendant

	if pendant, err = s.PendantInfo(c, pid); err != nil {
		return
	}
	if pendant == nil {
		err = ecode.PendantNotFound
		return
	}
	s.addNotify(func() {
		s.processBatchByMid(context.Background(), pid, expire, mids)
	})
	return
}

// BatchGrantPendantByPid batch grant pendant
func (s *Service) BatchGrantPendantByPid(c context.Context, mid int64, expires, pids []int64) (err error) {
	s.addNotify(func() {
		s.processBatchByPid(context.Background(), mid, expires, pids)
	})
	return
}

// PendantCallback pay call back func
func (s *Service) PendantCallback(c context.Context, arg *model.PendantOrderInfo) (err error) {
	var (
		order, o  *model.PendantOrderInfo
		tx        *sql.Tx
		argPkg, r *model.PendantPackage
	)
	if order, err = s.pendantDao.OrderInfoByID(c, arg.OrderID); err != nil {
		return
	} else if order == nil {
		err = ecode.PendantOrderNotFound
		return
	}

	if tx, err = s.pendantDao.BeginTran(c); err != nil {
		return
	}
	argPkg = &model.PendantPackage{Mid: order.Mid, Pid: order.Pid, Expires: order.TimeLength, Status: 1, Type: 1}
	if r, err = s.PackageByID(c, order.Mid, order.Pid); err != nil {
		return
	}
	if r != nil {
		if err = s.txUpdatePackagInfo(c, argPkg, r, tx); err != nil {
			tx.Rollback()
			log.Error("s.TxUpdatePackagInfo error %v", err)
			return
		}
	} else {
		t := time.Now().Unix()
		argPkg.Expires += t
		if _, err = s.pendantDao.TxAddPackage(c, argPkg, tx); err != nil {
			tx.Rollback()
			return
		}
	}
	o = &model.PendantOrderInfo{OrderID: order.OrderID, Stauts: 1, PayID: order.PayID, IsCallback: 1, CallbackTime: time.Now().Unix()}
	if _, err = s.pendantDao.TxUpdateOrderInfo(c, o, tx); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.pendantDao.DelPKGCache(c, order.Mid)
	s.addNotify(func() {
		if err = s.pendantDao.SetRedPointCache(context.Background(), order.Mid, order.Pid); err != nil {
			log.Error("s.pendantDao.SetRedPointCache(%d,%d) error(%+v)", order.Mid, order.Pid, err)
			err = nil
		}
	})

	return
}

// EquipPendant equip pendant
func (s *Service) EquipPendant(c context.Context, mid, pid int64, status int8, source int64) error {
	switch status {
	case model.PendantEquipOFF:
		if err := s.TakeOffPendant(c, mid); err != nil {
			return err
		}
	case model.PendantEquipON:
		if err := s.WearPendant(c, mid, pid, source); err != nil {
			return err
		}
	default:
		log.Warn("mid(%d) pid(%d) not exist status(%d)", mid, pid, status)
		return ecode.RequestErr
	}

	// 操作成功，通知
	s.addNotify(func() {
		s.accNotify(context.Background(), mid, model.AccountNotifyUpdatePendant)
	})
	return nil
}

// TakeOffPendant 卸下挂件,并删除缓存
func (s *Service) TakeOffPendant(c context.Context, mid int64) error {
	if _, err := s.pendantDao.UpEquipMID(c, mid); err != nil {
		log.Error("s.pendantDao.UpEquipMID(%d) error(%+v)", mid, err)
		return err
	}
	if cacheErr := s.pendantDao.DelEquipCache(c, mid); cacheErr != nil {
		log.Error("s.pendantDao.DelEquipCache(%d) error(%+v)", mid, cacheErr)
	}
	return nil
}

// WearPendant 佩戴挂件
func (s *Service) WearPendant(c context.Context, mid, pid, source int64) error {
	pd, err := s.PendantInfo(c, pid)
	if err != nil || pd == nil {
		return ecode.PendantNotFound
	}
	arg := &model.PendantEquip{
		Mid: mid,
		Pid: pid,
	}

	log.Info("wearPendant req params: mid=%v,pid=%v,source=%v", mid, pid, source)
	source = func(source int64) int64 {
		// 未按照规则传的挂件来源，设置为未知来源
		if !model.IsValidSource(source) {
			return model.UnknownEquipSource
		}
		// 不知道挂件来源,根据挂件属性组进行判断
		if source == model.UnknownEquipSource {
			source = model.EquipFromPackage
			if pd.Gid == _gtbig {
				return model.EquipFromVIP
			}
		}
		return source
	}(source)

	// 背包挂件：在背包里面则直接佩戴
	equipPkgPendant := func() error {
		// 从背包里取这个挂件出来
		var pkg *model.PendantPackage
		pkg, err = s.PackageByID(c, mid, pid)
		if err != nil {
			return err
		}
		// 挂件不在背包里面
		if pkg == nil {
			return ecode.PendantPackageNotFound
		}
		// 挂件在背包里，则直接佩戴
		arg.Expires = pkg.Expires
		_, err = s.pendantDao.AddEquip(c, arg)
		return err
	}

	// 大会员挂件：根据大会员是否过期判断是否佩戴
	equipVipPendant := func() error {
		// 判断是否是大会员挂件
		if pd.Gid != _gtbig {
			return ecode.PendantPackageNotFound
		}

		// 视用户大会员情况而定
		var vi *model.VipInfo
		vi, err = s.pendantDao.VipInfo(c, mid, metadata.String(c, metadata.RemoteIP))
		if err != nil {
			return err
		}
		if vi == nil {
			return ecode.PendantGetVIPErr
		}
		if vi.VipStatus != vipmdl.VipStatusNotOverTime {
			return ecode.VipUserInfoNotExit
		}
		arg.Expires = vi.VipDueDate / 1000
		_, err = s.pendantDao.AddEquip(c, arg)
		return err
	}

	// 佩戴过程
	switch source {
	case model.EquipFromPackage:
		if err := equipPkgPendant(); err != nil {
			log.Error("equipPkgPendant mid=%v,pid=%v,source=%v, err(%+v)", mid, pid, source, err)
			return err
		}
	case model.EquipFromVIP:
		if err := equipVipPendant(); err != nil {
			log.Error("equipVipPendant mid=%v,pid=%v,source=%v, err(%+v)", mid, pid, source, err)
			return err
		}
	default:
		log.Warn("mid=%v,pid=%v,source=%v", mid, pid, source)
		return ecode.RequestErr
	}

	// 佩戴成功删缓存
	if delCacheErr := s.pendantDao.DelEquipCache(c, mid); delCacheErr != nil {
		log.Error("s.pendantDao.DelEquipCache(%d) error(%+v)", mid, delCacheErr)
	}
	if addCacheErr := s.pendantDao.AddEquipCache(c, mid, arg); addCacheErr != nil {
		log.Error("s.pendantDao.AddEquipCache(%d) error(%+v)", mid, addCacheErr)
	}
	return nil
}

func (s *Service) assemblePendant(c context.Context, pendant *model.Pendant) (err error) {
	var (
		ok    bool
		gid   int64
		price map[int64]*model.PendantPrice
	)
	if price, err = s.pendantDao.PendantPrice(c, pendant.ID); err != nil {
		return
	}
	for _, v := range price {
		if v.Type == _coin {
			pendant.Coin = v.Price
		}
		if v.Type == _point {
			pendant.Point = v.Price
		}
		if v.Type == _bcoin {
			pendant.BCoin = v.Price
		}
	}
	if gid, ok = s.pidMap[pendant.ID]; !ok {
		log.Warn("s.pidMap pid(%d)", pendant.ID)
		return
	}
	pendant.Gid = gid
	return
}

func (s *Service) pendantInfoByGid(c context.Context, gid int64) (info []*model.Pendant, err error) {
	var (
		ok   bool
		ids  []int64
		temp []*model.Pendant
	)
	if ids, ok = s.gidMap[gid]; !ok || len(ids) == 0 {
		log.Warn("s.gidMap gid(%d)", gid)
		return
	}
	if temp, err = s.pendantDao.Pendants(c, ids); err != nil {
		return
	}
	for i := range temp {
		rp := temp[i]
		if err = s.assemblePendant(c, rp); err != nil {
			return
		}
		if rp.Gid == _gtbig || rp.Gid == _entryID {
			info = append(info, rp)
		} else if rp.BCoin != 0 || rp.Point != 0 || rp.Coin != 0 {
			info = append(info, rp)
		}
	}
	return
}

func (s *Service) grantOrderID() string {
	var b bytes.Buffer
	b.WriteString(time.Now().Format("20060102"))
	b.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	b.WriteString(strconv.FormatInt(rand.Int63n(999999), 10))
	return b.String()
}

func (s *Service) processOrder(c context.Context, p *orderParam) (err error) {
	var (
		order  *model.PendantOrderInfo
		tx     *sql.Tx
		arg, r *model.PendantPackage
	)
	if r, err = s.PackageByID(c, p.Mid, p.Pid); err != nil {
		return
	}
	if tx, err = s.pendantDao.BeginTran(c); err != nil {
		return
	}
	arg = &model.PendantPackage{Mid: p.Mid, Pid: p.Pid, Expires: p.Expires, Status: 1, Type: p.Tp}
	if r != nil {
		if err = s.txUpdatePackagInfo(c, arg, r, tx); err != nil {
			tx.Rollback()
			return
		}
	} else {
		t := time.Now().Unix()
		arg.Expires += t
		if _, err = s.pendantDao.TxAddPackage(c, arg, tx); err != nil {
			tx.Rollback()
			return
		}
	}
	order = &model.PendantOrderInfo{Mid: p.Mid, OrderID: p.OrderID, PayType: p.Tp, PayPrice: float64(p.Price), Stauts: 1, Pid: p.Pid, TimeLength: p.Expires, Cost: p.Cost, BuyTime: time.Now().Unix()}
	if _, err = s.pendantDao.TxAddOrderInfo(c, order, tx); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("OrderPendant tx.Commit(), error(%v)", err)
		return
	}
	s.pendantDao.DelPKGCache(c, p.Mid)
	s.addNotify(func() {
		if err = s.pendantDao.SetRedPointCache(context.Background(), p.Mid, p.Pid); err != nil {
			log.Error("s.pendantDao.SetRedPointCache(%d,%d) error(%+v)", p.Mid, p.Pid, err)
			err = nil
		}
	})
	return
}

// processBatch a routine process batch
func (s *Service) processBatchByPid(c context.Context, mid int64, expires, pids []int64) {
	var (
		pendant *model.Pendant
		err     error
	)
	for i, v := range pids {
		if expires[i] <= 0 || expires[i] > 3650 {
			log.Error("mid:%v pid:%v expire:%v error %v", mid, v, expires[i], err)
			continue
		}
		if pendant, err = s.PendantInfo(c, v); err != nil || pendant == nil {
			continue
		}
		s.grantPendant(c, v, expires[i], mid)
	}
}

func (s *Service) processBatchByMid(c context.Context, pid, expire int64, mids []int64) {
	for _, mid := range mids {
		s.grantPendant(c, pid, expire, mid)
	}
}

func (s *Service) grantPendant(c context.Context, pid, expire, mid int64) {
	var (
		tx      *sql.Tx
		p       *model.Pendant
		arg, r  *model.PendantPackage
		history *model.PendantHistory
		err     error
	)
	if r, err = s.PackageByID(c, mid, pid); err != nil {
		return
	}
	if p, err = s.PendantInfo(c, pid); err != nil {
		return
	}
	if tx, err = s.pendantDao.BeginTran(c); err != nil {
		return
	}
	expire *= _dayTime
	arg = &model.PendantPackage{Mid: mid, Pid: pid, Expires: expire, Status: 1, Type: 5}
	if p.Gid == _gtbig {
		arg.IsVIP = 1
	}
	if r != nil {
		if err = s.txUpdatePackagInfo(c, arg, r, tx); err != nil {
			tx.Rollback()
			return
		}
	} else {
		t := time.Now().Unix()
		arg.Expires += t
		if _, err = s.pendantDao.TxAddPackage(c, arg, tx); err != nil {
			tx.Rollback()
			return
		}
	}
	history = &model.PendantHistory{Mid: mid, Pid: pid, SourceType: 5, Expire: expire, OperatorName: "system", OperatorAction: 0}
	if _, err = s.pendantDao.TxAddHistory(c, history, tx); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.pendantDao.DelPKGCache(c, mid)
	s.addNotify(func() {
		if err = s.pendantDao.SetRedPointCache(context.Background(), mid, pid); err != nil {
			log.Error("s.pendantDao.SetRedPointCache(%d,%d) error(%+v)", mid, pid, err)
		}
	})
}

func (s *Service) txUpdatePackagInfo(c context.Context, arg *model.PendantPackage, item *model.PendantPackage, tx *sql.Tx) (err error) {
	var (
		n int64
		t = time.Now().Unix()
	)
	if item.Expires < t {
		arg.Expires += t
	} else {
		arg.Expires += item.Expires
	}
	if n, err = s.pendantDao.TxUpdatePackageInfo(c, arg, tx); err != nil {
		return
	} else if n <= 0 {
		err = ecode.PendantPackageNotFound
		return
	}
	return
}

// func (s *Service) vipOrder(c context.Context, mid, pid int64, pendantInfo *model.Pendant) (res *model.PayInfo, err error) {
// 	var (
// 		vipInfo *model.VipInfo
// 		r       *model.PendantPackage
// 		pkg     []*model.PendantPackage
// 		tx      *sql.Tx
// 	)
// 	if vipInfo, err = s.pendantDao.VipInfo(c, mid, metadata.String(c, metadata.RemoteIP)); err != nil {
// 		return
// 	}
// 	switch vipInfo.VipType {
// 	case 0:
// 	case 2:
// 		err = ecode.PendantCanNotBuy
// 		return
// 	case 1:
// 		if vipInfo.VipStatus != 1 {
// 			err = ecode.PendantVIPOverdue
// 			return
// 		}
// 		if pkg, err = s.PackageInfo(c, mid); err != nil {
// 			log.Error("s.PackageInfo error(%v)", err)
// 			return
// 		}
// 		for _, x := range pkg {
// 			if (x.Pendant.Gid == _gtbig) && x.Type != 4 && x.Type != 5 {
// 				err = ecode.PendantAlreadyGet
// 				return
// 			}
// 		}
// 		if r, err = s.PackageByID(c, mid, pid); err != nil {
// 			return
// 		}
// 		if tx, err = s.pendantDao.BeginTran(c); err != nil {
// 			return
// 		}
// 		arg := &model.PendantPackage{Mid: mid, Pid: pid, Expires: _expires, Status: 1, Type: 3}
// 		if r != nil { // 续期
// 			if err = s.txUpdatePackagInfo(c, arg, r, tx); err != nil {
// 				tx.Rollback()
// 				return
// 			}
// 		} else {
// 			arg.IsVIP = 1
// 			arg.Expires += time.Now().Unix()
// 			if _, err = s.pendantDao.TxAddPackage(c, arg, tx); err != nil {
// 				tx.Rollback()
// 				return
// 			}
// 		}
// 		OrderID := s.grantOrderID()
// 		order := &model.PendantOrderInfo{Mid: mid, OrderID: OrderID, PayType: 3, Stauts: 1, Pid: pid, TimeLength: _expires, BuyTime: time.Now().Unix()}
// 		if _, err = s.pendantDao.TxAddOrderInfo(c, order, tx); err != nil {
// 			tx.Rollback()
// 			return
// 		}
// 		if err = tx.Commit(); err != nil {
// 			return
// 		}
// 		eq := &model.PendantEquip{
// 			Mid:     mid,
// 			Pid:     pid,
// 			Expires: arg.Expires,
// 			Pendant: pendantInfo,
// 		}
// 		if _, err = s.pendantDao.AddEquip(c, eq); err != nil {
// 			return
// 		}
// 		s.pendantDao.DelPKGCache(c, mid)
// 		s.pendantDao.DelEquipCache(c, mid)
// 		s.pendantDao.AddEquipCache(c, mid, eq)
// 		s.addNotify(func() {
// 			s.accNotify(context.Background(), mid, model.AccountNotifyUpdatePendant)
// 			if err = s.pendantDao.SetRedPointCache(context.Background(), mid, pid); err != nil {
// 				log.Error("s.pendantDao.SetRedPointCache(%d,%d) error(%+v)", mid, pid, err)
// 				err = nil
// 			}
// 		})
// 		return
// 	default:
// 		break
// 	}
// 	return
// }

func (s *Service) bcoinOrder(c context.Context, mid, pid, expires int64, orderInfo *orderParam, pendant *model.Pendant) (res *model.PayInfo, err error) {
	var b bytes.Buffer
	if pendant.BCoin <= 0 {
		err = ecode.PendantPayTypeErr
		return
	}
	var payID, cashURL string
	bprice := float64(pendant.BCoin) / 100 * float64(expires)
	orderInfo.Price = pendant.BCoin * expires
	orderInfo.Expires = expires * _expires
	b.WriteString(strconv.FormatFloat(bprice, 'f', 2, 64))
	b.WriteString(_btype)
	orderInfo.Cost = b.String()
	orderInfo.OrderID = s.grantOrderID()
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("out_trade_no", orderInfo.OrderID)
	params.Set("money", strconv.FormatFloat(bprice, 'f', 2, 64))
	params.Set("subject", _subject)
	params.Set("remark", fmt.Sprintf(_subject+" - %s（%s个月）", strconv.FormatInt(pid, 10), strconv.FormatInt(expires, 10)))
	params.Set("merchant_id", s.merchantID)
	params.Set("merchant_product_id", s.merchantProductID)
	params.Set("platform_type", "3")
	params.Set("iap_pay_type", "0")
	params.Set("notify_url", s.callBackURL)
	if payID, cashURL, err = s.pendantDao.PayBcoin(c, params, orderInfo.IP); err != nil {
		return
	}
	order := &model.PendantOrderInfo{Mid: orderInfo.Mid, OrderID: orderInfo.OrderID, PayID: payID, PayType: orderInfo.Tp, PayPrice: float64(orderInfo.Price), Stauts: 2, Pid: orderInfo.Pid, TimeLength: orderInfo.Expires, Cost: orderInfo.Cost, BuyTime: time.Now().Unix()}
	if _, err = s.pendantDao.AddOrderInfo(c, order); err != nil {
		return
	}
	res = &model.PayInfo{OrderID: orderInfo.OrderID, OrderNum: payID, PayURL: cashURL}
	return
}

func (s *Service) pointOrder(c context.Context, mid, pid, expires int64, orderInfo *orderParam, pendant *model.Pendant) (res *model.PayInfo, err error) {
	var (
		status int8
		b      bytes.Buffer
	)
	if pendant.Point <= 0 {
		err = ecode.PendantPayTypeErr
		return
	}
	orderInfo.Price = pendant.Point * expires
	orderInfo.Expires = expires * _expires
	b.WriteString(strconv.FormatInt(orderInfo.Price, 10))
	b.WriteString(_ptype)
	orderInfo.Cost = b.String()
	orderInfo.OrderID = s.grantOrderID()
	arg := &pointmdl.ArgPointConsume{Mid: mid, ChangeType: 6, RelationID: strconv.FormatInt(pid, 10), Point: orderInfo.Price, Remark: fmt.Sprintf("兑换挂件(%s)", strconv.FormatInt(pid, 10))}
	if status, err = s.pointRPC.ConsumePoint(c, arg); err != nil {
		return
	}
	if status != pointmdl.SUCCESS {
		log.Warn("mid(%d) not consume point", mid)
		return
	}
	err = s.processOrder(c, orderInfo)
	return
}

func (s *Service) coinOrder(c context.Context, mid, pid, expires int64, orderInfo *orderParam, pendant *model.Pendant) (res *model.PayInfo, err error) {
	var b bytes.Buffer
	if pendant.Coin <= 0 {
		err = ecode.PendantPayTypeErr
		return
	}
	orderInfo.Price = pendant.Coin * expires
	orderInfo.Expires = expires * _expires
	b.WriteString(strconv.FormatInt(orderInfo.Price, 10))
	b.WriteString(_ctype)
	orderInfo.Cost = b.String()
	orderInfo.OrderID = s.grantOrderID()
	arg := &coinmdl.ArgModifyCoin{Mid: mid, Count: -float64(orderInfo.Price), Reason: _reason, CheckZero: 1}
	if s.coinRPC.ModifyCoin(c, arg); err != nil {
		return
	}
	err = s.processOrder(c, orderInfo)
	return
}
