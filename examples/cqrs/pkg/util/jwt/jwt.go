package jwt

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

// ExtractUserInfoFromClaims 从jwt载荷中提取出用户信息
func ExtractUserInfoFromClaims(ctx context.Context) uint64 {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return 0
	}

	str := claims.(jwtv4.MapClaims)["user_id"]
	if str == nil {
		return 0
	}
	switch str.(type) {
	case string:
		break
	default:
		return 0
	}

	userId, err := strconv.ParseUint(str.(string), 10, 64)
	if err != nil {
		return 0
	}

	return userId
}

// EncryptUserInfoToJwtToken 将用户ID加密到jwt载荷中
func EncryptUserInfoToJwtToken(key []byte, userId uint64) (string, error) {
	claims := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256,
		jwtv4.MapClaims{
			"user_id":      strconv.FormatUint(userId, 10),
			"authority_id": strconv.FormatUint(userId, 10),
		})
	return claims.SignedString(key)
}
