package flowmonitor

import (
	"errors"
	"time"
	"os"
	"strconv"
	"encoding/json"
	"net"

	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor/counter"
	"go-common/app/service/ops/log-agent/event"
)

type FlowMonitor struct {
	conf                *Config
	monitorLogChanSize  int
	currentVer          int
	verMap              map[int]*prometheus.CounterVec
	flowMonitorThrottle bool
	indexName           string
	conn                net.Conn
}

var argsError = errors.New("appId timeRangeKey status must be specified")
var throttledError = errors.New("flow monitor is throttled")
var notInitError = errors.New("flow monitor does not init")
var Fm *FlowMonitor

// InitFlowMonitor init flow monitor
func InitFlowMonitor(conf *Config) (err error) {
	fm := new(FlowMonitor)
	fm.conf = conf
	if err = fm.checkConfig(); err != nil {
		return err
	}
	fm.flowMonitorThrottle = false
	fm.verMap = make(map[int]*prometheus.CounterVec)
	ver1 := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "FlowMonitor1", Help: "help"}, []string{"appId", "timeRangeKey", "source", "kind", "status"})
	prometheus.MustRegister(ver1)
	fm.verMap[0] = ver1
	ver2 := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "FlowMonitor2", Help: "help"}, []string{"appId", "timeRangeKey", "source", "kind", "status"})
	prometheus.MustRegister(ver2)
	fm.verMap[1] = ver2
	fm.currentVer = 0
	fm.newConn()
	go fm.flowmonitorreport()
	Fm = fm
	return nil
}

// getCurrentVer get the ver being used
func (fm *FlowMonitor) getCurrentVer() *prometheus.CounterVec {
	return fm.verMap[fm.currentVer]
}

// rollVer roll to the next ver
func (fm *FlowMonitor) rollVer() {
	fm.currentVer = (fm.currentVer + 1) % 2
}

func (fm *FlowMonitor) AddEvent(e *event.ProcessorEvent, source string, kind string, status string) {
	if len(e.AppId) != 0 {
		fm.Add(string(e.AppId), source, e.TimeRangeKey, kind, status)
	}
	fm.Add(e.LogId, source, e.TimeRangeKey, kind, status)
}

// Add do the metric
func (fm *FlowMonitor) Add(appId string, source string, timeRangeKey string, kind string, status string) (error) {
	if fm == nil {
		return notInitError
	}
	if fm.flowMonitorThrottle {
		return throttledError
	}

	if appId == "" || source == "" || status == "" {
		return argsError
	}

	if timeRangeKey == "" {
		timeRangeKey = strconv.FormatInt(time.Now().Unix()/100*100, 10)
	}

	if counter, err := fm.getCurrentVer().GetMetricWithLabelValues(appId, timeRangeKey, source, kind, status); err != nil {
		return err
	} else {
		counter.Inc()
		return nil
	}
}

// readVec read metrics from one vec
func (fm *FlowMonitor) readVec(ver *prometheus.CounterVec) error {
	if fm.conn == nil {
		if err := fm.newConn(); err != nil {
			return err
		}
	}
	metrics := make(chan prometheus.Metric)
	go func() {
		ver.Collect(metrics)
		close(metrics)
	}()

	hostname, _ := os.Hostname()
	var ignore_metric bool = false
	for {
		select {
		case metric, ok := <-metrics:
			if !ok {
				return nil
			}
			if ignore_metric {
				continue
			}
			data := make(map[string]interface{})
			for _, label := range metric.(prometheus.Counter).Lables() {
				data[*label.Name] = *label.Value
			}
			if timeRangeKey, ok := data["timeRangeKey"]; ok {
				timeint64, _ := strconv.ParseInt(timeRangeKey.(string), 10, 64)
				data["time"] = time.Unix(timeint64, 0).UTC().Format("2006-01-02T15:04:05")
			}
			data["hostname"] = hostname
			data["counter"] = metric.(prometheus.Counter).Value()
			if dataSend, err := json.Marshal(data); err == nil {
				dataSend = append(dataSend, []byte("\n")...)
				n, err := fm.conn.Write(dataSend)
				if err == nil && n < len(dataSend) {
					log.Error("Error: flow monitor write error: short write")
				}
				if err != nil {
					log.Error("Error: flow monitor write error: %v", err)
					fm.conn.Close()
					fm.conn = nil
					// if conn write error, just ignore. ver.Collect must be finished or RLock will not be released
					ignore_metric = true
				}
			}
		}
	}
}

// linkmonitorreport report link monitor data periodicity
func (fm *FlowMonitor) flowmonitorreport() {
	for {
		time.Sleep(time.Duration(fm.conf.Interval))
		currentVer := fm.getCurrentVer()
		fm.rollVer()
		fm.readVec(currentVer)
		currentVer.Reset()
	}
}

// newConn make a conn to logstash(monitor data receiver)
func (fm *FlowMonitor) newConn() error {
	conn, err := net.DialTimeout("tcp", fm.conf.Addr, time.Duration(time.Second*5))
	if err == nil && conn != nil {
		fm.conn = conn
		fm.flowMonitorThrottle = false
		log.Info("init flow monitor conn to: %s", fm.conf.Addr)
		return nil
	} else {
		log.Error("flow monitor conn failed: %s: %v", fm.conf.Addr, err)
		fm.flowMonitorThrottle = true
		return err
	}
}
