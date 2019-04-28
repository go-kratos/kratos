package lancerroute

type Lancerroute struct {
	c *Config
}

var route *Lancerroute

func InitLancerRoute() error {
	route = new(Lancerroute)
	return route.InitConfig()
}

func GetLancerByLogid(logId string) string {
	if d, ok := route.c.LancerRoute[logId]; ok {
		return d
	}
	// lancer-common by default
	return "lancer-common"
}
