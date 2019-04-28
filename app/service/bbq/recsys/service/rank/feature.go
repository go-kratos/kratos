package rank

import (
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/retrieve"
	"go-common/app/service/bbq/recsys/service/util"
	"go-common/library/log"
	"math"
	"strconv"
	"strings"
	"time"
)

//rank feature names
const (
	MatchBBQTagLevel3     = "MatchBBQTagLevel3"
	MatchBBQTagLevel2     = "MatchBBQTagLevel2"
	MatchBBQTagCount      = "MatchBBQTagCount"
	MatchBBQTagCountScore = "MatchBBQTagCountScore"
	BBQPrefUp             = "BBQPrefUp"
	BBQFollow             = "BBQFollow"
	BBQBlack              = "BBQBlack"

	MatchBiliTagLevel3     = "MatchBiliTagLevel3"
	MatchBiliTagLevel2     = "MatchBiliTagLevel2"
	MatchBiliTagCount      = "MatchBiliTagCount"
	MatchBiliTagCountScore = "MatchBiliTagCountScore"
	BiliPrefUp             = "BiliPrefUp"

	SessionBBQFollow  = "SessionBBQFollow"
	SessionLikeTag    = "SessionLikeTag"
	SessionPosPlayTag = "SessionPosPlayTag"
	SessionNegPlayTag = "SessionNegPlayTag"
	PureNegPlayTag    = "PureNegPlayTag"

	OperationLevel = "OperationLevel"
	BiliPlayNum    = "BiliPlayNum"
	BiliFavRatio   = "BiliFavRatio"
	BiliLikeRatio  = "BiliLikeRatio"
	BiliShareRatio = "BiliShareRatio"
	BiliCoinRatio  = "BiliCoinRatio"
	BiliReplyRatio = "BiliReplyRatio"

	BBQPlayTotal  = "BBQPlayTotal"
	BBQPlayNum    = "BBQPlayNum"
	BBQFavRatio   = "BBQFavRatio"
	BBQLikeRatio  = "BBQLikeRatio"
	BBQShareRatio = "BBQShareRatio"
	BBQCoinRatio  = "BBQCoinRatio"
	BBQReplyRatio = "BBQReplyRatio"

	//recall
	HotRecall         = "HotRecall"
	SelectionRecall   = "SelectionRecall"
	BiliFollowsRecall = "BiliFollowsRecall"
	UserProfileBili   = "UserProfileBili"
	UserProfileBBQ    = "UserProfileBBQ"
	LikeI2IRecall     = "LikeI2IRecall"
	PosI2IRecall      = "PosI2IRecall"
	LikeTagRecall     = "LikeTagRecall"
	PosTagRecall      = "PosTagRecall"
	FollowRecall      = "FollowRecall"
	RandomRecall      = "RandomRecall"

	LikeUpTimeDiff  = "LikeUpTimeDiff"
	FollowTimeDiff  = "FollowTimeDiff"
	LikeI2ITimeDiff = "LikeI2ITimeDiff"
	LikeTagCount    = "LikeTagCount"
)

