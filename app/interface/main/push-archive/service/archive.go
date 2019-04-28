package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/log"
)

const (
	_insertAct = "insert"
	_updateAct = "update"
	_deleteAct = "delete"

	_archiveTable = "archive"

	_pushRetryTimes = 3
	_pushPartSize   = 500000 // 每次请求最多推50w mid

	_hbaseBatch = 100
)

func (s *Service) consumeArchiveproc() {
	defer s.wg.Done()
	var (
		err  error
		msgs = s.archiveSub.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("s.ArchiveSub has been closed.")
			return
		}
		msg.Commit()
		s.arcMo++
		log.Info("consume archive key(%s) offset(%d) message(%s)", msg.Key, msg.Offset, msg.Value)
		m := new(model.Message)
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.Table != _archiveTable {
			continue
		}
		n := new(model.Archive)
		if err = json.Unmarshal(m.New, n); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
			continue
		}
		// 稿件过审
		switch m.Action {
		case _insertAct:
			if !n.IsNormal() {
				log.Info("archive (%d) upper(%d) is not normal when insert, no need to push", n.ID, n.Mid)
				continue
			}
		case _updateAct:
			o := new(model.Archive)
			if err = json.Unmarshal(m.Old, o); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", m.Old, err)
				continue
			}
			if !n.IsNormal() || o.State == n.State || o.PubTime == n.PubTime {
				log.Info("archive (%d) upper(%d) is not normal when update, no need to push", n.ID, n.Mid)
				continue
			}
		default:
			continue
		}
		dao.PromInfo("archive_push")
		// 在指定周期内，up主只推送一次
		if s.limit(n.Mid) {
			log.Info("upper's push limit upper(%d) aid(%d)", n.Mid, n.ID)
			continue
		}
		log.Info("noticeFans start, mid(%d) aid(%d)", n.Mid, n.ID)
		err = s.noticeFans(n) // notice fans
		log.Info("noticeFans end, mid(%d) aid(%d) error(%v)", n.Mid, n.ID, err)
	}
}

// isForbidTime 当前时间是否被禁止推送消息
func (s *Service) isForbidTime() (forbid bool, err error) {
	var forbidStart, forbidEnd time.Time
	now := time.Now()
	for _, f := range s.ForbidTimes {
		forbidStart, err = s.getTodayTime(f.PushForbidStartTime)
		if err != nil {
			log.Error("isForbidTime getTodayTime(%s) error(%v)", f.PushForbidStartTime, err)
			return
		}
		forbidEnd, err = s.getTodayTime(f.PushForbidEndTime)
		if err != nil {
			log.Error("isForbidTime getTodayTime(%s) error(%v)", f.PushForbidEndTime, err)
			return
		}
		if now.After(forbidStart) && now.Before(forbidEnd) {
			forbid = true
			break
		}
	}
	return
}

func (s *Service) noticeFans(arc *model.Archive) (err error) {
	isPGC := s.isPGC(arc)
	mids, counter, err := s.fans(arc.Mid, arc.ID, isPGC)
	if err != nil {
		return
	}
	log.Info("noticeFans get fans mid(%d) aid(%d) fans number(%d)", arc.Mid, arc.ID, counter)
	arc.Author = s.name(arc.Mid)
	for gKey, list := range mids {
		g := s.dao.FanGroups[gKey]
		params := model.NewBatchParam(map[string]interface{}{
			"relationType": g.RelationType,
			"archive":      arc,
			"group":        s.getGroupParam(g),
			"msgTemplate":  g.MsgTemplate,
		}, model.PushParamHandler)
		dao.Batch(&list, _pushPartSize, _pushRetryTimes, params, s.dao.NoticeFans)
	}
	return
}

