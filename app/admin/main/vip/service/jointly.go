package service

import (
	"context"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
)

// JointlysByState jointlys by state.
func (s *Service) JointlysByState(c context.Context, state int8) (res []*model.Jointly, err error) {
	now := time.Now().Unix()
	res, err = s.dao.JointlysByState(c, state, time.Now().Unix())
	for _, v := range res {
		switch {
		case v.StartTime > now:
			v.State = model.WillEffect
		case v.StartTime < now && v.EndTime > now:
			v.State = model.Effect
		case v.EndTime < now:
			v.State = model.LoseEffect
		}
	}
	return
}

// AddJointly add jointly.
func (s *Service) AddJointly(c context.Context, j *model.ArgAddJointly) (err error) {
	if j.StartTime >= j.EndTime {
		err = ecode.VipStartTimeErr
		return
	}
	if len(j.Title) > _maxTitleLen {
		err = ecode.VipTitleTooLongErr
		return
	}
	if len(j.Content) > _maxContentLen {
		err = ecode.VipContentTooLongErr
		return
	}
	_, err = s.dao.AddJointly(c, &model.Jointly{
		Title:     j.Title,
		Content:   j.Content,
		Operator:  j.Operator,
		StartTime: j.StartTime,
		EndTime:   j.EndTime,
		Link:      j.Link,
		IsHot:     j.IsHot,
	})
	return
}

// ModifyJointly modify jointly.
func (s *Service) ModifyJointly(c context.Context, j *model.ArgModifyJointly) (err error) {
	if j.StartTime >= j.EndTime {
		err = ecode.VipStartTimeErr
		return
	}
	if len(j.Title) > _maxTitleLen {
		err = ecode.VipTitleTooLongErr
		return
	}
	if len(j.Content) > _maxContentLen {
		err = ecode.VipContentTooLongErr
		return
	}
	_, err = s.dao.UpdateJointly(c, &model.Jointly{
		Title:     j.Title,
		Content:   j.Content,
		Operator:  j.Operator,
		Link:      j.Link,
		IsHot:     j.IsHot,
		ID:        j.ID,
		StartTime: j.StartTime,
		EndTime:   j.EndTime,
	})
	return
}

// DeleteJointly delete jointly .
func (s *Service) DeleteJointly(c context.Context, id int64) (err error) {
	_, err = s.dao.DeleteJointly(c, id)
	return
}
