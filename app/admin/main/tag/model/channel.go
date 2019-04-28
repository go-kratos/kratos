package model

import "go-common/library/time"

// const const value.
const (
	ChanStateOffline  = int32(0) //频道下线
	ChanStateStop     = int32(1) //频道停用
	ChanStateCommon   = int32(2) //频道普通
	ChanStateRecomend = int32(3) //频道推荐

	// ChannelStateTop .
	ChannelStateTop        = int32(1) //频道置顶位
	ChannelStateShieldINT  = int32(1) // 频道海外屏蔽
	CategoryStateShieldINT = int32(1) // 频道分类海外屏蔽

	ChanRuleNormal = int32(0) //频道规则状态：正常
	ChanRuleDelete = int32(1) //频道规则状态：删除

	ChanGroupNormal = int32(0) //频道规则状态：正常
	ChanGroupDelete = int32(1) //频道规则状态：删除

	ChannelRuleMaxLen = int(2) //频道规则最长组成元素个数

	ChannelActivity   = int32(4) //是否为活动频道
	ChannelSynonymMax = int32(8) //频道最多相似频道

	ChannelTopNo  = int32(0) //非置顶频道
	ChannelTopYes = int32(1) //是置顶频道

	ChannelAttrCheckBack = uint(0) //回查频道位
	ChannelAttrTop       = uint(1) //置顶频道位
	ChannelAttrActivity  = uint(2) //活动频道位
	ChannelAttrINT       = uint(3) //频道海外版位

	ChannelCategoryAttrINT = uint(0) //频道分类海外版位

	ShortContentMax  = int(15)
	DetailContentMax = int(56)
)

// Channel Channel struct.
type Channel struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Type         int64     `json:"type"`
	Rank         int32     `json:"rank"`
	Operator     string    `json:"operator"`
	Editor       string    `json:"editor"`
	Cover        string    `json:"cover"`
	HeadCover    string    `json:"head_cover"`
	Content      string    `json:"detail_content"`
	ShortContent string    `json:"short_content"`
	Attr         int32     `json:"attr"`
	State        int32     `json:"state"`
	Top          int32     `json:"top"`
	TopRank      int32     `json:"-"`
	CheckBack    int32     `json:"check_back"`
	Activity     int32     `json:"activity"`
	INTShield    int32     `json:"int_shield"` // International Shield. 国际版是否屏蔽
	Count        *TagCount `json:"count,omitempty"`
	TagNums      int32     `json:"tag_nums"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"mtime"`
}

// AttrVal get attr flag.
func (t *Channel) AttrVal(bit uint) int32 {
	return (t.Attr >> bit) & int32(1)
}

// AttrSet channel attr set.
func (t *Channel) AttrSet(bit uint, v int32) {
	t.Attr = t.Attr&(^(1 << bit)) | (v << bit)
}

// ChannelInfo ChannelInfo.
type ChannelInfo struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	Type         int64             `json:"type"`
	Rank         int32             `json:"rank"`
	Operator     string            `json:"operator"`
	Cover        string            `json:"cover"`
	HeadCover    string            `json:"head_cover"`
	Content      string            `json:"detail_content"`
	ShortContent string            `json:"short_content"`
	Attr         int32             `json:"-"`
	CheckBack    int32             `json:"check_back"`
	Activity     int32             `json:"activity"`
	INTShield    int32             `json:"int_shield"` // International Shield. 国际版是否屏蔽
	State        int32             `json:"state"`
	Rules        []*ChannelRule    `json:"rules"`
	Synonyms     []*ChannelSynonym `json:"synonyms"`
}

// ChannelCategory ChannelCategory.
type ChannelCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Order     int32     `json:"order"`
	State     int32     `json:"state"`
	Attr      int32     `json:"-"`
	INTShield int32     `json:"int_shield"` // International Shield. 国际版是否屏蔽
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// AttrVal get attr flag.
func (t *ChannelCategory) AttrVal(bit uint) int32 {
	return (t.Attr >> bit) & int32(1)
}

// AttrSet get attr flag.
func (t *ChannelCategory) AttrSet(bit uint, v int32) {
	t.Attr = t.Attr&(^(1 << bit)) | (v << bit)
}

// ChannelRule ChannelRule.
type ChannelRule struct {
	ID            int64     `json:"id"`
	Tid           int64     `json:"tid"`
	InRule        string    `json:"in"`
	NotInRule     string    `json:"notin"`
	InRuleName    string    `json:"in_name"`
	NotInRuleName string    `json:"notin_name"`
	Name          string    `json:"name"`
	Editor        string    `json:"editor"`
	State         int32     `json:"state"`
	CTime         time.Time `json:"ctime"`
	MTime         time.Time `json:"mtime"`
}

// ChannelSynonym channel group.
type ChannelSynonym struct {
	ID       int64     `json:"id"`
	PTid     int64     `json:"-"`
	Tid      int64     `json:"tid"`
	TName    string    `json:"tname"`
	Alias    string    `json:"alias"`
	Rank     int32     `json:"rank"`
	Operator string    `json:"operator"`
	State    int32     `json:"state"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// ChannelSynonymSort ChannelSynonymSort.
type ChannelSynonymSort []*ChannelSynonym

// Len Len.
func (t ChannelSynonymSort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelSynonymSort) Less(i, j int) bool {
	if t[i].State < t[j].State {
		return true
	} else if t[i].State == t[j].State {
		if t[i].Rank < t[j].Rank {
			return true
		} else if t[i].Rank == t[j].Rank {
			if t[i].ID < t[j].ID {
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

// ChannelSort channel sort.
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

// RecommendChannelSort channel sort.
type RecommendChannelSort []*Channel

// Len Len.
func (t RecommendChannelSort) Len() int {
	return len(t)
}

// Less Less.
func (t RecommendChannelSort) Less(i, j int) bool {
	if t[i].TopRank < t[j].TopRank {
		return true
	} else if t[i].TopRank == t[j].TopRank {
		if t[i].MTime < t[j].MTime {
			return true
		} else if t[i].MTime == t[j].MTime {
			if t[i].ID < t[j].ID {
				return true
			}
		}
	}
	return false
}

// Swap Swap.
func (t RecommendChannelSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
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
			if t[i].ID > t[j].ID {
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

// ChannelRuleSort .
type ChannelRuleSort []*ChannelRule

// Len Len.
func (t ChannelRuleSort) Len() int {
	return len(t)
}

// Less Less.
func (t ChannelRuleSort) Less(i, j int) bool {
	if t[i].State < t[j].State {
		return true
	} else if t[i].State == t[j].State {
		if t[i].ID < t[j].ID {
			return true
		}
	}
	return false
}

// Swap Swap.
func (t ChannelRuleSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