func (s *Service) fans(upper int64, aid int64, isPGC bool) (mids map[string][]int64, counter int, err error) {
	mids = make(map[string][]int64)
	// 从 HBase 获取upper粉丝数据
	fans, err := s.dao.Fans(context.TODO(), upper, isPGC)
	ln := len(fans)
	if err != nil {
		log.Error("fans from hbase upper(%d) error(%v)", upper, err)
		return
	}
	log.Info("filter fans -- get all fans of upper(%d), fans num(%d)", upper, ln)
	if ln == 0 {
		return
	}
	dao.PromInfoAdd("archive_fans", int64(ln))
	// 筛选粉丝，分配到对应的实验组
	hitFans := s.group(upper, fans)
	hour := time.Now().Hour()
	for gKey, hits := range hitFans {
		log.Info("before filter by active time, upper(%d) group(%s) fans(%d)", upper, gKey, len(hits))
		if len(hits) <= 0 {
			continue
		}
		g := s.dao.FanGroups[gKey]
		var activeFans, excluded []int64
		if g.RelationType == model.RelationAttention {
			activeFans, excluded = s.dao.FansByActiveTime(hour, &hits)
			noLimitFans := s.noPushLimitFans(upper, gKey, &activeFans)
			for _, mid := range activeFans {
				if !s.filterUserSetting(mid, g.RelationType) {
					excluded = append(excluded, mid)
					log.Info("filter fans: excluded by usersettings, fan(%d), relationtype(%d), upper(%d) table(%s)", mid, g.RelationType, upper, g.HBaseTable)
					continue
				}
				if !s.pushLimit(mid, upper, g, &noLimitFans) {
					excluded = append(excluded, mid)
					continue
				}
				mids[gKey] = append(mids[gKey], mid)
				dao.PromInfo(fmt.Sprintf("%s__send", g.Name))
				log.Info("filter fans: included by pushlimit: upper(%d)'s fan(%d) group.name(%s)", upper, mid, g.Name)
			}
		} else {
			mids[gKey] = hits
		}
		counter += len(mids[gKey])
		log.Info("upper(%d) aid(%d) group(%s) fans(%d)", upper, aid, gKey, len(mids[gKey]))
		// 统计数据落表
		s.formStatisticsProc(aid, g.Name, mids[gKey], &excluded)
	}
	return
}

// group 粉丝分组: 分组优先级(优先级 & 分组的可用性)-->分组规则(默认default,hbase过滤走hbase)
func (s *Service) group(upper int64, fans map[int64]int) (list map[string][]int64) {
	list = make(map[string][]int64)
	// 将粉丝分为 特殊关注 和 普通关注 两类
	attentions, specials := s.dao.FansByProportion(upper, fans)
	log.Info("group count upper(%d) attentions(%d) specials(%d)", upper, len(attentions), len(specials))
	// 禁止推送时间判断, 免打扰用户(针对普通关注)
	isForbid, err := s.isForbidTime()
	if err != nil {
		log.Error("isForbidTime upper(%d) error(%v)", upper, err)
	}
	// 按顺序筛选出每个分组的数据
	for _, gkey := range s.dao.GroupOrder {
		g := s.dao.FanGroups[gkey]
		if g.RelationType == model.RelationAttention && isForbid {
			log.Info("forbid by time upper(%d) group(%+v)", upper, g)
			continue
		}
		// 优先处理abtest的流量
		if g.Hitby == model.GroupDataTypeAbtest {
			switch g.RelationType {
			case model.RelationAttention:
				var pool, exists, ex []int64
				pool, attentions = s.fansByAbtest(g, attentions)
				for _, gk := range s.dao.GroupOrder {
					gg := s.dao.FanGroups[gk]
					if gg.RelationType != model.RelationAttention ||
						gg.Hitby != model.GroupDataTypeHBase ||
						gg.Name == "ai:pushlist_follow_recent" {
						continue
					}
					ex, pool = s.dao.FansByHBase(upper, gk, &pool)
					exists = append(exists, ex...)
				}
				list[gkey] = s.filterByBlacklist(upper, exists)
			case model.RelationSpecial:
				list[gkey], specials = s.fansByAbtest(g, specials)
			}
			continue
		}
		if g.Hitby == model.GroupDataTypeAbComparison {
			// 对照组保持原来逻辑
			switch g.RelationType {
			case model.RelationAttention:
				var pool, exists []int64
				pool, attentions = s.fansByAbtest(g, attentions)
				for _, gk := range s.dao.GroupOrder {
					gg := s.dao.FanGroups[gk]
					if gg.RelationType != model.RelationAttention || gg.Hitby != model.GroupDataTypeHBase {
						continue
					}
					exists, pool = s.dao.FansByHBase(upper, gk, &pool)
					list[gkey] = append(list[gkey], exists...)
				}
			case model.RelationSpecial:
				list[gkey], specials = s.fansByAbtest(g, specials)
			}
			continue
		}
		switch {
		case g.RelationType == model.RelationAttention && g.Hitby == model.GroupDataTypeDefault:
			list[gkey] = attentions
			attentions = []int64{}
		case g.RelationType == model.RelationSpecial && g.Hitby == model.GroupDataTypeDefault:
			list[gkey] = specials
			specials = []int64{}
		case g.RelationType == model.RelationAttention && g.Hitby == model.GroupDataTypeHBase:
			list[gkey], attentions = s.dao.FansByHBase(upper, gkey, &attentions)
		case g.RelationType == model.RelationSpecial && g.Hitby == model.GroupDataTypeHBase:
			list[gkey], specials = s.dao.FansByHBase(upper, gkey, &specials)
		default:
			log.Error("group failed for grouporder(%s) & group.relationtype(%d) & hitby(%s)", gkey, g.RelationType, g.Hitby)
		}
	}
	return
}

