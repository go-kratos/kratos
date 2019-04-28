package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/apm/model/monitor"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// AddMonitor get monitor data and insert
func (s *Service) AddMonitor(c context.Context) (err error) {
	var (
		rpc       = make([]*monitor.Monitor, 0)
		http      = make([]*monitor.Monitor, 0)
		tcs       = make([]*monitor.Monitor, 0)
		dabs      = make([]*monitor.Monitor, 0)
		key       = make([]string, 0)
		value     = make([]interface{}, 0)
		mt        = &monitor.Monitor{}
		insertSQL = "INSERT INTO `monitor`(`app_id`, `interface`, `count`, `cost`) VALUES %s"
	)
	if rpc, err = s.RPCMonitor(c); err != nil {
		return
	}
	if http, err = s.HTTPMonitor(c); err != nil {
		return
	}
	if mt, err = s.Members(c); err != nil {
		return
	}
	if tcs, err = s.TenCent(c); err != nil {
		return
	}
	if dabs, err = s.DataBus(c); err != nil {
		return
	}
	rpc = append(rpc, http...)
	rpc = append(rpc, dabs...)
	rpc = append(rpc, tcs...)
	rpc = append(rpc, mt)
	if len(rpc) == 0 {
		err = fmt.Errorf("monitor data empty")
		log.Error("s.AddMonitor error(%v)", err)
		return
	}
	for _, mt := range rpc {
		key = append(key, "(?, ?, ?, ?)")
		value = append(value, mt.AppID, mt.Interface, mt.Count, mt.Cost)
	}
	return s.DB.Model(&monitor.Monitor{}).Exec(fmt.Sprintf(insertSQL, strings.Join(key, ",")), value...).Error
}

// AppNameList return app list
func (s *Service) AppNameList(c context.Context) (apps []string) {
	return s.c.Apps.Name
}

// PrometheusList return PrometheusList data
func (s *Service) PrometheusList(c context.Context, app, method, mType string) (ret *monitor.MoniRet, err error) {
	var (
		mt  = &monitor.Monitor{}
		mts = make([]*monitor.Monitor, 0)
	)
	if err = s.DB.Select("interface,count,cost,mtime").Where("app_id = ?", app+"-"+method).Group("mtime,interface").Order("interface,mtime").Find(&mts).Error; err != nil {
		log.Error("s.PrometheusList query all error(%v)", err)
		return
	}
	if len(mts) < 1 {
		return
	}
	if err = s.DB.Where("app_id = ?", app+"-"+method).First(mt).Error; err != nil {
		log.Error("s.Prometheus query first error(%v)", err)
		return
	}
	return merge(s.packing(mts), s.times(mt.MTime), mType), err
}

// BroadCastList return BroadCastList data
func (s *Service) BroadCastList(c context.Context) (ret *monitor.MoniRet, err error) {
	var (
		mt  = &monitor.Monitor{}
		mts = make([]*monitor.Monitor, 0)
	)
	if err = s.DB.Select("substring_index(interface, '_', -1) as temp_name,sum(count) as count,mtime").Where("interface REGEXP 'ESTAB|InBound|OutBound|InPacket|OutPacket$'").Group("mtime,temp_name").Order("temp_name,mtime").Find(&mts).Error; err != nil {
		log.Error("s.BroadCastList query  error(%v)", err)
		return
	}
	if err = s.DB.Where("interface REGEXP 'ESTAB|InBound|OutBound|InPacket|OutPacket$'").First(mt).Error; err != nil {
		log.Error("s.BroadCastList query first error(%v)", err)
		return
	}
	return merge(s.packing(mts), s.times(mt.MTime), "count"), err
}

