package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

// MapData .
type MapData map[string]interface{}

// StrID .
func (m MapData) StrID(attrs *Attrs) (string, bool) {
	if attrs.Index.IndexID == "UUID" {
		return "", true
	}
	var data []interface{}
	arr := strings.Split(attrs.Index.IndexID, ",")
	arrLen := len(arr)
	if arrLen >= 2 {
		for _, v := range arr[1:] {
			if item, ok := m[v].(*interface{}); ok {
				data = append(data, item)
				continue
			}
			if item, ok := m[v].(interface{}); ok {
				data = append(data, item)
				continue
			}
			log.Error("model.map_data.StrID err (%v)", v)
			return "", false
		}
		return fmt.Sprintf(arr[0], data...), true
	}
	return "", false
}

// Index .
func (m MapData) Index(attrs *Attrs) (indexAliasName, indexEntityName string) {
	switch attrs.Index.IndexSplit {
	case "single":
		indexAliasName = attrs.Index.IndexAliasPrefix
		indexEntityName = attrs.Index.IndexEntityPrefix
	case "int":
		if attrs.DataSQL.DataIndexSuffix != "" {
			s := strings.Split(attrs.DataSQL.DataIndexSuffix, ";")
			v := strings.Split(s[1], ":")
			if id, ok := m[v[0]].(*interface{}); ok {
				// indexAliasName = fmt.Sprintf("%s%d", attrs.Index.IndexAliasPrefix, (*id).(int64)%100) // mod
				divisor, _ := strconv.ParseInt(v[2], 10, 64)
				indexAliasName = fmt.Sprintf("%s"+s[0], attrs.Index.IndexAliasPrefix, (*id).(int64)%divisor)
				indexEntityName = fmt.Sprintf("%s"+s[0], attrs.Index.IndexEntityPrefix, (*id).(int64)%divisor)
			}
			if id, ok := m[v[0]].(interface{}); ok {
				divisor, _ := strconv.ParseInt(v[2], 10, 64)
				indexAliasName = fmt.Sprintf("%s"+s[0], attrs.Index.IndexAliasPrefix, id.(int64)%divisor)
				indexEntityName = fmt.Sprintf("%s"+s[0], attrs.Index.IndexEntityPrefix, id.(int64)%divisor)
			}
		}
	}
	//fmt.Println("indexname", indexAliasName, indexEntityName)
	return
}

// DtbIndex .
// func (m MapData) DtbIndex(attrs *Attrs) (indexName string) {
// 	if attrs.Index.IndexZero == "0" {
// 		indexName = attrs.Index.IndexAliasPrefix
// 		return
// 	}
// 	if attrs.DataSQL.DataIndexSuffix != "" {
// 		s := strings.Split(attrs.DataSQL.DataIndexSuffix, ";")
// 		v := strings.Split(s[1], ":")
// 		divisor, _ := strconv.ParseInt(v[2], 10, 64)
// 		indexName = fmt.Sprintf("%s"+s[0], attrs.Index.IndexAliasPrefix, int64(m[v[0]].(float64))%divisor)
// 	}
// 	return
// }

// PrimaryID .
func (m MapData) PrimaryID() int64 {
	if m["_id"] != nil {
		if id, ok := m["_id"].(*interface{}); ok {
			return (*id).(int64)
		}
	}
	return 0
}

// StrMTime .
func (m MapData) StrMTime() string {
	if m["_mtime"] != nil {
		if mtime, ok := m["_mtime"].(*interface{}); ok {
			return (*mtime).(time.Time).Format("2006-01-02 15:04:05")
		} else if mtime, ok := m["_mtime"].(string); ok {
			return mtime
		}
	}
	return ""
}

// StrCTime .
func (m MapData) StrCTime() string {
	if m["ctime"] != nil {
		if ctime, ok := m["ctime"].(*interface{}); ok {
			return (*ctime).(time.Time).Format("2006-01-02")
		} else if ctime, ok := m["ctime"].(string); ok {
			return ctime
		}
	}
	return ""
}

// InetNtoA int64 to string ip.
func (m MapData) InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// TransData transfer address into value
func (m MapData) TransData(attr *Attrs) {
	for k, v := range m {
		// transfer automaticlly
		if v2, ok := v.(*interface{}); ok {
			switch (*v2).(type) {
			case time.Time:
				m[k] = (*v2).(time.Time).Format("2006-01-02 15:04:05")
			case []uint, []uint8, []uint16, []uint32, []uint64:
				m[k] = string((*v2).([]byte))
			case int, int8, int16, int32, int64: // 一定要，用于extra_data查询
				m[k] = (*v2).(int64)
			case nil:
				m[k] = int64(0) //给个默认值，当查到为null时
			default:
				// other types
			}
		}
		// transfer again by custom
		if t, ok := attr.DataSQL.DataIndexFormatFields[k]; ok {
			if v3, ok := v.(*interface{}); ok {
				switch t {
				case "ip":
					if *v3 == nil {
						*v3 = int64(0)
					}
					ipFormat := m.InetNtoA((*v3).(int64))
					m[k+"_format"] = ipFormat
				case "arr":
					var arr []int
					binaryAttributes := strconv.FormatInt((*v3).(int64), 2)
					for i := len(binaryAttributes) - 1; i >= 0; i-- {
						b := fmt.Sprintf("%c", binaryAttributes[i])
						if b == "1" {
							arr = append(arr, len(binaryAttributes)-i)
						}
					}
					m[k+"_format"] = arr
				case "bin":
					var arr []int
					binaryAttributes := strconv.FormatInt((*v3).(int64), 2)
					for i := len(binaryAttributes) - 1; i >= 0; i-- {
						b := fmt.Sprintf("%c", binaryAttributes[i])
						if b == "1" {
							arr = append(arr, len(binaryAttributes)-i)
						}
					}
					m[k] = arr
				case "array_json":
					var arr []int64
					arr = []int64{}
					json.Unmarshal([]byte((*v3).([]uint8)), &arr) //如果不是json就是空数组
					// println(len(arr))
					m[k] = arr
				case "day":
					m[k] = (*v3).(time.Time).Format("2006-01-02")
				case "workflow":
					delete(m, k)
				default:
					// other types
				}
			}
		}
	}
}

// TransDtb transfer databus fields into es fields
func (m MapData) TransDtb(attr *Attrs) {
	// TODO	注释要打开，不然无法移除不要的dtb字段
	// for k := range m {
	// 	if _, ok := attr.DataSQL.DataDtbFields[k]; !ok {
	// 		if k == "index_name" {
	// 			continue
	// 		}
	// 		delete(m, k)
	// 	}
	// }
	res := map[string]interface{}{}
	for k, dv := range attr.DataSQL.DataDtbFields {
		for _, dk := range dv {
			if v, ok := m[k]; ok {
				switch v.(type) {
				case float64:
					res[dk] = int64(v.(float64))
				default:
					res[dk] = v
				}
			}
		}
	}
	for k := range res {
		m[k] = res[k]
	}
	id, okID := attr.DataSQL.DataFieldsV2["_id"]
	key, okKey := attr.DataSQL.DataDtbFields[id.Field]
	if len(key) >= 1 && okID && okKey {
		m["_id"] = m[key[0]]
	}
}
