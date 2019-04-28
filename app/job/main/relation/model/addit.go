package model

// AchieveFlag is
type AchieveFlag uint64

// const
var (
	EmptyAchieve         = AchieveFlag(0)
	FollowerAchieve1k    = AchieveFlag(1 << 0)
	FollowerAchieve5k    = AchieveFlag(1 << 1)
	FollowerAchieve10k   = AchieveFlag(1 << 2)
	FollowerAchieve100k  = AchieveFlag(1 << 3)
	FollowerAchieve1000k = AchieveFlag(1 << 12)
)

// AchieveFromFollower is
func AchieveFromFollower(count int64) AchieveFlag {
	if count <= 0 {
		return EmptyAchieve
	}
	if count >= 100000 {
		return AchieveFlag(1 << uint64(2+count/100000))
	}
	if count >= 10000 && count < 100000 {
		return FollowerAchieve10k
	}
	if count >= 5000 && count < 10000 {
		return FollowerAchieve5k
	}
	if count >= 1000 && count < 5000 {
		return FollowerAchieve1k
	}
	return EmptyAchieve
}

// AllAchieveFromFollower is
func AllAchieveFromFollower(count int64) []AchieveFlag {
	flags := []AchieveFlag{}
	if count <= 0 {
		return flags
	}
	if count >= 1000 {
		flags = append(flags, FollowerAchieve1k)
	}
	if count >= 5000 {
		flags = append(flags, FollowerAchieve5k)
	}
	if count >= 10000 {
		flags = append(flags, FollowerAchieve10k)
	}
	if count >= 100000 {
		remain := count / 100000
		for i := int64(1); i <= remain; i++ {
			flags = append(flags, AchieveFlag(1<<uint64(2+i)))
		}
	}
	return flags
}
