package middleware

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"go-common/library/log"
)

const (
	//TypeMiddleAll 全部操作都允许
	TypeMiddleAll = int8(0)
	//TypeMiddleEncode 只允许编码，后端返回给前端的响应处理
	TypeMiddleEncode = int8(1)
	//TypeMiddleDecode 只允许解码，前端请求后端的请求参数处理
	TypeMiddleDecode = int8(2)
)

//Aggregate 前端映射结构，对应于前后端交互的字段
type Aggregate struct {
	Hitn      string `json:"hitn"`      //name的多个结构体用.表示分级
	Hitv      string `json:"hitv"`      //枚举值
	Mapn      string `json:"mapn"`      //映射字段名，可以与hitname不同，也可以相同
	Mapv      string `json:"mapv"`      //映射字段值
	Delimiter string `json:"delimiter"` //映射字段值的分隔符号
	Order     int64  `json:"order"`     //顺序
	Type      int8   `json:"type"`
}

func (f *Aggregate) Process(data interface{}, encode bool) {
	var (
		field, fieldm     reflect.Value
		fieldExist, hited bool
	)

	defer func() {
		if errs := recover(); errs != nil {
			log.Error("Aggregate Process error(%+v)", errs)
		}
	}()
	hitn := f.Hitn
	hitv := strings.Split(f.Hitv, f.Delimiter)
	mapn := f.Mapn
	mapv := f.Mapv
	if !encode {
		hitn = f.Mapn
		hitv = strings.Split(f.Mapv, f.Delimiter)
		mapn = f.Hitn
		mapv = f.Hitv
	}

	//check fields exist
	fv := reflect.ValueOf(data)
	if field, fieldExist = getFieldByName(fv, hitn); !fieldExist {
		log.Warn("no field for hit(%s) data(%+v)", hitn, data)
		return
	}
	if mapn == hitn {
		fieldm = field
	} else if fieldm, fieldExist = getFieldByName(fv, mapn); !fieldExist || !fieldm.CanSet() {
		log.Warn("no field for map(%s) data(%+v)", mapn, data)
		return
	}

	fieldv := fmt.Sprintf("%v", field.Interface())
	for _, hit := range hitv {
		if fieldv == hit {
			hited = true
			break
		}
	}
	if !hited {
		return
	}

	log.Info("got hit field(%s) value(%s) config(%+v)", hitn, fieldv, f)
	switch fieldm.Kind() {
	case reflect.String:
		fieldm.SetString(mapv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vv, err := strconv.ParseInt(mapv, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", mapv, err)
			return
		}
		fieldm.SetInt(vv)
	default:
		log.Warn("not support field.kind(%s) for field(%s)", fieldm.Kind().String(), mapn)
	}
}

//getFieldByName 迭代遍历struct，获取指定名字的字段
func getFieldByName(v reflect.Value, name string) (res reflect.Value, ok bool) {
	tp := v.Type()
	if tp.Kind() == reflect.Ptr {
		v = v.Elem()
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Struct || !v.IsValid() {
		return
	}

	for i := 0; i < tp.NumField(); i++ {
		fn := strings.ToLower(tp.Field(i).Name)
		if fn == name {
			res = v.Field(i)
			ok = true
			return
		}

		if vres, vok := getFieldByName(v.Field(i), name); vok {
			ok = vok
			res = vres
			return
		}
	}
	return
}

//MiddleAggregate 处理聚合逻辑
type MiddleAggregate struct {
	Cfg    []*Aggregate
	Encode bool
}

//Process handle multi aggregate logists
func (m *MiddleAggregate) Process(data interface{}) {
	cfgs := []*Aggregate{}
	for _, item := range m.Cfg {
		if item.Type == TypeMiddleAll || (m.Encode && item.Type == TypeMiddleEncode) || (!m.Encode && item.Type == TypeMiddleDecode) {
			cfgs = append(cfgs, item)
		}
	}
	if len(cfgs) == 0 {
		return
	}

	sort.Sort(AggregateArr(cfgs))
	for _, item := range cfgs {
		item.Process(data, m.Encode)
	}
}

//AggregateArr arr
type AggregateArr []*Aggregate

//Len .
func (f AggregateArr) Len() int {
	return len(f)
}

//Less .
func (f AggregateArr) Less(i, j int) bool {
	return f[i].Order < f[j].Order
}

//Swap .
func (f AggregateArr) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
