package errors

import "fmt"

type WrapperError struct {
	msg string
	err error
}

func Wrap(err error, msg string, arg ...any) WrapperError {
	if len(arg) > 0 {
		msg = fmt.Sprintf(msg, arg...)
	}

	if err != nil {
		msg += ": " + err.Error()
	}

	return WrapperError{
		msg: msg,
		err: err,
	}
}

func (e WrapperError) Error() string {
	return e.msg
}

func (e WrapperError) Unwrap() error {
	return e.err
}
