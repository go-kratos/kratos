package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	text := "password"
	hash, _ := HashPassword(text)
	fmt.Println(hash)
}

func TestCheckPasswordHash(t *testing.T) {
	text := "123456"
	//hash1 := "$2a$10$4KoNdzgqllEiHsgTCCdtBu7RqyLrw.f7vR9cfhGJRowWiB7Q/2SjG"
	//hash2 := "$2a$10$BPQS8mjrm3DJZnLGW3hxH.bT/piDeKjxgl/etzxO92Id21BZt8eH."
	hash3 := "$2a$10$ygWrRwHCzg2GUpz0UK40kuWAGva121VkScpcdMNsDCih2U/bL2qYy"
	bMatched := CheckPasswordHash(text, hash3)
	assert.True(t, bMatched)

	bMatched = CheckPasswordHash(text, hash3)
	assert.True(t, bMatched)
}

func TestJwtToken(t *testing.T) {
	const bearerWord string = "Bearer"
	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowfQ.XgcKAAjHbA6o4sxxbEaMi05ingWvKdCNnyW9wowbJvs"
	auths := strings.SplitN(token, " ", 2)
	assert.Equal(t, len(auths), 2)
	assert.Equal(t, strings.EqualFold(auths[0], bearerWord), true, "JWT token is missing")
}
