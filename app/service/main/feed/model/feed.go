package model

import (
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/time"
)

// feed type
const (
	ArchiveType = 0
	BangumiType = 1

	TypeApp = iota
	TypeWeb
	TypeArt
)

// FeedType return feed type
func FeedType(app bool) int {
	if app {
		return TypeApp
	}
	return TypeWeb
}

// Feed struct of Feed
type Feed struct {
	Type int64 `json:"type"`
	// Data is *api.Arc or *Bangumi
	Archive *api.Arc `json:"archive"`
	Bangumi *Bangumi `json:"bangumi"`
	// ID is aid or SeasonID
	ID      int64      `json:"id"`
	PubDate time.Time  `json:"pubdate"`
	Fold    []*api.Arc `json:"fold"`
}

type Feeds []*Feed

func (as Feeds) Len() int { return len(as) }
func (as Feeds) Less(i, j int) bool {
	if as[i].PubDate != as[j].PubDate {
		return as[i].PubDate > as[j].PubDate
	}
	return as[i].ID > as[j].ID
}
func (as Feeds) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

type ArticleFeeds []*artmdl.Meta

func (as ArticleFeeds) Len() int { return len(as) }
func (as ArticleFeeds) Less(i, j int) bool {
	if as[i].PublishTime != as[j].PublishTime {
		return as[i].PublishTime > as[j].PublishTime
	}
	return as[i].ID > as[j].ID
}
func (as ArticleFeeds) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

// Arcs AidPubTime slice
type Arcs []*archive.AidPubTime

func (as Arcs) Len() int { return len(as) }
func (as Arcs) Less(i, j int) bool {
	if as[i].PubDate != as[j].PubDate {
		return as[i].PubDate > as[j].PubDate
	}
	return as[i].Aid > as[j].Aid
}
func (as Arcs) Swap(i, j int) { as[i], as[j] = as[j], as[i] }
