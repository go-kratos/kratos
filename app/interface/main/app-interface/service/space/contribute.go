package space

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/app/interface/main/app-interface/model/space"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

const (
	_androidAudio = 516009
	_iosAudio     = 6160
)

// Contribute func
func (s *Service) Contribute(c context.Context, plat int8, build int, vmid int64, pn, ps int, now time.Time) (res *space.Contributes, err error) {
	var (
		attrs *space.Attrs
		items []*space.Item
	)
	if pn == 1 {
		var (
			ctime  xtime.Time
			cached bool
		)
		size := ps
		if items, err = s.bplusDao.RangeContributeCache(c, vmid, pn, ps); err != nil {
			log.Error("%+v", err)
		} else if len(items) != 0 {
			ctime = items[0].CTime
		} else {
			size = 100
			cached = true
		}
		if res, err = s.firstContribute(c, vmid, size, now); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		if res != nil && len(res.Items) != 0 {
			if res.Items[0].CTime > ctime {
				if err = s.bplusDao.NotifyContribute(c, vmid, nil, ctime); err != nil {
					log.Error("%+v", err)
					err = nil
				}
			}
			if cached {
				ris := res.Items
				s.addCache(func() {
					s.bplusDao.AddContributeCache(context.Background(), vmid, nil, ris)
				})
			}
			if len(items) == 0 {
				ris := make([]*space.Item, 0, ps)
				for _, item := range res.Items {
					item.FormatKey()
					switch item.Goto {
					case model.GotoAudio:
						if (plat == model.PlatAndroid && build > _androidAudio) || (plat == model.PlatIPhone && build > _iosAudio) || plat == model.PlatAndroidB {
							ris = append(ris, item)
						}
					default:
						ris = append(ris, item)
					}
					if len(ris) == ps {
						break
					}
				}
				res.Items = ris
				return
			}
		}
	} else {
		if items, err = s.bplusDao.RangeContributeCache(c, vmid, pn, ps); err != nil {
			return
		}
	}
	if len(items) != 0 {
		if attrs, err = s.bplusDao.AttrCache(c, vmid); err != nil {
			log.Error("%+v", err)
		}
		// merge res
		if res, err = s.dealContribute(c, plat, build, vmid, attrs, items, now); err != nil {
			log.Error("%+v", err)
		}
	}
	if res == nil {
		res = &space.Contributes{Tab: &space.Tab{}, Items: []*space.Item{}, Links: &space.Links{}}
	}
	return
}

// Contribution func
func (s *Service) Contribution(c context.Context, plat int8, build int, vmid int64, cursor *model.Cursor, now time.Time) (res *space.Contributes, err error) {
	var (
		attrs *space.Attrs
		items []*space.Item
	)
	if cursor.Latest() {
		var (
			ctime  xtime.Time
			cached bool
		)
		size := cursor.Size
		if items, err = s.bplusDao.RangeContributeCache(c, vmid, 1, 1); err != nil {
			log.Error("%+v", err)
		} else if len(items) != 0 {
			ctime = items[0].CTime
		} else {
			size = 100
			cached = true
		}
		if res, err = s.firstContribute(c, vmid, size, now); err != nil {
			log.Error("%+v", err)
		}
		if res != nil && len(res.Items) != 0 {
			if cached {
				ris := res.Items
				s.addCache(func() {
					s.bplusDao.AddContributeCache(context.Background(), vmid, nil, ris)
				})
			}
			if res.Items[0].CTime > ctime {
				if len(items) != 0 {
					if attrs, err = s.bplusDao.AttrCache(c, vmid); err != nil {
						log.Error("%+v", err)
					}
				}
				if err = s.bplusDao.NotifyContribute(c, vmid, attrs, ctime); err != nil {
					log.Error("%+v", err)
					err = nil
				}
			}
			ris := make([]*space.Item, 0, cursor.Size)
			for _, item := range res.Items {
				item.FormatKey()
				ris = append(ris, item)
				if len(ris) == cursor.Size {
					break
				}
			}
			if len(ris) != 0 {
				res.Items = ris
				res.Links.Link(0, int64(ris[len(ris)-1].Member))
			}
			return
		}
	}
	if items, err = s.bplusDao.RangeContributionCache(c, vmid, cursor); err != nil {
		return
	}
	if len(items) != 0 {
		if attrs, err = s.bplusDao.AttrCache(c, vmid); err != nil {
			log.Error("%+v", err)
		}
		// merge res
		if res, err = s.dealContribute(c, plat, build, vmid, attrs, items, now); err != nil {
			log.Error("%+v", err)
		}
	}
	if res == nil {
		res = &space.Contributes{Tab: &space.Tab{}, Items: []*space.Item{}, Links: &space.Links{}}
	}
	return
}

