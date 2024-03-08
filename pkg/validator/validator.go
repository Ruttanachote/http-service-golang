package validator

import (
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils"
	"github.com/go-playground/validator"
)

func Validate(data interface{}) []*utils.ErrorResponse {
	var errors []*utils.ErrorResponse
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element utils.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
