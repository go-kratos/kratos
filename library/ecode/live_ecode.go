package ecode

//  直播　号段　　1000000 - 1999999
//  为防止定义重复，以及方便查找，各个服务错误码可直接选用一个段号
//  调用服务出错已选用一个段号，参见最后
var (
	// wallet 1000000 - 1001999
	CoinNotEnough = New(1000000)
	PayFailed     = New(1000001)

	// live-test 100 2000 - 100 2999
	RoomNotFound = New(1002001)

	InvalidParam                = New(1002002)
	GetAllListRPCError          = New(1002003)
	GetAllListReturnError       = New(1002004)
	GetAllListSimpleJSONError   = New(1002005)
	AttentionRPCError           = New(1002006)
	AttentionReturnError        = New(1002007)
	UserTagRPCError             = New(1002008)
	UserTagReturnError          = New(1002009)
	UserTagRoomListRPCError     = New(1002010)
	UserTagRoomListReturnError  = New(1002011)
	SeaPatrolRPCError           = New(1002012)
	SeaPatrolReturnError        = New(1002013)
	ActivityRPCError            = New(1002014)
	ActivityReturnError         = New(1002015)
	ChangeGetAllListRPCError    = New(1002016)
	ChangeGetAllListReturnError = New(1002017)
	ChangeGetAllListEmptyError  = New(1002018)

	SkyHorseError            = New(1002019)
	ChangeSkyHorseEmptyError = New(1002020)
	GetRoomError             = New(1002021)
	GetRoomEmptyError        = New(1002024)

	RoomPendantError                   = New(1002022)
	RoomPendantReturnError             = New(1002023)
	AttentionListRPCError              = New(1002100)
	RelationFrameWorkCallError         = New(1002101)
	RelationLiveRPCCodeError           = New(1002102)
	RelationFrameWorkGoRoutingError    = New(1002103)
	RoomGetStatusInfoRPCError          = New(1002104)
	GetMultipleRPCError                = New(1002105)
	RoomPendentRPCError                = New(1002106)
	PKIDRpcError                       = New(1002107)
	UnliveAnchorReqParamsError         = New(1002108)
	LiveAnchorReqParamsError           = New(1002109)
	NeedLogIn                          = New(1002110)
	RoomFrameWorkCallError             = New(1002111)
	RoomLiveRPCCodeError               = New(1002112)
	RoomFrameWorkGoRoutingError        = New(1002113)
	UserFrameWorkCallError             = New(1002114)
	UserLiveRPCCodeError               = New(1002115)
	UserFrameWorkGoRoutingError        = New(1002116)
	RecordRecordFrameWorkCallError     = New(1002117)
	RecordLiveRPCCodeError             = New(1002118)
	RecordFrameWorkGoRoutingError      = New(1002119)
	RoomNewsRecordFrameWorkCallError   = New(1002120)
	RoomNewsLiveRPCCodeError           = New(1002121)
	RoomNewsFrameWorkGoRoutingError    = New(1002122)
	RoomPendentFrameWorkCallError      = New(1002123)
	RoomPendentLiveRPCCodeError        = New(1002124)
	RoomPendentFrameWorkGoRoutingError = New(1002125)
	PkIDRecordFrameWorkCallError       = New(1002126)
	PkIDLiveRPCCodeError               = New(1002127)
	PkIDFrameWorkGoRoutingError        = New(1002128)
	FansMedalFrameWorkCallError        = New(1002129)
	FansMedalLiveRPCCodeError          = New(1002130)
	FansMedalFrameWorkGoRoutingError   = New(1002131)
	GiftFrameWorkCallError             = New(1002132)
	GiftLiveRPCCodeError               = New(1002133)
	RoomGetRoomIDCodeRPCError          = New(1002134)
	LiveAnchorReqParamsNil             = New(1002135)
	GetGrayRuleError                   = New(1002136)
	AccountGRPCError                   = New(1002137)
	AccountGRPCFrameError              = New(1002138)
	LiveAnchorReqV2ParamsNil           = New(1002139)
	LiveAnchorReqV2ParamsError         = New(1002140)
	UserDHHRPCError                    = New(1002141)
	UserDHHReturnError                 = New(1002142)
	UserDHHDataNil                     = New(1002143)
	XanchorGRPCError                   = New(1002150)

	// resource 1003000 - 1003199
	ResourceParamErr = New(1003001) // 参数传入错误
	TimeForErr       = New(1003002) // 时间格式错误
	AddResourceErr   = New(1003003) // 添加资源失败
	RepdAddErr       = New(1003004) // 添加平台已存在
	SeltResErr       = New(1003005) // 资源选择失败
	EditResErr       = New(1003006) // 编辑资源失败
	DeviceError      = New(1003007) // 参数<device>传入错误
	OfflineResErr    = New(1003008) // 下线资源失败
	GetListResErr    = New(1003009) // 获取资源列表失败
	GetBannerErr     = New(1003010) // 获取Banner配置失败
	GetSplashErr     = New(1003011) // 获取闪屏配置失败
	CheckURLErr      = New(1003012) // 链接格式错误
	GetConfAdminErr  = New(1003101) // 没有获取到配置
	SetConfAdminErr  = New(1003102) // 设置配置失败

	// dm 1003200 - 1003399
	DMallUser      = New(1003200) // 系统正在维护(全员弹幕禁言)
	DMUserLevel    = New(1003201) // 系统正在维护(全员指定等级禁言)
	RealName       = New(1003202) // 实名认证才可以发言
	PhoneBind      = New(1003203) // 根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作。
	PhoneReal      = New(1003204) // 根据国家实名制认证的相关要求，您需要换绑一个非170/171的手机号，才能继续进行操作。
	ShieldUser     = New(1003205) // u (被播主过滤的用户)
	ShieldContent  = New(1003206) // k (被播主过滤的内容)
	BlockUser      = New(1003207) // 你在本房间被禁言
	PayLive        = New(1003208) // 非常抱歉，本场直播需要购票，即可参与互动(付费直播)
	SecDMLimit     = New(1003209) // msg in 1s(每秒发言限制)
	DMSameMsgLimit = New(1003210) // msg repeat(消息重复)
	DMLimitPerRoom = New(1003211) // max limit(单房间每秒限制)
	MsgLengthLimit = New(1003212) // 超出限制长度
	FilterLimit    = New(1003213) // 内容非法（敏感词限制）
	CountryLimit   = New(1003214) // 你所在的地区暂无法发言(区域限制)
	RoomLeverLimit = New(1003215) // 房间等级限制
	RoomAllLimit   = New(1003216) // 房间全员限制
	RoomMedalLimit = New(1003217) // 房间勋章等级限制
	DMServiceERR   = New(1003218) // 依赖下游服务失败

	// user 1004000 - 1004999
	UidError     = New(1004000) //Uid错误
	UserNotFound = New(1004001) //找不到用户信息

	// dao-anchor  1005000 - 1005999
	DaoAnchorCheckAttrSubIdERROR = New(1005000) //找不到对应的标签子id

	// user_ex    1006000 - 1006999
	// gift       1007000 - 1007999
	// xuser 1008000 - 1008999
	XUserAddUserExpReqBizNotAllow   = New(1008001)
	XUserAddUserExpTypeNotAllow     = New(1008002)
	XUserAddUserExpNumNotAllow      = New(1008003)
	XUserAddUserExpUpdateDBError    = New(1008004)
	XUserAddUserExpParamsEmptyError = New(1008005)
	XUserAddUserExpQueryAfterFail   = New(1008006)
	XUserAddAnchorUpdateDBError     = New(1008007)
	XUserAddRExpUpdateDBError       = New(1008008)
	XUserAddUserRExpQueryAfterFail  = New(1008009)
	XUserAddRoomAdminOverLimitError = New(1008020)
	XUserAddRoomAdminIsAdminError   = New(1008021)
	XUserAddRoomAdminIsSilentError  = New(1008022)
	XUserAddRoomAdminNotAdminError  = New(1008023)

	XUserExpGetExpMcFail     = New(1008024) // 获取用户经验缓存失败
	XUserExpGetExpDBFail     = New(1008025) // 获取用户经验缓回源db失败
	XUserAddUExpGetParamsNil = New(1008026) // 添加用户经验缓入参为空

	XUserGuardFetchRecentTopListFail = New(1008027) // 获取主播最近总督失败

	// capusle 1009000 - 1008999
	XLotteryCapsuleAreaParamErr      = New(1009001)
	XLotteryCapsuleSystemErr         = New(1009002)
	XLotteryCapsuleCoinNotEnough     = New(1009003)
	XLotteryCapsuleCoinNotChange     = New(1009004)
	XLotteryCapsulePoolNotChange     = New(1009005)
	XLotteryCapsulePoolNotOffline    = New(1009006)
	XLotteryCapsuleOperationFrequent = New(1009007)

	// app-conf 1010000 - 1010100
	AppConfKeyErr = New(1010000)

	//验证码      1990000 - 1990100
	VerifyNeed = New(1990000)
	VerifyErr  = New(1990001)

	//网关层错误码 1100000 - 1101999
	FILTERNOTPASS     = New(1100000) //屏蔽词校验失败
	UploadTokenGenErr = New(1100010) // Upload token 创建失败
	UploadUploadErr   = New(1100011) // 上传失败
	UploadBucketErr   = New(1100012) // 上传bucket错误

	//调用服务出错 19001000 - 19001999
	CallRoomError       = New(19001000)
	CallUserError       = New(19001001)
	CallRelationError   = New(19001002)
	CallFansMedalError  = New(19001003)
	CallMainMemberError = New(19001004)
	CallResourceError   = New(19001005)
	CallMainFilterError = New(19001006)
	CallDaoAnchorError  = New(19001007)
)
