package dao

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go-common/app/job/main/search/model"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"gopkg.in/olivere/elastic.v5"
)

// BulkDatabusData 写入es数据来自databus.
func (d *Dao) BulkDatabusData(c context.Context, attrs *model.Attrs, writeEntityIndex bool, bulkData ...model.MapData) (err error) {
	// TODO 需要兼容
	var (
		request     elastic.BulkableRequest
		bulkRequest = d.ESPool[attrs.ESName].Bulk()
		//indexField  = ""
	)
	//s := strings.Split(attrs.DataSQL.DataIndexSuffix, ";")
	//if len(s) >= 2 {
	//	indexField = strings.Split(s[1], ":")[0]
	//}
	for _, b := range bulkData {
		var (
			indexName string
			strID     string
		)
		if name, ok := b["index_name"]; ok {
			if indexName, ok = name.(string); ok {
				delete(b, "index_name")
			} else {
				log.Error("dao.es.BulkDBData index_name err")
				continue
			}
		} else {
			if !writeEntityIndex {
				indexName, _ = b.Index(attrs)
			} else {
				_, indexName = b.Index(attrs)
			}
		}
		if id, ok := b["index_id"]; ok {
			if strID, ok = id.(string); !ok {
				log.Error("es.BulkDBData.strID(%v)", id)
				continue
			}
		} else {
			if strID, ok = b.StrID(attrs); !ok {
				log.Error("es.BulkDBData.strID")
				continue
			}
		}
		if indexName == "" {
			continue
		}
		for _, v := range attrs.DataSQL.DataIndexRemoveFields {
			delete(b, v)
		}
		if _, ok := b["index_field"]; ok {
			delete(b, "index_field")
			//delete(b, indexField)
			delete(b, "ctime")
			delete(b, "mtime")
		}
		for k := range b {
			if !d.Contain(k, attrs.DataSQL.DataIndexFormatFields) {
				delete(b, k)
			}
		}
		key := []string{}
		for k := range b {
			key = append(key, k)
		}
		for _, k := range key {
			customType, ok := attrs.DataSQL.DataIndexFormatFields[k]
			if ok {
				switch customType {
				case "ip":
					switch b[k].(type) {
					case float64:
						ipFormat := b.InetNtoA(int64(b[k].(float64)))
						b[k+"_format"] = ipFormat
					case int64:
						ipFormat := b.InetNtoA(b[k].(int64))
						b[k+"_format"] = ipFormat
					}
				case "arr":
					var arr []int
					binaryAttributes := strconv.FormatInt(b[k].(int64), 2)
					for i := len(binaryAttributes) - 1; i >= 0; i-- {
						b := fmt.Sprintf("%c", binaryAttributes[i])
						if b == "1" {
							arr = append(arr, len(binaryAttributes)-i)
						}
					}
					b[k+"_format"] = arr
				case "bin":
					var arr []int
					binaryAttributes := strconv.FormatInt(b[k].(int64), 2)
					for i := len(binaryAttributes) - 1; i >= 0; i-- {
						b := fmt.Sprintf("%c", binaryAttributes[i])
						if b == "1" {
							arr = append(arr, len(binaryAttributes)-i)
						}
					}
					b[k] = arr
				case "workflow":
					if state, ok := b[k].(int64); ok {
						b["state"] = state & 15
						b["business_state"] = state >> 4 & 15
						delete(b, k)
					}
				case "time":
					if v, ok := b[k].(string); ok {
						if v == "0000-00-00 00:00:00" {
							b[k] = "0001-01-01 00:00:00"
						}
					}
				default:
					// as long as you happy
				}
			}
		}
		if strID == "" {
			request = elastic.NewBulkIndexRequest().Index(indexName).Type(attrs.Index.IndexType).Doc(b)
		} else {
			request = elastic.NewBulkUpdateRequest().Index(indexName).Type(attrs.Index.IndexType).Id(strID).Doc(b).DocAsUpsert(true)
		}
		//fmt.Println(request)
		bulkRequest.Add(request)
	}
	if bulkRequest.NumberOfActions() == 0 {
		return
	}
	now := time.Now()
	// prom.BusinessInfoCount.Add("redis:bulk:doc", int64(bulkRequest.NumberOfActions()))
	for i := 0; i < bulkRequest.NumberOfActions(); i++ {
		prom.BusinessInfoCount.Incr("redis:bulk:doc")
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("appid(%s) bulk error(%v)", attrs.AppID, err)
	}
	prom.LibClient.Timing("redis:bulk", int64(time.Since(now)/time.Millisecond))
	return
}

// BulkDBData 写入es数据来自db.
func (d *Dao) BulkDBData(c context.Context, attrs *model.Attrs, writeEntityIndex bool, bulkData ...model.MapData) (err error) {
	var (
		indexName   string
		strID       string
		request     elastic.BulkableRequest
		bulkRequest = d.ESPool[attrs.ESName].Bulk()
	)
	for _, b := range bulkData {
		if name, ok := b["index_name"]; ok {
			if indexName, ok = name.(string); ok {
				delete(b, "index_name")
			} else {
				log.Error("dao.es.BulkDBData index_name err")
				continue
			}
		} else {
			if !writeEntityIndex {
				indexName, _ = b.Index(attrs)
			} else {
				_, indexName = b.Index(attrs)
			}
		}
		if id, ok := b["index_id"]; ok {
			if strID, ok = id.(string); !ok {
				log.Error("es.BulkDBData.strID(%v)", id)
				continue
			}
		} else {
			if strID, ok = b.StrID(attrs); !ok {
				log.Error("es.BulkDBData.strID")
				continue
			}
		}
		if indexName == "" || strID == "" {
			continue
		}
		//attr提供要去除掉的字段，不往ES中写
		for _, v := range attrs.DataSQL.DataIndexRemoveFields {
			delete(b, v)
		}
		request = elastic.NewBulkUpdateRequest().Index(indexName).Type(attrs.Index.IndexType).Id(strID).Doc(b).DocAsUpsert(true).RetryOnConflict(3)
		//fmt.Println(request)
		bulkRequest.Add(request)
	}

	if bulkRequest.NumberOfActions() == 0 {
		// 注意这里request格式问题，会引起action为0
		return
	}
	log.Info("insert number is %d", bulkRequest.NumberOfActions())
	now := time.Now()
	// prom.BusinessInfoCount.Add("redis:bulk:doc", int64(bulkRequest.NumberOfActions()))
	for i := 0; i < bulkRequest.NumberOfActions(); i++ {
		prom.BusinessInfoCount.Incr("redis:bulk:doc")
	}
	if _, err = bulkRequest.Do(c); err != nil {
		log.Error("appid(%s) bulk error(%v)", attrs.AppID, err)
	}
	prom.LibClient.Timing("redis:bulk", int64(time.Since(now)/time.Millisecond))
	return
}

// pingEsCluster ping es cluster
func (d *Dao) pingESCluster(ctx context.Context) (err error) {
	//for name, client := range d.ESPool {
	//	if _, _, err = client.Ping(d.c.Es[name].Addr[0]).Do(ctx); err != nil {
	//		d.PromError("Es:Ping", "%s:Ping error(%v)", name, err)
	//		return
	//	}
	//}
	return
}

// Contain .
func (d *Dao) Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
