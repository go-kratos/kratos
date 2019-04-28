package material

import (
	"context"
	"encoding/json"
	"fmt"
	appMdl "go-common/app/interface/main/creative/model/app"
	"go-common/app/interface/main/creative/model/music"
	"go-common/library/log"
	"go-common/library/xstr"
	"strings"
	"time"
)

const (
	_basicSQL        = "SELECT mtime,id,name,rank,extra,platform,build,type FROM material WHERE type in (%s) and state = 0 order by type asc, rank asc"
	_FiltersSQL      = "SELECT mtime,id,name,rank,extra,platform,build FROM material WHERE type=2 and state = 0 order by rank asc "
	_CategoryBindSQL = "SELECT a.new, a.type, a.name as cname,a.rank as crank,b.index as brank,b.material_id as mid, b.category_id as cid from material_category a join material_with_category b where a.id = b.category_id and a.state=0 and b.state=0 and a.type =%d order by a.rank asc, b.index asc"
	_VstickersSQL    = "SELECT mtime,id,name,rank,extra,platform,build FROM material WHERE type=7 and state = 0 order by rank asc "
	_CooperatesSQL   = "SELECT mtime,id,name,rank,extra,platform,build FROM material WHERE type=9 and state = 0 order by rank asc "
)

var basicTypes = []int64{
	int64(appMdl.TypeSubtitle),
	int64(appMdl.TypeFont),
	int64(appMdl.TypeHotWord),
	int64(appMdl.TypeSticker),
	int64(appMdl.TypeIntro),
	int64(appMdl.TypeTransition),
	int64(appMdl.TypeTheme),
}

