package model

// AchieveFlag is
type AchieveFlag uint64

// const
var (
	EmptyAchieve       = AchieveFlag(0)
	FollowerAchieve1k  = AchieveFlag(1 << 0)
	FollowerAchieve5k  = AchieveFlag(1 << 1)
	FollowerAchieve10k = AchieveFlag(1 << 2)
)
