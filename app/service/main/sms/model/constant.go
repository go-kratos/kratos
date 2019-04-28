package model

//go:generate $GOPATH/src/go-common/app/tool/warden/protoc.sh

const (
	// CountryChina .
	CountryChina = "86" // 中国地区码
)

// 短信模板状态
const (
	// TemplateStatusNew .
	TemplateStatusNew = 0 // 新建
	// TemplateStatusApprovel .
	TemplateStatusApprovel = 1 // 已审核
	// TemplateStatusCanceled .
	TemplateStatusCanceled = 2 // 已取消
)

// 运营商
const (
	// ProviderMengWang .
	ProviderMengWang = int32(2)
	// ProviderChuangLan
	ProviderChuangLan = int32(4)
)

// 短信类型
const (
	// TypeSms 验证码
	TypeSms = int32(1)
	// TypeActSms 营销短信
	TypeActSms = int32(2)
	// TypeActBatch 批量营销
	TypeActBatch = int32(3)
)

// 短信回执用户行为日志
const (
	// UserActionTypeSend 日志类型为发送
	UserActionTypeSend = int32(1)
	// UserActionCallback 日志类型为回执
	UserActionCallback = int32(2)

	// UserActionSendFailedStatus 发送失败的日志状态
	UserActionSendFailedStatus = "SUBMIT FAILED"
	// UserActionSendFailedDesc 发送失败的日志描述
	UserActionSendFailedDesc = "短信提交失败"
	// UserActionSendSuccessStatus 提交成功的日志状态
	UserActionSendSuccessStatus = "SUBMIT SUCCESS"
	// UserActionSendSuccessDesc 提交失败的日志描述
	UserActionSendSuccessDesc = "短信提交成功"
)

// 短信状态
const (
	// StatusNew .
	StatusNew = int32(0)
	// StatusSuccess .
	StatusSuccess = int32(1)
	// StatusFail .
	StatusFail = int32(2)
)

// 短信任务类型
const (
	// TaskTypeMobile 按手机号发送
	TaskTypeMobile = int32(1)
	// TaskTypeMid 按mid发送
	TaskTypeMid = int32(2)
)

// 短信任务状态
const (
	// TaskStatusPrepared 准备发
	TaskStatusPrepared = int32(1)
	// TaskStatusDoing 进行中
	TaskStatusDoing = int32(2)
	// TaskStatusSuccess 发送成功
	TaskStatusSuccess = int32(3)
	// TaskStatusFailed 发送失败
	TaskStatusFailed = int32(4)
	// TaskStatusStop 停止发送
	TaskStatusStop = int32(5)
)
