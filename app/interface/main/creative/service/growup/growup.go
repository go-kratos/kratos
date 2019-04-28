package growup

import (
	"context"
	"go-common/app/interface/main/creative/model/growup"
	"go-common/library/log"
	"time"
)

// UpStatus get up status. 前提:必须当前用户线上有可以观看的视频(0/-6) TODO deprecated
func (s *Service) UpStatus(c context.Context, mid int64, ip string) (us *growup.UpStatus, err error) {
	if us, err = s.growup.UpStatus(c, mid, ip); err != nil {
		log.Error("s.growup.UpStatus(%d) error(%v)", mid, err)
		return
	}
	return
}

// UpInfo get up info.
func (s *Service) UpInfo(c context.Context, mid int64, ip string) (ui *growup.UpInfo, err error) {
	if ui, err = s.growup.UpInfo(c, mid, ip); err != nil {
		log.Error("s.growup.UpInfo(%d) error(%v)", mid, err)
	}
	return
}

// Join join growup.
func (s *Service) Join(c context.Context, mid int64, accTy, signTy int, ip string) (err error) {
	if err = s.growup.Join(c, mid, accTy, signTy, ip); err != nil {
		log.Error("s.growup.Join(%d) error(%v)", mid, err)
	}
	return
}

// Quit quit growup.
func (s *Service) Quit(c context.Context, mid int64, ip string) (err error) {
	if err = s.growup.Quit(c, mid, ip); err != nil {
		log.Error("s.growup.Quit(%d) error(%v)", mid, err)
	}
	return
}

// Summary get up status.
func (s *Service) Summary(c context.Context, mid int64, ip string) (us *growup.Summary, err error) {
	if us, err = s.growup.Summary(c, mid, ip); err != nil {
		log.Error("s.growup.Summary(%d) error(%v)", mid, err)
	}
	return
}

// Stat get up info.
func (s *Service) Stat(c context.Context, mid int64, ip string) (ui *growup.Stat, err error) {
	if ui, err = s.growup.Stat(c, mid, ip); err != nil {
		log.Error("s.growup.Stat(%d) error(%v)", mid, err)
	}
	if ui == nil || len(ui.Tops) == 0 {
		log.Error("s.growup.Stat mid(%d) data(%v)", mid, ui)
		return
	}
	dt := time.Now().AddDate(0, 0, -2)
	year := dt.Format("2006")
	month := dt.Format("01")
	ui.Desc = "收入数据统计日期:" + year + "年" + month + "月"
	aids := make([]int64, 0, len(ui.Tops))
	for _, v := range ui.Tops {
		if v != nil {
			aids = append(aids, v.AID)
		}
	}
	avm, err := s.p.BatchArchives(c, mid, aids, ip)
	if err != nil {
		log.Error("s.p.BatchArchives aids (%v), ip(%s) err(%v)", aids, ip, err)
		return
	}
	if len(avm) == 0 {
		return
	}
	for _, v := range ui.Tops {
		if v != nil {
			if av, ok := avm[v.AID]; ok {
				v.Title = av.Archive.Title
			}
		}
	}
	return
}

// IncomeList get up status.
func (s *Service) IncomeList(c context.Context, mid int64, ty, pn, ps int, ip string) (l *growup.IncomeList, err error) {
	if l, err = s.growup.IncomeList(c, mid, ty, pn, ps, ip); err != nil {
		log.Error("s.growup.IncomeList(%d) error(%v)", mid, err)
		return
	}
	if l == nil || len(l.Data) == 0 {
		log.Error("s.growup.IncomeList mid(%d) type(%d) data(%v)", mid, ty, l)
		return
	}
	aids := make([]int64, 0, len(l.Data))
	for _, v := range l.Data {
		if v != nil {
			aids = append(aids, v.AID)
		}
	}
	avm, err := s.p.BatchArchives(c, mid, aids, ip)
	if err != nil {
		log.Error("s.p.BatchArchives aids (%v), ip(%s) err(%v)", aids, ip, err)
		return
	}
	if len(avm) == 0 {
		return
	}
	for _, v := range l.Data {
		if v != nil {
			if av, ok := avm[v.AID]; ok {
				v.Title = av.Archive.Title
			}
		}
	}
	return
}

// BreachList get up info.
func (s *Service) BreachList(c context.Context, mid int64, pn, ps int, ip string) (ui *growup.BreachList, err error) {
	if ui, err = s.growup.BreachList(c, mid, pn, ps, ip); err != nil {
		log.Error("s.growup.BreachList(%d) error(%v)", mid, err)
	}
	return
}
