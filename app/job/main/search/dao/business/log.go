package business

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/search/dao"
	"go-common/app/job/main/search/model"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"

	"gopkg.in/olivere/elastic.v5"
)

const _sql = "SELECT id, index_format, index_version, index_cluster, additional_mapping, data_center FROM digger_"

// Log .
type Log struct {
	d                 *dao.Dao
	appid             string
	attrs             *model.Attrs
	databus           *databus.Databus
	infoC             *infoc.Infoc
	infoCField        []string
	mapData           []model.MapData
	commits           map[int32]*databus.Message
	business          map[int]*info
	week              map[int]string
	additionalMapping map[int]map[string]string
	defaultMapping    map[string]string
	mapping           map[int]map[string]string
}

type info struct {
	Format     string
	Cluster    string
	Version    string
	DataCenter int8
}

// NewLog .
func NewLog(d *dao.Dao, appid string) (l *Log) {
	l = &Log{
		d:                 d,
		appid:             appid,
		attrs:             d.AttrPool[appid],
		databus:           d.DatabusPool[appid],
		infoC:             d.InfoCPool[appid],
		infoCField:        []string{},
		mapData:           []model.MapData{},
		commits:           map[int32]*databus.Message{},
		business:          map[int]*info{},
		additionalMapping: map[int]map[string]string{},
		mapping:           map[int]map[string]string{},
		week: map[int]string{
			0: "0107",
			1: "0815",
			2: "1623",
			3: "2431",
		},
	}
	switch appid {
	case "log_audit":
		l.defaultMapping = map[string]string{
			"uname":      "string",
			"uid":        "string",
			"business":   "string",
			"type":       "string",
			"oid":        "string",
			"action":     "string",
			"ctime":      "time",
			"int_0":      "int",
			"int_1":      "int",
			"int_2":      "int",
			"str_0":      "string",
			"str_1":      "string",
			"str_2":      "string",
			"extra_data": "string",
		}
		l.infoCField = []string{"uname", "uid", "business", "type", "oid", "action", "ctime",
			"int_0", "int_1", "int_2", "str_0", "str_1", "str_2", "str_3", "str_4", "extra_data"}
	case "log_user_action":
		l.defaultMapping = map[string]string{
			"mid":        "string",
			"platform":   "string",
			"build":      "string",
			"buvid":      "string",
			"business":   "string",
			"type":       "string",
			"oid":        "string",
			"action":     "string",
			"ip":         "string",
			"ctime":      "time",
			"int_0":      "int",
			"int_1":      "int",
			"int_2":      "int",
			"str_0":      "string",
			"str_1":      "string",
			"str_2":      "string",
			"extra_data": "string",
		}
		l.infoCField = []string{"mid", "platform", "build", "buvid", "business", "type", "oid", "action", "ip", "ctime",
			"int_0", "int_1", "int_2", "str_0", "str_1", "str_2", "extra_data"}
	default:
		log.Error("log appid error(%v)", appid)
		return
	}
	rows, err := d.SearchDB.Query(context.TODO(), _sql+appid)
	if err != nil {
		log.Error("log Query error(%v)", appid)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id                int
			additionalMapping string
		)
		info := &info{}
		if err = rows.Scan(&id, &info.Format, &info.Version, &info.Cluster, &additionalMapping, &info.DataCenter); err != nil {
			log.Error("Log New DB (%v)(%v)", id, err)
			continue
		}
		l.business[id] = info
		if additionalMapping != "" {
			var additionalMappingDict map[string]string
			if err = json.Unmarshal([]byte(additionalMapping), &additionalMappingDict); err != nil {
				log.Error("Log New Json (%v)(%v)", id, err)
				continue
			}
			l.additionalMapping[id] = additionalMappingDict
		}
	}
	for b := range l.business {
		l.mapping[b] = map[string]string{}
		for k, v := range l.defaultMapping {
			l.mapping[b][k] = v
		}
		if a, ok := l.additionalMapping[b]; ok {
			for k, v := range a {
				l.mapping[b][k] = v
			}
		}
	}
	return
}

// Business return business.
func (l *Log) Business() string {
	return l.attrs.Business
}

// InitIndex .
func (l *Log) InitIndex(c context.Context) {
}

// InitOffset .
func (l *Log) InitOffset(c context.Context) {
}

// Offset .
func (l *Log) Offset(c context.Context) {
}

// MapData .
func (l *Log) MapData(c context.Context) (mapData []model.MapData) {
	return l.mapData
}

// Attrs .
func (l *Log) Attrs(c context.Context) (attrs *model.Attrs) {
	return l.attrs
}

// SetRecover .
func (l *Log) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
}

