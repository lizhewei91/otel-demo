package response

// ErrResponse Err response 数据结构
type ErrResponse struct {
	Code    string      `json:"code"`                                // 错误代码 (以 错误模块_错误类型_ERR 大写英文字符进行标识)
	Message string      `json:"message"`                             // 提示信息 (一个人类可读的错误信息)
	Data    interface{} `json:"data,omitempty" swaggertype:"object"` // 错误详细数据 (可选)
}

func (s *ErrResponse) Error() string {
	return s.Message
}

func NewSvcError(code string, message string, data ...interface{}) error {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}

	return &ErrResponse{Code: code, Message: message, Data: responseData}
}
