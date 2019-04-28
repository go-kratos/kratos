package http

import (
	bm "go-common/library/net/http/blademaster"
)

func queryMachineLifeCycle(c *bm.Context) {
	c.JSON(svc.QueryMachineLifeCycle(c))
}

func queryMachineCount(c *bm.Context) {
	c.JSON(svc.QueryMachineCount(c))
}

func queryMachineTime(c *bm.Context) {
	c.JSON(svc.QueryMachineCreatedAndEndTime(c))
}

func queryMachineUsage(c *bm.Context) {
	c.JSON(svc.QueryMachineUsage(c))
}

func queryMobileMachineUsageCount(c *bm.Context) {
	c.JSON(svc.QueryMobileMachineUsageCount(c))
}

func queryMobileMachineModeCount(c *bm.Context) {
	c.JSON(svc.QueryMobileMachineModeCount(c))
}

func queryMobileMachineUsageTime(c *bm.Context) {
	c.JSON(svc.QueryMobileMachineUsageTime(c))
}
