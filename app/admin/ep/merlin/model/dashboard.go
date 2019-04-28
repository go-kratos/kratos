package model

import "time"

// MachineLifeCycle Machine Life Cycle.
type MachineLifeCycle struct {
	Duration string `json:"duration_day"`
	Count    int    `json:"count"`
}

// MachineCountGroupByBusiness Machine Count Group By App.
type MachineCountGroupByBusiness struct {
	BusinessUnit string
	Count        int
}

// MachineCountGroupResponse Machine Count Group Response.
type MachineCountGroupResponse struct {
	BusinessUnits []string                 `json:"departmentList"`
	Items         []*MachineCountGroupItem `json:"items"`
}

// MachineCountGroupItem Machine Count Group Item.
type MachineCountGroupItem struct {
	Type string `json:"type"`
	Data []int  `json:"data"`
}

// MachineCreatedAndEndTime Machine Created And End Time.
type MachineCreatedAndEndTime struct {
	ID          int64  `json:"id"`
	MachineName string `json:"machine_name"`
	App         string `json:"app"`
	Username    string `json:"username"`
	CreateTime  string `json:"created_time"`
	EndTime     string `json:"end_time"`
}

// MachineUsage Machine Usage.
type MachineUsage struct {
	ID            int64  `json:"id"`
	MachineName   string `json:"machine_name"`
	App           string `json:"app"`
	Username      string `json:"username"`
	CPURequest    int    `json:"cpu_request"`
	MemoryRequest int    `json:"memory_request"`
}

// MobileMachineUserUsageCount Mobile Machine user Usage Count.
type MobileMachineUserUsageCount struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}

// MobileMachineUsageCount Mobile Machine Usage Count.
type MobileMachineUsageCount struct {
	MobileMachineID   string `json:"mobile_machine_id"`
	MobileMachineName string `json:"mobile_machine_name"`
	Count             int    `json:"count"`
}

// MobileMachineTypeCount Mobile Machine Type Count.
type MobileMachineTypeCount struct {
	ModeName string `json:"mode_name"`
	Count    int    `json:"count"`
}

// MobileMachineUsageTime Mobile Machine Usage Time.
type MobileMachineUsageTime struct {
	//MobileMachineID   int64     `json:"mobile_machine_id"`
	//MobileMachineName string    `json:"mobile_machine_name"`
	//ModeName          string    `json:"mode_name"`
	Username  string    `json:"username"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  float64   `json:"duration_minutes"`
}

// MobileMachineUsageTimeResponse Mobile Machine Usage Time Response.
type MobileMachineUsageTimeResponse struct {
	MobileMachineID         int64                     `json:"mobile_machine_id"`
	MobileMachineName       string                    `json:"mobile_machine_name"`
	ModeName                string                    `json:"mode_name"`
	TotalDuration           float64                   `json:"total_duration_minutes"`
	MobileMachinesUsageTime []*MobileMachineUsageTime `json:"machine_usage_record"`
}
