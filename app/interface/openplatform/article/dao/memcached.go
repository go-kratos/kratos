package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/library/sync/errgroup"
)

const (
	_prefixArtMeta     = "art_mp_%d"
	_prefixArtContent  = "art_c_%d"
	_prefixArtKeywords = "art_kw_%d"
	_prefixArtStat     = "art_s_%d"
	_prefixCard        = "art_cards_"
	_bulkSize          = 50
)

func artMetaKey(id int64) string {
	return fmt.Sprintf(_prefixArtMeta, id)
}

func artContentKey(id int64) string {
	return fmt.Sprintf(_prefixArtContent, id)
}

func artKeywordsKey(id int64) string {
	return fmt.Sprintf(_prefixArtKeywords, id)
}

func artStatsKey(id int64) string {
	return fmt.Sprintf(_prefixArtStat, id)
}

func cardKey(id string) string {
	return _prefixCard + id
}

func hotspotsKey() string {
	return fmt.Sprintf("art_hotspots")
}

func mcHotspotKey(id int64) string {
	return fmt.Sprintf("art_hotspot_%d", id)
}

func mcAuthorKey(mid int64) string {
	return fmt.Sprintf("art_author_%d", mid)
}

func mcTagKey(tag int64) string {
	return fmt.Sprintf("tag_aids_%d", tag)
}

func mcUpStatKey(mid int64) string {
	var (
		hour int
		day  int
	)
	now := time.Now()
	hour = now.Hour()
	if hour < 7 {
		day = now.Add(time.Hour * -24).Day()
	} else {
		day = now.Day()
	}
	return fmt.Sprintf("up_stat_daily_%d_%d", mid, day)
}

// statsValue convert stats to string, format: "view,favorite,like,unlike,reply..."
func statsValue(s *model.Stats) string {
	if s == nil {
		return ",,,,,,"
	}
	ids := []int64{s.View, s.Favorite, s.Like, s.Dislike, s.Reply, s.Share, s.Coin}
	return xstr.JoinInts(ids)
}

func revoverStatsValue(c context.Context, s string) (res *model.Stats) {
	var (
		vs  []int64
		err error
	)
	res = new(model.Stats)
	if s == "" {
		return
	}
	if vs, err = xstr.SplitInts(s); err != nil || len(vs) < 7 {
		PromError("mc:stats解析")
		log.Error("dao.revoverStatsValue(%s) err: %+v", s, err)
		return
	}
	res = &model.Stats{
		View:     vs[0],
		Favorite: vs[1],
		Like:     vs[2],
		Dislike:  vs[3],
		Reply:    vs[4],
		Share:    vs[5],
		Coin:     vs[6],
	}
	return
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcArticleExpire}
	err = conn.Set(&item)
	return
}

//AddArticlesMetaCache add articles meta cache
func (d *Dao) AddArticlesMetaCache(c context.Context, vs ...*model.Meta) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		item := &memcache.Item{Key: artMetaKey(v.ID), Object: v, Flags: memcache.FlagProtobuf, Expiration: d.mcArticleExpire}
		if err = conn.Set(item); err != nil {
			PromError("mc:增加文章meta缓存")
			log.Error("conn.Store(%s) error(%+v)", artMetaKey(v.ID), err)
			return
		}
	}
	return
}

// ArticleMetaCache gets article's meta cache.
func (d *Dao) ArticleMetaCache(c context.Context, aid int64) (res *model.Meta, err error) {
	var (
		conn = d.mc.Get(c)
		key  = artMetaKey(aid)
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("article-meta")
			err = nil
			return
		}
		PromError("mc:获取文章meta缓存")
		log.Error("conn.Get(%v) error(%+v)", key, err)
		return
	}
	res = &model.Meta{}
	if err = conn.Scan(reply, res); err != nil {
		PromError("mc:文章meta缓存json解析")
		log.Error("reply.Scan(%s) error(%+v)", reply.Value, err)
		return
	}
	res.Strong()
	cachedCount.Incr("article-meta")
	return
}

