package gee

import "strings"

// Trie树节点结构
type node struct {
	pattern  string  // 待匹配的完整路由, 只有在子节点中会设置
	part     string  // 路由中的一部分, 以/为分割符
	children []*node // 子节点
	isWild   bool    // 标记part是否为通配符, part中含有 : 或 * 时为true
}

// 从子节点中找到第一个匹配的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 从子节点中找出所有匹配的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// @brief Trie的插入操作, 对应注册路由
// @param pattern 完整的路由路径
// @param parts   路径的各个部分
// @param height  当前递归到parts的第几层, 用于控制遍历路径的深度
func (n *node) insert(pattern string, parts []string, height int) {

	// 已经遍历到子节点, 将pattern赋值给n
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 继续处理当前part
	part := parts[height]

	// 若当前part在子节点中不存在, 那么创建一个
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	// 递归处理剩余部分
	child.insert(pattern, parts, height+1)
}

// @brief Trie的查找操作, 对应匹配路由
// @param parts  待匹配路径
// @param height 当前递归深度
// @return *node 返回匹配的节点
func (n *node) search(parts []string, height int) *node {

	// 若当前节点是通配符或者已经到底, 那么返回
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 若当前节点的pattern为空, 说明路径不完整或者不匹配, 返回空
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 继续处理当前part
	part := parts[height]

	// 找到所有匹配的路径
	children := n.matchChildren(part)

	// 向下递归查找
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
