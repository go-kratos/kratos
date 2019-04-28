package model

import (
	"go-common/app/service/main/member/model/block"
)

// BlockResult is
type BlockResult struct {
	MID         int64             `json:"mid"`
	BlockStatus block.BlockStatus `json:"block_status"`
	StartTime   int64             `json:"start_time"`
	EndTime     int64             `json:"end_time"`
}
