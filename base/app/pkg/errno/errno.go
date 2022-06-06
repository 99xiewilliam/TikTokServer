package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode         = 0
	ServiceErrCode      = 10001
	ParamErrCode        = 10002
	NeedLoginErrCode    = 10003
	TokenValidErrCode   = 10004
	TokenExpiredErrCode = 10005
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func New(code int64, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success    = New(SuccessCode, "Success")
	ServiceErr = New(ServiceErrCode, "Service is unable to start successfully")
	ParamErr   = New(ParamErrCode, "Wrong Parameter has been given")
)

func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
