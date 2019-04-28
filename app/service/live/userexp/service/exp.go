package service

import (
	"context"
	"go-common/app/service/live/userexp/model"
	"go-common/library/log"
	"strconv"
	"time"
)

// Level 获取目标用户的等级信息
func (s *Service) Level(c context.Context, uid int64) (level *model.Level, err error) {
	// 缓存查询
	//cacheHealth := true
	//level, err = s.dao.LevelCache(c, uid)
	//if err != nil {
	//	// 缓存异常
	//	log.Error("[service.exp|level] s.dao.LevelCache error(%v)", err)
	//	cacheHealth = false
	//}

	//if level != nil && level.Uid != 0 {
	//	// 命中缓存直接返回
	//	return
	//}

	// DB查询
	exp, err := s.dao.Exp(c, uid)
	if err != nil {
		log.Error("[service.exp|level] s.dao.Exp error(%v)", err)
		return
	}

	// 格式化Level结构
	level = model.FormatLevel(exp)

	//// 写入缓存
	//if cacheHealth {
	//	s.cache.Save(func() {
	//		s.dao.SetLevelCache(context.TODO(), level)
	//	})
	//}
	return
}

// Exp 获取目标用户的经验信息
func (s *Service) Exp(c context.Context, uid int64) (exp *model.Exp, err error) {
	exp, err = s.dao.Exp(c, uid)
	if err != nil {
		log.Error("[service.Exp|level] s.dao.Exp error(%v)", err)
		return
	}
	return
}

// MultiGetLevel 批量获取用户等级信息
func (s *Service) MultiGetLevel(c context.Context, uids []int64) (level []*model.Level, err error) {
	// 缓存查询
	//cacheHealth := true
	//level, missUids, err := s.dao.MultiLevelCache(c, uids)
	//if err != nil {
	//	// 缓存异常
	//	log.Error("[service.exp|MultiGetLevel] s.dao.MultiGetLevel error(%v)", err)
	//	cacheHealth = false
	//} else if len(missUids) == 0 {
	//	// 缓存全命中直接返回
	//	return
	//}

	// DB查询
	exps, err := s.dao.MultiExp(c, uids)
	if err != nil {
		log.Error("[service.exp|MultiGetLevel] s.dao.Exp error(%v)", err)
		return
	}
	for _, exp := range exps {
		// 格式化Level结构
		lv := model.FormatLevel(exp)
		//// 写入缓存
		//if cacheHealth {
		//	s.cache.Save(func() {
		//		s.dao.SetLevelCache(c, lv)
		//	})
		//}
		// 追加数据
		level = append(level, lv)
	}

	return
}

// AddUexp 添加主播经验
func (s *Service) AddUexp(c context.Context, uid int64, uexp int64, ric map[string]string) (err error) {
	_, err = s.dao.AddUexp(c, uid, uexp)
	return
}

func (s *Service) AddUExpLog(c context.Context, uid int64, uexp int64, nowuexp int64, nowrexp int64, ric map[string]string) (err error) {
	logParams := &model.ExpLog{
		Mid:  uid,
		Uexp: uexp,
		Rexp: 0,
		Ts:   time.Now().Unix(),
		Ip:   ric[RelInfocIP],
		Content: map[string]string{
			"type":    "增加用户经验",
			"add_num": strconv.FormatInt(uexp, 10),
			"uexp":    strconv.FormatInt(nowuexp, 10),
			"rexp":    strconv.FormatInt(nowrexp, 10),
		},
	}
	s.dao.AddUserExpLog(c, logParams)
	return
}

func (s *Service) AddRExpLog(c context.Context, uid int64, rexp int64, nowuexp int64, nowrexp int64, ric map[string]string) (err error) {
	logParams := &model.ExpLog{
		Mid:  uid,
		Uexp: 0,
		Rexp: rexp,
		Ts:   time.Now().Unix(),
		Ip:   ric[RelInfocIP],
		Content: map[string]string{
			"type":    "增加主播经验",
			"add_num": strconv.FormatInt(rexp, 10),
			"uexp":    strconv.FormatInt(nowuexp, 10),
			"rexp":    strconv.FormatInt(nowrexp, 10),
		},
	}
	s.dao.AddAnchorExpLog(c, logParams)
	return
}

// AddRexp 添加用户经验
func (s *Service) AddRexp(c context.Context, uid int64, rexp int64, ric map[string]string) (err error) {
	_, err = s.dao.AddRexp(c, uid, rexp)
	return
}
