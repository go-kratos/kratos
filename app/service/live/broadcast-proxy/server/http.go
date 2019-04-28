package server

import (
	"context"
	"encoding/json"
	"go-common/library/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type BroadcastService struct {
	wg       sync.WaitGroup
	server   *http.Server
	proxy    *BroadcastProxy
	dispatch *CometDispatcher
}

func NewBroadcastService(addr string, proxy *BroadcastProxy, dispatch *CometDispatcher) (*BroadcastService, error) {
	service := &BroadcastService{
		proxy:    proxy,
		dispatch: dispatch,
	}
	service.wg.Add(1)
	go func() {
		defer service.wg.Done()
		service.httpServerProcess(addr)
	}()
	return service, nil
}

func (service *BroadcastService) Close() {
	service.server.Shutdown(context.Background())
}

func (service *BroadcastService) httpServerProcess(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", service.Proxy)
	mux.HandleFunc("/monitor/ping", service.Ping)
	mux.HandleFunc("/dm/x/internal/v1/dispatch", service.Dispatch)
	mux.HandleFunc("/dm/x/internal/v1/set_angry_value", service.SetAngryValue)

	service.server = &http.Server{Addr: addr, Handler: mux}
	service.server.SetKeepAlivesEnabled(true)
	if err := service.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func writeJsonResult(w http.ResponseWriter, r *http.Request, begin time.Time, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("[Http] write result json.Marshal:%v error:%v", v, err)
		return
	}
	if _, err := w.Write([]byte(data)); err != nil {
		log.Error("[Http] write result socket error:%v", err)
		return
	}
	end := time.Now()
	log.Info("request %s, response:%s, time cost:%s", r.RequestURI, data, end.Sub(begin).String())
}

func (service *BroadcastService) Proxy(w http.ResponseWriter, r *http.Request) {
	service.proxy.HandleRequest(w, r)
}

func (service *BroadcastService) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (service *BroadcastService) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var result struct {
		Code int `json:"code"`
		Data struct {
			DanmakuServer []string `json:"dm_server"`
			DanmakuHost   []string `json:"dm_host"`
		} `json:"data"`
	}
	defer writeJsonResult(w, r, time.Now(), &result)
	ip := r.URL.Query().Get("ip")
	uid, _ := strconv.ParseInt(r.URL.Query().Get("uid"), 10, 64)
	result.Code = 0
	result.Data.DanmakuServer, result.Data.DanmakuHost = service.dispatch.Dispatch(ip, uid)
}

func (service *BroadcastService) SetAngryValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 气人值的workaround
	//go func() {
	responseCollection, errCollection := service.proxy.RequestAllBackend(r.Method, "/dm/1/num/change", requestBody)
	for i := range responseCollection {
		if errCollection[i] != nil {
			log.Error("SetAngryValue server:%d error:%+v", i, errCollection[i])
		} else {
			log.Info("SetAngryValue server:%d result:%s", i, responseCollection[i])
		}
	}
	//}()
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(map[string]interface{}{"ret": 1})
	w.Write(response)
	return
}

func (service *BroadcastService) SetAngryValueV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//go func() {
	responseCollection, errCollection := service.proxy.RequestAllBackend(r.Method,
		"/dm/x/internal/v2/set_angry_value", requestBody)
	for i := range responseCollection {
		if errCollection[i] != nil {
			log.Error("SetAngryValueV2 server:%d error:%+v", i, errCollection[i])
			w.WriteHeader(http.StatusServiceUnavailable)
			response, _ := json.Marshal(map[string]interface{}{"code": -1, "msg": errCollection[i].Error()})
			w.Write(response)
			return
		}
		var result struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(responseCollection[i]), &result); err != nil {
			log.Error("SetAngryValueV2 server:%d response:%s", i, responseCollection[i])
			w.WriteHeader(http.StatusServiceUnavailable)
			response, _ := json.Marshal(map[string]interface{}{"code": -2, "msg": responseCollection[i]})
			w.Write(response)
			return
		}
		if result.Code != 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseCollection[i]))
			return
		}
	}
	//}()
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(map[string]interface{}{"code": 0})
	w.Write(response)
	return
}
