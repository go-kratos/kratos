package jwt

import "time"

// TokenManager manager the jwt token.
type TokenManager interface {
	Token() string
}

//token store jwt token and manager that manage the token
type token struct {
	jwtToken string
	manager  TokenManager
}

// runRefreshTokenJob run a goroutine in order to refresh the jwt token periodically.
func (receiver *token) runRefreshTokenJob(refreshDuration time.Duration) {
	go func() {
		wait := refreshDuration
		timer := time.NewTimer(wait)
		defer timer.Stop()
		for {
			<-timer.C
			receiver.refreshToken()
			timer.Reset(wait)
		}
	}()
}

// refreshToken refresh the jwt token
func (receiver *token) refreshToken() {
	receiver.jwtToken = receiver.manager.Token()
}

// NewToken create token object,tokenManager must not be nil
func NewToken(tokenManager TokenManager, refreshDuration time.Duration) *token {
	target := &token{
		manager: tokenManager,
	}
	// init token
	target.refreshToken()
	// run a job to refresh jwt token periodically
	target.runRefreshTokenJob(refreshDuration)
	return target
}
