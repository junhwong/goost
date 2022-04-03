package result

import (
	"fmt"
	"net/http"

	"github.com/junhwong/goost/errors"
	"github.com/junhwong/goost/validate"
	"github.com/junhwong/goost/web"
)

type RenderResult interface {
	GetError() error
}

type resetfulResult struct {
	// Total      int64                             `json:"total,omitempty"`
	// Count      int64                             `json:"count,omitempty"`
	Data interface{} `json:"data,omitempty"`
	// Message    string                            `json:"message,omitempty"`
	// Error      bool                              `json:"-"`
	Details    interface{}                       `json:"details,omitempty"`
	statusCode int                               `json:"-"`
	filters    []func(web.Context, RenderResult) `json:"-"`
	failed     bool
	err        error
	errCode    string
}

func (r *resetfulResult) GetError() error {
	return r.err
}
func (r *resetfulResult) Render(ctx web.Context) error {
	if r.statusCode < 1 {
		r.statusCode = http.StatusOK
		if r.failed {
			r.statusCode = http.StatusInternalServerError
		}
	}
	for _, filter := range r.filters {
		filter(ctx, r)
	}
	if ctx.IsAborted() {
		return nil
	}
	if r.err != nil {
		fmt.Println("TODO 完成错误打印", r.err)
		// return r.err
	}
	if r.Data == nil {
		ctx.Status(r.statusCode)
		return nil
	}
	ctx.JSON(r.statusCode, r.Data) // TODO 完成协议
	return nil
}

// func Empty(options ...Option) web.Result {
// 	result := &resetfulResult{
// 		statusCode: http.StatusOK,
// 	}
// 	for _, opt := range options {
// 		opt.apply(result)
// 	}
// 	return result
// }
// func List(count int64, data interface{}, options ...Option) web.Result {
// 	result := &resetfulResult{
// 		// Total:      total,
// 		// Count:      count,
// 		Data:       data,
// 		statusCode: http.StatusOK,
// 	}
// 	for _, opt := range options {
// 		opt.apply(result)
// 	}
// 	return result
// }

func Result(status int, data interface{}, options ...Option) web.Result {
	result := &resetfulResult{
		statusCode: status,
		Data:       data,
	}
	for _, opt := range options {
		opt.apply(result)
	}
	return result
}

type SuccessOption interface {
	applySuccessOption()
}

func Success(data interface{}, options ...Option) web.Result {
	result := &resetfulResult{
		Data:       data,
		statusCode: http.StatusOK,
	}
	for _, opt := range options {
		opt.apply(result)
	}
	return result
}

type FailOption interface {
	apply(*resetfulResult)
}
type failOptionFunc func(*resetfulResult)

func (f failOptionFunc) apply(opt *resetfulResult) {
	f(opt)
}
func Fail(message string, options ...FailOption) web.Result {
	result := &resetfulResult{
		statusCode: http.StatusInternalServerError,
		failed:     true,
	}
	for _, opt := range options {
		opt.apply(result)
	}
	result.Data = errors.ErrorResult{
		Message: message,
		Code:    result.errCode,
	}
	return result
}
func Error(err error, options ...FailOption) web.Result {
	result := &resetfulResult{
		statusCode: http.StatusInternalServerError,
		failed:     true,
		err:        err,
	}
	if res, ok := err.(errors.ErrorResponse); ok {
		result.statusCode = res.ResponseStatusCode()
		result.Data = res.ResponseData()
	} else if err := validate.AsValidationError(err); err != nil {
		result.statusCode = err.ResponseStatusCode()
		result.Data = err.ResponseData()
	}

	for _, opt := range options {
		opt.apply(result)
	}
	return result
}

type Option interface {
	apply(*resetfulResult)
}
type optionFunc func(*resetfulResult)

func (f optionFunc) apply(opt *resetfulResult) {
	f(opt)
}

// 包含HTTP状态码
func WithStatus(code int) optionFunc {
	return func(opt *resetfulResult) {
		opt.statusCode = code
	}
}

func WithCode(c string) FailOption {
	return failOptionFunc(func(rr *resetfulResult) {
		rr.errCode = c
	})
}

// // 包含错误明细
// func WithDetails(details ...ErrorDetail) Option {
// 	return &resultOption{
// 		runner: func(target *resetfulResult) {
// 			target.Details = details
// 		},
// 	}
// }

// type resultOption struct {
// 	runner func(target *resetfulResult)
// }

// func (opt *resultOption) apply(target *resetfulResult) {
// 	opt.runner(target)
// }
