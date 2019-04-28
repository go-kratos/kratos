package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/model/signmodel"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	// URLUNames url for usernames to uid
	URLUNames = "/x/admin/manager/users/unames"
	// URLUids url for uids to username
	URLUids = "/x/admin/manager/users/uids"
	// URLAuditLog url for sign up change log
	URLAuditLog = "/x/admin/search/log"
)

// Dao is redis dao.
type Dao struct {
	c          *conf.Config
	managerDB  *sql.DB
	HTTPClient *bm.Client
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		managerDB: sql.NewMySQL(c.DB.Manager),
		// http client
		HTTPClient: bm.NewClient(c.HTTPClient.Normal),
	}
	return d
}

// Close fn
func (d *Dao) Close() {
	if d.managerDB != nil {
		d.managerDB.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.managerDB.Ping(c)
}

// GetUNamesByUids get name by uid
func (d *Dao) GetUNamesByUids(c context.Context, uids []int64) (res map[int64]string, err error) {
	var param = url.Values{}
	var uidStr = xstr.JoinInts(uids)
	param.Set("uids", uidStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[int64]string `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+URLUNames, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+URLUNames+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+URLUNames+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

// GetUIDByNames get uid by name
func (d *Dao) GetUIDByNames(c context.Context, names []string) (res map[string]int64, err error) {
	var param = url.Values{}
	var namesStr = strings.Join(names, ",")
	param.Set("unames", namesStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[string]int64 `json:"data"`
		Message string           `json:"message"`
	}

	err = d.HTTPClient.Get(c, d.c.Host.Manager+URLUids, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+URLUids+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+URLUids+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

// SignUpAuditLogs get sign up audit log .
func (d *Dao) SignUpAuditLogs(c context.Context, arg *signmodel.SignOpSearchArg) (res *signmodel.SignAuditListReply, err error) {
	params := url.Values{}
	params.Set("appid", "log_audit")
	params.Set("business", fmt.Sprintf("%d", signmodel.SignUpLogBizID))
	params.Set("order", arg.Order)
	params.Set("sort", arg.Sort)
	params.Set("int_0", fmt.Sprintf("%d", arg.SignID))
	params.Set("oid", fmt.Sprintf("%d", arg.Mid))
	params.Set("type", fmt.Sprintf("%d", arg.Tp))
	params.Set("pn", fmt.Sprintf("%d", arg.PN))
	params.Set("ps", fmt.Sprintf("%d", arg.PS))
	var httpRes struct {
		Code    int                           `json:"code"`
		Data    *signmodel.BaseAuditListReply `json:"data"`
		Message string                        `json:"message"`
	}
	err = d.HTTPClient.Get(c, d.c.Host.Manager+URLAuditLog, "", params, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.c.Host.Manager+URLAuditLog+"?"+params.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", d.c.Host.Manager+URLAuditLog+"?"+params.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = new(signmodel.SignAuditListReply)
	res.Page = arg.PN
	res.Size = arg.PS
	res.Order = arg.Order
	res.Sort = arg.Sort
	if httpRes.Data != nil && httpRes.Data.Pager != nil {
		res.TotalCount = httpRes.Data.Pager.TotalCount
	}
	for _, v := range httpRes.Data.Result {
		ctime, _ := time.ParseInLocation(upcrmmodel.TimeFmtDateTime, v.CTime, time.Local)
		var signAuit = &signmodel.SignAuditReply{
			Mid:      v.OID,
			SignID:   v.IntOne,
			Tp:       v.Tp,
			OperID:   v.UID,
			OperName: v.UName,
			CTime:    xtime.Time(ctime.Unix()),
		}
		var content signmodel.SignContentReply
		json.Unmarshal([]byte(v.ExtraData), &content)
		signAuit.Content = &content
		if signAuit.Content.New != nil {
			buildSignContractURL(signAuit.Content.New)
		}
		if signAuit.Content.Old != nil {
			buildSignContractURL(signAuit.Content.Old)
		}
		res.Result = append(res.Result, signAuit)
	}
	return
}

func buildSignContractURL(su *signmodel.SignUpArg) {
	if su.ContractInfo == nil {
		return
	}
	for _, v := range su.ContractInfo {
		v.Filelink = signmodel.BuildOrcBfsURL(v.Filelink)
		v.Filelink = signmodel.BuildDownloadURL(v.Filename, v.Filelink)
	}
}
