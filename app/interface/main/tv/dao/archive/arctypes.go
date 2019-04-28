package archive

import (
	"context"

	"go-common/app/interface/main/tv/model"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// call ArcRPC for types data
func (d *Dao) loadTypes(ctx context.Context) {
	var (
		res    = make(map[int32]*arcwar.Tp)
		resRel = make(map[int32][]*arcwar.Tp)
		reply  *arcwar.TypesReply
		err    error
	)
	if reply, err = d.arcClient.Types(ctx, &arcwar.NoArgRequest{}); err != nil {
		log.Error("arcRPC loadType Error %v", err)
		return
	}
	res = reply.Types
	log.Info("Reload Types Data! Len: %d", len(res))
	for _, v := range res {
		if v.Pid != 0 {
			if _, ok := resRel[v.Pid]; !ok {
				resRel[v.Pid] = []*arcwar.Tp{}
			}
			resRel[v.Pid] = append(resRel[v.Pid], v)
		}
	}
	if len(res) > 0 {
		d.arcTypes = res
	}
	if len(resRel) > 0 {
		d.arcTypesRel = resRel
	}
}

// GetPTypeName get first level of types name
func (d *Dao) GetPTypeName(typeID int32) (firstName string, secondName string) {
	var (
		second, first *arcwar.Tp
		ok            bool
	)
	if second, ok = d.arcTypes[typeID]; !ok {
		log.Error("can't find type for ID: %d ", typeID)
		return
	}
	secondName = second.Name
	if first, ok = d.arcTypes[second.Pid]; !ok {
		log.Error("can't find type for ID: %d, second Info: %v", second, second.Pid)
		return
	}
	firstName = first.Name
	return
}

// TargetTypes get all the ugc ranks that AI prepared for us
func (d *Dao) TargetTypes() (tids []int32, err error) {
	if len(d.arcTypesRel) == 0 {
		err = ecode.ServiceUnavailable
		return
	}
	for _, v := range d.conf.Cfg.ZonesInfo.TargetTypes {
		if children, ok := d.arcTypesRel[v]; ok { // second level types
			for _, child := range children {
				tids = append(tids, child.ID)
			}
		}
		tids = append(tids, v)
	}
	return
}

// FirstTypes returns only first level of types
func (d *Dao) FirstTypes() (typeMap map[int32]*model.ArcType, err error) {
	if len(d.arcTypes) == 0 {
		err = ecode.ServiceUnavailable
		return
	}
	typeMap = make(map[int32]*model.ArcType)
	for _, v := range d.arcTypes { // only pick first level of types
		if v.Pid == 0 {
			typeMap[v.ID] = &model.ArcType{
				ID:   v.ID,
				Name: v.Name,
			}
		}
	}
	return
}

// TypeInfo returns the type info
func (d *Dao) TypeInfo(typeid int32) (*arcwar.Tp, error) {
	if len(d.arcTypes) == 0 {
		return nil, ecode.ServiceUnavailable
	}
	info, ok := d.arcTypes[typeid]
	if !ok {
		return nil, ecode.NothingFound
	}
	return info, nil
}

// TypeChildren returns a first level's type children
func (d *Dao) TypeChildren(typeid int32) (children []*arcwar.Tp, err error) {
	if len(d.arcTypesRel) == 0 {
		err = ecode.ServiceUnavailable
		return
	}
	children, found := d.arcTypesRel[typeid]
	if !found {
		err = ecode.NothingFound
		return
	}
	return
}
