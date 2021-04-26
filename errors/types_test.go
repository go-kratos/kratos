package errors

import "testing"

func TestTypes(t *testing.T) {
	var (
		input = []error{
			BadRequest("domain_400", "reason_400", "message_400"),
			Unauthorized("domain_401", "reason_401", "message_401"),
			Forbidden("domain_403", "reason_403", "message_403"),
			NotFound("domain_404", "reason_404", "message_404"),
			Conflict("domain_409", "reason_409", "message_409"),
			InternalServer("domain_500", "reason_500", "message_500"),
			ServiceUnavailable("domain_503", "reason_503", "message_503"),
		}
		output = []func(error) bool{
			IsBadRequest,
			IsUnauthorized,
			IsForbidden,
			IsNotFound,
			IsConflict,
			IsInternalServer,
			IsServiceUnavailable,
		}
	)

	for i, in := range input {
		if !output[i](in) {
			t.Errorf("not expect: %v", in)
		}
	}
}