func (s *Service) firstContribute(c context.Context, vmid int64, size int, now time.Time) (res *space.Contributes, err error) {
	res = &space.Contributes{Tab: &space.Tab{}, Items: []*space.Item{}, Links: &space.Links{}}
	g, ctx := errgroup.WithContext(c)
	var arcItem, artItem, clipItem, albumItem, audioItem, items []*space.Item
	g.Go(func() (err error) {
		var arcs []*api.Arc
		if arcs, err = s.arcDao.UpArcs3(ctx, vmid, 1, size); err != nil {
			log.Error("s.arcDao.UpArcs3(%d,%d,%d) error(%v)", vmid, 1, size, err)
			err = nil
			return
		}
		if len(arcs) != 0 {
			arcItem = make([]*space.Item, 0, len(arcs))
			for _, v := range arcs {
				if v.IsNormal() {
					si := &space.Item{}
					si.FromArc3(v)
					arcItem = append(arcItem, si)
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		var arts []*article.Meta
		if arts, _, err = s.artDao.UpArticles(ctx, vmid, 1, size); err != nil {
			log.Error("s.artDao.UpArticles(%d,%d,%d) error(%v)", vmid, 1, size, err)
			err = nil
			return
		}
		if len(arts) != 0 {
			artItem = make([]*space.Item, 0, len(arts))
			for _, v := range arts {
				if v.AttrVal(article.AttrBitNoDistribute) {
					continue
				}
				si := &space.Item{}
				si.FromArticle(v)
				artItem = append(artItem, si)
			}
		}
		return
	})
	g.Go(func() (err error) {
		var clips []*bplus.Clip
		if clips, _, err = s.bplusDao.AllClip(c, vmid, size); err != nil {
			log.Error("s.bplusDao.AllClip(%d,%d) error(%v)", vmid, size, err)
			err = nil
			return
		}
		if len(clips) != 0 {
			clipItem = make([]*space.Item, 0, len(clips))
			for _, v := range clips {
				si := &space.Item{}
				si.FromClip(v)
				clipItem = append(clipItem, si)
			}
		}
		return
	})
	g.Go(func() (err error) {
		var album []*bplus.Album
		if album, _, err = s.bplusDao.AllAlbum(c, vmid, size); err != nil {
			log.Error("s.bplusDao.AllAlbum(%d,%d) error(%v)", vmid, size, err)
			err = nil
			return
		}
		if len(album) != 0 {
			albumItem = make([]*space.Item, 0, len(album))
			for _, v := range album {
				si := &space.Item{}
				si.FromAlbum(v)
				albumItem = append(albumItem, si)
			}
		}
		return
	})
	g.Go(func() (err error) {
		var audio []*audio.Audio
		if audio, err = s.audioDao.AllAudio(c, vmid); err != nil {
			log.Error("s.audioDao.AllAudio(%d) error(%v)", vmid, err)
			err = nil
			return
		}
		if len(audio) != 0 {
			audioItem = make([]*space.Item, 0, len(audio))
			for _, v := range audio {
				si := &space.Item{}
				si.FromAudio(v)
				audioItem = append(audioItem, si)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("Contribute errgroup.WithContext error(%v)", err)
	}
	items = make([]*space.Item, 0, len(arcItem)+len(artItem)+len(clipItem)+len(albumItem)+len(audioItem))
	if len(arcItem) != 0 {
		res.Tab.Archive = true
		items = append(items, arcItem...)
	}
	if len(artItem) != 0 {
		res.Tab.Article = true
		items = append(items, artItem...)
	}
	if len(clipItem) != 0 {
		res.Tab.Clip = true
		items = append(items, clipItem...)
	}
	if len(albumItem) != 0 {
		res.Tab.Album = true
		items = append(items, albumItem...)
	}
	if len(audioItem) != 0 {
		res.Tab.Audios = true
		items = append(items, audioItem...)
	}
	sort.Sort(space.Items(items))
	res.Items = items
	return
}

func (s *Service) dealContribute(c context.Context, plat int8, build int, vmid int64, attrs *space.Attrs, items []*space.Item, now time.Time) (res *space.Contributes, err error) {
	res = &space.Contributes{Tab: &space.Tab{}, Items: []*space.Item{}, Links: &space.Links{}}
	var aids, cvids, clids, alids, auids []int64
	if attrs == nil {
		attrs = &space.Attrs{}
	} else if !((plat == model.PlatAndroid && build > _androidAudio) || (plat == model.PlatIPhone && build > _iosAudio) || plat == model.PlatAndroidB) {
		attrs.Audio = false
	}
	for _, item := range items {
		if item.ID == 0 {
			continue
		}
		switch item.Goto {
		case model.GotoAv:
			aids = append(aids, item.ID)
		case model.GotoArticle:
			cvids = append(cvids, item.ID)
		case model.GotoClip:
			clids = append(clids, item.ID)
		case model.GotoAlbum:
			alids = append(alids, item.ID)
		case model.GotoAudio:
			if (plat == model.PlatAndroid && build > _androidAudio) || (plat == model.PlatIPhone && build > _iosAudio) || plat == model.PlatAndroidB {
				auids = append(auids, item.ID)
			}
		}
	}
	var (
		am  map[int64]*api.Arc
		atm map[int64]*article.Meta
		clm map[int64]*bplus.Clip
		alm map[int64]*bplus.Album
		aum map[int64]*audio.Audio
	)
	g, ctx := errgroup.WithContext(c)
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.arcDao.Archives2(ctx, aids); err != nil {
				log.Error("s.arcDao.Archives(%v) error(%v)", aids, err)
				err = nil
			}
			return
		})
	}
	if len(cvids) != 0 {
		g.Go(func() (err error) {
			if atm, err = s.artDao.Articles(ctx, cvids); err != nil {
				log.Error("s.artDao.Articles(%v) error(%v)", cvids, err)
				err = nil
			}
			return
		})
	}
	if len(clids) != 0 {
		g.Go(func() (err error) {
			if clm, err = s.bplusDao.ClipDetail(ctx, clids); err != nil {
				log.Error("s.bplusDao.ClipDetail(%v) error(%v)", clids, err)
				err = nil
			}
			return
		})
	}
	if len(alids) != 0 {
		g.Go(func() (err error) {
			if alm, err = s.bplusDao.AlbumDetail(ctx, vmid, alids); err != nil {
				log.Error("s.bplusDao.AlbumDetail(%v) error(%v)", alids, err)
				err = nil
			}
			return
		})
	}
	if len(auids) != 0 {
		g.Go(func() (err error) {
			if aum, err = s.audioDao.AudioDetail(c, auids); err != nil {
				log.Error("s.audioDao.AudioDetail(%v) error(%v)", auids, err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("Contribute errgroup.WithContext error(%v)", err)
		return
	}
	if len(am) != 0 || attrs.Archive {
		res.Tab.Archive = true
	}
	if len(atm) != 0 || attrs.Article {
		res.Tab.Article = true
	}
	if len(clm) != 0 || attrs.Clip {
		res.Tab.Clip = true
	}
	if len(alm) != 0 || attrs.Album {
		res.Tab.Album = true
	}
	if len(aum) != 0 || attrs.Audio {
		res.Tab.Audios = true
	}
	ris := make([]*space.Item, 0, len(items))
	for _, item := range items {
		ri := &space.Item{}
		switch item.Goto {
		case model.GotoAv:
			if a, ok := am[item.ID]; ok && a.IsNormal() {
				ri.FromArc3(a)
			}
		case model.GotoArticle:
			if at, ok := atm[item.ID]; ok {
				ri.FromArticle(at)
			}
		case model.GotoClip:
			if cl, ok := clm[item.ID]; ok {
				ri.FromClip(cl)
			}
		case model.GotoAlbum:
			if al, ok := alm[item.ID]; ok {
				ri.FromAlbum(al)
			}
		case model.GotoAudio:
			if au, ok := aum[item.ID]; ok {
				ri.FromAudio(au)
			}
		}
		if ri.Goto != "" {
			ri.FormatKey()
			ris = append(ris, ri)
		}
	}
	res.Items = ris
	res.Links.Link(int64(items[0].Member), int64(items[len(items)-1].Member))
	return
}

// Clip func
func (s *Service) Clip(c context.Context, vmid int64, pos, size int) (res *space.ClipList) {
	var (
		clips []*bplus.Clip
		err   error
	)
	res = &space.ClipList{Item: []*space.Item{}}
	if clips, res.More, res.Offset, err = s.bplusDao.Clips(c, vmid, pos, size); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(clips) > 0 {
		res.Item = make([]*space.Item, 0, len(clips))
		for _, v := range clips {
			si := &space.Item{}
			si.FromClip(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// Album func
func (s *Service) Album(c context.Context, vmid int64, pos, size int) (res *space.AlbumList) {
	var (
		albums []*bplus.Album
		err    error
	)
	res = &space.AlbumList{Item: []*space.Item{}}
	if albums, res.More, res.Offset, err = s.bplusDao.Albums(c, vmid, pos, size); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(albums) > 0 {
		res.Item = make([]*space.Item, 0, len(albums))
		for _, v := range albums {
			si := &space.Item{}
			si.FromAlbum(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// AddContribute func
func (s *Service) AddContribute(c context.Context, vmid int64, attrs *space.Attrs, items []*space.Item) (err error) {
	if err = s.bplusDao.AddContributeCache(c, vmid, attrs, items); err != nil {
		log.Error("%+v", err)
	}
	return
}
