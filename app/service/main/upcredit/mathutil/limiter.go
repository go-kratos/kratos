package mathutil

import "time"

//Limiter speed limiter
type Limiter struct {
	Rate  float64 // 每秒多少个
	token chan time.Time
	timer *time.Ticker
}

//Token get token
func (l *Limiter) Token() (c <-chan time.Time) {
	return l.token
}

func (l *Limiter) putToken() {
	for t := range l.timer.C {
		l.token <- t
	}
}

//NewLimiter create new limiter
func NewLimiter(rate float64) *Limiter {
	var l = &Limiter{
		Rate:  rate,
		token: make(chan time.Time, 1),
		timer: time.NewTicker(time.Duration(1.0 / rate * float64(time.Second))),
	}
	go l.putToken()
	return l
}
