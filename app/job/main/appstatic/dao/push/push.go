package push

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go-common/app/job/main/appstatic/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_diffFinish = "SELECT COUNT(1) FROM resource_file WHERE resource_id = ? AND url = ? AND file_type = ? AND is_deleted = 0"
	_pushMsg    = "SELECT resource.id, resource.pool_id, resource_pool.`name` FROM resource LEFT JOIN resource_pool ON resource.pool_id = resource_pool.id WHERE resource.id = ? LIMIT 1"
	_platform   = "SELECT b.value FROM resource_config a LEFT JOIN resource_limit b ON a.id = b.config_id WHERE a.resource_id = ? AND b.`column` = 'mobi_app' AND a.is_deleted = 0 AND b.is_deleted = 0"
	_diffPkg    = 1
)

// CallPush calls the push server api
func (d *Dao) CallPush(ctx context.Context, platform string, msg string, ip string) (err error) {
	var (
		cfg    = d.c.Cfg.Push
		params = url.Values{}
	)
	params.Set("operation", fmt.Sprintf("%d", cfg.Operation))
	params.Set("platform", platform)
	params.Set("message", msg)
	params.Set("speed", fmt.Sprintf("%d", cfg.QPS))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(ctx, cfg.URL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), cfg.URL+"?"+params.Encode())
	}
	return
}

// DiffFinish checks whether the resource's diff calculation has been finished or not
func (d *Dao) DiffFinish(c context.Context, resID string) (res bool, err error) {
	count := 0
	row := d.db.QueryRow(c, _diffFinish, resID, "", _diffPkg)
	if err = row.Scan(&count); err != nil {
		log.Error("d.DiffFinish err(%v)", err)
		return
	}
	if count == 0 {
		res = true
	}
	return
}

// PushMsg combines the resource pool info to prepare the msg to call PUSH
func (d *Dao) PushMsg(c context.Context, resID string) (res string, err error) {
	var (
		msg  model.PushMsg
		data []byte
	)
	row := d.db.QueryRow(c, _pushMsg, resID)
	if err = row.Scan(&msg.ResID, &msg.ModID, &msg.ModName); err != nil {
		log.Error("d.PushMsg err(%v)", err)
	}
	if data, err = json.Marshal(msg); err != nil {
		log.Error("PushMsg Info ResID %d, Json Err %v", resID, err)
		return
	}
	res = string(data)
	return
}

// Platform picks the mobi_app value to distinguish the platform to push
func (d *Dao) Platform(c context.Context, resID string) (res []string, err error) {
	rows, err := d.db.Query(c, _platform, resID)
	if err != nil {
		log.Error("db.Query(%d) error(%v)", resID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mobiApp string
		if err = rows.Scan(&mobiApp); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, mobiApp)
	}
	return
}
