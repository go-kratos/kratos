package data

var (
	//HBaseVideoTablePrefix 播放流失分布
	HBaseVideoTablePrefix = "video_play_churn_"
	//HBaseArchiveTablePrefix 分类分端播放
	HBaseArchiveTablePrefix = "video_play_category_"
	//HBaseAreaTablePrefix 地区播放
	HBaseAreaTablePrefix = "video_play_area_"
	//HBaseUpStatTablePrefix up主概况
	HBaseUpStatTablePrefix = "up_stats_"
	//HBaseUpViewerBase 观众数据,性别年龄分布 + 设备分布
	HBaseUpViewerBase = "up_viewer_base_"
	//HBaseUpViewerArea 地区分布
	HBaseUpViewerArea = "up_viewer_area_"
	//HBaseUpViewerTrend 内容倾向
	HBaseUpViewerTrend = "up_viewer_trend_"
	//HBaseUpViewerActionHour 行为时间分布
	HBaseUpViewerActionHour = "up_viewer_action_hour_"
	//HBaseUpRelationFansDay 日维度 最近30天 只保留31天
	HBaseUpRelationFansDay = "up_relation_fans_day"
	// HBaseUpRelationFansHistory 日维度 各月份每日数据,日更,永久保存
	HBaseUpRelationFansHistory = "up_relation_fans_history"
	//HBaseUpRelationFansMonth 年维度 2017.8月以后的数据永久保存
	HBaseUpRelationFansMonth = "up_relation_fans_month"
	//HBaseUpPlayInc 我的概况 播放相关
	HBaseUpPlayInc = "up_play_inc_"
	//HBaseUpDmInc 弹幕相关
	HBaseUpDmInc = "up_dm_inc_"
	//HBaseUpReplyInc 评论相关
	HBaseUpReplyInc = "up_reply_inc_"
	//HBaseUpShareInc 分享相关
	HBaseUpShareInc = "up_share_inc_"
	//HBaseUpCoinInc 投币相关
	HBaseUpCoinInc = "up_coin_inc_"
	//HBaseUpFavInc 收藏相关
	HBaseUpFavInc = "up_fav_inc_"
	//HBaseUpElecInc 充电相关
	HBaseUpElecInc = "up_elec_inc_"
	//HBaseUpFansAnalysis  粉丝管理
	HBaseUpFansAnalysis = "up_fans_analysis"
	//HBaseUpPlaySourceAnalysis  播放来源
	HBaseUpPlaySourceAnalysis = "up_play_analysis"
	//HBaseUpArcPlayAnalysis  平均观看时长、播放用户数、留存率
	HBaseUpArcPlayAnalysis = "up_archive_play_analysis"
	//HBaseUpArcQuery  稿件索引表
	HBaseUpArcQuery = "up_archive_query"

	//HBasePlayArc 播放相关 archive for 30 days
	HBasePlayArc = "up_play_trend"
	//HBaseDmArc 弹幕相关
	HBaseDmArc = "up_dm_trend"
	//HBaseReplyArc 评论相关
	HBaseReplyArc = "up_reply_trend"
	//HBaseShareArc 分享相关
	HBaseShareArc = "up_share_trend"
	//HBaseCoinArc 投币相关
	HBaseCoinArc = "up_coin_trend"
	//HBaseFavArc 收藏相关
	HBaseFavArc = "up_fav_trend"
	//HBaseElecArc 充电相关
	HBaseElecArc = "up_elec_trend"
	//HBaseLikeArc 点赞相关
	HBaseLikeArc = "up_like_trend"

	//HBaseFamilyPlat  family
	HBaseFamilyPlat = []byte("v")
	//HBaseColumnAid aid
	HBaseColumnAid = []byte("avid")
	//HBaseColumnWebPC pc
	HBaseColumnWebPC = []byte("plat0")
	//HBaseColumnWebH5 h5
	HBaseColumnWebH5 = []byte("plat1")
	//HBaseColumnOutsite out
	HBaseColumnOutsite = []byte("plat2")
	//HBaseColumnIOS ios
	HBaseColumnIOS = []byte("plat3")
	//HBaseColumnAndroid android
	HBaseColumnAndroid = []byte("plat4")
	//HBaseColumnElse else
	HBaseColumnElse = []byte("else")
	//HBaseColumnFans fans
	HBaseColumnFans = []byte("fans")
	//HBaseColumnGuest guest
	HBaseColumnGuest = []byte("guest")
	//HBaseColumnAll all
	HBaseColumnAll = []byte("all")
	//HBaseColumnCoin coin
	HBaseColumnCoin = []byte("coin")
	//HBaseColumnElec elec
	HBaseColumnElec = []byte("elec")
	//HBaseColumnFav fav
	HBaseColumnFav = []byte("fav")
	//HBaseColumnShare share
	HBaseColumnShare = []byte("share")
)
