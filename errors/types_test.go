package errors

import (
	"testing"
)

func TestTypes(t *testing.T) {
	var (
		input = []error{
			BadRequest("reason_400", "message_400"),
			Unauthorized("reason_401", "message_401"),
			Forbidden("reason_403", "message_403"),
			NotFound("reason_404", "message_404"),
			Conflict("reason_409", "message_409"),
			InternalServer("reason_500", "message_500"),
			ServiceUnavailable("reason_503", "message_503"),
			GatewayTimeout("reason_504", "message_504"),
			ClientClosed("reason_499", "message_499"),
		}
		output = []func(error) bool{
			IsBadRequest,
			IsUnauthorized,
			IsForbidden,
			IsNotFound,
			IsConflict,
			IsInternalServer,
			IsServiceUnavailable,
			IsGatewayTimeout,
			IsClientClosed,
		}
	)

	for i, in := range input {
		if !output[i](in) {
			t.Errorf("not expect: %v", in)
		}
	}
}
