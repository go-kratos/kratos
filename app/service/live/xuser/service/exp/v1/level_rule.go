package v1

import (
	expm "go-common/app/service/live/xuser/model/exp"
)

const (
	_AnchorLevelMax = int64(40)
	_UserLevelMax   = int64(60)

	_ColorLevel0_10  = int64(9868950)
	_ColorLevel10_20 = int64(6406234)
	_ColorLevel20_30 = int64(5805790)
	_ColorLevel30_40 = int64(10512625)
	_ColorLevel40_50 = int64(16746162)
	_ColorLevel50_60 = int64(16752445)

	_AnchorColorLevelDefault = int64(6406234)
	_AnchorColorLevel10_20   = int64(5805790)
	_AnchorColorLevel20_30   = int64(10512625)
	_AnchorColorLevel30_40   = int64(16746162)
	_AnchorColorLevel40_50   = int64(16752445)
)

type userLevelInfo struct {
	// 当前等级
	level int64
	// 下一等级
	nextLevel int64
	// 当前等级对应的经验
	userExpLeft int64
	// 下一等级对应的经验
	userExpRight int64
	// 升级到下一等级对应的经验
	userExpNextLevel int64
}

var (
	_anchorLevelMap = []userLevelInfo{
		1:  {level: 1, nextLevel: 2, userExpLeft: 0, userExpRight: 5000, userExpNextLevel: 5000},
		2:  {level: 2, nextLevel: 3, userExpLeft: 5000, userExpRight: 20000, userExpNextLevel: 15000},
		3:  {level: 3, nextLevel: 4, userExpLeft: 20000, userExpRight: 47000, userExpNextLevel: 27000},
		4:  {level: 4, nextLevel: 5, userExpLeft: 47000, userExpRight: 92000, userExpNextLevel: 45000},
		5:  {level: 5, nextLevel: 6, userExpLeft: 92000, userExpRight: 210000, userExpNextLevel: 118000},
		6:  {level: 6, nextLevel: 7, userExpLeft: 210000, userExpRight: 406000, userExpNextLevel: 196000},
		7:  {level: 7, nextLevel: 8, userExpLeft: 406000, userExpRight: 716000, userExpNextLevel: 310000},
		8:  {level: 8, nextLevel: 9, userExpLeft: 716000, userExpRight: 1176000, userExpNextLevel: 460000},
		9:  {level: 9, nextLevel: 10, userExpLeft: 1176000, userExpRight: 1806000, userExpNextLevel: 630000},
		10: {level: 10, nextLevel: 11, userExpLeft: 1806000, userExpRight: 2716000, userExpNextLevel: 910000},
		11: {level: 11, nextLevel: 12, userExpLeft: 2716000, userExpRight: 3961000, userExpNextLevel: 1245000},
		12: {level: 12, nextLevel: 13, userExpLeft: 3961000, userExpRight: 5641000, userExpNextLevel: 1680000},
		13: {level: 13, nextLevel: 14, userExpLeft: 5641000, userExpRight: 7881000, userExpNextLevel: 2240000},
		14: {level: 14, nextLevel: 15, userExpLeft: 7881000, userExpRight: 10981000, userExpNextLevel: 3100000},
		15: {level: 15, nextLevel: 16, userExpLeft: 10981000, userExpRight: 15481000, userExpNextLevel: 4500000},
		16: {level: 16, nextLevel: 17, userExpLeft: 15481000, userExpRight: 22681000, userExpNextLevel: 7200000},
		17: {level: 17, nextLevel: 18, userExpLeft: 22681000, userExpRight: 31981000, userExpNextLevel: 9300000},
		18: {level: 18, nextLevel: 19, userExpLeft: 31981000, userExpRight: 44281000, userExpNextLevel: 12300000},
		19: {level: 19, nextLevel: 20, userExpLeft: 44281000, userExpRight: 60281000, userExpNextLevel: 16000000},
		20: {level: 20, nextLevel: 21, userExpLeft: 60281000, userExpRight: 81681000, userExpNextLevel: 21400000},
		21: {level: 21, nextLevel: 22, userExpLeft: 81681000, userExpRight: 113881000, userExpNextLevel: 32200000},
		22: {level: 22, nextLevel: 23, userExpLeft: 113881000, userExpRight: 159481000, userExpNextLevel: 45600000},
		23: {level: 23, nextLevel: 24, userExpLeft: 159481000, userExpRight: 221481000, userExpNextLevel: 62000000},
		24: {level: 24, nextLevel: 25, userExpLeft: 221481000, userExpRight: 300481000, userExpNextLevel: 79000000},
		25: {level: 25, nextLevel: 26, userExpLeft: 300481000, userExpRight: 398481000, userExpNextLevel: 98000000},
		26: {level: 26, nextLevel: 27, userExpLeft: 398481000, userExpRight: 522981000, userExpNextLevel: 124500000},
		27: {level: 27, nextLevel: 28, userExpLeft: 522981000, userExpRight: 690981000, userExpNextLevel: 168000000},
		28: {level: 28, nextLevel: 29, userExpLeft: 690981000, userExpRight: 901381000, userExpNextLevel: 210400000},
		29: {level: 29, nextLevel: 30, userExpLeft: 901381000, userExpRight: 1188381000, userExpNextLevel: 287000000},
		30: {level: 30, nextLevel: 31, userExpLeft: 1188381000, userExpRight: 1561381000, userExpNextLevel: 373000000},
		31: {level: 31, nextLevel: 32, userExpLeft: 1561381000, userExpRight: 2061381000, userExpNextLevel: 500000000},
		32: {level: 32, nextLevel: 33, userExpLeft: 2061381000, userExpRight: 2731381000, userExpNextLevel: 670000000},
		33: {level: 33, nextLevel: 34, userExpLeft: 2731381000, userExpRight: 3641381000, userExpNextLevel: 910000000},
		34: {level: 34, nextLevel: 35, userExpLeft: 3641381000, userExpRight: 4781381000, userExpNextLevel: 1140000000},
		35: {level: 35, nextLevel: 36, userExpLeft: 4781381000, userExpRight: 6201381000, userExpNextLevel: 1420000000},
		36: {level: 36, nextLevel: 37, userExpLeft: 6201381000, userExpRight: 7951381000, userExpNextLevel: 1750000000},
		37: {level: 37, nextLevel: 38, userExpLeft: 7951381000, userExpRight: 9951381000, userExpNextLevel: 2000000000},
		38: {level: 38, nextLevel: 39, userExpLeft: 9951381000, userExpRight: 12201381000, userExpNextLevel: 2250000000},
		39: {level: 39, nextLevel: 40, userExpLeft: 12201381000, userExpRight: 14701381000, userExpNextLevel: 2500000000},
		40: {level: 40, nextLevel: 40, userExpLeft: 14701381000, userExpRight: 14701381000, userExpNextLevel: 0},
	}

	_userLevelMap = []userLevelInfo{
		0:  {level: 0, nextLevel: 1, userExpLeft: 0, userExpRight: 100000, userExpNextLevel: 100000},
		1:  {level: 1, nextLevel: 2, userExpLeft: 100000, userExpRight: 200000, userExpNextLevel: 100000},
		2:  {level: 2, nextLevel: 3, userExpLeft: 200000, userExpRight: 300000, userExpNextLevel: 100000},
		3:  {level: 3, nextLevel: 4, userExpLeft: 300000, userExpRight: 400000, userExpNextLevel: 100000},
		4:  {level: 4, nextLevel: 5, userExpLeft: 400000, userExpRight: 500000, userExpNextLevel: 100000},
		5:  {level: 5, nextLevel: 6, userExpLeft: 500000, userExpRight: 600000, userExpNextLevel: 100000},
		6:  {level: 6, nextLevel: 7, userExpLeft: 600000, userExpRight: 700000, userExpNextLevel: 100000},
		7:  {level: 7, nextLevel: 8, userExpLeft: 700000, userExpRight: 800000, userExpNextLevel: 100000},
		8:  {level: 8, nextLevel: 9, userExpLeft: 800000, userExpRight: 900000, userExpNextLevel: 100000},
		9:  {level: 9, nextLevel: 10, userExpLeft: 900000, userExpRight: 1000000, userExpNextLevel: 100000},
		10: {level: 10, nextLevel: 11, userExpLeft: 1000000, userExpRight: 1800000, userExpNextLevel: 800000},
		11: {level: 11, nextLevel: 12, userExpLeft: 1800000, userExpRight: 2600000, userExpNextLevel: 800000},
		12: {level: 12, nextLevel: 13, userExpLeft: 2600000, userExpRight: 3400000, userExpNextLevel: 800000},
		13: {level: 13, nextLevel: 14, userExpLeft: 3400000, userExpRight: 4200000, userExpNextLevel: 800000},
		14: {level: 14, nextLevel: 15, userExpLeft: 4200000, userExpRight: 5000000, userExpNextLevel: 800000},
		15: {level: 15, nextLevel: 16, userExpLeft: 5000000, userExpRight: 6000000, userExpNextLevel: 1000000},
		16: {level: 16, nextLevel: 17, userExpLeft: 6000000, userExpRight: 7000000, userExpNextLevel: 1000000},
		17: {level: 17, nextLevel: 18, userExpLeft: 7000000, userExpRight: 8000000, userExpNextLevel: 1000000},
		18: {level: 18, nextLevel: 19, userExpLeft: 8000000, userExpRight: 9000000, userExpNextLevel: 1000000},
		19: {level: 19, nextLevel: 20, userExpLeft: 9000000, userExpRight: 10000000, userExpNextLevel: 1000000},
		20: {level: 20, nextLevel: 21, userExpLeft: 10000000, userExpRight: 18000000, userExpNextLevel: 8000000},
		21: {level: 21, nextLevel: 22, userExpLeft: 18000000, userExpRight: 26000000, userExpNextLevel: 8000000},
		22: {level: 22, nextLevel: 23, userExpLeft: 26000000, userExpRight: 34000000, userExpNextLevel: 8000000},
		23: {level: 23, nextLevel: 24, userExpLeft: 34000000, userExpRight: 42000000, userExpNextLevel: 8000000},
		24: {level: 24, nextLevel: 25, userExpLeft: 42000000, userExpRight: 50000000, userExpNextLevel: 8000000},
		25: {level: 25, nextLevel: 26, userExpLeft: 50000000, userExpRight: 60000000, userExpNextLevel: 10000000},
		26: {level: 26, nextLevel: 27, userExpLeft: 60000000, userExpRight: 70000000, userExpNextLevel: 10000000},
		27: {level: 27, nextLevel: 28, userExpLeft: 70000000, userExpRight: 80000000, userExpNextLevel: 10000000},
		28: {level: 28, nextLevel: 29, userExpLeft: 80000000, userExpRight: 90000000, userExpNextLevel: 10000000},
		29: {level: 29, nextLevel: 30, userExpLeft: 90000000, userExpRight: 100000000, userExpNextLevel: 10000000},
		30: {level: 30, nextLevel: 31, userExpLeft: 100000000, userExpRight: 110000000, userExpNextLevel: 10000000},
		31: {level: 31, nextLevel: 32, userExpLeft: 110000000, userExpRight: 120000000, userExpNextLevel: 10000000},
		32: {level: 32, nextLevel: 33, userExpLeft: 120000000, userExpRight: 130000000, userExpNextLevel: 10000000},
		33: {level: 33, nextLevel: 34, userExpLeft: 130000000, userExpRight: 140000000, userExpNextLevel: 10000000},
		34: {level: 34, nextLevel: 35, userExpLeft: 140000000, userExpRight: 150000000, userExpNextLevel: 10000000},
		35: {level: 35, nextLevel: 36, userExpLeft: 150000000, userExpRight: 180000000, userExpNextLevel: 30000000},
		36: {level: 36, nextLevel: 37, userExpLeft: 180000000, userExpRight: 210000000, userExpNextLevel: 30000000},
		37: {level: 37, nextLevel: 38, userExpLeft: 210000000, userExpRight: 240000000, userExpNextLevel: 30000000},
		38: {level: 38, nextLevel: 39, userExpLeft: 240000000, userExpRight: 270000000, userExpNextLevel: 30000000},
		39: {level: 39, nextLevel: 40, userExpLeft: 270000000, userExpRight: 300000000, userExpNextLevel: 30000000},
		40: {level: 40, nextLevel: 41, userExpLeft: 300000000, userExpRight: 340000000, userExpNextLevel: 40000000},
		41: {level: 41, nextLevel: 42, userExpLeft: 340000000, userExpRight: 380000000, userExpNextLevel: 40000000},
		42: {level: 42, nextLevel: 43, userExpLeft: 380000000, userExpRight: 420000000, userExpNextLevel: 40000000},
		43: {level: 43, nextLevel: 44, userExpLeft: 420000000, userExpRight: 460000000, userExpNextLevel: 40000000},
		44: {level: 44, nextLevel: 45, userExpLeft: 460000000, userExpRight: 500000000, userExpNextLevel: 40000000},
		45: {level: 45, nextLevel: 46, userExpLeft: 500000000, userExpRight: 550000000, userExpNextLevel: 50000000},
		46: {level: 46, nextLevel: 47, userExpLeft: 550000000, userExpRight: 600000000, userExpNextLevel: 50000000},
		47: {level: 47, nextLevel: 48, userExpLeft: 600000000, userExpRight: 700000000, userExpNextLevel: 100000000},
		48: {level: 48, nextLevel: 49, userExpLeft: 700000000, userExpRight: 800000000, userExpNextLevel: 100000000},
		49: {level: 49, nextLevel: 50, userExpLeft: 800000000, userExpRight: 1000000000, userExpNextLevel: 200000000},
		50: {level: 50, nextLevel: 51, userExpLeft: 1000000000, userExpRight: 1200000000, userExpNextLevel: 200000000},
		51: {level: 51, nextLevel: 52, userExpLeft: 1200000000, userExpRight: 1400000000, userExpNextLevel: 200000000},
		52: {level: 52, nextLevel: 53, userExpLeft: 1400000000, userExpRight: 1600000000, userExpNextLevel: 200000000},
		53: {level: 53, nextLevel: 54, userExpLeft: 1600000000, userExpRight: 1800000000, userExpNextLevel: 200000000},
		54: {level: 54, nextLevel: 55, userExpLeft: 1800000000, userExpRight: 2000000000, userExpNextLevel: 200000000},
		55: {level: 55, nextLevel: 56, userExpLeft: 2000000000, userExpRight: 2200000000, userExpNextLevel: 200000000},
		56: {level: 56, nextLevel: 57, userExpLeft: 2200000000, userExpRight: 2400000000, userExpNextLevel: 200000000},
		57: {level: 57, nextLevel: 58, userExpLeft: 2400000000, userExpRight: 2600000000, userExpNextLevel: 200000000},
		58: {level: 58, nextLevel: 59, userExpLeft: 2600000000, userExpRight: 2800000000, userExpNextLevel: 200000000},
		59: {level: 59, nextLevel: 60, userExpLeft: 2800000000, userExpRight: 3000000000, userExpNextLevel: 200000000},
		60: {level: 60, nextLevel: 60, userExpLeft: 3000000000, userExpRight: 3000000000, userExpNextLevel: 0},
	}
)

