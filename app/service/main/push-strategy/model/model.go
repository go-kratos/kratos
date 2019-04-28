package model

import (
	pushmdl "go-common/app/service/main/push/model"
)

// MidChan .
type MidChan struct {
	Task *pushmdl.Task
	Data *string
}
