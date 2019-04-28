package service

import (
	"context"
	"fmt"
	"sort"

	creditMDL "go-common/app/admin/main/credit/model"
	"go-common/app/admin/main/credit/model/blocked"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

// Infos  deal info data.
func (s *Service) Infos(c context.Context, arg *blocked.ArgBlockedSearch) (list []*blocked.Info, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.Blocked(c, arg)
	if err != nil {
		log.Error("s.searchDao.Blocked error (%v)", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	var (
		infoMap map[int64]*account.Info
		uids    = make([]int64, len(ids))
	)
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
		return
	}
	for _, v := range list {
		uids = append(uids, v.UID)
	}
	if infoMap, err = s.accDao.RPCInfos(c, uids); err != nil {
		log.Error("s.accDao.RPCInfos(%s) error(%v)", xstr.JoinInts(uids), err)
		err = nil
	}
	for _, v := range list {
		v.OPName = s.Managers[v.OperID]
		if v.OPName == "" {
			v.OPName = v.OOPName
		}
		if in, ok := infoMap[v.UID]; ok {
			v.UName = in.Name
		}
		v.ReasonTypeDesc = blocked.ReasonTypeDesc(v.ReasonType)
		v.PublishStatusDesc = blocked.PStatusDesc[v.PublishStatus]
		v.OriginTypeDesc = blocked.OriginTypeDesc[v.OriginType]
		v.BlockedTypeDesc = blocked.BTypeDesc[v.BlockedType]
		v.BlockedDaysDesc = blocked.BDaysDesc(v.BlockedDays, v.MoralNum, v.PunishType, v.BlockedForever)
	}
	return
}

// InfosEx  export info list
func (s *Service) InfosEx(c context.Context, arg *blocked.ArgBlockedSearch) (list []*blocked.Info, err error) {
	var (
		count int
		pager *blocked.Pager
		g     errgroup.Group
		ps    = 500
	)
	if list, pager, err = s.Infos(c, arg); err != nil {
		log.Error("s.Infos(%+v) error(%v)", arg, err)
		return
	}
	if pager == nil {
		log.Warn("arg(%+v) info search data empty!", arg)
		return
	}
	count = pager.Total / ps
	if pager.Total%ps != 0 {
		count++
	}
	lCh := make(chan []*blocked.Info, count)
	for pn := 1; pn <= count; pn++ {
		tmpPn := pn
		g.Go(func() (err error) {
			var gInfo []*blocked.Info
			gArg := &blocked.ArgBlockedSearch{
				Keyword:       arg.Keyword,
				UID:           arg.UID,
				OPID:          arg.OPID,
				OriginType:    arg.OriginType,
				BlockedType:   arg.BlockedType,
				PublishStatus: arg.PublishStatus,
				Start:         arg.Start,
				End:           arg.End,
				PN:            tmpPn,
				PS:            ps,
				Order:         arg.Order,
				Sort:          arg.Sort,
			}
			gInfo, _, err = s.Infos(c, gArg)
			if err != nil {
				log.Error("s.Infos(%+v) error(%v)", gArg, err)
				err = nil
				return
			}
			lCh <- gInfo
			return
		})
	}
	g.Wait()
	close(lCh)
	for bInfo := range lCh {
		list = append(list, bInfo...)
	}
	sort.Slice(list, func(i int, j int) bool {
		return list[i].ID < list[j].ID
	})
	return
}

// Publishs get publishs data
func (s *Service) Publishs(c context.Context, arg *blocked.ArgPublishSearch) (list []*blocked.Publish, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.Publish(c, arg)
	if err != nil {
		log.Error("s.searchDao.Publish error (%v)", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
	}
	for _, v := range list {
		v.OPName = s.Managers[v.OPID]
		v.PublishTypeDesc = blocked.PTypeDesc[v.Type]
		v.PublishStatusDesc = blocked.PStatusDesc[v.PublishStatus]
		v.StickStatusDesc = blocked.SStatusDesc[v.StickStatus]
	}
	return
}
