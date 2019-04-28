package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_userRandomType = "用户随机-尾号"
)

// Policy return policy info.
func (s *Service) Policy(c context.Context, id, mid int64) (item *model.Pitem, err error) {
	var policy *model.Policy
	if policy, err = s.checkPolicy(id); err != nil {
		log.Error("s.getPolicy(%d) err(%v)", id, err)
		return
	}
	switch policy.Type {
	case _userRandomType:
		if item, err = s.userPolicy(mid, policy); err != nil {
			log.Error("s.userPolicy(%d) err(%v)", mid, err)
			return
		}
	}
	return
}

func (s *Service) checkPolicy(id int64) (policy *model.Policy, err error) {
	if id != 1 {
		err = ecode.PLayerPolicyNotExist
		return
	}
	policy = s.c.Policy
	if time.Now().Unix() < policy.StartTime.Unix() {
		err = ecode.PLayerPolicyNotStart
		return
	}
	if time.Now().Unix() > policy.EndTime.Unix() {
		err = ecode.PLayerPolicyEnded
		return
	}
	return
}

// 用户随机-尾号 策略方法
func (s *Service) userPolicy(mid int64, policy *model.Policy) (res *model.Pitem, err error) {
	var itemMap = make(map[string]*model.Pitem, len(s.c.Pitem))
	for _, item := range s.c.Pitem {
		item.Ver = policy.MtimeTime.Unix()
		itemMap[item.ExtData] = item
	}
	if mid > 0 {
		utail := int(mid % 100)
		for _, item := range itemMap {
			var (
				begin       int
				end         int
				beginAndEnd []string
			)
			if item.ExtData == "default" {
				continue
			}
			beginAndEnd = strings.Split(item.ExtData, "-")
			if len(beginAndEnd) != 2 {
				log.Error("item.ExtData error")
				return
			}
			if begin, err = strconv.Atoi(beginAndEnd[0]); err != nil {
				log.Error("item.ExtData error")
				return
			}
			if end, err = strconv.Atoi(beginAndEnd[1]); err != nil {
				log.Error("item.ExtData error")
				return
			}
			if utail >= begin && utail < end {
				res = item
				return
			}
		}
	} else {
		res = itemMap["default"]
	}
	return
}
