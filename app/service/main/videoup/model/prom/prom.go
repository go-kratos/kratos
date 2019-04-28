package prom

import (
	"fmt"
)

const (
	RouteNormalMode = "normal_mode"
	RouteCodeMode   = "code_mode"
	RouteUpFrom     = "upfrom_"

	RouteNumFormatVids = "%s_%dp"
	RouteStrFormatVids = "%s_%sp"

	RouteDmIndexTry = "dm_index_retry"
	RouteDatabusTry = "databus_retry"
)

// FormatVideoKey format videos prom key
func FormatVideoKey(count int, tp string) (key string) {
	if count >= 0 && count <= 5 {
		key = fmt.Sprintf(RouteNumFormatVids, tp, count)
	} else if count >= 6 && count <= 10 {
		key = fmt.Sprintf(RouteStrFormatVids, tp, "6_10")
	} else {
		key = fmt.Sprintf(RouteStrFormatVids, tp, "above_10")
	}
	return
}
