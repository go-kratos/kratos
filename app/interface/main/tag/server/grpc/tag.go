package grpc

import (
	"context"

	pb "go-common/app/interface/main/tag/api"
	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
)

// Tag get a tag info by id.
func (s *grpcServer) Tag(c context.Context, arg *pb.TagReq) (res *pb.TagReply, err error) {
	res = &pb.TagReply{
		Tag: new(pb.Tag),
	}
	if arg.Tid <= 0 {
		err = ecode.RequestErr
		return
	}
	tag, err := s.svr.InfoByID(c, arg.Mid, arg.Tid)
	if err != nil {
		return
	}
	res.Tag = &pb.Tag{
		Id:           tag.ID,
		Name:         tag.Name,
		Type:         int32(tag.Type),
		Cover:        tag.Cover,
		HeadCover:    tag.HeadCover,
		Content:      tag.Content,
		ShortContent: tag.ShortContent,
		Attention:    int32(tag.IsAtten),
		Sub:          int64(tag.Count.Atten),
		Bind:         int64(tag.Count.Use),
		Liked:        int32(tag.Liked),
		Hated:        int32(tag.Hated),
		Likes:        tag.Likes,
		Hates:        tag.Hates,
		Ctime:        tag.CTime,
		State:        int32(tag.State),
		Mtime:        tag.MTime,
	}
	return
}

// TagByName get a tag info by name.
func (s *grpcServer) TagByName(c context.Context, arg *pb.TagByNameReq) (res *pb.TagReply, err error) {
	res = &pb.TagReply{
		Tag: new(pb.Tag),
	}
	if arg.Tname, err = s.svr.CheckName(arg.Tname); err != nil {
		return
	}
	tag, err := s.svr.InfoByName(c, arg.Mid, arg.Tname)
	if err != nil {
		return
	}
	res.Tag = &pb.Tag{
		Id:           tag.ID,
		Name:         tag.Name,
		Type:         int32(tag.Type),
		Cover:        tag.Cover,
		HeadCover:    tag.HeadCover,
		Content:      tag.Content,
		ShortContent: tag.ShortContent,
		Attention:    int32(tag.IsAtten),
		Sub:          int64(tag.Count.Atten),
		Bind:         int64(tag.Count.Use),
		Liked:        int32(tag.Liked),
		Hated:        int32(tag.Hated),
		Likes:        tag.Likes,
		Hates:        tag.Hates,
		Ctime:        tag.CTime,
		State:        int32(tag.State),
		Mtime:        tag.MTime,
	}
	return
}

// Tags get tags info by ids.
func (s *grpcServer) Tags(c context.Context, arg *pb.TagsReq) (res *pb.TagsReply, err error) {
	res = &pb.TagsReply{
		Tags: make(map[int64]*pb.Tag, len(arg.Tids)),
	}
	if len(arg.Tids) <= 0 || len(arg.Tids) > model.MaxTagNum {
		err = ecode.RequestErr
		return
	}
	tags, err := s.svr.MinfoByIDs(c, arg.Mid, arg.Tids)
	if err != nil {
		return
	}
	for _, tag := range tags {
		res.Tags[tag.ID] = &pb.Tag{
			Id:           tag.ID,
			Name:         tag.Name,
			Type:         int32(tag.Type),
			Cover:        tag.Cover,
			HeadCover:    tag.HeadCover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			Attention:    int32(tag.IsAtten),
			Sub:          int64(tag.Count.Atten),
			Bind:         int64(tag.Count.Use),
			Liked:        int32(tag.Liked),
			Hated:        int32(tag.Hated),
			Likes:        tag.Likes,
			Hates:        tag.Hates,
			Ctime:        tag.CTime,
			State:        int32(tag.State),
			Mtime:        tag.MTime,
		}
	}
	return
}

// TagByNames get tags info by names.
func (s *grpcServer) TagByNames(c context.Context, arg *pb.TagByNamesReq) (res *pb.TagsReply, err error) {
	res = &pb.TagsReply{
		Tags: make(map[int64]*pb.Tag, len(arg.Tnames)),
	}
	tnames := make([]string, 0, len(arg.Tnames))
	for _, name := range arg.Tnames {
		if name, err = s.svr.CheckName(name); err != nil {
			continue
		}
		tnames = append(tnames, name)
	}
	if len(tnames) <= 0 || len(tnames) > model.MaxTagNum {
		err = ecode.RequestErr
		return
	}
	tags, err := s.svr.MinfoByNames(c, arg.Mid, tnames)
	if err != nil {
		return
	}
	for _, tag := range tags {
		res.Tags[tag.ID] = &pb.Tag{
			Id:           tag.ID,
			Name:         tag.Name,
			Type:         int32(tag.Type),
			Cover:        tag.Cover,
			HeadCover:    tag.HeadCover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			Attention:    int32(tag.IsAtten),
			Sub:          int64(tag.Count.Atten),
			Bind:         int64(tag.Count.Use),
			Liked:        int32(tag.Liked),
			Hated:        int32(tag.Hated),
			Likes:        tag.Likes,
			Hates:        tag.Hates,
			Ctime:        tag.CTime,
			State:        int32(tag.State),
			Mtime:        tag.MTime,
		}
	}
	return
}
