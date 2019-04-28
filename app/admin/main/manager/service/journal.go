package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// SearchLogAudit .
func (s *Service) SearchLogAudit(c *bm.Context) (result []byte, err error) {
	var (
		res *model.LogRes
	)
	if res, err = s.dao.SearchLogAudit(c); err != nil {
		log.Error("s.dao.SearchLogAduit error (%v)", err)
		return
	}
	result, err = outData(res)
	return
}

// SearchLogAction .
func (s *Service) SearchLogAction(c *bm.Context) (result []byte, err error) {
	var (
		res *model.LogRes
	)
	if res, err = s.dao.SearchLogAction(c); err != nil {
		log.Error("s.dao.SearchLogAduit error (%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("res.Code error (%d)", res.Code)
		return
	}
	result, err = outData(res)
	return
}

// OutData .
func outData(res *model.LogRes) (result []byte, err error) {
	var resultMap []map[string]interface{}
	//Output the data
	for _, j := range res.Data.Result {
		item := map[string]interface{}{}
		decoder := json.NewDecoder(bytes.NewReader(j))
		decoder.UseNumber()
		decoder.Decode(&item)
		item = numberToInt64(item)
		resultMap = append(resultMap, item)
	}

	for parentKey, parentValue := range resultMap {
		if extraValue, ok := parentValue["extra_data"]; ok {
			if extraData, ok := extraValue.(string); ok {
				p := make(map[string]interface{})
				if e := json.Unmarshal([]byte(extraData), &p); e == nil {
					for childKey, childValue := range backtrace(p) {
						resultMap[parentKey][childKey] = childValue
					}
				}
			}
		}
	}

	titleMap := make(map[string]string)
	// Iterator title collections
	for _, v := range resultMap {
		for key := range v {
			titleMap[key] = key
		}
	}
	// Get the titles
	titleStr := []string{}
	for _, v := range titleMap {
		titleStr = append(titleStr, v)
	}
	data := [][]string{}
	data = append(data, titleStr)
	for _, value := range resultMap {
		fields := []string{}
		for _, parentTitle := range titleStr {
			if value[parentTitle] != nil {
				fields = append(fields, fmt.Sprintf("%v", value[parentTitle]))
			} else {
				fields = append(fields, "")
			}
		}
		data = append(data, fields)
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	for _, record := range data {
		if err = w.Write(record); err != nil {
			log.Error("w Write (%v) error (%v)", record, err)
			return
		}
	}
	w.Flush()
	result = buf.Bytes()
	return
}

// backtrace .
func backtrace(in map[string]interface{}) (out map[string]interface{}) {
	out = make(map[string]interface{})
	for k, v := range in {
		if z, ok := v.(map[string]interface{}); ok {
			for childKey, childValue := range backtrace(z) {
				out[childKey] = childValue
			}
		} else {
			out[k] = v
		}
	}
	return
}

// numberToInt64 .
func numberToInt64(in map[string]interface{}) (out map[string]interface{}) {
	var err error
	out = map[string]interface{}{}
	for k, v := range in {
		if integer, ok := v.(json.Number); ok {
			if out[k], err = integer.Int64(); err != nil {
				log.Error("service.log.numberToInt64(%v)(%v)", integer, err)
			}
		} else {
			out[k] = v
		}
	}
	return
}
