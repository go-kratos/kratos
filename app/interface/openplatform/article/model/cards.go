package model

import (
	"fmt"
)

// Cards
const (
	CardPrefixBangumi   = "ss"
	CardPrefixBangumiEp = "ep"
	CardPrefixTicket    = "pw"
	CardPrefixMall      = "sp"
	CardPrefixAudio     = "au"
	CardPrefixArchive   = "av"
	CardPrefixArticle   = "cv"
)

// TicketCard .
type TicketCard struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"performance_image"`
	StartTime int64   `json:"start_time"`
	EndTime   int64   `json:"end_time"`
	Province  string  `json:"province_name"`
	City      string  `json:"city_name"`
	District  string  `json:"district_name"`
	Venue     string  `json:"venue_name"`
	PriceLow  float64 `json:"price_low"`
	URL       string  `json:"url"`
}

// MallCard .
type MallCard struct {
	ID     int64    `json:"itemsId"`
	Name   string   `json:"name"`
	Brief  string   `json:"brief"`
	Images []string `json:"img"`
	Price  int64    `json:"price"`
	Type   int      `json:"type"`
}

// AudioCard .
type AudioCard struct {
	ID       int64  `json:"song_id"`
	Title    string `json:"title"`
	UpMid    int64  `json:"up_mid"`
	UpName   string `json:"up_name"`
	Play     int64  `json:"play_num"`
	Reply    int64  `json:"reply_num"`
	CoverURL string `json:"cover_url"`
}

// BangumiCard .
type BangumiCard struct {
	ID       int64  `json:"season_id"`
	Type     int64  `json:"season_type"`
	TypeName string `json:"season_type_name"`
	Image    string `json:"cover"`
	Title    string `json:"title"`
	Rating   struct {
		Score float64 `json:"score"`
		Count int64   `json:"count"`
	} `json:"rating"`
	Playable    bool  `json:"playable"`
	FollowCount int64 `json:"follow_count"`
	PlayCount   int64 `json:"play_count"`
}

// Cards .
type Cards struct {
	Type        string       `json:"type,omitempty"`
	TicketCard  *TicketCard  `json:"ticket_card,omitempty"`
	BangumiCard *BangumiCard `json:"bangumi_card,omitempty"`
	MallCard    *MallCard    `json:"mall_card,omitempty"`
	AudioCard   *AudioCard   `json:"audio_card,omitempty"`
}

// Key .
func (c *Cards) Key() string {
	var id int64
	if c.TicketCard != nil {
		id = c.TicketCard.ID
	} else if c.BangumiCard != nil {
		id = c.BangumiCard.ID
	} else if c.MallCard != nil {
		id = c.MallCard.ID
	} else if c.AudioCard != nil {
		id = c.AudioCard.ID
	}
	return fmt.Sprintf("%s%d", c.Type, id)
}

// Item .
func (c *Cards) Item() interface{} {
	if c.TicketCard != nil {
		return c.TicketCard
	} else if c.BangumiCard != nil {
		return c.BangumiCard
	} else if c.MallCard != nil {
		return c.MallCard
	}
	return c.AudioCard.ID
}
