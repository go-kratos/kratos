package model

import "time"

// const var
const (
	TaskStateSearch = int32(2) //查询执行中
	TaskStateDelDM  = int32(3) //删除执行中
	TaskStateFail   = int32(4) //执行失败
	TaskStatePause  = int32(5) //任务中断
	TaskStateSuc    = int32(6) //执行成功
	TaskStateDel    = int32(8) //任务被删除
	TaskStateWait   = int32(9) //等待执行删除

	// 数据平台返回的弹幕任务查询状态
	TaskSearchSuc  = int32(1) // 查询完成
	TaskSearchFail = int32(2) // 查询失败

	// 企业微信通知
	TaskNoticeTitle   = "弹幕任务删除过多告警"
	TaskNoticeContent = "弹幕任务(id:%d, title:%s)已删除%d条弹幕，已经被暂停，请前往管理后台查看"
)

// TaskInfo .
type TaskInfo struct {
	ID        int64
	Topic     string
	State     int32
	Count     int64
	Result    string
	Sub       int32
	LastIndex int32
	Priority  int64
	Title     string
	Creator   string
	Reviewer  string
}

// SubTask .
type SubTask struct {
	ID        int64
	Operation int32
	Rate      int32
	Tcount    int64 //删除总数
	Start     time.Time
	End       time.Time
}
