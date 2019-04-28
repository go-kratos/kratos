package model

//缓存key常量
const (
	CacheKeyUserShareToken    = "user:share:token:%d" //用户分享token
	CacheExpireUserShareToken = 600

	CacheKeyLastPubtime    = "%d:last:pubtime"
	CacheExpireLastPubtime = 432000
)
