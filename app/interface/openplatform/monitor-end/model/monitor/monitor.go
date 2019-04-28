package monitor

import (
	"bytes"
	"errors"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"go-common/library/log"
)

const (
	_regex     = `.*\d{2,}`
	_fieldNums = 33
)

// ErrIncompleteLog .
var ErrIncompleteLog = errors.New("Incomplete log")

// Log .
type Log struct {
	// from app log
	RequestURI   string           `form:"request_uri" json:"request_uri"`
	TimeIso      string           `form:"time_iso" json:"time_iso"`
	IP           string           `form:"ip" json:"ip"`
	Version      string           `form:"version" json:"version"`
	Buvid        string           `form:"buvid" json:"buvid"`
	Fts          string           `form:"fts" json:"fts"`
	Proid        string           `form:"proid" json:"proid"`
	Chid         string           `form:"chid" json:"chid"`
	Pid          string           `form:"pid" json:"pid"`
	Brand        string           `form:"brand" json:"brand"`
	Deviceid     string           `form:"deviceid" json:"deviceid"`
	Model        string           `form:"model" json:"model"`
	Osver        string           `form:"osver" json:"osver"`
	Ctime        string           `form:"ctime" json:"ctime"`
	Mid          string           `form:"mid" json:"mid"`
	Ver          string           `form:"ver" json:"ver"`
	Net          string           `form:"net" json:"net"`
	Oid          string           `form:"oid" json:"oid"`
	Product      string           `form:"product" json:"product"`
	Createtime   string           `form:"createtime" json:"createtime"`
	Event        string           `form:"event" json:"event"`
	SubEvent     string           `form:"sub_event" json:"sub_event"`
	LogType      string           `form:"log_type" json:"log_type"`
	Duration     string           `form:"duration" json:"duration"`
	Message      string           `form:"message" json:"message"`
	Result       string           `form:"result" json:"result"`
	ExtJSON      string           `form:"ext_json" json:"ext_json"`
	Traceid      string           `form:"traceid" json:"traceid"`
	Desc         string           `form:"desc" json:"desc"`
	Network      string           `form:"network" json:"network"`
	TraceidEnd   string           `form:"traceid_end" json:"traceid_end"`
	HTTPCode     string           `form:"http_code" json:"http_code"`
	SubProduct   string           `form:"sub_product" json:"sub_product"`
	Codes        string           `json:"codes"`
	BusinessCode string           `json:"business_code"`
	InnerCode    string           `json:"inner_code"`
	TraceidSvr   string           `json:"traceid_svr"`
	Type         string           `json:"type"`
	Query        string           `json:"query"`
	UserAgent    string           `json:"user_agent"`
	Details      map[string]int64 `json:"details"`
}

// Codes .
type Codes struct {
	HTTPCode         interface{} `json:"http_code"`
	HTTPInnerCode    interface{} `json:"http_inner_code"`
	HTTPBusinessCode interface{} `json:"http_business_code"`
}

// LogData .
func (a *Log) LogData() (appID string, logType string, kv []log.D, err error) {
	// if a.Product == "" {
	// 	a.Product = "nameless"
	// }
	//TODO 去除非法appid
	var (
		ok    bool
		regex = "^[A-Za-z_]{1,}$"
	)
	if ok, err = regexp.Match(regex, []byte(a.Product)); err != nil || !ok {
		a.Product = "nameless"
	}
	appID = a.Product
	logType = a.LogType
	kv = a.kv()
	return
}

func (a *Log) kv() (kv []log.D) {
	t := reflect.TypeOf(*a)
	v := reflect.ValueOf(*a)
	for i := 0; i < t.NumField()-1; i++ {
		if v.Field(i).String() == "" {
			continue
		}
		tag := t.Field(i).Tag.Get("json")
		// s, _ := url.QueryUnescape(v.Field(i).Interface().(string))
		kv = append(kv, log.KV(tag, v.Field(i).Interface()))
	}
	if a.Details != nil {
		kv = append(kv, log.KV("detail", a.Details))
	}
	return
}

// LogFromBytes .
func LogFromBytes(msg []byte) *Log {
	var (
		e error
		d string
	)
	a := &Log{}
	data := bytes.Split(msg, []byte("|"))
	v := reflect.ValueOf(a).Elem()
	l := _fieldNums
	if len(data) < l {
		l = len(data)
	}
	for i := 0; i < l; i++ {
		if d, e = url.QueryUnescape(string(data[i])); e != nil {
			tmp := string(data[i])
			tmp = tmp[:strings.LastIndex(tmp, "%")]
			if d, e = url.QueryUnescape(tmp); e != nil {
				d = string(data[i])
				log.Warn("s.LogFromBytes url decode error(%+v), data(%s)", e, string(data[i]))
			}
		}
		if len(d) > 512 {
			d = d[:512]
		}
		v.Field(i).SetString(d)
	}
	return a
}
