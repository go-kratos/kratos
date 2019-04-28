package dao

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_moralLogID = 12
)

// Moral get user moral from cache,if miss get from db.
func (d *Dao) Moral(c context.Context, mid int64) (moral *model.Moral, err error) {
	if moral, err = d.moralCache(c, mid); err == nil {
		return
	}
	if err != nil && err != memcache.ErrNotFound {
		log.Error("Failed to get moral from cache, mid: %d, error:%v", mid, err)
		return
	}
	if moral, err = d.MoralDB(c, mid); err != nil {
		return
	}
	if moral == nil {
		moral = &model.Moral{Mid: mid, Moral: model.DefaultMoral}
	}
	d.SetMoralCache(c, mid, moral)
	return
}

// MoralLog is
func (d *Dao) MoralLog(ctx context.Context, mid int64) ([]*model.UserLog, error) {
	ip := metadata.String(ctx, metadata.RemoteIP)
	t := time.Now().Add(-time.Hour * 24 * 7)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appid", "log_user_action")
	params.Set("business", strconv.FormatInt(_moralLogID, 10))
	params.Set("pn", "1")
	params.Set("ps", "1000")
	params.Set("ctime_from", t.Format("2006-01-02 00:00:00"))
	params.Set("sort", "desc")
	params.Set("order", "ctime")

	res := &model.SearchResult{}
	if err := d.client.Get(ctx, d.c.Host.Search+_searchLogURI, ip, params, res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, ecode.Int(res.Code)
	}
	logs := asMoralLog(res)
	return logs, nil
}

// MoralLogByID get distinct log by log id
func (d *Dao) MoralLogByID(ctx context.Context, logID string) (*model.UserLog, error) {
	ip := metadata.String(ctx, metadata.RemoteIP)
	params := url.Values{}
	params.Set("str_0", logID)
	params.Set("appid", "log_user_action")
	params.Set("business", strconv.FormatInt(_moralLogID, 10))
	params.Set("pn", "1")
	params.Set("ps", "1")
	params.Set("sort", "desc")
	params.Set("order", "ctime")

	res := &model.SearchResult{}
	if err := d.client.Get(ctx, d.c.Host.Search+_searchLogURI, ip, params, res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, ecode.Int(res.Code)
	}

	logs := asMoralLog(res)
	if len(logs) == 0 {
		return nil, ecode.NothingFound
	}
	return logs[0], nil
}

// DeleteMoralLog is
func (d *Dao) DeleteMoralLog(ctx context.Context, logID string) error {
	ip := metadata.String(ctx, metadata.RemoteIP)
	return d.deleteLogReport(ctx, _moralLogID, logID, ip)
}

// DeleteLogReport is
func (d *Dao) deleteLogReport(ctx context.Context, business int, logID string, ip string) error {
	if logID == "" {
		return errors.New("Failed to delete log with empty logID")
	}

	params := url.Values{}
	params.Set("str_0", logID)
	params.Set("appid", "log_user_action")
	params.Set("business", strconv.FormatInt(int64(business), 10))

	res := &model.SearchResult{}
	if err := d.client.Post(ctx, d.c.Host.Search+_deleteLogURI, ip, params, res); err != nil {
		return err
	}
	if res.Code != 0 {
		return ecode.Int(res.Code)
	}
	return nil
}

func asMoralLog(res *model.SearchResult) []*model.UserLog {
	logs := make([]*model.UserLog, 0, len(res.Data.Result))
	for _, r := range res.Data.Result {
		ts, err := time.ParseInLocation("2006-01-02 15:04:05", r.Ctime, time.Local)
		if err != nil {
			log.Warn("Failed to parse log ctime: ctime: %s: %+v", r.Ctime, err)
			continue
		}
		content := map[string]string{
			"from_moral": "",
			"log_id":     "",
			"mid":        "",
			"operater":   "",
			"origin":     "",
			"reason":     "",
			"remark":     "",
			"status":     "",
			"to_moral":   "",
		}
		if err := json.Unmarshal([]byte(r.ExtraData), &content); err != nil {
			log.Warn("Failed to parse extra data in moral log: mid: %d, extra_data: %s: %+v", r.Mid, r.ExtraData, err)
			continue
		}
		if content["mid"] == "" {
			content["mid"] = strconv.FormatInt(r.Mid, 10)
		}
		l := &model.UserLog{
			Mid:     r.Mid,
			IP:      r.IP,
			TS:      ts.Unix(),
			LogID:   content["log_id"],
			Content: content,
		}
		logs = append(logs, l)
	}
	return logs
}
