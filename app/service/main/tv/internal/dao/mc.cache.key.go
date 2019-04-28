package dao

import "fmt"

func userInfoKey(mid int64) string {
	return fmt.Sprintf("cache_tv_vip_ui_%d", mid)
}

func payParamKey(token string) string {
	return token
}
