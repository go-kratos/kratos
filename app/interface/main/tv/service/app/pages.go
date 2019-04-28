package service

import (
	"context"
	"time"

	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// loadPagesproc loads the mod&zone&home pages
func (s *Service) loadPagesproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.PageReload))
		s.loadPages()
	}
}

func (s *Service) loadPages() {
	defer elapsed("loadPages")() // record page loading time
	pCtx := context.TODO()
	// prepare pgc & ugc data
	g, errCtx := errgroup.WithContext(pCtx)
	g.Go(func() (err error) {
		return s.filterIntervs(errCtx)
	})
	g.Go(func() (err error) {
		return s.prepareUGCData(errCtx)
	})
	g.Go(func() (err error) {
		s.preparePGCData(errCtx)
		return nil
	})
	if err := g.Wait(); err != nil {
		log.Error("loadPages PrepareData Err %v", err)
		return
	}
	// load page data
	if err := s.zonesData(pCtx); err != nil {
		log.Error("loadPages, zonesData Err %v", err)
		return
	} // zonepage depends on PGCData refresh
	if err := s.loadHome(pCtx); err != nil {
		log.Error("loadPages, homeData Err %v", err)
		return
	} // homepage depends on the zone data from Zone's data
	if err := s.loadMods(pCtx); err != nil {
		log.Error("loadPages, Module Data Err %v", err)
		return
	} // modpage depends on home's recom and zone's pgc data
}
