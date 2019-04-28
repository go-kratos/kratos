package v2

// cacheTag tag.
type cacheTag struct {
	Tag     int64
	ConfIDs []int64
	Force   int8
}

// curTag current tag version.
type curTag struct {
	O *cacheTag
	C *cacheTag
}

func (tag *curTag) diff() (diffs []int64) {
	if tag.O == nil || tag.C == nil {
		return nil
	}
	oIDs := tag.O.ConfIDs
	tmp := make(map[int64]struct{}, len(oIDs))
	for _, oID := range tag.O.ConfIDs {
		tmp[oID] = struct{}{}
	}
	for _, ID := range tag.C.ConfIDs {
		if _, ok := tmp[ID]; !ok {
			diffs = append(diffs, ID)
		}
	}
	return
}

func (tag *curTag) old() int64 {
	if tag.O == nil {
		return 0
	}
	return tag.O.Tag
}

func (tag *curTag) cur() int64 {
	if tag.C == nil {
		return 0
	}
	return tag.C.Tag
}
