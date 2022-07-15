package util

// HWOperatorType 宽高比较运算符
type HWOperatorType int32

const (
	GreaterOrEqualThan HWOperatorType = 0 // 大于等于
	LessOrEqualThan    HWOperatorType = 1 // 小于等于
)

// ErrorCode 错误码
type ErrorCode struct {
	Code int
	Msg  string
}

var (
	ErrParseReq = ErrorCode{
		Code: 1001,
		Msg:  "failed to parse req.",
	}
	ErrFilterImage = ErrorCode{
		Code: 1002,
		Msg:  "failed to filter image.",
	}
	ErrDownloadImage = ErrorCode{
		Code: 1003,
		Msg:  "failed to download image.",
	}
)
