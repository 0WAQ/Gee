package gee

import (
	"log"
	"net/http"
	"strings"
)

// router
// roots key: roots['GET'] roots['POST']
// handlers key: handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node       // 存储每种请求方式的根节点
	handlers map[string]HandlerFunc // 路由映射表, 映射路径与处理函数
}

// router的构造函数
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// @brief 解析pattern为parts, 并且路径中只允许一个*号
func parsePattern(pattern string) []string {
	// 以 / 分割pattern
	vs := strings.Split(pattern, "/")

	// 将分割后的[]string加入到parts
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)

			// 只允许出现一次*号
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)

	// 若不存在当前方法, 那么创建一个
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	// 解析pattern
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// @brief 匹配路由
// @param method HTTP方法
// @param path   待匹配的请求路径
// @return *node 匹配到的路由节点
// @return map[string]string 匹配到的路径参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil // 若找不到, 说明当前方法没有注册路由
	}

	searchParts := parsePattern(path) // 解析路径
	n := root.search(searchParts, 0)  // 查找匹配的节点
	if n == nil {
		return nil, nil
	}

	parts := parsePattern(n.pattern)  // 解析节点的parts
	params := make(map[string]string) // 用于存储动态参数的解析结果

	// 遍历并解析路径参数
	for idx, part := range parts {

		// :开头, 例如 :lang表示动态参数, 键为lang
		if part[0] == ':' {
			params[part[1:]] = searchParts[idx]
		}

		// *开头, 例如 *filepath, 键为filepath
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[idx:], "/")
			break
		}
	}
	return n, params
}

// 路由执行函数
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	} else {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	}
}
