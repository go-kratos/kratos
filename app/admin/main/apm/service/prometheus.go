package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/apm/model/monitor"
	"go-common/library/log"

	"github.com/gogo/protobuf/sortkeys"
)

const (
	countQuery     = "sum(rate(go_%s_server_count{app='%s'}[2m])) by (method)"
	costQuery      = "avg(increase(go_%s_server_sum{app='%s'}[5m]) >0) by (method)"
	inPacketQuery  = "irate(node_network_receive_packets{instance_name=~'%s', device!~'^(lo|bond).*'}[5m]) or irate(node_network_receive_packets_total{instance_name='~%s', device!~'^(lo|bond).*'}[5m])"
	outPacketQuery = "irate(node_network_transmit_packets{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m]) or irate(node_network_transmit_packets_total{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m])"
	inBoundQuery   = "irate(node_network_receive_bytes{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m]) * 8 or irate(node_network_receive_bytes_total{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m]) * 8"
	outBoundQuery  = "irate(node_network_transmit_bytes{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m]) * 8 or irate(node_network_transmit_bytes_total{instance_name=~'%s', device!~\"^(lo|bond).*\"}[5m]) * 8"
	tcpStatQuery   = "max(node_tcp_stat{instance_name=~'%s', stat='ESTAB'}) by (stat)"
	producerQuery  = "sum(increase(go_databus_counter{operation='producer_msg_speed'}[2m]))"
	consumerQuery  = "sum(increase(go_databus_counter{operation='consumer_msg_speed'}[2m]))"
)

// PrometheusRes http result
type PrometheusRes struct {
	RetCode int           `json:"RetCode"`
	Data    []*Prometheus `json:"data"`
}

// Prometheus .
type Prometheus struct {
	Metric struct {
		Method string `json:"method"`
	} `json:"metric"`
	Values [][]interface{} `json:"values"`
}

// Max max value
type Max struct {
	K string
	V int64
}

// CommonRes packet or bound result data
type CommonRes struct {
	RetCode int       `json:"RetCode"`
	Data    []*Common `json:"data"`
}

// Common .
type Common struct {
	Metric struct {
		InstanceName string `json:"instance_name"`
	} `json:"metric"`
	Value []interface{} `json:"value"`
}

// auth calc sign
func (s *Service) auth(params url.Values) string {
	var (
		sortKey   = make([]string, 0)
		hash      = sha1.New()
		signature string
		str       string
	)
	for key := range params {
		sortKey = append(sortKey, key)
	}
	sortkeys.Strings(sortKey)
	for _, key := range sortKey {
		str += key + params.Get(key)
	}
	str += s.c.Prometheus.Secret
	hash.Write([]byte(str))
	signature = hex.EncodeToString(hash.Sum(nil))
	return signature
}

// PrometheusProxy Get Prometheus Data
func (s *Service) PrometheusProxy(c context.Context, params url.Values, res interface{}) (err error) {
	var (
		req *http.Request
		uri = s.c.Prometheus.URL
	)
	if req, err = s.client.NewRequest(http.MethodPost, uri, "", params); err != nil {
		log.Error("s.PrometheusProxy.NewRequest error(%v) params(%s)", err, params.Encode())
		return
	}
	if err = s.client.Do(c, req, res); err != nil {
		log.Error("s.PrometheusProxy.Client.Do error(%v) params(%s)", err, params.Encode())
	}
	return
}

// prometheus get prometheus data
func (s *Service) prometheus(c context.Context, method string) (mts []*monitor.Monitor, err error) {
	var (
		sign, ins string
		names     []string
		params    = url.Values{}
	)
	ins, _ = s.GetInterfaces(c)
	mts = make([]*monitor.Monitor, 0)
	params.Add("Action", "GetPromDataRange")
	params.Add("PublicKey", s.c.Prometheus.Key)
	params.Add("DataSource", "app")
	sign = s.auth(params)
	params.Add("Signature", sign)
	date := time.Now().Format("2006-01-02")
	params.Set("Start", date+" 23:00:00")
	params.Set("End", date+" 23:00:10")
	params.Set("Step", "30")
	names = s.c.Apps.Name
	for _, name := range names {
		var (
			costRet  = &PrometheusRes{}
			countRet = &PrometheusRes{}
		)
		params.Set("Query", fmt.Sprintf(costQuery, method, name))
		if err = s.PrometheusProxy(c, params, costRet); err != nil {
			return
		}
		params.Set("Query", fmt.Sprintf(countQuery, method, name))
		if err = s.PrometheusProxy(c, params, countRet); err != nil {
			return
		}
		for _, val := range costRet.Data {
			var (
				count float64
				api   = val.Metric.Method
			)
			if api == "inner.Ping" || len(val.Values) < 1 || len(val.Values[0]) < 1 {
				continue
			}
			cost, _ := strconv.ParseFloat(val.Values[0][1].(string), 64)
			if int64(cost) < s.c.Apps.Max && !strings.Contains(ins, api) {
				continue
			}
			for _, v := range countRet.Data {
				if api == v.Metric.Method {
					count, _ = strconv.ParseFloat(v.Values[0][1].(string), 64)
					break
				}
			}
			mt := &monitor.Monitor{
				AppID:     name + "-" + method,
				Interface: api,
				Count:     int64(count),
				Cost:      int64(cost),
			}
			mts = append(mts, mt)
		}
	}
	return
}

