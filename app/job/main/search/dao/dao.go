package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/model"
	"go-common/library/xstr"
	// "go-common/database/hbase"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"

	"gopkg.in/olivere/elastic.v5"
)

var errorsCount = prom.BusinessErrCount

const (
	// business

	// search db name. for table attr,offset,manager.
	_searchDB = "search"
)

// App .
type App interface {
	Business() string
	InitIndex(c context.Context)
	InitOffset(c context.Context)
	Offset(c context.Context)
	SetRecover(c context.Context, recoverID int64, recoverTime string, i int)
	IncrMessages(c context.Context) (length int, err error)
	AllMessages(c context.Context) (length int, err error)
	BulkIndex(c context.Context, start, end int, writeEntityIndex bool) (err error)
	Commit(c context.Context) (err error)
	Sleep(c context.Context)
	Size(c context.Context) (size int)
}

// Dao .
type Dao struct {
	c *conf.Config
	// smsClient
	sms *sms
	// search db
	SearchDB *xsql.DB
	// hbase        *hbase.Client
	BusinessPool map[string]model.BsnAppInfo
	AttrPool     map[string]*model.Attrs
	AppPool      map[string]App
	DBPool       map[string]*xsql.DB
	ESPool       map[string]*elastic.Client
	DatabusPool  map[string]*databus.Databus
	InfoCPool    map[string]*infoc.Infoc
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		DBPool: newDbPool(c),
	}
	// check search db
	if d.SearchDB = d.DBPool[_searchDB]; d.SearchDB == nil {
		panic("SearchDB must config")
	}
	d.sms = newSMS(d)
	d.BusinessPool = newBusinessPool(d)
	d.AttrPool = newAttrPool(d)
	d.ESPool = newEsPool(c, d)
	// consumer
	d.DatabusPool = newDatabusPool(c, d)
	d.InfoCPool = newInfoCPool(c, d)
	return
}

// newDatabusPool .
func newDatabusPool(c *conf.Config, d *Dao) (pool map[string]*databus.Databus) {
	pool = make(map[string]*databus.Databus)
	if c.Business.Index {
		return
	}
	for name := range d.BusinessPool {
		if config, ok := c.Databus[name]; ok {
			pool[name] = databus.New(config)
		}
	}
	return
}

// newInfoCPool .
func newInfoCPool(c *conf.Config, d *Dao) (pool map[string]*infoc.Infoc) {
	pool = map[string]*infoc.Infoc{}
	if c.Business.Index {
		return
	}
	for k := range d.BusinessPool {
		if n, ok := c.InfoC[k]; ok {
			pool[k] = infoc.New(n)
		}
	}
	return
}

// newBusinessPool all appid info from one business
func newBusinessPool(d *Dao) (pool map[string]model.BsnAppInfo) {
	pool = map[string]model.BsnAppInfo{}
	if bns, err := newBusiness(d, d.c.Business.Env); err == nil {
		for _, v := range bns.bInfo.AppInfo {
			if v.AppID != "" {
				pool[v.AppID] = v
			}
		}
	}
	return
}

// newAttrPool .
func newAttrPool(d *Dao) (pool map[string]*model.Attrs) {
	pool = make(map[string]*model.Attrs)
	for k := range d.BusinessPool {
		ar := newAttr(d, k)
		pool[k] = ar.attrs
	}
	//fmt.Println("strace:attr-pool>", pool)
	return
}

// SetRecover set recover.
func (d *Dao) SetRecover(c context.Context, appid string, recoverID int64, recoverTime string, i int) {
	d.AppPool[appid].SetRecover(c, recoverID, recoverTime, i)
}

// newDbPool db combo
func newDbPool(c *conf.Config) (pool map[string]*xsql.DB) {
	pool = make(map[string]*xsql.DB)
	for dbName, config := range c.DB {
		pool[dbName] = xsql.NewMySQL(config)
	}
	return
}

