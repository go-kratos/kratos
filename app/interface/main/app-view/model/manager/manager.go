package manager

import (
	"encoding/json"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

type Relate struct {
	ID          int64               `json:"id,omitempty"`
	Param       int64               `json:"param,omitempty"`
	Goto        string              `json:"goto,omitempty"`
	Title       string              `json:"title,omitempty"`
	ResourceIDs string              `json:"resource_ids,omitempty"`
	TagIDs      string              `json:"tag_ids,omitempty"`
	ArchiveIDs  string              `json:"archive_ids,omitempty"`
	RecReason   string              `json:"rec_reason,omitempty"`
	Position    int                 `json:"position,omitempty"`
	STime       xtime.Time          `json:"stime,omitempty"`
	ETime       xtime.Time          `json:"etime,omitempty"`
	PlatVer     json.RawMessage     `json:"plat_ver,omitempty"`
	Versions    map[int8][]*Version `json:"versions,omitempty"`
	Aids        map[int64]struct{}
	Tids        map[int64]struct{}
	Rids        map[int64]struct{}
}

type Version struct {
	Plat      int8   `json:"plat,omitempty"`
	Build     int    `json:"build,omitempty"`
	Condition string `json:"conditions,omitempty"`
}

func (r *Relate) Change() {
	var (
		vs  []*Version
		err error
	)
	if r.ArchiveIDs != "" {
		var aids []int64
		if aids, err = xstr.SplitInts(r.ArchiveIDs); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", r.ArchiveIDs, err)
			return
		}
		r.Aids = make(map[int64]struct{}, len(aids))
		for _, aid := range aids {
			r.Aids[aid] = struct{}{}
		}
	}
	if r.TagIDs != "" {
		var tids []int64
		if tids, err = xstr.SplitInts(r.TagIDs); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", r.TagIDs, err)
			return
		}
		r.Tids = make(map[int64]struct{}, len(tids))
		for _, tid := range tids {
			r.Tids[tid] = struct{}{}
		}
	}
	if r.ResourceIDs != "" {
		var rids []int64
		if rids, err = xstr.SplitInts(r.ResourceIDs); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", r.ResourceIDs, err)
			return
		}
		r.Rids = make(map[int64]struct{}, len(rids))
		for _, rid := range rids {
			r.Rids[rid] = struct{}{}
		}
	}
	if err = json.Unmarshal(r.PlatVer, &vs); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.PlatVer, err)
		return
	}
	vm := make(map[int8][]*Version, len(vs))
	for _, v := range vs {
		vm[v.Plat] = append(vm[v.Plat], v)
	}
	r.Versions = vm
}
