package model

import (
	"errors"

	"go-common/library/log"
)

//Stra 实验策略
type Stra struct {
	//精度
	Precision int `json:"precision"`
	//依次比例
	Ratio []int `json:"ratio"`
}

func (s *Stra) check() (isValid bool) {
	sum := 0
	for _, r := range s.Ratio {
		sum += r
	}
	isValid = (sum == s.Precision)
	return
}

//Check ensure stra valid
func (s *Stra) Check() (isValid bool) {
	return s.check()
}

//Version calculate version by score
func (s *Stra) Version(score int) (version int, err error) {
	if !s.check() {
		err = errors.New("the sum of ratio is not equal to precision")
		log.Error("[model.stra|Version] s.check failed")
		return
	}

	if score >= s.Precision || score < 0 {
		err = errors.New("score should between 0 and s.Precision")
		return
	}

	for i, r := range s.Ratio {
		if score >= r {
			score -= r
		} else {
			version = i
			break
		}
	}
	return
}
