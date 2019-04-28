package service

import (
	"context"
	"fmt"

	artmdl "go-common/app/interface/openplatform/article/model"
)

var (
	_likeMessage = int64(1)
)

// SendMessage send message to uppper
func (s *Service) SendMessage(c context.Context, aid int64, stat *artmdl.Stats) (err error) {
	var (
		title, msg string
		meta       *artmdl.Meta
		max        int64
	)
	if exist, _ := s.dao.ExpireMaxLikeCache(c, aid); exist {
		max, _ = s.dao.MaxLikeCache(c, aid)
	}
	if (stat.Like <= max) || (!shouldNofify(stat.Like)) {
		return
	}
	if meta, err = s.ArticleMeta(c, aid); (err != nil) || (meta == nil) {
		return
	}
	mid := meta.Author.Mid
	if len(s.c.Article.MessageMids) > 0 {
		var exist bool
		for _, m := range s.c.Article.MessageMids {
			if m == mid {
				exist = true
				break
			}
		}
		if !exist {
			return
		}
	}
	title = fmt.Sprintf("有%v人点赞了你的专栏文章", stat.Like)
	msg = fmt.Sprintf("有%v个小伙伴点赞你投稿的专栏文章“#{%s}{\"http://www.bilibili.com/read/cv%d\"}”～快去看看吧！#{点击前往}{\"http://www.bilibili.com/read/cv%d\"}", stat.Like, meta.Title, aid, aid)
	err = s.dao.SendMessage(c, _likeMessage, mid, aid, title, msg)
	cache.Save(func() {
		s.dao.SetMaxLikeCache(context.TODO(), aid, stat.Like)
	})
	return
}

func shouldNofify(n int64) (res bool) {
	switch {
	case n <= 0:
		res = false
	case n <= 10:
		res = true
	case n <= 100:
		res = (n%10 == 0)
	case n <= 1000:
		res = (n%100 == 0)
	default:
		res = (n%10000 == 0)
	}
	return
}