// newEsCluster cluster action
func newEsPool(c *conf.Config, d *Dao) (esCluster map[string]*elastic.Client) {
	esCluster = make(map[string]*elastic.Client)
	for esName, e := range c.Es {
		if client, err := elastic.NewClient(elastic.SetURL(e.Addr...)); err == nil {
			esCluster[esName] = client
		} else {
			d.PromError("es:集群连接失败", "cluster: %s, %v", esName, err)
			if err := d.SendSMS(fmt.Sprintf("[search-job]%s集群连接失败", esName)); err != nil {
				d.PromError("es:集群连接短信失败", "cluster: %s, %v", esName, err)
			}
		}
	}
	return
}

// PromError .
func (d *Dao) PromError(name string, format string, args ...interface{}) {
	errorsCount.Incr(name)
	log.Error(format, args)
}

// Close close dao
func (d *Dao) Close() {
	for _, db := range d.DBPool {
		db.Close()
	}
}

// Ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	// TODO 循环ping
	if err = d.SearchDB.Ping(c); err != nil {
		d.PromError("db:ping", "")
		return
	}
	if err = d.pingESCluster(c); err != nil {
		d.PromError("es:ping", "d.pingESCluster error(%v)", err)
		return
	}
	return
}

// GetAliases get all aliases by indexAliasPrefix
func (d *Dao) GetAliases(esName, indexAliasPrefix string) (aliases map[string]bool, err error) {
	aliases = map[string]bool{}
	if _, ok := d.ESPool[esName]; !ok {
		log.Error("GetAliases 集群不存在 (%s)", esName)
		return
	}
	if aliasesRes, err := d.ESPool[esName].Aliases().Index(indexAliasPrefix + "*").Do(context.TODO()); err != nil {
		log.Error("GetAliases(%s*) failed", indexAliasPrefix)
	} else {
		for _, indexDetails := range aliasesRes.Indices {
			for _, v := range indexDetails.Aliases {
				if v.AliasName != "" {
					aliases[v.AliasName] = true
				}
			}
		}
	}
	return
}

// InitIndex create entity indecies & aliases if necessary
func (d *Dao) InitIndex(c context.Context, aliases map[string]bool, esName, indexAliasName, indexEntityName, indexMapping string) {
	if indexMapping == "" {
		log.Error("indexEntityName(%s) mapping is epmty", indexEntityName)
		return
	}
	for {
		exists, err := d.ESPool[esName].IndexExists(indexEntityName).Do(c)
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		if !exists {
			if _, err := d.ESPool[esName].CreateIndex(indexEntityName).Body(indexMapping).Do(c); err != nil {
				log.Error("indexEntityName(%s) create err(%v)", indexEntityName, err)
				time.Sleep(time.Second * 3)
				continue
			}
		}
		break
	}
	// add aliases if necessary
	if aliases != nil && indexAliasName != indexEntityName {
		if _, ok := aliases[indexAliasName]; !ok {
			if _, err := d.ESPool[esName].Alias().Add(indexEntityName, indexAliasName).Do(context.TODO()); err != nil {
				log.Error("indexEntityName(%s) failed to add alias indexAliasName(%s) err(%v)", indexEntityName, indexAliasName, err)
			}
		}
	}
}

