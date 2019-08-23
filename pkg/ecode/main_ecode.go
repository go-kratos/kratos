package ecode

// main ecode interval is [0,990000]
var (
	// appeal
	AppealNotExist     = New(10101) // 不存在该申诉
	AppealAlreadyClose = New(10102) // 申诉工单已经被关闭
	AppealInterval     = New(10103) // 申诉间隔时间内，不能再发起申诉
	AppealOpen         = New(10104) // 该稿件已处于申诉中
	AppealOwner        = New(10105) // 只能申诉自己的稿件
	AppealHasStar      = New(10106) // 该申诉已评分
	AppealLimit        = New(10107) // 仅支持打回和锁定的稿件申诉
	AppealedInDay      = New(10108) // 24小时内已发起过申诉
	// favorite
	FavNameTooLong        = New(11001) // 收藏夹名称过长
	FavMaxFolderCount     = New(11002) // 达到最大收藏夹数
	FavCanNotDelDefault   = New(11003) // 不能删除默认收藏夹
	FavFolderNoPublic     = New(11004) // 收藏夹目录未公开
	FavMaxVideoCount      = New(11005) // 视频数达到目录最大收藏数
	FavFolderExist        = New(11006) // 已经存在该收藏夹
	FavVideoExist         = New(11007) // 已经存在该视频了
	FavOnlyPublic         = New(11008) // 仅仅只能设置公开0 或 非公开1
	FavDefaultFolder      = New(11009) // 默认收藏夹
	FavFolderNotExist     = New(11010) // 没有该收藏夹
	FavSearchReqErr       = New(11011) // 请求搜索出错
	FavFloderAlreadyDel   = New(11012) // 已经删除该收藏夹了
	FavVideoAlreadyDel    = New(11013) // 已经取消收藏该视频了
	FavMaxOperNum         = New(11014) // 超出允许的最大操作数75
	FavFolderSame         = New(11015) // 一样的收藏夹
	FavFolderMoveFailed   = New(11016) // 收藏夹视频移动失败
	FavFolderSortErr      = New(11017) // 收藏夹列表信息错误
	FavSearchWordIllegal  = New(11018) // 收藏夹视频搜索关键词非法
	FavMaintenance        = New(11019) // 收藏服务维护中
	FavFolderBanned       = New(11020) // 收藏夹名称包含敏感词
	FavCleaneInProgress   = New(11021) // 正在删除失效视频…请过段时间再来访问
	FavCleanedLocked      = New(11022) // 清除操作锁定中
	FavResourceExist      = New(11201) // 已经收藏过了
	FavResourceAlreadyDel = New(11202) // 已经取消收藏了
	FavResourceOverflow   = New(11203) // 达到收藏上限
	FavDescTooLang        = New(11204) // 收藏夹描述过长
	FavRetryLater         = New(11205) // 请稍后重试
	FavHitSensitive       = New(11206) // 收藏夹命中敏感词
	FavMaxSortCount       = New(11207) // 内容太多啦！超过1000不支持排序
	FavFolderHidden       = New(11208) // 用户隐藏了他的收藏夹

	// elec
	ElecUserForbid    = New(13001) // 用户被禁止参与充电计划
	ElecUserAudit     = New(13002) // 用户的是否充电正在审核中
	ElecNotUpper      = New(13003) // 不是up主
	ElecArchiveForbid = New(13004) // 稿件被屏蔽

	// stat
	ClickAesDecryptErr  = New(14001) // aes解密失败
	ClickQueryFormatErr = New(14002) // 解密粗来的参数格式有问题
	// ClickServerTimeout  = New(14003) // 服务端时间超时
	ClickQuerySignErr = New(14004) // 参数中的sign错误
	ClickHmacSignErr  = New(14005) // hmac算的签名错误

	// topic
	TopicNotExist      = New(15001) // 不存在该话题
	FavTopicExist      = New(15002) // 已经存在该话题了
	FavTopicAlreadyDel = New(15003) // 已经取消收藏该话题了

	// zlimit
	ZlimitAllow     = New(17001) // 稿件无限制
	ZlimitForbidden = New(17002) // 稿件禁止查看
	ZlimitFormal    = New(17003) // 正式会员
	ZlimitPay       = New(17004) // 付费会员
	ZlimitShared    = New(17005) // 共享地址错误提示

	// short utl
	ShortURLAlreadyExist = New(30001) // 短链已存在
	ShortURLNotFound     = New(30002) // 短链不存在
	ShortURLIllegalSrc   = New(30003) // 生成短链的长链来源不合法

	// captcha
	CaptchaSignErr            = New(33001) // 签名错误
	CaptchaTSOverTolerant     = New(33002) // 时间戳超过规定时间范围
	CaptchaBusinessNotAllowed = New(33003) // 业务id不存在
	CaptchaCodeNotFound       = New(33004) // code不存在，一般是token已过期
	CaptchaTokenErr           = New(33005) // 验证时token错误
	CaptchaNotCreate          = New(33006) // 验证码未创建
	CaptchaTokenNotExist      = New(33007) // 验证码Token不存在
	CaptchaTokenExpired       = New(33008) // 验证码已过期

	// account
	AccountOverdue     = New(35001) // token过期 -658
	AccountNotLogin    = New(35002) // 未登录 -400
	AccountInexistence = New(35003) // 用户不存在 -626
	AccountAKNotFound  = New(35004) // accessKey不存在 -2

	// relation
	RelFollowSelfBanned          = New(22001) // 不能关注自己
	RelFollowBlacked             = New(22002) // 被用户拉黑，无法关注
	RelFollowAlreadyBlack        = New(22003) // 已经拉黑用户，无法关注
	RelFollowAttrAlreadySet      = New(22004) // 已经设置该属性了
	RelFollowAttrNotSet          = New(22005) // 未设置该属性，不能取消
	RelFollowReachTelLimit       = New(22006) // 关注已达上限，答题成为正式会员或者绑定手机号才能继续关注
	RelFollowingGuestLimit       = New(22007) // 访客只限制访问前五页
	RelBlackReachMaxLimit        = New(22008) // 黑名单达到上限
	RelFollowReachMaxLimit       = New(22009) // 关注失败，已达关注上限
	RelBatchFollowAlreadyBlack   = New(22010) // 部分拉黑用户未成功关注
	RelTagExistNotAllowedWords   = New(22101) // 分组名称存在不允许的字符
	RelTagNumLimit               = New(22102) // 分组数量超过限制
	RelTagLenLimit               = New(22103) // 分组名称长度超过限制
	RelTagNotExist               = New(22104) // 分组不存在
	RelTagAddFollowingFirst      = New(22105) // 请先公开关注后再添加分组
	RelTagExisted                = New(22106) // 分组已存在
	RelAwardPhoneRequired        = New(22107) // 该账号未通过手机认证
	RelAwardIsBlocked            = New(22108) // 该账号处于封禁状态
	RelAwardInsufficientFollower = New(22109) // 该账号粉丝数不符合满足10000的标准
	RelAwardGetFailed            = New(22110) // 获取粉丝成就奖励失败
	RelAwardInfoFailed           = New(22111) // 获取成就信息失败
	RelAwardInsufficientRank     = New(22112) // 很抱歉您的账号为非转正会员

	// dm
	DMFilterIllegalType   = New(36001) // 不支持该屏蔽类型
	DMFilterTooLong       = New(36002) // 屏蔽词超过上限啦（关键词50字，正则200字）
	DMFilterOverMax       = New(36003) // 屏蔽规则超过条数限制
	DMFitlerIllegalRegex  = New(36004) // 屏蔽规则正则格式不对
	DMFilterExist         = New(36005) // 屏蔽规则已经存在
	DMFilterIsEmpty       = New(36006) // 屏蔽规则不允许为空
	DMAdvNotAllow         = New(36007) // 不允许购买
	DMAdvConfirm          = New(36009) // 正在确认中
	DMAdvBought           = New(36010) // 已购买
	DMAdvNoFound          = New(36011) // 高级弹幕购买记录不存在
	DMPADMNotOwner        = New(36101) // 别人的弹幕不可以申请弹幕保护~
	DMPAUserLevel         = New(36102) // 目前仅限 lv4及以上的用户可以直接申请哦~
	DMPAUserLimit         = New(36103) // 一个人一天最多只能申请保护100条哦
	DMPADMLimit           = New(36104) // 该弹幕已经申请保护~
	DMPADMProtected       = New(36105) // 本弹幕已经被保护了~
	DMNotFound            = New(36106) // 该弹幕已被删除
	DMPAFailed            = New(36107) // 申请失败
	DMPoolLimit           = New(36108) // 弹幕池超过大小
	DMReportNotExist      = New(36201) // 举报弹幕不存在
	DMReportReasonTooLong = New(36202) // 举报原因过长
	DMReportReasonError   = New(36203) // 举报原因类型错误
	DMReportExist         = New(36204) // 已举报
	DMReportLimit         = New(36205) // 操作过于频繁，请稍后再试
	DMRecallTimeout       = New(36301) // 撤回失败，弹幕发送已过2分钟
	DMRecallDeleted       = New(36302) // 撤回失败，弹幕已经被删除或撤回
	DMRecallLimit         = New(36303) // 撤回失败，今天撤回的机会已经用完
	DMRecallError         = New(36304) // 撤回失败，服务器出错
	DMAssistNo            = New(36401) // 不是协管
	DMAssistLimit         = New(36402) // 操作次数不足
	DMTransferSame        = New(36501) // 弹幕转移源cid和目标cid相等
	DMTransferNotFound    = New(36502) // 弹幕转移cid不存在
	DMTransferNotBelong   = New(36503) // 弹幕转移cid不属于该用户
	DMTransferRepet       = New(36504) // 弹幕转移任务重复
	DMActSilence          = New(36601) // 弹幕点赞被禁言
	DMUpgrading           = New(36700) // 系统升级中
	DMMsgIlleagel         = New(36701) // 弹幕包含被禁止的内容
	DMMsgTooLong          = New(36702) // 您的弹幕长度大于100
	DMMsgPubTooFast       = New(36703) // 您发送弹幕的频率过快
	DMArchiveIlleagel     = New(36704) // 禁止向未审核的视频发送弹幕
	DMMsgNoPubPerm        = New(36705) // 您的等级不足，不能发送弹幕
	DMMsgNoPubTopPerm     = New(36706) // 您的等级不足，不能发送顶端弹幕
	DMMsgNoPubBottomPerm  = New(36707) // 您的等级不足，不能发送底端弹幕
	DMMsgNoColorPerm      = New(36708) // 您的等级不足，不能发送彩色弹幕
	DMMsgNoPubAdvancePerm = New(36709) // 您的等级不足，不能发送高级弹幕
	DMMsgNoPubStylePerm   = New(36710) // 您的权限不足，不能发送这种样式的弹幕
	DMForbidPost          = New(36711) // 该视频禁止发送弹幕
	DMMsgTooLongLevel1    = New(36712) // level 1用户发送弹幕的最大长度为20
	DMNotpayForPost       = New(36713) // 稿件未付费，不能发送弹幕
	DMProgressTooBig      = New(36714) // 弹幕发送时间不合法
	DMAssistOpToMuch      = New(36715) // 当日操作数量超过上限，请明天再试
	DMTaskRegexTooLong    = New(36800) // 任务正则过长
	DMTaskRegexIllegal    = New(36801) // 任务正则不合法
	// article
	ArtLikeDupErr            = New(37001) // 重复点赞
	ArtCancelLikeErr         = New(37002) // 取消点赞失败 用户未点赞
	ArtDislikeDupErr         = New(37003) // 重复不喜欢
	ArtCancelDislikeErr      = New(37004) // 取消不喜欢失败 用户未不喜欢
	ArtCreationNoPrivilege   = New(37101) // 创作中心:用户没有权限发文章
	ArtCreationStateErr      = New(37102) // 创作中心:文章状态错误
	ArtCreationIDErr         = New(37103) // 创作中心:文章ID错误
	ArtCreationMIDErr        = New(37104) // 创作中心:非文章作者
	ArtCreationDelPendingErr = New(37105) // 创作中心:审核中的文章不能删除
	ArtCreationDraftFull     = New(37106) // 创作中心:草稿数已达最大上限
	ArtCreationTplErr        = New(37107) // 创作中心:模板和图片数量不匹配
	ArtCreationDraftDeleted  = New(37108) // 创作中心:草稿已被删除，不可编辑
	ArtCreationArticleFull   = New(37109) // 创作中心：当日投稿数量已到达上限
	ArtUserDisabled          = New(37200) // 用户被封禁 无法操作
	ArtNoCategory            = New(37300) // 文章分区错误
	ArtApplyPass             = New(37400) // 申请已通过
	ArtApplyReject           = New(37401) // 已经申请处于冷冻期
	ArtApplySubmit           = New(37402) // 已经申请待审
	ArtApplyClose            = New(37403) // 关闭申请
	ArtApplyFull             = New(37404) // 今日申请名额已满
	ArtApplyVerify           = New(37405) // 用户未实名认证
	ArtApplyForbid           = New(37406) // 用户已被封禁
	ArtApplyPhone            = New(37407) // 用户未绑定手机
	ArtApplyPhoneVirtual     = New(37408) // 绑定的手机号是虚拟号码
	ArtCannotEditErr         = New(37409) // 文章不能被编辑
	ArtAuthorReject          = New(37410) // 申请被拒绝
	ArtNoActivity            = New(37411) // 活动未开始
	ArtMaxListErr            = New(37412) // 达到文集上限 无法再增加新文集
	ArtListNameErr           = New(37413) // 文集标题不合法
	ArtArtAddListErr         = New(37414) // 文章已存在于其他文集或者文章不存在
	ArtAddListLimitErr       = New(37415) // 达到文集文章数量上限 无法再增加新文章
	ArtPermClosedErr         = New(37416) // 操作失败，你的专栏权限已被关闭
	ArtLevelFailedErr        = New(37417) // 等级未达到要求
	ArtMediaExistedErr       = New(37418) // 已经存在长评了
	ArtUpdateFullErr         = New(37419) // 重复编辑次数已用完
	ArtOriginalEditErr       = New(37420) // 原创文章不能进行重复编辑
	ArtCreativeLimitErr      = New(37421) // 操作太快了，休息一下吧
	ArtTagBindErr            = New(27422) // 您输入的tag不合法哟
	ArtAllowErr              = New(27423) // 系统升级中
	// Member
	UpdateBirthdayFormat      = New(40001) // 出生日期格式不正确
	UpdateUnameSensitive      = New(40002) // 昵称包含敏感词
	UpdateSexError            = New(40003) // 请选择正常的性别
	UpdateUnameFormat         = New(40004) // 昵称不可包含除-和_以外的特殊字符
	UpdateUnameTooLong        = New(40005) // 昵称过长，不能修改
	UpdateUnameTooShort       = New(40006) // 昵称过短，不能修改
	UpdateUnameMoneyIsNot     = New(40007) // 硬币不足,改昵称需要6个硬币
	UpdateUnameHadVerified    = New(40008) // 已过实名验证，不能修改
	UpdateUnameHadLocked      = New(40009) // 昵称已锁定不能修改
	UpdateUnameHadOfficial    = New(40010) // 认证账号不得随意修改昵称，如有需要请联系客服娘~
	UpdateFaceFormat          = New(40012) // 头像格式错误，允许：png/jpg/jpeg/jp2
	UpdateFaceSize            = New(40013) // 头像超过限制的大小，允许2M
	UpdateUnameRepeated       = New(40014) // 昵称已存在
	MemberSignSensitive       = New(40015) // 签名包含敏感词
	MemberPhoneRequired       = New(40016) // 根据国家实名制认证的相关要求，需要绑定手机号
	MemberRealPhoneRequired   = New(40017) // 根据国家实名制认证的相关要求，需要绑定非虚拟手机号
	MemberSignHasEmoji        = New(40021) // 签名不能包含表情图片
	MemberSignOverLimit       = New(40022) // 签名最多支持70个字
	BirthdayInfoIsNull        = New(40043) // 该用户没有生日信息 // 答题系统使用
	MemberUpdate              = New(40050) // 系统维护中
	MemberBlocked             = New(40051) // 用户被封禁
	MemberNameFormatErr       = New(40052) // 用户名不合法
	MemberNameOverLimit       = New(40053) // 用户名长度超过限制
	MemberNameUnmodify        = New(40054) // 用户名未修改
	MemberNameHasEmoji        = New(40055) // 用户名包含emoji
	MemberNameCoinErr         = New(40056) // 扣除硬币失败
	MemberUnRealName          = New(40058) // 用户名未实名
	MemberCerted              = New(40059) // 用户名包含敏感词
	MemberOverLimit           = New(40060) // 批量请求超过限制
	MemberNotExist            = New(40061) // 用户不存在
	MemberUpdateBirthdayFaild = New(40071) // 修改生日失败
	MemberBirthdayNotAllow    = New(40072) // 生日信息不合法
	MemberBirthdayInfoIsNull  = New(40073) // 该用户没有生日信息
	MemberTagsOverLen         = New(40080) // 用户 Tag 不合法
	MemberTagsOverCount       = New(40081) // 用户 Tag 不合法
	SubmitOfficialDocFailed   = New(40083) // 提交官方认证请求失败
	NoOfficialDoc             = New(40084) // 未提交过官方认证请求
	SearchMidOverLimit        = New(40085) // Mid查询数量过大
	OfficialDocReasonLimit    = New(40086) // 官方认证拒绝理由字数不能超过200
	NoOfficial                = New(40087) // 当前没有认证
	OfficialDocWait           = New(40088) // 已存在待审核的官方认证
	OfficialConditionFailed   = New(40089) // 不满足官方认证条件
	MemberChangeNameToRetry   = New(40090) // 修改昵称后重试
	// Answer
	QuestionStrNotAllow     = New(41001) // 分院帽题目不合法
	QuestionAnsNotAllow     = New(41002) // 分院帽题目答案不合法
	QuestionTipsNotAllow    = New(41003) // 分院帽题目提示不合法
	QuestionTypeNotAllow    = New(41004) // 分院帽题目类型不合法
	AnswerDenied            = New(41010) // 用户答题非法访问
	AnswerTimeExpire        = New(41011) // 用户答题时间已超时
	AnswerIdsErr            = New(41012) // 用户答题提交题目id不合法
	AnswerQsNumErr          = New(41013) // 用户答题提交题目数量不合法
	AnswerBlock             = New(41014) // 用户自选题提交过快（2分钟内）被封禁12小时
	AnswerSorceZero         = New(41016) // 该用户答题分数为0
	AnswerGeetestErr        = New(41017) // 答题验证码错误
	AnswerFormalFailed      = New(41018) // 答题转正失败
	AnswerBasePassed        = New(41020) // 用户基础题已通过
	AnswerBaseNotPassed     = New(41021) // 用户基础题未通过
	AnswerHistoryNotFound   = New(41023) // 用户答题记录不存在
	AnswerMidCacheQidsErr   = New(41024) // 获取用户题目id缓存异常
	AnswerQidDiffRequestErr = New(41025) // 用户答题提交题目ID和实际用户的答题id不一致
	AnswerMidDBQueErr       = New(41026) // 获取用户DB题目信息异常
	AnswerCheckFaild        = New(41027) // 基础题检查不通过
	AnswerProNoPass         = New(41031) // 自选题未通过
	AnswerCaptchaPassed     = New(41050) // 用户答题验证码已通过
	AnswerCaptchaNoPassed   = New(41051) // 用户答题验证码未通过
	AnswerTypeIDsErr        = New(41052) // 用户题目类型不合法
	AnswerGeetestVaErr      = New(41053) // 极验验证异常
	AnswerExtraHadPass      = New(41054) // 基础附加题已通过
	AnswerExtraNoPass       = New(41055) // 基础附加题未通过
	AnswerAccCallErr        = New(41090) // 调用账号异常
	AnswerNeedBindTel       = New(41091) // 答题需要绑定手机
	AnswerDynamicErr        = New(41092) // 答题分享动态失败

	// bfs upload
	BfsUploadCodeErr                     = New(42001) // bfs响应code错误
	BfsUploadStatusErr                   = New(42002) // 返回状态错误（非常规捕捉）
	BfsRequestErr                        = New(42400) // bfs参数错误
	BfsUploadAuthErr                     = New(42401) // 上传验证错误
	BfsUplaodBucketNotExist              = New(42404) // bucket不存在
	BfsUploadServiceUnavailable          = New(42503) // 服务不可用
	BfsUploadFileTooLarge                = New(42601) // 上传的文件太大
	BfsUploadFilePixelError              = New(42602) // 不能获取图片的像素信息
	BfsUploadFilePixelWidthIllegal       = New(42603) // 宽像素不合法
	BfsUploadFilePixelHeightIllegal      = New(42604) // 高像素不合法
	BfsUploadFilePixelAspectRatioIllegal = New(42605) // 像素宽高比不合法
	BfsUploadFileContentTypeIllegal      = New(42606) // 文件类型不合法

	// remote login
	RemoteLoginStatusQueryError = New(43001) // 查询失败
	RemoteLoginFeedBackError    = New(43002) // 反馈失败
	RemoteLoginWarnCloseError   = New(43003) // 关闭失败

	// Spy
	SpyEventNotExist          = New(50001) // 反作弊事件类型不存在
	SpyServiceNotExist        = New(50002) // 反作弊服务不存在
	SpyFactorNotExist         = New(50003) // 反作弊因子不存在
	SpySettingUnknown         = New(50004) // 反作弊配置类型不存在
	SpySettingValTypeError    = New(50005) // 反作弊配置值类型错误
	SpySettingValueOutOfRange = New(50006) // 反作弊配置值超出范围
	SpyRuleNotExist           = New(50007) // 反作弊规则不存在
	SpyRulesNotMatch          = New(50008) // 反作弊规则不匹配
	// filter-service and filter-job
	FilterHitLimitBlack                   = New(52001) // 命中黑名单
	FilterHitRubLimit                     = New(52002) // 超过发送次数
	FilterLimitTypeNotExist               = New(52003) // 限制类型不存
	FilterLimitContentNotExist            = New(52004) // 限制关键词不存在
	FilterHitStrictLimit                  = New(52005) // 命中严格限制
	FilterIllegalRegexp                   = New(52006) // 非法正则
	FilterIllegalArea                     = New(52007) // 业务不存在
	FilterWhiteSampleHit                  = New(52010) // 敏感词可能误杀较大
	FilterBlackSampleHit                  = New(52011) // 敏感词导致高危内容失效
	FilterDuplicateContent                = New(52012) // 已存在相同内容敏感词/白名单
	FilterRegexpError1                    = New(52013) // 含有.*容易引起误伤，请换用.{0,10}
	FilterRegexpError2                    = New(52014) // 含有||容易引起误伤
	FilterInvalidAreaGroupName            = New(52020) // 不合法的业务组命名
	FilterDuplicateAreaGroup              = New(52021) // 业务组重复
	FilterAreaGroupNotFound               = New(52022) // 业务组不存在
	FilterInvalidAreaShowName             = New(52023) // 不合法的业务模块命名
	FilterInvalidAreaName                 = New(52024) // 不合法的业务id命名
	FilterDuplicateArea                   = New(52025) // 业务重复
	FilterInvalidArea                     = New(52026) // 业务不存在
	FilterInvalidAIWhiteUID               = New(52027) // AI过滤白名单重复添加
	FilterDuplicatedContentFromOtherAreas = New(52028) // 已存在相同内容敏感词/白名单, 在其他area中
	FilterMultiDuplicatedContents         = New(52029) // 存在多个相同内容的敏感词

	// search
	SearchArchiveCheckFailed     = New(54001) // 搜索稿件管理失败
	SearchArticleDataFailed      = New(54002) // 搜索专栏数据失败
	SearchReplyRecordFailed      = New(54003) // 搜索个人中心评论记录数据失败
	SearchBlockedListFailed      = New(54010) // 搜索风纪委封禁列表失败
	SearchBlockedPublishFailed   = New(54011) // 搜索公告列表失败
	SearchBlockedCaseFailed      = New(54012) // 搜索案件列表失败
	SearchBlockedJuryFailed      = New(54013) // 搜索委员列表失败
	SearchBlockedOpinionFailed   = New(54014) // 搜索观点列表失败
	SearchWorkflowGroupFailed    = New(54015) // 工作流获取反馈列表失败
	SearchWorkflowChaFailed      = New(54016) // 工作流获取用户工单失败
	SearchWorkflowTagFailed      = New(54017) // 工作流获取Tag列表失败
	SearchWorkflowLogFailed      = New(54018) // 工作流获取日志列表失败
	SearchWorkflowCommonFaild    = New(54019) // 工作流举报列表获取失败
	SearchUpdateIndexFailed      = New(54900) // 更新索引失败
	SearchDmFailed               = New(54020) // 弹幕列表获取失败
	SearchVideoFailed            = New(54021) // 弹幕列表获取失败
	SearchMusicSongsFailed       = New(54022) // 音乐审核列表获取失败
	SearchKpiPointFailed         = New(54023) // 搜索kpi评分列表失败
	SearchWorkflowFeedbackFailed = New(54024) // 工作流获取feedback列表失败
	SearchUNameFailed            = New(54025) // 用户昵称查询失败
	SearchDmmonitorFailed        = New(54026) // 弹幕监控查询失败
	SearchPgcMediaFailed         = New(54027) // pgc影视查询失败
	SearchFeedbackFailed         = New(54028) // 用户反馈查询失败
	SearchFeedbackReplyFailed    = New(54029) // 用户反馈报告查询失败
	SearchLogAuditFailed         = New(54030) // 审核日志查询失败
	SearchLogAuditOidFailed      = New(54031) // 根据oid查询审核日志查询失败
	SearchLogUserActionFailed    = New(54032) // 用户操作日志查询失败
	SearchUserApplyReviewsFailed = New(54033) // 用户头像挂件查询失败
	SearchBusinessFailed         = New(54901) // 后台接口Business报错
	SearchAppidFailed            = New(54902) // 后台接口Appid报错
	SearchBusinessExistErr       = New(54903) // 该业务已经存在
	SearchAssetExistErr          = New(54904) // 该数据源已经存在
	SearchAppExistErr            = New(54905) // 该应用已经存在
	// figure 信用分服务
	FigureNotFound = New(55001) // 未找到用户信用分
	// workflow 工作流
	WkfGroupNotFound                  = New(56001) // 未找到工单
	WkfChallNotFound                  = New(56002) // 未找到工单详情
	WkfAppealNotFound                 = New(56003) // 未找到申诉单
	WkfAppealNotUserOwned             = New(56201) // 只能查看或操作自己的申诉
	WkfAppealTransferStateIllegal     = New(56202) // 不合法的申诉流转状态
	WkfBusinessNotFound               = New(56401) // 未找到业务id
	WkfBusinessNotConsistent          = New(56402) // 工单业务id不一致
	WkfTagNotFound                    = New(56403) // 未找到tid数据
	WkfBusinessCallbackConfigNotFound = New(56404) // 未找到业务回调配置
	WkfBanNotSupportBatchOperate      = New(56405) // 不支持批量封禁账号
	WkfBidNotSupportPublicReferee     = New(56406) // 业务不支持移交众裁
	WkfBidNotSupportQuerySource       = New(56407) // 业务不支持查询来源
	WkfBidNotSupportQueryContentState = New(56408) // 业务不支持查询内容状态
	WkfSearchGroupFailed              = New(56501) // es search工单失败
	WkfSearchChallFailed              = New(56502) // es search工单详情失败
	WkfGetBlockInfoFailed             = New(56503) // 获取封禁信息失败
	WkfSetPublicRefereeFailed         = New(56504) // 提交众裁失败
	WkfPlatformGetLockFailed          = New(56505) // 工作台获取锁失败
	WkfPlatformDelLockFailed          = New(56506) // 工作台释放锁失败
	WkfPlatformSearchEmpty            = New(56507) // 工作台没有搜索到内容
	WkfPlatformNoOnline               = New(56508) // 工作台不在线
	WkfBusinessAttrNotFound           = New(56509) // 未找到业务属性
	// account common
	UserLoginInvalid      = New(61000) // 使用登录状态访问了，并且登录状态无效，客服端可以／需要删除登录状态
	UserCheckNoPhone      = New(61001) // 根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作
	UserCheckInvalidPhone = New(61002) // 根据国家实名制认证的相关要求，您需要换绑一个非170/171的手机号，才能继续进行操作

	// usersuit
	UsersuitInviteLevelLow               = New(64001) // 你还不满足购买激活码的条件哦，升级到Lv5再来吧~
	UsersuitInviteReachCurrentMonthLimit = New(64002) // 当月邀请码申请数达到上限
	UsersuitInviteAlreadyFormal          = New(64003) // 已经转正不能使用邀请码
	UsersuitInviteCodeNotExists          = New(64004) // 邀请码不存在
	UsersuitInviteCodeUsed               = New(64005) // 邀请码已使用
	UsersuitInviteCodeExpired            = New(64006) // 邀请码已过期
	UsersuitInviteReachMaxGeneLimit      = New(64007) // 超过批量生成邀请码上限（最多1000个）
	UsersuitInviteCodeNotUsed            = New(64008) // 邀请码未使用
	UsersuitInviteCodeImidError          = New(64009) // 受邀请人mid不匹配
	UsersuitInviterError                 = New(64010) // 邀请人信息错误
	UsersuitInviteLocked                 = New(64011) // 邀请码被锁定，请稍后再试
	// pendant 挂件相关
	PendantNotFound          = New(64101) // 挂件不存在
	PendantCanNotBuy         = New(64102) // 大会员挂件不能购买
	PendantAlreadyGet        = New(64103) // 大会员挂件已领取过
	PendantGetVIPErr         = New(64104) // 获取大会员信息错误
	PendantPayErr            = New(64105) // 订单接口错误
	PendantOrderNotFound     = New(64106) // 订单不存在
	PendantPackageNotFound   = New(64107) // 背包里无此挂件
	PendantPayTypeErr        = New(64108) // 该挂件无此种支付方式
	PendantVIPOverdue        = New(64109) // 大会员过期
	PendantAboveSendMaxLimit = New(64110) // 超过批量发放挂件上限 (最多1000个)
	// usersuit medal 勋章
	MedalNotFound = New(64201) // 勋章不存在
	MedalNotGet   = New(64202) // 不拥有该勋章
	MedalHasGet   = New(64203) // 已拥有该勋章
	// thumbup
	ThumbupBusinessBlankErr = New(65001) //业务id不存在
	ThumbupOriginErr        = New(65002) //origin id 错误
	ThumbupBusinessErr      = New(65003) //未开通此业务
	ThumbupCancelLikeErr    = New(65004) // 取消点赞失败 用户未点赞
	ThumbupCancelDislikeErr = New(65005) // 取消不喜欢失败 用户未不喜欢
	ThumbupDupLikeErr       = New(65006) // 重复点赞
	ThumbupDupDislikeErr    = New(65007) // 重复点踩
	// sms
	SmsTemplateNotExist       = New(66001) // 模版不存在
	SmsTemplateParamNotEnough = New(66002) // 模版参数不足
	SmsTemplateCodeExist      = New(66003) // 模版code已存在
	SmsTemplateParamIllegal   = New(66004) // 模版参数值不合法
	SmsTemplateModifyForbind  = New(66005) // 修改已审核的模版必须提供approver3
	SmsTemplateNotAct         = New(66006) // 模版不是营销短信
	SmsSendBatchOverLimit     = New(66023) // 批量发送超出限制
	SmsSendBothMidAndMobile   = New(66024) // mid,mobile只能传一个参数
	SmsMobilePatternErr       = New(66031) // 手机号码格式不正确
	// growup admin and interface and job
	GrowupDisabled                  = New(68001) // up主在黑名单
	GrowupTagForbit                 = New(68002) // 不允许操作该标签
	GrowupNotExist                  = New(68003) // 有不存在的UP主账号
	GrowupAuthorityUserNotFound     = New(68004) // 用户在权限系统中不存在
	GrowupTagAddForbit              = New(68005) // 该标签在本业务分区下已存在
	GrowupAuthorityExist            = New(68006) // 用户名/任务组名/角色组名/权限点名已存在
	GrowupBodyTooLarge              = New(68007) // 上传的文件太大
	GrowupBodyNotExist              = New(68008) // 上传的文件无内容
	GrowupGetTypeError              = New(68009) // 获取视频全部分区失败
	GrowupGetActivityError          = New(68010) // 根据活动ID获取稿件失败
	GrowupPriceErr                  = New(68020) // 购买大会员价格错误
	GrowupPriceNotEnough            = New(68021) // 购买大会员余额不足
	GrowupBuyErr                    = New(68022) // 激励兑换购买失败
	GrowupGoodsNotExist             = New(68023) // 激励兑换商品不存在
	GrowupGoodsTimeErr              = New(68024) // 不在兑换时间内
	GrowupAlreadyApplied            = New(68025) // up主已申请过创作激励计划
	GrowupAccTypeMismatching        = New(68026) // up主添加类型与原类型不匹配
	GrowupArchiveNotYours           = New(68101) // 稿件不属于此人
	GrowupSubTidNotExist            = New(68102) // 此二级分类不存在
	GrowupActivityCountNotEnough    = New(68103) // 活动列表不存在
	GrowupRecommendUpNotExist       = New(68104) // 推荐UP主列表不存在
	GrowupUpInfoNotExist            = New(68105) // UP主信息不存在
	GrowupRecommendUpInfoNotExist   = New(68106) // 推荐UP主信息不存在
	GrowupTidNotExist               = New(68107) // 此一级分类不存在
	GrowupSpecialAwardJoined        = New(68201) // 已参加过专项奖
	GrowupSpecialAwardUnqualified   = New(68202) // 没有资格参加专项奖
	GrowupSpecialAwardUnreserve     = New(68203) // 无法预约下期专项奖
	GrowupFansNotExist              = New(68204) // 粉丝数据不存在
	GrowupUpNoIncome                = New(68205) // up主没有收入
	GrowupWithdrawTypeErr           = New(68206) // 提现方式错误
	GrowupManualWithdrawDateErr     = New(68207) // 手动提现时间错误
	GrowupManualWithdrawIncomeErr   = New(68208) // 手动提现金额错误
	GrowupWithdrawUpdateTypeDateErr = New(68209) // 每月2到6号不能修改提现方式

	SvenRepeat                  = New(70001) // sven 数据重复
	CanalAddrFmtErr             = New(70002) //canal地址格式错误{host:port}!
	CanalAddrExist              = New(70003) //canal地址已存在
	CanalAddrNotFound           = New(70004) //canal地址不存在
	CanalApplyUpdateErr         = New(70005) //canal申请信息更新失败
	CanalApplyErr               = New(70006) //canal申请失败
	DatabasesUnmarshalErr       = New(70007) //Databases解析失败
	GetConfigByNameErr          = New(70008) //根据name获取config信息失败
	DatabusGroupErr             = New(70009) //databus group信息获取失败
	DatabusAppErr               = New(70010) //databus app信息获取失败
	DatabusActionErr            = New(70011) //databus action信息获取失败
	ConfigCreateErr             = New(70012) //配置中心生成配置文件失败
	ConfigUpdateErr             = New(70013) //配置中心更新配置文件失败
	SetConfigIDErr              = New(70014) //canal更新configId失败
	CheckMasterErr              = New(70015) //canalcheckMaster验证失败
	ConfigParseErr              = New(70016) //canal config信息解析失败
	NeedInfoErr                 = New(70017) // 需求不存在
	NeedEditErr                 = New(70018) // 需求不符合编辑要求
	NeedVerifyErr               = New(70019) //需求已审核
	ConfigNotNow                = New(70020) //配置中心配置源文件非最新
	DatabusDuplErr              = New(70021) //databus group重复
	// share
	ShareAlreadyAdd = New(71000) // 已经分享过了
	// push
	PushSensitiveWordsErr   = New(72001) // 推送信息中有敏感词
	PushUUIDErr             = New(72002) // 调用添加推送任务接口请求重放了
	PushBizAuthErr          = New(72003) // 业务方调用接口时token校验未通过
	PushSilenceErr          = New(72004) // 业务方处于免打扰时间段，不允许推送
	PushBizForbiddenErr     = New(72005) // 业务方被禁止推送
	PushUploadInvalidErr    = New(72006) // 上传的文件内容不符合规范
	PushRecordRepeatErr     = New(72007) // 该记录已经存在
	PushServiceBusyErr      = New(72008) // 系统繁忙，请稍后再试
	PushServiceFileSizeErr  = New(72009) // 图片大小超过限制
	PushServiceFileExtErr   = New(72010) // 图片格式不支持
	PushAdminDPNoDataErr    = New(72101) // 数据平台数据未准备好
	PushAdminNotPreparedErr = New(72102) // 未预处理完成，不可以上线
	PushAdminPreparingErr   = New(72103) // 预处理中，请勿重复
	PushAdminAllNotAllowErr = New(72104) // 全量不支持预处理

	// realname 实名认证
	RealnameCaptureErr         = New(74001) // 实名验证码输入错误
	RealnameCaptureSendTooMany = New(74002) // 实名验证码发送次数过于频繁
	RealnameCaptureInvalid     = New(74003) // 实名验证码未发送或已失效
	RealnameCaptureErrTooMany  = New(74004) // 实名验证码错误次数过多
	RealnameApplyAlready       = New(74005) // 实名认证已提交申请
	RealnameInvalidName        = New(74006) // 实名姓名错误
	RealnameImageIDErr         = New(74007) // 实名照片错误
	RealnameImageExpired       = New(74008) // 实名认证照片过期
	RealnameCardNumErr         = New(74009) // 实名证件号码错误
	RealnameCardBindAlready    = New(74010) // 实名证件已绑定
	RealnameAlipayAntispam     = New(74011) // 实名触发芝麻认证防刷
	RealnameAlipayErr          = New(74012) // 实名芝麻认证服务错误
	RealnameAlipayApplyInvalid = New(74013) // 实名芝麻认证申请不合法
	// manager
	ManagerTagDelErr     = New(77001) // manager tag不能删除
	ManagerUIDNOTExist   = New(77002) // manager管理员不存在
	ManagerFlowForbidden = New(77003) // manager 工作流被禁用
	ManagerTagTypeDelErr = New(77004) // manager tag类型不能删除

	// subtitle
	SubtitleDrfatExist            = New(79001) // 当前已存在未过审字幕
	SubtitleDrfatNotExist         = New(79002) // 当前字幕草稿不存在
	SubtitleDrfatUnSubmit         = New(79003) // 当前字幕未提交
	SubtitleDelUnExist            = New(79004) // 删除不存在的字幕
	SubtitlePermissionDenied      = New(79005) // 字幕操作权限不足
	SubtitleDenied                = New(79006) // 视频禁止投稿
	SubtileLanLocked              = New(79007) // 当前语言版本锁定
	SubtitleUnValid               = New(79008) // 字幕不合法
	SubtitleWaveFormFailed        = New(79009) // 波形图调用失败
	SubtitlePubNotExist           = New(79010) // 发布的字幕不存在
	SubtitleIllegalLanguage       = New(79011) // 不合法的语言
	SubtitleNotPublish            = New(79012) // 当前字幕未发布
	SubtitleTimeUnValid           = New(79013) // 字幕时间不合法
	SubtitleSizeLimit             = New(79014) // 字幕超过限制
	SubtitleOriginUnValid         = New(79015) // 原字幕不合法
	SubtitleLocationUnValid       = New(79016) // 字幕位置不合法
	SubtitleUserBalcked           = New(79017) // 账号黑名单
	SubtitleStatusUnValid         = New(79018) // 字幕状态不合法
	SubtitleAlreadyHasDraft       = New(79019) // 前存在草稿或者待审核状态的字幕
	SubtitleVideoDurationOverFlow = New(79020) // 字幕时间点超过视频时间长度
	SubtitleDuarionMustThanZero   = New(79021) // 字幕的持续时间必须大于0
	SubtitleVideoNotExist         = New(79022) // 当前视频不存在
	SubtitleDmPostBalcked         = New(79023) // 弹幕发送黑名单
	// MCN
	// --82000~82499 前台错误
	MCNNotAllowed                        = New(82001) // 没有权限操作
	MCNUpCannotBind                      = New(82002) // 无法绑定Up主，已被绑定
	MCNUpBindUpAlreadyInProgress         = New(82003) // 已存在与up的绑定在审核中
	MCNUpBindUpDuplicatePhone            = New(82004) // 该电话号码已绑定
	MCNUpBindUpDuplicateIDCard           = New(82005) // 该身份证号码已绑定
	MCNUpBindUpDuplicateCompanyLicenseID = New(82006) // 该营业执照号码已绑定
	MCNUpBindUpDuplicateCompanyName      = New(82007) // 该公司名称已绑定
	MCNUpBindUpSTimeLtETime              = New(82008) // up主绑定的开始时间必须小于结束时间
	MCNUpBindUpIsBlocked                 = New(82009) // up主已被封禁
	MCNUpBindUpDateError                 = New(82010) // 日期错误，请检查
	MCNStateInvalid                      = New(82011) // MCN状态异常
	MCNUpBindInvalid                     = New(82012) // 该绑定请求已失效
	MCNUpBindInvalidURL                  = New(82013) // 该绑定的站外up主链接错误，请检查
	MCNUpOutSiteIsNotQualified           = New(82014) // 站外Up主需要满足（1.粉丝数≤100  或  2. 投稿数＜2及90天内未投稿）
	MCNUpBindUpIsBlueUser                = New(82015) // 该up主为蓝V用户，不能被绑定
	MCNUpSignStateInvalid                = New(82016) // 该Up主签约状态异常
	MCNChangePermissionAlreadyInProgress = New(82020) // 等待UP主授权或运营审核
	MCNChangePermissionLackPermission    = New(82021) // 您缺少要申请的权限
	MCNChangePermissionSamePermission    = New(82022) // 权限没有任何变更
	MCNPublicationFailTimeLimit          = New(82030) // 刊例价每月只能更改一次
	MCNPermissionInsufficient            = New(82040) // 权限不足
	McnUpNotFound                        = New(82041) // 未找到mcn绑定up主记录
	McnUpBusinessOperateIllegal          = New(82042) // mcn操作up主商单权限不合法
	McnUpAddPermitAlreadyInProgress      = New(82043) // up主当前有权限变更中，暂时不能取消授权哦
	McnUpRemovePermitAlreadyInProgress   = New(82044) // 您添加的up主有部分权限变更中，请刷新页面后再操作哦

	// --82500~82999 后台错误
	MCNSignCycleNotUQErr        = New(82501) // mcn签约周期冲突
	MCNUnknownFileTypeErr       = New(82502) // 只允许jpg、png、pdf、word格式的文件上传
	MCNSignUnknownReviewErr     = New(82503) // mcn签约未知审核方式
	MCNSignOnlyReviewOpErr      = New(82504) // 只有待审核中的mcn才能操作
	MCNUpUnknownReviewErr       = New(82505) // mcn和up主绑定未知审核方式
	MCNUpOnlyReviewOpErr        = New(82506) // 有待审核中的mcn和up主绑定才能操作
	MCNContractFileSize         = New(82507) // mcn合同上传的大小，允许20M
	MCNCSignUnknownInfoErr      = New(82508) // 未知mcn信息
	MCNCUpUnknownInfoErr        = New(82509) // 未知mcn-up主绑定信息
	MCNUnknownFileExt           = New(82510) // 未知上传文件后缀名
	MCNSignNoOkState            = New(82511) // 处于封禁中、审核中、未申请中、驳回中、签约中、待开启中的状态不可录入
	MCNUpPassOnEffectSign       = New(82512) // up主通过必须是有效的签约状态
	MCNSignIsBlocked            = New(82513) // mcn管理用户已被封禁
	MCNSignEtimeNLEQNowTime     = New(82514) // mcn签约结束时间不能小于等于当前时间
	MCNSignStateFlowErr         = New(82515) // mcn状态流转错误,请联系开发人员
	MCNRecommendUpStateFlowErr  = New(82516) // 推荐的up状态流转错误,请联系开发人员
	MCNRecommendUpInPool        = New(82517) // 该up主已经被推荐
	MCNRecommendUpMidsIsEmpty   = New(82518) // 推荐池被执行的up主不能为空
	MCNUpPermitUnknownReviewErr = New(82519) // 未知mcn-up的权限变更审核方式
	MCNUpAbnormalDataErr        = New(82520) // 异常数据
	MCNUpPermitStateFlowErr     = New(82521) // mcn-up的权限变更的状态流转错误,请联系开发人员

	MCNRenewalNotInRangeErr        = New(82601) // MCN续约不在范围内
	MCNRenewalAlreadyErr           = New(82602) // MCN已续过约
	MCNRenewalDateErr              = New(82603) // MCN续约周期错误
	MCNUPRewardMarkErr             = New(82610) // 未在签约状态下的mcn-up、mcn不在签约状态不可标记发奖
	MCNUPFansMarkUnMeetErr         = New(82611) // 未达标的up主不能标记奖励
	MCNUPFansMarkAlreadyMarkErr    = New(82612) // up主已经标记过该奖励，不能再次标记
	MCNUANotExists                 = New(82620) // 协议不存在
	MCNUPFileOrUpMidsMustExistsErr = New(82621) // csv文件或者up mid至少存在一个
	MCNUPBindMidsMustSmallErr      = New(82622) // 批量的绑定up主必须小于1000

	//passport
	//passport-login 86000~86299
	PassportLoginRsaDecryptErr = New(86000) //rsa解密错误
	//passport-user  86300~86599
	//passport-sns	 86600~
	PassportSnsMidAlreadyBindQQ    = New(86600) //mid已绑定QQ号
	PassportSnsMidAlreadyBindWEIBO = New(86601) //mid已绑定微博
	PassportSnsQQAlreadyBind       = New(86610) //QQ号已被绑定
	PassportSnsWEIBOAlreadyBind    = New(86611) //微博已经绑定
	PassportSnsRequestErr          = New(86660) //请求第三方失败

	// passport-auth 86700
	PassportAuthInvalidTmpCode = New(86701) // 无效的tmpCode

	// safe-center 86800
	SafeCenterWrongTelFormat      = New(86801) // 手机号格式错误
	SafeCenterTelUsed             = New(86802) // 该手机号已被占用
	SafeCenterUserNotBindTel      = New(86803) // 该用户未绑定手机
	SafeCenterCaptchaHasSend      = New(86804) // 验证码已发送
	SafeCenterWrongSmsCaptcha     = New(86805) // 短信验证码错误
	SafeCenterEmptySmsCaptcha     = New(86806) // 短信验证码为空
	SafeCenterExpiredSmsCaptcha   = New(86807) // 短信验证码已过期
	SafeCenterEmptyCid            = New(86808) // 国家码为空
	SafeCenterEmptyTel            = New(86809) // 手机号为空
	SafeCenterSmsSendToMuch       = New(86810) // 短信发送次数已达上限
	SafeCenterUserHasBindTel      = New(86811) // 该用户已绑定手机
	SafeCenterInvalidTel          = New(86812) // 无效的手机号
	SafeCenterUserNotBindEmail    = New(86813) // 该用户未绑定邮箱
	SafeCenterEmptyEmail          = New(86814) // 邮箱为空
	SafeCenterUnverifiedEmail     = New(86815) // 邮箱未验证
	SafeCenterWrongEmailCaptcha   = New(86816) // 邮箱验证码错误
	SafeCenterExpiredEmailCaptcha = New(86817) // 邮箱验证码已过期
	SafeCenterEmptyEmailCaptcha   = New(86818) // 邮箱验证码为空
	SafeCenterEmailSendToMuch     = New(86819) // 邮件发送次数已达上限
	SafeCenterEmailSendErr        = New(86820) // 邮件发送错误
	SafeCenterEmptyVerifyType     = New(86821) // 请传入验证时的邮箱类型或者验证码
	SafeCenterExpiredIdentify     = New(86822) // 身份验证过期
	SafeCenterEmptyCaptcha        = New(86823) // 空的验证码
	SafeCenterExpiredCaptcha      = New(86824) // 验证码已过期
	SafeCenterGeeServerErr        = New(86825) // 极验服务器错误
	SafeCenterGeeValidateErr      = New(86826) // 验证极验服务错误
	SafeCenterSettingHasEnabled   = New(86827) // 此设置已开启，请勿重复操作
	SafeCenterSettingHasDisabled  = New(86828) // 此设置已关闭，请勿重复操作
	SafeCenterDeviceReachLimit    = New(86829) // 设备数已达上线
	SafeCenterForbidDelCurDevice  = New(86830) // 禁止删除本设备
	SafeCenterNeedHaveBindTel     = New(86831) // 该操作需要已绑定手机号
	SafeCenterSettingDisabled     = New(86832) // 该设置未启用
	SafeCenterDeviceExist         = New(86833) // 用户设备已存在，请勿重复添加
	SafeCenterDeviceNotExist      = New(86834) // 用户设备不存在，非法操作
	SafeCenterIllegalOpDevice     = New(86835) // 不合法的操作设备

	// uprating
	UpRatingNoPermission = New(91000) // 用户没有权限访问
	UpRatingNoData       = New(91001) // UP主无评分数据
	UpRatingScoreLimit   = New(91002) // UP主评分未达标
	UpRatingDataErr      = New(91003) // UP主等级体系获取数据错误

	//aegis
	AegisUniqueAlreadyExist = New(92001) //当前流程%s %s 已存在
	AegisTokenNotAssign     = New(92002) //当前绑定令牌%s不是赋值语句
	AegisTokenNotFound      = New(92003) //当前令牌不存在
	AegisTextComputeFail    = New(92004) //当前文本解析失败
	AegisFlowNotFound       = New(92005) //当前节点不存在
	AegisFlowDisabled       = New(92006) //当前节点已被禁用
	AegisFlowNoEnableTran   = New(92007) //当前节点所有下游变化不可用
	AegisFlowNoFromDir      = New(92008) //当前节点上游没有有向线
	AegisFlowBinded         = New(92009) //当前流程节点已被关联，请取消后再禁用
	AegisTranBinded         = New(92010) //当前流程变化已被关联，请取消后再禁用
	AegisTranNotFound       = New(92011) //当前变化不存在
	AegisTranDisabled       = New(92012) //当前变化已被禁用
	AegisTranNoFlow         = New(92013) //当前变化没有可用输出节点
	AegisRunInDiffNet       = New(92014) //当前资源在不同网内运行
	AegisNotRunInFlow       = New(92015) //当前资源不在当前节点上运行
	AegisTextErr            = New(92016) //文本配置错误
	AegisNotRunInRange      = New(92017) //当前资源不在该业务或网内运行
	AegisNotTriggerFlow     = New(92018) //当前资源没有触发当前节点上的变化
	AegisDirOrderConflict   = New(92019) //当前有向线顺序(%s)冲突:%s
	AegisNetErr             = New(92020) //当前网配置错误

	AegisTaskErr         = New(92021) //任务操作错误
	AegisResourceErr     = New(92022) //资源操作失败
	AegisBusinessCfgErr  = New(92023) //业务配置错误
	AegisSearchErr       = New(92024) //搜索接口错误
	AegisDuplicateErr    = New(92025) //资源重复
	AegisReservedCfgErr  = New(92026) //保留字段配置错误
	AegisBusinessSyncErr = New(92027) //业务回调错误
	AegisTaskBusy        = New(92028) //任务调度繁忙，请重试
	AegisTaskRelease     = New(92029) //任务已被释放
	AegisESLimitExceeded = New(92030) //搜索内容超过1w条
	AegisNetBusy         = New(92031) //流程网更新繁忙，请重试

	// up admin special group
	UpGroupExist                        = New(95000) // 特殊分组已存在，请重新选择
	UpGroupNoPermit                     = New(95001) // 没有权限处理此分组
	UpGroupsNotBelongToBusiness         = New(95002) // 所选分组都不属于此业务类型
	UpGroupNotExist                     = New(95003) // 特殊分组不存在
	UpMarketPoolExist                   = New(95100) // 营销号已存在
	UpMarketPoolNotExist                = New(95101) // 营销号不存在
	UpMarketPoolOnlyAllowedToAddConfirm = New(95102) // 营销号只允许被添加到确认池中，请重新选择状态
	UpMarketPoolNotConfirm              = New(95103) // 营销号不是确认状态，不允许操作
	UpMarketPoolActionNotExist          = New(95104) // 此操作行为不存在
	UpMarketPoolTypeNotExist            = New(95105) // 此池子类型不存在
	UpMarketPoolRemoveNoReason          = New(95106) // 从确认池里剔除，请选择理由
	UpMarketPoolNotMisjudge             = New(95107) // 此营销号不是误判不能加入白名单
	// upcrm  auth
	UpcrmAuthUserNotExists        = New(95200)
	UpcrmAuthUserExists           = New(95201)
	UpcrmAuthGroupNotExists       = New(95202)
	UpcrmAuthGroupExists          = New(95203)
	UpcrmAuthUserExistsInGroup    = New(95204)
	UpcrmAuthUserNotExistsInGroup = New(95205)
	UpcrmAuthRoleNotExists        = New(95206)
	UpcrmAuthRoleExists           = New(95207)
	UpcrmAuthUserExistsAsRole     = New(95208)
	UpcrmAuthUserNotExistsAsRole  = New(95209)
	UpcrmAuthPrivilegeExists      = New(95210)
	UpcrmAuthPrivilegeNotExists   = New(95211)
	UpcrmAuthMenuItemNotExists    = New(95212)
	UpcrmAuthMenuItemExists       = New(95213)
	UpcrmAuthMenuItemNeedTitle    = New(95214)
	UpcrmAuthPageMenuItemNeedURL  = New(95215)
	// upcrm contract templ
	UpCrmContractTemplateNotExist = New(95300) // 合同模版不存在
	// upmark
	UpMarkNotExist      = New(95350) // 此up主未与该用户绑定
	UpMarkGroupNotExist = New(95351) // 不存在的标记分组
	UpMarMustContainMID = New(95352) // 输入uid或者上传csv文件

	// account recovery
	AccountRecoveryAppealExisted = New(96000) // 有未处理的账号找回申诉
	// steins gate
	SteinsCidNotMatch    = New(99000) // 节点Cid信息不匹配
	NonValidGraph        = New(99001) // 互动视频：当前稿件无可用剧情图
	NotSteinsGateArc     = New(99002) // 当前稿件不是互动视频
	GraphInvalid         = New(99003) // 剧情图被修改已失效
	GraphAidEmpty        = New(99004) // 剧情图缺少aid参数
	GraphNotOwner        = New(99005) // 你不是该稿件的作者
	GraphAidAttrErr      = New(99006) // 该稿件类型不是互动视频
	GraphScriptEmpty     = New(99007) // 剧情图缺少scpirt数据
	GraphNodeCntErr      = New(99008) // 剧情图节点数错误
	GraphNodeCidEmpty    = New(99009) // 剧情图节点缺少cid
	GraphNodeNameErr     = New(99010) // 剧情图节点名称长度不对
	GraphNodeNameExist   = New(99011) // 剧情图节点名称重复
	GraphDefaultNodeErr  = New(99012) // 剧情图有多个默认节点
	GraphEdgeCntErr      = New(99013) // 剧情图节点分支选项数错误
	GraphLackStartNode   = New(99014) // 剧情图缺少开始节点
	GraphEdgeNameErr     = New(99015) // 剧情图节点名称长度不对
	GraphDefaultEdgeErr  = New(99016) // 剧情图节点多个默认分支选项
	GraphFilterHitErr    = New(99017) // 剧情图有内容命中敏感词
	GraphNodeOtypeErr    = New(99018) // 剧情图节点类型错误
	GraphNodeCircle      = New(99019) // 剧情图节点有回环结构
	GraphEdgeToNodeErr   = New(99020) // 剧情图节点分支无到达节点
	GraphFilterErr       = New(99021) // 请求过滤词服务失败
	GraphShowTimeEdgeErr = New(99022) // 剧情图直连只支持一个选项
	GraphArcStateErr     = New(99023) // 该稿件暂未过审，请耐心等待稿件过审后再提交
	GraphPageWidthErr    = New(99024) // 暂不支持竖屏视频，请更换含有【%s】的剧情后再进行提交

	// SilverBullet [100000,105000]
	// recaptcha [100000,100200]
	RecaptchaRegisterFailed     = New(100000) // 验证码获取失败
	RecaptchaValidateFailed     = New(100001) // 验证码校验失败
	RecaptchaCodeErr            = New(100002) // 验证码输入错误
	RecaptchaExpiredErr         = New(100003) // 验证码过期
	RecaptchaGeeTestServerErr   = New(100004) // 极验服务器错误
	RecaptchaGeeTestRegisterErr = New(100005) // 注册极验服务错误
	RecaptchaGeeTestValidateErr = New(100006) // 校验极验服务错误
)
