package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func convert(params map[string][]string) (res map[string][]interface{}) {
	res = make(map[string][]interface{})
	for k, v := range params {
		var arr []interface{}
		for _, m := range v {
			if m != "" {
				for _, u := range strings.Split(m, ",") {
					arr = append(arr, u)
				}
			}
		}
		if len(arr) > 0 {
			res[k] = arr
		}
	}
	return res
}

// Check .
func (s *Service) Check(appID string, businessID int) (business *model.Business, ok bool) {
	business, ok = s.dao.GetLogInfo(appID, businessID)
	return
}

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

// 获取部门
func (s *Service) uDepTs(c context.Context, res *model.SearchResult, sp *model.LogParams) *model.SearchResult {
	var (
		uids   []string
		result []map[string]interface{}
		err    error
	)
	result = []map[string]interface{}{}
	for _, j := range res.Result {
		item := map[string]interface{}{}
		decoder := json.NewDecoder(bytes.NewReader(j))
		decoder.UseNumber()
		if err = decoder.Decode(&item); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索LogAudit JSON失败", sp.Bsp.AppID), "s.dao.LogAudit(%v) json error(%v)", sp, err)
		}
		item = numberToInt64(item)
		result = append(result, item)
		if _, ok := item["uid"]; ok {
			uids = append(uids, fmt.Sprintf("%v", item["uid"]))
		}
	}
	var depRs = &model.UDepTsData{
		Data: map[string]string{},
	}
	if len(uids) > 0 {
		if depRs, err = s.dao.UDepTs(c, uids); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索LogAudit失败", sp.Bsp.AppID), "s.dao.LogAudit(%v) error(%v)", sp, err)
			err = ecode.SearchLogAuditFailed
			return res
		}
	}
	for i, j := range result {
		result[i]["department"] = ""
		if _, ok := j["uid"]; ok {
			if m, sok := depRs.Data[fmt.Sprintf("%v", j["uid"])]; sok {
				result[i]["department"] = m
			}
		}
		if res.Result[i], err = json.Marshal(j); err != nil {
			log.Error("s.dao.LogAudit(%v) json res(%v)", err, j)
		}
	}
	return res
}

// 获取IP地址
func (s *Service) IP(c context.Context, res *model.SearchResult, sp *model.LogParams) *model.SearchResult {
	var (
		ip     []string
		result []map[string]interface{}
		err    error
	)
	result = []map[string]interface{}{}
	for _, j := range res.Result {
		item := map[string]interface{}{}
		decoder := json.NewDecoder(bytes.NewReader(j))
		decoder.UseNumber()
		if err = decoder.Decode(&item); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索LogUserAction JSON失败", sp.Bsp.AppID), "s.dao.LogUserAction(%v) json error(%v)", sp, err)
		}
		item = numberToInt64(item)
		result = append(result, item)
		if _, ok := item["ip"]; ok {
			if v := fmt.Sprintf("%v", item["ip"]); v != "" {
				ip = append(ip, v)
			}
		}
	}
	var ipData = &model.IPData{
		Data: map[string]struct {
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
			Isp      string `json:"isp"`
		}{},
	}
	if len(ip) > 0 {
		if ipData, err = s.dao.IP(c, ip); err != nil {
			dao.PromError(fmt.Sprintf("es:%s 搜索LogUserAction失败", sp.Bsp.AppID), "s.dao.LogUserAction(%v) error(%v)", sp, err)
			err = ecode.SearchLogAuditFailed
			return res
		}
	}
	for i, j := range result {
		result[i]["location"] = ""
		if _, ok := j["ip"]; ok {
			if m, sok := ipData.Data[fmt.Sprintf("%v", j["ip"])]; sok {
				location := make([]string, 0, 4)
				for _, v := range []string{m.Country, m.Province, m.City, m.Isp} {
					if v != "" {
						location = append(location, v)
					}
				}
				result[i]["location"] = strings.Join(location, "-")
			}
		}
		if res.Result[i], err = json.Marshal(j); err != nil {
			log.Error("s.dao.LogUserAction(%v) json res(%v)", err, j)
		}
	}
	return res
}

// LogAudit .
func (s *Service) LogAudit(c context.Context, params map[string][]string, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	p := convert(params)
	if res, err = s.dao.LogAudit(c, p, sp, business); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索LogAudit失败", sp.Bsp.AppID), "s.dao.LogAudit(%v) error(%v)", sp, err)
		err = ecode.SearchLogAuditFailed
		return
	}
	res = s.uDepTs(c, res, sp)
	return
}

// LogAuditGroupBy .
func (s *Service) LogAuditGroupBy(c context.Context, params map[string][]string, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	p := convert(params)
	if v, ok := p["group"]; !ok || len(v) == 0 {
		err = ecode.RequestErr
		return
	}
	if res, err = s.dao.LogAuditGroupBy(c, p, sp, business); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索LogAuditGroupBy失败", sp.Bsp.AppID), "s.dao.LogAuditGroupBy(%v) error(%v)", sp, err)
		err = ecode.SearchLogAuditOidFailed
		return
	}
	if res.Page.Total < int64(res.Page.Ps*(res.Page.Pn-1)) || len(res.Result) == 0 {
		res.Result = []json.RawMessage{}
	} else if int64(res.Page.Ps*(res.Page.Pn-1)) <= res.Page.Total && res.Page.Total <= int64(res.Page.Ps*res.Page.Pn) {
		res.Result = res.Result[res.Page.Ps*(res.Page.Pn-1) : res.Page.Total]
	} else {
		res.Result = res.Result[res.Page.Ps*(res.Page.Pn-1) : res.Page.Ps*res.Page.Pn]
	}
	res = s.uDepTs(c, res, sp)
	return
}

// LogAuditDelete .
func (s *Service) LogAuditDelete(c context.Context, params map[string][]string, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	p := convert(params)
	if res, err = s.dao.LogAuditDelete(c, p, sp, business); err != nil {
		dao.PromError(fmt.Sprintf("es:%s LogAuditDelete失败", sp.Bsp.AppID), "s.dao.LogAuditDelete(%v) error(%v)", sp, err)
		err = ecode.SearchLogAuditFailed
		return
	}
	return
}

// LogUserAction .
func (s *Service) LogUserAction(c context.Context, params map[string][]string, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	p := convert(params)
	if res, err = s.dao.LogUserAction(c, p, sp, business); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索LogUserAction失败", sp.Bsp.AppID), "s.dao.LogUserAction(%v) error(%v)", sp, err)
		err = ecode.SearchLogUserActionFailed
		return
	}
	res = s.IP(c, res, sp)
	return
}

// LogUserActionDelete .
func (s *Service) LogUserActionDelete(c context.Context, params map[string][]string, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	p := convert(params)
	if res, err = s.dao.LogUserActionDelete(c, p, sp, business); err != nil {
		dao.PromError(fmt.Sprintf("es:%s LogUserActionDelete失败", sp.Bsp.AppID), "s.dao.LogUserActionDelete(%v) error(%v)", sp, err)
		err = ecode.SearchLogAuditFailed
		return
	}
	return
}

func (s *Service) LogCount(c context.Context, name string, business int, uid interface{}) {
	s.dao.LogCount(c, name, business, uid)
}