//ArticlesMetaCache articles meta cache
func (d *Dao) ArticlesMetaCache(c context.Context, ids []int64) (cached map[int64]*model.Meta, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.Meta, len(ids))
	allKeys := make([]string, 0, len(ids))
	idmap := make(map[string]int64, len(ids))
	for _, id := range ids {
		k := artMetaKey(id)
		allKeys = append(allKeys, k)
		idmap[k] = id
	}

	group, errCtx := errgroup.WithContext(c)
	mutex := sync.Mutex{}
	keysLen := len(allKeys)
	for i := 0; i < keysLen; i += _bulkSize {
		var keys []string
		if (i + _bulkSize) > keysLen {
			keys = allKeys[i:]
		} else {
			keys = allKeys[i : i+_bulkSize]
		}

		group.Go(func() (err error) {
			conn := d.mc.Get(errCtx)
			defer conn.Close()
			replys, err := conn.GetMulti(keys)
			if err != nil {
				PromError("mc:获取文章meta缓存")
				log.Error("conn.Gets(%v) error(%+v)", keys, err)
				err = nil
				return
			}
			for key, item := range replys {
				art := &model.Meta{}
				if err = conn.Scan(item, art); err != nil {
					PromError("mc:文章meta缓存json解析")
					log.Error("item.Scan(%s) error(%+v)", item.Value, err)
					err = nil
					continue
				}
				mutex.Lock()
				cached[idmap[key]] = art.Strong()
				delete(idmap, key)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	missed = make([]int64, 0, len(idmap))
	for _, id := range idmap {
		missed = append(missed, id)
	}
	missedCount.Add("article-meta", int64(len(missed)))
	cachedCount.Add("article-meta", int64(len(cached)))
	return
}

// AddArticleStatsCache batch set article cache.
func (d *Dao) AddArticleStatsCache(c context.Context, id int64, v *model.Stats) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	bs := []byte(statsValue(v))
	item := &memcache.Item{Key: artStatsKey(id), Value: bs, Expiration: d.mcStatsExpire}
	if err = conn.Set(item); err != nil {
		PromError("mc:增加文章统计缓存")
		log.Error("conn.Store(%s) error(%+v)", artStatsKey(id), err)
	}
	return
}

//AddArticleContentCache add article content cache
func (d *Dao) AddArticleContentCache(c context.Context, id int64, content string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	var bs = []byte(content)
	item := &memcache.Item{Key: artContentKey(id), Value: bs, Expiration: d.mcArticleExpire, Flags: memcache.FlagGzip}
	if err = conn.Set(item); err != nil {
		PromError("mc:增加文章内容缓存")
		log.Error("conn.Store(%s) error(%+v)", artContentKey(id), err)
	}
	return
}

// AddArticleKeywordsCache add article keywords cache.
func (d *Dao) AddArticleKeywordsCache(c context.Context, id int64, keywords string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	var bs = []byte(keywords)
	item := &memcache.Item{Key: artKeywordsKey(id), Value: bs, Expiration: d.mcArticleExpire, Flags: memcache.FlagGzip}
	if err = conn.Set(item); err != nil {
		PromError("mc:增加文章关键字缓存")
		log.Error("conn.Store(%s) error(%+v)", artKeywordsKey(id), err)
	}
	return
}

// ArticleContentCache article content cache
func (d *Dao) ArticleContentCache(c context.Context, id int64) (res string, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	reply, err := conn.Get(artContentKey(id))
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("article-content")
			err = nil
			return
		}
		PromError("mc:获取文章内容缓存")
		log.Error("conn.Get(%v) error(%+v)", artContentKey(id), err)
		return
	}
	err = conn.Scan(reply, &res)
	return
}

// ArticleKeywordsCache article Keywords cache
func (d *Dao) ArticleKeywordsCache(c context.Context, id int64) (res string, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	reply, err := conn.Get(artKeywordsKey(id))
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("article-keywords")
			err = nil
			return
		}
		PromError("mc:获取文章关键字缓存")
		log.Error("conn.Get(%v) error(%+v)", artKeywordsKey(id), err)
		return
	}
	err = conn.Scan(reply, &res)
	return
}

//DelArticleMetaCache delete article meta cache
func (d *Dao) DelArticleMetaCache(c context.Context, id int64) (err error) {
	var (
		key  = artMetaKey(id)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			PromError("mc:删除文章meta缓存")
			log.Error("key(%v) error(%+v)", key, err)
		}
	}
	return
}

// DelArticleStatsCache delete article stats cache
func (d *Dao) DelArticleStatsCache(c context.Context, id int64) (err error) {
	var (
		key  = artStatsKey(id)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			PromError("mc:删除文章stats缓存")
			log.Error("key(%v) error(%+v)", key, err)
		}
	}
	return
}

//DelArticleContentCache delete article content cache
func (d *Dao) DelArticleContentCache(c context.Context, id int64) (err error) {
	var (
		key  = artContentKey(id)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			PromError("mc:删除文章content缓存")
			log.Error("key(%v) error(%+v)", key, err)
		}
	}
	return
}

