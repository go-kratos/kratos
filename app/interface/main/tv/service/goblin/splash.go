package goblin

import (
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) loadSphproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.PageReload))
		log.Info("Reload Splash Data!")
		s.loadSph()
	}
}

func (s *Service) loadSph() {
	var (
		err       error
		chls      []*model.Channel
		chlSplash = make(map[string]string)
	)
	// pick channel's splash data
	if chls, err = s.dao.ChlInfo(ctx); err != nil {
		log.Error("LoadSph Error (%v)", err)
		return
	}
	if len(chls) == 0 {
		log.Error("loadSph Channel Data is Empty!")
		return
	}
	// travel the channels to make the map
	for _, v := range chls {
		chlSplash[v.Title] = v.Splash
	}
	s.ChlSplash = chlSplash
	log.Info("Reload %d Channel Data", len(chlSplash))
}

// PickSph picks the splash data from memory map
func (s *Service) PickSph(channel string) (sph string, err error) {
	var ok bool
	if len(s.ChlSplash) == 0 {
		log.Error("Channel Data is Nil")
		return "", ecode.ServiceUnavailable
	}
	if sph, ok = s.ChlSplash[channel]; !ok {
		sph = s.conf.Cfg.DefaultSplash
	}
	return sph, nil
}
