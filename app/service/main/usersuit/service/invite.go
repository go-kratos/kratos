package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"

	coinmdl "go-common/app/service/main/coin/model"
	memmdl "go-common/app/service/main/member/model"
	blkmdl "go-common/app/service/main/member/model/block"
	"go-common/app/service/main/usersuit/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

var (
	_buyFlag      = "1"
	_applyFlag    = "1"
	_emptyInvites = make([]*model.Invite, 0)
)

const (
	_inviteCodeExpireSeconds = int64(86400 * 3)
	_checkZero               = int8(1) // 硬币可为负数
	_level5BuyLimit          = 1
	_level6BuyLimit          = 2
)

// BuyInvite buy invite
func (s *Service) BuyInvite(c context.Context, mid int64, num int64, ip string) (res []*model.Invite, err error) {
	var (
		levelInfo *memmdl.LevelInfo
		moralInfo *memmdl.Moral
		blockInfo *blkmdl.RPCResInfo
		coins     float64
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		if levelInfo, err = s.memRPC.Level(errCtx, &memmdl.ArgMid2{Mid: mid, RealIP: ip}); err != nil {
			log.Error("s.memRPC.Level(%+v) error(%v)", &memmdl.ArgMid2{Mid: mid, RealIP: ip}, err)
		}
		return
	})
	eg.Go(func() (err error) {
		if moralInfo, err = s.memRPC.Moral(errCtx, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil {
			log.Error("s.memRPC.Moral(%+v) error(%v)", &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}, err)
		}
		return
	})
	eg.Go(func() (err error) {
		if blockInfo, err = s.memRPC.BlockInfo(errCtx, &blkmdl.RPCArgInfo{MID: mid}); err != nil {
			log.Error("s.memRPC.BlockInfo(%d) error(%v)", mid, err)
		}
		return
	})
	eg.Go(func() (err error) {
		if coins, err = s.coinRPC.UserCoins(errCtx, &coinmdl.ArgCoinInfo{Mid: mid}); err != nil {
			log.Error("d.coinRP.UserCoins(%d) err(%+v)", mid, err)
		}
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}
	switch {
	case levelInfo == nil:
		err = ecode.UsersuitInviteLevelLow
		return
	case moralInfo == nil:
		err = ecode.LackOfScores
		return
	case blockInfo == nil:
		err = ecode.UserDisabled
		return
	}
	if blockInfo.BlockStatus != 0 {
		err = ecode.UserDisabled
		return
	}
	if moralInfo.Moral/100 < 60 {
		err = ecode.LackOfScores
		return
	}
	cost := float64(num * 50)
	if coins < cost {
		err = ecode.LackOfCoins
		return
	}
	level := levelInfo.Cur
	var limit int64
	var expireSeconds int64
	switch level {
	case 5:
		limit = _level5BuyLimit
		expireSeconds = _inviteCodeExpireSeconds
	case 6:
		limit = _level6BuyLimit
		expireSeconds = _inviteCodeExpireSeconds
	default:
		err = ecode.UsersuitInviteLevelLow
		return
	}
	var ok bool
	if ok, err = s.d.SetBuyFlagCache(c, mid, _buyFlag); err != nil {
		return
	}
	if !ok {
		err = ecode.ServiceUnavailable
		return
	}
	defer func() {
		s.d.DelBuyFlagCache(context.Background(), mid)
	}()
	now := time.Now()
	start, end := rangeMonth(now)
	var count int64
	if count, err = s.d.CurrentCount(c, mid, start, end); err != nil {
		return
	}
	if count+num > limit {
		err = ecode.UsersuitInviteReachCurrentMonthLimit
		return
	}
	arg := &coinmdl.ArgModifyCoin{Mid: mid, Count: -cost, Reason: fmt.Sprintf("购买邀请码%d个", num), CheckZero: _checkZero, IP: ip}
	if _, err = s.coinRPC.ModifyCoin(c, arg); err != nil {
		return
	}
	nowTs := now.Unix()
	codem := make(map[string]int)
	for int64(len(codem)) < num {
		codem[geneInviteCode(mid, nowTs)] = 1
	}
	invs := make([]*model.Invite, 0)
	buyIP := net.ParseIP(ip)
	for code := range codem {
		invs = append(invs, &model.Invite{
			Status:  model.StatusOk,
			Mid:     mid,
			Code:    code,
			IP:      IPv4toN(buyIP),
			IPng:    buyIP,
			Expires: nowTs + expireSeconds,
			Ctime:   xtime.Time(nowTs),
		})
	}
	var tx *xsql.Tx
	if tx, err = s.d.Begin(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			log.Error("发放邀请码失败，用户ID：%d，已扣减硬币数：%.2f，未发放成功数量：%d，用户IP：%s", mid, cost, num, ip)
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
			log.Error("发放邀请码失败，用户ID：%d，已扣减硬币数：%.2f，未发放成功数量：%d，用户IP：%s", mid, cost, num, ip)
		}
	}()
	for _, inv := range invs {
		if _, err = s.d.TxAddInvite(c, tx, inv); err != nil {
			return
		}
	}
	res = invs
	return
}

