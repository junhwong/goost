package errors

type ErrorResponse interface {
	ResponseStatusCode() int
	ResponseData() interface{}
}
type ErrorResponseDetail interface {
	ResponseDetailData() ErrorDetail
}

type ErrorResult struct {
	Code    string        `json:"code,omitempty" xml:"code,omitempty"`
	Message string        `json:"message,omitempty" xml:"message,omitempty"`
	Details []ErrorDetail `json:"details,omitempty" xml:"details,omitempty"`
}

// 表示错误明细
type ErrorDetail struct {
	Code    string `json:"code,omitempty" xml:"code,omitempty"`
	Field   string `json:"field,omitempty" xml:"field,omitempty"`
	Message string `json:"message,omitempty" xml:"message,omitempty"`
}
