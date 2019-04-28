package music

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/creative/model/music"
	"go-common/library/log"
	"go-common/library/xstr"
	"sort"
	"strings"
)

const (
	_CategorysSQL  = "SELECT id,pid,name,`index`,camera_index FROM music_category WHERE id IN (%s) and state = 0 order by `index` asc "
	_MCategorysSQL = "SELECT id,sid,tid,`index`,ctime FROM music_with_category where state = 0 order by tid asc, `index` asc "
	_MusicsSQL     = "SELECT cooperate,id,sid,name,frontname,musicians,mid,cover,playurl,state,duration,filesize,ctime,mtime,pubtime,tags,timeline FROM music WHERE sid IN (%s) and state = 0 "
)

// Categorys fn
func (d *Dao) Categorys(c context.Context, ids []int64) (res []*music.Category, resMap map[int]*music.Category, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_CategorysSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*music.Category, 0)
	resMap = make(map[int]*music.Category, len(ids))
	for rows.Next() {
		v := &music.Category{}
		if err = rows.Scan(&v.ID, &v.PID, &v.Name, &v.Index, &v.CameraIndex); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		res = append(res, v)
		resMap[v.ID] = v
	}
	return
}

// MCategorys fn
func (d *Dao) MCategorys(c context.Context) (res []*music.Mcategory, err error) {
	rows, err := d.db.Query(c, _MCategorysSQL)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*music.Mcategory, 0)
	for rows.Next() {
		v := &music.Mcategory{}
		if err = rows.Scan(&v.ID, &v.SID, &v.Tid, &v.Index, &v.CTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		res = append(res, v)
	}
	return
}

// Music fn
func (d *Dao) Music(c context.Context, sids []int64) (res map[int64]*music.Music, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_MusicsSQL, xstr.JoinInts(sids)))
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*music.Music, len(sids))
	for rows.Next() {
		v := &music.Music{}
		var fName string
		if err = rows.Scan(&v.Cooperate, &v.ID, &v.SID, &v.Name, &fName, &v.Musicians, &v.UpMID, &v.Cover, &v.Playurl, &v.State, &v.Duration, &v.FileSize, &v.CTime, &v.MTime, &v.Pubtime, &v.TagsStr, &v.Timeline); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if len(fName) > 0 {
			v.Name = fName
		}
		v.CooperateURL = d.c.H5Page.Cooperate
		v.Tl = make([]*music.TimePoint, 0)
		if len(v.Timeline) > 0 {
			if err = json.Unmarshal([]byte(v.Timeline), &v.Tl); err != nil {
				log.Error("json.Unmarshal Timeline failed error(%v)", err)
				continue
			}
			sort.Slice(v.Tl, func(i, j int) bool {
				return v.Tl[i].Point < v.Tl[j].Point
			})
			if len(v.Tl) > 0 {
				for _, point := range v.Tl {
					if point.Recommend == 1 {
						v.RecommendPoint = point.Point
						break
					}
				}
			}
		}
		v.Tags = make([]string, 0)
		if len(v.TagsStr) > 0 {
			v.Tags = append(v.Tags, strings.Split(v.TagsStr, ",")...)
		}
		res[v.SID] = v
	}
	return
}
