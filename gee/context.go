package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 用于json
type H map[string]interface{}

// Context, 封装HTTP请求和响应的处理逻辑
type Context struct {

	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string

	// response info
	StatusCode int // 记录当前响应的状态码
}

// Context的构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 获取POST请求中表单数据的值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取URL的查询参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 发送响应文本
func (c *Context) String(code int, format string, values ...interface{}) {

	// 设置响应头为 Content-Type: text/plain
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)

	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 发送json格式的响应
func (c *Context) JSON(code int, obj interface{}) {

	// 设置响应头为 Content-Type: text/plain
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)

	// 创建编码器
	encoder := json.NewEncoder(c.Writer)

	// 将obj编码为json格式发送
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 发送二进制数据, 用于文件、图片等
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 发送HTML响应
func (c *Context) HTML(code int, html string) {

	// 设置响应头为 Content-Type: text/plain
	c.SetHeader("Content-type", "text/plain")
	c.Status(code)

	c.Writer.Write([]byte(html))
}