// DataBusList return DataBusList data
func (s *Service) DataBusList(c context.Context) (ret *monitor.MoniRet, err error) {
	var (
		mts = make([]*monitor.Monitor, 0)
		mt  = &monitor.Monitor{}
	)
	if err = s.DB.Select("interface,count,mtime").Where("app_id=?", "kafka-databus").Group("mtime,interface").Order("interface,mtime").Find(&mts).Error; err != nil {
		log.Error("s.MonitorList query error(%v)", err)
		return
	}
	if len(mts) < 1 {
		return
	}
	if err = s.DB.Where("app_id=?", "kafka-databus").First(mt).Error; err != nil {
		log.Error("s.DataBusList query first error(%v)", err)
		return
	}
	return merge(s.packing(mts), s.times(mt.MTime), "count"), err
}

// OnlineList return online data
func (s *Service) OnlineList(c context.Context) (ret *monitor.MoniRet, err error) {
	var (
		mts = make([]*monitor.Monitor, 0)
		mt  = &monitor.Monitor{}
	)
	if err = s.DB.Select("interface,count,mtime").Where("app_id=?", "online").Find(&mts).Error; err != nil {
		log.Error("s.OnlineList query error(%v)", err)
	}
	if len(mts) < 1 {
		return
	}
	if err = s.DB.Where("app_id=?", "online").First(mt).Error; err != nil {
		log.Error("s.OnlineList query error(%v)", err)
	}
	return merge(s.packing(mts), s.times(mt.MTime), "count"), err
}

// merge .
func merge(dts []*monitor.Data, strs []string, mType string) (ret *monitor.MoniRet) {
	items := make([]*monitor.Items, 0)
	for _, dt := range dts {
		var (
			yAxis []int64
			item  = &monitor.Items{}
		)
		if mType == "count" {
			yAxis = formatArray(strs, dt.Times, dt.Counts)
		} else {
			yAxis = formatArray(strs, dt.Times, dt.Costs)
		}
		item.Interface = dt.Interface
		item.YAxis = yAxis
		items = append(items, item)
	}
	if len(items) > 0 {
		ret = &monitor.MoniRet{
			XAxis: strs,
			Items: items,
		}
	}
	return
}

// formatArray formatArray missing data by time
func formatArray(strs, times []string, counts []int64) []int64 {
	var newCounts []int64
	for _, str := range strs {
		if len(counts) < 1 {
			break
		}
		if ok, index := inArray(times, str); ok {
			newCounts = append(newCounts, counts[index])
		} else {
			newCounts = append(newCounts, 0)
		}
	}
	return newCounts
}

// inArray check key in array or not and return position
func inArray(arrays []string, key string) (bool, int) {
	for index, arr := range arrays {
		if key == arr {
			return true, index
		}
	}
	return false, 0
}

// times return a standard string time array
func (s *Service) times(t xtime.Time) []string {
	var (
		nextDay  string
		tList    = []string{t.Time().Format("2006-01-02")}
		curDay   = time.Now().Format("2006-01-02")
		nextTime = t.Time().Add(24 * time.Hour)
	)
	for {
		nextDay = nextTime.Format("2006-01-02")
		tList = append(tList, nextDay)
		nextTime = nextTime.Add(24 * time.Hour)
		if nextDay == curDay {
			break
		}
	}
	return tList
}

// packing .
func (s *Service) packing(mts []*monitor.Monitor) (data []*monitor.Data) {
	d := &monitor.Data{}
	for k, mt := range mts {
		if mt.Interface == "" {
			mt.Interface = mt.TempName
		}
		if d.Interface != mt.Interface {
			if d.Interface != "" {
				data = append(data, d)
			}
			d = &monitor.Data{
				Interface: mt.Interface,
				Counts:    []int64{mt.Count},
				Costs:     []int64{mt.Cost},
				Times:     []string{mt.MTime.Time().Format("2006-01-02")},
			}
			continue
		}
		d.Counts = append(d.Counts, mt.Count)
		d.Costs = append(d.Costs, mt.Cost)
		d.Times = append(d.Times, mt.MTime.Time().Format("2006-01-02"))
		if k == len(mts)-1 {
			data = append(data, d)
		}
	}
	return
}
