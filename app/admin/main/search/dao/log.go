package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/search/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
)

const (
	_sql     = "SELECT id, name, index_format, index_cluster, additional_mapping, permission_point FROM digger_"
	_count   = "INSERT INTO digger_count (`business`,`type`,`time`,`count`) values (?, 'inc', ?, 1) ON DUPLICATE KEY UPDATE count=count+1"
	_percent = "INSERT INTO digger_count (`business`,`type`,`time`,`name`,`count`) values (?, 'inc', ?, ?, 1) ON DUPLICATE KEY UPDATE count=count+1"
)

var (
	logAuditBusiness      map[int]*model.Business
	logUserActionBusiness map[int]*model.Business
)

func (d *Dao) NewLogProcess() {
	for {
		if err := d.NewLog(); err != nil {
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Minute)
	}
}

// NewLog .
func (d *Dao) NewLog() (err error) {
	if logAuditBusiness, err = d.initMapping("log_audit"); err != nil {
		return
	}
	for k, v := range logAuditBusiness {
		if _, ok := d.esPool[v.IndexCluster]; !ok {
			log.Error("logAudit esPool no exist(%v)", k)
			delete(logAuditBusiness, k)
		}
	}
	if logUserActionBusiness, err = d.initMapping("log_user_action"); err != nil {
		return
	}
	for k, v := range logUserActionBusiness {
		if _, ok := d.esPool[v.IndexCluster]; !ok {
			log.Error("logUserAction esPool no exist(%v)", k)
			delete(logUserActionBusiness, k)
		}
	}
	return
}

// GetLogInfo .
func (d *Dao) GetLogInfo(appID string, id int) (business *model.Business, ok bool) {
	switch appID {
	case "log_audit":
		business, ok = logAuditBusiness[id]
		return
	case "log_user_action":
		business, ok = logUserActionBusiness[id]
		return
	}
	return &model.Business{}, false
}

func (d *Dao) initMapping(appID string) (business map[int]*model.Business, err error) {
	defaultMapping := map[string]string{}
	switch appID {
	case "log_audit":
		defaultMapping = model.LogAuditDefaultMapping
	case "log_user_action":
		defaultMapping = model.LogUserActionDefaultMapping
	}
	business = map[int]*model.Business{}
	rows, err := d.db.Query(context.Background(), _sql+appID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var value = &model.Business{
			AppID:   appID,
			Mapping: map[string]string{},
		}
		if err = rows.Scan(&value.ID, &value.Name, &value.IndexFormat, &value.IndexCluster, &value.AdditionalMapping, &value.PermissionPoint); err != nil {
			log.Error("Log New DB (%v)(%v)", appID, err)
			continue
		}
		if appID == "log_audit" {
			value.IndexCluster = "log"
		}
		for k, v := range defaultMapping {
			value.Mapping[k] = v
		}
		if value.AdditionalMapping != "" {
			var additionalMappingDict map[string]string
			if err = json.Unmarshal([]byte(value.AdditionalMapping), &additionalMappingDict); err != nil {
				log.Error("Log New Json (%v)(%v)", value.ID, err)
				continue
			}
			for k, v := range additionalMappingDict {
				value.Mapping[k] = v
			}
		}
		business[value.ID] = value
	}
	err = rows.Err()
	return
}

