package binding

import (
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// MIME
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

// Binding http binding request interface.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// StructValidator http validator interface.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error

	// RegisterValidation adds a validation Func to a Validate's map of validators denoted by the key
	// NOTE: if the key already exists, the previous validation function will be replaced.
	// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
	RegisterValidation(string, validator.Func) error
}

// Validator default validator.
var Validator StructValidator = &defaultValidator{}

// Binding
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
)

// Default get by binding type by method and contexttype.
func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	}

	contentType = stripContentTypeParam(contentType)
	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
		return Form
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}

func stripContentTypeParam(contentType string) string {
	i := strings.Index(contentType, ";")
	if i != -1 {
		contentType = contentType[:i]
	}
	return contentType
}
