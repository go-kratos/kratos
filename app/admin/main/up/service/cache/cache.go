package cache

import "time"

//LoadCache load cache
func LoadCache() {
	var now = time.Now()
	RefreshUpType(now)
}