// IncrMessages .
func (l *Log) IncrMessages(c context.Context) (length int, err error) {
	var jErr error
	ticker := time.NewTicker(time.Duration(time.Millisecond * time.Duration(l.attrs.Databus.Ticker)))
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-l.databus.Messages():
			if !ok {
				log.Error("databus: %s binlog consumer exit!!!", l.attrs.Databus)
				break
			}
			l.commits[msg.Partition] = msg
			var result map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(msg.Value))
			decoder.UseNumber()
			if jErr = decoder.Decode(&result); jErr != nil {
				log.Error("appid(%v) json.Unmarshal(%s) error(%v)", l.appid, msg.Value, jErr)
				continue
			}
			// json.Number转int64
			for k, v := range result {
				switch t := v.(type) {
				case json.Number:
					if result[k], jErr = t.Int64(); jErr != nil {
						log.Error("appid(%v) log.bulkDatabusData.json.Number(%v)(%v)", l.appid, t, jErr)
					}
				}
			}
			l.mapData = append(l.mapData, result)
			if len(l.mapData) < l.attrs.Databus.AggCount {
				continue
			}
		case <-ticker.C:
		}
		break
	}
	// todo: 额外的参数
	length = len(l.mapData)
	return
}

// AllMessages .
func (l *Log) AllMessages(c context.Context) (length int, err error) {
	return
}

// BulkIndex .
func (l *Log) BulkIndex(c context.Context, start, end int, writeEntityIndex bool) (err error) {
	partData := l.mapData[start:end]
	if err = l.bulkDatabusData(c, l.attrs, writeEntityIndex, partData...); err != nil {
		log.Error("appid(%v) json.bulkDatabusData error(%v)", l.appid, err)
		return
	}
	return
}

// Commit .
func (l *Log) Commit(c context.Context) (err error) {
	for k, msg := range l.commits {
		if err = msg.Commit(); err != nil {
			log.Error("appid(%v) Commit error(%v)", l.appid, err)
			continue
		}
		delete(l.commits, k)
	}
	l.mapData = []model.MapData{}
	return
}

// Sleep .
func (l *Log) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(l.attrs.Other.Sleep))
}

// Size .
func (l *Log) Size(c context.Context) (size int) {
	return l.attrs.Other.Size
}

func (l *Log) bulkDatabusData(c context.Context, attrs *model.Attrs, writeEntityIndex bool, bulkData ...model.MapData) (err error) {
	var (
		request     elastic.BulkableRequest
		bulkRequest map[string]*elastic.BulkService
		businessID  int
	)
	bulkRequest = map[string]*elastic.BulkService{}
	for _, b := range bulkData {
		indexName := ""
		if business, ok := b["business"].(int64); ok {
			businessID = int(business)
			if v, ok := b["ctime"].(string); ok {
				if cTime, timeErr := time.Parse("2006-01-02 15:04:05", v); timeErr == nil {
					if info, ok := l.business[businessID]; ok {
						suffix := strings.Replace(cTime.Format(info.Format), "week", l.week[cTime.Day()/8], -1) + "_" + info.Version
						if !writeEntityIndex {
							indexName = attrs.Index.IndexAliasPrefix + "_" + strconv.Itoa(businessID) + "_" + suffix
						} else {
							indexName = attrs.Index.IndexEntityPrefix + "_" + strconv.Itoa(businessID) + "_" + suffix
						}
					}
				}
			}
		}
		if indexName == "" {
			log.Error("appid(%v) ac.d.bulkDatabusData business business(%v) data(%+v)", l.appid, b["business"], b)
			continue
		}
		esCluster := l.business[businessID].Cluster // 上方已经判断l.business[businessID]是否存在
		if _, ok := bulkRequest[esCluster]; !ok {
			if _, eok := l.d.ESPool[esCluster]; eok {
				bulkRequest[esCluster] = l.d.ESPool[esCluster].Bulk()
			} else {
				log.Error("appid(%v) ac.d.bulkDatabusData cluster no find error(%v)", l.appid, esCluster)
				continue //忽略这条数据
			}
		}
		//发送数据中心
		if l.business[businessID].DataCenter == 1 {
			arr := make([]interface{}, len(l.infoCField))
			for i, f := range l.infoCField {
				if v, ok := b[f]; ok {
					arr[i] = fmt.Sprintf("%v", v)
				}
			}
			if er := l.infoC.Info(arr...); er != nil {
				log.Error("appid(%v) ac.infoC.Info error(%v)", l.appid, er)
			}
		}
		//数据处理
		for k, v := range b {
			if t, ok := l.mapping[businessID][k]; ok {
				switch t {
				case "int_to_bin":
					if item, ok := v.(int64); ok {
						item := int(item)
						arr := []string{}
						for i := 0; item != 0; i++ {
							if item&1 == 1 {
								arr = append(arr, strconv.Itoa(item&1<<uint(i)))
							}
							item = item >> 1
						}
						b[k] = arr
					} else {
						delete(b, k)
					}
				case "array":
					if arr, ok := v.([]interface{}); ok {
						b[k] = arr
					} else {
						delete(b, k)
					}
				}
			} else {
				delete(b, k)
			}
		}
		request = elastic.NewBulkIndexRequest().Index(indexName).Type(attrs.Index.IndexType).Doc(b)
		bulkRequest[esCluster].Add(request)
	}
	for _, v := range bulkRequest {
		if v.NumberOfActions() == 0 {
			continue
		}
		if _, err = v.Do(c); err != nil {
			log.Error("appid(%s) bulk error(%v)", attrs.AppID, err)
		}
	}
	return
}
