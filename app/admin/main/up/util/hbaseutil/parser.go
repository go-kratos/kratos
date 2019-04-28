package hbaseutil

import (
	"encoding/binary"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"go-common/library/log"
	"reflect"
	"strconv"
)

/*
	add tag for struct fields:
	you can add tag:
		family: for hbase family, can be omitted, if omitted, the qualifier would be set at whatever famliy
		qualifier: for hbase qualifier
			if omitted, the fields must be map[string]int or map[string]string
	see parser_test.go for detail
*/

type field struct {
	parser      *Parser
	name        string
	structField reflect.StructField
	fieldValue  reflect.Value
	family      string
}

func (f *field) isValid() bool {
	return f.fieldValue.IsValid()
}

func (f *field) setValue(c *hrpc.Cell) (err error) {
	if c == nil {
		return
	}

	if f.fieldValue.Kind() == reflect.Ptr {
		f.fieldValue.Set(reflect.New(f.fieldValue.Type().Elem()))
		f.fieldValue = f.fieldValue.Elem()
	}
	switch f.fieldValue.Kind() {
	case reflect.Map:
		err = f.setMap(c)
	default:
		err = setBasicValue(c.Value, f.fieldValue, f.name, f.parser.ParseIntFunc)
	}
	return
}

func setBasicValue(value []byte, rv reflect.Value, name string, parsefunc ParseIntFunc) (err error) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		rv.Set(reflect.ValueOf(string(value)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i, e = parsefunc(value, rv, name)
		if e != nil {
			err = fmt.Errorf("field=%s, fail to convert: %s", name, e)
			return
		}
		if rv.OverflowInt(int64(i)) {
			log.Warn("field overflow, field=%s, value=%d, field type=%s", name, i, rv.Type().Name())
		}
		rv.SetInt(int64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i, e = parsefunc(value, rv, name)
		if e != nil {
			err = fmt.Errorf("field=%s, fail to convert: %s", name, e)
			return
		}
		if rv.OverflowUint(i) {
			log.Warn("field overflow, field=%s, value=%d, field type=%s", name, i, rv.Type().Name())
		}
		rv.SetUint(i)
	default:
		err = fmt.Errorf("cannot convert type:%s, kind()=%v, field=%s", rv.Type().Name(), rv.Kind(), name)
	}
	return
}

func (f *field) setMap(c *hrpc.Cell) (err error) {
	var fieldType = f.fieldValue.Type()
	if f.fieldValue.IsNil() {
		f.fieldValue.Set(reflect.MakeMap(fieldType))
		fieldType = f.fieldValue.Type()
	}
	var keyType = fieldType.Key()
	if keyType.Kind() != reflect.String {
		err = fmt.Errorf("cannot convert to map, only support map key: (string), but get type=%s", keyType.Name())
		return
	}

	var val = reflect.Indirect(reflect.New(fieldType.Elem()))
	err = setBasicValue(c.Value, val, f.name, f.parser.ParseIntFunc)
	if err != nil {
		err = fmt.Errorf("cannot convert to map, only support map value: (integer, string), type=%s, err=%v", fieldType.Name(), err)
		return
	}

	var key = indirect(reflect.New(fieldType.Key()))
	key.SetString(string(c.Qualifier))
	f.fieldValue.SetMapIndex(key, val)
	return
}

// ParseIntFunc function to parse []byte to uint64
// if not set, will assume []byte is big endian form of integer, length of 1/2/4/8 bytes
type ParseIntFunc func(v []byte, rv reflect.Value, fieldname string) (result uint64, err error)

//Parser parser for hbase cell
type Parser struct {
	ParseIntFunc ParseIntFunc
}

func getOrCreateFieldMapByFamily(familyMap map[string]map[string]field, key string) (result map[string]field) {
	var ok bool
	if result, ok = familyMap[key]; !ok {
		result = make(map[string]field)
		familyMap[key] = result
	}
	return result
}

func getField(familyMap map[string]map[string]field, family string, qualifier string) (result field) {
	var ok bool
	var qualifierMap map[string]field
	if qualifierMap, ok = familyMap[family]; !ok {
		qualifierMap, ok = familyMap[""]
		if !ok {
			return
		}
	}

	if result, ok = qualifierMap[qualifier]; !ok {
		qualifierMap, ok = familyMap[""]
		if ok {
			result, ok = qualifierMap[qualifier]
		}
		if !ok {
			return
		}
	}
	return
}

//Parse parse cell to struct
// supported type:
// 	integer from 16 ~ 64 bit, the cell's value must be big endian form of the integer, length could be 2 or 4 or 8 bytes
// 	string
func (p *Parser) Parse(cell []*hrpc.Cell, ptr interface{}) (err error) {
	if len(cell) == 0 {
		log.Warn("cell length = 0, nothing to parse")
		return
	}
	var familyFieldMap = make(map[string]map[string]field)
	// field only have family, and type is map[string]{integer,string}
	var familyOnlyMap = make(map[string]field)
	//var noFamilyFieldMap = make(map[string]reflect.Value)

	var ptrType = reflect.TypeOf(ptr)
	// if it's ptr
	if ptrType.Kind() == reflect.Ptr {
		var value = reflect.ValueOf(ptr)
		value = indirect(value)
		var valueType = value.Type()
		var valueKind = valueType.Kind()
		if valueKind == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				fieldInfo := valueType.Field(i) // a reflect.StructField
				tag := fieldInfo.Tag            // a reflect.StructTag
				//fmt.Printf("tag for field: %s, tag: %s\n", fieldInfo.Name, tag)
				family := tag.Get("family")
				qualifier := tag.Get("qualifier")

				var field = field{
					family:      family,
					name:        fieldInfo.Name,
					structField: fieldInfo,
					fieldValue:  value.Field(i),
					parser:      p,
				}
				// if no qualifier, or star, we create only family field
				if qualifier == "" || qualifier == "*" {
					if fieldInfo.Type.Kind() != reflect.Map {
						log.Warn("%s.%s, family-only field only support map, but get(%s)", ptrType.Name(), fieldInfo.Name, fieldInfo.Type.Name())
						continue
					}
					familyOnlyMap[family] = field
				} else {
					// save field info
					var fieldMapForFamily = getOrCreateFieldMapByFamily(familyFieldMap, family)
					fieldMapForFamily[qualifier] = field
				}
			}
		} else {
			log.Warn("cannot decode, unsupport type(%s)", valueKind.String())
		}
	}
	if p.ParseIntFunc == nil {
		p.ParseIntFunc = ByteBigEndianToUint64
	}
	// parse
	for _, c := range cell {
		var family = string(c.Family)
		var qualifier = string(c.Qualifier)
		//log.Info("parse cell, family=%s, qualifier=%s", family, qualifier)
		var fieldValue = getField(familyFieldMap, family, qualifier)
		if !fieldValue.isValid() {
			fieldValue = familyOnlyMap[family]
			if !fieldValue.isValid() {
				//log.Warn("no field for cell, family=%s, qualifier=%s", family, qualifier)
				continue
			}
		}
		if e := fieldValue.setValue(c); e != nil {
			log.Warn("fail to set value, err=%v", e)
			continue
		}
	}

	return
}

