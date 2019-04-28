package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"time"
)

// CheckSecretkey 鉴权
func TXCheckSecretkey(c *bm.Context) {
	req := c.Request
	q := req.URL.Query()
	wsSecret := q.Get("wsSecret")
	wsTime := q.Get("wsTime")
	url := fmt.Sprintf("http://%s%s", c.Request.Host, req.URL.Path)
	log.Warn("request url = %s", url)

	key := "xtCTceP0fdH8"

	err := check(wsSecret, wsTime, url, key)

	if err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.AccessTokenExpires)
		c.Abort()
		return
	}
}

// BvcCheckSecret 鉴权
func BvcCheckSecret(c *bm.Context) {
	req := c.Request
	q := req.URL.Query()
	wsSecret := q.Get("wsSecret")
	wsTime := q.Get("wsTime")
	url := fmt.Sprintf("http://%s%s", c.Request.Host, req.URL.Path)
	log.Warn("request url = %s", url)

	key := "cMRvgcQXZdph"

	err := check(wsSecret, wsTime, url, key)

	if err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.AccessTokenExpires)
		c.Abort()
		return
	}
}

func check(wsSecret string, wsTime string, url string, key string) error {
	if wsSecret == "" || wsTime == "" {
		return errors.New("secret or time is empty")
	}

	wsTimeInt, err := strconv.ParseInt(wsTime, 10, 64)
	//log.Warn("%d, %d", time.Now().Unix(), wsTimeInt)
	if err != nil || time.Now().Unix() > wsTimeInt {
		return errors.New("request time expired")
	}

	notCheckStr := fmt.Sprintf("%s%s%s", key, url, wsTime)
	log.Warn("%s", notCheckStr)
	h := md5.New()
	h.Write([]byte(notCheckStr))
	cipherStr := h.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	if wsSecret != sign {
		return errors.New("auth failed")
	}
	return nil
}
