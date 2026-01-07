package invitation

import (
	"github.com/dyxj/bigbackend/pkg/errorx"
	"github.com/dyxj/bigbackend/pkg/validx"
)

type CreateRequest struct {
	Email string `json:"email"`
}

func (r CreateRequest) Validate() *errorx.ValidationError {
	errors := make(map[string]string)

	if !validx.IsEmail(r.Email) {
		errors["email"] = "is not a valid email"
	}

	if len(errors) > 0 {
		return &errorx.ValidationError{Properties: errors}
	}

	return nil
}

type CreateResponse struct {
	Email string `json:"email"`
}
