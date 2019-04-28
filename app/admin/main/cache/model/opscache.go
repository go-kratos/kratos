package model

// OpsCacheMemcache ops cache mc
type OpsCacheMemcache struct {
	Labels struct {
		Name    string `json:"name"`
		Project string `json:"project"`
	}
	Targets []string `json:"targets"`
}

// OpsCacheRedis ops cache redis
type OpsCacheRedis struct {
	Labels struct {
		Name    string `json:"name"`
		Project string `json:"project"`
	}
	Type    string   `json:"type"`
	Targets []string `json:"master_targets"`
}
