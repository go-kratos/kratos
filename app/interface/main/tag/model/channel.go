package model

import (
	"go-common/library/time"
)

// const const value.
const (
	ChanStateOffline   = int32(0) //频道下线
	ChanStateStop      = int32(1) //频道停用
	ChanStateCommon    = int32(2) //频道普通
	ChanStateRecommend = int32(3) //频道推荐

	DiscoveryChannelNum = int32(3) //3个发现频道
	ResMaxNum           = int(50)  // 最多可支持一次50个查询稿件

	DefaultPageNum  = int(1)
	DefaultPageSize = int(20)

	BatchSize = int32(100)

	TagChannelNo  = int32(0) // tag是运营频道
	TagChannelYes = int32(1) // tag是运营频道

	ChannelCategoryStateOK  = int32(0)
	ChannelCategoryStateDel = int32(1)
	ChannelRuleStateOK      = int32(0)
	ChannelRuleStateDel     = int32(1)

	ChannelRuleFlagSingle = int32(1)
	ChannelRuleFlagPlus   = int32(2)
	ChannelRuleFlagMinus  = int32(3)

	AIRecommandChannel = int32(9)
	AIRecommandTag     = int32(54)
	NoneUserID         = int64(0)

	ChannelFromApp = int32(0) // 移动端频道详情页
	ChannelFromH5  = int32(1) // 频道H5页面

	StateChannelCheckOK   = int32(0)
	StateChannelChecking  = int32(1)
	StateChannelCheckNone = int32(2)

	StateResCheckBackNo  = int32(0)
	StateResCheckBackYes = int32(1)

	ChannelActivity      = int32(4)
	StateChannelActivity = int32(1)

	ManagerYes = int32(1)
	ManagerNo  = int32(0)

	ChannelFromINT               = int32(1) //频道来源于：国际版
	ChanStateINTField            = int32(1) //频道国际版本屏蔽
	ChannelCategoryStateINTField = int32(1) //频道分类国际版本屏蔽

	ChannelStateTop       = int32(1) //频道置顶位
	ChannelStateShieldINT = int32(1) // 频道海外屏蔽

	ChannelAttrCheckBack = uint(0) //回查频道位
	ChannelAttrTop       = uint(1) //置顶频道位
	ChannelAttrActivity  = uint(2) //活动频道位
	ChannelAttrINT       = uint(3) //频道海外版位

)

