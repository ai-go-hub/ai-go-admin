package tree

import (
	"strings"
	"unicode/utf8"

	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

// Build 根据 id 与 pid 字段将扁平数据组装为带 children 的树状结构
//
// data: 扁平数据列表，每个元素为 map
// idField: 主键字段名
// pidField: 父级字段名
// childrenField: children 字段名
func Build(data []map[string]any, idField, pidField, childrenField string) []map[string]any {
	// 创建 id map，子找父模式实现（子级找父级，并挂到父级 children 数组中）
	idMap := make(map[any]map[string]any, len(data))
	for _, item := range data {
		idMap[item[idField]] = item
	}

	roots := make([]map[string]any, 0)
	for _, item := range data {
		id := item[idField]
		pid, hasPid := item[pidField]
		if !hasPid || util.IsZero(pid) || pid == id {
			roots = append(roots, item)
			continue
		}

		parent, ok := idMap[pid]
		if !ok {
			roots = append(roots, item)
			continue
		}

		children, _ := parent[childrenField].([]map[string]any)
		parent[childrenField] = append(children, item)
	}

	return roots
}

// Render 根据 id / pid 关系在 title 字段上渲染树状分支符号
//
// 输出仍然是平铺列表（不生成 children 嵌套）
// 顺序为前序遍历顺序：父级在前，子级紧随其后，同一父级下的子级保持输入时的原顺序
//
// idField: 主键字段名
// pidField: 父级字段名
// titleField: 需要渲染树状分支符号的字段名
func Render(data []map[string]any, idField, pidField, titleField string) []map[string]any {
	if len(data) == 0 {
		return make([]map[string]any, 0)
	}

	childrenMap := make(map[any][]map[string]any)

	roots := make([]map[string]any, 0)
	for _, item := range data {
		id := item[idField]
		pid, hasPid := item[pidField]
		if !hasPid || util.IsZero(pid) || pid == id {
			roots = append(roots, item)
			continue
		}
		childrenMap[pid] = append(childrenMap[pid], item)
	}

	result := make([]map[string]any, 0, len(data))
	var traverse func(items []map[string]any, depth int, isLasts []bool)
	traverse = func(items []map[string]any, depth int, isLasts []bool) {
		for i, item := range items {
			isLast := i == len(items)-1
			item[titleField] = treePrefix(depth, append(isLasts, isLast)) + cleanTreeTitle(item[titleField].(string))

			result = append(result, item)

			id := item[idField]
			if children, ok := childrenMap[id]; ok && len(children) > 0 {
				// 根节点的 isLast 不影响子级前缀，因此从根往下时重置为空
				if depth == 0 {
					traverse(children, depth+1, []bool{})
				} else {
					traverse(children, depth+1, append(isLasts, isLast))
				}
			}
		}
	}

	traverse(roots, 0, []bool{})
	return result
}

// treePrefix 根据层级深度生成树状前缀
// depth: 当前节点深度，根节点为 0
// isLasts: 记录了每级是否为最后一个子级
func treePrefix(depth int, isLasts []bool) string {
	if depth == 0 {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < depth-1; i++ {
		if isLasts[i] {
			sb.WriteString("     ") // 末级祖先：5 个空格
		} else {
			sb.WriteString("    │") // 非末级祖先：4 空格 + │
		}
	}
	if isLasts[depth-1] {
		sb.WriteString("    └") // 当前为末级：4 空格 + └
	} else {
		sb.WriteString("    ├") // 当前非末级：4 空格 + ├
	}
	return sb.String()
}

// cleanTreeTitle 清理 title 前导的空白与树状分支符号，返回干净的标题
func cleanTreeTitle(title string) string {
	branchChars := "├└│─┼┤┴┬┌┐┘┛┗┏┓"
	i := 0

	for i < len(title) {
		r, size := utf8.DecodeRuneInString(title[i:])
		if r == ' ' || r == '\t' || strings.ContainsRune(branchChars, r) {
			i += size
		} else {
			break
		}
	}

	return strings.TrimLeft(title[i:], " \t")
}
