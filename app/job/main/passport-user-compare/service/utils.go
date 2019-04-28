package service

import (
	"crypto/md5"
	"math/big"
	"net"
)

func (s *Service) doEncrypt(param string) ([]byte, error) {
	var (
		err error
		res = make([]byte, 0)
	)
	if len(param) == 0 {
		return res, nil
	}
	input := []byte(param)
	if res, err = Encrypt(input, []byte(aesKey)); err != nil {
		return nil, nil
	}
	return res, nil
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
	hash.Write([]byte(md5slat))
	res = hash.Sum(nil)
	return res
}

// InetAtoN .
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}
