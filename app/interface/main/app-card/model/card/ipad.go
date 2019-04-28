package card

import (
	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/operate"
	tag "go-common/app/interface/main/tag/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/log"
)

func ipadHandle(cardGoto model.CardGt, cardType model.CardType, rcmd *ai.Item, tagm map[int64]*tag.Tag, isAttenm map[int64]int8, statm map[int64]*relation.Stat, cardm map[int64]*account.Card) (hander Handler) {
	base := &Base{CardGoto: cardGoto, Rcmd: rcmd, Tagm: tagm, IsAttenm: isAttenm, Statm: statm, Cardm: cardm, Columnm: model.ColumnSvrSingle}
	switch cardType {
	default:
		switch cardGoto {
		case model.CardGotoAv, model.CardGotoBangumi, model.CardGotoLive, model.CardGotoPGC:
			base.CardType = model.LargeCoverV1
			base.CardLen = 1
			hander = &LargeCoverV1{Base: base}
		case model.CardGotoBangumiRcmd:
			base.CardType = model.SmallCoverV1
			hander = &SmallCoverV1{Base: base}
		case model.CardGotoRank:
			base.CardType = model.FourItemHV3
			hander = &FourItemV3{Base: base}
		case model.CardGotoLogin:
			base.CardType = model.CoverOnlyV3
			base.CardLen = 1
			hander = &CoverOnly{Base: base}
		case model.CardGotoBanner:
			base.CardType = model.BannerV3
			hander = &Banner{Base: base}
		case model.CardGotoAdAv:
			base.CardType = model.CmV1
			base.CardLen = 1
			hander = &LargeCoverV1{Base: base}
		case model.CardGotoAdWebS:
			base.CardType = model.CmV1
			base.CardLen = 1
			hander = &SmallCoverV1{Base: base}
		case model.CardGotoAdWeb:
			base.CardType = model.CmV1
			base.CardLen = 2
			hander = &SmallCoverV1{Base: base}
		case model.CardGotoSearchUpper:
			base.CardType = model.ThreeItemAll
			hander = &ThreeItemAll{Base: base}
		}
	}
	return
}

type FourItemV3 struct {
	*Base
	TitleIcon   model.Icon        `json:"title_icon,omitempty"`
	BannerCover string            `json:"banner_cover,omitempty"`
	BannerURI   string            `json:"banner_uri,omitempty"`
	MoreURI     string            `json:"more_uri,omitempty"`
	MoreText    string            `json:"more_text,omitempty"`
	Items       []*FourItemV3Item `json:"items,omitempty"`
}

type FourItemV3Item struct {
	Base
	CoverLeftText string     `json:"cover_left_text,omitempty"`
	CoverLeftIcon model.Icon `json:"cover_left_icon,omitempty"`
	Desc1         string     `json:"desc_1,omitempty"`
	Desc2         string     `json:"desc_2,omitempty"`
	Badge         string     `json:"badge,omitempty"`
}

func (c *FourItemV3) From(main interface{}, op *operate.Card) {
	switch main.(type) {
	case map[model.Gt]interface{}:
		intfcm := main.(map[model.Gt]interface{})
		if op == nil {
			return
		}
		switch op.CardGoto {
		case model.CardGotoRank:
			const (
				_title = "全站排行榜"
				_limit = 4
			)
			c.Base.from("0", "", _title, "", "", nil)
			// c.TitleIcon = model.IconRank
			c.MoreURI = model.FillURI(op.Goto, op.URI, nil)
			c.MoreText = "查看更多"
			c.Items = make([]*FourItemV3Item, 0, _limit)
			for _, v := range op.Items {
				if v == nil {
					continue
				}
				intfc, ok := intfcm[v.Goto]
				if !ok {
					continue
				}
				am := intfc.(map[int64]*archive.ArchiveWithPlayer)
				a, ok := am[v.ID]
				if !ok || !model.AvIsNormal(a) {
					continue
				}
				item := &FourItemV3Item{
					CoverLeftText: model.DurationString(a.Duration),
					Desc1:         model.ScoreString(v.Score),
				}
				item.Base.from(v.Param, a.Pic, a.Title, model.GotoAv, v.URI, model.AvPlayHandler(a.Archive3, a.PlayerInfo, op.TrackID))
				item.Args.fromArchive(a.Archive3, nil)
				c.Items = append(c.Items, item)
				if len(c.Items) == _limit {
					break
				}
			}
			if len(c.Items) < _limit {
				return
			}
			c.Items[0].CoverLeftIcon = model.IconGoldMedal
			c.Items[1].CoverLeftIcon = model.IconSilverMedal
			c.Items[2].CoverLeftIcon = model.IconBronzeMedal
		default:
			log.Warn("FourItemV3 From: unexpected card_goto %s", op.CardGoto)
			return
		}
	default:
		log.Warn("FourItemV3 From: unexpected type %T", main)
		return
	}
	c.Right = true
}

func (c *FourItemV3) Get() *Base {
	return c.Base
}