func geneInviteCode(mid int64, ts int64) string {
	data := md5.Sum([]byte(fmt.Sprintf("%d,%d,%d", ts, mid, rand.Int63())))
	h := hex.EncodeToString(data[:])
	return h[8:24]
}

func rangeMonth(now time.Time) (start, end time.Time) {
	year := now.Year()
	month := now.Month()
	loc := now.Location()
	start = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end = time.Date(year, month+1, 0, 23, 59, 59, 0, loc)
	return
}

// ApplyInvite apply invite
func (s *Service) ApplyInvite(c context.Context, mid int64, code string, cookie, ip string) (err error) {
	var memInfo *memmdl.BaseInfo
	if memInfo, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil {
		log.Error("s.memRPC.Level(%+v) error(%v)", &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}, err)
		return
	}
	if memInfo == nil || memInfo.Rank >= 10000 {
		err = ecode.UsersuitInviteAlreadyFormal
		return
	}
	var ok bool
	if ok, err = s.d.SetApplyFlagCache(c, code, _applyFlag); err != nil {
		return
	}
	if !ok {
		err = ecode.ServiceUnavailable
		return
	}
	defer func() {
		s.d.DelApplyFlagCache(context.Background(), code)
	}()
	var inv *model.Invite
	if inv, err = s.d.Invite(c, code); err != nil {
		return
	}
	if inv == nil {
		err = ecode.UsersuitInviteCodeNotExists
		return
	}
	if inv.Used() {
		err = ecode.UsersuitInviteCodeUsed
		return
	}
	nowTs := time.Now().Unix()
	if inv.Expired(nowTs) {
		err = ecode.UsersuitInviteCodeExpired
		return
	}
	if _, err = s.d.UpdateInvite(c, mid, nowTs, code); err != nil {
		return
	}
	if err = s.beFormal(c, mid, cookie, ip); err != nil {
		log.Error("service.beFormal(%s, %d)  error(%v)", code, mid, err)
		log.Error("使用邀请码转正失败，用户ID：%d，用户IP：%s，邀请码：%+v", mid, ip, inv)
	}
	return
}

func (s *Service) beFormal(c context.Context, mid int64, cookie, ip string) (err error) {
	for i := 0; i < 3; i++ {
		if err = s.d.BeFormal(c, mid, cookie, ip); err == nil {
			return
		}
	}
	return
}

// Stat stat
func (s *Service) Stat(c context.Context, mid int64, ip string) (res *model.InviteStat, err error) {
	var (
		level     int32
		levelInfo *memmdl.LevelInfo
	)
	if levelInfo, err = s.memRPC.Level(c, &memmdl.ArgMid2{Mid: mid, RealIP: ip}); err != nil {
		log.Error("s.memRPC.Level(%+v) error(%v)", &memmdl.ArgMid2{Mid: mid, RealIP: ip}, err)
		return
	}
	if levelInfo == nil {
		level = 0
	} else {
		level = levelInfo.Cur
	}
	var limit int64
	switch level {
	case 5:
		limit = _level5BuyLimit
	case 6:
		limit = _level6BuyLimit
	default:
		res = &model.InviteStat{
			Mid:           mid,
			CurrentLimit:  0,
			CurrentBought: 0,
			TotalUsed:     0,
			TotalBought:   0,
			InviteCodes:   _emptyInvites,
		}
		return
	}
	var invs []*model.Invite
	if invs, err = s.d.Invites(c, mid); err != nil {
		return
	}
	now := time.Now()
	nowTs := now.Unix()
	sort.Sort(model.SortInvitesByCtimeDesc(invs))
	start, end := rangeMonth(now)
	startTs, endTs := start.Unix(), end.Unix()
	curBought := int64(0)
	totalUsed := int64(0)
	for _, inv := range invs {
		inv.FillStatus(nowTs)
		buyTs := int64(inv.Ctime)
		if buyTs >= startTs && buyTs <= endTs {
			curBought++
		}
		if inv.Status == model.StatusUsed {
			totalUsed++
		}
	}
	res = &model.InviteStat{
		Mid:           mid,
		CurrentLimit:  limit,
		CurrentBought: curBought,
		TotalBought:   int64(len(invs)),
		TotalUsed:     totalUsed,
		InviteCodes:   invs,
	}
	return
}

// IPv4toN is
func IPv4toN(ip net.IP) (sum uint32) {
	v4 := ip.To4()
	if v4 == nil {
		return
	}
	sum += uint32(v4[0]) << 24
	sum += uint32(v4[1]) << 16
	sum += uint32(v4[2]) << 8
	sum += uint32(v4[3])
	return sum
}
