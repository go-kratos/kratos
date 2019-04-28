package grpc

import (
	"context"

	v1 "go-common/app/service/main/tag/api"
	"go-common/app/service/main/tag/model"
	"go-common/library/ecode"
)

func (s *grpcServer) TagMap(c context.Context, arg *v1.TagMapByIDReq) (res *v1.TagMapByIDReply, err error) {
	res = &v1.TagMapByIDReply{
		Tags: make(map[int64]*v1.Tag),
	}
	var (
		tagMap   map[int64]*model.Tag
		countMap map[int64]*model.Count
	)
	if tagMap, err = s.svr.InfoMap(c, arg.Mid, arg.Tids); err != nil {
		return
	}
	if countMap, err = s.svr.Counts(c, arg.Tids); err != nil {
		return
	}
	for _, tag := range tagMap {
		var (
			sub  int64
			bind int64
		)
		if k, ok := countMap[tag.ID]; ok {
			sub = k.Sub
			bind = k.Bind
		}
		t := &v1.Tag{
			Id:           tag.ID,
			Name:         tag.Name,
			Cover:        tag.Cover,
			HeadCover:    tag.HeadCover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			Type:         tag.Type,
			Verify:       tag.Verify,
			Attr:         tag.Attr,
			Attention:    tag.Attention,
			State:        tag.State,
			Sub:          sub,
			Bind:         bind,
			Ctime:        tag.CTime,
			Mtime:        tag.MTime,
		}
		res.Tags[t.Id] = t
	}
	return
}

func (s *grpcServer) ResTag(c context.Context, arg *v1.ResTagReq) (res *v1.ResTagReply, err error) {
	var (
		resTagMap map[int64]*model.Resource
		req       = &model.ArgResTag{
			Oid:  arg.Oid,
			Type: arg.Type,
		}
	)
	res = &v1.ResTagReply{
		Resource: make([]*v1.Resource, 0),
	}
	if resTagMap, err = s.svr.ResTagMap(c, req); err != nil {
		return
	}
	for _, v := range resTagMap {
		t := &v1.Resource{
			Id:    v.ID,
			Oid:   v.Oid,
			Tid:   v.Tid,
			Type:  v.Type,
			Mid:   v.Mid,
			Attr:  v.Attr,
			Role:  v.Role,
			Like:  v.Like,
			Hate:  v.Hate,
			State: v.State,
			Ctime: v.CTime,
			Mtime: v.MTime,
		}
		res.Resource = append(res.Resource, t)
	}
	return
}

func (s *grpcServer) ResTagMap(c context.Context, arg *v1.ResTagReq) (res *v1.ResTagMapReply, err error) {
	var (
		resTagMap map[int64]*model.Resource
		req       = &model.ArgResTag{
			Oid:  arg.Oid,
			Type: arg.Type,
		}
	)
	res = &v1.ResTagMapReply{
		Resource: make(map[int64]*v1.Resource),
	}
	if resTagMap, err = s.svr.ResTagMap(c, req); err != nil {
		return
	}
	for _, v := range resTagMap {
		t := &v1.Resource{
			Id:    v.ID,
			Oid:   v.Oid,
			Tid:   v.Tid,
			Type:  v.Type,
			Mid:   v.Mid,
			Attr:  v.Attr,
			Role:  v.Role,
			Like:  v.Like,
			Hate:  v.Hate,
			State: v.State,
			Ctime: v.CTime,
			Mtime: v.MTime,
		}
		res.Resource[v.Tid] = t
	}
	return
}

func (s *grpcServer) ResTags(c context.Context, arg *v1.ResTagsReq) (res *v1.ResTagsReply, err error) {
	var (
		req = &model.ArgMutiResTag{
			Oids: arg.Oids,
			Type: arg.Type,
		}
		rsMap map[int64][]*model.Resource
	)
	res = &v1.ResTagsReply{
		Resource: make(map[int64]*v1.ResTag, len(arg.Oids)),
	}
	if rsMap, err = s.svr.MutiResTagMap(c, req); err != nil {
		return
	}
	for oid, rs := range rsMap {
		resTag := &v1.ResTag{
			Resource: make([]*v1.Resource, 0, len(rs)),
		}
		for _, v := range rs {
			t := &v1.Resource{
				Id:    v.ID,
				Oid:   v.Oid,
				Tid:   v.Tid,
				Type:  v.Type,
				Mid:   v.Mid,
				Attr:  v.Attr,
				Role:  v.Role,
				Like:  v.Like,
				Hate:  v.Hate,
				State: v.State,
				Ctime: v.CTime,
				Mtime: v.MTime,
			}
			resTag.Resource = append(resTag.Resource, t)
		}
		res.Resource[oid] = resTag
	}
	return
}