// ChannelCategory channel category.
type ChannelCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	State     int32     `json:"state"`
	Order     int32     `json:"-"`
	INTShield int32     `json:"int_shield"` // International Shield. 国际版是否屏蔽
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"-"`
}

// ChannelRule channel rule.
type ChannelRule struct {
	Tid      int64  `json:"tid"`
	TidA     int64  `json:"a_tid"`
	TidAName string `json:"a_tname"`
	TidB     int64  `json:"b_tid"`
	TidBName string `json:"b_tname"`
	Flag     int32  `json:"flag"`
}

// ChannelRuleClassifier channel rule classifier.
type ChannelRuleClassifier struct {
	Single []*ChannelRule `json:"single"`
	Plus   []*ChannelRule `json:"plus"`
	Minus  []*ChannelRule `json:"minus"`
}

// HitRule rule hit.
// TODO 更好的计算方式.
func (rule *ChannelRule) HitRule(tids []int64) bool {
	switch rule.Flag {
	case ChannelRuleFlagSingle:
		for _, tid := range tids {
			if rule.TidA == tid {
				return true
			}
		}
	case ChannelRuleFlagPlus:
		var times int
		for _, tid := range tids {
			if rule.TidA == tid || rule.TidB == tid {
				times++
			}
		}
		return times == 2
	case ChannelRuleFlagMinus:
		var b bool
		times := len(tids)
		for _, tid := range tids {
			if rule.TidA == tid {
				b = true
			}
			if rule.TidA != tid {
				times = times - 1
			}
		}
		if times == 0 && b {
			return true
		}
	default:
	}
	return false
}

// RuleCALC rule calc.
// TODO 更好的计算方式.
func (rc *ChannelRuleClassifier) RuleCALC(tids []int64, mng int32) (res []*ChannelRule) {
	for _, single := range rc.Single {
		if single.HitRule(tids) {
			res = append(res, single)
			if mng == ManagerNo {
				return
			}
		}
	}
	for _, plus := range rc.Plus {
		if plus.HitRule(tids) {
			res = append(res, plus)
			if mng == ManagerNo {
				return
			}
		}
	}
	for _, minus := range rc.Minus {
		if minus.HitRule(tids) {
			res = append(res, minus)
			if mng == ManagerNo {
				return
			}
		}
	}
	return
}

// AttrVal get attr flag.
func (t *Channel) AttrVal(bit uint) int32 {
	return (t.Attr >> bit) & int32(1)
}

// Top channel is top.
func (t *Channel) Top() bool {
	return t.AttrVal(ChannelAttrTop) == ChannelStateTop
}

// Attent IsAttent is attent.
func (t *Channel) Attent() bool {
	return t.Attention == 1
}

// Recommend IsRecommend is recommend.
func (t *Channel) Recommend() bool {
	return t.State == ChanStateRecommend
}

// ChannelCategorySort .
type ChannelCategorySort []*ChannelCategory

// Len Len.
func (t ChannelCategorySort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelCategorySort) Less(i, j int) bool {
	if t[i].Order < t[j].Order {
		return true
	} else if t[i].Order == t[j].Order {
		if t[i].CTime > t[j].CTime {
			return true
		} else if t[i].CTime == t[j].CTime {
			if t[i].ID < t[j].ID {
				return true
			}
		}
	}
	return false
}

// Swap Swap.
func (t ChannelCategorySort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// BinocularRule BinocularRule.
type BinocularRule struct {
	A int64
	B int64
}

// ChanRule ChanRule.
type ChanRule struct {
	Single      []int64          // a
	AllContains []*BinocularRule // a+b
	Contains    []*BinocularRule // a-b
}

// ResourceChannel resource channels.
type ResourceChannel struct {
	Channel []*Channel `json:"channel"`
	Tags    []*Channel `json:"tag"`
}

// ResChannelCheckBack ResChannelCheckBack.
type ResChannelCheckBack struct {
	Tids      []int64                `json:"tids"`
	Channels  map[int64]*ChannelInfo `json:"channels"`
	CheckBack int32                  `json:"check_back"`
}

// ChannelInfo channel basic info.
type ChannelInfo struct {
	Tid       int64    `json:"tid"`
	TName     string   `json:"tname"`
	HitRule   string   `json:"hit_rule"`
	HitTNames []string `json:"hit_tnames"`
	HitRules  []string `json:"hit_rules"`
}

// ChannelResource ChannelResource.
type ChannelResource struct {
	Oids      []int64 `json:"resource"`
	Failover  bool    `json:"failover"`
	IsChannel bool    `json:"is_channel"`
	Pages     *Page   `json:"page"`
}

// AIChannelRecommand AIChannelRecommand.
type AIChannelRecommand struct {
	Tid  int64  `json:"tid"`
	Oid  int64  `json:"id"`
	Type string `json:"goto"`
}

// CustomSortChannel CustomSortChannel.
type CustomSortChannel struct {
	Custom   []*TagInfo `json:"custom"`
	Standard []*TagInfo `json:"standard"`
	Total    int        `json:"total"`
}

// ChannelSquare channel square.
type ChannelSquare struct {
	Channels []*Channel        `json:"tags"`
	Oids     map[int64][]int64 `json:"oids"`
}

// ChannelDetail channel detail.
type ChannelDetail struct {
	Tag     *TagInfo          `json:"tag"`
	Synonym []*ChannelSynonym `json:"synonyms"`
}

// Page page.
type Page struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pagesize"`
	Total    int64 `json:"count"`
}

// ChannelSort ChannelSort.
type ChannelSort []*Channel

// Len Len.
func (t ChannelSort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelSort) Less(i, j int) bool {
	if t[i].Rank < t[j].Rank {
		return true
	} else if t[i].Rank == t[j].Rank {
		if t[i].CTime < t[j].CTime {
			return true
		} else if t[i].CTime == t[j].CTime {
			if t[i].ID < t[j].ID {
				return true
			}
		}
	}
	return false
}

// Swap Swap.
func (t ChannelSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// ChannelRecomendSort ChannelRecomendSort.
type ChannelRecomendSort []*Channel

// Len Len.
func (t ChannelRecomendSort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelRecomendSort) Less(i, j int) bool {
	if t[i].TopRank < t[j].TopRank {
		return true
	} else if t[i].TopRank == t[j].TopRank {
		if t[i].CTime < t[j].CTime {
			return true
		} else if t[i].CTime == t[j].CTime {
			if t[i].ID < t[j].ID {
				return true
			}
		}
	}
	return false
}

// Swap Swap.
func (t ChannelRecomendSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// ChannelSynonymSort channel synonym sort by rank.
type ChannelSynonymSort []*ChannelSynonym

// Len Len.
func (t ChannelSynonymSort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelSynonymSort) Less(i, j int) bool {
	if t[i].Rank < t[j].Rank {
		return true
	} else if t[i].Rank == t[j].Rank {
		if t[i].Id < t[j].Id {
			return true
		} else if t[i].Id == t[j].Id {
			if t[i].CTime < t[j].CTime {
				return true
			}
		}
	}
	return false
}

// Swap Swap.
func (t ChannelSynonymSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
