package favorite

import (
	"context"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/favorite"
	fav "go-common/app/service/main/favorite/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_av       = "av"       //视频（ipad没有播单还是视频）
	_playlist = "playlist" // 播单
	_bangumi  = "bangumi"  // 追番
	_cinema   = "cinema"   // 追剧
	_topic    = "topic"    // 话题
	_article  = "article"  // 专栏
	_menu     = "menu"     // 歌单
	_pgcMenu  = "pgc_menu" // 专辑
	_clips    = "clips"    // 小视频
	_albums   = "albums"   // 相簿
	_product  = "product"  // 商品
	_ticket   = "ticket"   // 展演
	_favorite = "favorite"
)

var tabMap = map[string]*favorite.TabItem{
	_av:       {Name: "视频", Uri: "bilibili://main/favorite/video", Tab: _favorite},
	_playlist: {Name: "播单", Uri: "bilibili://main/favorite/playlist", Tab: _favorite},
	_bangumi:  {Name: "追番", Uri: "bilibili://pgc/favorite/bangumi", Tab: _bangumi},
	_cinema:   {Name: "追剧", Uri: "bilibili://pgc/favorite/cinema", Tab: _cinema},
	_topic:    {Name: "话题", Uri: "bilibili://main/favorite/topic", Tab: _topic},
	_article:  {Name: "专栏", Uri: "bilibili://column/favorite/article", Tab: _article},
	_menu:     {Name: "歌单", Uri: "bilibili://music/favorite/menu", Tab: _menu},
	_pgcMenu:  {Name: "专辑", Uri: "bilibili://music/favorite/album", Tab: _pgcMenu},
	_clips:    {Name: "小视频", Uri: "bilibili://clip/favorite", Tab: _clips},
	_albums:   {Name: "相簿", Uri: "bilibili://pictureshow/favorite", Tab: _albums},
	_product:  {Name: "商品", Uri: "bilibili://mall/favorite/goods", Tab: _product},
	_ticket:   {Name: "展演", Uri: "bilibili://mall/favorite/ticket", Tab: _ticket},
}
var tabArr = []string{_av, _playlist, _bangumi, _cinema, _topic, _article, _menu, _pgcMenu, _clips, _albums, _product, _ticket}

