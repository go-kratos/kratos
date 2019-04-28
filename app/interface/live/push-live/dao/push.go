package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"go-common/library/xstr"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type _response struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

// BatchPush 批量推送，失败重试
func (d *Dao) BatchPush(fans *[]int64, task *model.ApPushTask) (total int) {
	limit := d.c.Push.PushOnceLimit
	retry := d.c.Push.PushRetryTimes
	var times int
	for {
		var (
			mids []int64
			err  error
		)
		uuid := d.GetUUID(task, times)
		l := len(*fans)
		if l == 0 {
			break
		} else if l <= limit {
			mids = (*fans)[:l]
		} else {
			mids = (*fans)[:limit]
			l = limit
		}
		*fans = (*fans)[l:]

		for i := 0; i < retry; i++ {
			// 每次投递成功结束循环
			if err = d.Push(mids, task, uuid); err == nil {
				total += len(mids) //单次投递成功数
				break
			}
			time.Sleep(time.Duration(time.Second * 3))
		}
		times++
		if err != nil {
			// 重试若干次仍然失败，需要记录日志并且配置elk告警
			log.Error("[dao.push|BatchPush] retry push failed. error(%+v), retry times(%d), task(%+v)", err, retry, task)
		}
	}
	return
}

// Push 调用推送接口
func (d *Dao) Push(fans []int64, task *model.ApPushTask, uuid string) (err error) {
	if len(fans) == 0 {
		log.Info("[dao.push|Push] empty fans. task(%+v)", task)
		return
	}

	// 业务参数
	businessID, token := d.getPushBusiness(task.Group)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("app_id", strconv.Itoa(d.c.Push.AppID))
	w.WriteField("business_id", strconv.Itoa(businessID))
	w.WriteField("alert_title", d.GetPushTemplate(task.Group, task.AlertTitle))
	w.WriteField("alert_body", task.AlertBody)
	w.WriteField("mids", xstr.JoinInts(fans))
	w.WriteField("link_type", strconv.Itoa(task.LinkType))
	w.WriteField("link_value", task.LinkValue)
	w.WriteField("expire_time", strconv.Itoa(task.ExpireTime))
	w.WriteField("group", task.Group)
	w.WriteField("uuid", uuid)
	w.Close()
	// 签名
	query := map[string]string{
		"ts":     strconv.FormatInt(time.Now().Unix(), 10),
		"appkey": d.c.HTTPClient.Key,
	}
	query["sign"] = d.getSign(query, d.c.HTTPClient.Secret)
	requestURL := fmt.Sprintf("%s?ts=%s&appkey=%s&sign=%s", d.c.Push.MultiAPI, query["ts"], query["appkey"], query["sign"])
	// request
	req, err := http.NewRequest(http.MethodPost, requestURL, buf)
	if err != nil {
		log.Error("[dao.push|Push] http.NewRequest error(%+v), url(%s), uuid(%s), task(%+v)",
			err, requestURL, uuid, task)
		PromError("[dao.push|Push] http:NewRequest")
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "token="+token)
	res := &_response{}
	if err = d.httpClient.Do(context.TODO(), req, res); err != nil {
		log.Error("[dao|push|Push] httpClient.Do error(%+v), url(%s), uuid(%s), task(%+v)",
			err, requestURL, uuid, task)
		PromError("[dao.push|Push] http:Do")
		return
	}
	// response
	if res.Code != 0 || res.Data == 0 {
		log.Error("[dao.push|Push] push failed. url(%s), uuid(%s), response(%+v), task(%+v)", requestURL, uuid, res, task)
		err = fmt.Errorf("[dao.push|Push] push failed. url(%s), uuid(%s), response(%+v), task(%+v)", requestURL, uuid, res, task)
	} else {
		log.Info("[dao.push|Push] push success. url(%s), uuid(%s), response(%+v), task(%+v)", requestURL, uuid, res, task)
	}
	return
}

// GetPushTemplate 根据类型返回不同的推送文案
func (d *Dao) GetPushTemplate(group string, part string) (template string) {
	switch group {
	case model.SpecialGroup:
		template = fmt.Sprintf(d.c.Push.SpecialCopyWriting, part)
	case model.AttentionGroup:
		template = fmt.Sprintf(d.c.Push.DefaultCopyWriting, part)
	default:
		template = part
	}
	return
}

// GetUUID 构造一个每次请求的uuid
func (d *Dao) GetUUID(task *model.ApPushTask, times int) string {
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(times))
	b.WriteString(task.Group) // Group必须加入uuid计算，区分单次开播提醒关注与特别关注
	b.WriteString(strconv.Itoa(d.c.Push.BusinessID))
	b.WriteString(strconv.FormatInt(task.TargetID, 10))
	b.WriteString(strconv.Itoa(task.ExpireTime))
	b.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	mh := md5.Sum(b.Bytes())
	uuid := hex.EncodeToString(mh[:])
	return uuid
}

// getSign 获取签名
func (d *Dao) getSign(params map[string]string, secret string) (sign string) {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k) + "=")
		buf.WriteString(url.QueryEscape(params[k]))
	}
	hash := md5.New()
	io.WriteString(hash, buf.String()+secret)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// getPushBusiness 获取推送配置
func (d *Dao) getPushBusiness(group string) (businessID int, token string) {
	// 预约走单独的白名单通道，business id 和token不一样
	if group == "activity_appointment" {
		businessID = 41
		token = "13aoowdzm0u8pcqdoulvj5vdkihohtcj"
	} else {
		businessID = d.c.Push.BusinessID
		token = d.c.Push.BusinessToken
	}
	return
}
