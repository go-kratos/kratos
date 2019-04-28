package region

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"go-common/app/interface/main/app-feed/conf"
	tagdao "go-common/app/interface/main/app-feed/dao/tag"
	"go-common/library/log"
)

type Service struct {
	c *conf.Config
	// dao
	tg *tagdao.Dao
	// tick
	tick time.Duration
	// infoc
	logCh chan interface{}
}

// New a region service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		tg:   tagdao.New(c),
		tick: time.Duration(c.Tick),
		// infoc
		logCh: make(chan interface{}, 1024),
	}
	go s.infocproc()
	return
}

func (s *Service) md5(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return "region_version"
	}
	hs := md5.Sum(bs)
	return hex.EncodeToString(hs[:])
}
