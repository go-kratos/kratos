package usersuit

import (
	"context"
	"math"

	"go-common/app/interface/main/account/model"
	cmdl "go-common/app/service/main/coin/model"
	memmdl "go-common/app/service/main/member/model"
	usmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Equip user pendant equip.
func (s *Service) Equip(c context.Context, mid, pid int64, status int8, source int64) (err error) {
	return s.usRPC.Equip(c, &usmdl.ArgEquip{Mid: mid, Pid: pid, Status: status, Source: source})
}

// Equipment get pendant current equipment
func (s *Service) Equipment(c context.Context, mid int64) (equipPHP *model.EquipPHP, err error) {
	var equip *usmdl.PendantEquip
	ip := metadata.String(c, metadata.RemoteIP)
	if equip, err = s.usRPC.Equipment(c, &usmdl.ArgEquipment{Mid: mid, IP: ip}); err != nil {
		log.Error("s.usRPC.Equipment(%d) error(%v)", mid, err)
		return
	}
	var coin float64
	if coin, err = s.coinRPC.UserCoins(c, &cmdl.ArgCoinInfo{Mid: mid}); err != nil {
		log.Error("s.coinRPC.UserCoins(%d) error(%v)", mid, err)
		return
	}
	var base *memmdl.BaseInfo
	if base, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: mid, RemoteIP: ip}); err != nil {
		log.Error("s.memRPC.Base(%d) error(%v)", mid, err)
		return
	}
	equipPHP = &model.EquipPHP{
		Coins:   coin,
		FaceURL: base.Face,
	}
	if equip == nil || equip.Pendant == nil {
		log.Info("s.Equipment(%d) usequip(%+v) or usequip.Pendant(%+v) is nil", equip, equip.Pendant)
		return
	}
	equipPHP.Pid = equip.Pid
	equipPHP.Image = model.FormatImgURL(mid, equip.Pendant.Image)
	equipPHP.ImageModel = model.FormatImgURL(mid, equip.Pendant.ImageModel)
	return
}

// Pendant return pendant info.
func (s *Service) Pendant(c context.Context, pid int64) (pendantPHP *model.PendantPHP, err error) {
	var pendant *usmdl.Pendant
	ip := metadata.String(c, metadata.RemoteIP)
	if pendant, err = s.dao.Pendant(c, pid, ip); err != nil {
		log.Error("s.dao.Group(%d) error(%v)", pid, err)
		return
	}
	pendantPHP = &model.PendantPHP{}
	pendantPHP.Name = pendant.Name
	pendantPHP.Pid = pendant.ID
	pendantPHP.Image = model.FormatImgURL(pid, pendant.Image)
	pendantPHP.ImageModel = model.FormatImgURL(pid, pendant.ImageModel)
	return
}

// Group return pendant group info.
func (s *Service) Group(c context.Context, mid int64) (groupPHP []*model.GroupPHP, err error) {
	var groups []*usmdl.PendantGroupInfo
	ip := metadata.String(c, metadata.RemoteIP)
	if groups, err = s.dao.Group(c, ip); err != nil {
		log.Error("s.dao.Group(%d) error(%v)", mid, err)
		return
	}
	for _, g := range groups {
		if g.ID == 30 || g.ID == 31 {
			continue
		}
		for _, p := range g.SubPendant {
			p.BCoin = p.BCoin / 100
			p.Image = model.FormatImgURL(mid, p.Image)
			p.ImageModel = model.FormatImgURL(mid, p.ImageModel)
		}
		gh := &model.GroupPHP{}
		gh.Name = g.Name
		gh.Count = g.Number
		gh.Pendant = g.SubPendant
		groupPHP = append(groupPHP, gh)
	}
	return
}

// GroupEntry return vip pendant.
func (s *Service) GroupEntry(c context.Context, mid int64) (entryPHP []*model.GroupEntryPHP, err error) {
	var group *usmdl.PendantGroupInfo
	ip := metadata.String(c, metadata.RemoteIP)
	if group, err = s.dao.GroupEntry(c, ip); err != nil {
		log.Error("s.dao.GroupEntry(%d) error(%v)", mid, ip)
		return
	}
	if group == nil {
		log.Info("s.dao.GroupEntry(%d) result value is nil", mid)
		return
	}
	for _, p := range group.SubPendant {
		entry := &model.GroupEntryPHP{}
		entry.Pid = p.ID
		entry.Name = p.Name
		entry.Money = p.Point
		entry.Image = model.FormatImgURL(mid, p.Image)
		entry.ImageModel = model.FormatImgURL(mid, p.ImageModel)
		entryPHP = append(entryPHP, entry)
	}
	return
}

// GroupVIP return vip pendant.
func (s *Service) GroupVIP(c context.Context, mid int64) (vipPHP []*model.GroupVipPHP, err error) {
	var group *usmdl.PendantGroupInfo
	ip := metadata.String(c, metadata.RemoteIP)
	if group, err = s.dao.GroupVip(c, ip); err != nil {
		log.Error("s.dao.GroupVip(%d) error(%v)", mid, ip)
		return
	}
	if group == nil {
		log.Info("s.dao.GroupEntry(%d) result value is nil", mid)
		return
	}
	for _, p := range group.SubPendant {
		vip := &model.GroupVipPHP{}
		vip.Pid = p.ID
		vip.Name = p.Name
		vip.Money = 0
		vip.MoneyType = 3
		vip.Expire = 2678400
		vip.Image = model.FormatImgURL(mid, p.Image)
		vip.ImageModel = model.FormatImgURL(mid, p.ImageModel)
		vipPHP = append(vipPHP, vip)
	}
	return
}

