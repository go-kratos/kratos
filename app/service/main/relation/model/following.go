package model

var (
	_emptyFollowings = make([]*Following, 0)
)

// Black get if black.
func (f *Following) Black() bool {
	return AttrBlack == Attr(f.Attribute)
}

// Friend get if both way following.
func (f *Following) Friend() bool {
	return AttrFriend == Attr(f.Attribute)
}

// Following get if following.
func (f *Following) Following() bool {
	return AttrFollowing == Attr(f.Attribute) || Attr(f.Attribute) == AttrFriend
}

// Whisper get if whisper.
func (f *Following) Whisper() bool {
	return AttrWhisper == Attr(f.Attribute)
}

// Filter filter followings by the given attribute.
func Filter(fs []*Following, attr uint32) (res []*Following) {
	for _, f := range fs {
		// NOTE: if current attribute evaluated by Attr() matched, then continue,
		// this includes the situation that matches black, friend, whisper, and no-relation directly.
		// Now we have following to deal with, since we know that the attribute friend
		// can either do not exist or exists with following at the same time,
		// to deal with this situation, we need to filter for items which have 1 on the bit that attr stands for,
		// and especially, the attribute it self cannot be black because the attribute black has the highest priority,
		// when it exists, it shadows other bits, including friend, following, whisper, no-relation,
		// there is no need to do further calculate,
		// more specifically, black when black included, the value of f.Attribute&attr may greater than 0
		// when f.Attribute is 128+2 or 128+1 and the corresponding attr is 2 or 1,
		// which is not as we expected.
		if f.Attribute == 4 {
			f.Attribute = 6
		}
		if (Attr(f.Attribute) == attr) || (!f.Black() && f.Attribute&attr > 0) {
			res = append(res, f)
		}
	}
	if len(res) == 0 {
		res = _emptyFollowings
	}
	return
}

// SortFollowings sort followings by the mtime desc.
type SortFollowings []*Following

func (fs SortFollowings) Len() int {
	return len(fs)
}
func (fs SortFollowings) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
func (fs SortFollowings) Less(i, j int) bool {
	if fs[i].MTime == fs[j].MTime {
		return fs[i].Mid < fs[j].Mid
	}
	return fs[i].MTime.Time().After(fs[j].MTime.Time())
}
