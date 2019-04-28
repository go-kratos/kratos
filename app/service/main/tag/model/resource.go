package model

const (
	// ResTypeArticle res type .
	ResTypeArticle = int32(1) // 文章类型 	// resource type
	// ResTypeMusic res type .
	ResTypeMusic = int32(2) // 音乐类型
	// ResTypeArchive res type .
	ResTypeArchive = int32(3) // 稿件类型
	// ResTypeOpenMall res type .
	ResTypeOpenMall = int32(4) // 开放平台电商

	// ResTagUser res tag user .
	ResTagUser = int32(0) // tag membership
	// ResTagUpper res tag uper .
	ResTagUpper = int32(1)

	// ResRoleUpper res tag uper .
	ResRoleUpper = int32(0) // up主  resource role
	// ResRoleUser res tag uper .
	ResRoleUser = int32(1) // 用户
	// ResRoleAdmin res tag uper .
	ResRoleAdmin = int32(2) // 管理员

	// ResStateNormal res tag uper .
	ResStateNormal = int32(0) // resource state
	// ResStateDelete res tag uper .
	ResStateDelete = int32(1)
	// ResStateDefault res tag default.
	ResStateDefault = int32(3)

	// ResAttrLocked resource attr of locked
	ResAttrLocked = uint(0)

	// UserActionNormal res tag uper .
	UserActionNormal = int32(0)
	// UserActionDel res tag uper .
	UserActionDel = int32(1)
	// UserActionAdd res tag uper .
	UserActionAdd = int32(2)
	// UserActionLike res tag uper .
	UserActionLike = int32(3)
	// UserActionHate res tag uper .
	UserActionHate = int32(4)

	// LimitAttrAllowAdd limit allow add
	LimitAttrAllowAdd = uint(0)
	// LimitAttrAllowDel limit allow delete
	LimitAttrAllowDel = uint(1)

	// LogActionAdd limit allow delete
	LogActionAdd = int32(0)
	// LogActionDel limit allow delete
	LogActionDel = int32(1)

	// ResTagLogAdd res tag uper .
	ResTagLogAdd = int32(0)
	// ResTagLogDel res tag uper .
	ResTagLogDel = int32(1)
	// ResTagLogOpen res tag uper .
	ResTagLogOpen = int32(0)
	// ResTagLogClose res tag uper .
	ResTagLogClose = int32(1)
)

// AttrVal gets attr val by bit.
func (r *Resource) AttrVal(bit uint) int32 {
	return (r.Attr >> bit) & int32(1)
}

// AttrSet sets attr value by bit.
func (r *Resource) AttrSet(v int32, bit uint) {
	r.Attr = r.Attr&(^(1 << bit)) | (v << bit)
}

// Locked resource locked state
func (r *Resource) Locked() bool {
	return r.AttrVal(ResAttrLocked) == AttrYes
}

// Resources .
type Resources []*Resource

func (t Resources) Len() int { return len(t) }

func (t Resources) Less(i, j int) bool {
	if t[i].Like > t[j].Like {
		return true
	} else if t[i].Like == t[j].Like {
		if t[i].Type > t[j].Type {
			return true
		} else if t[i].Type == t[j].Type {
			if (t[i].Role == 0 && t[j].Role != 0) || (t[i].Role == 2 && t[j].Role == 1) {
				return true
			}
			if t[i].Role == t[j].Role && t[i].Hate < t[j].Hate {
				return true
			}
			if t[i].Hate == t[j].Hate && t[i].CTime < t[j].CTime {
				return true
			}
		}
	}
	return false
}

func (t Resources) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
