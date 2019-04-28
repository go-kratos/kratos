package dao

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_defaultExpiration = 10
	_tokenNamespace = "upload_token:"
)

// RequestUploadToken generates a token for subsequent upload.
// Token will expire in a specific duration.
func (d *Dao) RequestUploadToken(ctx context.Context, bucket, operator string, now int64) (token string, err error) {
	token = genToken(bucket, operator, now)

	conn := d.redis.Get(ctx)
	defer conn.Close()

	nsToken := namespaceToken(token)
	if _, err = conn.Do("SETEX", nsToken, _defaultExpiration, "1"); err != nil {
		log.Error("conn.Do(SETEX %s %d) failure(%v)", nsToken, _defaultExpiration, err)
	}

	return
}

// VerifyUploadToken verifies if a token is legal.
func (d *Dao) VerifyUploadToken(ctx context.Context, token string) bool {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	nsToken := namespaceToken(token)
	valid, err := redis.Bool(conn.Do("GET", nsToken))
	if err != nil && err != redis.ErrNil {
		log.Warn("conn.Do(GET %s) failure(%v)", nsToken, err)
		return false
	}

	return valid
}

func genToken(bucket, operator string, now int64) string {
	sha := sha1.New()
	sha.Write([]byte(fmt.Sprintf("i love bilibili:%s:%s:%d", bucket, operator, now)))
	return fmt.Sprintf("%s:%d", hex.EncodeToString(sha.Sum([]byte(""))), now)
}

// Avoid key collision.
func namespaceToken(token string) string {
	return _tokenNamespace + token
}
