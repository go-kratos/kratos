package http

import (
	"fmt"

	bm "go-common/library/net/http/blademaster"
)

// searchLogAudit .
func searchLogAudit(c *bm.Context) {
	var (
		err    error
		result []byte
	)
	if result, err = mngSvc.SearchLogAudit(c); err != nil {
		c.Bytes(-500, "application/csv", nil)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", "Authority_Record"))
	c.Writer.Write([]byte("\xEF\xBB\xBF"))
	c.Writer.Write(result)
}

// searchLogAction .
func searchLogAction(c *bm.Context) {
	var (
		err    error
		result []byte
	)
	if result, err = mngSvc.SearchLogAction(c); err != nil {
		c.Bytes(-500, "application/csv", nil)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", "Authority_Record"))
	c.Writer.Write([]byte("\xEF\xBB\xBF"))
	c.Writer.Write(result)
}
