package errors

type Error struct {
	Message string // 错误信息
	Raise   error  // 由什么错误引起的
	File    string // 源码文件名
	Line    int    // 源码文件行号
	Stack   []byte // 原始栈
	Code    string // 错误码, 用于识别错误并处理
}

func (err *Error) Error() string {
	switch {
	case err.Message != "" && err.Code != "":
		return err.Code + ":" + err.Message
	case err.Message != "":
		return err.Message
	case err.Raise != nil:
		return err.Raise.Error()
	case err.Code != "":
		return err.Code
	}
	return ""
}

// MarshalJSON impl json.Marshaller to starand-logs error fields
func (err *Error) MarshalJSON() ([]byte, error) {
	return nil, nil
}