// VipGet pc vip install pendant.
func (s *Service) VipGet(c context.Context, mid, pid int64, activated int8) (err error) {
	err = s.Equip(c, mid, pid, int8(activated), usmdl.EquipFromVIP)
	return
}

// CheckOrder check order by oid.
func (s *Service) CheckOrder(c context.Context, mid int64, orderID string) (err error) {
	var hs []*usmdl.PendantOrderInfo
	ip := metadata.String(c, metadata.RemoteIP)
	if hs, _, err = s.dao.OrderHistory(c, mid, 1, 0, orderID, ip); err != nil {
		log.Error("s.dao.OrderHistory(%d) error(%v)", mid, err)
		return
	}
	if len(hs) == 0 {
		err = ecode.PendantOrderNotFound
		log.Info("s.dao.OrderHistory(%d) orderID(%d) error(%v)", mid, orderID, err)
		return
	}
	if hs[0].Stauts != 1 {
		err = ecode.PendantOrderNotFound
		log.Info("s.dao.OrderHistory(%d) orderID(%d) order not complete", mid, orderID, err)
		return
	}
	return
}

// Order pay pandent by coin/bcoin/point.
func (s *Service) Order(c context.Context, mid, pid, timeLength int64, moneyType int8) (res interface{}, err error) {
	var payInfo *usmdl.PayInfo
	ip := metadata.String(c, metadata.RemoteIP)
	if payInfo, err = s.dao.Order(c, mid, pid, timeLength, moneyType, ip); err != nil {
		log.Error("s.dao.Order(%d) error(%v)", mid, err)
		return
	}
	if payInfo != nil {
		payInfo.PayURL = "https://pay.bilibili.com" + payInfo.PayURL
		res = payInfo
		return
	}
	if moneyType == 2 {
		log.Info("s.dao.Order(%d) pid(%d) buy type with point", mid, pid)
		s.Equip(c, mid, pid, 2, usmdl.EquipFromPackage)
	}
	var pkgs []*usmdl.PendantPackage
	if pkgs, err = s.dao.Packages(c, mid, ip); err != nil {
		log.Error("s.dao.Packages(%d) error(%v)", mid, err)
		return
	}
	var pendant *usmdl.PendantPackage
	for _, pkg := range pkgs {
		if pkg.Pid == pid {
			pendant = pkg
		}
	}
	if pendant != nil {
		res = &struct {
			Msg    string `json:"msg"`
			Expire int64  `json:"expire"`
		}{
			Msg:    "您已成功购买" + pendant.Pendant.Name,
			Expire: pendant.Expires,
		}
	}
	return
}

// My get my pandent
func (s *Service) My(c context.Context, mid int64) (my []*model.MyPHP, err error) {
	var equip *usmdl.PendantEquip
	ip := metadata.String(c, metadata.RemoteIP)
	if equip, err = s.usRPC.Equipment(c, &usmdl.ArgEquipment{Mid: mid, IP: ip}); err != nil {
		log.Error("s.usRPC.Equipment(%d) error(%v)", mid, err)
		return
	}
	var pkgs []*usmdl.PendantPackage
	if pkgs, err = s.dao.Packages(c, mid, ip); err != nil {
		log.Error("s.dao.Packages(%d) error(%v)", mid, err)
		return
	}
	for _, pkg := range pkgs {
		m := &model.MyPHP{}
		m.Pid = pkg.Pid
		m.Name = pkg.Pendant.Name
		m.MoneyType = int8(pkg.Type)
		m.Image = model.FormatImgURL(mid, pkg.Pendant.Image)
		m.ImageModel = model.FormatImgURL(mid, pkg.Pendant.ImageModel)
		m.Expire = pkg.Expires
		m.IsOnline = 1
		if equip != nil && equip.Pid == pkg.Pid {
			m.IsActivated = 1
		}
		my = append(my, m)
	}
	return
}

// MyHistory get my pandent buy history.
func (s *Service) MyHistory(c context.Context, mid, page int64) (res map[string]interface{}, err error) {
	var (
		hs    []*usmdl.PendantOrderInfo
		myhs  []*model.MyHistoryPHP
		count map[string]int64
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	if hs, count, err = s.dao.OrderHistory(c, mid, page, 0, "", ip); err != nil {
		log.Error("s.dao.OrderHistory(%d) error(%v)", mid, err)
		return
	}
	if len(hs) == 0 {
		log.Info("s.dao.OrderHistory(%d) result len eq(0)", mid)
		return
	}
	for _, h := range hs {
		my := &model.MyHistoryPHP{}
		my.Pid = h.Pid
		my.Image = model.FormatImgURL(mid, h.Image)
		my.Name = h.Name
		my.BuyTime = h.BuyTime
		my.PayID = h.PayID
		my.Cost = h.Cost
		my.TimeLength = h.TimeLength
		myhs = append(myhs, my)
	}
	res = make(map[string]interface{})
	count["page_count"] = int64(math.Ceil(float64(count["result_count"]) / float64(count["page_size"])))
	res["page"] = count
	res["result"] = myhs
	return
}
