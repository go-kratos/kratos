package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"
	"go-common/app/interface/openplatform/monitor-end/model/prom"
	"go-common/library/log"

	"github.com/json-iterator/go"
)

const _regex = `.*\d{2,}`

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	_typeAPP       = "app"
	_typeWeb       = "web/h5"
	_typeAPPH      = "app_h"
	_eventPage     = "page"
	_eventAPI      = "api"
	_eventResource = "resource"
	_levelInfo     = "info"
	_levelWarning  = "warning"
	_levelError    = "error"
	_logTypeOne    = "1"
	_logTypeTwo    = "2"
)

// Report .
func (s *Service) Report(c context.Context, params *model.LogParams, mid int64, ip string, buvid string, userAgent string) (err error) {
	var l *monitor.Log
	log.Info("get log report(%+v)", params)
	if params.IsAPP == 1 {
		_, err = s.nativeLog(c, params.Log)
	} else {
		l, err = s.frontendLog(c, params.Source, params.Log, mid, ip, buvid, userAgent)
		go s.promsFE(context.TODO(), l)
	}
	return
}
func induceSubEvent(s string) string {
	if strings.Contains(s, "?") {
		s = s[:strings.Index(s, "?")]
	}
	if strings.Contains(s, "://") {
		s = s[(strings.Index(s, "://") + 3):]
	}
	if len(s) > 512 {
		s = s[:512]
	}
	var (
		ok  bool
		err error
	)
	for {
		if ok, err = regexp.Match(_regex, []byte(s)); err != nil {
			log.Error("s.Report.regexp error(%+v), data(%s)", err, s)
			return s
		}
		if !ok {
			break
		}
		if !strings.Contains(s, "/") {
			return s
		}
		s = s[:strings.LastIndex(s, "/")]
	}
	return s
}

func (s *Service) promsFE(c context.Context, l *monitor.Log) {
	if !s.c.Prom.Promed || l == nil {
		return
	}
	var (
		ok  bool
		err error
	)
	l.SubEvent = induceSubEvent(l.SubEvent)
	if l.Event == _eventResource {
		regex := `.*i\d{1}.hdslb.com`
		if ok, err = regexp.Match(regex, []byte(l.SubEvent)); err != nil {
			log.Error("s.Report.regexp error(%+v), data(%s)", err, l.SubEvent)
			return
		}
		if ok {
			l.SubEvent = l.SubEvent[:strings.Index(l.SubEvent, "hdslb.com")+9]
		}
		prom.AddCode(l.Type, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.BusinessCode)
		return
	}
	prom.AddHTTPCode(l.Type, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.HTTPCode)
	prom.AddCode(l.Type, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.BusinessCode)
	if l.Details == nil {
		var cost int64
		if l.Duration == "" {
			return
		}
		if cost, err = strconv.ParseInt(l.Duration, 10, 64); err != nil {
			log.Warn("s.Info.ParseInt can not convert duration(%s) to int64", l.Duration)
			return
		}
		prom.AddCommonLog(l.Type, l.SubProduct, l.SubEvent, l.Event, l.Ver, cost)
	} else {
		prom.AddDetailedLog(l.Type, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.Details)
	}
}

func (s *Service) nativeLog(c context.Context, data string) (l *monitor.Log, err error) {
	l = &monitor.Log{}
	if err = json.Unmarshal([]byte(data), l); err != nil {
		log.Error("s.nativeLog.unmarshal error(%+v), data(%s)", err, data)
		return
	}
	l.Type = _typeAPPH
	s.handleLog(c, l)
	return
}

// HandleMsg .
func (s *Service) HandleMsg(msg []byte) {
	l := monitor.LogFromBytes(msg)
	l.Type = _typeAPP
	s.handleLog(context.TODO(), l)
}

