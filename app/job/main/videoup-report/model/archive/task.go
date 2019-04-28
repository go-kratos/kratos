package archive

import (
	"sync"
	"time"
)

var (
	// TookTypeMinute video task took time in 1 minute
	TookTypeMinute = int8(1)
	// TookTypeHalfHour video task took time in 10 minutes
	TookTypeHalfHour = int8(2)
	// TaskStateUnclaimed video task belongs to nobody
	TaskStateUnclaimed = int8(0)
	// TaskStateUntreated video task not submit
	TaskStateUntreated = int8(1)
	// TaskStateCompleted video task completed
	TaskStateCompleted = int8(2)
	// TaskStateDelayed video task delayed
	TaskStateDelayed = int8(3)
	// TaskStateClosed video task closed
	TaskStateClosed = int8(4)
)

// TaskCache store task video
type TaskCache struct {
	Task  map[int64]*Task
	Took  []*TaskTook
	Sort  []int
	Mtime time.Time
	sync.Mutex
}

// Task video task entity
type Task struct {
	ID      int64     `json:"id"`
	Subject int8      `json:"subject"`
	Adminid int64     `json:"adminid"`
	Pool    int8      `json:"pool"`
	Aid     int64     `json:"aid"`
	Cid     int64     `json:"cid"`
	State   int8      `json:"state"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"-"`
}

// TaskTook video task take time
type TaskTook struct {
	ID     int64     `json:"id"`
	M90    int       `json:"m90"`
	M80    int       `json:"m80"`
	M60    int       `json:"m60"`
	M50    int       `json:"m50"`
	TypeID int8      `json:"type"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"-"`
}