// Basic fn
func (d *Dao) Basic(c context.Context) (basicMap map[string]interface{}, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_basicSQL, xstr.JoinInts(basicTypes)))
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	basicMap = make(map[string]interface{})
	var (
		subs   = make([]*music.Subtitle, 0)
		fons   = make([]*music.Font, 0)
		hots   = make([]*music.Hotword, 0)
		stis   = make([]*music.Sticker, 0)
		ints   = make([]*music.Intro, 0)
		trans  = make([]*music.Transition, 0)
		themes = make([]*music.Theme, 0)
	)
	basicMap["subs"] = subs
	basicMap["fons"] = fons
	basicMap["hots"] = hots
	basicMap["stis"] = stis
	basicMap["ints"] = ints
	basicMap["trans"] = trans
	basicMap["themes"] = themes
	for rows.Next() {
		b := &music.Basic{
			Material: music.Material{},
		}
		if err = rows.Scan(&b.MTime, &b.ID, &b.Name, &b.Rank, &b.Extra, &b.Material.Platform, &b.Material.Build, &b.Material.Type); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		var extra *struct {
			Cover       string `json:"cover"`
			DownloadURL string `json:"download_url"`
			Max         int    `json:"max"`
			SubType     int64  `json:"sub_type"`
			Tip         string `json:"tip"`
			White       int8   `json:"white_list"`
			TagsStr     string `json:"tags"`
		}
		if err = json.Unmarshal([]byte(b.Extra), &extra); err != nil {
			log.Error("json.Unmarshal extra failed error(%v)", err)
			continue
		}
		b.Cover = extra.Cover
		b.DownloadURL = extra.DownloadURL
		b.Max = extra.Max
		b.White = extra.White
		if len(b.Material.Build) > 0 {
			var buildComps []*music.BuildComp
			if err = json.Unmarshal([]byte(b.Build), &buildComps); err != nil {
				log.Error("json.Unmarshal buildComps failed error(%v)", err)
				continue
			}
			b.Material.BuildComps = buildComps
		}
		if b.MTime.Time().AddDate(0, 0, 7).Unix() >= time.Now().Unix() {
			b.New = 1
		}
		if len(extra.TagsStr) > 0 {
			b.Tags = append(b.Tags, strings.Split(extra.TagsStr, ",")...)
		}
		log.Info("b info v(%+v)", b)
		switch b.Material.Type {
		case appMdl.TypeSubtitle:
			subs = append(subs, &music.Subtitle{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Max:         b.Max,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeFont:
			fons = append(fons, &music.Font{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeHotWord:
			hots = append(hots, &music.Hotword{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeSticker:
			stis = append(stis, &music.Sticker{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				SubType:     extra.SubType,
				Tip:         extra.Tip,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeTheme:
			themes = append(themes, &music.Theme{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeIntro:
			ints = append(ints, &music.Intro{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		case appMdl.TypeTransition:
			trans = append(trans, &music.Transition{
				ID:          b.ID,
				Name:        b.Name,
				Cover:       b.Cover,
				DownloadURL: b.DownloadURL,
				Rank:        b.Rank,
				Material:    b.Material,
				New:         b.New,
				Tags:        b.Tags,
				MTime:       b.MTime,
			})
		}
	}
	basicMap["subs"] = subs
	basicMap["fons"] = fons
	basicMap["hots"] = hots
	basicMap["stis"] = stis
	basicMap["ints"] = ints
	basicMap["trans"] = trans
	basicMap["themes"] = themes
	return
}

// CategoryBind fn
func (d *Dao) CategoryBind(c context.Context, tp int8) (res []*music.MaterialBind, err error) {
	res = make([]*music.MaterialBind, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_CategoryBindSQL, tp))
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &music.MaterialBind{}
		if err = rows.Scan(&v.New, &v.Tp, &v.CName, &v.CRank, &v.BRank, &v.MID, &v.CID); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		res = append(res, v)
	}
	return
}

// Filters fn
func (d *Dao) Filters(c context.Context) (res []*music.Filter, resMap map[int64]*music.Filter, err error) {
	res = make([]*music.Filter, 0)
	resMap = make(map[int64]*music.Filter)
	rows, err := d.db.Query(c, _FiltersSQL)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &music.Filter{
			Material: music.Material{},
		}
		if err = rows.Scan(&v.MTime, &v.ID, &v.Name, &v.Rank, &v.Extra, &v.Material.Platform, &v.Material.Build); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		var extra *struct {
			Cover       string `json:"cover"`
			DownloadURL string `json:"download_url"`
			TagsStr     string `json:"tags"`
			FilterType  int8   `json:"filter_type"`
		}
		if err = json.Unmarshal([]byte(v.Extra), &extra); err != nil {
			log.Error("json.Unmarshal failed error(%v)", err)
			continue
		}
		v.Cover = extra.Cover
		v.DownloadURL = extra.DownloadURL
		v.FilterType = extra.FilterType
		if len(v.Material.Build) > 0 {
			var buildComps []*music.BuildComp
			if err = json.Unmarshal([]byte(v.Build), &buildComps); err != nil {
				log.Error("json.Unmarshal buildComps failed error(%v)", err)
				continue
			}
			v.Material.BuildComps = buildComps
		}
		if v.MTime.Time().AddDate(0, 0, 7).Unix() >= time.Now().Unix() {
			v.New = 1
		}
		if len(extra.TagsStr) > 0 {
			v.Tags = append(v.Tags, strings.Split(extra.TagsStr, ",")...)
		}
		log.Info("v info v(%+v)", v)
		res = append(res, v)
		resMap[v.ID] = v
	}
	return
}

// Vstickers fn
func (d *Dao) Vstickers(c context.Context) (res []*music.VSticker, resMap map[int64]*music.VSticker, err error) {
	res = make([]*music.VSticker, 0)
	resMap = make(map[int64]*music.VSticker)
	rows, err := d.db.Query(c, _VstickersSQL)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &music.VSticker{
			Material: music.Material{},
		}
		if err = rows.Scan(&v.MTime, &v.ID, &v.Name, &v.Rank, &v.Extra, &v.Material.Platform, &v.Material.Build); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		var extra *struct {
			Cover       string `json:"cover"`
			DownloadURL string `json:"download_url"`
			TagsStr     string `json:"tags"`
		}
		if err = json.Unmarshal([]byte(v.Extra), &extra); err != nil {
			log.Error("json.Unmarshal failed error(%v)", err)
			continue
		}
		v.Cover = extra.Cover
		v.DownloadURL = extra.DownloadURL
		if len(v.Material.Build) > 0 {
			var buildComps []*music.BuildComp
			if err = json.Unmarshal([]byte(v.Build), &buildComps); err != nil {
				log.Error("json.Unmarshal buildComps failed error(%v)", err)
				continue
			}
			v.Material.BuildComps = buildComps
		}
		if v.MTime.Time().AddDate(0, 0, 7).Unix() >= time.Now().Unix() {
			v.New = 1
		}
		if len(extra.TagsStr) > 0 {
			v.Tags = append(v.Tags, strings.Split(extra.TagsStr, ",")...)
		}
		log.Info("v info v(%+v)", v)
		res = append(res, v)
		resMap[v.ID] = v
	}
	return
}

// Cooperates fn
func (d *Dao) Cooperates(c context.Context) (res []*music.Cooperate, daids []int64, err error) {
	res = make([]*music.Cooperate, 0)
	rows, err := d.db.Query(c, _CooperatesSQL)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &music.Cooperate{
			Material: music.Material{},
		}
		if err = rows.Scan(&v.MTime, &v.ID, &v.Name, &v.Rank, &v.Extra, &v.Material.Platform, &v.Material.Build); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		var extra *struct {
			Cover       string `json:"cover"`
			MaterialAID int64  `json:"material_aid"`
			MaterialCID int64  `json:"material_cid"`
			DemoAID     int64  `json:"demo_aid"`
			DemoCID     int64  `json:"demo_cid"`
			MissionID   int64  `json:"mission_id"`
			TagsStr     string `json:"tags"`
			SubType     int    `json:"sub_type"`
			Style       int    `json:"style"`
			DownloadURL string `json:"download_url"`
		}
		if err = json.Unmarshal([]byte(v.Extra), &extra); err != nil {
			log.Error("json.Unmarshal failed error(%v)", err)
			continue
		}
		daids = append(daids, extra.MaterialAID)
		v.Cover = extra.Cover
		v.MaterialAID = extra.MaterialAID
		v.MaterialCID = extra.MaterialCID
		v.DemoAID = extra.DemoAID
		v.DemoCID = extra.DemoCID
		v.MissionID = extra.MissionID
		v.SubType = extra.SubType
		v.Style = extra.Style
		v.DownloadURL = extra.DownloadURL
		if len(v.Material.Build) > 0 {
			var buildComps []*music.BuildComp
			if err = json.Unmarshal([]byte(v.Build), &buildComps); err != nil {
				log.Error("json.Unmarshal buildComps failed error(%v)", err)
				continue
			}
			v.Material.BuildComps = buildComps
		}
		if v.MTime.Time().AddDate(0, 0, 7).Unix() >= time.Now().Unix() {
			v.New = 1
		}
		if len(extra.TagsStr) > 0 {
			v.Tags = append(v.Tags, strings.Split(extra.TagsStr, ",")...)
		}
		log.Info("v info v(%+v)", v)
		res = append(res, v)
	}
	return
}
