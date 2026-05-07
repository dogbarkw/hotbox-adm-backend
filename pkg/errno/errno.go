package errno

import "fmt"

type ErrNo struct {
	statusCode      int32
	statusMessage   string
	suggestHttpCode int32
}

func (e *ErrNo) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.statusCode, e.statusMessage)
}

func (e *ErrNo) GetStatusCode() int32 {
	return e.statusCode
}

func (e *ErrNo) GetStatusMessage() string {
	return e.statusMessage
}

var (
	errMap = make(map[int32]*ErrNo)

	ErrOK = genErrNo(0, "OK", 200)
	// common error code
	ErrTokenInvalid     = genErrNo(1001, "token invalid", 401)
	ErrParamWrong       = genErrNo(1002, "param error", 400)
	ErrRpcFailed        = genErrNo(1003, "rpc failed", 500)
	ErrDBError          = genErrNo(1004, "db error", 500)
	ErrDBRecordNotFound = genErrNo(1005, "record not found", 400)

	ErrBccInvalidKey = genErrNo(1011, "bcc invalid key", 500)
)

func genErrNo(code int32, message string, httpCode int32) *ErrNo {
	e := &ErrNo{
		statusCode:      code,
		statusMessage:   message,
		suggestHttpCode: httpCode,
	}
	errMap[code] = e
	return e
}
