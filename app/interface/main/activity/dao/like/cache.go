package like

import (
	"context"
	"fmt"

	likemdl "go-common/app/interface/main/activity/model/like"
)

// likeKey likes table line cache
func likeKey(id int64) string {
	return fmt.Sprintf("go_l_id_%d", id)
}

// actSubjectKey act_subject table line cache .
func actSubjectKey(id int64) string {
	return fmt.Sprintf("go_s_id_%d", id)
}

// actSubjectMaxIDKey act_subject table max id cache
func actSubjectMaxIDKey() string {
	return "go_sub_id_max"
}

// likeMaxIDKey likes table max id cache
func likeMaxIDKey() string {
	return "go_like_id_max"
}

// likeMissionBuffKey .
func likeMissionBuffKey(sid, mid int64) string {
	return fmt.Sprintf("go_l_m_a_%d_%d", sid, mid)
}

// likeMissionGroupIDkey .
func likeMissionGroupIDkey(lid int64) string {
	return fmt.Sprintf("go_l_m_g_id_%d", lid)
}

// likeActMissionKey flag has buff or not.
func likeActMissionKey(sid, lid, mid int64) string {
	return fmt.Sprintf("go:b-a:m:l:%d:%d:%d", sid, lid, mid)
}

// actAchieveKey .
func actAchieveKey(sid int64) string {
	return fmt.Sprintf("go:a:achs:%d", sid)
}

// actMissionFriendsKey .
func actMissionFriendsKey(sid, lid int64) string {
	return fmt.Sprintf("go:a:m:frd:%d:%d", sid, lid)
}

// actUserAchieveKey .
func actUserAchieveKey(id int64) string {
	return fmt.Sprintf("go:a:u:m:%d", id)
}

// actUserAchieveAwardKey .
func actUserAchieveAwardKey(id int64) string {
	return fmt.Sprintf("go:a:u:a:%d", id)
}

func subjectStatKey(sid int64) string {
	return fmt.Sprintf("ob_s_%d", sid)
}

func viewRankKey(sid int64) string {
	return fmt.Sprintf("v_r_%d", sid)
}

func likeContentKey(lid int64) string {
	return fmt.Sprintf("go_l_ct_%d", lid)
}

func sourceItemKey(sid int64) string {
	return fmt.Sprintf("so_i_%d", sid)
}

