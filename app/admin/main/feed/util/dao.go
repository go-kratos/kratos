package util

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/feed/model/common"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus/report"
)

//AddLog add action log
func AddLog(id int, uname string, uid int64, oid int64, action string, obj interface{}) (err error) {
	report.Manager(&report.ManagerInfo{
		Uname:    uname,
		UID:      uid,
		Business: id,
		Type:     0,
		Oid:      oid,
		Action:   action,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{},
		Content: map[string]interface{}{
			"json": obj,
		},
	})
	return
}

//AddLogs add action logs
func AddLogs(logtype int, uname string, uid int64, oid int64, action string, obj interface{}) (err error) {
	report.Manager(&report.ManagerInfo{
		Uname:    uname,
		UID:      uid,
		Business: common.BusinessID,
		Type:     logtype,
		Oid:      oid,
		Action:   action,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{},
		Content: map[string]interface{}{
			"json": obj,
		},
	})
	return
}

//UserInfo get login userinfo
func UserInfo(c *bm.Context) (uid int64, username string) {
	if nameInter, ok := c.Get("username"); ok {
		username = nameInter.(string)
	}
	if uidInter, ok := c.Get("uid"); ok {
		uid = uidInter.(int64)
	}
	if username == "" {
		cookie, _ := c.Request.Cookie("username")
		if cookie == nil || cookie.Value == "" {
			return
		}
		username = cookie.Value
		cookie, _ = c.Request.Cookie("uid")
		if cookie == nil || cookie.Value == "" {
			return
		}
		uidInt, _ := strconv.Atoi(cookie.Value)
		uid = int64(uidInt)
	}
	return
}

//TrimStrSpace trim string space
func TrimStrSpace(v string) string {
	return strings.TrimSpace(v)
}

//CTimeStr current time string
func CTimeStr() (cTime string) {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}
