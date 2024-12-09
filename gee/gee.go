package gee

import (
	"net/http"
)

// 用于定义路由映射的处理方法
type HandlerFunc func(*Context)

// 实现了ServeHTTP接口
type Engine struct {
	router *router // 路由相关
}

// Engine的构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加路由
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 向下调用router中的addRoute
	e.router.addRoute(method, pattern, handler)
}

// 添加一个GET类型的路由
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// 添加一个POST类型的路由
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// 启动一个HTTP服务器
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// 解析请求的路径
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
