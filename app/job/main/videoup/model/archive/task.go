package archive

const (
	// PoolForFirst 一审
	PoolForFirst = int8(0)
	// PoolForSecond 二审
	PoolForSecond = int8(1)
	// SubjectForNormal 普通任务
	SubjectForNormal = int8(0)
	// SubjectForTask 指派任务
	SubjectForTask = int8(1)
	// StateForTaskDefault 初始化状态（未认领）
	StateForTaskDefault = int8(0)
	// StateForTaskWork 已认领，未处理
	StateForTaskWork = int8(1)
	// StateForTaskDelay 延迟审核
	StateForTaskDelay = int8(3)
	// StateForTaskUserDeleted 被释放
	StateForTaskUserDeleted = int8(6)
)