/**
获取es索引名 多个用逗号分隔
按日分最多7天
按周分最多2月
按月分最多6月
按年分最多3年
*/
func (d *Dao) logIndexName(c context.Context, p *model.LogParams, business *model.Business) (res string, err error) {
	var (
		sTime  = time.Now()
		eTime  = time.Now()
		resArr []string
	)
	if p.CTimeFrom != "" {
		sTime, err = time.Parse("2006-01-02 15:04:05", p.CTimeFrom)
		if err != nil {
			log.Error("d.LogAuditIndexName(%v)", p.CTimeFrom)
			return
		}
	}
	if p.CTimeTo != "" {
		eTime, err = time.Parse("2006-01-02 15:04:05", p.CTimeTo)
		if err != nil {
			log.Error("d.LogAuditIndexName(p.CTimeTo)(%v)", p.CTimeTo)
			return
		}
	}
	resDict := map[string]bool{}
	if strings.Contains(business.IndexFormat, "02") {
		for a := 0; a <= 60; a++ {
			resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, eTime)] = true
			eTime = eTime.AddDate(0, 0, -1)
			if (p.CTimeFrom == "" && a >= 1) || (p.CTimeFrom != "" && sTime.After(eTime)) {
				break
			}
		}
	} else if strings.Contains(business.IndexFormat, "week") {
		for a := 0; a <= 366; a++ {
			resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, eTime)] = true
			eTime = eTime.AddDate(0, 0, -1)
			if (p.CTimeFrom == "" && a >= 1) || (p.CTimeFrom != "" && sTime.After(eTime)) {
				resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, sTime)] = true
				break
			}
		}
	} else if strings.Contains(business.IndexFormat, "01") {
		// 1月31日时AddDate(0, -1, 0)会出现错误
		year, month, _ := eTime.Date()
		hour, min, sec := eTime.Clock()
		eTime = time.Date(year, month, 1, hour, min, sec, 0, eTime.Location())
		for a := 0; a <= 360; a++ {
			resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, eTime)] = true
			eTime = eTime.AddDate(0, -1, 0)
			if (p.CTimeFrom == "" && a >= 1) || p.CTimeFrom != "" && sTime.After(eTime) {
				break
			}
		}
	} else if strings.Contains(business.IndexFormat, "2006") {
		// 2月29日时AddDate(-1, 0, 0)会出现错误
		year, _, _ := eTime.Date()
		hour, min, sec := eTime.Clock()
		eTime = time.Date(year, 1, 1, hour, min, sec, 0, eTime.Location())
		for a := 0; a <= 100; a++ {
			resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, eTime)] = true
			eTime = eTime.AddDate(-1, 0, 0)
			if (p.CTimeFrom == "" && a >= 1) || (p.CTimeFrom != "" && sTime.After(eTime)) {
				break
			}
		}
	} else if business.IndexFormat == "all" {
		resDict[getLogAuditIndexName(p.Business, business.AppID, business.IndexFormat, eTime)] = true
	}
	for k := range resDict {
		if exist, e := d.ExistIndex(c, business.IndexCluster, k); exist && e == nil {
			resArr = append(resArr, k)
		}
	}
	res = strings.Join(resArr, ",")
	return
}

func getLogAuditIndexName(business int, indexName string, format string, time time.Time) (index string) {
	var (
		week = map[int]string{
			0: "0107",
			1: "0815",
			2: "1623",
			3: "2431",
		}
	)
	format = strings.Replace(time.Format(format), "week", week[time.Day()/8], -1)
	index = indexName + "_" + strconv.Itoa(business) + "_" + format
	return
}

