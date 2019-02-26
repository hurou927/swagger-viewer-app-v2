package common

import (
	"fmt"
	"runtime"
)

type Error struct {
	Message  string `json:"message"`
	Code     int    `json:"code"`
	Function string `json:"functionName"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Err      error  `json:"error"`
}

func NewError(code int, message string, err error) *Error {

	pc, file, line, _ := runtime.Caller(1)
	return &Error{
		Message:  message,
		Code:     code,
		Function: runtime.FuncForPC(pc).Name(),
		File:     file,
		Line:     line,
		Err:      err,
	}
}

func (err *Error) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("ERROR::%v::%v::%v::%v::%v::%v", err.File, err.Function, err.Line, "unspecified error", err.Code, err.Message)
	}
	return fmt.Sprintf("ERROR::%v::%v::%v::%v::%v::%v", err.File, err.Function, err.Line, err.Err.Error(), err.Code, err.Message)
}
