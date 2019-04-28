package dao

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

type _response struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

// NoticeFans pushs the notification to fans.
func (d *Dao) NoticeFans(fans *[]int64, params map[string]interface{}) (err error) {
	arc := params["archive"].(*model.Archive)
	group := strings.TrimSpace(params["group"].(string))
	msgTemplate := params["msgTemplate"].(string)
	uuid := params["uuid"].(string)
	relationType := params["relationType"].(int)
	author := "UP主"
	if arc.Author != "" {
		author = fmt.Sprintf(`“%s”`, arc.Author)
	}
	// 普通关注和特殊关注用不同的业务组推
	businessID := d.c.Push.BusinessID
	businessToken := d.c.Push.BusinessToken
	if relationType == model.RelationSpecial {
		businessID = d.c.Push.BusinessSpecialID
		businessToken = d.c.Push.BusinessSpecialToken
	}
	msg := fmt.Sprintf(msgTemplate, author, arc.Title)
	sp := strings.SplitN(msg, "\r\n", 2)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("group", group) // 实验组名，值为实验组数据表名
	w.WriteField("app_id", "1")  // 1表示哔哩哔哩动画
	w.WriteField("business_id", strconv.Itoa(businessID))
	w.WriteField("alert_title", sp[0])
	w.WriteField("alert_body", sp[1])
	w.WriteField("mids", xstr.JoinInts(*fans))
	w.WriteField("link_type", "2") // 2代表视频稿件播放页
	w.WriteField("link_value", strconv.FormatInt(arc.ID, 10))
	w.WriteField("uuid", uuid)
	// 1、v5.20.0 后客户端才接特殊关注  2、 iPad版本没更新，不推
	w.WriteField("builds", `{"2":{"Build":6500,"Condition":"gte"}, "3":{"Build":0,"Condition":"lt"}, "1":{"Build":519010,"Condition":"gte"}}`)
	w.Close()
	query := map[string]string{
		"ts":     strconv.FormatInt(time.Now().Unix(), 10),
		"appkey": d.c.HTTPClient.Key,
	}
	query["sign"] = d.signature(query, d.c.HTTPClient.Secret)
	url := fmt.Sprintf("%s?ts=%s&appkey=%s&sign=%s", d.c.Push.AddAPI, query["ts"], query["appkey"], query["sign"])
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v) uuid(%s)", url, err, uuid)
		PromError("http:NewRequest")
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("token=%s", businessToken))
	res := &_response{}
	if err = d.httpClient.Do(context.TODO(), req, &res); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		PromError("http:Do")
		return
	}
	if res.Code != 0 || res.Data == 0 {
		log.Error("push failed archive(%d) upper(%d) fans_total(%d) group(%s) response(%+v)", arc.ID, arc.Mid, len(*fans), group, res)
	} else {
		log.Info("push success archive(%d) upper(%d) fans_total(%d) group(%s) response(%+v)", arc.ID, arc.Mid, len(*fans), group, res)
	}
	return
}
