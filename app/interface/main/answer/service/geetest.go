package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"go-common/app/interface/main/answer/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// preProcess getGeetestChal
func (s *Service) preProcess(c context.Context, mid int64, ip, clientType string, newCaptcha int) (res *model.ProcessRes, err error) {
	var pre string
	res = &model.ProcessRes{}
	gc, geeType := s.geetestDao.GeeConfig(clientType, s.c.Geetest)
	res.CaptchaID = gc.CaptchaID
	res.NewCaptcha = newCaptcha
	if pre, err = s.geetestDao.PreProcess(c, mid, ip, geeType, gc, newCaptcha); err != nil || pre == "" {
		log.Error("s.geetestDao.PreProcess(%d) err(%v)", mid, err)
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		challenge := hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
		res.Challenge = challenge
		err = nil
		return
	}
	res.Success = 1
	slice := md5.Sum([]byte(pre + gc.PrivateKEY))
	res.Challenge = hex.EncodeToString(slice[:])
	return
}

// validate recheck the seccode
func (s *Service) validate(c context.Context, challenge, validate, seccode, clientType, ip string, success int, mid int64) (stat bool) {
	if len(validate) != 32 {
		log.Error("s.Validate(%s,%s,%s,%d) err(validate not eq 32byte)", challenge, validate, seccode, mid)
		stat = s.failbackValidate(c, challenge, validate, seccode)
		log.Info("s.failbackValidate(%s,%s,%s,%d), stat(%t)", challenge, validate, seccode, mid, stat)
		return
	}
	if success != 1 {
		slice := md5.Sum([]byte(challenge))
		stat = hex.EncodeToString(slice[:]) == validate
		return
	}
	gc, geeType := s.geetestDao.GeeConfig(clientType, s.c.Geetest)
	slice := md5.Sum([]byte(gc.PrivateKEY + "geetest" + challenge))
	if hex.EncodeToString(slice[:]) != validate {
		log.Error("s.Validate(%s,%s,%s,%d) err(challenge not found)", challenge, validate, seccode, mid)
		return
	}
	res, err := s.geetestDao.Validate(c, challenge, seccode, geeType, ip, gc.CaptchaID, mid)
	if err != nil {
		if errors.Cause(err) == context.DeadlineExceeded { // for geetest timeout
			stat = true
			return
		}
		log.Error("s.geetestDao.Validate(%d) err(%v)", mid, err)
		return
	}
	slice = md5.Sum([]byte(seccode))
	stat = hex.EncodeToString(slice[:]) == res.Seccode
	return
}

//failbackValidate geetest failback model.
func (s *Service) failbackValidate(c context.Context, challenge, validate, seccode string) bool {
	varr := strings.Split(validate, "_")
	if len(varr) < 3 {
		return false
	}
	encodeAns := varr[0]
	encodeFbii := varr[1]
	encodeIgi := varr[2]
	decodeAns := s.decodeResponse(challenge, encodeAns)
	decodeFbii := s.decodeResponse(challenge, encodeFbii)
	decodeIgi := s.decodeResponse(challenge, encodeIgi)
	return s.validateFailImage(decodeAns, decodeFbii, decodeIgi)
}

func (s *Service) decodeResponse(challenge, userresponse string) (res int) {
	if len(userresponse) > 100 {
		return
	}
	digits := []int{1, 2, 5, 10, 50}
	key := make(map[rune]int)
	for _, i := range challenge {
		if _, exist := key[i]; exist {
			continue
		}
		value := digits[len(key)%5]
		key[i] = value
	}
	for _, i := range userresponse {
		res += key[i]
	}
	res -= s.decodeRandBase(challenge)
	return
}

func (s *Service) decodeRandBase(challenge string) int {
	baseStr := challenge[32:]
	var tempList []int
	for _, char := range baseStr {
		tempChar := int(char)
		result := tempChar - 48
		if tempChar > 57 {
			result = tempChar - 87
		}
		tempList = append(tempList, result)
	}
	return tempList[0]*36 + tempList[1]
}

func (s *Service) md5Encode(values []byte) string {
	return fmt.Sprintf("%x", md5.Sum(values))
}

func (s *Service) validateFailImage(ans, fullBgIndex, imgGrpIndex int) bool {
	var thread float64 = 3
	fullBg := s.md5Encode([]byte(strconv.Itoa(fullBgIndex)))[0:10]
	imgGrp := s.md5Encode([]byte(strconv.Itoa(imgGrpIndex)))[10:20]
	var answerDecode []byte
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			answerDecode = append(answerDecode, fullBg[i])
		} else if i%2 == 1 {
			answerDecode = append(answerDecode, imgGrp[i])
		}
	}
	xDecode := answerDecode[4:]
	xInt64, err := strconv.ParseInt(string(xDecode), 16, 32)
	if err != nil {
		log.Error("%+v", err)
	}
	xInt := int(xInt64)
	result := xInt % 200
	if result < 40 {
		result = 40
	}
	return math.Abs(float64(ans-result)) < thread
}
