package mi

// Xiaomi push service document: https://dev.mi.com/doc/cat=35/index.html

const (
	// VipHost VIP host.
	VipHost = "https://vip.api.xmpush.xiaomi.com"
	// DevHost dev host.
	DevHost = "https://sandbox.xmpush.xiaomi.com"
	// ProductionHost production host.
	ProductionHost = "https://api.xmpush.xiaomi.com"
	// feedbackHost host to get invalid token.
	feedbackHost = "https://feedback.xmpush.xiaomi.com"
	// emqHost message queue
	emqHost = "https://emq.xmpush.xiaomi.com"
	// AuthPrefix auth prefix.
	AuthPrefix = "key="
	// ResultOk result status.
	ResultOk = "ok" // "ok" means success, "error" means failed.
	// ResultError result status.
	ResultError = "error"
	// ResultCodeOk result status code.
	ResultCodeOk = 0
	// ResultCodeNoValidTargets no valid token.
	ResultCodeNoValidTargets = 20301
	// ResultCodeNoMsgInEmq no message in emq.
	ResultCodeNoMsgInEmq = 80002

	// RegURL 向某个regid或一组regid列表推送某条消息
	RegURL = "/v3/message/regid"
	// AccountURL 根据account，发送消息到指定account上
	AccountURL = "/v2/message/user_account"
	// MultiRegIDURL  针对不同的regid推送不同的消息
	MultiRegIDURL = "/v2/multi_messages/regids"
	// MultiAliasURL 针对不同的aliases推送不同的消息
	MultiAliasURL = "/v2/multi_messages/aliases"
	// MultiUserAccountURL 针对不同的accounts推送不同的消息
	MultiUserAccountURL = "/v2/multi_messages/user_accounts"
	// AliasURL 根据alias，发送消息到指定设备上
	AliasURL = "/v3/message/alias"
	// MultiPackageNameMultiTopicURL 根据topic，发送消息到指定一组设备上
	MultiPackageNameMultiTopicURL = "/v3/message/multi_topic"
	// MultiTopicURL 根据topic，发送消息到指定一组设备上
	MultiTopicURL = "/v2/message/topic"
	// MultiPackageNameAllURL 向所有设备推送某条消息
	MultiPackageNameAllURL = "/v3/message/all"
	// AllURL 向所有设备推送某条消息
	AllURL = "/v2/message/all"
	// TopicURL 向多个topic广播消息
	TopicURL = "/v3/message/multi_topic"
	// ScheduleJobExistURL 检测定时消息的任务是否存在
	ScheduleJobExistURL = "/v2/schedule_job/exist"
	// ScheduleJobDeleteURL 删除指定的定时消息
	ScheduleJobDeleteURL = "/v2/schedule_job/delete"
	// ScheduleJobDeleteByJobKeyURL 删除指定的定时消息
	ScheduleJobDeleteByJobKeyURL = "/v3/schedule_job/delete"
	// feedbackURI 获取无效token列表
	feedbackURI = "/v1/feedback/fetch_invalid_regids"
	// uninstalledURI 获取卸载token列表
	uninstalledURI = "/app/uninstall/regid"
	// StatusURL 追踪消息
	StatusURL = "/v1/trace/message/status"

	// NotifyTypeDefaultAll 包括下面三种(notify type 可以是以下几种的OR组合)
	NotifyTypeDefaultAll = -1
	// NotifyTypeDefaultNone 声音、振动、led灯全关
	NotifyTypeDefaultNone = 0
	// NotifyTypeDefaultSound 使用默认提示音提示
	NotifyTypeDefaultSound = 1
	// NotifyTypeDefaultVibration 使用默认震动提示
	NotifyTypeDefaultVibration = 2
	// NotifyTypeDefaultLight 使用默认led灯光提示
	NotifyTypeDefaultLight = 4

	// NotPassThrough 显示通知
	NotPassThrough = 0
	// PassThrough 静默推送
	PassThrough = 1

	// CallbackURL 客户端收到后回调
	CallbackURL = "https://api.bilibili.com/x/push/callback/xiaomi"
	// CallbackBarStatusEnable .
	CallbackBarStatusEnable = 1
	// CallbackBarStatusDisable .
	CallbackBarStatusDisable = 2
	// CallbackBarStatusUnknown .
	CallbackBarStatusUnknown = 3
	// CallbackBarStatusEnableStr .
	CallbackBarStatusEnableStr = "Enable"
	// CallbackBarStatusDisableStr .
	CallbackBarStatusDisableStr = "Disable"
	// CallbackBarStatusUnknownStr .
	CallbackBarStatusUnknownStr = "Unknown"
)
