package service

import (
	"context"

	"go-common/app/interface/main/space/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_typeFavAlbum = 2
	_typeFavMovie = 2
)

var _emptyArcFavFolder = make([]*favmdl.VideoFolder, 0)

// FavNav get fav info.
func (s *Service) FavNav(c context.Context, mid int64, vmid int64) (res *model.FavNav, err error) {
	var (
		folder                                              []*favmdl.VideoFolder
		plData, topicData, artData                          *favmdl.Favorites
		albumCount, movieCount                              int
		folErr, plErr, topicErr, artErr, albumErr, movieErr error
		ip                                                  = metadata.String(c, metadata.RemoteIP)
	)
	group, errCtx := errgroup.WithContext(c)
	res = new(model.FavNav)
	// video fav folder
	if mid == vmid || s.privacyCheck(c, vmid, model.PcyFavVideo) == nil {
		group.Go(func() error {
			if folder, folErr = s.dao.FavFolder(errCtx, mid, vmid); folErr != nil {
				log.Error("s.dao.FavFolder(%d) error(%v)", vmid, folErr)
			} else {
				res.Archive = folder
			}
			return nil
		})
	}
	// playlist
	group.Go(func() error {
		arg := &favmdl.ArgFavs{Type: favmdl.TypePlayList, Mid: vmid, Pn: _samplePn, Ps: _samplePs, RealIP: ip}
		if plData, plErr = s.fav.Favorites(errCtx, arg); plErr != nil {
			log.Error("s.fav.Favorites TypePlayVideo (%d) error(%v)", vmid, plErr)
		} else if plData != nil {
			res.Playlist = plData.Page.Count
		}
		return nil
	})
	// topic
	group.Go(func() error {
		arg := &favmdl.ArgFavs{Type: favmdl.TypeTopic, Mid: vmid, Pn: _samplePn, Ps: _samplePs, RealIP: ip}
		if topicData, topicErr = s.fav.Favorites(errCtx, arg); topicErr != nil {
			log.Error("s.fav.Favorites TypeTopic (%d) error(%v)", vmid, topicErr)
		} else if topicData != nil {
			res.Topic = topicData.Page.Count
		}
		return nil
	})
	// article
	group.Go(func() error {
		arg := &favmdl.ArgFavs{Type: favmdl.Article, Mid: vmid, Pn: _samplePn, Ps: _samplePs, RealIP: ip}
		if artData, artErr = s.fav.Favorites(errCtx, arg); artErr != nil {
			log.Error("s.fav.Favorites Article (%d) error(%v)", vmid, artErr)
		} else if artData != nil {
			res.Article = artData.Page.Count
		}
		return nil
	})
	// album
	group.Go(func() error {
		if albumCount, albumErr = s.dao.LiveFavCount(errCtx, vmid, _typeFavAlbum); albumErr != nil {
			log.Error("s.dao.LiveFavCount(%d,%d) error(%v)", vmid, _typeFavAlbum, albumErr)
		} else if albumCount > 0 {
			res.Album = albumCount
		}
		return nil
	})
	// movie
	if mid > 0 {
		group.Go(func() error {
			if movieCount, movieErr = s.dao.MovieFavCount(errCtx, mid, _typeFavMovie); movieErr != nil {
				log.Error("s.dao.MovieFavCount(%d,%d) error(%v)", vmid, _typeFavMovie, movieErr)
			} else if movieCount > 0 {
				res.Movie = movieCount
			}
			return nil
		})
	}
	group.Wait()
	if len(res.Archive) == 0 {
		res.Archive = _emptyArcFavFolder
	}
	return
}

// FavArchive get favorite archive.
func (s *Service) FavArchive(c context.Context, mid int64, arg *model.FavArcArg) (res *favmdl.SearchArchive, err error) {
	if mid != arg.Vmid {
		if err = s.privacyCheck(c, arg.Vmid, model.PcyFavVideo); err != nil {
			return
		}
	}
	return s.dao.FavArchive(c, mid, arg)
}
