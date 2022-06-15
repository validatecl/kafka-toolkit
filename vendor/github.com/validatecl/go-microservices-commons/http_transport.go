package commons

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

//MakeEncodeHTTPResponseFunc crea encode function
func MakeEncodeHTTPResponseFunc(err2wFunc Error2WrapperFunc) func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	encoderFunc := MakeServerErrorEncoderFunc(err2wFunc)
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			encoderFunc(ctx, f.Failed(), w)
			return nil
		}

		return json.NewEncoder(w).Encode(response)

	}
}

//MakeDefaultServerErrorEncoderFunc default error encoder
func MakeDefaultServerErrorEncoderFunc() func(ctx context.Context, err error, w http.ResponseWriter) {
	return MakeServerErrorEncoderFunc(Error2Wrapper)
}

//MakeServerErrorEncoderFunc funcion de encode de errores
func MakeServerErrorEncoderFunc(err2wFunc Error2WrapperFunc) func(ctx context.Context, err error, w http.ResponseWriter) {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		status, errorBody := err2wFunc(err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(errorBody)
	}
}
