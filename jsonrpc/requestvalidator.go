package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// A JSON-RPC 2.0 request type that needs validation should implement this interface.
//
// Note:
// 1. Request validation should NOT use information other than the request itself.
// 2. Request validation is server-implementation independent. It is the "lower bound" of any
//    implemenation. Additional checking logic should go to the RPC method implementation.
// 3. Usually raw data from the request should not be put into `data` field of the error - usually
//    that can be retrieved from RPC logs or so.
type RequestValidator interface {
	Validate() *InvalidParamsError
}

type InvalidParamsError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewInvalidParamsError(message string) *InvalidParamsError {
	return NewInvalidParamsErrorWithData(message, nil)
}

func NewInvalidParamsErrorWithData(message string, data interface{}) *InvalidParamsError {
	return &InvalidParamsError{
		Code:    -32602,
		Message: message,
		Data:    data,
	}
}

func (e *InvalidParamsError) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"code":%d,"message":%s,"data:(failed to marshal)"}`, e.Code, e.Message)
	}
	return string(bytes)
}