func (d *Dao) getQuery(pr map[string][]interface{}, indexMapping map[string]string) (query *elastic.BoolQuery) {
	query = elastic.NewBoolQuery()
	for k, t := range indexMapping {
		switch t {
		case "int", "int64":
			if v, ok := pr[k]; ok {
				query = query.Filter(elastic.NewTermsQuery(k, v...))
			}
			if v, ok := pr[k+"_from"]; ok {
				query = query.Filter(elastic.NewRangeQuery(k).Gte(v[0]))
			}
			if v, ok := pr[k+"_to"]; ok {
				query = query.Filter(elastic.NewRangeQuery(k).Lte(v[0]))
			}
		case "string":
			if v, ok := pr[k]; ok {
				query = query.Filter(elastic.NewTermsQuery(k, v...))
			}
			if v, ok := pr[k+"_like"]; ok {
				likeMap := []model.QueryBodyWhereLike{
					{
						KWFields: []string{k},
						KW:       []string{fmt.Sprintf("%v", v)},
						Level:    model.LikeLevelHigh,
					},
				}
				if o, e := d.queryBasicLike(likeMap, ""); e == nil {
					query = query.Must(o...)
				}
			}
		case "time":
			if v, ok := pr[k+"_from"]; ok {
				query = query.Filter(elastic.NewRangeQuery(k).Gte(v[0]))
			}
			if v, ok := pr[k+"_to"]; ok {
				query = query.Filter(elastic.NewRangeQuery(k).Lte(v[0]))
			}
		case "int_to_bin":
			if v, ok := pr[k]; ok {
				var arr []elastic.Query
				for _, i := range v {
					item, err := strconv.ParseUint(i.(string), 10, 64)
					if err != nil {
						break
					}
					arr = append(arr, elastic.NewTermsQuery(k, 1<<(item-1)))
				}
				query = query.Filter(arr...)
			}
		case "array":
			if v, ok := pr[k+"_and"]; ok {
				for _, n := range v {
					query = query.Filter(elastic.NewTermsQuery(k, n))
				}
			}
			if v, ok := pr[k+"_or"]; ok {
				query = query.Filter(elastic.NewTermsQuery(k, v...))
			}
		}
	}
	return query
}

// LogAudit .
func (d *Dao) LogAudit(c context.Context, pr map[string][]interface{}, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	var indexName string
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	query := d.getQuery(pr, business.Mapping)
	indexName, err = d.logIndexName(c, sp, business)
	if err != nil {
		log.Error("d.LogAudit.logIndexName(%v)(%v)", err, indexName)
		return
	}
	if indexName == "" {
		return
	}
	if res, err = d.searchResult(c, business.IndexCluster, indexName, query, sp.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Bsp.AppID), "%v", err)
	}
	return
}

// LogAuditGroupBy .
func (d *Dao) LogAuditGroupBy(c context.Context, pr map[string][]interface{}, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	var (
		indexName    = ""
		searchResult *elastic.SearchResult
	)
	group := pr["group"][0].(string)
	if _, ok := d.esPool[business.IndexCluster]; !ok {
		PromError(fmt.Sprintf("es:集群不存在%s", "LogAuditGroupBy"), "s.dao.LogAuditGroupBy indexName:%s", "LogAuditGroupBy")
		return
	}
	query := d.getQuery(pr, business.Mapping)
	indexName, err = d.logIndexName(c, sp, business)
	if err != nil {
		log.Error("d.LogAuditGroupBy.logIndexName(%v)(%v)", err, indexName)
		return
	}
	if indexName == "" {
		return
	}
	collapse := elastic.NewCollapseBuilder(group).MaxConcurrentGroupRequests(1)
	searchResult, err = d.esPool[business.IndexCluster].Search().Index(indexName).Type("base").Query(query).
		Sort("ctime", false).Collapse(collapse).Size(1000).Do(c)
	if err != nil {
		log.Error("d.LogAuditGroupBy(%v)", err)
		return
	}
	for _, hit := range searchResult.Hits.Hits {
		var t json.RawMessage
		err = json.Unmarshal(*hit.Source, &t)
		if err != nil {
			log.Error("es:%s 返回不是json!!!", business.IndexCluster)
			return
		}
		res.Result = append(res.Result, t)
	}
	res.Page.Ps = sp.Bsp.Ps
	res.Page.Pn = sp.Bsp.Pn
	res.Page.Total = int64(len(res.Result))
	return
}

// LogAuditDelete .
func (d *Dao) LogAuditDelete(c context.Context, pr map[string][]interface{}, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	var (
		indexName    string
		searchResult *elastic.BulkIndexByScrollResponse
	)
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	query := d.getQuery(pr, business.Mapping)
	indexName, err = d.logIndexName(c, sp, business)
	if err != nil {
		log.Error("d.LogAuditDelete.logIndexName(%v)(%v)", err, indexName)
		return
	}
	if indexName == "" {
		return
	}
	searchResult, err = d.esPool[business.IndexCluster].DeleteByQuery().Index(indexName).Type("base").Query(query).Size(10000).Do(c)
	if err != nil {
		log.Error("d.LogAuditDelete.DeleteByQuery(%v)(%v)", err, indexName)
		return
	}
	res.Page.Total = searchResult.Total
	return
}

