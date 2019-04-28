/*
rebuild: user portrait score reset to normal if it's score large than punishment threshold score
*/

package service

import (
	"context"
	"time"

	spy "go-common/app/service/main/spy/model"
	"go-common/library/log"
)

const (
	_normal = 0
	_ps     = 100
)

func (s *Service) reBuild() {
	var (
		err   error
		count int64
	)
	current := time.Now()
	before30d, _ := time.ParseDuration("-720h")
	before31d, _ := time.ParseDuration("-744h")
	start := current.Add(before31d)
	end := current.Add(before30d)
	log.Info("ReBuild task start: start:(%s) end:(%s))", start, end)
	for t := 0; t < int(s.c.Property.UserInfoShard); t++ {
		if count, err = s.dao.ReBuildMidCount(context.TODO(), t, _normal, start, end); err != nil {
			log.Error("s.dao.ReBuildMidCount(%s, %s), err(%v)", start, end, err)
			continue
		}
		log.Info("ReBuild task: index:%d, count:%d)", t, count)
		if count <= 0 {
			continue
		}
		total := count / _ps
		log.Info("ReBuild task: shard:%d, count:%d, total:%d)", t, count, total)
		for i := 0; int64(i) <= total; i++ {
			midList, err := s.dao.ReBuildMidList(context.TODO(), t, _normal, start, end, _ps)
			if err != nil {
				log.Error("s.dao.ReBuildMidList(%s, %s, %d, %d)", start, end, i, _ps)
				continue
			}
			for _, mid := range midList {
				if err := s.spyRPC.ReBuildPortrait(context.TODO(), &spy.ArgReBuild{Mid: mid, Reason: "自动恢复行为得分"}); err != nil {
					log.Error("s.spyRPC.ReBuildPortrait(%d), err:%v", mid, err)
					continue
				}
				log.Info("ReBuild task: mid(%d) ReBuild Portrait success)", mid)
			}
		}
	}
	log.Info("ReBuild task end: start:(%s) end:(%s))", start, end)
}
