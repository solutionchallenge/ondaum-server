package http

import (
	"context"
	"fmt"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

var _ error = &Error{}

type Error struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(ctx context.Context, cause error, fmtstr string, args ...interface{}) *Error {
	requestID := utils.GetRequestID(ctx)
	utils.Log(utils.ErrorLevel).Ctx(ctx).Err(cause).RID(requestID).BT(1).Send(fmtstr, args...)
	return &Error{
		Cause:   cause,
		Message: fmt.Sprintf(fmtstr, args...),
	}
}
