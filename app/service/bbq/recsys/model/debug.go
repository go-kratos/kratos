package model

//rank model, rank feature
const (
	ResponseDownGrade   = "ResponseDownGrade" // 1:recall service 2: redis
	ResponseRecallCount = "ResponseRecallCount"
	ResponseCount       = "ResponseCount"
	ResponseRecallStat  = "ResponseRecallStat"
	RankModelName       = "RankModelName"
	RankModelScore      = "RankModelScore"
	QueryID             = "QueryID"

	ScoreMessage  = "scoreMessage"
	FeatureString = "feature"

	OrderRecall           = "Order01Recall"
	OrderRanker           = "Order02Ranker"
	OrderWeakIntervention = "Order03WeakIntervention"
	OrderFinal            = "Order04Final"
	OrderPostProcess      = "OrderPostProcess"

	//rank score

	ScoreTotalScore = "TotalScore"

	ScoreBiliZone = "ScoreBiliZone"
	ScoreBiliTag  = "ScoreBiliTag"

	ScoreLikeTag = "scoreLikeTag"
	ScorePosTag  = "scorePosTag"
	ScoreNegTag  = "scoreNegTag"

	ScoreMatchTitle     = "scoreMatchTitle"
	ScoreFollowUP       = "scoreFollowUp"
	ScoreOperationLevel = "scoreOperationLevel"

	BiliPlayNum    = "BiliPlayNum"
	BiliFavRatio   = "BiliFavRatio"
	BiliLikeRatio  = "BiliLikeRatio"
	BiliShareRatio = "BiliShareRatio"
	BiliCoinRatio  = "BiliCoinRatio"
	BiliReplyRatio = "BiliReplyRatio"

	ScoreRelevant    = "scoreRelevant"
	ScoreRetrieveTag = "scoreRetrieveTag"

	//recall tag
)