// ArticleStatsCache article stats cache
func (d *Dao) ArticleStatsCache(c context.Context, id int64) (res *model.Stats, err error) {
	if id == 0 {
		err = ecode.NothingFound
		return
	}
	var (
		conn     = d.mc.Get(c)
		key      = artStatsKey(id)
		statsStr string
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			res = nil
			err = nil
			return
		}
		PromError("mc:获取文章计数缓存")
		log.Error("conn.Get(%v) error(%+v)", key, err)
		return
	}
	if err = conn.Scan(reply, &statsStr); err == nil {
		res = revoverStatsValue(c, statsStr)
	} else {
		PromError("mc:获取文章计数缓存")
		log.Error("dao.ArticleStatsCache.reply.Scan(%v, %v) error(%+v)", key, statsStr, err)
	}
	return
}

// ArticlesStatsCache articles stats cache
func (d *Dao) ArticlesStatsCache(c context.Context, ids []int64) (cached map[int64]*model.Stats, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.Stats, len(ids))
	allKeys := make([]string, 0, len(ids))
	idmap := make(map[string]int64, len(ids))
	for _, id := range ids {
		k := artStatsKey(id)
		allKeys = append(allKeys, k)
		idmap[k] = id
	}

	group, errCtx := errgroup.WithContext(c)
	mutex := sync.Mutex{}
	keysLen := len(allKeys)
	for i := 0; i < keysLen; i += _bulkSize {
		var keys []string
		if (i + _bulkSize) > keysLen {
			keys = allKeys[i:]
		} else {
			keys = allKeys[i : i+_bulkSize]
		}

		group.Go(func() (err error) {
			conn := d.mc.Get(errCtx)
			defer conn.Close()
			replys, err := conn.GetMulti(keys)
			if err != nil {
				PromError("mc:获取文章计数缓存")
				log.Error("conn.Gets(%v) error(%+v)", keys, err)
				err = nil
				return
			}
			for _, reply := range replys {
				var info string
				if e := conn.Scan(reply, &info); e != nil {
					PromError("mc:获取文章计数缓存scan")
					continue
				}
				art := revoverStatsValue(c, info)
				mutex.Lock()
				cached[idmap[reply.Key]] = art
				delete(idmap, reply.Key)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	missed = make([]int64, 0, len(idmap))
	for _, id := range idmap {
		missed = append(missed, id)
	}
	missedCount.Add("article-stats", int64(len(missed)))
	cachedCount.Add("article-stats", int64(len(cached)))
	return
}

// AddCardsCache .
func (d *Dao) addCardsCache(c context.Context, vs ...*model.Cards) (err error) {
	if len(vs) == 0 {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		key := cardKey(v.Key())
		item := memcache.Item{Key: key, Object: v, Expiration: d.mcCardsExpire, Flags: memcache.FlagJSON}
		if err = conn.Set(&item); err != nil {
			PromError("mc:增加卡片缓存")
			log.Error("conn.Set(%s) error(%+v)", key, err)
			return
		}
	}
	return
}

// CardsCache ids like cv123 av123 au123
func (d *Dao) cardsCache(c context.Context, ids []string) (res map[string]*model.Cards, err error) {
	if len(ids) == 0 {
		return
	}
	res = make(map[string]*model.Cards, len(ids))
	var keys []string
	for _, id := range ids {
		keys = append(keys, cardKey(id))
	}
	conn := d.mc.Get(c)
	replys, err := conn.GetMulti(keys)
	defer conn.Close()
	if err != nil {
		PromError("mc:获取cards缓存")
		log.Error("conn.Gets(%v) error(%+v)", keys, err)
		err = nil
		return
	}
	for _, reply := range replys {
		s := model.Cards{}
		if err = conn.Scan(reply, &s); err != nil {
			PromError("获取cards缓存json解析")
			log.Error("json.Unmarshal(%v) error(%+v)", reply.Value, err)
			err = nil
			continue
		}
		res[strings.TrimPrefix(reply.Key, _prefixCard)] = &s
	}
	return
}

// AddBangumiCardsCache .
func (d *Dao) AddBangumiCardsCache(c context.Context, vs map[int64]*model.BangumiCard) (err error) {
	var cards []*model.Cards
	for _, v := range vs {
		cards = append(cards, &model.Cards{Type: model.CardPrefixBangumi, BangumiCard: v})
	}
	err = d.addCardsCache(c, cards...)
	return
}

// BangumiCardsCache .
func (d *Dao) BangumiCardsCache(c context.Context, ids []int64) (vs map[int64]*model.BangumiCard, err error) {
	var cards map[string]*model.Cards
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, model.CardPrefixBangumi+strconv.FormatInt(id, 10))
	}
	if cards, err = d.cardsCache(c, idsStr); err != nil {
		return
	}
	vs = make(map[int64]*model.BangumiCard)
	for _, card := range cards {
		if (card != nil) && (card.BangumiCard != nil) {
			vs[card.BangumiCard.ID] = card.BangumiCard
		}
	}
	return
}

// AddBangumiEpCardsCache .
func (d *Dao) AddBangumiEpCardsCache(c context.Context, vs map[int64]*model.BangumiCard) (err error) {
	var cards []*model.Cards
	for _, v := range vs {
		cards = append(cards, &model.Cards{Type: model.CardPrefixBangumiEp, BangumiCard: v})
	}
	err = d.addCardsCache(c, cards...)
	return
}

// BangumiEpCardsCache .
func (d *Dao) BangumiEpCardsCache(c context.Context, ids []int64) (vs map[int64]*model.BangumiCard, err error) {
	var cards map[string]*model.Cards
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, model.CardPrefixBangumiEp+strconv.FormatInt(id, 10))
	}
	if cards, err = d.cardsCache(c, idsStr); err != nil {
		return
	}
	vs = make(map[int64]*model.BangumiCard)
	for _, card := range cards {
		if (card != nil) && (card.BangumiCard != nil) {
			vs[card.BangumiCard.ID] = card.BangumiCard
		}
	}
	return
}

