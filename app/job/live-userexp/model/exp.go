package model

const (
	_MasterLevelMax = int32(40)
	_UserLevelMax   = int32(60)

	_ColorLevel1 = int32(9868950)
	_ColorLevel2 = int32(6406234)
	_ColorLevel3 = int32(5805790)
	_ColorLevel4 = int32(9868950)
)

var (
	_masterLevelMap = []int64{
		0,           // 0
		5000,        // 1
		20000,       // 2
		47000,       // 3
		92000,       // 4
		210000,      // 5
		406000,      // 6
		716000,      // 7
		1176000,     // 8
		1806000,     // 9
		2716000,     // 10
		3961000,     // 11
		5641000,     // 12
		7881000,     // 13
		10981000,    // 14
		15481000,    // 15
		22681000,    // 16
		31981000,    // 17
		44281000,    // 18
		60281000,    // 19
		81681000,    // 20
		113881000,   // 21
		159481000,   // 22
		221481000,   // 23
		300481000,   // 24
		398481000,   // 25
		522981000,   // 26
		690981000,   // 27
		901381000,   // 28
		1188381000,  // 29
		1561381000,  // 30
		2061381000,  // 31
		2731381000,  // 32
		3641381000,  // 33
		4781381000,  // 34
		6201381000,  // 35
		7951381000,  // 36
		9951381000,  // 37
		12201381000, // 38
		14701381000, // 39
	}

	_userLevelMap = []int64{
		100000,     // 0
		200000,     // 1
		300000,     // 2
		400000,     // 3
		500000,     // 4
		600000,     // 5
		700000,     // 6
		800000,     // 7
		900000,     // 8
		1000000,    // 9
		1800000,    // 10
		2600000,    // 11
		3400000,    // 12
		4200000,    // 13
		5000000,    // 14
		6000000,    // 15
		7000000,    // 16
		8000000,    // 17
		9000000,    // 18
		10000000,   // 19
		18000000,   // 20
		26000000,   // 21
		34000000,   // 22
		42000000,   // 23
		50000000,   // 24
		60000000,   // 25
		70000000,   // 26
		80000000,   // 27
		90000000,   // 28
		100000000,  // 29
		110000000,  // 30
		120000000,  // 31
		130000000,  // 32
		140000000,  // 33
		150000000,  // 34
		180000000,  // 35
		210000000,  // 36
		240000000,  // 37
		270000000,  // 38
		300000000,  // 39
		340000000,  // 40
		380000000,  // 41
		420000000,  // 42
		460000000,  // 43
		500000000,  // 44
		550000000,  // 45
		600000000,  // 46
		700000000,  // 47
		800000000,  // 48
		1000000000, // 49
		1200000000, // 50
		1400000000, // 51
		1600000000, // 52
		1800000000, // 53
		2000000000, // 54
		2200000000, // 55
		2400000000, // 56
		2600000000, // 57
		2800000000, // 58
		3000000000, // 59
		4000000000, // 60
	}
)

// FormatLevel 计算等级
func FormatLevel(exp *Exp) (level *Level) {
	level = &Level{Uid: exp.Uid, Uexp: exp.Uexp, Rexp: exp.Rexp, Ulevel: -1, Rlevel: -1, Color: 0}
	level.Uid = exp.Uid
	level.Uexp = exp.Uexp
	level.Rexp = exp.Rexp

	// 主播等级
	for rlevel, v := range _masterLevelMap {
		if exp.Rexp < v {
			level.Rlevel = int32(rlevel)
			level.Rnext = v - exp.Rexp
			break
		}
	}
	if level.Rlevel == -1 {
		level.Rlevel = _MasterLevelMax
	}

	// 用户等级
	for ulevel, v := range _userLevelMap {
		if exp.Uexp < v {
			level.Ulevel = int32(ulevel)
			level.Unext = v - exp.Uexp
			break
		}
	}
	if level.Ulevel == -1 {
		level.Ulevel = _UserLevelMax
	}

	// 等级颜色
	switch {
	case level.Ulevel <= 10:
		level.Color = _ColorLevel1
	case level.Ulevel <= 20:
		level.Color = _ColorLevel2
	case level.Ulevel <= 40:
		level.Color = _ColorLevel3
	case level.Ulevel < 50:
		level.Color = _ColorLevel4
	default:
		level.Color = _ColorLevel4
	}
	return
}