//FeatureLog is feature for log
type FeatureLog struct {

	// Record
	record *recsys.RecsysRecord

	// score
	Score float64 `json:"score,omitempty"`

	// user feature
	MID    int64  `json:"mid,omitempty"`
	BUVID  string `json:"buvid,omitempty"`
	Gender int8   `json:"gender,omitempty"`
	Age    int8   `json:"age,omitempty"`

	ViewVideoNum int `json:"ViewVideoNum,omitempty"`

	// item feature
	//item feature: attribute
	SVID int64 `json:"svid,omitempty"`
	AVID int64 `json:"avid,omitempty"`
	CID  int64 `json:"cid,omitempty"`

	PubTime      int64 `json:"pubtime,omitempty"`
	PubTimeToNow int64 `json:"pubtimetonow,omitempty"`
	TagID1       int64 `json:"tagid1,omitempty"`
	TagID2       int64 `json:"tagid2,omitempty"`
	ZoneID       int64 `json:"ZoneID,omitempty"`
	Duration     int64 `json:"duration,omitempty"`
	Width        int64 `json:"width,omitempty"`
	Height       int64 `json:"height,omitempty"`
	Rotate       int64 `json:"rotate,omitempty"`
	State        int64 `json:"State,omitempty"`

	//item feature: feedback
	// bili
	PlayB  int64 `json:"playb,omitempty"`
	FavB   int64 `json:"favb,omitempty"`
	LikeB  int64 `json:"likeb,omitempty"`
	ShareB int64 `json:"shareb,omitempty"`
	ReplyB int64 `json:"replyb,omitempty"`
	CoinB  int64 `json:"coinb,omitempty"`

	// bbq
	PlayBBQTotal  int64 `json:"PlayBBQTotal,omitempty"`
	PlayBBQ       int64 `json:"PlayBBQ,omitempty"`
	PlayBBQFinish int64 `json:"PlayBBQFinish,omitempty"`
	LikeBBQ       int64 `json:"LikeBBQ,omitempty"`
	ShareBBQ      int64 `json:"ShareBBQ,omitempty"`
	ReplyBBQ      int64 `json:"ReplyBBQ,omitempty"`

	// retrieve feature
	RecallClasses  string `json:"RecallClasses,omitempty"`
	RetrieveName   string `json:"retrievename,omitempty"`
	RetrieveNum    int64  `json:"retrievenum,omitempty"`
	OperationLevel int64  `json:"operationlevel,omitempty"`

	// user-item feature
	MatchBiliTagLevel1 int64 `json:"matchtaglevel1,omitempty"`
	MatchBiliTagLevel2 int64 `json:"matchtaglevel2,omitempty"`
	MatchBiliTagLevel3 int64 `json:"matchtaglevel3,omitempty"`
	MatchTitle         int64 `json:"matchtitle,omitempty"`
	MatchBBQTagLevel2  int64 `json:"MatchBBQTagLevel2,omitempty"`
	MatchBBQTagLevel3  int64 `json:"MatchBBQTagLevel3,omitempty"`

	MatchBiliTagCount int `json:"MatchBiliTagCount,omitempty"`
	MatchBBQTagCount  int `json:"MatchBBQTagCount,omitempty"`

	// user-item-up feature
	BBQPrefUp  int64 `json:"BBQPrefUp,omitempty"`
	BiliPrefUp int64 `json:"BiliPrefUp,omitempty"`
	BBQFollow  int64 `json:"BBQFollow,omitempty"`
	BBQBlack   int64 `json:"BBQBlack,omitempty"`
	LikeAuthor int64 `json:"likeauthor,omitempty"`
	PlayAuthor int64 `json:"playauthor,omitempty"`

	// user-item feature: bbq session feature
	SessionBBQFollow int64 `json:"SessionBBQFollow,omitempty"`
	SessionLikeI2I   int64 `json:"SessionLikeI2I,omitempty"`
	SessionLikeTag   int64 `json:"SessionLikeTag,omitempty"`
	SessionLikeTag1  int64 `json:"SessionLikeTag1,omitempty"`
	SessionLikeTag2  int64 `json:"SessionLikeTag2,omitempty"`
	SessionLikeTag3  int64 `json:"SessionLikeTag3,omitempty"`

	SessionPosPlayTag  int64 `json:"SessionPosPlayTag,omitempty"`
	SessionPosPlayTag1 int64 `json:"SessionPosPlayTag1,omitempty"`
	SessionPosPlayTag2 int64 `json:"SessionPosPlayTag2,omitempty"`
	SessionPosPlayTag3 int64 `json:"SessionPosPlayTag3,omitempty"`

	SessionNegPlayTag  int64 `json:"SessionNegPlayTag,omitempty"`
	PureNegPlayTag     int64 `json:"PureNegPlayTag,omitempty"`
	SessionNegPlayTag1 int64 `json:"SessionNegPlayTag1,omitempty"`
	SessionNegPlayTag2 int64 `json:"SessionNegPlayTag2,omitempty"`
	SessionNegPlayTag3 int64 `json:"SessionNegPlayTag3,omitempty"`

	MatchLast1VideoTag1 int64 `json:"matchlast1videotag1,omitempty"`
	MatchLast1VideoTag2 int64 `json:"matchlast1videotag2,omitempty"`
	MatchLast1VideoTag3 int64 `json:"matchlast1videotag3,omitempty"`

	Last1VideoTag3 int64

	// context feature
	HourOfDay int64 `json:"hourofday,omitempty"`
	DayOfWeek int64 `json:"dayofweek,omitempty"`
}

