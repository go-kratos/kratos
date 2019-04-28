package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) parseResTagMessage(msg *model.ResTagMessage) (res *model.ResTag, err error) {
	business := s.businessCache.Load().(map[string]*model.Business)
	b, ok := business[msg.Type]
	if !ok || b == nil || (msg.Appkey != b.Appkey) {
		log.Error("parseResTagMessage error: type or appkey not right. message:%+v", msg)
		err = ecode.RequestErr
		return
	}
	if len(msg.Tids) <= 0 && len(msg.TNames) <= 0 {
		log.Error("parseResTagMessage error: tids or tnames null slice. message:%+v", msg)
		err = ecode.RequestErr
		return
	}
	if msg.Oid <= 0 {
		log.Error("parseResTagMessage error: msg.Oid <= 0, message:%+v", msg)
		err = ecode.RequestErr
		return
	}
	if msg.Mid <= 0 && msg.Role != model.ResTagRoleAdmin {
		log.Error("parseResTagMessage error: role!=admin, but mid=0, message:%+v", msg)
		err = ecode.RequestErr
		return
	}
	var role int32
	switch msg.Role {
	case model.ResTagRoleUp:
		role = model.RoleUp
	case model.ResTagRoleUser:
		role = model.RoleUser
	case model.ResTagRoleAdmin:
		role = model.RoleAdmin
		if msg.Mid <= 0 {
			msg.Mid = 0
		}
	default:
		log.Error("parseResTagMessage error: Wrong role, message:%+v", msg)
		err = ecode.RequestErr
		return
	}
	res = &model.ResTag{
		Oid:   msg.Oid,
		Type:  b.Type,
		Mid:   msg.Mid,
		Role:  role,
		MTime: msg.MTime,
	}
	return
}

func (s *Service) resTagActionConsumeproc() {
	var (
		err error
		c   = context.Background()
	)
	for {
		msg, ok := <-s.tagSub.Messages()
		if !ok {
			log.Error("s.resTagActionConsumeproc() consumer exit")
			return
		}
		log.Warn("res_tag_action accept: partition:%d,offset:%d,key:%s,value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
		if err = msg.Commit(); err != nil {
			log.Error("s.resTagActionConsumeproc() commit offset(%v) error(%v)", msg, err)
			continue
		}
		m := &model.ResTagMessage{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("s.resTagActionConsumeproc().json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.MTime <= 0 {
			m.MTime = msg.Timestamp
		}
		if err = s.resTagHandle(c, m); err != nil {
			log.Error("res_tag_action have not dealed, value:%v, error(%v)", m, err)
			continue
		}
		log.Warn("res_tag_action have dealed, value:%v", m)
	}
}

func (s *Service) resTagHandle(c context.Context, msg *model.ResTagMessage) (err error) {
	var (
		missedNames []string
		missedTags  []*model.Tag
		tags        []*model.Tag
		rt          *model.ResTag
	)
	if rt, err = s.parseResTagMessage(msg); err != nil {
		return
	}
	if len(msg.TNames) != 0 {
		tnames := make([]string, 0, len(msg.TNames))
		for _, name := range msg.TNames {
			if name, err = s.checkName(name); err == nil {
				tnames = append(tnames, name)
			}
		}
		if len(tnames) <= 0 {
			log.Error("resTagHandle msg tags is null, message:%+v", msg)
			return ecode.RequestErr
		}
		tags, missedNames, err = s.dao.TagByNames(c, tnames)
		if err != nil {
			return
		}
		if len(missedNames) > 0 {
			if missedTags, err = s.createTags(c, missedNames); err == nil {
				tags = append(tags, missedTags...)
			}
		}
	} else if len(msg.Tids) != 0 {
		tags, err = s.dao.Tags(c, msg.Tids)
	}
	if err != nil {
		return
	}
	for _, t := range tags {
		if t.ID <= 0 || t.State != model.TagStateNormal {
			continue
		}
		rt.Tids = append(rt.Tids, t.ID)
	}
	if len(rt.Tids) == 0 {
		log.Error("resTagHandle, tids is null, msg:%+v", msg)
		return
	}
	for retry := 0; retry < model.MaxRetryTimes; retry++ {
		switch msg.Action {
		case model.ResTagBind:
			rt.State = model.ResTagStateNormal
			err = s.resTagBind(c, rt)
		case model.ResTagDelete:
			rt.State = model.ResTagStateDelete
			err = s.resTagDelete(c, rt)
		default:
		}
		if err == nil {
			return
		}
		time.Sleep(time.Millisecond * 50)
	}
	return
}
