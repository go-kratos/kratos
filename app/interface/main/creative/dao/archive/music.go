package archive

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"go-common/app/interface/main/creative/model/music"
	"go-common/library/log"
)

const (
	_AllMusicsSQL = "SELECT cooperate,id,sid,name,frontname,musicians,mid,cover,playurl,state,duration,filesize,ctime,pubtime,tags,timeline FROM music"
)

// AllMusics fn
func (d *Dao) AllMusics(c context.Context) (res map[int64]*music.Music, err error) {
	rows, err := d.db.Query(c, _AllMusicsSQL)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*music.Music)
	for rows.Next() {
		v := &music.Music{}
		var fName string
		if err = rows.Scan(&v.Cooperate, &v.ID, &v.SID, &v.Name, &fName, &v.Musicians, &v.UpMID, &v.Cover, &v.Playurl, &v.State, &v.Duration, &v.FileSize, &v.CTime, &v.Pubtime, &v.TagsStr, &v.Timeline); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if len(fName) > 0 {
			v.Name = fName
		}
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
			v.Tags = strings.Split(v.TagsStr, ",")
		}
		res[v.SID] = v
	}
	return
}
