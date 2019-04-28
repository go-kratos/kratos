package service

import (
	"context"
	"fmt"
	"reflect"

	"go-common/library/log"
)

// 每个ID单独查询 strict严格模式下一次错误，直接返回
func (s *Service) singleIDtoName(c context.Context, list interface{}, singletrans func(context.Context, int64) ([]interface{}, error), strict bool, ID string, Names ...string) (err error) {
	var (
		lV, itemI, itemIE, idFiled, nameFiled, valueField reflect.Value
		id                                                int64
		values                                            []interface{}
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
		for _, Name := range Names {
			if nameFiled = itemIE.FieldByName(Name); !nameFiled.IsValid() || !nameFiled.CanSet() {
				return fmt.Errorf("invalid nameFiled")
			}
		}

		if id = idFiled.Int(); id != 0 {
			if values, err = singletrans(c, id); err != nil || len(values) != len(Names) {
				log.Error("s.sigleIDtoName error(%v) len(values)=%d len(Names)=%d", err, len(values), len(Names))
				if strict {
					return
				}
				err = nil
				continue
			}
			for i, value := range values {
				nameFiled = itemIE.FieldByName(Names[i])
				valueField = reflect.ValueOf(value)
				if nameFiled.Kind() != valueField.Kind() {
					log.Error("singletrans return %s while need %s", valueField.Kind().String(), nameFiled.Kind().String())
					continue
				}
				nameFiled.Set(valueField)
			}
		}
	}
	return
}

/* 批量查询，批量转换
 * list 	[]*struct{}
 * multrans 转化器，根据ID查出其他值
 * ID   	id字段名称，id字段类型必须是int64
 * Names 	查出来的各个字段名称
 */
func (s *Service) mulIDtoName(c context.Context, list interface{}, multrans func(context.Context, []int64) (map[int64][]interface{}, error), ID string, Names ...string) (err error) {
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
					log.Error("multrans return %v while need %v", ids)
					continue
				}
				itemIE.FieldByName(Names[i]).Set(reflect.ValueOf(name))
			}
		}
	}
	return
}
