package reply

import (
	"context"

	"go-common/app/interface/main/creative/model/music"
	"go-common/app/interface/main/creative/model/reply"
	seamdl "go-common/app/interface/main/creative/model/search"
	"go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// Replies get reply list.
func (s *Service) Replies(c context.Context, p *seamdl.ReplyParam) (res *seamdl.Replies, err error) {
	if res, err = s.sear.ReplyES(c, p); err != nil {
		return
	}
	if res == nil {
		return
	}
	var (
		g, ctx       = errgroup.WithContext(c)
		mids         = res.Repliers
		oids         = res.Oids
		tyOids       = res.TyOids
		replies      map[int64]*reply.Reply
		elecRelation map[int64]int
		followers    map[int64]int
		users        map[int64]*account.Info
		arcs         map[int64]*api.Arc
		arts         map[int64]*model.Meta
		auds         map[int64]*music.Audio
	)
	log.Info("s.sear.Replies mid(%d)|type(%d)|mids(%+v)|res(%+v)", p.OMID, p.Type, mids, res)
	g.Go(func() error { //获取具体评论信息
		if len(res.DeriveIds) > 0 && len(res.DeriveOids) > 0 {
			replies, _ = s.reply.ReplyMinfo(ctx, p.Ak, p.Ck, p.OMID, int64(p.Type), res.DeriveIds, res.DeriveOids, p.IP)
		}
		return nil
	})
	g.Go(func() error { //获取被充电状态
		if len(mids) > 0 {
			elecRelation, _ = s.elec.ElecRelation(ctx, p.OMID, mids, p.IP)
		}
		return nil
	})
	g.Go(func() error { //获取被关注状态
		if len(mids) > 0 {
			followers, _ = s.acc.Followers(ctx, p.OMID, mids, p.IP)
		}
		return nil
	})
	g.Go(func() error { //获取用户信息
		if len(mids) > 0 {
			users, _ = s.acc.Infos(ctx, mids, p.IP)
		}
		return nil
	})
	g.Go(func() error { //获取各种查询对象信息
		switch p.Type {
		case seamdl.All: //查询所有
			if v, ok := tyOids[seamdl.Archive]; ok { //稿件
				g.Go(func() error {
					arcs, _ = s.arc.Archives(ctx, v, p.IP)
					return nil
				})
			}
			if v, ok := tyOids[seamdl.Article]; ok { //文章
				g.Go(func() error {
					arts, _ = s.art.ArticleMetas(ctx, v, p.IP)
					return nil
				})
			}
			if v, ok := tyOids[seamdl.Audio]; ok { //音频
				g.Go(func() error {
					auds, _ = s.mus.Audio(c, v, 0, p.IP)
					return nil
				})
			}
		case seamdl.Archive: //稿件
			arcs, _ = s.arc.Archives(ctx, oids, p.IP)
		case seamdl.SmallVideo: //小视频
		case seamdl.Article: //文章
			arts, _ = s.art.ArticleMetas(ctx, oids, p.IP)
		case seamdl.Audio: //音频
			auds, _ = s.mus.Audio(c, oids, 0, p.IP)
		}
		return nil
	})
	g.Wait()
	for _, v := range res.Result {
		if v == nil {
			continue
		}
		if p, ok := replies[v.Parent]; ok { //设置父级评论信息
			v.RootInfo = p
			v.ParentInfo = p
		}
		if elec, ok := elecRelation[v.Mid]; ok { //设置充电状态
			v.IsElec = elec
		}
		if fl, ok := followers[v.Mid]; ok { //设置关注状态
			v.Relation = fl
		}
		if u, ok := users[v.Mid]; ok { //设置图像和用户名
			v.Replier = u.Name
			v.Uface = u.Face
		}
		switch v.Type {
		case seamdl.Archive: //稿件
			if av, ok := arcs[v.Oid]; ok && av != nil {
				v.Title = av.Title
				v.Cover = av.Pic
			}
		case seamdl.SmallVideo: //小视频
		case seamdl.Article: //文章
			if art, ok := arts[v.Oid]; ok && art != nil {
				var cover string
				if len(art.ImageURLs) > 0 {
					cover = art.ImageURLs[0]
				}
				v.Title = art.Title
				v.Cover = cover
			}
		case seamdl.Audio: //音频
			if au, ok := auds[v.Oid]; ok && au != nil {
				v.Title = au.Title
				v.Cover = au.Cover
			}
		}
	}
	return
}

// AppIndexReplies get newest reply list.
func (s *Service) AppIndexReplies(c context.Context, ak, ck string, mid, oid int64, isReport, isHidden, tp, resMdlPlat int8, filterStr, kw, order, ip string, pn, ps int64) (res *seamdl.Replies, err error) {
	p := &seamdl.ReplyParam{
		Ak:          ak,
		Ck:          ck,
		OMID:        mid,
		OID:         oid,
		IsReport:    isReport,
		Type:        tp,
		FilterCtime: filterStr,
		Kw:          kw,
		Order:       order,
		IP:          metadata.String(c, metadata.RemoteIP),
		Ps:          int(ps),
		Pn:          int(pn),
		ResMdlPlat:  resMdlPlat,
	}
	if res, err = s.Replies(c, p); err != nil {
		return
	}
	if res == nil || len(res.Result) == 0 {
		return
	}
	replies := make([]*seamdl.Reply, 0, len(res.Result))
	for _, v := range res.Result {
		if v.Type == seamdl.Audio {
			continue
		}
		replies = append(replies, v)
		if len(replies) == 2 {
			break
		}
	}
	res.Result = replies
	return
}
