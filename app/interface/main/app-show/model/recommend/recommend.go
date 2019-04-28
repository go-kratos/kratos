package recommend

import (
	"go-common/app/interface/main/app-show/model/card"
)

// Arc is index show recommend.
type Arc struct {
	Aid         interface{} `json:"aid"`
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Pic         string      `json:"pic"`
	Views       interface{} `json:"play"`
	Comments    int64       `json:"review"`
	Coins       int64       `json:"coins"`
	Danmaku     int         `json:"video_review"`
	Favorites   int64       `json:"favorites"`
	Pts         int64       `json:"pts"`
	Others      []*Arc      `json:"others"`
}

type List struct {
	Aid        int64  `json:"aid"`
	Desc       string `json:"desc"`
	CornerMark int8   `json:"corner_mark"`
}

type CardList struct {
	ID         int64            `json:"id"`
	Goto       string           `json:"goto"`
	FromType   string           `json:"from_type"`
	Desc       string           `json:"desc"`
	CornerMark int8             `json:"corner_mark"`
	Condition  []*CardCondition `json:"condition"`
}

type CardCondition struct {
	Plat      int8   `json:"plat"`
	Condition string `json:"conditions"`
	Build     int    `json:"build"`
}

func (c *CardList) CardListChange() (p *card.PopularCard) {
	p = &card.PopularCard{
		Value:      c.ID,
		Type:       c.Goto,
		FromType:   c.FromType,
		Reason:     c.Desc,
		CornerMark: c.CornerMark,
	}
	if p.Reason != "" {
		p.ReasonType = 3
	}
	if len(c.Condition) > 0 {
		tmpcondition := map[int8][]*card.PopularCardPlat{}
		for _, condition := range c.Condition {
			tmpcondition[condition.Plat] = append(tmpcondition[condition.Plat], &card.PopularCardPlat{
				Plat:      condition.Plat,
				Condition: condition.Condition,
				Build:     condition.Build,
			})
		}
		p.PopularCardPlat = tmpcondition
	}
	return
}