func (s *Service) filterByBlacklist(upper int64, fans []int64) (result []int64) {
	for {
		var mids []int64
		l := len(fans)
		if l == 0 {
			return
		}
		if l <= _hbaseBatch {
			mids = fans[:]
			fans = []int64{}
		} else {
			mids = fans[:_hbaseBatch]
			fans = fans[_hbaseBatch:]
		}
		exists, notExists := s.dao.ExistsInBlacklist(context.Background(), upper, mids)
		result = append(result, notExists...) // 不在黑名单中的用户直接可推
		exists, notExists = s.dao.ExistsInWhitelist(context.Background(), upper, exists)
		result = append(result, exists...) // 存在黑名单中但是后来被恢复到白名单中的用户可推
		for _, mid := range notExists {
			// 只出现在黑名单中的用户不推
			log.Warn("filter by blacklist, mid(%d)", mid)
		}
	}
}

// getGroupParam g.name被_分隔后的第1个/2个值（比如：ai:pushlist_offline_up的是offline, special的就是special）
func (s *Service) getGroupParam(g *dao.FanGroup) string {
	p := strings.Split(g.Name, "_")
	if len(p) == 1 {
		return p[0]
	}
	if g.RelationType == model.RelationSpecial {
		return "s_" + p[1]
	}
	return p[1]
}

// filterUserSetting 普通关注up主的粉丝中，推送开启配置为全部推送的人很少，因此，只要其开启了特殊关注/全部关注／没有设置的就发送；特殊关注等同
func (s *Service) filterUserSetting(mid int64, relationType int) bool {
	if relationType == model.RelationSpecial && s.userSettings[mid] != nil && s.userSettings[mid].Type == model.PushTypeForbid {
		return false
	}
	if relationType == model.RelationAttention && s.userSettings[mid] != nil && s.userSettings[mid].Type == model.PushTypeForbid {
		return false
	}
	return true
}

func (s *Service) name(mid int64) (name string) {
	arg := &accmdl.ArgMid{Mid: mid}
	info, err := s.accRPC.Info3(context.TODO(), arg)
	if err != nil {
		dao.PromError("archive:获取作者信息")
		log.Error("s.accRPC.Info3(%+v) error(%v)", arg, err)
		return
	}
	name = info.Name
	return
}

func (s *Service) isPGC(arc *model.Archive) bool {
	return (arc.Attribute >> model.AttrBitIsPGC & 1) == 1
}

// Test for test
// TODO delete, for test
func (s *Service) Test(arc *model.Archive) {
	if arc.ID == 0 || arc.Mid == 0 {
		return
	}
	go s.noticeFans(arc)
} // TODO delete, for test