//BuildFeature ...
func BuildFeature(record *recsys.RecsysRecord, userProfile *model.UserProfile) (featureLog *FeatureLog, featureValueMap map[string]float64) {

	now := time.Now().Unix()

	featureValueMap = make(map[string]float64)
	featureLog = &FeatureLog{}

	featureLog.record = record
	featureLog.MID = userProfile.Mid
	featureLog.BUVID = userProfile.Buvid
	featureLog.SVID = record.Svid
	featureLog.AVID, _ = strconv.ParseInt(record.Map[model.AVID], 10, 64)
	featureLog.CID, _ = strconv.ParseInt(record.Map[model.CID], 10, 64)

	featureLog.Duration, _ = strconv.ParseInt(record.Map[model.Duration], 10, 64)
	featureLog.ZoneID, _ = strconv.ParseInt(record.Map[model.ZoneID], 10, 64)

	// recall feature
	featureLog.RecallClasses = record.Map[model.RecallClasses]
	recallClasses := strings.Split(record.Map[model.RecallClasses], "|")
	for _, recallClass := range recallClasses {
		featureValueMap[recallClass] = 1
	}
	recallTags := strings.Split(record.Map[model.RecallTags], "|")
	for _, recallTag := range recallTags {
		if strings.HasPrefix(recallTag, retrieve.RecallKeyTagIDPrefix) {
			fields := strings.Split(recallTag, ":")
			if len(fields) >= 3 {
				tagStr := strings.Split(recallTag, ":")[2]
				sourceTagID, _ := strconv.ParseInt(tagStr, 10, 64)
				if count, ok := userProfile.LikeTagIDs[sourceTagID]; ok {
					featureValueMap[LikeTagCount] = util.ScoreCount(float64(count))
				}
			} else {
				log.Error("feature error like tag recall")
			}
		}

		if strings.HasPrefix(recallTag, retrieve.RecallKeyI2IPrefix) {
			fields := strings.Split(recallTag, ":")
			if len(fields) >= 3 {
				I2IStr := strings.Split(recallTag, ":")[2]
				sourceID, _ := strconv.ParseInt(I2IStr, 10, 64)
				if sourceTimestamp, ok := userProfile.LikeVideos[sourceID]; ok {
					timeDiff := math.Max(float64(now-sourceTimestamp), 0)
					timeDiffScore := math.Max(util.ScoreTimeDiff(timeDiff), 0)
					featureValueMap[LikeI2ITimeDiff] = timeDiffScore
					record.Map[model.SourceTimeToNow] = strconv.Itoa(int(timeDiff))
				}
			} else {
				log.Error("feature error like i2i recall")
			}
		}
		if strings.HasPrefix(recallTag, retrieve.RecallKeyUpIDPrefix) {
			fields := strings.Split(recallTag, ":")
			if len(fields) >= 3 {
				str := strings.Split(recallTag, ":")[2]
				sourceUpID, _ := strconv.ParseInt(str, 10, 64)
				if sourceTimestamp, ok := userProfile.BBQFollowAction[sourceUpID]; ok {
					timeDiff := math.Max(float64(now-sourceTimestamp), 0)
					timeDiffScore := math.Max(util.ScoreTimeDiff(timeDiff), 0)
					featureValueMap[FollowTimeDiff] = timeDiffScore
					record.Map[model.SourceTimeToNow] = strconv.Itoa(int(timeDiff))
				}
				if sourceTimestamp, ok := userProfile.LikeUPs[sourceUpID]; ok {
					timeDiff := math.Max(float64(now-sourceTimestamp), 0)
					timeDiffScore := math.Max(util.ScoreTimeDiff(timeDiff), 0)
					featureValueMap[LikeUpTimeDiff] = timeDiffScore
					record.Map[model.SourceTimeToNow] = strconv.Itoa(int(timeDiff))
				}
			} else {
				log.Error("feature error follow recall")
			}
		}
	}

	// user feature
	featureLog.ViewVideoNum = len(userProfile.ViewVideos)

	// user tag && item tag
	matchBiliTagCount := 0
	matchTagCount := 0
	itemTagIDs := strings.Split(record.Map[model.TagsID], "|")
	tagCount := len(itemTagIDs)
	for _, tagIDStr := range itemTagIDs {

		//bili user userProfile tag
		if tagScore, ok := userProfile.BiliTags[tagIDStr]; ok {
			if tagScore > 0 {
				featureLog.MatchBiliTagLevel3 = 1
				featureValueMap[MatchBiliTagLevel3] = 1
				matchBiliTagCount++
			}
		}

		if tagScore, ok := userProfile.Zones2[tagIDStr]; ok {
			if tagScore > 0 {
				featureLog.MatchBiliTagLevel2 = 1
				featureValueMap[MatchBiliTagLevel2] = 1
				matchBiliTagCount++
			}
		}

		//bbq user userProfile tag
		if tagScore, ok := userProfile.BBQZones[tagIDStr]; ok {
			if tagScore > 0 {
				featureLog.MatchBBQTagLevel2 = 1
				featureValueMap[MatchBBQTagLevel2] = 1
				matchTagCount++
			}
		}

		if tagScore, ok := userProfile.BBQTags[tagIDStr]; ok {
			if tagScore > 0 {
				featureLog.MatchBBQTagLevel3 = 1
				featureValueMap[MatchBBQTagLevel3] = 1
				matchTagCount++
			}
		}

		// bbq user session tag FIXME
		tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
		if timestamp, ok := userProfile.LikeTagIDs[tagID]; ok {
			featureLog.SessionLikeTag = 1
			timeDiff := math.Max(float64(now-timestamp), 0)
			timeDiffScore := math.Max(util.ScoreTimeDiff(timeDiff), 0)
			featureValueMap[SessionLikeTag] = timeDiffScore
		}

		if count, ok := userProfile.PosTagIDs[tagID]; ok {
			featureLog.SessionPosPlayTag = count
			featureValueMap[SessionPosPlayTag] = util.ScoreCount(float64(count))
		}

		if count, ok := userProfile.NegTagIDs[tagID]; ok {
			featureLog.SessionNegPlayTag = count
			featureValueMap[SessionNegPlayTag] = util.ScoreCount(float64(count))

			// pure negative tag
			if _, ok := userProfile.PosTagIDs[tagID]; !ok {
				featureLog.PureNegPlayTag = count
				featureValueMap[PureNegPlayTag] = util.ScoreCount(float64(count))
			}
		}
	}

	featureLog.MatchBBQTagCount = matchTagCount
	featureValueMap[MatchBBQTagCount] = float64(matchTagCount)
	if matchTagCount > 0 {
		featureValueMap[MatchBBQTagCountScore] = (float64(matchTagCount) + 1.0) / (float64(tagCount) + 1.0)
	}

	featureLog.MatchBiliTagCount = matchBiliTagCount
	featureValueMap[MatchBiliTagCount] = float64(matchBiliTagCount)
	if matchBiliTagCount > 0 {
		featureValueMap[MatchBiliTagCountScore] = (float64(matchBiliTagCount) + 1.0) / (float64(tagCount) + 1.0)
	}

	// user-up feature
	upMidStr := record.Map[model.UperMid]
	upMid, _ := strconv.ParseInt(upMidStr, 10, 64)
	if _, ok := userProfile.FollowUps[upMid]; ok {
		featureLog.BiliPrefUp = 1
		featureValueMap[BiliPrefUp] = 1
	}
	if _, ok := userProfile.BBQPrefUps[upMid]; ok {
		featureLog.BBQPrefUp = 1
		featureValueMap[BBQPrefUp] = 1
	}
	//最近关注
	if _, ok := userProfile.BBQFollowAction[upMid]; ok {
		featureLog.SessionBBQFollow = 1
		featureValueMap[SessionBBQFollow] = 1
	}
	//全部关注
	if _, ok := userProfile.BBQFollow[upMid]; ok {
		featureLog.BBQFollow = 1
		featureValueMap[BBQFollow] = 1
	}
	//拉黑
	if _, ok := userProfile.BBQBlack[upMid]; ok {
		featureLog.BBQBlack = 1
		featureValueMap[BBQBlack] = 1
	}

	//Up feature TODO

	// item feature

	stateStr := record.Map[model.State]
	state, _ := strconv.ParseInt(stateStr, 10, 64)
	if state == model.State5 {
		featureLog.OperationLevel = 1
		featureValueMap[OperationLevel] = 1
	}
	featureLog.State = state

	pubTime, _ := strconv.ParseInt(record.Map[model.PubTime], 10, 64)
	featureLog.PubTime = pubTime
	featureLog.PubTimeToNow = now - pubTime

	if play, ok := record.Map[model.PlayHive]; ok {
		playNum, err := strconv.ParseFloat(play, 64)
		if err == nil {
			playNumScore := math.Log10(math.Min(playNum+1.0, 1000000.0)) / math.Log10(1000000.0)
			featureLog.PlayB = int64(playNum)
			featureValueMap[BiliPlayNum] = playNumScore
		}

		if fav, ok := record.Map[model.FavHive]; ok {
			favNum, _ := strconv.ParseFloat(fav, 64)
			favScore := (math.Min(favNum, playNum) + 1.0) / (playNum + 200.0)
			favScore = math.Min(favScore, 0.1)
			featureLog.FavB = int64(favNum)
			featureValueMap[BiliFavRatio] = favScore
		}

		if likes, ok := record.Map[model.LikesHive]; ok {
			likesNum, _ := strconv.ParseFloat(likes, 64)
			likesScore := (math.Min(likesNum, playNum) + 1.0) / (playNum + 100.0)
			likesScore = math.Min(likesScore, 0.1)
			featureLog.LikeB = int64(likesNum)
			featureValueMap[BiliLikeRatio] = likesScore
		}

		if share, ok := record.Map[model.ShareHive]; ok {
			shareNum, _ := strconv.ParseFloat(share, 64)
			shareScore := (math.Min(shareNum, playNum) + 1.0) / (playNum + 500.0)
			shareScore = math.Min(shareScore, 0.1)
			featureLog.ShareB = int64(shareNum)
			featureValueMap[BiliShareRatio] = shareScore
		}

		if coin, ok := record.Map[model.CoinHive]; ok {
			coinNum, _ := strconv.ParseFloat(coin, 64)
			coinScore := (math.Min(coinNum, playNum) + 1.0) / (playNum + 200.0)
			coinScore = math.Min(coinScore, 0.1)
			featureLog.CoinB = int64(coinNum)
			featureValueMap[BiliCoinRatio] = coinScore
		}

		if reply, ok := record.Map[model.ReplyHive]; ok {
			replyNum, _ := strconv.ParseFloat(reply, 64)
			replyScore := (math.Min(replyNum, playNum) + 1.0) / (playNum + 500.0)
			replyScore = math.Min(replyScore, 0.1)
			featureLog.ReplyB = int64(replyNum)
			featureValueMap[BiliReplyRatio] = replyScore
		}
	}

	// bbq video feature
	playMonthTotal, _ := strconv.ParseFloat(record.Map[model.PlayMonthTotal], 64)
	featureLog.PlayBBQTotal = int64(playMonthTotal)
	featureValueMap[BBQPlayTotal] = playMonthTotal

	if play, ok := record.Map[model.PlayMonth]; ok {
		playNum, _ := strconv.ParseFloat(play, 64)
		playNumScore := math.Log10(math.Min(playNum+1.0, 1000000.0)) / math.Log10(1000000.0)
		featureLog.PlayBBQ = int64(playNum)
		featureValueMap[BBQPlayNum] = playNumScore
		featureLog.PlayBBQ = int64(playNum)

		if likes, ok := record.Map[model.LikesMonth]; ok {
			likesNum, _ := strconv.ParseFloat(likes, 64)
			likesScore := (math.Min(likesNum, playNum) + 1.0) / (playNum + 100.0)
			likesScore = math.Min(likesScore, 0.1)
			featureLog.LikeBBQ = int64(likesNum)
			featureValueMap[BBQLikeRatio] = likesScore
		}

		if share, ok := record.Map[model.ShareMonth]; ok {
			shareNum, _ := strconv.ParseFloat(share, 64)
			shareScore := (math.Min(shareNum, playNum) + 1.0) / (playNum + 500.0)
			shareScore = math.Min(shareScore, 0.1)
			featureLog.ShareBBQ = int64(shareNum)
			featureValueMap[BBQShareRatio] = shareScore
		}

		if reply, ok := record.Map[model.ReplyMonth]; ok {
			replyNum, _ := strconv.ParseFloat(reply, 64)
			replyScore := (math.Min(replyNum, playNum) + 1.0) / (playNum + 500.0)
			replyScore = math.Min(replyScore, 0.1)
			featureLog.ReplyBBQ = int64(replyNum)
			featureValueMap[BBQReplyRatio] = replyScore
		}
	}
	return
}