// FormatLevel ...
// 等级转换
func (s *UserExpService) FormatLevel(exps []*expm.Exp) (level map[int64]*expm.LevelInfo) {
	level = make(map[int64]*expm.LevelInfo)
	if len(exps) <= 0 {
		return
	}
	for _, v := range exps {
		uid := v.UID
		level[uid] = &expm.LevelInfo{}
		level[uid].UID = uid
		level[uid].UserLevel = s.getUserLevel(v)
		level[uid].AnchorLevel = s.getAnchorLevel(v)
		level[uid].CTime = v.CTime
		level[uid].MTime = v.MTime
	}
	return
}

func (s *UserExpService) getUserLevel(originDBResult *expm.Exp) (userLevel expm.UserLevelInfo) {
	userLevel = expm.UserLevelInfo{Level: -1}
	for k, value := range _userLevelMap {
		if originDBResult.Uexp < value.userExpRight {
			userLevel.Level = value.level
			userLevel.NextLevel = value.nextLevel
			userLevel.UserExpLeft = value.userExpLeft
			userLevel.UserExpRight = value.userExpRight
			userLevel.UserExp = originDBResult.Uexp
			userLevel.UserExpNextLevel = value.userExpNextLevel
			nextLevel := k + 1
			if nextLevel < len(_userLevelMap) {
				userLevel.UserExpNextLeft = _userLevelMap[nextLevel].userExpLeft
				userLevel.UserExpNextRight = _userLevelMap[nextLevel].userExpRight
				userLevel.IsLevelTop = 0
			}

			break
		}
	}
	if userLevel.Level == -1 {
		userLevel.Level = _userLevelMap[_UserLevelMax].level
		userLevel.NextLevel = _userLevelMap[_UserLevelMax].nextLevel
		userLevel.UserExpLeft = _userLevelMap[_UserLevelMax].userExpLeft
		userLevel.UserExpRight = _userLevelMap[_UserLevelMax].userExpRight
		userLevel.UserExp = originDBResult.Uexp
		userLevel.UserExpNextLevel = _userLevelMap[_UserLevelMax].userExpNextLevel
		userLevel.UserExpNextLeft = _userLevelMap[_UserLevelMax].userExpLeft
		userLevel.UserExpNextRight = _userLevelMap[_UserLevelMax].userExpRight
		userLevel.IsLevelTop = 1
	}
	userLevel.Color = s.getUserLevelColor(userLevel.Level)
	return
}

