package common

import "go-common/library/xstr"

/*
https://www.tapd.cn/66539426/prong/stories/view/1166539426001102539
state:
视频的状态有以下几种，如何选取合适的状态是该部分的关键
我们拆出了两个层次：业务层&详情层，先通过业务层获取svid，再从详情层获取svid的详情
	详情层：根据svid获取相应的视频base信息、play信息、user信息等
	业务层：根据不同的业务，设置相应的state集合，根据state去获取相应的svid（当然这里也有两种业务层区分，在于从其他服务获取svid（这里还需要做个后过滤，保证视频展示是正确的），还是在video-c中获取）
也就是说，state的选出在于业务层而不是详情层，需要业务层自己去保护，而详情层在这里仅仅做一次较宽松的保证（也就是完全不可见视频的过滤，如下架状态）

仅自己可见的状态说明：该状态可能在同一个业务场景出现主客态区分
个人空间：根据是否主客态选择select的state
通知中心：选择所有可见视频，再根据当前视频是否属于自己可见进行过滤
//点赞列表：暂时不区分主客态

业务说明：
个人空间页：作品列表区分主客态
关注页：除自己可见外的所有状态
搜索：
推荐feed页：
通知中心，含详情中转页：
评论
点赞
分享
*/

//视频状态集合
const (
	//VideoStRecommend 精选，在APP端加权露出
	VideoStRecommend = 5
	//VideoStHighGrade 优质，回查被选为优质，在APP端普通露出
	VideoStHighGrade = 4
	//VideoStCanPlay 回查可放出，在APP端普通露出
	VideoStCanPlay = 3
	//VideoStCheckBack 待冷启动回查，在APP端部分区域露出
	VideoStCheckBack = 2
	//VideoStPassReview 新鲜，安全审核通过，在APP端普通露出
	VideoStPassReview = 1
	//VideoStPendingPassReview 新鲜，未安全审核，在APP端普通露出
	VideoStPendingPassReview = 0
	//VideoStPassReviewReject 待安全审核，在APP端仅自见
	VideoStPassReviewReject = -1
	//VideoStCheckBackPatialPlay 回查不放出，在APP部分放出
	VideoStCheckBackPatialPlay = -2
	//VideoInActive 安全审核不通过，在APP端不可见，待物理删除
	VideoInActive = -3
	//VideoDeleted Up主删除，在APP端不可见，待物理删除
	VideoDeleted = -4
)

//SvAllOutState APP全部可露出状态
var SvAllOutState = []int16{
	VideoStPendingPassReview,
	VideoStPassReview,
	VideoStCanPlay,
	VideoStHighGrade,
	VideoStRecommend,
}

/*
以下用于最后根据svid获取详情时的过滤，用于那些从其他服务获取svid的业务：推荐页、搜索页、点赞列表等
*/

// IsSvStateAvailable 广义上是否可见，包含用户自见state，用于获取详情，当前和owner available一致
func IsSvStateAvailable(state int64) (available bool) {
	return IsSvStateOwnerAvailable(state)
}

// IsSvStateGuestAvailable 客态可见的视频状态
func IsSvStateGuestAvailable(state int64) (available bool) {
	_, available = svGuestAvailableState[state]
	return
}

// IsSvStateOwnerAvailable 主态可见的视频状态
func IsSvStateOwnerAvailable(state int64) (available bool) {
	if state == VideoStPassReviewReject {
		return true
	}
	return IsSvStateGuestAvailable(state)
}

// IsRecommendSvStateAvailable 推荐页中的状态过滤
func IsRecommendSvStateAvailable(state int64) (available bool) {
	_, available = svRecommendAvailableState[state]
	return
}

// IsSearchSvStateAvailable 搜索页中的状态过滤
func IsSearchSvStateAvailable(state int64) (available bool) {
	// 暂时复用推荐
	return IsRecommendSvStateAvailable(state)
}

// IsTopicSvStateAvailable 话题页中的状态过滤
func IsTopicSvStateAvailable(state int64) (available bool) {
	// 暂时复用推荐
	return IsRecommendSvStateAvailable(state)
}

var svGuestAvailableState = map[int64]bool{
	VideoStCheckBackPatialPlay: true,
	VideoStPendingPassReview:   true,
	VideoStPassReview:          true,
	VideoStCanPlay:             true,
	VideoStHighGrade:           true,
	VideoStRecommend:           true,
	VideoStCheckBack:           true,
}
var svRecommendAvailableState = map[int64]bool{
	VideoStPendingPassReview: true,
	VideoStPassReview:        true,
	VideoStCanPlay:           true,
	VideoStHighGrade:         true,
	VideoStRecommend:         true,
}

/*
以下用于业务逻辑在select语句中state in，用于video-c服务中自己进行选取svid的业务，如：关注页、个人空间页
*/

// FeedStates .
var FeedStates = xstr.JoinInts(svFeedOutStates)

// SpaceOwnerStates .
var SpaceOwnerStates = xstr.JoinInts(svSpaceOwnerOutStates)

// SpaceFanStates .
var SpaceFanStates = xstr.JoinInts(svSpaceFanOutStates)

// svFeedOutStates .
var svFeedOutStates = []int64{
	VideoStCheckBackPatialPlay,
	VideoStPendingPassReview,
	VideoStPassReview,
	VideoStCanPlay,
	VideoStHighGrade,
	VideoStRecommend,
	VideoStCheckBack,
}

// svSpaceOwnerOutStates .
var svSpaceOwnerOutStates = []int64{
	VideoStPassReviewReject,
	VideoStCheckBackPatialPlay,
	VideoStPendingPassReview,
	VideoStPassReview,
	VideoStCanPlay,
	VideoStHighGrade,
	VideoStRecommend,
	VideoStCheckBack,
}

// svSpaceFanOutStates .
var svSpaceFanOutStates = []int64{
	VideoStCheckBackPatialPlay,
	VideoStPendingPassReview,
	VideoStPassReview,
	VideoStCanPlay,
	VideoStHighGrade,
	VideoStRecommend,
	VideoStCheckBack,
}
