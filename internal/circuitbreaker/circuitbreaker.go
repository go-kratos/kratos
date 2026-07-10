package circuitbreaker

import "errors"

// ErrNotAllowed is returned when the circuit breaker is open.
var ErrNotAllowed = errors.New("circuitbreaker: not allowed for circuit open")

// CircuitBreaker is a circuit breaker.
type CircuitBreaker interface {
	Allow() error
	MarkSuccess()
	MarkFailed()
}
