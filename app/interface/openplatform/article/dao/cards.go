package dao

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// TicketCard get ticket card from api
func (d *Dao) TicketCard(c context.Context, ids []int64) (resp map[int64]*model.TicketCard, err error) {
	params := url.Values{}
	params.Set("id", xstr.JoinInts(ids))
	params.Set("for", "2")
	params.Set("tag", "0")
	params.Set("price", "1")
	params.Set("imgtype", "2")
	params.Set("rettype", "1")
	var res struct {
		Code int                         `json:"errno"`
		Msg  string                      `json:"msg"`
		Data map[int64]*model.TicketCard `json:"data"`
	}
	err = d.httpClient.Get(c, d.c.Cards.TicketURL, "", params, &res)
	if err != nil {
		PromError("cards:ticket接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.TicketURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError("cards:ticket接口")
		log.Error("cards: url(%s) res code(%d) msg: %s", d.c.Cards.TicketURL+"?"+params.Encode(), res.Code, res.Msg)
		err = ecode.Int(res.Code)
		return
	}
	resp = res.Data
	return
}

// MallCard .
func (d *Dao) MallCard(c context.Context, ids []int64) (resp map[int64]*model.MallCard, err error) {
	idsStr := `{"itemsIdList":[` + xstr.JoinInts(ids) + "]}"
	req, err := http.NewRequest("POST", d.c.Cards.MallURL, strings.NewReader(idsStr))
	if err != nil {
		PromError("cards:mall接口")
		log.Error("cards: NewRequest(%s) error(%+v)", d.c.Cards.MallURL+"?"+idsStr, err)
		return
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*model.MallCard `json:"list"`
		} `json:"data"`
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "")
	err = d.httpClient.Do(c, req, &res)
	if err != nil {
		PromError("cards:mall接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.MallURL+"?"+idsStr, err)
		return
	}
	if res.Code != 0 {
		PromError("cards:mall接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.MallURL+"?"+idsStr, err)
		err = ecode.Int(res.Code)
		return
	}
	resp = make(map[int64]*model.MallCard)
	for _, l := range res.Data.List {
		resp[l.ID] = l
	}
	return
}

// AudioCard .
func (d *Dao) AudioCard(c context.Context, ids []int64) (resp map[int64]*model.AudioCard, err error) {
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	params.Set("level", "1")
	var res struct {
		Code int                        `json:"code"`
		Data map[int64]*model.AudioCard `json:"data"`
	}
	err = d.httpClient.Get(c, d.c.Cards.AudioURL, "", params, &res)
	if err != nil {
		PromError("cards:audio接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.AudioURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError("cards:audio接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.AudioURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	resp = res.Data
	return
}

// BangumiCard .
func (d *Dao) BangumiCard(c context.Context, seasonIDs []int64, episodeIDs []int64) (resp map[int64]*model.BangumiCard, err error) {
	params := url.Values{}
	params.Set("season_ids", xstr.JoinInts(seasonIDs))
	params.Set("episode_ids", xstr.JoinInts(episodeIDs))
	var res struct {
		Code int `json:"code"`
		Data struct {
			SeasonMap  map[int64]*model.BangumiCard `json:"season_map"`
			EpisodeMap map[int64]*model.BangumiCard `json:"episode_map"`
		} `json:"result"`
	}
	err = d.httpClient.Post(c, d.c.Cards.BangumiURL, "", params, &res)
	if err != nil {
		PromError("cards:bangumi接口")
		log.Error("cards: d.client.Get(%s) error(%+v)", d.c.Cards.BangumiURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		PromError("cards:bangumi接口")
		log.Error("cards: url(%s) res code(%d)", d.c.Cards.BangumiURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	resp = make(map[int64]*model.BangumiCard)
	for id, item := range res.Data.EpisodeMap {
		resp[id] = item
	}
	for id, item := range res.Data.SeasonMap {
		resp[id] = item
	}
	return
}