// GetInterfaces .
func (s *Service) GetInterfaces(c context.Context) (string, error) {
	mt := &monitor.Monitor{}
	if err := s.DB.Select("group_concat(interface) as interface").Find(mt).Error; err != nil {
		log.Error("s.GetInterfaces query error(%v)", err)
		return "", err
	}
	return mt.Interface, nil
}

// RPCMonitor get rpc monitor data
func (s *Service) RPCMonitor(c context.Context) ([]*monitor.Monitor, error) {
	return s.prometheus(c, "rpc")
}

// HTTPMonitor get http monitor data
func (s *Service) HTTPMonitor(c context.Context) ([]*monitor.Monitor, error) {
	return s.prometheus(c, "http")
}

// DataBus return DataBus monitor data
func (s *Service) DataBus(c context.Context) (mts []*monitor.Monitor, err error) {
	var (
		sign   string
		params = url.Values{}
		proRet = &CommonRes{}
		conRet = &CommonRes{}
	)
	mts = make([]*monitor.Monitor, 0)
	sign = s.auth(params)
	params.Add("Action", "GetCurrentPromData")
	params.Add("PublicKey", s.c.Prometheus.Key)
	params.Add("DataSource", "app")
	params.Add("Signature", sign)
	params.Set("Query", producerQuery)
	if err = s.PrometheusProxy(c, params, proRet); err != nil {
		return
	}
	mts = append(mts, pack(proRet.Data, "kafka-databus", "producer")...)
	params.Set("Query", consumerQuery)
	if err = s.PrometheusProxy(c, params, conRet); err != nil {
		return
	}
	mts = append(mts, pack(conRet.Data, "kafka-databus", "consumer")...)
	return
}

// TenCent return TenCent monitor data
func (s *Service) TenCent(c context.Context) (mts []*monitor.Monitor, err error) {
	return s.BroadCast(c, s.c.BroadCast.TenCent, "tencent_main")
}

// KingSoft return KingSoft monitor data
func (s *Service) KingSoft(c context.Context) (mts []*monitor.Monitor, err error) {
	return s.BroadCast(c, s.c.BroadCast.KingSoft, "kingsoft_main")
}

// BroadCast return monitor data
func (s *Service) BroadCast(c context.Context, appName []string, app string) (mts []*monitor.Monitor, err error) {
	var (
		sign      string
		names     string
		params    = url.Values{}
		inPacket  = &CommonRes{}
		outPacket = &CommonRes{}
		inBound   = &CommonRes{}
		outBound  = &CommonRes{}
		tcpStat   = &CommonRes{}
	)
	mts = make([]*monitor.Monitor, 0)
	params.Add("Action", "GetCurrentPromData")
	params.Add("PublicKey", s.c.Prometheus.Key)
	sign = s.auth(params)
	params.Add("Signature", sign)
	names = strings.Join(appName, "|")
	params.Set("Query", fmt.Sprintf(inPacketQuery, names, names))
	params.Set("DataSource", app)
	if err = s.PrometheusProxy(c, params, inPacket); err != nil {
		return
	}
	mts = append(mts, pack(inPacket.Data, "", "InPacket")...)
	params.Set("Query", fmt.Sprintf(outPacketQuery, names, names))
	if err = s.PrometheusProxy(c, params, outPacket); err != nil {
		return
	}
	mts = append(mts, pack(outPacket.Data, "", "OutPacket")...)
	params.Set("Query", fmt.Sprintf(inBoundQuery, names, names))
	if err = s.PrometheusProxy(c, params, inBound); err != nil {
		return
	}
	mts = append(mts, pack(inBound.Data, "", "InBound")...)
	params.Set("Query", fmt.Sprintf(outBoundQuery, names, names))
	if err = s.PrometheusProxy(c, params, outBound); err != nil {
		return
	}
	mts = append(mts, pack(outBound.Data, "", "OutBound")...)
	for _, name := range appName {
		params.Set("Query", fmt.Sprintf(tcpStatQuery, name))
		if err = s.PrometheusProxy(c, params, tcpStat); err != nil {
			return
		}
		mts = append(mts, pack(tcpStat.Data, name, "ESTAB")...)
	}
	return
}

// pack .
func pack(pts []*Common, K, V string) (mts []*monitor.Monitor) {
	for _, pt := range pts {
		if K != "" {
			pt.Metric.InstanceName = K
		}
		count, _ := strconv.ParseFloat(pt.Value[1].(string), 64)
		if count == 0 {
			continue
		}
		mt := &monitor.Monitor{
			AppID:     pt.Metric.InstanceName,
			Interface: pt.Metric.InstanceName + "_" + V,
			Count:     int64(count),
		}
		mts = append(mts, mt)
	}
	return mts
}