// AddAudioCardsCache .
func (d *Dao) AddAudioCardsCache(c context.Context, vs map[int64]*model.AudioCard) (err error) {
	var cards []*model.Cards
	for _, v := range vs {
		cards = append(cards, &model.Cards{Type: model.CardPrefixAudio, AudioCard: v})
	}
	err = d.addCardsCache(c, cards...)
	return
}

// AudioCardsCache .
func (d *Dao) AudioCardsCache(c context.Context, ids []int64) (vs map[int64]*model.AudioCard, err error) {
	var cards map[string]*model.Cards
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, model.CardPrefixAudio+strconv.FormatInt(id, 10))
	}
	if cards, err = d.cardsCache(c, idsStr); err != nil {
		return
	}
	vs = make(map[int64]*model.AudioCard)
	for _, card := range cards {
		if (card != nil) && (card.AudioCard != nil) {
			vs[card.AudioCard.ID] = card.AudioCard
		}
	}
	return
}

// AddMallCardsCache .
func (d *Dao) AddMallCardsCache(c context.Context, vs map[int64]*model.MallCard) (err error) {
	var cards []*model.Cards
	for _, v := range vs {
		cards = append(cards, &model.Cards{Type: model.CardPrefixMall, MallCard: v})
	}
	err = d.addCardsCache(c, cards...)
	return
}

// MallCardsCache .
func (d *Dao) MallCardsCache(c context.Context, ids []int64) (vs map[int64]*model.MallCard, err error) {
	var cards map[string]*model.Cards
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, model.CardPrefixMall+strconv.FormatInt(id, 10))
	}
	if cards, err = d.cardsCache(c, idsStr); err != nil {
		return
	}
	vs = make(map[int64]*model.MallCard)
	for _, card := range cards {
		if (card != nil) && (card.MallCard != nil) {
			vs[card.MallCard.ID] = card.MallCard
		}
	}
	return
}

// AddTicketCardsCache .
func (d *Dao) AddTicketCardsCache(c context.Context, vs map[int64]*model.TicketCard) (err error) {
	var cards []*model.Cards
	for _, v := range vs {
		cards = append(cards, &model.Cards{Type: model.CardPrefixTicket, TicketCard: v})
	}
	err = d.addCardsCache(c, cards...)
	return
}

// TicketCardsCache .
func (d *Dao) TicketCardsCache(c context.Context, ids []int64) (vs map[int64]*model.TicketCard, err error) {
	var cards map[string]*model.Cards
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, model.CardPrefixTicket+strconv.FormatInt(id, 10))
	}
	if cards, err = d.cardsCache(c, idsStr); err != nil {
		return
	}
	vs = make(map[int64]*model.TicketCard)
	for _, card := range cards {
		if (card != nil) && (card.TicketCard != nil) {
			vs[card.TicketCard.ID] = card.TicketCard
		}
	}
	return
}

// CacheHotspots .
func (d *Dao) CacheHotspots(c context.Context) (res []*model.Hotspot, err error) {
	res, err = d.cacheHotspots(c)
	for _, r := range res {
		if r.TopArticles == nil {
			r.TopArticles = []int64{}
		}
	}
	return
}
