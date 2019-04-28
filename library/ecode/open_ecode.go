package ecode

// common  ecode
// 开放平台 2000000~2999999
// 票务的code码 2000000~2099999
var (
	//销售-营销
	TicketUnKnown               = New(2000000) //未知错误
	TicketParamInvalid          = New(2000001) //参数错误
	TicketRecordDupli           = New(2000002) //重复插入
	TicketRecordLost            = New(2000003) //数据不存在
	TicketPromotionLost         = New(2000004) //活动不存在
	TicketPromotionEnd          = New(2000005) //活动结束
	TicketPromotionRepeatJoin   = New(2000006) //活动重复参加
	TicketPromotionGroupLost    = New(2000007) //拼团不存在
	TicketPromotionGroupFull    = New(2000008) //拼团人数已满
	TicketPromotionGroupNotFull = New(2000009) //拼团人数未满
	TicketPromotionOrderLost    = New(2000010) //拼团订单不存在
	TicketPromoExistSameTime    = New(2000011) //同时间段存在已上架拼团活动
	TicketAddPromoOrderFail     = New(2000012) //添加活动订单失败
	TicketAddPromoGroupFail     = New(2000013) //添加拼团 团订单失败
	TicketPromoGroupEnd         = New(2000014) //拼团 团订单已失效
	TicketUpdatePromoOrderFail  = New(2000015) //更新拼团订单失败
	TicketUpdatePromoGroupFail  = New(2000016) //更新拼团 团订单失败
	IllegalPromoOperate         = New(2000017) //拼团 不支持的操作类型
	PromoStatusChanged          = New(2000018) //拼团状态无法变更
	TicketPromoGroupStatusErr   = New(2000019) //拼团状态不对
	TicketPromoOrderTypeErr     = New(2000020) //订单类型不对
	PromoEditNotALlowed         = New(2000021) //不可编辑
	PromoEditFieldNotALlowed    = New(2000022) //不可编辑部分字段
	PromoExists                 = New(2000023) //拼团已存在

	//销售-交易
	TicketGetOidFail   = New(2000101) //获取订单号失败
	TicketExceedLimit  = New(2000102) //超过购买限制
	TicketMissData     = New(2000103) //信息不完整
	TicketSaleNotStart = New(2000104) //没开售
	TicketSaleEnd      = New(2000105) //已结束
	TicketNoPriv       = New(2000106) //无权操作
	TicketInvalidUser  = New(2000107) //无效用户
	TicketPriceChanged = New(2000108) //价格变化

	//销售-库存
	TicketStockLack        = New(2000201) //库存不足
	TicketStockLogNotFound = New(2000202) //没有库存操作记录
	TicketStockUpdateFail  = New(2000203) //库存更新失败

	//番剧推荐
	SugEsSearchErr   = New(2002000) //es搜索错误
	SugSearchTypeErr = New(2002001) //搜索类型错误
	SugOpTypeErr     = New(2002002) //操作类型错误
	SugOpErr         = New(2002003) //add or del match fail
	SugItemNone      = New(2002004) //商品不存在
	SugSeasonNone    = New(2002005) //番剧不存在

	//防刷工具
	ParamInvalid          = New(2001000) //参数错误
	UpdateError           = New(2001002) //更新失败
	QusbNotFound          = New(2001003) //找不到题库
	QusIDInvalid          = New(2001005) //题目id错误
	BankUsing             = New(2001007) //题目正在使用
	BindBankNotFound      = New(2001009) //未找到题库绑定关系
	AnswerError           = New(2001010) //答案错误
	GetQusBankInfoCache   = New(2001011) //获取题库缓存失败
	GetComponentTimesErr  = New(2001012) //获取组件缓存失败
	SetComponentTimesErr  = New(2001013) //设置答题次数缓存失败
	SetComponentIDErr     = New(2001014) //设置组件缓存失败
	GetComponentIDErr     = New(2001015) //获取组件ID缓存失败
	SameCompentErr        = New(2001016) //相同组件
	GetQusIDsErr          = New(2001017) //获取题目失败
	AnswerPoiError        = New(2001018) //答案错误
	NotEnoughQuestion     = New(2001019) //部分题库不足3题，无法绑定，请绑定别的题库，或者修改题库
	AntiSalesTimeErr      = New(2001020) //售卖时间有错
	AntiIPChangeLimit     = New(2001021) //用户IP变更
	AntiLimitNumUpper     = New(2001022) //次数达到上限
	AntiCheckVoucherErr   = New(2001023) //用户凭证验证失败
	AntiValidateFailed    = New(2001024) //验证失败
	AntiGeetestCountUpper = New(2001025) //极验总数达到上线
	AntiCustomerErr       = New(2001026) //业务方错误
	AntiBlackErr          = New(2001027) //黑名单用户

	//项目
	TicketCannotDelTk      = New(2004000) //无法删除票价
	TicketDelTkFailed      = New(2004001) //删除票价失败
	TicketLkTkNotFound     = New(2004002) //关联票种不存在
	TicketLkTkTypeNotFound = New(2004003) //关联票种类型不存在
	TicketLkScNotFound     = New(2004004) //关联场次不存在
	TicketCannotDelSc      = New(2004005) //无法删除场次
	TicketLkScTimeNotFound = New(2004006) //关联的场次时间不存在
	TicketPidIsEmpty       = New(2004007) //项目id为空
	TicketMainInfoTooLarge = New(2004008) //项目版本详情信息量过大
	TicketDelTkExFailed    = New(2004009) //删除票价额外信息失败
	TicketAddVersionFailed = New(2004010) //添加版本信息失败
	TicketAddVerExtFailed  = New(2004011) //添加版本详情失败
	TicketBannerIDEmpty    = New(2004012) //BannerID为空
	TicketVerCannotEdit    = New(2004013) //版本不可编辑
	TicketVerCannotReview  = New(2004014) //无法审核 非待审核版本
	TicketAddTagFailed     = New(2004015) //添加项目标签失败

)
