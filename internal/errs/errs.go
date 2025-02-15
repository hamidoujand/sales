package errs

import (
	"fmt"
	"runtime"
)

type Error struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	FuncName string            `json:"-"`
	Filename string            `json:"-"`
	Fields   map[string]string `json:"fields,omitempty"`
}

func New(code int, err error) *Error {
	pc, filename, line, _ := runtime.Caller(1)
	file := fmt.Sprintf("%s:%d", filename, line)
	funcName := runtime.FuncForPC(pc).Name()

	return &Error{
		Code:     code,
		Message:  err.Error(),
		FuncName: funcName,
		Filename: file,
	}
}

func Newf(code int, format string, args ...any) *Error {
	pc, filename, line, _ := runtime.Caller(1)
	file := fmt.Sprintf("%s:%d", filename, line)
	funcName := runtime.FuncForPC(pc).Name()

	msg := fmt.Errorf(format, args...)
	return &Error{
		Code:     code,
		Message:  msg.Error(),
		FuncName: funcName,
		Filename: file,
	}
}

func (e *Error) Error() string {
	return e.Message
}
