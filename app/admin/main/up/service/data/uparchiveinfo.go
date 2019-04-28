package data

import (
	"context"
	"go-common/app/admin/main/up/dao/data"
	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/util"
	"go-common/library/log"
	"sort"
	"strconv"
)

//GetUpArchiveInfo get up archive info
func (s *Service) GetUpArchiveInfo(c context.Context, arg *datamodel.GetUpArchiveInfoArg) (result *datamodel.GetUpArchiveInfoResult, err error) {
	if arg.DataType == 0 {
		arg.DataType = datamodel.DataType30Day
	}
	result = &datamodel.GetUpArchiveInfoResult{}
	var mids = util.ExplodeInt64(arg.Mids, ",")
	var length = len(mids)
	if length == 0 {
		log.Info("no mids specified")
		return
	} else if length > 100 {
		// 每次最多100个
		mids = mids[0:100]
	}

	dataMap, err := s.data.UpArchiveInfo(c, mids, data.UpArchiveDataType(arg.DataType))
	if err != nil {
		log.Error("get up archive info fail, err=%v, arg=%+v", err, arg)
		return
	}

	for mid, v := range dataMap {
		(*result)[mid] = v
	}
	log.Info("get up archive info ok, type=%d", arg.DataType)
	return
}

//GetUpArchiveTagInfo get up archive tag info
func (s *Service) GetUpArchiveTagInfo(c context.Context, arg *datamodel.GetUpArchiveTagInfoArg) (result []*datamodel.ViewerTagData, err error) {

	tagData, err := s.data.UpArchiveTagInfo(c, arg.Mid)
	if err != nil {
		log.Error("get up archive tag fail, err=%v", err)
		return
	}

	var tagResultMap = make(map[int64]*datamodel.ViewerTagData)
	var tagIds []int64
	for idxstr, tid := range tagData.TagMap {
		tagIds = append(tagIds, tid)
		var idx, _ = strconv.Atoi(idxstr)
		var tag = &datamodel.ViewerTagData{
			Idx:   idx,
			TagID: int(tid),
		}
		tagResultMap[tid] = tag
	}
	var tagMeta = s.GetTags(c, tagIds...)
	for tid, meta := range tagMeta {
		tag, ok := tagResultMap[tid]
		if !ok {
			continue
		}
		tag.Name = meta.TagName
	}

	for _, tag := range tagResultMap {
		result = append(result, tag)
	}

	if len(result) > 1 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Idx < result[j].Idx
		})
	}

	log.Info("get up archive tag info ok, mid=%d", arg.Mid)
	return
}

//GetUpArchiveTypeInfo get type info
func (s *Service) GetUpArchiveTypeInfo(c context.Context, arg *datamodel.GetUpArchiveTypeInfoArg) (result *datamodel.UpArchiveTypeData, err error) {
	res, err := s.data.UpArchiveTypeInfo(c, arg.Mid)
	result = &res
	if err != nil {
		log.Error("fail to get up type, err=%v", err)
		return
	}

	log.Info("get up archive type info ok, mid=%d", arg.Mid)
	return
}
