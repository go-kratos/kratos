package service

import (
	"context"

	"go-common/app/service/main/account/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Relation get user relation.
func (s *Service) Relation(c context.Context, mid, fid int64) (r *model.Relation, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgRelation{Mid: mid, Fid: fid, RealIP: ip}
	res, err := s.relRPC.Relation(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service relation")
		return
	}
	r = &model.Relation{}
	r.Following = res.Following()
	return
}

// Relations get user relations.
func (s *Service) Relations(c context.Context, mid int64, fids []int64) (rs map[int64]*model.Relation, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}
	res, err := s.relRPC.Relations(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service relations")
		return
	}
	rs = make(map[int64]*model.Relation, len(res))
	for _, fid := range fids {
		if f, ok := res[fid]; ok {
			rs[fid] = &model.Relation{Following: f.Following()}
		} else {
			rs[fid] = &model.Relation{Following: false}
		}
	}
	return
}

// RichRelations2 get mid relations of fids.
func (s *Service) RichRelations2(c context.Context, mid int64, fids []int64) (rs map[int64]int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}
	res, err := s.relRPC.Relations(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service rich relation")
		return
	}
	rs = make(map[int64]int, len(res))
	for mid, f := range res {
		rs[mid] = int(f.Attribute)
	}
	return
}

// Blacks get user balcks.
func (s *Service) Blacks(c context.Context, mid int64) (rs map[int64]struct{}, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgMid{Mid: mid, RealIP: ip}
	res, err := s.relRPC.Blacks(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service blacks")
		return
	}
	rs = make(map[int64]struct{}, len(res))
	for _, r := range res {
		rs[r.Mid] = struct{}{}
	}
	return
}

// Attentions get all attentions list ,include followings and whispers.
func (s *Service) Attentions(c context.Context, mid int64) (mids []int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgMid{Mid: mid, RealIP: ip}
	rls, err := s.relRPC.Followings(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service attentions")
		return
	}
	for _, rl := range rls {
		mids = append(mids, rl.Mid)
	}
	whs, err := s.relRPC.Whispers(c, arg)
	if err != nil {
		err = errors.Wrap(err, "service whispers")
		return
	}
	for _, wh := range whs {
		mids = append(mids, wh.Mid)
	}
	return
}
