package elec

import (
	"context"
	"time"

	model "go-common/app/interface/main/creative/model/elec"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

const (
	ftime = "2006-01-02"
)

// UserInfo get user elec info.
func (s *Service) UserInfo(c context.Context, mid int64, ip string) (data *model.UserInfo, err error) {
	if data, err = s.elec.UserInfo(c, mid, ip); err != nil {
		log.Error("s.elec.UserInfo(%d) error(%v)", mid, err)
	}
	return
}

// UserUpdate user join or exit elec.
func (s *Service) UserUpdate(c context.Context, mid int64, st int8, ip string) (data *model.UserInfo, err error) {
	if data, err = s.elec.UserUpdate(c, mid, st, ip); err != nil {
		log.Error("s.elec.UserUpdate(%d) error(%v)", mid, err)
	}
	return
}

// ArcUpdate arc open or close elec.
func (s *Service) ArcUpdate(c context.Context, mid, aid int64, st int8, ip string) (err error) {
	if err = s.elec.ArcUpdate(c, mid, aid, st, ip); err != nil {
		log.Error("s.elec.ArcInfo(%d,%d) error(%v)", mid, aid, err)
	}
	return
}

// Notify get up-to-date notice.
func (s *Service) Notify(c context.Context, ip string) (data *model.Notify, err error) {
	if data, err = s.elec.Notify(c, ip); err != nil {
		log.Error("s.elec.Notify error(%v)", err)
	}
	return
}

// Status get elec setting status.
func (s *Service) Status(c context.Context, mid int64, ip string) (data *model.Status, err error) {
	if data, err = s.elec.Status(c, mid, ip); err != nil {
		log.Error("s.elec.Status error(%d, %v)", mid, err)
	}
	return
}

// UpStatus update elec setting status.
func (s *Service) UpStatus(c context.Context, mid int64, spday int, ip string) (err error) {
	if err = s.elec.UpStatus(c, mid, spday, ip); err != nil {
		log.Error("s.elec.UpStatus(%d) error(%v)", mid, err)
	}
	return
}

// RecentRank recent rank.
func (s *Service) RecentRank(c context.Context, mid, size int64, ip string) (data []*model.Rank, err error) {
	if data, err = s.elec.RecentRank(c, mid, size, ip); err != nil {
		log.Error("s.elec.RecentRank error(%d, %v)", mid, err)
	}
	if len(data) == 0 {
		log.Error("s.elec.TotalRank (%d, %v)", mid, data)
		return
	}
	//data, _ = s.CheckIsFriend(c, data, mid, ip)
	return
}

// CurrentRank current rank.
func (s *Service) CurrentRank(c context.Context, mid int64, ip string) (data []*model.Rank, err error) {
	if data, err = s.elec.CurrentRank(c, mid, ip); err != nil {
		log.Error("s.elec.CurrentRank error(%d, %v)", mid, err)
	}
	if len(data) == 0 {
		log.Error("s.elec.TotalRank (%d, %v)", mid, data)
		return
	}
	//data, _ = s.CheckIsFriend(c, data, mid, ip)
	return
}

// TotalRank total rank.
func (s *Service) TotalRank(c context.Context, mid int64, ip string) (data []*model.Rank, err error) {
	if data, err = s.elec.TotalRank(c, mid, ip); err != nil {
		log.Error("s.elec.TotalRank error(%d, %v)", mid, err)
		return
	}
	if len(data) == 0 {
		log.Error("s.elec.TotalRank (%d, %v)", mid, data)
		return
	}
	//data, _ = s.CheckIsFriend(c, data, mid, ip)
	return
}

// DailyBill daily settlement.
func (s *Service) DailyBill(c context.Context, mid int64, pn, ps int, bg, end, ip string) (data *model.BillList, err error) {
	if bg == "" {
		bg = time.Now().Add(-7 * 24 * time.Hour).Format(ftime)
	}
	if end == "" {
		end = time.Now().Format(ftime)
	}
	if data, err = s.elec.DailyBill(c, mid, pn, ps, bg, end, ip); err != nil {
		log.Error("s.elec.DailyBill error(%d, %v)", mid, err)
	}
	return
}

// Balance get battery balance.
func (s *Service) Balance(c context.Context, mid int64, ip string) (data *model.Balance, err error) {
	if data, err = s.elec.Balance(c, mid, ip); err != nil {
		log.Error("s.elec.Balance error(%d, %v)", mid, err)
	}
	return
}

// AppDailyBill daily settlement.
func (s *Service) AppDailyBill(c context.Context, mid int64, pn, ps int, ip string) (cb *model.ChargeBill, err error) {
	var data *model.BillList
	bg := time.Now().Add(-24 * 30 * 12 * time.Hour).Format(ftime)
	end := time.Now().Format(ftime)
	if data, err = s.elec.DailyBill(c, mid, pn, ps, bg, end, ip); err != nil {
		log.Error("s.elec.DailyBill error(%d, %v)", mid, err)
		return
	}
	cb = &model.ChargeBill{}
	if data == nil {
		log.Error("s.elec.DailyBill mid(%d) data(%v)", mid, data)
		return
	}
	bls := make([]*model.Bill, 0, len(data.List))
	for _, v := range data.List {
		bl := &model.Bill{}
		bl.ID = v.ID
		bl.MID = v.MID
		bl.ChannelType = v.ChannelType
		bl.ChannelTyName = v.ChannelTyName
		bl.AddNum = v.AddNum
		bl.ReduceNum = v.ReduceNum
		bl.WalletBalance = v.WalletBalance
		bl.DateVersion = v.DateVersion
		bl.Remark = v.Remark
		bl.MonthBillResp = v.MonthBillResp
		t, _ := time.Parse(ftime, bl.DateVersion)
		bl.Weekday = model.Weekday(t)
		bls = append(bls, bl)
	}
	cb.List = bls
	cb.Pager.Current = data.Pn
	cb.Pager.Size = data.Ps
	cb.Pager.Total = data.TotalCount
	return
}

// RecentElec get recent charge info.
func (s *Service) RecentElec(c context.Context, mid int64, pn, ps int, ip string) (l *model.RecentElecList, err error) {
	if l, err = s.elec.RecentElec(c, mid, pn, ps, ip); err != nil {
		log.Error("s.elec.RecentElec error(%d, %v)", mid, err)
		return
	}
	if l == nil || len(l.List) == 0 {
		return
	}
	var (
		mids, aids []int64
		a          map[int64]*api.Arc
		u          map[int64]*account.Info
	)
	for _, v := range l.List {
		mids = append(mids, v.MID)
		if v.AID > 0 {
			aids = append(aids, v.AID)
		}
	}
	if len(aids) > 0 {
		if a, err = s.arc.Archives(c, aids, ip); err != nil {
			log.Error("s.arc.Archives aids(%v), ip(%s) err(%v)", aids, ip, err)
			return
		}
	}
	if len(mids) > 0 {
		if u, err = s.acc.Infos(c, mids, ip); err != nil {
			log.Error("s.acc.Infos mids(%v), ip(%s) err(%v)", mids, ip, err)
			return
		}
	}
	els := make([]*model.RecentElec, 0, len(l.List))
	for _, v := range l.List {
		el := &model.RecentElec{}
		el.AID = v.AID
		el.MID = v.MID
		el.ElecNum = v.ElecNum
		el.Avatar = v.Avatar
		el.CTime = v.CTime
		if ui, ok := u[el.MID]; ok && ui != nil {
			el.Avatar = ui.Face
			el.Uname = ui.Name
		}
		if el.AID > 0 {
			if av, ok := a[el.AID]; ok && av != nil {
				el.Title = av.Title
			}
		}
		els = append(els, el)
	}
	l.List = els
	return
}

// RemarkList get remark list.
func (s *Service) RemarkList(c context.Context, mid int64, pn, ps int, bg, end, ip string) (rms *model.RemarkList, err error) {
	if rms, err = s.elec.RemarkList(c, mid, pn, ps, bg, end, ip); err != nil {
		log.Error("s.elec.RemarkList error(%d, %v)", mid, err)
	}
	return
}

// RemarkDetail get remark detail.
func (s *Service) RemarkDetail(c context.Context, mid, id int64, ip string) (rm *model.Remark, err error) {
	if rm, err = s.elec.RemarkDetail(c, mid, id, ip); err != nil {
		log.Error("s.elec.RemarkDetail error(%d, %d,%v)", mid, id, err)
	}
	return
}

// Remark reply a msg.
func (s *Service) Remark(c context.Context, mid, id int64, msg, ak, ck, ip string) (status int, err error) {
	if status, err = s.elec.Remark(c, mid, id, msg, ak, ck, ip); err != nil {
		log.Error("s.elec.Remark error(%d, %d, %v)", mid, id, err)
	}
	return
}
