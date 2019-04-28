package service

import (
	"context"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/sync/errgroup"
)

// FindCard find card from source
func (s *Service) FindCard(c context.Context, idStr string) (res interface{}, err error) {
	var id int64
	if len(idStr) < 3 {
		err = ecode.RequestErr
		return
	}
	id, _ = strconv.ParseInt(idStr[2:], 10, 64)
	prefix := idStr[:2]
	if prefix == model.CardPrefixAudio {
		resp, err1 := s.dao.AudioCard(c, []int64{id})
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixBangumi {
		resp, err1 := s.dao.BangumiCard(c, []int64{id}, nil)
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixBangumiEp {
		resp, err1 := s.dao.BangumiCard(c, nil, []int64{id})
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixTicket {
		resp, err1 := s.dao.TicketCard(c, []int64{id})
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixMall {
		resp, err1 := s.dao.MallCard(c, []int64{id})
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixArchive {
		resp, err1 := s.Archives(c, []int64{id}, "")
		if err1 != nil {
			return nil, err1
		}
		if resp != nil {
			res = resp[id]
		}
		return res, err1
	}
	if prefix == model.CardPrefixArticle {
		resp, err1 := s.ArticleMeta(c, id)
		if err1 != nil {
			return nil, err1
		}
		return resp, err1
	}
	err = ecode.RequestErr
	return
}

// FindCards find cards
func (s *Service) FindCards(c context.Context, ids []string) (res map[string]interface{}, err error) {
	var (
		bangumis, eps, audios, malls, tickets, archives, articles []int64
		bangumiRes                                                map[int64]*model.BangumiCard
		epRes                                                     map[int64]*model.BangumiCard
		audioRes                                                  map[int64]*model.AudioCard
		mallRes                                                   map[int64]*model.MallCard
		ticketRes                                                 map[int64]*model.TicketCard
		archiveRes                                                map[int64]*api.Arc
		articleRes                                                map[int64]*model.Meta
	)
	for _, idStr := range ids {
		if len(idStr) < 2 {
			continue
		}
		id, _ := strconv.ParseInt(idStr[2:], 10, 64)
		switch idStr[:2] {
		case model.CardPrefixAudio:
			audios = append(audios, id)
		case model.CardPrefixBangumi:
			bangumis = append(bangumis, id)
		case model.CardPrefixBangumiEp:
			eps = append(eps, id)
		case model.CardPrefixMall:
			malls = append(malls, id)
		case model.CardPrefixTicket:
			tickets = append(tickets, id)
		case model.CardPrefixArchive:
			archives = append(archives, id)
		case model.CardPrefixArticle:
			articles = append(articles, id)
		}
	}
	group := errgroup.Group{}
	group.Go(func() (err error) {
		if len(bangumis) < 1 {
			return nil
		}
		if bangumiRes, err = s.dao.BangumiCard(c, bangumis, nil); err == nil {
			cache.Save(func() {
				s.dao.AddBangumiCardsCache(context.TODO(), bangumiRes)
			})
		} else {
			bangumiRes, _ = s.dao.BangumiCardsCache(c, bangumis)
		}
		return nil
	})
	group.Go(func() (err error) {
		if len(eps) < 1 {
			return nil
		}
		if epRes, err = s.dao.BangumiCard(c, nil, eps); err == nil {
			cache.Save(func() {
				s.dao.AddBangumiEpCardsCache(context.TODO(), epRes)
			})
		} else {
			epRes, _ = s.dao.BangumiEpCardsCache(c, eps)
		}
		return nil
	})
	group.Go(func() (err error) {
		if len(audios) < 1 {
			return nil
		}
		if audioRes, err = s.dao.AudioCard(c, audios); err == nil {
			cache.Save(func() {
				s.dao.AddAudioCardsCache(context.TODO(), audioRes)
			})
		} else {
			audioRes, _ = s.dao.AudioCardsCache(c, audios)
		}
		return nil
	})
	group.Go(func() (err error) {
		if len(malls) < 1 {
			return nil
		}
		if mallRes, err = s.dao.MallCard(c, malls); err == nil {
			cache.Save(func() {
				s.dao.AddMallCardsCache(context.TODO(), mallRes)
			})
		} else {
			mallRes, _ = s.dao.MallCardsCache(c, malls)
		}
		return nil
	})
	group.Go(func() (err error) {
		if len(tickets) < 1 {
			return nil
		}
		if ticketRes, err = s.dao.TicketCard(c, tickets); err == nil {
			cache.Save(func() {
				s.dao.AddTicketCardsCache(context.TODO(), ticketRes)
			})
		} else {
			ticketRes, _ = s.dao.TicketCardsCache(c, tickets)
		}
		return nil
	})
	group.Go(func() (err error) {
		if len(archives) < 1 {
			return nil
		}
		archiveRes, _ = s.Archives(c, archives, "")
		return nil
	})
	group.Go(func() (err error) {
		if len(articles) < 1 {
			return nil
		}
		articleRes, _ = s.ArticleMetas(c, articles)
		return nil
	})
	group.Wait()
	res = make(map[string]interface{})
	for id, v := range bangumiRes {
		res[model.CardPrefixBangumi+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range epRes {
		res[model.CardPrefixBangumiEp+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range ticketRes {
		res[model.CardPrefixTicket+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range mallRes {
		res[model.CardPrefixMall+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range audioRes {
		res[model.CardPrefixAudio+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range archiveRes {
		res[model.CardPrefixArchive+strconv.FormatInt(id, 10)] = v
	}
	for id, v := range articleRes {
		res[model.CardPrefixArticle+strconv.FormatInt(id, 10)] = v
	}
	return
}
