package service

import (
	"context"
	"encoding/json"
	"time"

	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/main/app/model"
	"go-common/app/job/main/app/model/space"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

const _sleep = 100 * time.Millisecond

// contributeConsumeproc consumer contribute
func (s *Service) contributeConsumeproc() {
	var (
		msg *databus.Message
		ok  bool
		err error
	)
	msgs := s.contributeSub.Messages()
	for {
		if msg, ok = <-msgs; !ok {
			close(s.contributeChan)
			log.Info("arc databus Consumer exit")
			break
		}
		var ms = &model.ContributeMsg{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		s.contributeChan <- ms
		msg.Commit()
	}
	s.waiter.Done()
}

func (s *Service) contributeroc() {
	var (
		ms *model.ContributeMsg
		ok bool
	)
	for {
		if ms, ok = <-s.contributeChan; !ok {
			log.Error("s.contributeChan id closed")
			break
		}
		s.contributeCache(ms.Vmid, ms.Attrs, ms.CTime, ms.IP)
	}
	s.waiter.Done()
}

func (s *Service) contributeCache(vmid int64, attrs *space.Attrs, ctime xtime.Time, ip string) (err error) {
	if vmid == 0 {
		return
	}
	var (
		items    []*space.Item
		archives []*api.Arc
		articles []*article.Meta
		clips    []*space.Clip
		albums   []*space.Album
		audios   []*space.Audio
	)
	c := context.Background()
	if attrs == nil {
		attrs = &space.Attrs{}
		if err = s.spdao.DelContributeCache(c, vmid); err != nil {
			log.Error("%+v", err)
		}
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		pn, ps := 1, 20
		for {
			var as []*api.Arc
			if err = retry(func() (err error) {
				if as, err = s.spdao.UpArchives(ctx, vmid, pn, ps, ip); err != nil {
					log.Error("%+v", err)
				}
				return
			}, _upContributeRetry, _sleep); err != nil {
				log.Error("%+v", err)
				return
			} else if len(as) == 0 {
				break
			}
			archives = append(archives, as...)
			if attrs.Archive {
				a := as[len(as)-1]
				if a != nil && a.PubDate < ctime {
					break
				}
			}
			pn++
		}
		return
	})
	g.Go(func() (err error) {
		pn, ps := 1, 20
		for {
			var ats []*article.Meta
			if err = retry(func() (err error) {
				if ats, _, err = s.spdao.UpArticles(ctx, vmid, pn, ps); err != nil {
					log.Error("%+v", err)
				}
				return
			}, _upContributeRetry, _sleep); err != nil {
				log.Error("%+v", err)
				return
			} else if len(ats) == 0 {
				break
			}
			articles = append(articles, ats...)
			if attrs.Article {
				at := ats[len(ats)-1]
				if at != nil && at.PublishTime < ctime {
					break
				}
			}
			pn++
		}
		return
	})
	g.Go(func() (err error) {
		pos, size := 0, 20
		for {
			var (
				cls    []*space.Clip
				more   int
				offset int
			)
			if err = retry(func() (err error) {
				if cls, more, offset, err = s.spdao.ClipList(ctx, vmid, pos, size, ip); err != nil {
					log.Error("%+v", err)
				}
				return
			}, _upContributeRetry, _sleep); err != nil {
				log.Error("%+v", err)
				return
			}
			if len(cls) != 0 {
				clips = append(clips, cls...)
				if attrs.Clip {
					cl := cls[len(cls)-1]
					if cl != nil && cl.CTime < ctime {
						break
					}
				}
			}
			if more != 1 {
				break
			}
			pos = offset
		}
		return
	})
	g.Go(func() (err error) {
		pos, size := 0, 20
		for {
			var (
				als    []*space.Album
				more   int
				offset int
			)
			if err = retry(func() (err error) {
				if als, more, offset, err = s.spdao.AlbumList(ctx, vmid, pos, size, ip); err != nil {
					log.Error("%+v", err)
				}
				return
			}, _upContributeRetry, _sleep); err != nil {
				log.Error("%+v", err)
				return
			}
			if len(als) != 0 {
				albums = append(albums, als...)
				if attrs.Album {
					al := als[len(als)-1]
					if al != nil && al.CTime < ctime {
						break
					}
				}
			}
			if more != 1 {
				break
			}
			pos = offset
		}
		return
	})
	g.Go(func() (err error) {
		pn, ps := 1, 20
		for {
			var (
				aus      []*space.Audio
				hasNext  bool
				nextPage int
			)
			if err = retry(func() (err error) {
				if aus, hasNext, nextPage, err = s.spdao.AudioList(c, vmid, pn, ps, ip); err != nil {
					log.Error("%+v", err)
				}
				return
			}, _upContributeRetry, _sleep); err != nil {
				log.Error("%+v", err)
				return
			}
			if len(aus) != 0 {
				audios = append(audios, aus...)
				if attrs.Audio {
					au := aus[len(aus)-1]
					if au != nil && au.CTime < ctime {
						break
					}
				}
			}
			if !hasNext {
				break
			}
			pn = nextPage
		}
		return
	})
	if err = g.Wait(); err != nil {
		s.syncRetry(c, model.ActionUpContributeAid, vmid, 0, attrs, nil, 0, "")
		log.Error("%+v", err)
		return
	}
	if len(archives) != 0 {
		attrs.Archive = true
	}
	if len(articles) != 0 {
		attrs.Article = true
	}
	if len(clips) != 0 {
		attrs.Clip = true
	}
	if len(albums) != 0 {
		attrs.Album = true
	}
	if len(audios) != 0 {
		attrs.Audio = true
	}
	items = make([]*space.Item, 0, len(archives)+len(articles)+len(clips)+len(albums)+len(audios))
	for _, a := range archives {
		if a != nil {
			item := &space.Item{ID: a.Aid, Goto: model.GotoAv, CTime: a.PubDate}
			items = append(items, item)
		}
	}
	for _, a := range articles {
		if a != nil {
			item := &space.Item{ID: a.ID, Goto: model.GotoArticle, CTime: a.PublishTime}
			items = append(items, item)
		}
	}
	for _, c := range clips {
		if c != nil {
			item := &space.Item{ID: c.VideoID, Goto: model.GotoClip, CTime: c.CTime}
			items = append(items, item)
		}
	}
	for _, a := range albums {
		if a != nil {
			item := &space.Item{ID: a.DocID, Goto: model.GotoAlbum, CTime: a.CTime}
			items = append(items, item)
		}
	}
	for _, a := range audios {
		if a != nil {
			item := &space.Item{ID: a.ID, Goto: model.GotoAudio, CTime: a.CTime}
			items = append(items, item)
		}
	}
	s.contributeUpdate(vmid, attrs, items)
	return
}

func (s *Service) contributeUpdate(vmid int64, attrs *space.Attrs, items []*space.Item) {
	var err error
	c := context.Background()
	if err = retry(func() (err error) {
		if err = s.spdao.UpContributeCache(c, vmid, attrs, items); err != nil {
			log.Error("vmid(%d) attrs(%v) contribute update failed error(%v)", vmid, attrs, err)
		} else {
			log.Info("vmid(%d) attrs(%v) contribute update success items len(%d)", vmid, attrs, len(items))
		}
		return
	}, _upContributeRetry, time.Second); err != nil {
		s.syncRetry(c, model.ActionUpContribute, vmid, 0, attrs, items, 0, "")
	}
}
