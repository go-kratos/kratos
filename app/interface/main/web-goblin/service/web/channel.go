package web

import (
	"context"

	tagmdl "go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/web-goblin/model/web"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_chRqCnt     = 40
	_chDisplayID = 1
	_chTypeArc   = 3
	_chFrom      = 1
)

var _emptyArcs = make([]*api.Arc, 0)

// Channel .
func (s *Service) Channel(c context.Context, id, mid int64, buvid string) (channel *web.Channel, err error) {
	var (
		aids   []int64
		arcs   map[int64]*api.Arc
		tagErr error
	)
	ip := metadata.String(c, metadata.RemoteIP)
	channel = new(web.Channel)
	if cards, ok := s.channelCards[id]; ok {
		for _, card := range cards {
			aids = append(aids, card.Value)
		}
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		arg := &tagmdl.ArgChannelResource{
			Tid:        id,
			Mid:        mid,
			RequestCNT: int32(_chRqCnt),
			DisplayID:  _chDisplayID,
			Type:       _chTypeArc,
			Buvid:      buvid,
			From:       _chFrom,
			RealIP:     ip,
		}
		if channelResource, chErr := s.tag.ChannelResources(errCtx, arg); chErr != nil {
			log.Error("Channel s.tag.Resources error(%v)", chErr)
		} else if channelResource != nil {
			aids = append(aids, channelResource.Oids...)
		}
		return nil
	})
	group.Go(func() error {
		if channel.Tag, tagErr = s.tag.InfoByID(errCtx, &tagmdl.ArgID{ID: id, Mid: mid}); tagErr != nil {
			log.Error("Channel s.tag.InfoByID(%d, %d) error(%v)", id, mid, err)
			return tagErr
		}
		return nil
	})
	if err = group.Wait(); err != nil {
		return
	}
	if len(aids) == 0 {
		channel.Archives = _emptyArcs
		return
	}
	if arcs, err = s.arc.Archives3(c, &archive.ArgAids2{Aids: aids, RealIP: ip}); err != nil {
		log.Error("Channel s.arc.Archives3(%v) error(%v)", aids, err)
		err = nil
		channel.Archives = _emptyArcs
		return
	}
	for _, aid := range aids {
		if arc, ok := arcs[aid]; ok && arc.IsNormal() {
			channel.Archives = append(channel.Archives, arc)
		}
	}
	if len(channel.Archives) > _chRqCnt {
		channel.Archives = channel.Archives[:_chRqCnt]
	}
	return
}
