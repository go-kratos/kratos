package service

import (
	"context"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
)

// SyncAllUser 同步旧user——info到新db.
// FIXME 切新db后删除.
func (s *Service) SyncAllUser(c context.Context) {
	var (
		err      error
		maxID    int
		size     = s.c.Property.BatchSize
		ids      = []int64{}
		ousers   = make(map[int64]*model.VipUserInfoOld, size)
		nusers   = make(map[int64]*model.VipUserInfo, _ps)
		updateDB = s.c.Property.UpdateDB
		nu       *model.VipUserInfo
		ok       bool
	)
	if maxID, err = s.dao.SelOldUserInfoMaxID(context.TODO()); err != nil {
		log.Error("sync job s.dao.SelOldUserInfoMaxID err(%+v)", err)
		return
	}
	page := maxID / size
	if maxID%size != 0 {
		page++
	}
	log.Info("sync job vip_user_info total(%d)", page)
	for i := 0; i < page; i++ {
		log.Info("sync job vip_user_info page index(%d) total(%d)", i, page)
		startID := i * size
		endID := (i + 1) * size
		if endID > maxID {
			endID = maxID
		}
		if ousers, err = s.dao.SelOldUserInfoMaps(context.TODO(), startID, endID); err != nil {
			log.Error("sync job s.dao.SelOldUserInfoMaps(%d, %d) err(%+v)", startID, endID, err)
			return
		}
		j := 1
		for _, v := range ousers {
			ids = append(ids, v.Mid)
			if j%_ps == 0 || j == len(ousers) {
				if nusers, err = s.dao.SelVipByIds(context.TODO(), ids); err != nil {
					return
				}
				for _, mid := range ids {
					var ou *model.VipUserInfoOld
					if ou, ok = ousers[mid]; !ok {
						log.Warn("sync job old not found %d", mid)
						continue
					}
					if nu, ok = nusers[mid]; !ok {
						log.Warn("sync job need insert to new %d, old(%+v), toNew(%+v)", mid, ou, ou.ToNew())
						if updateDB {
							s.dao.SyncAddUser(context.Background(), ou.ToNew())
						}
						continue
					}
					if ou.RecentTime <= 0 {
						ou.RecentTime = ou.Mtime
					}
					if nu.Type != ou.Type ||
						nu.Status != ou.Status ||
						!nu.StartTime.Time().Equal(ou.StartTime.Time()) ||
						!nu.OverdueTime.Time().Equal(ou.OverdueTime.Time()) ||
						!nu.AnnualVipOverdueTime.Time().Equal(ou.AnnualVipOverdueTime.Time()) ||
						!nu.Ctime.Time().Equal(ou.Ctime.Time()) ||
						!nu.Mtime.Time().Equal(ou.Mtime.Time()) ||
						nu.PayType != ou.IsAutoRenew ||
						nu.PayChannelID != ou.PayChannelID ||
						!nu.IosOverdueTime.Time().Equal(ou.IosOverdueTime.Time()) ||
						nu.Ver != ou.Ver ||
						!nu.RecentTime.Time().Equal(ou.RecentTime.Time()) {
						log.Warn("sync job need update to new %d, old(%+v), new(%+v), toNew(%+v)", mid, ou, nu, ou.ToNew())
						if updateDB {
							s.dao.SyncUpdateUser(context.Background(), ou.ToNew(), nu.Ver)
						}
						continue
					}
				}
				log.Info("sync job vip_user_info page index(%d) ids(%+v)", j, ids)
				// reset
				ids = []int64{}
			}
			j++
		}
		log.Info("sync job vip_user_info page index(%d) end", i)
		time.Sleep(time.Millisecond * _defsleepmsec)
	}
}
