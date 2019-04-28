package geetest

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/geetest"
	"go-common/app/interface/main/account/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Service is
type Service struct {
	c          *conf.Config
	geetestDao *geetest.Dao
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		geetestDao: geetest.New(c),
	}
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close dao.
func (s *Service) Close() {}

// PreProcess preprocessing the geetest and get to challenge
func (s *Service) PreProcess(c context.Context, req *model.GeeCaptchaRequest) (res *model.ProcessRes, err error) {
	var pre string
	res = &model.ProcessRes{}
	gc, clientType := s.geetestDao.GeeConfig(req.ClientType, s.c.Geetest)
	res.CaptchaID = gc.CaptchaID
	res.NewCaptcha = 1
	if pre, err = s.geetestDao.PreProcess(c, req.MID, clientType, gc, 1); err != nil || pre == "" {
		log.Error("s.geetestDao.PreProcess(%+v) err(%v)", req, err)
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		challenge := hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
		res.Challenge = challenge
		return
	}
	res.Success = 1
	slice := md5.Sum([]byte(pre + gc.PrivateKEY))
	res.Challenge = hex.EncodeToString(slice[:])
	log.Info("PreProcess success(%+v) ", res)
	return
}

// Validate recheck the challenge validate seccode
func (s *Service) Validate(c context.Context, req *model.GeeCheckRequest) (stat bool) {
	if len(req.Validate) != 32 {
		log.Error("s.Validate(%+v) err(validate not eq 32byte)", req)
		stat = s.failbackValidate(c, req.Challenge, req.Validate, req.Seccode)
		return
	}
	gc, clientType := s.geetestDao.GeeConfig(req.ClientType, s.c.Geetest)
	slice := md5.Sum([]byte(gc.PrivateKEY + "geetest" + req.Challenge))
	if hex.EncodeToString(slice[:]) != req.Validate {
		log.Error("s.Validate(%+v) err(challenge not found)", req)
		return
	}
	res, err := s.geetestDao.Validate(c, req.Challenge, req.Seccode, clientType, gc.CaptchaID, req.MID)
	if err != nil {
		if errors.Cause(err) == context.DeadlineExceeded { // for geetest timeout
			stat = true
			return
		}
		log.Error("s.geetestDao.Validate(%+v) err(%v)", req, err)
		return
	}
	slice = md5.Sum([]byte(req.Seccode))
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
