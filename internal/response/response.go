package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
	Data    any    `json:"data,omitempty"`
}

// ---------- Functional Options ----------

// Option 响应选项函数
type Option func(*Response)

// WithCode 自定义业务码
func WithCode(code int) Option {
	return func(r *Response) { r.Code = code }
}

// WithMessage 自定义消息
func WithMessage(msg string) Option {
	return func(r *Response) { r.Message = msg }
}

// WithData 附带数据
func WithData(data any) Option {
	return func(r *Response) { r.Data = data }
}

// Success 成功响应，默认 code = 0
func Success(c *gin.Context, opts ...Option) {
	r := &Response{Code: 0, Message: "ok", Time: time.Now().Unix()}
	for _, opt := range opts {
		opt(r)
	}
	c.JSON(http.StatusOK, r)
}

// Fail 失败响应，默认 code = 1
func Fail(c *gin.Context, opts ...Option) {
	r := &Response{Code: 1, Time: time.Now().Unix()}
	for _, opt := range opts {
		opt(r)
	}
	c.JSON(http.StatusOK, r)
}

// ---------- 链式调用 ----------

// Writer 链式响应构建器
type Writer struct {
	c    *gin.Context
	resp *Response
}

// New 开启链式构建，默认 code = 0
func New(c *gin.Context) *Writer {
	return &Writer{c: c, resp: &Response{Code: 0, Time: time.Now().Unix()}}
}

// Code 自定义业务码
func (w *Writer) Code(code int) *Writer {
	w.resp.Code = code
	return w
}

// Message 自定义消息
func (w *Writer) Message(msg string) *Writer {
	w.resp.Message = msg
	return w
}

// Data 附带数据
func (w *Writer) Data(data any) *Writer {
	w.resp.Data = data
	return w
}

// Send 发送响应
func (w *Writer) Send() {
	w.c.JSON(http.StatusOK, w.resp)
}