func (s *UserExpService) getAnchorLevel(originDBResult *expm.Exp) (anchorLevel expm.AnchorLevelInfo) {
	anchorLevel = expm.AnchorLevelInfo{Level: -1}
	for k, value := range _anchorLevelMap {
		if originDBResult.Rexp < value.userExpRight {
			anchorLevel.Level = value.level
			anchorLevel.NextLevel = value.nextLevel
			anchorLevel.UserExpLeft = value.userExpLeft
			anchorLevel.UserExpRight = value.userExpRight
			anchorLevel.UserExp = originDBResult.Rexp
			anchorLevel.UserExpNextLevel = value.userExpNextLevel
			if anchorLevel.UserExp == 0 {
				anchorLevel.AnchorScore = 0
			} else {
				anchorLevel.AnchorScore = anchorLevel.UserExp / 100
			}

			nextLevel := k + 1
			if nextLevel < len(_userLevelMap) {
				anchorLevel.UserExpNextLeft = _userLevelMap[nextLevel].userExpLeft
				anchorLevel.UserExpNextRight = _userLevelMap[nextLevel].userExpRight
				anchorLevel.IsLevelTop = 0
			}
			break
		}
	}
	if anchorLevel.Level == -1 {
		anchorLevel.Level = _anchorLevelMap[_AnchorLevelMax].level
		anchorLevel.NextLevel = _anchorLevelMap[_AnchorLevelMax].nextLevel
		anchorLevel.UserExpLeft = _anchorLevelMap[_AnchorLevelMax].userExpLeft
		anchorLevel.UserExpRight = _anchorLevelMap[_AnchorLevelMax].userExpRight
		anchorLevel.UserExp = originDBResult.Rexp
		anchorLevel.UserExpNextLevel = _anchorLevelMap[_AnchorLevelMax].userExpNextLevel
		anchorLevel.UserExpNextLeft = _userLevelMap[_UserLevelMax].userExpLeft
		anchorLevel.UserExpNextRight = _userLevelMap[_UserLevelMax].userExpRight
		anchorLevel.IsLevelTop = 1
		if anchorLevel.UserExp == 0 {
			anchorLevel.AnchorScore = 0
		} else {
			anchorLevel.AnchorScore = anchorLevel.UserExp / 100
		}
	}
	anchorLevel.Color = s.getAnchorLevelColor(anchorLevel.Level)
	return
}

func (s *UserExpService) getUserLevelColor(level int64) (color int64) {
	switch {
	case level <= 10:
		color = _ColorLevel0_10
	case level <= 20:
		color = _ColorLevel10_20
	case level <= 30:
		color = _ColorLevel20_30
	case level <= 40:
		color = _ColorLevel30_40
	case level <= 50:
		color = _ColorLevel40_50
	default:
		color = _ColorLevel50_60
	}
	return
}

func (s *UserExpService) getAnchorLevelColor(level int64) (color int64) {
	switch {
	case level <= 10:
		color = _AnchorColorLevelDefault
	case level <= 20:
		color = _AnchorColorLevel10_20
	case level <= 30:
		color = _AnchorColorLevel20_30
	case level <= 40:
		color = _AnchorColorLevel30_40
	case level <= 50:
		color = _AnchorColorLevel40_50
	default:
		color = _AnchorColorLevel40_50
	}
	return
}