func (s *Service) handleLog(c context.Context, l *monitor.Log) {
	var (
		appID   string
		kv      []log.D
		err     error
		logtype int
		isInt   bool
		cost    int64
		ok      bool
		app     = "android"
	)
	l.TraceidSvr = l.Traceid
	l.Traceid = ""
	l.CalCode()
	if appID, _, kv, err = l.LogData(); err != nil {
		log.Error("s.nativeLog.LogData error(%+v), data(%+v)", err, l)
		return
	}
	if l.Result == "1" || l.Result == "" {
		s.mh.Info(c, appID, kv...)
	} else {
		s.mh.Error(c, appID, kv...)
	}
	if l.SubProduct == "" {
		l.SubProduct = l.Product
	}
	if strings.Contains(l.RequestURI, "ios") {
		app = "ios"
		isInt = true
	}
	if logtype, err = bStrToInt(l.LogType, isInt); err != nil {
		return
	}
	if s.c.Prom.IgnoreNA && !s.checkProduct(l) {
		return
	}
	//	丢弃老版本的network日志
	if l.Event == "network" && l.Codes == "" {
		return
	}
	l.SubEvent = induceSubEvent(l.SubEvent)
	regex := `.*i\d{1}.hdslb.com`
	if ok, err = regexp.Match(regex, []byte(l.SubEvent)); err != nil {
		log.Error("s.Report.regexp error(%+v), data(%s)", err, l.SubEvent)
		return
	}
	if ok {
		l.SubEvent = l.SubEvent[:strings.Index(l.SubEvent, "hdslb.com")+9]
	}
	if (logtype & 1) == 1 {
		// 性能日志
		if cost, err = strconv.ParseInt(l.Duration, 10, 64); err != nil {
			log.Warn("s.handleLog.ParseInt can not convert duration(%s) to int64", l.Duration)
			return
		}
		prom.AddCommonLog(app, l.SubProduct, l.SubEvent, l.Event, l.Ver, cost)
	}
	if (logtype >> 2 & 1) == 1 {
		// 成功/失败日志
		if l.Event == "network" {
			// 网络类型 上报业务码和http状态
			prom.AddHTTPCode(app, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.HTTPCode)
			prom.AddCode(app, l.SubProduct, l.SubEvent, l.Event, l.Ver, l.BusinessCode)
		} else {
			// 其他类型 只上报失败成功 result=1 表示成功
			res := "999"
			if l.Result == "1" {
				res = "0"
			}
			prom.AddCode(app, l.SubProduct, l.SubEvent, l.Event, l.Ver, res)
		}
	}
}

func (s *Service) frontendLog(c context.Context, source string, data string, mid int64, ip string, buvid string, userAgent string) (l *monitor.Log, err error) {
	var (
		ms                []interface{}
		kv                []log.D
		logtype, loglevel string
		url, query        string
	)
	if err = json.Unmarshal([]byte(data), &ms); err != nil {
		log.Error("s.frontendLog.unmarshal error(%+v), data(%s)", err, data)
		return
	}
	for _, m := range ms {
		switch m := m.(type) {
		case map[string]interface{}:
			logtype = stringValueByKey(m, "1", "logtype")
			loglevel = stringValueByKey(m, "info", "level")
			tmp := stringValueByKey(m, "", "url")
			if strings.Contains(tmp, "?") {
				url = strings.Split(tmp, "?")[0]
				query = strings.Split(tmp, "?")[1]
			} else {
				url = tmp
			}
			switch logtype {
			case _logTypeOne:
				l = &monitor.Log{
					LogType:      _logTypeOne,
					Type:         _typeWeb,
					Product:      source,
					IP:           ip,
					Buvid:        buvid,
					UserAgent:    userAgent,
					Event:        _eventAPI,
					Mid:          strconv.FormatInt(mid, 10),
					SubEvent:     url,
					Query:        query,
					Duration:     stringValueByKey(m, "", "cost"),
					TraceidSvr:   stringValueByKey(m, "", "traceid_svr"),
					TraceidEnd:   stringValueByKey(m, "", "traceid_end"),
					HTTPCode:     stringValueByKey(m, "200", "status"),
					BusinessCode: stringValueByKey(m, "0", "errno", "code"),
					SubProduct:   stringValueByKey(m, source, "sub_product"),
					Result:       "0",
				}
				srcURL := stringValueByKey(m, "", "srcUrl")
				if srcURL != "" {
					l.Event = _eventResource
					l.SubEvent = srcURL
					l.BusinessCode = "-999"
				}
				if loglevel == _levelInfo {
					l.Result = "1"
				}
				if _, _, kv, err = l.LogData(); err != nil {
					return
				}
				kv = append(kv, extKV(m)...)
			case _logTypeTwo:
				var msg []byte
				l = &monitor.Log{
					LogType:      _logTypeTwo,
					Type:         _typeWeb,
					Product:      source,
					Event:        _eventPage,
					IP:           ip,
					Buvid:        buvid,
					UserAgent:    userAgent,
					Mid:          strconv.FormatInt(mid, 10),
					SubEvent:     url,
					Query:        query,
					HTTPCode:     "200",
					Details:      details(m),
					Result:       "1",
					BusinessCode: "0",
					SubProduct:   stringValueByKey(m, source, "sub_product"),
				}
				if msg, err = json.Marshal(m); err != nil {
					log.Error("s.frontendLog.Marshal error(%v), data(%v)", err, m)
					continue
				}
				l.Message = string(msg)
				if _, _, kv, err = l.LogData(); err != nil {
					return
				}
			default:
				log.Error("s.frontendLog error logype(%v)", logtype)
				return
			}
			switch loglevel {
			case _levelInfo:
				s.mh.Info(c, source, kv...)
			case _levelWarning:
				s.mh.Warn(c, source, kv...)
			case _levelError:
				s.collectFE(c, l)
				s.mh.Error(c, source, kv...)
			default:
				s.mh.Info(c, source, kv...)
			}
		default:
			log.Error("s.frontendLog log data is not json type: " + data)
		}
	}
	return
}

