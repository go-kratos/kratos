package report

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/library/conf/env"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

type conf struct {
	Secret string
	Addr   string
}

var (
	mn   *databus.Databus
	user *databus.Databus

	// ErrInit report init error
	ErrInit = errors.New("report initialization failed")
)

const (
	_timeFormat = "2006-01-02 15:04:05"

	_uname    = "uname"
	_uid      = "uid"
	_business = "business"
	_type     = "type"
	_oid      = "oid"
	_action   = "action"
	_ctime    = "ctime"
	_platform = "platform"
	_build    = "build"
	_buvid    = "buvid"
	_ip       = "ip"
	_mid      = "mid"

	_indexInt = "int_"
	_indexStr = "str_"
	_extra    = "extra_data"
)

// ManagerInfo manager report info.
type ManagerInfo struct {
	// common
	Uname    string
	UID      int64
	Business int
	Type     int
	Oid      int64
	Action   string
	Ctime    time.Time
	// extra
	Index   []interface{}
	Content map[string]interface{}
}

// UserInfo user report info
type UserInfo struct {
	Mid      int64
	Platform string
	Build    int64
	Buvid    string
	Business int
	Type     int
	Oid      int64
	Action   string
	Ctime    time.Time
	IP       string
	// extra
	Index   []interface{}
	Content map[string]interface{}
}

// UserActionLog 用户行为日志
type UserActionLog struct {
	Uname    string `json:"uname"`
	UID      int64  `json:"uid"`
	Business int    `json:"business"`
	Type     int    `json:"type"`
	Oid      int64  `json:"oid"`
	Action   string `json:"action"`
	Platform string `json:"platform"`
	Build    int64  `json:"build"`
	Buvid    string `json:"buvid"`
	IP       string `json:"ip"`
	Mid      int64  `json:"mid"`
	Int0     int64  `json:"int_0"`
	Int1     int64  `json:"int_1"`
	Int2     int64  `json:"int_2"`
	Str0     string `json:"str_0"`
	Str1     string `json:"str_1"`
	Str2     string `json:"str_2"`
	Ctime    string `json:"ctime"`
	Extra    string `json:"extra_data"`
}

// AuditLog 审核日志
type AuditLog struct {
	Uname    string `json:"uname"`
	UID      int64  `json:"uid"`
	Business int    `json:"business"`
	Type     int    `json:"type"`
	Oid      int64  `json:"oid"`
	Action   string `json:"action"`
	Int0     int64  `json:"int_0"`
	Int1     int64  `json:"int_1"`
	Int2     int64  `json:"int_2"`
	Str0     string `json:"str_0"`
	Str1     string `json:"str_1"`
	Str2     string `json:"str_2"`
	Ctime    string `json:"ctime"`
	Extra    string `json:"extra_data"`
}

// InitManager init manager report log agent.
func InitManager(c *databus.Config) {
	if c == nil {
		c = _managerConfig
		if d, ok := _defaultManagerConfig[env.DeployEnv]; ok {
			c.Secret = d.Secret
			c.Addr = d.Addr
		}
	}
	mn = databus.New(c)
}

// InitUser init user report log agent.
func InitUser(c *databus.Config) {
	if c == nil {
		c = _userConfig
		if d, ok := _defaultUserConfig[env.DeployEnv]; ok {
			c.Secret = d.Secret
			c.Addr = d.Addr
		}
	}
	user = databus.New(c)
}

// Manager log a message for manager, xx-admin.
func Manager(m *ManagerInfo) error {
	if mn == nil || m == nil {
		return ErrInit
	}
	v := map[string]interface{}{}
	if len(m.Content) > 0 {
		extraData, _ := json.Marshal(m.Content)
		v[_extra] = string(extraData)
	}
	v[_business] = m.Business
	v[_type] = m.Type
	v[_uid] = m.UID
	v[_oid] = m.Oid
	v[_uname] = m.Uname
	v[_action] = m.Action
	v[_ctime] = m.Ctime.Format(_timeFormat)
	return report(mn, v, m.Index...)
}

// User log a message for user, xx-interface.
func User(u *UserInfo) error {
	if user == nil || u == nil {
		return ErrInit
	}
	v := map[string]interface{}{}
	if len(u.Content) > 0 {
		extraData, _ := json.Marshal(u.Content)
		v[_extra] = string(extraData)
	}
	v[_business] = u.Business
	v[_type] = u.Type
	v[_mid] = u.Mid
	v[_oid] = u.Oid
	v[_build] = u.Build
	v[_action] = u.Action
	v[_platform] = u.Platform
	v[_buvid] = u.Buvid
	v[_ip] = u.IP
	v[_ctime] = u.Ctime.Format(_timeFormat)
	return report(user, v, u.Index...)
}

func report(h *databus.Databus, v map[string]interface{}, extras ...interface{}) error {
	var i, j int
	for _, extra := range extras {
		switch ex := extra.(type) {
		case string:
			v[_indexStr+strconv.Itoa(i)] = ex
			i++
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			v[_indexInt+strconv.Itoa(j)] = ex
			j++
		}
	}
	return h.Send(context.Background(), v[_ctime].(string), v)
}
