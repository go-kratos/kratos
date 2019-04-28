package http

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// audit season with content
func audit(c *bm.Context) {
	if err := auditT(c); err != nil { // if some error, return it
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func auditT(c *bm.Context) (err error) {
	var (
		audit model.Audit
		req   = c.Request
	)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	if err = json.Unmarshal(body, &audit); err != nil {
		log.Error("audit json(%s) error:(%v)", string(body), err)
		err = ecode.RequestErr
		return
	}
	if !validateJSONData(audit) {
		log.Error("audit msg (%s), missing field", string(body))
		err = ecode.RequestErr
		return
	}
	if !checkSign(c, string(body)) {
		log.Error("audit msg (%s), sign error", string(body))
		err = ecode.RequestErr
		return
	}
	return auditSvc.HandleAudits(c, audit.IDList)
}

// validateJSONData check json format whether valid
func validateJSONData(a model.Audit) bool {
	if a.OpType == "" {
		return false
	}
	for _, v := range a.IDList {
		if v.Type == "" || v.VID == "" || v.Action == "" {
			return false
		}
	}
	return a.Count > 0
}

// checkSign check sign whether valid
func checkSign(c *bm.Context, body string) bool {
	var (
		req   = c.Request.Form
		query = make(map[string]string)
		ts    = req.Get("ts")
		key   = req.Get("key")
		sign  = req.Get("sign")
	)
	if key != signCfg.Key {
		log.Error("The appkey not exists")
		return false
	}
	if ts == "" {
		log.Error("The timestamp not exists")
		return false
	}
	query["ts"] = ts
	query["body"] = body
	query["appkey"] = key
	if sign == "" {
		log.Error("The sign not exists")
		return false
	}
	getSign := signature(query)
	if sign != getSign {
		log.Error("The expected signature is :(%s)", getSign)
		return false
	}
	return sign == getSign
}

func signature(query map[string]string) string {
	secret := signCfg.Secret
	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var str string
	for _, v := range keys {
		str += string(v) + "=" + query[v] + "&"
	}
	str = str[:len(str)-1] + secret
	hash := md5.New()
	hash.Write([]byte(str))
	sign := fmt.Sprintf("%x", hash.Sum(nil))
	return sign
}
