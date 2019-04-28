package model

// attribute bit. priority black > following > whisper > no relation.
const (
	AttrNoRelation = uint32(0)
	AttrWhisper    = uint32(1)
	AttrFollowing  = uint32(1) << 1
	AttrFriend     = uint32(1) << 2
	AttrBlack      = uint32(1) << 7
	// 128，129,130 变为 0 时候，status = 1
	StatusOK  = 0
	StatusDel = 1
)

// relation act type.
const (
	ActAddFollowing = int8(1)
	ActDelFollowing = int8(2)
	ActAddWhisper   = int8(3)
	ActDelWhisper   = int8(4)
	ActAddBlack     = int8(5)
	ActDelBalck     = int8(6)
	ActDelFollower  = int8(7)
)

// Attr get real attribute by the specified priority.
func Attr(attribute uint32) uint32 {
	if attribute&AttrBlack > 0 {
		return AttrBlack
	}
	if attribute&AttrFriend > 0 {
		return AttrFriend
	}
	if attribute&AttrFollowing > 0 {
		return AttrFollowing
	}
	if attribute&AttrWhisper > 0 {
		return AttrWhisper
	}
	return AttrNoRelation
}

// SetAttr set attribute.
func SetAttr(attribute uint32, mask uint32) uint32 {
	return attribute | mask
}

// UnsetAttr unset attribute.
func UnsetAttr(attribute uint32, mask uint32) uint32 {
	return attribute & ^mask // ^ 按位取反
}
