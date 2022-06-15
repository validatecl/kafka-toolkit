package commons

import (
	"encoding/json"
	"net/http"
)

// Error2Status Convierte errores en status HTTP
func Error2Status(err error) int {
	switch err.(type) {
	case *InputError:
		return http.StatusBadRequest
	case *json.UnmarshalTypeError:
		return http.StatusBadRequest
	case *json.SyntaxError:
		return http.StatusBadRequest
	case *NotFoundError:
		return http.StatusNotFound
	case *GatewayError:
		return http.StatusBadGateway
	case *BackendCodedError:
		return http.StatusInternalServerError
	case *CustomError:
		ceErr := err.(*CustomError)
		return ceErr.StatusCode
	}

	return http.StatusInternalServerError
}
