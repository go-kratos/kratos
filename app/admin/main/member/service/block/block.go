package block

import (
	"context"
	"sync"
	"time"

	model "go-common/app/admin/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// Search .
func (s *Service) Search(c context.Context, mids []int64) (infos []*model.BlockInfo, err error) {
	var (
		users       []*model.DBUser
		userDetails []*model.DBUserDetail
		eg          errgroup.Group
		mapMu       sync.Mutex
		userMap     = make(map[int64]*model.BlockInfo)
	)
	if users, err = s.dao.Users(c, mids); err != nil {
		return
	}
	if userDetails, err = s.dao.UserDetails(c, mids); err != nil {
		return
	}
	infos = make([]*model.BlockInfo, 0, len(mids))
	for _, m := range mids {
		mid := m
		eg.Go(func() (err error) {
			info := &model.BlockInfo{
				MID: mid,
			}
			// 1. account 数据
			if info.Nickname, info.TelStatus, info.Level, info.RegTime, err = s.AccountInfo(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			if info.Nickname == "" {
				log.Info("user mid(%d) not found", info.MID)
				return
			}
			// 2. 封禁状态
			for i := range users {
				if users[i].MID == mid {
					info.ParseStatus(users[i])
					break
				}
			}
			// 3. 封禁次数
			for i := range userDetails {
				if userDetails[i].MID == mid {
					info.BlockCount = userDetails[i].BlockCount
					break
				}
			}
			// 4. spy 分值
			if info.SpyScore, err = s.SpyScore(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
				info.SpyScore = -1
			}
			// 5. figure 排名
			if info.FigureRank, err = s.FigureRank(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
				info.FigureRank = -1
			}
			// 6. extra 额外账号信息
			if info.Tel, err = s.telInfo(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
				info.Tel = "N/A"
			}
			if info.Mail, err = s.mailInfo(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
				info.Mail = "N/A"
			}
			if info.Username, err = s.userID(context.Background(), mid); err != nil {
				log.Error("%+v", err)
				err = nil
				info.Username = "N/A"
			}

			mapMu.Lock()
			userMap[info.MID] = info
			mapMu.Unlock()
			return
		})
	}
	eg.Wait()
	for _, mid := range mids {
		if info, ok := userMap[mid]; ok {
			infos = append(infos, info)
		}
	}
	return
}

// History .
func (s *Service) History(c context.Context, mid int64, ps, pn int, desc bool) (history *model.BlockDetail, err error) {
	var (
		start     = (pn - 1) * ps
		limit     = ps
		dbHistory []*model.DBHistory
		dbUser    *model.DBUser
	)
	history = &model.BlockDetail{}
	if dbUser, err = s.dao.User(c, mid); err != nil {
		return
	}
	if dbUser != nil {
		history.Status = dbUser.Status
	}
	if history.Total, err = s.dao.HistoryCount(c, mid); err != nil {
		return
	}
	if dbHistory, err = s.dao.History(c, mid, start, limit, desc); err != nil {
		return
	}
	history.History = make([]*model.BlockHistory, 0, len(dbHistory))
	for i := range dbHistory {
		his := &model.BlockHistory{}
		his.ParseDB(dbHistory[i])
		history.History = append(history.History, his)
	}
	return
}

// BatchBlock .
func (s *Service) BatchBlock(c context.Context, p *model.ParamBatchBlock) (err error) {
	var (
		tx       *xsql.Tx
		duration = time.Duration(p.Duration) * time.Hour * 24
		source   model.BlockSource
		stime    = time.Now()
	)
	if tx, err = s.dao.BeginTX(c); err != nil {
		return
	}
	if p.Source == model.BlockMgrSourceSys {
		// 系统封禁
		source = model.BlockSourceSys
	} else if p.Source == model.BlockMgrSourceCredit {
		// 小黑屋封禁
		source = model.BlockSourceBlackHouse
		if err = s.dao.BlackhouseBlock(context.Background(), p); err != nil {
			return
		}
	}
	for _, mid := range p.MIDs {
		if err = s.action(c, tx, mid, p.AdminID, p.AdminName, source, p.Area, p.Reason, p.Comment, p.Action, duration, p.Notify, stime); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	// 发送站内信
	if p.Notify {
		s.mission(func() {
			if notifyErr := s.notifyMSG(context.Background(), p.MIDs, source, p.Action, p.Area, p.Reason, p.Duration); notifyErr != nil {
				log.Error("%+v", notifyErr)
				return
			}
		})
	}
	s.mission(func() {
		s.AddAuditLog(context.Background(), p.Action, p.AdminID, p.AdminName, p.MIDs, duration, source, p.Area, p.Reason, p.Comment, p.Notify, stime)
	})
	s.asyncPurgeCache(p.MIDs)
	return
}

// BatchRemove .
func (s *Service) BatchRemove(c context.Context, p *model.ParamBatchRemove) (err error) {
	var (
		tx    *xsql.Tx
		stime = time.Now()
	)
	if tx, err = s.dao.BeginTX(c); err != nil {
		return
	}

	for _, mid := range p.MIDs {
		if err = s.action(c, tx, mid, p.AdminID, p.AdminName, model.BlockSourceManager, model.BlockAreaNone, "", p.Comment, model.BlockActionAdminRemove, 0, p.Notify, stime); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	// 发送站内信
	if p.Notify {
		s.mission(func() {
			if notifyErr := s.notifyMSG(context.Background(), p.MIDs, model.BlockSourceManager, model.BlockActionAdminRemove, model.BlockAreaNone, "", 0); notifyErr != nil {
				log.Error("%+v", notifyErr)
				return
			}
		})
	}
	s.mission(func() {
		s.AddAuditLog(context.Background(), model.BlockActionAdminRemove, p.AdminID, p.AdminName, p.MIDs, 0, model.BlockSourceManager, model.BlockAreaNone, "", p.Comment, p.Notify, stime)
	})
	s.asyncPurgeCache(p.MIDs)
	return
}

// notifyMSG .
func (s *Service) notifyMSG(c context.Context, mids []int64, source model.BlockSource, action model.BlockAction, area model.BlockArea, reason string, days int64) (err error) {
	code, title, content := s.MSGInfo(source, action, area, reason, days)
	log.Info("block admin title : %s , content : %s , mids : %+v", title, content, mids)
	if err = s.dao.SendSysMsg(context.Background(), code, mids, title, content, ""); err != nil {
		return
	}
	return
}

func (s *Service) action(c context.Context, tx *xsql.Tx, mid int64, adminID int64, adminName string, source model.BlockSource, area model.BlockArea, reason, comment string, action model.BlockAction, duration time.Duration, notify bool, stime time.Time) (err error) {
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
			StartTime: stime,
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
		case model.BlockSourceBlackHouse:
			blockStatus = model.BlockStatusCredit
		default:
			blockStatus = model.BlockStatusLimit
		}
		s.mission(func() {
			if err = s.dao.UpdateAddBlockCount(context.Background(), mid); err != nil {
				log.Error("%+v", err)
			}
		})
	case model.BlockActionForever:
		blockStatus = model.BlockStatusForever
		s.mission(func() {
			if err = s.dao.UpdateAddBlockCount(context.Background(), mid); err != nil {
				log.Error("%+v", err)
			}
		})
	default:
		err = errors.Errorf("unknown block action [%d]", action)
		return
	}
	if err = s.dao.TxUpdateUser(c, tx, mid, blockStatus); err != nil {
		return
	}
	return
}

func (s *Service) userID(c context.Context, mid int64) (id string, err error) {
	return "N/A", nil
}

func (s *Service) mailInfo(c context.Context, mid int64) (mail string, err error) {
	if mail, err = s.dao.MailInfo(c, mid); mail == "" {
		mail = "N/A"
	}
	return
}

// TelInfo .
func (s *Service) telInfo(c context.Context, mid int64) (tel string, err error) {
	if tel, err = s.dao.TelInfo(c, mid); err != nil {
		return
	}
	if len(tel) == 0 {
		tel = "N/A"
		return
	}
	if len(tel) < 4 {
		tel = tel[:1] + "****"
		return
	}
	tel = tel[:3] + "****" + tel[len(tel)-4:]
	return
}

func (s *Service) asyncPurgeCache(mids []int64) {
	_ = s.cache.Do(context.Background(), func(ctx context.Context) {
		for _, mid := range mids {
			if cacheErr := s.dao.DeleteUserCache(ctx, mid); cacheErr != nil {
				log.Error("%+v", cacheErr)
			}
			if cacheErr := s.dao.DeleteUserDetailCache(ctx, mid); cacheErr != nil {
				log.Error("%+v", cacheErr)
			}
			if databusErr := s.accountNotify(ctx, mid); databusErr != nil {
				log.Error("%+v", databusErr)
			}
		}
	})
}
