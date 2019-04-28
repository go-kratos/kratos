package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/library/log"
)

type multransFunc func(context.Context, []int64) (map[int64][]interface{}, error)

/* 批量查询，批量转换
 * list 	[]*struct{}
 * multrans 转化器，根据ID查出其他值
 * ID   	id字段名称，id字段类型必须是int64
 * Names 	查出来的各个字段名称
 */
func (s *Service) mulIDtoName(c context.Context, list interface{}, multrans multransFunc, ID string, Names ...string) (err error) {
	var (
		lV, itemI, itemIE, idFiled, nameFiled, valueField reflect.Value
		id                                                int64
		ids                                               []int64
		hashIDName                                        = make(map[int64][]interface{})
	)

	if lV = reflect.ValueOf(list); !lV.IsValid() || lV.IsNil() || lV.Kind() != reflect.Slice {
		return fmt.Errorf("invalid list")
	}

	count := lV.Len()
	for i := 0; i < count; i++ {
		if itemI = lV.Index(i); !itemI.IsValid() || itemI.IsNil() || itemI.Kind() != reflect.Ptr {
			return fmt.Errorf("invalid itemI")
		}
		if itemIE = itemI.Elem(); !itemIE.IsValid() || itemIE.Kind() != reflect.Struct {
			return fmt.Errorf("invalid itemIE")
		}

		if idFiled = itemIE.FieldByName(ID); !idFiled.IsValid() || idFiled.Kind() != reflect.Int64 {
			return fmt.Errorf("invalid idFiled")
		}
		for _, name := range Names {
			if nameFiled = itemIE.FieldByName(name); !nameFiled.IsValid() || !nameFiled.CanSet() {
				return fmt.Errorf("invalid nameFiled")
			}
		}
		if id = idFiled.Int(); id != 0 {
			if _, ok := hashIDName[id]; !ok {
				hashIDName[id] = []interface{}{}
				ids = append(ids, id)
			}
		}
	}
	if hashIDName, err = multrans(c, ids); err != nil {
		log.Error("multrans error(%v)", ids)
		return
	}

	for i := 0; i < count; i++ {
		itemIE = lV.Index(i).Elem()
		id = itemIE.FieldByName(ID).Int()
		if names, ok := hashIDName[id]; ok && len(names) == len(Names) {
			for i, name := range names {
				nameFiled = itemIE.FieldByName(Names[i])
				valueField = reflect.ValueOf(name)
				if nameFiled.Kind() != valueField.Kind() {
					log.Error("multrans(%s) return %v while need %v", name, valueField.Kind(), nameFiled.Kind())
					continue
				}
				itemIE.FieldByName(Names[i]).Set(reflect.ValueOf(name))
			}
		}
	}
	return
}

func (s *Service) getUserGroup(c context.Context, ids []int64) (group map[int64]*common.Group) {
	group = make(map[int64]*common.Group)
	for _, id := range ids {
		group[id] = s.groupCache[id]
	}
	return
}

func mergeSlice(arr1 []int64, arr2 []int64) (arr []int64) {
	for _, id1 := range arr1 {
		for _, id2 := range arr2 {
			if id1 == id2 {
				arr = append(arr, id1)
			}
		}
	}
	return
}

func fnvhash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func joinstr(maphash []string, sep string, max int) (msgs []string) {
	var msg string
	for _, v := range maphash {
		if len(msg) == 0 {
			msg = v
		} else if len(msg)+len(v) < max {
			msg += sep + v
		} else {
			msgs = append(msgs, msg)
			msg = v
		}
	}
	msgs = append(msgs, msg)
	return
}

//stringset 过滤重复字符串
func stringset(arr []string) (res []string) {
	mapHash := make(map[uint32]struct{})
	for _, item := range arr {
		hf := fnvhash32(item)
		if _, ok := mapHash[hf]; ok {
			continue
		}
		mapHash[hf] = struct{}{}
		res = append(res, item)
	}
	return
}