func (s *grpcServer) Tag(c context.Context, arg *v1.TagReq) (res *v1.TagReply, err error) {
	res = &v1.TagReply{
		Tag: new(v1.Tag),
	}
	if arg.Tid <= 0 {
		err = ecode.RequestErr
		return
	}
	var (
		tag   *model.Tag
		count *model.Count
	)
	if tag, err = s.svr.Info(c, arg.Mid, arg.Tid); err != nil {
		return
	}
	if tag == nil {
		err = ecode.TagNotExist
		return
	}
	if count, err = s.svr.Count(c, arg.Tid); err != nil {
		return
	}
	res.Tag = &v1.Tag{
		Id:           tag.ID,
		Name:         tag.Name,
		Cover:        tag.Cover,
		HeadCover:    tag.HeadCover,
		Content:      tag.Content,
		ShortContent: tag.ShortContent,
		Type:         tag.Type,
		Verify:       tag.Verify,
		Attr:         tag.Attr,
		Attention:    tag.Attention,
		State:        tag.State,
		Ctime:        tag.CTime,
		Mtime:        tag.MTime,
	}
	if count != nil {
		res.Tag.Sub = count.Sub
		res.Tag.Bind = count.Bind
	}
	return
}

func (s *grpcServer) TagByName(c context.Context, arg *v1.TagByNameReq) (res *v1.TagReply, err error) {
	res = &v1.TagReply{
		Tag: new(v1.Tag),
	}
	if arg.Tname == "" {
		err = ecode.RequestErr
		return
	}
	var (
		tag   *model.Tag
		count *model.Count
	)
	if tag, err = s.svr.InfoByName(c, arg.Mid, arg.Tname); err != nil {
		return
	}
	if tag == nil {
		err = ecode.TagNotExist
		return
	}
	if count, err = s.svr.Count(c, tag.ID); err != nil {
		return
	}
	res.Tag = &v1.Tag{
		Id:           tag.ID,
		Name:         tag.Name,
		Cover:        tag.Cover,
		HeadCover:    tag.HeadCover,
		Content:      tag.Content,
		ShortContent: tag.ShortContent,
		Type:         tag.Type,
		Verify:       tag.Verify,
		Attr:         tag.Attr,
		Attention:    tag.Attention,
		State:        tag.State,
		Ctime:        tag.CTime,
		Mtime:        tag.MTime,
	}
	if count != nil {
		res.Tag.Sub = count.Sub
		res.Tag.Bind = count.Bind
	}
	return
}

// Like user like oid-type-tid relation.
func (s *grpcServer) Like(c context.Context, arg *v1.ResTagActionReq) (res *v1.ResTagActionReply, err error) {
	if arg.Type != model.ResTypeArchive {
		err = ecode.RequestErr
		return
	}
	return &v1.ResTagActionReply{}, s.svr.Like(c, arg.Mid, arg.Oid, arg.Tid, arg.Type, "")
}

// Hate user hate oid-type-tid relation.
func (s *grpcServer) Hate(c context.Context, arg *v1.ResTagActionReq) (res *v1.ResTagActionReply, err error) {
	if arg.Type != model.ResTypeArchive {
		err = ecode.RequestErr
		return
	}
	return &v1.ResTagActionReply{}, s.svr.Hate(c, arg.Mid, arg.Oid, arg.Tid, arg.Type, "")
}

// ResTagActionMap resource tag action map by given tids.
func (s *grpcServer) ResTagActionMap(c context.Context, arg *v1.ResTagActionMapReq) (res *v1.ResTagActionMapReply, err error) {
	if arg.Type != model.ResTypeArchive {
		err = ecode.RequestErr
		return
	}
	res = &v1.ResTagActionMapReply{}
	res.ActionMap, err = s.svr.ActionMap(c, arg.Mid, arg.Oid, arg.Type, arg.Tids)
	return
}

// UpBind res-tag upbind.
func (s *grpcServer) UpBind(c context.Context, arg *v1.UpBindReq) (res *v1.UpBindReply, err error) {
	return &v1.UpBindReply{}, s.svr.PlatformUpBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, "")
}

// AdminBind res-tag admin bind.
func (s *grpcServer) AdminBind(c context.Context, arg *v1.AdminBindReq) (res *v1.AdminBindReply, err error) {
	return &v1.AdminBindReply{}, s.svr.PlatformAdminBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, "")
}

// DefaultUpBind res-tag default upbind.
func (s *grpcServer) DefaultUpBind(c context.Context, arg *v1.UpBindReq) (res *v1.UpBindReply, err error) {
	return &v1.UpBindReply{}, s.svr.DefaultUpBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, "")
}

// DefaultAdminBind res-tag default admin bind.
func (s *grpcServer) DefaultAdminBind(c context.Context, arg *v1.AdminBindReq) (res *v1.AdminBindReply, err error) {
	return &v1.AdminBindReply{}, s.svr.DefaultAdminBind(c, arg.Oid, arg.Mid, arg.Tids, arg.Type, "")
}
