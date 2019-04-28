package dao

import (
	"context"

	"go-common/app/interface/main/tag/model"
	taGrpcModel "go-common/app/service/main/tag/api"
	"go-common/library/log"
)

// AddReport add a report to tag-serivce.
func (d *Dao) AddReport(c context.Context, oid, tid, mid int64, typ, partID, reason, score int32) (err error) {
	arg := &taGrpcModel.AddReportReq{
		Oid:      oid,
		Mid:      mid,
		Tid:      tid,
		Type:     typ,
		PartId:   partID,
		ReasonId: reason,
		Score:    score,
	}
	if _, err = d.tagRPC.AddReport(c, arg); err != nil {
		log.Error("d.tagRPC.AddReport()Arg:%+v, error(%v)", arg, err)
	}
	return
}

// Tag get tag info by id.
func (d *Dao) Tag(c context.Context, tid int64, mid int64) (*taGrpcModel.Tag, error) {
	arg := &taGrpcModel.TagReq{
		Mid: mid,
		Tid: tid,
	}
	reply, err := d.tagRPC.Tag(c, arg)
	if err != nil {
		log.Error("d.dao.Tag(%v) error: %v", arg, err)
		return nil, err
	}
	return reply.Tag, err
}

// TagByName get tag info by name.
func (d *Dao) TagByName(c context.Context, mid int64, tname string) (*taGrpcModel.Tag, error) {
	arg := &taGrpcModel.TagByNameReq{
		Mid:   mid,
		Tname: tname,
	}
	reply, err := d.tagRPC.TagByName(c, arg)
	if err != nil {
		log.Error("d.dao.TagByName(%v) error: %v", arg, err)
		return nil, err
	}
	return reply.Tag, err
}

// TagMap get tag info map.
func (d *Dao) TagMap(c context.Context, tids []int64, mid int64) (res map[int64]*taGrpcModel.Tag, err error) {
	var (
		reply   *taGrpcModel.TagMapByIDReply
		n       = model.MaxTagNum
		tidsMap = make(map[int64]struct{}, len(tids))
		newTids = make([]int64, 0, len(tids))
	)
	for _, tid := range tids {
		if _, ok := tidsMap[tid]; !ok {
			tidsMap[tid] = struct{}{}
			newTids = append(newTids, tid)
		}
	}
	res = make(map[int64]*taGrpcModel.Tag, len(newTids))
	for len(newTids) > 0 {
		if n > len(newTids) {
			n = len(newTids)
		}
		arg := &taGrpcModel.TagMapByIDReq{
			Mid:  mid,
			Tids: newTids[:n],
		}
		newTids = newTids[n:]
		if reply, err = d.tagRPC.TagMap(c, arg); err != nil {
			log.Error("d.dao.TagMapByID(%v) error: %v", arg, err)
			return
		}
		for k, v := range reply.Tags {
			res[k] = v
		}
	}
	return
}

// ResTag res tag.
func (d *Dao) ResTag(c context.Context, oid int64, tp int32) (res []*taGrpcModel.Resource, err error) {
	var (
		reply *taGrpcModel.ResTagReply
		arg   = &taGrpcModel.ResTagReq{
			Oid:  oid,
			Type: tp,
		}
	)
	if reply, err = d.tagRPC.ResTag(c, arg); err != nil {
		log.Error("d.dao.ResTags(%d,%d) error: %v", oid, tp, err)
		return
	}
	return reply.Resource, nil
}

// ResTagMap res tag map.
func (d *Dao) ResTagMap(c context.Context, oid int64, tp int32) (res map[int64]*taGrpcModel.Resource, err error) {
	var (
		reply *taGrpcModel.ResTagMapReply
		arg   = &taGrpcModel.ResTagReq{
			Oid:  oid,
			Type: tp,
		}
	)
	if reply, err = d.tagRPC.ResTagMap(c, arg); err != nil {
		log.Error("d.dao.ResTags(%d,%d) error: %v", oid, tp, err)
		return
	}
	return reply.Resource, nil
}

// ResTags res tags.
func (d *Dao) ResTags(c context.Context, oids []int64, tp int32) (res map[int64][]*taGrpcModel.Resource, err error) {
	var (
		reply *taGrpcModel.ResTagsReply
		n     = model.ResMaxNum
	)
	res = make(map[int64][]*taGrpcModel.Resource, len(oids))
	for len(oids) > 0 {
		if n > len(oids) {
			n = len(oids)
		}
		arg := &taGrpcModel.ResTagsReq{
			Oids: oids[:n],
			Type: tp,
		}
		oids = oids[n:]
		if reply, err = d.tagRPC.ResTags(c, arg); err != nil {
			log.Error("d.dao.ResTags(%v) error: %v", arg, err)
			return
		}
		for k, values := range reply.Resource {
			resTag := make([]*taGrpcModel.Resource, 0, len(values.Resource))
			resTag = append(resTag, values.Resource...)
			res[k] = resTag
		}
	}
	return
}