func subjectProtocolKey(sid int64) string {
	return fmt.Sprintf("go_s_pt_%d", sid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -sync=true
	Like(c context.Context, id int64) (*likemdl.Item, error)
	// cache: -sync=true
	Likes(c context.Context, ids []int64) (map[int64]*likemdl.Item, error)
	// cache: -sync=true
	ActSubject(c context.Context, id int64) (*likemdl.SubjectItem, error)
	//cache: -sync=true -nullcache=-1 -check_null_code=$==-1
	LikeMissionBuff(ctx context.Context, sid int64, mid int64) (res int64, err error)
	//cache: -sync=true
	MissionGroupItems(ctx context.Context, lids []int64) (map[int64]*likemdl.MissionGroup, error)
	//cache: -sync=true -nullcache=-1 -check_null_code=$!=nil&&$==-1
	ActMission(ctx context.Context, sid int64, lid int64, mid int64) (res int64, err error)
	//cache:-sync=true
	ActLikeAchieves(ctx context.Context, sid int64) (res *likemdl.Achievements, err error)
	//cache:-sync=true
	ActMissionFriends(ctx context.Context, sid int64, lid int64) (res *likemdl.ActMissionGroups, err error)
	//cache:-sync=true
	ActUserAchieve(ctx context.Context, id int64) (res *likemdl.ActLikeUserAchievement, err error)
	// cache
	MatchSubjects(c context.Context, ids []int64) (map[int64]*likemdl.Object, error)
	// cache:-sync=true
	LikeContent(c context.Context, ids []int64) (map[int64]*likemdl.LikeContent, error)
	// cache
	SourceItemData(c context.Context, sid int64) ([]int64, error)
	// cache:-sync=true
	ActSubjectProtocol(c context.Context, sid int64) (res *likemdl.ActSubjectProtocol, err error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=likeKey
	CacheLike(c context.Context, id int64) (*likemdl.Item, error)
	// mc: -key=likeKey
	CacheLikes(c context.Context, id []int64) (map[int64]*likemdl.Item, error)
	// mc: -key=likeKey -expire=d.mcPerpetualExpire -encode=json
	AddCacheLikes(c context.Context, items map[int64]*likemdl.Item) error
	// mc: -key=likeKey -expire=d.mcPerpetualExpire -encode=json
	AddCacheLike(c context.Context, key int64, value *likemdl.Item) error
	// mc: -key=actSubjectKey
	CacheActSubject(c context.Context, id int64) (*likemdl.SubjectItem, error)
	// mc: -key=actSubjectKey -expire=d.mcPerpetualExpire -encode=pb
	AddCacheActSubject(c context.Context, key int64, value *likemdl.SubjectItem) error
	// mc: -key=actSubjectMaxIDKey
	CacheActSubjectMaxID(c context.Context) (res int64, err error)
	// mc: -key=actSubjectMaxIDKey -expire=d.mcPerpetualExpire -encode=raw
	AddCacheActSubjectMaxID(c context.Context, sid int64) error
	// mc: -key=likeMaxIDKey
	CacheLikeMaxID(c context.Context) (res int64, err error)
	// mc: -key=likeMaxIDKey -expire=d.mcPerpetualExpire -encode=raw
	AddCacheLikeMaxID(c context.Context, lid int64) error
	//mc: -key=likeMissionBuffKey
	CacheLikeMissionBuff(c context.Context, sid int64, mid int64) (res int64, err error)
	//mc: -key=likeMissionBuffKey
	AddCacheLikeMissionBuff(c context.Context, sid int64, val int64, mid int64) error
	//mc: -key=likeMissionGroupIDkey
	CacheMissionGroupItems(ctx context.Context, lids []int64) (map[int64]*likemdl.MissionGroup, error)
	//mc: -key=likeMissionGroupIDkey -expire=d.mcItemExpire -encode=pb
	AddCacheMissionGroupItems(ctx context.Context, val map[int64]*likemdl.MissionGroup) error
	//mc: -key=likeActMissionKey
	CacheActMission(c context.Context, sid int64, lid int64, mid int64) (res int64, err error)
	//mc: -key=likeActMissionKey -expire=d.mcPerpetualExpire -encode=raw
	AddCacheActMission(c context.Context, sid int64, val int64, lid int64, mid int64) error
	//mc: -key=actAchieveKey
	CacheActLikeAchieves(c context.Context, sid int64) (res *likemdl.Achievements, err error)
	//mc: -key=actAchieveKey -expire=d.mcItemExpire -encode=pb
	AddCacheActLikeAchieves(c context.Context, sid int64, res *likemdl.Achievements) error
	//mc: -key=actMissionFriendsKey
	CacheActMissionFriends(c context.Context, sid int64, lid int64) (res *likemdl.ActMissionGroups, err error)
	//mc: -key=actMissionFriendsKey
	DelCacheActMissionFriends(c context.Context, sid int64, lid int64) error
	//mc: -key=actMissionFriendsKey -expire=d.mcItemExpire -encode=pb
	AddCacheActMissionFriends(c context.Context, sid int64, res *likemdl.ActMissionGroups, lid int64) error
	//mc: -key=actUserAchieveKey
	CacheActUserAchieve(c context.Context, id int64) (res *likemdl.ActLikeUserAchievement, err error)
	//mc: -key=actUserAchieveKey -expire=d.mcItemExpire -encode=pb
	AddCacheActUserAchieve(c context.Context, id int64, val *likemdl.ActLikeUserAchievement) error
	//mc: -key=actUserAchieveAwardKey
	CacheActUserAward(c context.Context, id int64) (res int64, err error)
	//mc: -key=actUserAchieveAwardKey -expire=d.mcPerpetualExpire -encode=raw
	AddCacheActUserAward(c context.Context, id int64, val int64) error
	// mc: -key=subjectStatKey
	CacheSubjectStat(c context.Context, sid int64) (*likemdl.SubjectStat, error)
	// mc: -key=subjectStatKey -expire=d.mcSubStatExpire -encode=json
	AddCacheSubjectStat(c context.Context, sid int64, value *likemdl.SubjectStat) error
	// mc: -key=viewRankKey
	CacheViewRank(c context.Context, sid int64) (string, error)
	// mc: -key=viewRankKey -expire=d.mcViewRankExpire -encode=raw
	AddCacheViewRank(c context.Context, sid int64, value string) error
	// mc: -key=likeContentKey
	CacheLikeContent(c context.Context, lids []int64) (res map[int64]*likemdl.LikeContent, err error)
	// mc: -key=likeContentKey -expire=d.mcPerpetualExpire -encode=pb
	AddCacheLikeContent(c context.Context, val map[int64]*likemdl.LikeContent) error
	// mc: -key=sourceItemKey
	CacheSourceItemData(c context.Context, sid int64) ([]int64, error)
	// mc: -key=sourceItemKey -expire=d.mcSourceItemExpire -encode=json
	AddCacheSourceItemData(c context.Context, sid int64, lids []int64) error
	// mc: -key=subjectProtocolKey
	CacheActSubjectProtocol(c context.Context, sid int64) (res *likemdl.ActSubjectProtocol, err error)
	// mc: -key=subjectProtocolKey -expire=d.mcProtocolExpire -encode=pb
	AddCacheActSubjectProtocol(c context.Context, sid int64, value *likemdl.ActSubjectProtocol) error
}
