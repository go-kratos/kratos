package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/log"
)

type BroadcastProxy struct {
	backend         []string
	probePath       string
	reverseProxy    []*httputil.ReverseProxy
	bestClientIndex int32
	probeSample     int
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewBroadcastProxy(backend []string, probePath string, maxIdleConns int, probeSample int) (*BroadcastProxy, error) {
	if len(backend) == 0 {
		return nil, errors.New("Require at least one backend")
	}
	proxy := new(BroadcastProxy)
	if probeSample > 0 {
		proxy.probeSample = probeSample
	} else {
		proxy.probeSample = 1
	}
	for _, addr := range backend {
		proxy.backend = append(proxy.backend, addr)
		proxy.probePath = probePath
		p := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   addr,
		})
		p.Transport = &http.Transport{
			DisableKeepAlives:   false,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConns,
		}
		proxy.reverseProxy = append(proxy.reverseProxy, p)
	}
	proxy.ctx, proxy.cancel = context.WithCancel(context.Background())
	proxy.wg.Add(1)
	go func() {
		defer proxy.wg.Done()
		proxy.mainProbeProcess()
	}()

	return proxy, nil
}

func (proxy *BroadcastProxy) Close() {
	proxy.cancel()
	proxy.wg.Wait()
}

func (proxy *BroadcastProxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	i := atomic.LoadInt32(&proxy.bestClientIndex)
	proxy.reverseProxy[i].ServeHTTP(w, r)
	t2 := time.Now()
	log.V(3).Info("proxy process req:%s,backend id:%d, timecost:%s", r.RequestURI, i, t2.Sub(t1).String())
}

func (proxy *BroadcastProxy) RequestAllBackend(method, uri string, requestBody []byte) ([]string, []error) {
	responseCollection := make([]string, len(proxy.reverseProxy))
	errCollection := make([]error, len(proxy.reverseProxy))
	var wg sync.WaitGroup
	for i := range proxy.reverseProxy {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			req, err := http.NewRequest(method, uri, bytes.NewReader(requestBody))
			if err != nil {
				errCollection[index] = err
				return
			}
			httpRecorder := httptest.NewRecorder()
			proxy.reverseProxy[index].ServeHTTP(httpRecorder, req)
			resp := httpRecorder.Result()
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				errCollection[index] = errors.New(fmt.Sprintf("http response:%s", resp.Status))
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCollection[index] = err
				return
			}
			responseCollection[index] = string(body)
		}(i)
	}
	wg.Wait()
	return responseCollection, errCollection
}

func (proxy *BroadcastProxy) mainProbeProcess() {
	var wg sync.WaitGroup
	interval := make([]time.Duration, len(proxy.backend))
	for {
		fast := -1
		for i, probe := range proxy.reverseProxy {
			wg.Add(1)
			go func(p *httputil.ReverseProxy, index int) {
				defer wg.Done()
				interval[index] = proxy.unitProbeProcess(p, index)
			}(probe, i)
		}
		wg.Wait()
		for i, v := range interval {
			if fast < 0 {
				fast = i
			} else {
				if v < interval[fast] {
					fast = i
				}
			}
		}
		atomic.StoreInt32(&proxy.bestClientIndex, int32(fast))
		log.Info("[probe result]best server id:%d,addr:%s", fast, proxy.backend[fast])
		for i, d := range interval {
			log.Info("[probe log]server id:%d,addr:%s,avg time cost:%fms", i, proxy.backend[i], 1000*d.Seconds())
		}
		select {
		case <-time.After(time.Second):
		case <-proxy.ctx.Done():
			return
		}
	}
}

func (proxy *BroadcastProxy) unitProbeProcess(p *httputil.ReverseProxy, backendIndex int) time.Duration {
	var (
		wg       sync.WaitGroup
		duration int64
	)
	for i := 0; i < proxy.probeSample; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			timeout := time.Second
			<-time.After(time.Duration(rand.Intn(proxy.probeSample)) * time.Millisecond)
			if timeCost, err := proxy.checkMonitorPing(p, backendIndex, timeout); err == nil {
				atomic.AddInt64(&duration, timeCost.Nanoseconds())
			} else {
				atomic.AddInt64(&duration, timeout.Nanoseconds())
			}
		}()
	}
	wg.Wait()
	return time.Duration(duration/int64(proxy.probeSample)) * time.Nanosecond
}

func (proxy *BroadcastProxy) checkMonitorPing(p *httputil.ReverseProxy, backendIndex int, timeout time.Duration) (time.Duration, error) {
	req, err := http.NewRequest("GET", proxy.probePath, nil)
	if err != nil {
		return timeout, err
	}
	ctx, cancel := context.WithTimeout(proxy.ctx, timeout)
	defer cancel()
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	beginTime := time.Now()
	p.ServeHTTP(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("probe:server id:%d,addr:%s,send requset error:%s", backendIndex,
			proxy.backend[backendIndex], resp.Status)
		return timeout, errors.New("http response:" + resp.Status)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("probe:server id:%d,addr:%s,read response error:%v", backendIndex,
			proxy.backend[backendIndex], err)
		return timeout, err
	}
	if !bytes.Equal(result, []byte("pong")) {
		log.Error("probe:server id:%d,addr:%s,not match response:%s", backendIndex,
			proxy.backend[backendIndex], result)
		return timeout, err
	}
	endTime := time.Now()
	return endTime.Sub(beginTime), nil
}
