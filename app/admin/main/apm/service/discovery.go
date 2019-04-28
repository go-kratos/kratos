package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"sync"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/discovery"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/naming"

	errgroup "go-common/library/sync/errgroup.v2"
)

// DiscoveryProxy discovery proxy.
func (s *Service) DiscoveryProxy(c context.Context, method, loc string, params url.Values) (data interface{}, err error) {
	var bs json.RawMessage
	uri := conf.Conf.Host.APICo + "/discovery/" + loc
	if loc == "nodes" {
		api := conf.Conf.Discovery.API
		v1 := make(map[string]interface{})
		for _, v := range api {
			sData := make([]*discovery.Addr, 0, 3)
			uri = "http://" + v + "/discovery/" + loc
			bs, err = s.curl(c, method, uri, params)
			if err != nil {
				log.Error("apmSvc.DiscoveryProxy curl url:"+uri+" params:(%v) error(%v)", params.Encode(), err)
				return nil, err
			}
			if err = json.Unmarshal(bs, &sData); err != nil {
				return nil, err
			}
			Others := make([]*discovery.Addr, 0, len(sData))
			tmpData := &discovery.Status{}
			tmpData.Status = 0
			for _, vv := range sData {
				if vv.Addr == v {
					tmpData.Status = vv.Status
					continue
				}
				Others = append(Others, vv)
			}
			tmpData.Others = Others
			v1[v] = tmpData
		}
		data = v1
	} else if loc == "polling" {
		api := conf.Conf.Discovery.API
		v1 := make([]string, 0)
		tmp := make(map[string]struct{})
		for _, v := range api {
			sData := make([]string, 0)
			uri = "http://" + v + "/discovery/" + loc
			bs, err = s.curl(c, method, uri, params)
			if err != nil {
				log.Error("apmSvc.DiscoveryProxy curl url:"+uri+" params:(%v) error(%v)", params.Encode(), err)
				return nil, err
			}
			if err = json.Unmarshal(bs, &sData); err != nil {
				return nil, err
			}
			for _, vv := range sData {
				if _, ok := tmp[vv]; !ok {
					v1 = append(v1, vv)
					tmp[vv] = struct{}{}
				}
			}
		}
		sort.Strings(v1)
		data = v1
	} else {
		data, err = s.curl(c, method, uri, params)
	}
	return
}

func (s *Service) curl(c context.Context, method, uri string, params url.Values) (data json.RawMessage, err error) {
	var res struct {
		Code    int             `json:"code"`
		Data    json.RawMessage `json:"data"`
		Message string          `json:"message"`
	}
	if method == "GET" {
		if err = s.client.Get(c, uri, "", params, &res); err != nil {
			log.Error("apmSvc.DiscoveryProxy get url:"+uri+" params:(%v) error(%v)", params.Encode(), err)
			return
		}
	} else {
		if err = s.client.Post(c, uri, "", params, &res); err != nil {
			log.Error("apmSvc.DiscoveryProxy post url:"+uri+" params:(%v) error(%v)", params.Encode(), err)
			return
		}
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("apmSvc.DiscoveryProxy url:"+uri+" params:(%v) ecode(%v)", params.Encode(), err)
		return
	}
	data = res.Data
	return
}

// DatabusConsumerAddrs databus consumer addrs.
func (s *Service) DatabusConsumerAddrs(c context.Context, group string) (addrs []string, err error) {
	uri := conf.Conf.Host.APICo + "/discovery/fetch"
	params := url.Values{}
	params.Set("env", env.DeployEnv)
	params.Set("appid", "middleware.databus")
	params.Set("status", "1")
	data, err := s.curl(c, "GET", uri, params)
	if err != nil {
		log.Error("curl discovery fetch error:%v", err)
		return
	}
	var res struct {
		ZoneInstances map[string][]*naming.Instance `json:"zone_instances"`
	}
	if err = json.Unmarshal(data, &res); err != nil {
		log.Error("json unmarshal discovery data error:%v", err)
		return
	}
	var (
		wg   errgroup.Group
		lock sync.Mutex
	)
	for _, instances := range res.ZoneInstances {
		for _, instance := range instances {
			for _, addr := range instance.Addrs {
				u, _ := url.Parse(addr)
				if u.Scheme == "http" {
					wg.Go(func(c context.Context) error {
						p := url.Values{}
						p.Set("group", group)
						var as struct {
							Code int      `json:"code"`
							Data []string `json:"data"`
						}
						if e := s.client.Get(context.Background(), fmt.Sprintf("http://%s/databus/consumer/addrs", u.Host), "", p, &as); e != nil {
							log.Error("curl databus consumer addrs error:%v", e)
							return nil
						}
						lock.Lock()
						addrs = append(addrs, as.Data...)
						lock.Unlock()
						return nil
					})
				}
			}
		}
	}
	wg.Wait()
	return
}
