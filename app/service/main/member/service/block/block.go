package block

import (
	"context"
	"time"

	member "go-common/app/service/main/member/model"
	model "go-common/app/service/main/member/model/block"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Block .
func (s *Service) Block(c context.Context, mids []int64, source model.BlockSource, area model.BlockArea, action model.BlockAction, startTime int64, dur time.Duration, operator, reason, comment string, notify bool) (err error) {
	var (
		tx       *sql.Tx
		stime    = time.Unix(startTime, 0)
		duration = time.Duration(dur)
	)
	if mids, err = s.cleanMIDs(c, action, startTime, dur, mids); err != nil {
		return
	}
	log.Info("Block actual mids (%+v)", mids)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	for _, mid := range mids {
		if operator == "" {
			operator = "sys"
		}
		if err = s.action(c, tx, mid, -1, operator, source, area, reason, comment, action, stime, duration, notify); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	// 发送站内信
	if notify {
		s.mission(func() {
			if notifyErr := s.notifyMSG(context.Background(), mids, source, action, area, reason, int64(dur/time.Duration(time.Hour*24))); notifyErr != nil {
				log.Error("%+v", notifyErr)
				return
			}
		})
	}
	s.mission(func() {
		s.AddAuditLog(context.Background(), action, -1, operator, mids, duration, source, area, reason, comment, notify, stime)
	})
	s.asyncPurgeCache(mids)
	return
}

func (s *Service) cleanMIDs(c context.Context, action model.BlockAction, startTime int64, duration time.Duration, mids []int64) (cmids []int64, err error) {
	if action == model.BlockActionForever || action == model.BlockActionLimit {
		var (
			foreverMIDMap map[int64]struct{}
		)
		// 清洗永久封禁用户
		if foreverMIDMap, err = s.dao.UserStatusMapWithMIDs(c, model.BlockStatusForever, mids); err != nil {
			return
		}
		for _, mid := range mids {
			if _, ok := foreverMIDMap[mid]; ok {
				continue
			}
			cmids = append(cmids, mid)
		}
		// 清洗限时封禁用户
		if action == model.BlockActionLimit {
			var (
				infos   []*model.BlockInfo
				endTime = time.Unix(startTime, 0).Add(duration)
			)
			if infos, err = s.Infos(c, cmids); err != nil {
				return
			}
			cmids = []int64{}
			for _, info := range infos {
				if info.EndTime > 0 {
					if endTime.After(time.Unix(info.EndTime, 0)) {
						cmids = append(cmids, info.MID)
					}
				} else {
					cmids = append(cmids, info.MID)
				}
			}
		}
	} else {
		cmids = mids
	}
	return
}

// Remove .
func (s *Service) Remove(c context.Context, mids []int64, source model.BlockSource, area model.BlockArea, operator, reason, comment string, notify bool) (err error) {
	var (
		tx    *sql.Tx
		stime = time.Now()
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	for _, mid := range mids {
		if operator == "" {
			operator = "sys"
		}
		if err = s.action(c, tx, mid, -1, operator, source, area, reason, comment, model.BlockActionSelfRemove, stime, 0, notify); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	// 发送站内信
	if notify {
		s.mission(func() {
			if notifyErr := s.notifyMSG(context.Background(), mids, source, model.BlockActionSelfRemove, area, reason, 0); notifyErr != nil {
				log.Error("%+v", notifyErr)
				return
			}
		})
	}
	s.mission(func() {
		s.AddAuditLog(context.Background(), model.BlockActionSelfRemove, -1, operator, mids, 0, source, area, reason, comment, notify, stime)
	})
	s.asyncPurgeCache(mids)
	return
}

// Infos 获取用户封禁信息数据
// 1. mc中取得 2. db中取 3. 返回默认值（未封禁）
func (s *Service) Infos(c context.Context, mids []int64) (infos []*model.BlockInfo, err error) {
	defer func() {
		// 白名单用户永远不会被封禁，防止在bug时封禁涉政用户
		for _, info := range infos {
			if info == nil {
				continue
			}
			if _, ok := s.whiteMap[info.MID]; ok {
				info.BlockStatus = model.BlockStatusFalse
				info.StartTime = -1
				info.EndTime = -1
			}
		}
	}()
	infos = make([]*model.BlockInfo, len(mids))
	if len(mids) == 0 {
		return
	}
	var (
		mcFlag     = true
		res        map[int64]*model.MCBlockInfo
		missMidMap = make(map[int64]struct{})
		missInfos  []*model.BlockInfo
	)
	// 1. get from mc
	if res, err = s.dao.UsersCache(c, mids); err != nil {
		log.Error("%+v", err)
		mcFlag = false
	}
	for i, mid := range mids {
		if mcInfo, ok := res[mid]; ok {
			// mc hit
			info := &model.BlockInfo{}
			info.ParseMC(mcInfo, mid)
			infos[i] = info
		} else {
			missMidMap[mid] = struct{}{}
		}
	}
	// 3. get from db
	for mid := range missMidMap {
		var dbInfo *model.DBHistory
		if dbInfo, err = s.dao.UserLastHistory(c, mid); err != nil {
			return
		}
		if dbInfo == nil {
			// 加入空缓存
			missInfos = append(missInfos, s.DefaultUser(mid))
		} else {
			info := &model.BlockInfo{}
			info.ParseDB(dbInfo)
			missInfos = append(missInfos, info)
		}
	}
	// 4. fill nil infos with info
	for i, info := range infos {
		if info == nil {
			for _, missInfo := range missInfos {
				if missInfo.MID == mids[i] {
					infos[i] = missInfo
					break
				}
			}
		}
	}
	// 6. set cache
	if mcFlag && len(missInfos) != 0 {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			for _, missInfo := range missInfos {
				if theErr := s.dao.SetUserCache(ctx, missInfo.MID, missInfo.BlockStatus, missInfo.StartTime, missInfo.EndTime); theErr != nil {
					log.Error("%+v", theErr)
					return
				}
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	}
	return
}

// UserDetails 批量查询用户封禁详情
func (s *Service) UserDetails(c context.Context, mids []int64) (map[int64]*model.BlockUserDetail, error) {
	details := make(map[int64]*model.BlockUserDetail, len(mids))
	if len(mids) == 0 {
		return details, nil
	}
	var (
		mcFlag      = true
		mcRes       map[int64]*model.MCUserDetail
		dbRes       map[int64]*model.DBUserDetail
		missMids    = make([]int64, 0, len(mids))
		missDetails []*model.BlockUserDetail
		err         error
	)

	// get from mc
	if mcRes, err = s.dao.UserDetailsCache(c, mids); err != nil {
		log.Error("%+v", err)
		mcFlag = false
	}
	for _, mid := range mids {
		mcInfo, ok := mcRes[mid]
		if !ok {
			missMids = append(missMids, mid)
			continue
		}
		details[mid] = &model.BlockUserDetail{
			MID:        mid,
			BlockCount: mcInfo.BlockCount,
		}
	}

	// get from db
	if len(missMids) == 0 {
		return details, nil
	}
	if dbRes, err = s.dao.UserDetails(c, missMids); err != nil {
		log.Error("%+v", err)
		return nil, err
	}
	for _, mid := range missMids {
		blockUserDetail := s.DefaultUserDetail(mid)
		if dbDetail, ok := dbRes[mid]; ok {
			blockUserDetail = &model.BlockUserDetail{
				MID:        mid,
				BlockCount: dbDetail.BlockCount,
			}
		}
		missDetails = append(missDetails, blockUserDetail)
		details[mid] = blockUserDetail
	}

	// set missDetails to cache
	if mcFlag && len(missDetails) != 0 {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			for _, missDetail := range missDetails {
				if theErr := s.dao.SetUserDetailCache(ctx, missDetail.MID, &model.MCUserDetail{BlockCount: missDetail.BlockCount}); theErr != nil {
					log.Error("%+v", theErr)
					continue
				}
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	}
	return details, nil
}

func (s *Service) action(c context.Context, tx *sql.Tx, mid int64, adminID int64, adminName string, source model.BlockSource, area model.BlockArea, reason, comment string, action model.BlockAction, startTime time.Time, duration time.Duration, notify bool) (err error) {
	var (
		db = &model.DBHistory{
			MID:       mid,
			AdminID:   adminID,
			AdminName: adminName,
			Source:    source,
			Area:      area,
			Reason:    reason,
			Comment:   comment,
			Action:    action,
			StartTime: startTime,
			Duration:  int64(duration / time.Second),
			Notify:    notify,
		}
		blockStatus model.BlockStatus
	)
	if err = s.dao.TxInsertHistory(c, tx, db); err != nil {
		return
	}
	switch action {
	case model.BlockActionAdminRemove, model.BlockActionSelfRemove:
		blockStatus = model.BlockStatusFalse
	case model.BlockActionLimit:
		switch source {
		case model.BlockSourceBlackHouse, model.BlockSourceBplus:
			blockStatus = model.BlockStatusCredit
		default:
			blockStatus = model.BlockStatusLimit
		}
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			if err = s.dao.UpdateAddBlockCount(ctx, mid); err != nil {
				log.Error("%+v", err)
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	case model.BlockActionForever:
		blockStatus = model.BlockStatusForever
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			if err = s.dao.UpdateAddBlockCount(ctx, mid); err != nil {
				log.Error("%+v", err)
			}
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	default:
		err = errors.Errorf("unknown block action [%d]", action)
		return
	}
	if err = s.dao.TxUpdateUser(c, tx, mid, blockStatus); err != nil {
		return
	}
	return
}

// DefaultUser .
func (s *Service) DefaultUser(mid int64) (info *model.BlockInfo) {
	return &model.BlockInfo{
		MID:         mid,
		BlockStatus: model.BlockStatusFalse,
		StartTime:   -1,
		EndTime:     -1,
	}
}

// DefaultUserDetail .
func (s *Service) DefaultUserDetail(mid int64) (info *model.BlockUserDetail) {
	return &model.BlockUserDetail{
		MID:        mid,
		BlockCount: 0,
	}
}

func (s *Service) notifyMSG(c context.Context, mids []int64, source model.BlockSource, action model.BlockAction, area model.BlockArea, reason string, days int64) (err error) {
	log.Info("Block send msg mids : %+v , source : %s , action : %s , area : %s , reason : %s , days : %d", mids, source, action, area, reason, days)
	code, title, content := s.MSGInfo(source, action, area, reason, days)
	if err = s.dao.SendSysMsg(context.Background(), code, mids, title, content, ""); err != nil {
		return
	}
	return
}

func (s *Service) asyncPurgeCache(mids []int64) {
	fanoutErr := s.cache.Do(context.Background(), func(ctx context.Context) {
		for _, mid := range mids {
			if cacheErr := s.dao.DeleteUserCache(ctx, mid); cacheErr != nil {
				log.Error("%+v", cacheErr)
			}
			if cacheErr := s.dao.DeleteUserDetailCache(ctx, mid); cacheErr != nil {
				log.Error("%+v", cacheErr)
			}
			if databusErr := s.dao.NotifyPurgeCache(ctx, mid, member.ActBlockUser); databusErr != nil {
				log.Error("%+v", databusErr)
			}
		}
	})
	if fanoutErr != nil {
		log.Error("fanout do err(%+v)", fanoutErr)
	}
}
