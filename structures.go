package recipe

import (
	"net/http"
	"runtime"

	"github.com/uber-go/zap"
)

// Response is used to serialize the data to return.
type Response struct {
	Data       interface{} `json:"data"`
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
}

// ErrorResponse wraps an error in the right response format.
type ErrorResponse struct {
	Message    string `json:"message"`
	Reference  string `json:"reference"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

// NewErrorResponse creates a new error response message
func NewErrorResponse(msg string, code int, err ...error) *ErrorResponse {
	pc, file, line, _ := runtime.Caller(1)
	Logger.Debug("error response",
		zap.String("handlerName", runtime.FuncForPC(pc).Name()),
		zap.String("file", file),
		zap.Int("line", line),
	)
	var e error
	if err != nil {
		e = err[0]
	}
	Logger.Error("error response",
		zap.String("handlerName", runtime.FuncForPC(pc).Name()),
		zap.Error(e),
	)
	return &ErrorResponse{
		Message:    msg,
		StatusCode: code,
		Status:     http.StatusText(code),
		Reference:  "https://httpstatuses.com/",
	}
}

// Error implement the Error interface
func (e ErrorResponse) Error() string {
	return e.Message
}

// NewResponse creates a new response
func NewResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}

// Health conatins the heakth metric of the service
type Health struct {
	Status int
}
