package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"go-common/app/admin/ep/tapd/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

var rwmLock = new(sync.RWMutex)

// CallHookUrl Call Hook Url.
func (d *Dao) CallHookUrl(c context.Context, URL string, body interface{}) (err error) {
	var (
		req *http.Request
		res = make(map[string]interface{})
	)

	if req, err = d.newRequest(http.MethodPost, URL, body); err != nil {
		log.Error("d.CallHookUrl url(%s) err(%v)", URL, err)
		return
	}

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.CallHookUrl url(%s) err(%v)", URL, err)
	}
	return
}

// CallHookUrlAsForm Call Hook Url As Form.
func (d *Dao) CallHookUrlAsForm(c context.Context, URL string, body map[string]interface{}) (err error) {
	var (
		res = make(map[string]interface{})
		req *http.Request
	)

	data := make(url.Values)
	for mapKey := range body {
		typeKind := reflect.TypeOf(body[mapKey]).Kind()

		switch typeKind {

		case reflect.Int:
			data[mapKey] = []string{fmt.Sprintf("%d", body[mapKey].(int))}
		case reflect.Int64:
			data[mapKey] = []string{fmt.Sprintf("%d", body[mapKey].(int64))}
		case reflect.Float64:
			data[mapKey] = []string{fmt.Sprintf("%.0f", body[mapKey].(float64))}
		case reflect.String:
			data[mapKey] = []string{body[mapKey].(string)}
		default:
			data[mapKey] = []string{fmt.Sprint(body[mapKey])}
		}
	}

	if req, err = http.NewRequest(http.MethodPost, URL, strings.NewReader(data.Encode())); err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = d.httpClient.Do(c, req, &res)
	return
}

// SaveEnableHookURLToCache Save Enable Hook URL To Cache.
func (d *Dao) SaveEnableHookURLToCache() (eventMap map[string][]*model.HookUrl, err error) {
	rwmLock.Lock()
	defer rwmLock.Unlock()

	var urlEvents []*model.UrlEvent

	eventMap = make(map[string][]*model.HookUrl)

	// 获取有效event
	if urlEvents, err = d.QueryURLEventByStatus(model.HookEventStatusEnable); err != nil {
		return
	}

	//倒查url 过滤不存在或无效状态
	for _, urlEvent := range urlEvents {
		var hookURL *model.HookUrl
		if hookURL, err = d.QueryHookURLByID(urlEvent.UrlID); err != nil || hookURL.ID == 0 || hookURL.Status != model.HookURLStatusEnable {
			continue
		}
		eventMap[urlEvent.Event] = append(eventMap[urlEvent.Event], hookURL)
	}

	for eventMapKey := range eventMap {
		item := &memcache.Item{Key: eventMapKey, Object: eventMap[eventMapKey], Flags: memcache.FlagJSON, Expiration: d.expire}
		d.tokenCacheSave(context.Background(), item)
	}

	return
}

// GetEnableHookURLFromCacheAndSaveIfNot Get Enable Hook URL From Cache and save if not.
func (d *Dao) GetEnableHookURLFromCacheAndSaveIfNot(event model.Event) (hookURLs []*model.HookUrl, err error) {
	var (
		conn = d.mc.Get(context.Background())
		item *memcache.Item
	)

	rwmLock.RLock()

	defer rwmLock.RUnlock()
	defer conn.Close()

	if item, err = conn.Get(string(event)); err == nil {
		if err = json.Unmarshal(item.Value, &hookURLs); err != nil {
			log.Error("GetEnableHookURLFromCache json parse error(%v)", err)
		}
		return
	}

	var eventMap map[string][]*model.HookUrl
	if eventMap, err = d.SaveEnableHookURLToCache(); err != nil {
		return
	}
	hookURLs = eventMap[string(event)]
	return

}

// GetEnableHookURLFromCache Get Enable Hook URL From Cache.
func (d *Dao) GetEnableHookURLFromCache(event model.Event) (hookURLs []*model.HookUrl, err error) {
	var (
		conn = d.mc.Get(context.Background())
		item *memcache.Item
	)

	rwmLock.RLock()

	defer rwmLock.RUnlock()
	defer conn.Close()

	if item, err = conn.Get(string(event)); err == nil {
		if err = json.Unmarshal(item.Value, &hookURLs); err != nil {
			log.Error("GetEnableHookURLFromCache json parse error(%v)", err)
		}
		return
	}
	return
}

// GetEnableHookURLFromDB Get Enable Hook URL From DB.
func (d *Dao) GetEnableHookURLFromDB(event model.Event) (hookURLs []*model.HookUrl, err error) {
	var urlEvents []*model.UrlEvent
	if urlEvents, err = d.QueryURLEventByEventAndStatus(string(event), model.HookEventStatusEnable); err != nil {
		return
	}

	for _, urlEvent := range urlEvents {
		var hookURL *model.HookUrl
		if hookURL, err = d.QueryHookURLByID(urlEvent.UrlID); err != nil {
			return
		}
		hookURLs = append(hookURLs, hookURL)
	}
	return

}