// indirect returns the value pointed to by a pointer.
// Pointers are followed until the value is not a pointer.
// New values are allocated for each nil pointer.
//
// An exception to this rule is if the value satisfies an interface of
// interest to us (like encoding.TextUnmarshaler).
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return v
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	} else if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return indirect(reflect.Indirect(v))
}

//StringToUint parse string to uint
func StringToUint(value []byte, rv reflect.Value, fieldname string) (result uint64, err error) {
	if len(value) == 0 {
		return
	}
	if value[0] == '-' {
		i64, e := strconv.ParseInt(string(value), 10, 64)
		err = e
		result = uint64(i64)
	} else {
		result, err = strconv.ParseUint(string(value), 10, 64)
	}

	return
}

//ByteBigEndianToUint64 convert big endian to uint64
func ByteBigEndianToUint64(value []byte, rv reflect.Value, fieldname string) (result uint64, err error) {
	var length = len(value)
	switch length {
	case 4:
		result = uint64(binary.BigEndian.Uint32(value))
	case 8:
		result = uint64(binary.BigEndian.Uint64(value))
	case 2:
		result = uint64(binary.BigEndian.Uint16(value))
	case 1:
		result = uint64(value[0])
	default:
		err = fmt.Errorf("cannot decode to integer, byteslen=%d, only support (1,2,4,8)", length)
	}
	if err == nil {
		var vlen = len(value)
		var rvType = rv.Type()
		if rvType.Size() != uintptr(vlen) {
			log.Error("field=%s type=%s length=%d, cell length=%d, doesn't match, may yield wrong value!",
				fieldname, rvType.Name(), rvType.Size(), vlen)
		}
	}

	return
}
