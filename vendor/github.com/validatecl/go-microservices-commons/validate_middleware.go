package commons

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-kit/kit/endpoint"
	validator "github.com/go-playground/validator/v10"
)

var (
	v9 *validator.Validate
)

//ValidateMiddleware Validate middleware endpoint
func ValidateMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, reqIn interface{}) (reqVal interface{}, err error) {
			if err = validate(reqIn); err != nil {
				return nil, err
			}
			return next(ctx, reqIn)
		}
	}
}

func validate(req interface{}) error {
	v9 = validator.New()
	v9.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("fieldName"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	if err := v9.Struct(req); err != nil {
		fieldError, oneofError, minError, defaultError := make([]string, 0), make([]string, 0), make([]string, 0), make([]string, 0)

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				fieldError = append(fieldError, err.Field())
			case "oneof":
				oneofError = append(oneofError, err.Field())
			case "min", "len", "lte", "max":
				minError = append(minError, err.Field())
			default:
				defaultError = append(defaultError, err.Field())
			}
		}

		if len(fieldError) > 0 {
			return &InputError{Message: fmt.Sprintf("Los siguientes campos son requeridos: %v", strings.Join(fieldError, ", "))}
		} else if len(oneofError) > 0 {
			return &InputError{Message: fmt.Sprintf("Los siguientes campos no hacen match con el enumerado: %v", strings.Join(oneofError, ", "))}
		} else if len(minError) > 0 {
			return &InputError{Message: fmt.Sprintf("Los siguientes campos no cumplen con la longitud de caracteres requeridos: %v", strings.Join(minError, ", "))}
		}
		return &InputError{Message: fmt.Sprintf("Los siguientes campos son inválidos: %v", strings.Join(defaultError, ", "))}
	}
	return nil
}
