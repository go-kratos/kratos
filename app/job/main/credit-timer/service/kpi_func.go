package service

// 根据投准率计算投准系数
func (s *Service) voteRightRatio(vr float64) (vf float64) {
	switch {
	case vr >= 0.9:
		vf = float64(1.2)
	case vr >= 0.8 && vr < 0.9:
		vf = float64(1.1)
	case vr >= 0.7 && vr < 0.8:
		vf = float64(0.9)
	case vr >= 0.6 && vr < 0.7:
		vf = float64(0.6)
	case vr >= 0.5 && vr < 0.6:
		vf = float64(0.3)
	case vr >= 0.4 && vr < 0.5:
		vf = float64(0.1)
	default:
		vf = float64(0)
	}
	return
}

// 根据活跃天数计算活跃系数
func (s *Service) activeDaysRatio(activeDays int64) (af float64) {
	switch {
	case activeDays >= 26:
		af = float64(1.3)
	case activeDays >= 21 && activeDays <= 25:
		af = float64(1.2)
	case activeDays >= 16 && activeDays <= 20:
		af = float64(1.1)
	case activeDays >= 11 && activeDays <= 15:
		af = float64(1.0)
	case activeDays >= 6 && activeDays <= 10:
		af = float64(0.9)
	case activeDays >= 1 && activeDays <= 5:
		af = float64(0.7)
	default:
		af = float64(0)
	}
	return
}

// 根据观点数量计算观点数量系数
func (s *Service) opinionNumsRatio(opinionNums int64) (of float64) {
	switch {
	case opinionNums >= 31:
		of = 1.3
	case opinionNums >= 16 && opinionNums <= 30:
		of = 1.2
	case opinionNums >= 6 && opinionNums <= 15:
		of = 1.1
	case opinionNums >= 1 && opinionNums <= 5:
		of = 1
	default:
		of = 0.8
	}
	return
}

// 根据观点（赞-踩）数计算观点质量系数
func (s *Service) opinionQualityRatio(opinionQuality int64) (oqf float64) {
	switch {
	case opinionQuality >= 16:
		oqf = 1.3
	case opinionQuality >= 6 && opinionQuality <= 15:
		oqf = 1.2
	case opinionQuality >= 1 && opinionQuality <= 5:
		oqf = 1.1
	case opinionQuality == 0:
		oqf = 1
	case opinionQuality >= -10 && opinionQuality <= -1:
		oqf = 0.8
	case opinionQuality >= -20 && opinionQuality <= -11:
		oqf = 0.7
	case opinionQuality <= -21:
		oqf = 0.5
	}
	return
}