func (s *Service) collectFE(c context.Context, l *monitor.Log) {
	if s.c.CollectFE {
		regex := `.*i\d{1}.hdslb.com`
		if ok, err := regexp.Match(regex, []byte(l.SubEvent)); err != nil {
			log.Error("s.Report.regexp error(%+v), data(%s)", err, l.SubEvent)
			return
		} else if ok {
			l.SubEvent = l.SubEvent[:strings.Index(l.SubEvent, "hdslb.com")+9]
		}
		go s.Collect(context.TODO(), l)
	}
}

func stringValueByKey(m map[string]interface{}, defValue string, keys ...string) (r string) {
	r = defValue
	if m == nil {
		return
	}
	var (
		res interface{}
		ok  bool
	)
	for _, key := range keys {
		if res, ok = m[key]; ok {
			break
		}
		return
	}
	if res == nil {
		return
	}
	switch d := res.(type) {
	case int:
		r = strconv.Itoa(d)
	case float64:
		r = strconv.FormatFloat(d, 'f', -1, 64)
	case int64:
		r = strconv.FormatInt(d, 10)
	case string:
		r = d
	default:
		log.Warn("s.stringValueByKey unexcept type of value(%v)", d)
	}
	return
}

func extKV(m map[string]interface{}) (kv []log.D) {
	if m == nil {
		return
	}
	var keys = []string{"level", "logtype", "url", "cost", "traceid_svr", "traceid_end", "status", "code", "errno", "sub_product"}
	for _, key := range keys {
		delete(m, key)
	}
	for k, v := range m {
		kv = append(kv, log.KV(k, v))
	}
	return
}

func details(m map[string]interface{}) (res map[string]int64) {
	res = make(map[string]int64)
	if m == nil {
		return
	}
	delete(m, "level")
	delete(m, "logtype")
	delete(m, "url")
	for k, v := range m {
		var d int64
		if v == nil {
			res[k] = 0
			continue
		}
		switch v := v.(type) {
		case int64:
			d = v
		case float64:
			d = int64(v)
		case int:
			d = int64(v)
		case int32:
			d = int64(v)
		case string:
			d, _ = strconv.ParseInt(v, 10, 64)
		default:
			log.Warn("s.details unexcept type of value(%v)", v)
		}
		res[k] = d
	}
	return
}

func bStrToInt(str string, isInt bool) (r int, err error) {
	if isInt {
		if r, err = strconv.Atoi(str); err != nil {
			log.Warn("s.bStrToInt error(%+v), string(%s)", err, str)
		}
		return
	}
	for _, i := range str {
		if i == '1' {
			r = r*2 + 1
		} else {
			r *= 2
		}
	}
	return
}

func (s *Service) checkProduct(l *monitor.Log) bool {
	if s.naProducts == nil {
		return true
	}
	if s.naProducts[l.Product] {
		return true
	}
	return false
}
