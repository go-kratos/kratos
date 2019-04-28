package archive

const (
	// Donghua 动画
	Donghua = int16(1)
	// Music 音乐
	Music = int16(3)
	// Game 游戏
	Game = int16(4)
	// Ent 娱乐
	Ent = int16(5)
	// Tv 电视剧
	Tv = int16(11)
	// Bangumi 番剧
	Bangumi = int16(13)
	// Movie 电影
	Movie = int16(23)
	// Tech 科技
	Tech = int16(36)
	// Kichiku 鬼畜
	Kichiku = int16(119)
	// Dance 舞蹈
	Dance = int16(129)
	// Fashion 时尚
	Fashion = int16(155)
	// Life 生活
	Life = int16(160)
	// Ad 广告
	Ad = int16(165)
	// Guochuang 国创
	Guochuang = int16(167)
	// Filmwithtele 影视
	Filmwithtele = int16(181)
	// Documentary 纪录片
	Documentary = int16(177)
)

// Type info
type Type struct {
	ID   int16  `json:"id"`
	PID  int16  `json:"pid"`
	Name string `json:"name"`
	Desc string `json:"description"`
}
