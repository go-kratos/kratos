package model

import (
	smodel "go-common/app/service/main/relation/model"
	"go-common/library/time"
	"sort"
)

// Relation is
type Relation struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Mid       int64     `json:"mid" gorm:"column:mid"`
	Fid       int64     `json:"fid" gorm:"column:fid"`
	Attribute uint32    `json:"attribute" gorm:"column:attribute"`
	Status    int8      `json:"status" gorm:"column:status"`
	Source    int8      `json:"source" gorm:"column:source"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`

	Relation uint32 `json:"relation"`
}

// Stat is
type Stat struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Mid       int64     `json:"mid" gorm:"column:mid"`
	Following int64     `json:"following" gorm:"column:following"`
	Whisper   int64     `json:"whisper" gorm:"column:whisper"`
	Black     int64     `json:"black" gorm:"column:black"`
	Follower  int64     `json:"follower" gorm:"column:follower"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
}

// ParseRelation is
func (r *Relation) ParseRelation() {
	r.Relation = smodel.Attr(r.Attribute)
}

// Follower is
type Follower struct {
	*Relation
	MemberName   string `json:"member_name"`
	FollowerName string `json:"follower_name"`
}

// Following is
type Following struct {
	*Relation
	MemberName    string `json:"member_name"`
	FollowingName string `json:"following_name"`
}

// RelationList is
type RelationList []*Relation

// FollowersList is
type FollowersList []*Follower

// FollowingsList is
type FollowingsList []*Following

func (rl RelationList) Len() int {
	return len(rl)
}

func (rl RelationList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

func (rl RelationList) Less(i, j int) bool {
	return rl[i].MTime < rl[j].MTime
}

// Paginate is
func (rl RelationList) Paginate(skip int, size int) RelationList {
	if skip > len(rl) {
		skip = len(rl)
	}

	end := skip + size
	if end > len(rl) {
		end = len(rl)
	}

	return rl[skip:end]
}

// FilterMTimeFrom is
func (rl RelationList) FilterMTimeFrom(from time.Time) RelationList {
	res := make(RelationList, 0)
	for _, r := range rl {
		if r.MTime >= from {
			res = append(res, r)
		}
	}
	return res
}

// FilterMTimeTo is
func (rl RelationList) FilterMTimeTo(to time.Time) RelationList {
	res := make(RelationList, 0)
	for _, r := range rl {
		if r.MTime <= to {
			res = append(res, r)
		}
	}
	return res
}

// OrderByMTime is
func (rl RelationList) OrderByMTime(desc bool) {
	sort.Sort(rl)
}

// FollowersList is
func (rl RelationList) FollowersList() FollowersList {
	res := make(FollowersList, 0, len(rl))
	for _, r := range rl {
		res = append(res, &Follower{
			Relation: r,
		})
	}
	return res
}

// FollowingsList is
func (rl RelationList) FollowingsList() FollowingsList {
	res := make(FollowingsList, 0, len(rl))
	for _, r := range rl {
		res = append(res, &Following{
			Relation: r,
		})
	}
	return res
}
