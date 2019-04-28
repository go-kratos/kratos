package service

import (
	"context"
	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
)

func (s *Service) pushLimit(fan int64, upper int64, g *dao.FanGroup, noLimitFans *map[int64]int) (allow bool) {
	if _, ok := (*noLimitFans)[fan]; ok {
		log.Info("included by pushlimit(%d) upper(%d) group.name(%s) without pushlimit)", fan, upper, g.Name)
		allow = true
		return
	}
	if !s.fanLimit(fan, g) {
		log.Info("excluded by fanlimit(%d) upper(%d) group.name(%s)", fan, upper, g.Name)
		return
	}
	if !s.perUpperLimit(fan, upper, g) {
		log.Info("excluded by perupperlimit(%d) upper(%d) group.name(%s)", fan, upper, g.Name)
		return
	}

	allow = true
	return
}

//perUpperLimit 粉丝的在指定周期内的次数限制
func (s *Service) perUpperLimit(fan int64, upper int64, g *dao.FanGroup) (allow bool) {
	limit := g.PerUpperLimit
	//没有次数限制
	if limit <= 0 {
		allow = true
		return
	}
	//有次数限制
	var (
		now int
		err error
	)
	if now, err = s.dao.GetPerUpperLimitCache(context.TODO(), fan, upper); err != nil {
		log.Error("s.dao.GetPerUpperLimitCache err(%v), fan(%d), upper(%d) group.name(%s)", err, fan, upper, g.Name)
		return
	}

	now = now + 1
	if limit < now {
		return
	}
	if err = s.dao.AddPerUpperLimitCache(context.TODO(), fan, upper, now, g.LimitExpire); err != nil {
		log.Error("s.dao.AddPerUpperLimitCache err(%v), fan(%d), upper(%d), value(%d) group.name(%s)", err, fan, upper, now, g.Name)
		return
	}
	allow = true
	return
}

//fanLimit 粉丝的在指定周期内的次数限制
func (s *Service) fanLimit(fan int64, g *dao.FanGroup) (allow bool) {
	limit := g.Limit
	//没有次数限制
	if limit <= 0 {
		allow = true
		return
	}
	//有次数限制
	var (
		now int
		err error
	)
	if now, err = s.dao.GetFanLimitCache(context.TODO(), fan, g.RelationType); err != nil {
		log.Error("s.dao.GetFanLimitCache err(%v), fan(%d), group.name(%s)", err, fan, g.Name)
		return
	}

	now = now + 1
	if limit < now {
		return
	}
	if err = s.dao.AddFanLimitCache(context.TODO(), fan, g.RelationType, now, g.LimitExpire); err != nil {
		log.Error("s.dao.AddFanLimitCache err(%v), fan(%d), value(%d), group.name(%s)", err, fan, now, g.Name)
		return
	}
	allow = true
	return
}

// limit limits push frequency.
func (s *Service) limit(upper int64) (limit bool) {
	if s.dao.UpperLimitExpire == 0 {
		return
	}
	limit = true
	exist, err := s.dao.ExistUpperLimitCache(context.TODO(), upper)
	if err != nil {
		log.Error("s.dao.ExistUpperLimitCache(%d) error(%v)", upper, err)
		return
	}
	if exist {
		return
	}

	if err = s.dao.AddUpperLimitCache(context.TODO(), upper); err != nil {
		log.Error("s.dao.AddUpperLimitCache(%d) error(%v)", upper, err)
		return
	}

	limit = false
	return
}

func (s *Service) noPushLimitFans(upper int64, fanGroupKey string, fans *[]int64) (noLimitFans map[int64]int) {
	noLimitFans = map[int64]int{}
	g := s.dao.FanGroups[fanGroupKey]
	// 没有频率限制，没有免限制范围的概念
	if g.Limit <= 0 {
		return
	}
	// 只有特殊关注，才有免限制范围
	if g.RelationType != model.RelationSpecial {
		return
	}
	// 没有hbase表，没有免限制的概念
	if len(g.HBaseTable) == 0 {
		return
	}
	// abtest 不走免限制逻辑
	if g.Hitby == model.GroupDataTypeAbtest || g.Hitby == model.GroupDataTypeAbComparison {
		return
	}
	f := *fans
	hit, _ := s.dao.FansByHBase(upper, fanGroupKey, &f)
	for _, mid := range hit {
		noLimitFans[mid] = 1
	}
	return
}
