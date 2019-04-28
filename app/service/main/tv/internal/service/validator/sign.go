package validator

import (
	"go-common/app/service/main/tv/internal/pkg"
	"go-common/library/ecode"
)

type SignerValidator struct {
	Signer *pkg.Signer
	Sign   string
	Val    interface{}
}

func (sv *SignerValidator) Validate() error {
	sign, err := sv.Signer.Sign(sv.Val)
	if err != nil {
		return err
	}
	if sign != sv.Sign {
		return ecode.TVIPSignErr
	}
	return nil
}
