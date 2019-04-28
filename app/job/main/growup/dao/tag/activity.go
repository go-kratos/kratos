package tag

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/ecode"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_activityInfoSQL = "SELECT archive_id FROM activity_info"
	// insert
	_insertActivityInfoSQL = "INSERT INTO activity_info(archive_id,mid,activity_id,category,tag_id) VALUES %s"
)

// InsertActivityInfo insert
func (d *Dao) InsertActivityInfo(c context.Context, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_insertActivityInfoSQL, vals))
	if err != nil {
		log.Error("dao.InsertActivityInfo exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// ListActivityInfo list
func (d *Dao) ListActivityInfo(c context.Context) (avs map[int64]bool, err error) {
	avs = make(map[int64]bool)
	rows, err := d.db.Query(c, _activityInfoSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avID int64
		err = rows.Scan(&avID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		avs[avID] = true
	}
	return
}

// GetVideoActivityInfo get activity by api.
func (d *Dao) GetVideoActivityInfo(c context.Context, activities []int64, pn, ps int) (info []*model.ActivityInfo, err error) {
	info = make([]*model.ActivityInfo, 0)
	params := url.Values{}
	params.Set("search_type", "manager_arc")
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("page", strconv.Itoa(pn))
	params.Set("mission_id", xstr.JoinInts(activities))

	var res struct {
		Code   int                   `json:"code"`
		Result []*model.ActivityInfo `json:"result"`
		Total  int                   `json:"total"`
	}

	if err = d.client.Get(c, d.archiveURL, "", params, &res); err != nil {
		log.Error("d.client.Get url(%s) error(%v)", d.archiveURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetVideoActivityInfo code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, d.archiveURL+"?"+params.Encode(), err)
		err = ecode.GrowupGetActivityError
		return
	}
	info = res.Result
	return
}

// GetCmActivityInfo get column activity
func (d *Dao) GetCmActivityInfo(c context.Context, activityID int64, pn, ps int) (info []*model.ActivityInfo, err error) {
	info = make([]*model.ActivityInfo, 0)
	params := url.Values{}
	params.Set("order", "ctime")
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("page", strconv.Itoa(pn))

	var res struct {
		Code int              `json:"code"`
		Data *model.ColumnAct `json:"data"`
	}
	url := fmt.Sprintf("%s%d", d.columnActURL, activityID)
	if err = d.client.Get(c, url, "", params, &res); err != nil {
		log.Error("d.client.Get url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetCmActivityInfo code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), err)
		err = ecode.GrowupGetActivityError
		return
	}

	if res.Data.List == nil {
		return
	}
	for _, cm := range res.Data.List {
		info = append(info, &model.ActivityInfo{
			ActivityID: cm.SID,
			AvID:       cm.ID,
			MID:        cm.MID,
			TypeID:     cm.Category.ID,
			CDate:      cm.CTime,
		})
	}
	return
}

// GetVideoTypes get all types.
func (d *Dao) GetVideoTypes(c context.Context) (rmap map[int64]int64, err error) {
	var res struct {
		Code    int                        `json:"code"`
		Message string                     `json:"message"`
		Data    map[int16]*model.TypesInfo `json:"data"`
	}
	if err = d.client.Get(c, d.typeURL, "", nil, &res); err != nil {
		log.Error("dao.GetAllTypes GET error(%v) | typesURI(%s)", err, d.typeURL)
		return
	}
	if res.Code != 0 {
		log.Error("dao.GetAllTypes code != 0. res.Code(%d) | typesURI(%s) res(%v)", res.Code, d.typeURL, res)
		err = ecode.GrowupGetTypeError
		return
	}
	rmap = make(map[int64]int64, len(res.Data))
	for _, v := range res.Data {
		if v.PID != 0 {
			rmap[v.ID] = v.PID
		}
	}
	return
}

// GetColumnTypes get column types.
func (d *Dao) GetColumnTypes(c context.Context) (rmap map[int64]int64, err error) {
	var res struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []*model.ColumnType `json:"data"`
	}
	if err = d.client.Get(c, d.columnURL, "", nil, &res); err != nil {
		log.Error("dao.GetColumnTypes GET error(%v) | typesURI(%s)", err, d.columnURL)
		return
	}
	if res.Code != 0 {
		log.Error("dao.GetColumnTypes code != 0. res.Code(%d) | typesURI(%s) res(%v)", res.Code, d.columnURL, res)
		err = ecode.GrowupGetTypeError
		return
	}
	rmap = make(map[int64]int64)
	for _, parent := range res.Data {
		if parent.Children != nil {
			for _, child := range parent.Children {
				rmap[child.ID] = child.ParentID
			}
		}
	}
	return
}
