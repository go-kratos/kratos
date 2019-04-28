package service

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"net"
)

func (s *Service) doEncrypt(param string) []byte {
	var (
		err error
		res = make([]byte, 0)
	)
	if param == "" || len(param) == 0 {
		return nil
	}
	input := []byte(param)
	if res, err = s.CBCEncrypt(input); err != nil {
		return input
	}
	return res
}

//func (s *Service) doDecrypt(param []byte) string {
//	var (
//		err error
//		res = make([]byte, 0)
//	)
//	if param == nil {
//		return ""
//	}
//	input := []byte(param)
//	if res, err = s.CBCDecrypt(input); err != nil {
//		return string(param)
//	}
//	return string(res)
//}

// InetNtoA convert int64 to ip addr .
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// InetAtoN convert ip addr to int64.
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func (s *Service) doHash(plaintext string) []byte {
	var res = make([]byte, 0)
	if plaintext == "" {
		return res
	}
	hash := md5.New()
	hash.Write([]byte(plaintext))
	hash.Write([]byte(s.salt))
	res = hash.Sum(nil)
	return res
}
