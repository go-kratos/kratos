package http

import (
	"fmt"

	bm "go-common/library/net/http/blademaster"
)

func outBuf(c *bm.Context, buf []byte, name string) {
	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", name))
	c.Writer.Write(buf)

}
