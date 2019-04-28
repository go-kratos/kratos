package service

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"go-common/library/log"
)

var (
	checkTelReg = regexp.MustCompile("^[0-9]*$")
	base64Reg   = regexp.MustCompile(`[a-zA-Z0-9_\-+/=]*`)
)

func (s *Service) doEncrypt(param string) []byte {
	var res = make([]byte, 0)
	if param == "" || len(param) == 0 {
		return res
	}
	input := []byte(param)
	res, _ = Encrypt(input, []byte(aesKey))
	return res
}

func (s *Service) doDecrypt(param []byte) (string, error) {
	var (
		err error
		res = make([]byte, 0)
	)
	if len(param) == 0 {
		return "", nil
	}
	if res, err = Decrypt(param, []byte(aesKey)); err != nil {
		return "", err
	}
	return string(res), nil
}

func (s *Service) doHash(plaintext string) []byte {
	var res = make([]byte, 0)
	if plaintext == "" {
		return res
	}
	hash := md5.New()
	hash.Write([]byte(plaintext))
	hash.Write([]byte(securitySalt))
	res = hash.Sum(nil)
	return res
}

// IsTel isTel
func IsTel(tel string) bool {
	return checkTelReg.MatchString(tel)
}

// IsMail is mail
func IsMail(mail string) bool {
	return strings.Contains(mail, "@")
}

func addBase64Padding(value string) string {
	strArr := base64Reg.FindAllString(value, -1)
	value = strings.Replace(strings.Trim(fmt.Sprint(strArr), "[]"), " ", "", -1)
	m := len(value) % 4
	if m != 0 {
		before := value
		value += strings.Repeat("=", 4-m)
		log.Error("addBase64Padding before(len:%d)(%s) after(len:%d)(%s)", len(before), before, len(value), value)
	}
	return value
}
