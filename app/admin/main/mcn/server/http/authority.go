package http

import (
	"net/http"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// check username and dashboard sessionid
func checkCookie(c *bm.Context) (username, sid string, err error) {
	var r = c.Request
	var name *http.Cookie
	if name, err = r.Cookie("username"); err == nil {
		username = name.Value
	}
	var session *http.Cookie
	if session, err = r.Cookie("_AJSESSIONID"); err == nil {
		sid = session.Value
	}
	if username == "" || sid == "" {
		err = ecode.Unauthorized
	}
	return
}