// Folder get my favorite.
func (s *Service) Folder(c context.Context, accessKey, actionKey, device, mobiApp, platform string, build int, aid, vmid, mid int64) (rs *favorite.MyFavorite, err error) {
	var pn, ps int = 1, 5
	rs = &favorite.MyFavorite{
		Tab: &favorite.Tab{
			Fav: true,
		},
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var (
			mediaList bool
			folders   []*favorite.Folder
		)
		plat := model.Plat(mobiApp, device)
		// 双端版本号限制，符合此条件显示为“默认收藏夹”：
		// iPhone <5.36.1(8300) 或iPhone>5.36.1(8300)
		// Android <5360001或Android>5361000
		// 双端版本号限制，符合此条件显示为“默认播单”：
		// iPhone=5.36.1(8300)
		// 5360001 <=Android <=5361000
		if (plat == model.PlatIPhone && build == 8300) || (plat == model.PlatAndroid && build >= 5360001 && build <= 5361000) {
			mediaList = true
		}
		if folders, err = s.favDao.Folders(ctx, mid, vmid, mobiApp, build, mediaList); err != nil {
			log.Error("%+v", err)
			return
		}
		if len(folders) != 0 {
			rs.Favorite = &favorite.FavList{
				Count: len(folders),
				Items: make([]*favorite.FavItem, 0, len(folders)),
			}
			for _, v := range folders {
				fi := &favorite.FavItem{}
				fi.FromFav(v)
				rs.Favorite.Items = append(rs.Favorite.Items, fi)
			}
		}
		return
	})
	g.Go(func() (err error) {
		var topic *fav.UserFolderReply
		if topic, err = s.topicDao.UserFolder(ctx, mid, 4); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if topic != nil && topic.Res != nil && topic.Res.Count > 0 {
			rs.Tab.Topic = true
		}
		return
	})
	g.Go(func() error {
		article := s.Article(ctx, mid, pn, ps)
		if article != nil && article.Count > 0 {
			rs.Tab.Article = true
		}
		return nil
	})
	g.Go(func() error {
		clips := s.Clips(ctx, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
		if clips != nil && clips.PageInfo != nil && clips.Count > 0 {
			rs.Tab.Clips = true
		}
		return nil
	})
	g.Go(func() error {
		albums := s.Albums(ctx, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
		if albums != nil && albums.PageInfo != nil && albums.Count > 0 {
			rs.Tab.Albums = true
		}
		return nil
	})
	g.Go(func() error {
		specil := s.Specil(ctx, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
		if specil != nil && specil.Count > 0 {
			rs.Tab.Specil = true
		}
		return nil
	})
	g.Go(func() (err error) {
		var has bool
		if has, err = s.bangumiDao.HasFollows(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		rs.Tab.Cinema = has
		return
	})
	g.Go(func() (err error) {
		fav, err := s.audioDao.Fav(ctx, mid)
		if err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if fav != nil {
			rs.Tab.Menu = fav.Menu
			rs.Tab.PGCMenu = fav.PGCMenu
			rs.Tab.Audios = fav.Song
		}
		return
	})
	g.Go(func() (err error) {
		var ticket int32
		if ticket, err = s.ticketDao.FavCount(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if ticket > 0 {
			rs.Tab.Ticket = true
		}
		return
	})
	g.Go(func() (err error) {
		var product int32
		if product, err = s.mallDao.FavCount(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if product > 0 {
			rs.Tab.Product = true
		}
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}

func (s *Service) FolderVideo(c context.Context, accessKey, actionKey, device, mobiApp, platform, keyword, order string, build, tid, pn, ps int, mid, fid, vmid int64) (folder *favorite.FavideoList) {
	video, err := s.favDao.FolderVideo(c, accessKey, actionKey, device, mobiApp, platform, keyword, order, build, tid, pn, ps, mid, fid, vmid)
	if err != nil {
		folder = &favorite.FavideoList{Items: []*favorite.FavideoItem{}}
		log.Error("%+v", err)
		return
	}
	folder = &favorite.FavideoList{
		Count: video.Total,
		Items: make([]*favorite.FavideoItem, 0, len(video.Archives)),
	}
	if video != nil {
		for _, v := range video.Archives {
			fi := &favorite.FavideoItem{}
			fi.FromFavideo(v)
			folder.Items = append(folder.Items, fi)
		}
	}
	return
}

func (s *Service) Topic(c context.Context, accessKey, actionKey, device, mobiApp, platform string, build, ps, pn int, mid int64) (topic *favorite.TopicList) {
	topics, err := s.topicDao.Topic(c, accessKey, actionKey, device, mobiApp, platform, build, ps, pn, mid)
	if err != nil {
		topic = &favorite.TopicList{Items: []*favorite.TopicItem{}}
		log.Error("%+v", err)
		return
	}
	topic = &favorite.TopicList{
		Count: topics.Total,
		Items: make([]*favorite.TopicItem, 0, len(topics.Lists)),
	}
	if topics != nil {
		for _, v := range topics.Lists {
			fi := &favorite.TopicItem{}
			fi.FromTopic(v)
			topic.Items = append(topic.Items, fi)
		}
	}
	return
}

func (s *Service) Article(c context.Context, mid int64, pn, ps int) (article *favorite.ArticleList) {
	articleTmp, err := s.artDao.Favorites(c, mid, pn, ps)
	if err != nil {
		article = &favorite.ArticleList{Items: []*favorite.ArticleItem{}}
		log.Error("%+v", err)
		return
	}
	article = &favorite.ArticleList{
		Count: len(articleTmp),
		Items: make([]*favorite.ArticleItem, 0, len(articleTmp)),
	}
	if len(articleTmp) != 0 {
		for _, v := range articleTmp {
			fi := &favorite.ArticleItem{}
			fi.FromArticle(v)
			article.Items = append(article.Items, fi)
		}
	}
	return
}

// Clips
func (s *Service) Clips(c context.Context, mid int64, accessKey, actionKey, device, mobiApp, platform string, build, pn, ps int) (clips *favorite.ClipsList) {
	clipsTmp, err := s.bplusDao.FavClips(c, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
	if err != nil {
		clips = &favorite.ClipsList{Items: []*favorite.ClipsItem{}}
		log.Error("%+v", err)
		return
	}
	clips = &favorite.ClipsList{
		PageInfo: clipsTmp.PageInfo,
		Items:    make([]*favorite.ClipsItem, 0, len(clipsTmp.List)),
	}
	if clipsTmp != nil {
		for _, v := range clipsTmp.List {
			fi := &favorite.ClipsItem{}
			fi.FromClips(v)
			clips.Items = append(clips.Items, fi)
		}
	}
	return
}

func (s *Service) Albums(c context.Context, mid int64, accessKey, actionKey, device, mobiApp, platform string, build, pn, ps int) (albums *favorite.AlbumsList) {
	albumsTmp, err := s.bplusDao.FavAlbums(c, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
	if err != nil {
		albums = &favorite.AlbumsList{Items: []*favorite.AlbumItem{}}
		log.Error("%+v", err)
		return
	}
	albums = &favorite.AlbumsList{
		PageInfo: albumsTmp.PageInfo,
		Items:    make([]*favorite.AlbumItem, 0, len(albumsTmp.List)),
	}
	if albumsTmp != nil {
		for _, v := range albumsTmp.List {
			fi := &favorite.AlbumItem{}
			fi.FromAlbum(v)
			albums.Items = append(albums.Items, fi)
		}
	}
	return
}

func (s *Service) Specil(c context.Context, accessKey, actionKey, device, mobiApp, platform string, build, pn, ps int) (specil *favorite.SpList) {
	specilTmp, err := s.spDao.Specil(c, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
	if err != nil {
		specil = &favorite.SpList{Items: []*favorite.SpItem{}}
		log.Error("%+v", err)
		return
	}
	specil = &favorite.SpList{
		Count: len(specilTmp.Items),
		Items: make([]*favorite.SpItem, 0, len(specilTmp.Items)),
	}
	if specilTmp != nil {
		for _, v := range specilTmp.Items {
			fi := &favorite.SpItem{}
			fi.FromSp(v)
			specil.Items = append(specil.Items, fi)
		}
	}
	return
}

func (s *Service) Audio(c context.Context, accessKey string, mid int64, pn, ps int) (audio *favorite.AudioList) {
	audioTmp, err := s.audioDao.FavAudio(c, accessKey, mid, pn, ps)
	if err != nil {
		audio = &favorite.AudioList{Items: []*favorite.AudioItem{}}
		log.Error("%+v", err)
		return
	}
	audio = &favorite.AudioList{
		Count: len(audioTmp),
		Items: make([]*favorite.AudioItem, 0, len(audioTmp)),
	}
	for _, v := range audioTmp {
		fi := &favorite.AudioItem{}
		fi.FromAudio(v)
		audio.Items = append(audio.Items, fi)
	}
	return
}

// Tab fav tab.
func (s *Service) Tab(c context.Context, accessKey, actionKey, device, mobiApp, platform, filtered string, build int, mid int64) (tab []*favorite.TabItem, err error) {
	var (
		pn, ps     = 1, 5
		tabDisplay = []string{_playlist}
	)
	plat := model.Plat(mobiApp, device)
	if model.IsIPad(plat) {
		tabDisplay = []string{_av}
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var bangumiFav, cinemaFav int
		if bangumiFav, cinemaFav, err = s.bangumiDao.FavDisplay(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if bangumiFav == 1 {
			tabDisplay = append(tabDisplay, _bangumi)
		}
		if cinemaFav == 1 {
			tabDisplay = append(tabDisplay, _cinema)
		}
		return
	})
	if !model.IsIPad(plat) {
		if filtered != "1" {
			g.Go(func() (err error) {
				var topic *fav.UserFolderReply
				if topic, err = s.topicDao.UserFolder(ctx, mid, 4); err != nil {
					log.Error("%+v", err)
					err = nil
					return
				}
				if topic != nil && topic.Res != nil && topic.Res.Count > 0 {
					tabDisplay = append(tabDisplay, _topic)
				}
				return
			})
		}
		g.Go(func() error {
			article := s.Article(ctx, mid, pn, ps)
			if article != nil && article.Count > 0 {
				tabDisplay = append(tabDisplay, _article)
			}
			return nil
		})
		g.Go(func() error {
			clips := s.Clips(ctx, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
			if clips != nil && clips.PageInfo != nil && clips.Count > 0 {
				tabDisplay = append(tabDisplay, _clips)
			}
			return nil
		})
		g.Go(func() error {
			albums := s.Albums(ctx, mid, accessKey, actionKey, device, mobiApp, platform, build, pn, ps)
			if albums != nil && albums.PageInfo != nil && albums.Count > 0 {
				tabDisplay = append(tabDisplay, _albums)
			}
			return nil
		})
		g.Go(func() (err error) {
			fav, err := s.audioDao.Fav(ctx, mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if fav != nil {
				tabDisplay = append(tabDisplay, _menu)
				tabDisplay = append(tabDisplay, _pgcMenu)
			}
			return
		})
		g.Go(func() (err error) {
			var ticket int32
			if ticket, err = s.ticketDao.FavCount(ctx, mid); err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if ticket > 0 {
				tabDisplay = append(tabDisplay, _ticket)
			}
			return
		})
		g.Go(func() (err error) {
			var product int32
			if product, err = s.mallDao.FavCount(ctx, mid); err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if product > 0 {
				tabDisplay = append(tabDisplay, _product)
			}
			return
		})
	}
	g.Wait()
	for _, t := range tabArr {
		for _, dt := range tabDisplay {
			if t == dt {
				tab = append(tab, tabMap[t])
			}
		}
	}
	return
}
