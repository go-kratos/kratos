package service

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/log"
)

// Register 注册extension信息，包括关联话题等
func (s *Service) Register(ctx context.Context, req *api.VideoExtension) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	svid := req.Svid
	// 0. check
	if svid == 0 {
		log.Warnw(ctx, "log", "svid is 0")
		return
	}
	if len(req.Extension) == 0 {
		return
	}
	var extension api.Extension
	err = json.Unmarshal([]byte(req.Extension), &extension)
	if err != nil {
		log.Errorw(ctx, "log", "unmarshal extension fail", "req", req)
		return
	}

	// 2. 获取新的title_extra
	extension.TitleExtra, err = s.registerTopic(ctx, svid, extension.TitleExtra)
	if err != nil {
		log.Warnw(ctx, "log", "regist topic fail")
		return
	}

	// 3. 存入extension中
	num, err := s.dao.InsertExtension(ctx, svid, model.ExtensionTypeTitleExtra, &extension)
	if err != nil {
		log.Warnw(ctx, "log", "insert extension fail", "extension", extension, "svid", svid)
		return
	}
	// 插入无效主要是因为已经存在，所以这里默认为成功，但是打error日志
	if num == 0 {
		log.Errorw(ctx, "log", "insert extension fail due to svid already exists", "svid", svid)
	}
	return
}

// ListExtension 获取视频的extension信息
func (s *Service) ListExtension(ctx context.Context, req *api.ListExtensionReq) (res *api.ListExtensionReply, err error) {
	res = new(api.ListExtensionReply)
	videoExtensions, err := s.dao.VideoExtension(ctx, req.Svids)
	if err != nil {
		log.Warnw(ctx, "log", "get video extension fail", "svids", req.Svids)
		return
	}

	for _, extension := range videoExtensions {
		res.List = append(res.List, extension)
	}

	return
}
