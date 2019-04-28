package service

import (
	"context"
	"time"

	"go-common/app/service/main/passport-game/model"
)

const (
	_originPubKey = "-----BEGIN PUBLIC KEY-----\n" +
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdScM09sZJqFPX7bvmB2y6i08J\n" +
		"bHsa0v4THafPbJN9NoaZ9Djz1LmeLkVlmWx1DwgHVW+K7LVWT5FV3johacVRuV98\n" +
		"37+RNntEK6SE82MPcl7fA++dmW2cLlAjsIIkrX+aIvvSGCuUfcWpWFy3YVDqhuHr\n" +
		"NDjdNcaefJIQHMW+sQIDAQAB\n" +
		"-----END PUBLIC KEY-----\n"

	_originPrivKey = "-----BEGIN PRIVATE KEY-----\n" +
		"MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJ1JwzT2xkmoU9ft\n" +
		"u+YHbLqLTwlsexrS/hMdp89sk302hpn0OPPUuZ4uRWWZbHUPCAdVb4rstVZPkVXe\n" +
		"OiFpxVG5X3zfv5E2e0QrpITzYw9yXt8D752ZbZwuUCOwgiStf5oi+9IYK5R9xalY\n" +
		"XLdhUOqG4es0ON01xp58khAcxb6xAgMBAAECgYAsBeQ8I8HWBeYJrsGDnZpiD/G8\n" +
		"On+uP1XbtdYtKT+SsTs1RfTW0jhtvJex2yJPFTjzDIeew6fxk22jMgLlLTyLtc+H\n" +
		"LFmNt6/5DIxeSslx0cZOJW2w1KLyzPBIVQ+q1rAu4wZx4LSTFzAHQGWkGH3352Nj\n" +
		"40GIMOMOLrOReb8/rQJBAM4S8mChP9NAWvjSyQY4jqHQsCmcDeV1Tuwo5jxeMMZ5\n" +
		"YSrtwEEaLN4mkvpMeQR60RtjdJXxfySJY0lF/e5RZ7sCQQDDZQjUmry41Ar6uSGj\n" +
		"OrZADOSmqZy0qCijgthkN2Gblx9eyFFMFcakJECF/m2zUo5MeqKu5RTDXBZt0b3l\n" +
		"JV6DAkEAnIv0KMgWbmsDMOcf42PvpqmcSd/NBrU5AVqInO+I6h2nXS9Dz7EMyK5R\n" +
		"FWgmvup2E/JXzNiql5zvGejb4MFipQJAVScl5wmcb2wxcLzXtPw0SsuTpjJK0cxr\n" +
		"EX9HcL1V82mzySnBjEf9LrGB0SNliX3T9+6GEXRSTSVHvQpoGIHlowJBAKJBTxKi\n" +
		"9ypIdK1mA7kw8+g+YdXliov5B2MMe1yaGyGjz3YGQ8N1YLBK6Yp4KOwodFurmgcX\n" +
		"ozLanqOwPdNW1Nk=\n" +
		"-----END PRIVATE KEY-----\n"
)

// RSAKeyOrigin get rsa key via model api.
func (s *Service) RSAKeyOrigin(c context.Context) (res *model.RSAKey, err error) {
	return s.d.RSAKeyOrigin(c)
}

// RSAKey get cloud rsa pub key and seconds ts hash.
func (s *Service) RSAKey(c context.Context) *model.RSAKey {
	return &model.RSAKey{
		Hash: TsSeconds2Hash(time.Now().Unix()),
		Key:  _originPubKey,
	}
}
