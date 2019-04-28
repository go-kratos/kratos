package abtest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/experiment"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selExpLimit = `SELECT experiment_id,build,conditions FROM experiment_limit`
	_selExpByIDs = `SELECT id,name,plat,strategy,description,traffic_group FROM experiment WHERE state=1 AND id IN (%s) ORDER BY uptime DESC`
)

// Dao is notice dao.
type Dao struct {
	db       *sql.DB
	limit    *sql.Stmt
	dataPath string
	client   *bm.Client
}

// New new a notice dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:       sql.NewMySQL(c.MySQL.Show),
		client:   bm.NewClient(conf.Conf.HTTPClient),
		dataPath: c.Host.Data + "/abserver/v1/app/match-exp",
	}
	d.limit = d.db.Prepared(_selExpLimit)
	return
}

func (d *Dao) ExperimentLimit(c context.Context) (lm map[int64][]*experiment.Limit, err error) {
	rows, err := d.limit.Query(c)
	if err != nil {
		log.Error("d.limit.Query error (%v)", err)
		return
	}
	defer rows.Close()
	lm = map[int64][]*experiment.Limit{}
	for rows.Next() {
		limit := &experiment.Limit{}
		if err = rows.Scan(&limit.ExperimentID, &limit.Build, &limit.Condition); err != nil {
			log.Error("rows.Scan err (%v)", err)
			continue
		}
		lm[limit.ExperimentID] = append(lm[limit.ExperimentID], limit)
	}
	return
}

func (d *Dao) ExperimentByIDs(c context.Context, ids []int64) (eps []*experiment.Experiment, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selExpByIDs, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.expByIDs.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ep := &experiment.Experiment{}
		if err = rows.Scan(&ep.ID, &ep.Name, &ep.Plat, &ep.Strategy, &ep.Desc, &ep.TrafficGroup); err != nil {
			log.Error("rows.Scan err (%v)", err)
			continue
		}
		eps = append(eps, ep)
	}
	return
}

// AbServer  http://info.bilibili.co/pages/viewpage.action?pageId=8741843 大数据abtest
func (d *Dao) AbServer(c context.Context, buvid, device, mobiAPP, filteredStr string, build int, mid int64) (res json.RawMessage, err error) {
	params := url.Values{}
	params.Set("buvid", buvid)
	params.Set("device", device)
	params.Set("mobi_app", mobiAPP)
	params.Set("build", strconv.Itoa(build))
	if mid > 0 {
		params.Set("mid", strconv.FormatInt(mid, 10))
	}
	if filteredStr != "" {
		params.Set("filtered", filteredStr)
	}
	var data struct {
		Code int `json:"errorCode"`
	}
	if err = d.client.Get(c, d.dataPath, "", params, &res); err != nil {
		return
	}
	if err = json.Unmarshal(res, &data); err != nil {
		err = errors.Wrap(err, "json.Unmarshal")
		return
	}
	if data.Code != ecode.OK.Code() {
		log.Warn("code(%d) path:(%s)", data.Code, d.dataPath+params.Encode())
		return
	}
	return
}
