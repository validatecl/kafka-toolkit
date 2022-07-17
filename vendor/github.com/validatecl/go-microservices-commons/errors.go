package commons

import (
	"fmt"
)

// StringPointer .
func StringPointer(s string) *string {
	return &s
}

//InputError error de input de services
type InputError struct {
	Message string
}

func (i *InputError) Error() string {
	return fmt.Sprintf("Input inválido: %v", i.Message)
}

//ServiceError error de services
type ServiceError struct {
	Message string
}

func (s *ServiceError) Error() string {
	return fmt.Sprintf("Error de service : %v", s.Message)
}

//GatewayError error de service externo
type GatewayError struct {
	Message string
	Cause   error
}

func (g *GatewayError) Error() string {
	return fmt.Sprintf("Error de Gateway : %v", g.Message)
}

// BadGatewayError error de service externo
type BadGatewayError struct {
	Message string
	Cause   error
}

func (b *BadGatewayError) Error() string {
	return fmt.Sprintf("Error de BadGateway : %v", b.Message)
}

//BackendCodedError errores codificados de backend
type BackendCodedError struct {
	Code    string
	Message string
}

func (g *BackendCodedError) Error() string {
	return fmt.Sprintf("Error de backend code: %q, message: %q", g.Code, g.Message)
}

//NotFoundError errores codificados de backend
type NotFoundError struct {
	Code    string
	Message string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("Error de NotFound : %v", n.Message)
}

//CustomError errores customizables
type CustomError struct {
	Code       string
	Message    string
	StatusCode int
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("Error : %v", c.Message)
}

type ErrItem struct {
	Msg      *string `json:"msg,omitempty"`
	Param    *string `json:"param,omitempty"`
	Value    *string `json:"value,omitempty"`
	Location *string `json:"location,omitempty"`
}

// WMError .
type WMError struct {
	Success    bool
	StatusCode int
	Message    string
	Errors     []*ErrItem `json:"errors,omitempty"`
	Docs       []*string  `json:"docs,omitempty"`
}

// WMError .
func (wm *WMError) Error() string {
	if wm.Message != "" {
		return wm.Message
	}
	return "an error has occurred"
}
