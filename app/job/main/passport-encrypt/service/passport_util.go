package service

import "go-common/app/job/main/passport-encrypt/model"

// EncryptAccount Encrypt Account
func EncryptAccount(a *model.OriginAccount) (res *model.EncryptAccount) {
	return &model.EncryptAccount{
		Mid:            a.Mid,
		UserID:         a.UserID,
		Uname:          a.Uname,
		Pwd:            a.Pwd,
		Salt:           a.Salt,
		Email:          a.Email,
		Tel:            doEncrypt(a.Tel),
		CountryID:      a.CountryID,
		MobileVerified: a.MobileVerified,
		Isleak:         a.Isleak,
		Mtime:          a.Mtime,
	}
}

// DecryptAccount Decrypt Account
func DecryptAccount(a *model.EncryptAccount) (res *model.OriginAccount) {
	return &model.OriginAccount{
		Mid:            a.Mid,
		UserID:         a.UserID,
		Uname:          a.Uname,
		Pwd:            a.Pwd,
		Salt:           a.Salt,
		Email:          a.Email,
		Tel:            doDecrypt(a.Tel),
		CountryID:      a.CountryID,
		MobileVerified: a.MobileVerified,
		Isleak:         a.Isleak,
		Mtime:          a.Mtime,
	}
}

func doEncrypt(param string) []byte {
	var (
		err error
		res = make([]byte, 0)
	)
	if param == "" || len(param) == 0 {
		return nil
	}
	input := []byte(param)
	if res, err = Encrypt(input); err != nil {
		return input
	}
	return res
}

func doDecrypt(param []byte) string {
	var (
		err error
		res = make([]byte, 0)
	)
	if param == nil {
		return ""
	}
	input := []byte(param)
	if res, err = Decrypt(input); err != nil {
		return string(param)
	}
	return string(res)
}
