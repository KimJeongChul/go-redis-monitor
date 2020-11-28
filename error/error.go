package error

import "strconv"

type ErrCode int

const (
	// OS_OPEN_ERR cannot open file os.Open() error
	OS_OPEN_ERR = iota

	// JSON_UNMARSHAL_ERR JSON decoding json.Unmarshal() error
	JSON_UNMARSHAL_ERR

	// REDIS_DB_ERR Redis DB error
	REDIS_DB_ERR
)

type CError struct {
	Code    ErrCode
	Message string
}

func NewCError(code ErrCode, msg string) *CError {
	return &CError{
		Code:    code,
		Message: msg,
	}
}

func (e *CError) Error() string {
	return "CODE:" + strconv.Itoa(int(e.Code)) + ", MSG:" + e.Message
}
