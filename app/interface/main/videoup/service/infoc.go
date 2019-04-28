package service

import (
	"context"
	"encoding/json"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"strconv"
)

type infoc struct {
	Aid        string          `json:"aid"`
	Ext2       json.RawMessage `json:"ext2"`
	Ext1       json.RawMessage `json:"ext1"`
	Ext3       json.RawMessage `json:"ext3"`
	Mid        string          `json:"mid"`
	Cid        string          `json:"cid"`
	Filename   string          `json:"filename"`
	Upfrom     string          `json:"upfrom"`
	PicCount   string          `json:"pic_count"`
	VideoCount string          `json:"video_count"`
	Build      string          `json:"build"`
	Platform   string          `json:"platform"`
	Device     string          `json:"device"`
	MobiApp    string          `json:"mobi_app"`
	// none business fields
	IP    string `json:"ip"`
	LogID string `json:"logid"`
	Name  string `json:"name"`
}

// VideoInfoc fn
func (s *Service) VideoInfoc(c context.Context, ap *archive.ArcParam, ar *archive.AppRequest) (err error) {
	log.Warn("infocproc begin ap(%+v) ar(%+v)", ap, ar)
	ip := metadata.String(c, metadata.RemoteIP)
	name := "APP投稿分P的视频和图片的计数"
	logID := "001729"
	for _, v := range ap.Videos {
		if v.Editor == nil || v.Cid == 0 {
			continue
		}
		infoc := &infoc{
			Name:       name,
			Mid:        strconv.FormatInt(ap.Mid, 10),
			Aid:        strconv.FormatInt(ap.Aid, 10),
			Cid:        strconv.FormatInt(v.Cid, 10),
			Filename:   v.Filename,
			Upfrom:     strconv.Itoa(int(ap.UpFrom)),
			PicCount:   strconv.Itoa(int(v.Editor.PicCount)),
			VideoCount: strconv.Itoa(int(v.Editor.VideoCount)),
			MobiApp:    ar.MobiApp,
			Platform:   ar.Platform,
			Build:      ar.Build,
			Device:     ar.Device,
			IP:         ip,
			LogID:      logID,
		}
		log.Warn("infocproc create infoc ap(%+v) ar(%+v) infoc(%+v)", ap, ar, infoc)
		err = s.infoc.Info(
			infoc.Aid,
			"",
			"",
			"",
			infoc.Mid,
			infoc.Cid,
			infoc.Filename,
			infoc.Upfrom,
			infoc.PicCount,
			infoc.VideoCount,
			infoc.Build,
			infoc.Platform,
			infoc.Device,
			infoc.MobiApp,
		)
		log.Warn("infocproc end infoc ap(%+v) ar(%+v) infoc(%+v)|err(%+v)", ap, ar, infoc, err)
	}
	return
}
