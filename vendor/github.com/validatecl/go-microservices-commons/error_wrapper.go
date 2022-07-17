package commons

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

//DOCSURL variable para setear DOCS de response
var DOCSURL string = ""

//Error2WrapperFunc tipo de funcion error2Wrapper
type Error2WrapperFunc func(err error) (status int, errBody interface{})

//Error2WrapperMiddleware middleware de error2WrapperFunc
type Error2WrapperMiddleware func(Error2WrapperFunc) Error2WrapperFunc

//ErrorWrapper wrapper de response error
type ErrorWrapper struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Docs    []string `json:"docs,omitempty"`
}

//ErrorWrapperWM wrapper de response error WM
type ErrorWrapperWM struct {
	Success bool       `json:"success"`
	Errors  []*ErrItem `json:"errors"`
	Docs    []*string  `json:"docs,omitempty"`
}

//NewErrorWrapper crea un nuevo error wrapper
func NewErrorWrapper(code, message string, docs ...string) ErrorWrapper {

	if len(docs) < 1 {
		docs = []string{DOCSURL}
	}

	return ErrorWrapper{
		Code:    code,
		Message: message,
		Docs:    docs,
	}
}

//NewErrorWrapperWM crea un nuevo error wrapper
func NewErrorWrapperWM(s bool, e []*ErrItem, docs ...*string) ErrorWrapperWM {
	//if len(docs) < 1 {
	//	docs = []*string{StringPointer(DOCSURL)}
	//}

	return ErrorWrapperWM{
		Success: s,
		Errors:  e,
		Docs:    docs,
	}
}

//Error2Wrapper Convierte error a status code y error wrapper
func Error2Wrapper(err error) (status int, errBody interface{}) {
	switch err.(type) {
	case *InputError:
		return http.StatusBadRequest, NewErrorWrapper("400", err.Error())
	case *json.UnmarshalTypeError:
		umErr := err.(*json.UnmarshalTypeError)
		return http.StatusBadRequest, UnmarshalError2Wrapper(umErr)
	case *json.SyntaxError:
		return http.StatusBadRequest, NewErrorWrapper("400", err.Error())
	case *NotFoundError:
		return http.StatusNotFound, NewErrorWrapper("404", err.Error())
	case *GatewayError:
		return http.StatusBadGateway, NewErrorWrapper("502", err.Error())
	case *BackendCodedError:
		beErr := err.(*BackendCodedError)
		return http.StatusInternalServerError, NewErrorWrapper(beErr.Code, beErr.Message)
	case *CustomError:
		ceErr := err.(*CustomError)
		return ceErr.StatusCode, NewErrorWrapper(ceErr.Code, ceErr.Message)
	case *WMError:
		ceErr := err.(*WMError)
		return ceErr.StatusCode, NewErrorWrapperWM(ceErr.Success, ceErr.Errors, ceErr.Docs...)
	default:
		return http.StatusInternalServerError, NewErrorWrapper("500", err.Error())
	}
}

// UnmarshalError2Wrapper Convierte el error del Unmarshal a su equivalente HTTP
func UnmarshalError2Wrapper(err *json.UnmarshalTypeError) ErrorWrapper {

	errW := NewErrorWrapper("400", "")
	k := err.Type.Kind()
	switch k {
	case reflect.Map:
		errW.Message = fmt.Sprintf("Error field '%v' expected object not %v", err.Field, err.Value)
	case reflect.Slice, reflect.Array:
		errW.Message = fmt.Sprintf("Error field '%v' expected array not %v", err.Field, err.Value)
	case reflect.Struct:
		errW.Message = fmt.Sprintf("Error field '%v' expected object not %v", err.Field, err.Value)
	default:
		errW.Message = fmt.Sprintf("Error field '%v' expected %v not %v", err.Field, err.Type, err.Value)
	}

	return errW
}