// LogUserAction .
func (d *Dao) LogUserAction(c context.Context, pr map[string][]interface{}, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	var indexName string
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	query := d.getQuery(pr, business.Mapping)
	indexName, err = d.logIndexName(c, sp, business)
	if err != nil {
		log.Error("d.LogUserAction.logIndexName(%v)(%v)", err, indexName)
		return
	}
	if indexName == "" {
		return
	}
	if res, err = d.searchResult(c, business.IndexCluster, indexName, query, sp.Bsp); err != nil {
		PromError(fmt.Sprintf("es:%s ", sp.Bsp.AppID), "%v", err)
	}
	return
}

// LogUserActionDelete .
func (d *Dao) LogUserActionDelete(c context.Context, pr map[string][]interface{}, sp *model.LogParams, business *model.Business) (res *model.SearchResult, err error) {
	var (
		indexName    string
		searchResult *elastic.BulkIndexByScrollResponse
	)
	res = &model.SearchResult{
		Result: []json.RawMessage{},
		Page:   &model.Page{},
	}
	query := d.getQuery(pr, business.Mapping)
	indexName, err = d.logIndexName(c, sp, business)
	if err != nil {
		log.Error("d.LogUserActionDelete.logIndexName(%v)(%v)", err, indexName)
		return
	}
	if indexName == "" {
		return
	}
	searchResult, err = d.esPool[business.IndexCluster].DeleteByQuery().Index(indexName).Type("base").Query(query).Size(10000).Do(c)
	if err != nil {
		log.Error("d.LogUserActionDelete.DeleteByQuery(%v)(%v)", err, indexName)
		return
	}
	res.Page.Total = searchResult.Total
	return
}

// UDepTs .
func (d *Dao) UDepTs(c context.Context, uids []string) (res *model.UDepTsData, err error) {
	params := url.Values{}
	params.Set("uids", strings.Join(uids, ","))
	if err = d.client.Get(c, d.managerDep, "", params, &res); err != nil {
		err = errors.Wrapf(err, "d.httpSearch url(%s)", d.managerDep+"?"+params.Encode())
		log.Error("d.httpSearch url(%s)", d.managerDep+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "response url(%s) code(%d)", d.managerDep+"?"+params.Encode(), res.Code)
		log.Error("response url(%s) code(%d)", d.managerDep+"?"+params.Encode(), res.Code)
		return
	}
	return
}

// IP .
func (d *Dao) IP(c context.Context, ip []string) (res *model.IPData, err error) {
	params := url.Values{}
	params.Set("ips", strings.Join(ip, ","))
	if err = d.client.Get(c, d.managerIP, "", params, &res); err != nil {
		err = errors.Wrapf(err, "d.httpSearch url(%s)", d.managerIP+"?"+params.Encode())
		log.Error("d.httpSearch url(%s)", d.managerIP+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "response url(%s) code(%d)", d.managerDep+"?"+params.Encode(), res.Code)
		log.Error("response url(%s) code(%d)", d.managerIP+"?"+params.Encode(), res.Code)
		return
	}
	return
}

// LogCount .
func (d *Dao) LogCount(c context.Context, name string, business int, uid interface{}) {
	date := time.Now().Format("2006-01-02")
	if _, err := d.db.Exec(c, _count, name+"_access", date); err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	if _, err := d.db.Exec(c, _percent, name+"_uid", date, uid); err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
	if _, err := d.db.Exec(c, _percent, name+"_business", date, business); err != nil {
		log.Error("d.db.Exec err(%v)", err)
		return
	}
}
