package service

import (
	"context"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) getVideo(ctx context.Context, aid int64, cid int64) (video *model.Video, err error) {
	if video, err = s.dao.Video(ctx, aid, cid); err != nil {
		log.Error("getVideo s.arc.Video error(%v) aid(%d) cid(%d)", err, aid, cid)
		return
	}
	if video == nil {
		err = ecode.NothingFound
		log.Error("getVideo s.arc.Video error(%v) aid(%d) cid(%d)", err, aid, cid)
		return
	}

	video.AttributeList = model.AttributeList(video.Attribute)
	if video.XcodeState >= model.VideoXcodeHDFinish {
		video.Encoding = 1
	}
	return
}

func (s *Service) getVideoOperInfo(ctx context.Context, vid int64) (list []*model.VideoOperInfo, err error) {
	var (
		vopers []*model.VOper
		uids   []int64
		users  map[int64]*model.UserDepart
	)

	list = []*model.VideoOperInfo{}
	if vopers, uids, err = s.dao.VideoOpers(ctx, vid); err != nil {
		log.Error("getVideoOperInfo s.dao.VideoOpers(%d) error(%v)", vid, err)
		return
	}
	if users, err = s.dao.GetUsernameAndDepartment(ctx, uids); err != nil {
		log.Error("getVideoOperInfo s.dao.GetUsernameAndRol(%d) error(%v) uids(%v)", vid, err, uids)
		return
	}
	for _, op := range vopers {
		u := users[op.UID]
		if u == nil {
			u = &model.UserDepart{UID: op.UID}
		}
		info := &model.VideoOperInfo{
			VOper:      *op,
			UserDepart: *u,
		}
		list = append(list, info)
	}
	return
}
