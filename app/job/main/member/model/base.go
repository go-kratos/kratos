package model

import (
	"fmt"
	"math/rand"
)

const (
	//CacheKeyBase is.
	CacheKeyBase = "bs_%d" // key of baseInfo
	//CacheKeyMoral is.
	CacheKeyMoral = "moral_%d" // key of detail
	//CacheKeyInfo is.
	CacheKeyInfo = "i_"
	//URLNoFace is.
	URLNoFace = "http://static.hdslb.com/images/member/noface.gif"
	//TableExpLog is.
	TableExpLog = "ugc:ExpLog"
	//TableMoralLog is.
	TableMoralLog = "ugc:MoralLog"
)

// RandFaceURL get face URL
func (b *BaseInfo) RandFaceURL() {
	if b.Face == "" {
		b.Face = URLNoFace
		return
	}
	b.Face = fmt.Sprintf("http://i%d.hdslb.com%s", rand.Int63n(3), b.Face)
}
