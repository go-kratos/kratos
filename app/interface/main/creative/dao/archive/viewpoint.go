package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/log"
)

const (
	_viewPointKey    = "viewpoint_%d_%d"
	_viewPointExp    = 300
	_videoViewPoints = "SELECT id,aid,cid,content,state,ctime,mtime FROM video_viewpoint WHERE aid=? AND cid=? AND state=1 ORDER BY mtime DESC LIMIT ?"
)

// viewPointCacheKey 高能看点MC缓存key
func viewPointCacheKey(aid, cid int64) string {
	return fmt.Sprintf(_viewPointKey, aid, cid)
}

// RawViewPoint get video highlight viewpoint
func (d *Dao) RawViewPoint(c context.Context, aid, cid int64) (vp *archive.ViewPointRow, err error) {
	vps, err := d.RawViewPoints(c, aid, cid, 1)
	if err != nil {
		return
	}
	if len(vps) == 0 {
		return
	}
	vp = vps[0]
	return
}

// RawViewPoints 获取多个版本的高能看点
func (d *Dao) RawViewPoints(c context.Context, aid, cid int64, count int) (vps []*archive.ViewPointRow, err error) {
	rows, err := d.db.Query(c, _videoViewPoints, aid, cid, count)
	if err != nil {
		log.Error("d.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			p struct {
				ID      int64
				AID     int64
				CID     int64
				Content string
				State   int32
				CTime   string
				MTime   string
			}
			points struct {
				Points []*archive.ViewPoint `json:"points"`
			}
		)
		if err = rows.Scan(&p.ID, &p.AID, &p.CID, &p.Content, &p.State, &p.CTime, &p.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if err = json.Unmarshal([]byte(p.Content), &points); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", p.Content, err)
			return
		}
		for i := 0; i < len(points.Points); i++ {
			if points.Points[i].State != 1 {
				points.Points = append(points.Points[:i], points.Points[i+1:]...)
				i--
			}
		}
		vps = append(vps, &archive.ViewPointRow{
			ID:     p.ID,
			AID:    p.AID,
			CID:    p.CID,
			Points: points.Points,
			State:  p.State,
			CTime:  p.CTime,
			MTime:  p.MTime,
		})
	}

	return
}
