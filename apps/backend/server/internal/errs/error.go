package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Kind Kind
	Err  error  // 底层错误
	Msg  string // 错误展示消息
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *Error) Unwrap() error { return e.Err }

func New(kind Kind, msg string) *Error {
	return &Error{
		Kind: kind,
		Msg:  msg,
	}
}

func Wrap(kind Kind, msg string, err error) *Error {
	return &Error{Kind: kind, Msg: msg, Err: err}
}

func KindOf(err error) Kind {
	if e, ok := errors.AsType[*Error](err); ok {
		return e.Kind
	}
	return KindInternal
}

func IsKind(err error, kind Kind) bool {
	return KindOf(err) == kind
}
