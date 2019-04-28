package notice

import (
	"go-common/app/job/main/reply/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao activity dao.
type Dao struct {
	c                 *conf.Config
	urlLiveSmallVideo string
	urlLiveActivity   string
	urlLiveNotice     string
	urlLivePicture    string
	urlCredit         string
	urlTopic          string
	urlActivity       string
	urlActivitySub    string
	urlDrwayoo        string
	urlDynamic        string
	urlNotice         string
	urlBan            string
	urlBangumi        string
	urlAudio          string
	urlAudioPlaylist  string
	httpClient        *bm.Client
	drawyooHTTPClient *bm.Client
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	return &Dao{
		c: c,
		// http
		urlLiveSmallVideo: c.Host.LiveVC + "/clip/v1/video/detail",
		urlLiveActivity:   c.Host.LiveAct + "/comment/v1/relation/get_by_id",
		urlLiveNotice:     c.Host.LiveVC + "/news/v1/notice/info",
		urlLivePicture:    c.Host.LiveVC + "/link_draw/v1/doc/detail",
		urlCredit:         c.Host.API + "/x/internal/credit/blocked/cases",
		urlTopic:          c.Host.Activity + "/activity/page/one/%d",
		urlActivity:       c.Host.Activity + "/activity/page/one/%d",
		urlActivitySub:    c.Host.Activity + "/activity/subject/url",
		urlDrwayoo:        c.Host.DrawYoo + "/api/pushS",
		urlDynamic:        c.Host.LiveVC + "/dynamic_repost/v0/dynamic_repost/ftch_rp_cont?dynamic_ids[]=%d",
		urlNotice:         c.Host.API + "/x/internal/credit/publish/infos",
		urlBan:            c.Host.API + "/x/internal/credit/blocked/infos",
		urlBangumi:        c.Host.Bangumi + "/api/inner/aid_episodes_v2",
		urlAudio:          c.Host.API + "/x/internal/v1/audio/songs/batch",
		urlAudioPlaylist:  c.Host.API + "/x/internal/v1/audio/menus/%d",
		httpClient:        bm.NewClient(c.HTTPClient),
		drawyooHTTPClient: bm.NewClient(c.DrawyooHTTPClient),
	}
}
