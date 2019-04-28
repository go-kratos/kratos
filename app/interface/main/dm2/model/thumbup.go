package model

// All const variable used in thumbup
const (
	ThumbupLike       int8 = 1 // 点赞
	ThumbupLikeCancel int8 = 2 // 取消赞
	ThumbupHate       int8 = 3 // 点踩
	ThumbupHateCancel int8 = 4 // 取消踩
)

// ThumbupStat thumbup state
type ThumbupStat struct {
	Likes    int64 `json:"likes"`
	UserLike int8  `json:"user_like"`
}

// CheckThumbup check thumbup.
func CheckThumbup(Thumbup int8) bool {
	if Thumbup == ThumbupLikeCancel || Thumbup == ThumbupLike || Thumbup == ThumbupHateCancel || Thumbup == ThumbupHate {
		return true
	}
	return false
}
