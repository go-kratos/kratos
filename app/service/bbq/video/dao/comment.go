package dao

import (
	"bytes"
	"context"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	//DefaultCmType 默认评论类型
	DefaultCmType = 23
)

// ReplyCounts 批量评论数
func (d *Dao) ReplyCounts(c context.Context, ids []int64, t int64) (res map[int64]*model.ReplyCount, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	oidStr := strings.Replace(strings.Trim(fmt.Sprint(ids), "[]"), " ", ",", -1)
	req := map[string]interface{}{
		"type": t,
		"oid":  oidStr,
	}
	res = make(map[int64]*model.ReplyCount)
	var r []byte
	r, err = replyHTTPCommon(c, d.httpClient, d.c.URLs["reply_counts"], "GET", req, ip)
	if err != nil {
		log.Infov(c,
			log.KV("log", fmt.Sprintf("replyHTTPCommon err [%v]", err)),
		)
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("json unmarlshal err data[%s]", string(r))))
	}
	return
}
