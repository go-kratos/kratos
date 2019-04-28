package model

import (
	"encoding/json"
	"fmt"
	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/library/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// UT Data
var (
	DataID     = int64(75)
	DataIDs    = []int64{75, 80}
	DataSIDs   = []int64{90, 133, 136}
	DataTIDs   = []int64{1179, 1180, 1368, 1360}
	NoDataID   = int64(100000000)
	NoDataIDs  = []int64{100000000, 100000001}
	NoDataSIDs = []int64{100000000, 100000001}
	NoDataTIDs = []int64{100000000, 100000001}
)

// JSONEncode 仿phpJSONEncode
func JSONEncode(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Error("JSONEncode error(%v)", err)
		return ""
	}
	return string(b)
}

// String2Int64 convert string slice([]string) to int64 slice([]int64)
func String2Int64(arr []string) (r []int64) {
	var (
		id  int64
		err error
	)
	for _, v := range arr {
		if id, err = strconv.ParseInt(v, 10, 64); err != nil {
			continue
		}
		r = append(r, id)
	}
	return
}

// UniqueInt64 Ints returns a unique subset of the int slice provided.
func UniqueInt64(input []int64) []int64 {
	u := make([]int64, 0, len(input))
	m := make(map[int64]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// Min 获取数组中最小值
func Min(vars []int32) (minVar int32) {
	if vars != nil {
		minVar = vars[0]
		for _, v := range vars {
			if v < minVar {
				minVar = v
			}
		}
	}
	return
}

// Max 获取数组中最大值
func Max(vars []int32) (maxVar int32) {
	for _, v := range vars {
		if v > maxVar {
			maxVar = v
		}
	}
	return
}

// GetTicketIDFromBase baseCenter获取票价id
func GetTicketIDFromBase() (int64, error) {
	params := url.Values{}
	params.Add("count", "1")
	params.Add("biz_tag", "price")
	params.Add("app_id", conf.Conf.BASECenter.AppID)
	params.Add("app_token", conf.Conf.BASECenter.AppToken)

	reqParam := params.Encode()

	resp, err := http.Get(fmt.Sprintf(conf.Conf.BASECenter.URL+"orderid/get?%s", reqParam))

	if err != nil {
		log.Error("获取票价id HTTP REQUEST失败")
		return 0, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("读取HTTP RESPONSE失败")
		return 0, err
	}
	var result Response
	if err := json.Unmarshal(content, &result); err != nil {
		log.Error("json解析失败")
	}
	return result.Data[0], nil
}

var alphabetTable []string

// AlphabetTable 获取票价所需symbol的字母表
func AlphabetTable() []string {
	if alphabetTable != nil {
		return alphabetTable
	}
	result := make([]string, 52)
	var i int

	ch := 97
	for i = 0; i < 26; i++ {
		result[i] = string(ch + i)
	}

	ch = 65
	j := i
	for i = 0; i < 26; i++ {
		result[j] = string(ch + i)
		j++
	}

	alphabetTable = result
	return alphabetTable
}

// ClassifyIDs 获取已经存在和需要被删除的id list
func ClassifyIDs(oldIDs []int64, newIDs []int64) (needDel []int64, existed []int64) {
	newIDsMap := make(map[int64]int64)
	for _, newID := range newIDs {
		newIDsMap[newID] = newID
	}

	for _, oldID := range oldIDs {
		if oldID == 0 {
			continue
		}
		if _, ok := newIDsMap[oldID]; !ok {
			needDel = append(needDel, oldID)
		} else {
			existed = append(existed, oldID)
		}
	}
	return
}

// Implode 仅支持不同类型的数组
func Implode(glue string, list interface{}) string {
	listValue := reflect.Indirect(reflect.ValueOf(list))
	if listValue.Kind() != reflect.Slice {
		// 数组以外类型返回空字符串
		return ""
	}
	count := listValue.Len()
	listStr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		str := fmt.Sprint(listValue.Index(i).Interface())
		listStr = append(listStr, str)
	}
	return strings.Join(listStr, glue)
}
