package dao

import (
	"fmt"
)

// AnniversaryKey format anniversary key
func AnniversaryKey(mid int64) string {
	return fmt.Sprintf("art_anniversary_%d", mid)
}
