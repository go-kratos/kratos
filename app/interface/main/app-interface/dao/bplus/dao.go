package bplus

import (
	"go-common/app/interface/main/app-interface/conf"
	"go-common/library/cache/redis"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao is favorite dao
type Dao struct {
	client        *httpx.Client
	favorPlus     string
	clips         string
	albums        string
	allClip       string
	allAlbum      string
	clipDetail    string
	albumDetail   string
	groupsCount   string
	dynamic       string
	dynamicCount  string
	dynamicDetail string
	// redis
	redis *redis.Pool
	// databus
	pub *databus.Databus
}

// New initial favorite dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:        httpx.NewClient(c.HTTPBPlus),
		favorPlus:     c.Host.APILiveCo + _favorPlus,
		clips:         c.Host.VC + _clips,
		albums:        c.Host.VC + _allbums,
		allClip:       c.Host.VC + _allClip,
		allAlbum:      c.Host.VC + _allAlbum,
		clipDetail:    c.Host.VC + _clipDetail,
		albumDetail:   c.Host.VC + _albumDetail,
		groupsCount:   c.Host.VC + _groupsCount,
		dynamic:       c.Host.VC + _dynamic,
		dynamicCount:  c.Host.VC + _dunamicCount,
		dynamicDetail: c.Host.VC + _dynamicDetail,
		redis:         redis.NewPool(c.Redis.Contribute.Config),
		pub:           databus.New(c.ContributePub),
	}
	return
}
