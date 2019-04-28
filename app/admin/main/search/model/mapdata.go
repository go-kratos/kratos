package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go-common/library/log"
)

// MapData .
type MapData map[string]interface{}

// StrID .
func (m MapData) StrID(indexID string) string {
	if indexID == "base" { // 需要改配置
		return ""
	}
	var data []interface{}
	arr := strings.Split(indexID, ",")
	for _, v := range arr[1:] {
		v = strings.TrimSpace(v)
		if item, ok := m[v].(interface{}); ok {
			if reflect.TypeOf(item).Kind() == reflect.Float64 {
				item = int64(item.(float64))
			}
			data = append(data, item)
			continue
		}
		log.Error("model.MapData.StrID err (%v)", v)
	}
	if len(data) == 0 {
		return ""
	}
	return fmt.Sprintf(arr[0], data...)
}

func (m MapData) NumberToInt64() (err error) {
	for k, v := range m {
		if integer, ok := v.(json.Number); ok {
			if m[k], err = integer.Int64(); err != nil {
				log.Error("service.log.numberToInt64(%v)(%v)", integer, err)
			}
		}
	}
	return
}