// InitOffset init offset to offset table .
func (d *Dao) InitOffset(c context.Context, offset *model.LoopOffset, attrs *model.Attrs, arr []string) {
	for {
		if err := d.bulkInitOffset(c, offset, attrs, arr); err != nil {
			log.Error("project(%s) initOffset(%v)", attrs.AppID, err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
}

// InitMapData init each field struct
func InitMapData(fields []string) (item model.MapData, row []interface{}) {
	item = make(map[string]interface{})
	for _, v := range fields {
		item[v] = new(interface{})
	}
	for _, v := range fields {
		row = append(row, item[v])
	}
	return
}

// UpdateOffsetByMap .
func UpdateOffsetByMap(offsets *model.LoopOffset, mapData ...model.MapData) {
	var (
		id    int64
		mtime string
	)
	length := len(mapData)
	if length == 0 {
		return
	}
	offsetTime := offsets.OffsetTime
	lastRes := mapData[length-1]
	id = lastRes.PrimaryID()
	lastMtime := lastRes.StrMTime()
	//fmt.Println("real", lastMtime, id, offsets.OffsetID)
	if (id != offsets.OffsetID) && (offsetTime == lastMtime) {
		offsets.IsLoop = true
	} else {
		if offsets.IsLoop {
			for _, p := range mapData {
				tempMtime := p.StrMTime()
				if tempMtime == offsetTime {
					continue
				}
				id = p.PrimaryID()
				mtime = tempMtime
				break
			}
		} else {
			mtime = lastMtime
		}
		offsets.IsLoop = false
	}
	offsets.SetTempOffset(id, mtime)
}

// CommitOffset .
func (d *Dao) CommitOffset(c context.Context, offset *model.LoopOffset, appid, tableName string) (err error) {
	if offset.TempOffsetID != 0 {
		offset.SetOffset(offset.TempOffsetID, "")
	}
	if offset.TempOffsetTime != "" {
		offset.SetOffset(0, offset.TempOffsetTime)
	}
	if offset.TempRecoverID >= 0 {
		offset.SetRecoverOffset(offset.TempRecoverID, "")
	}
	if offset.TempRecoverTime != "" {
		offset.SetRecoverOffset(-1, offset.TempRecoverTime)
	}
	err = d.updateOffset(c, offset, appid, tableName)
	return
}

// JSON2map json to map.
func (d *Dao) JSON2map(rowJSON json.RawMessage) (result map[string]interface{}, err error) {
	decoder := json.NewDecoder(bytes.NewReader(rowJSON))
	decoder.UseNumber()
	if err = decoder.Decode(&result); err != nil {
		log.Error("JSON2map.Unmarshal(%s) error(%v)", rowJSON, err)
		return nil, err
	}
	// json.Number转int64
	for k, v := range result {
		switch t := v.(type) {
		case json.Number:
			if result[k], err = t.Int64(); err != nil {
				log.Error("JSON2map.json.Number(%v)(%v)", t, err)
				return nil, err
			}
		}
	}
	return
}

// ExtraData .
func (d *Dao) ExtraData(c context.Context, mapData []model.MapData, attrs *model.Attrs, way string, tags []string) (md []model.MapData, err error) {
	md = mapData
	switch way {
	case "db":
		for i, item := range mapData {
			item.TransData(attrs)
			for k, v := range item {
				md[i][k] = v
			}
		}
	case "dtb":
		for i, item := range mapData {
			item.TransDtb(attrs)
			for k, v := range item {
				md[i][k] = v
			}
		}
	}
	for _, ex := range attrs.DataExtras {
		// db exists or not
		if _, ok := d.DBPool[ex.DBName]; !ok {
			log.Error("ExtraData d.DBPool excludes:%s", ex.DBName)
			continue
		}
		if len(tags) != 0 {
			for _, v := range tags {
				if v != ex.Tag {
					continue
				}
				switch ex.Type {
				case "slice":
					md, err = d.extraDataSlice(c, md, attrs, ex)
				default:
					md, err = d.extraDataDefault(c, md, attrs, ex)
				}
			}
		} else {
			switch ex.Type {
			case "slice":
				md, err = d.extraDataSlice(c, md, attrs, ex)
			default:
				md, err = d.extraDataDefault(c, md, attrs, ex)
			}
		}
	}
	return
}

// extraData-default
func (d *Dao) extraDataDefault(c context.Context, mapData []model.MapData, attrs *model.Attrs, ex model.AttrDataExtra) (md []model.MapData, err error) {
	md = mapData
	// filter ids from in_fields
	var (
		ids     []int64
		items   map[int64]model.MapData
		include []string
	)
	cdtInField := ex.Condition["in_field"]
	items = make(map[int64]model.MapData)
	if cld, ok := ex.Condition["include"]; ok {
		include = strings.Split(cld, "=")
	}
	var rows *xsql.Rows
	if cdtInFields := strings.Split(cdtInField, ","); len(cdtInFields) == 1 { //FIXME 支持主键多个条件定位一条数据
		for _, m := range mapData {
			if v, ok := m[cdtInField]; ok {
				if len(include) >= 2 { //TODO 支持多种
					if cldVal, ok := m[include[0]]; ok && strconv.FormatInt(cldVal.(int64), 10) == include[1] {
						ids = append(ids, v.(int64))
					}
				} else {
					ids = append(ids, v.(int64)) //TODO 加去重
				}
			}
		}
		// query extra data
		//TODO 如果分表太多的业务，单次循环size设置过大一下子来50万的数据，where in一个表会拒绝请求或超时
		if len(ids) > 0 {
			if tableFormat := strings.Split(ex.TableFormat, ","); ex.TableFormat == "" || tableFormat[0] == "single" {
				i := 0
				flag := false
				//TODO 缺点：耗内存
				for {
					var id []int64
					if (i+1)*200 < len(ids) {
						id = ids[i*200 : (i+1)*200]
					} else {
						id = ids[i*200:]
						flag = true
					}
					rows, err = d.DBPool[ex.DBName].Query(c, fmt.Sprintf(ex.SQL, xstr.JoinInts(id))+" and 1 = ? ", 1)
					if err != nil {
						log.Error("extraDataDefault db.Query error(%v)", err)
						return
					}
					for rows.Next() {
						item, row := InitMapData(ex.Fields)
						if err = rows.Scan(row...); err != nil {
							log.Error("extraDataDefault rows.Scan() error(%v)", err)
							continue
						}
						if v, ok := item[ex.InField]; ok {
							if v2, ok := v.(*interface{}); ok {
								item.TransData(attrs)
								items[(*v2).(int64)] = item
							}
						}
						// fmt.Println(item)
					}
					rows.Close()
					i++
					if flag {
						break
					}
				}
			} else if tableFormat[0] == "int" {
				formatData := make(map[int64][]int64)
				var dbid = []int64{}
				if len(tableFormat) >= 6 { // 弹幕举报根据文章id来分表 dmid进行匹配
					for _, m := range mapData {
						if v, ok := m[tableFormat[5]]; ok {
							dbid = append(dbid, v.(int64)) // 加去重
						}
					}
				} else {
					dbid = ids
				}
				if len(dbid) != len(ids) {
					log.Error("tableFormat[5] len error(%v)(%v)", len(dbid), len(ids))
					return
				}
				for i := 0; i < len(ids); i++ {
					d, e := strconv.ParseInt(tableFormat[2], 10, 64)
					if e != nil {
						log.Error("extraDataDefault strconv.Atoi() error(%v)", e)
						continue
					}
					d = dbid[i] % (d + 1)
					if d < 0 { //可能有脏数据
						continue
					}
					formatData[d] = append(formatData[d], ids[i])
				}
				for v, k := range formatData {
					rows, err = d.DBPool[ex.DBName].Query(c, fmt.Sprintf(ex.SQL, v, xstr.JoinInts(k))+" and 1 = ? ", 1)
					if err != nil {
						log.Error("extraDataDefaultTableFormat db.Query error(%v)", err)
						return
					}
					for rows.Next() {
						item, row := InitMapData(ex.Fields)
						if err = rows.Scan(row...); err != nil {
							log.Error("extraDataDefaultTableFormat rows.Scan() error(%v)", err)
							continue
						}
						if v, ok := item[ex.InField]; ok {
							if v2, ok := v.(*interface{}); ok {
								item.TransData(attrs)
								items[(*v2).(int64)] = item
							}
						}
					}
					rows.Close()
				}
			}
		}
		// fmt.Println("ids:", ids, "items:", items)
		// merge data
		for i, m := range mapData {
			if len(include) >= 2 { //TODO 支持多种
				if cldVal, ok := m[include[0]]; !ok || strconv.FormatInt(cldVal.(int64), 10) != include[1] {
					continue
				}
			}
			if k, ok := m[cdtInField]; ok {
				if item, ok := items[k.(int64)]; ok {
					for _, v := range ex.RemoveFields {
						delete(item, v)
					}
					item.TransData(attrs)
					for k, v := range item {
						md[i][k] = v
					}
				}
			}
		}
		//fmt.Println(md)
	} else {
		for i, m := range mapData {
			var value []interface{}
			for _, v := range cdtInFields {
				value = append(value, m[v])
			}
			rows, err = d.DBPool[ex.DBName].Query(c, ex.SQL, value...)
			if err != nil {
				log.Error("extraDataDefault db.Query error(%v)", err)
				return
			}
			for rows.Next() {
				item, row := InitMapData(ex.Fields)
				if err = rows.Scan(row...); err != nil {
					log.Error("extraDataDefault rows.Scan() error(%v)", err)
					continue
				}
				item.TransData(attrs)
				for _, v := range ex.RemoveFields {
					delete(item, v)
				}
				for k, v := range item {
					md[i][k] = v
				}
			}
			rows.Close()
		}
	}
	return
}

// extraData-slice
func (d *Dao) extraDataSlice(c context.Context, mapData []model.MapData, attrs *model.Attrs, ex model.AttrDataExtra) (md []model.MapData, err error) {
	md = mapData
	// filter ids from in_fields
	var (
		ids     []int64
		items   map[string]map[string][]interface{}
		include []string
	)
	cdtInField := ex.Condition["in_field"]
	items = make(map[string]map[string][]interface{})
	sliceFields := strings.Split(ex.SliceField, ",")
	if cld, ok := ex.Condition["include"]; ok {
		include = strings.Split(cld, "=")
	}
	for _, m := range mapData {
		if v, ok := m[cdtInField]; ok {
			if len(include) >= 2 { //TODO 支持多种
				if cldVal, ok := m[include[0]]; ok && strconv.FormatInt(cldVal.(int64), 10) == include[1] {
					ids = append(ids, v.(int64))
				}
			} else {
				ids = append(ids, v.(int64)) //TODO 加去重
			}
		}
	}
	// query extra data
	if len(ids) > 0 {
		var rows *xsql.Rows
		rows, err = d.DBPool[ex.DBName].Query(c, fmt.Sprintf(ex.SQL, xstr.JoinInts(ids))+" and 1 = ? ", 1)
		if err != nil {
			log.Error("extraDataSlice db.Query error(%v)", err)
			return
		}
		for rows.Next() {
			item, row := InitMapData(ex.Fields)
			if err = rows.Scan(row...); err != nil {
				log.Error("extraDataSlice rows.Scan() error(%v)", err)
				continue
			}
			if v, ok := item[ex.InField]; ok {
				if v2, ok := v.(*interface{}); ok {
					var key string
					switch (*v2).(type) {
					case int, int8, int16, int32, int64:
						key = strconv.FormatInt((*v2).(int64), 10)
					case []uint, []uint8, []uint16, []uint32, []uint64:
						key = string((*v2).([]byte))
					}
					for _, sf := range sliceFields {
						if _, ok := items[key]; !ok {
							items[key] = make(map[string][]interface{})
						}
						var res interface{}
						if v3, ok := item[sf].(*interface{}); ok {
							switch (*v3).(type) {
							case []uint, []uint8, []uint16, []uint32, []uint64:
								res = string((*v3).([]byte))
							default:
								res = v3
							}
						}
						items[key][sf] = append(items[key][sf], res)
					}
				}
			}
		}
		rows.Close()
	}
	//log.Info("items:%v", items)
	// merge data
	for i, m := range mapData {
		if len(include) >= 2 { //TODO 支持多种
			if cldVal, ok := m[include[0]]; !ok || strconv.FormatInt(cldVal.(int64), 10) != include[1] {
				continue
			}
		}
		if v, ok := m[cdtInField]; ok {
			if item, ok := items[strconv.FormatInt(v.(int64), 10)]; ok {
				for _, sf := range sliceFields {
					if list, ok := item[sf]; ok {
						md[i][sf] = list
					}
				}
			} else {
				for _, sf := range sliceFields {
					md[i][sf] = []int64{}
				}
			}
		}
	}
	// for _, v := range md {
	// 	log.Info("md:%v", v)
	// }
	return
}

// GetConfig .
func (d *Dao) GetConfig(c context.Context) *conf.Config {
	return d.c
}
