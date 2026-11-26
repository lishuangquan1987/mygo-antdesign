package sample

import "strings"

// Extract 从嵌入的源码里提取 // sample:<name> … // /sample 标记段，
// 返回展示给用户的真实代码（示例即源码，永不漂移）。
// 缺少标记段时返回注释提示，而不是空串，便于发现遗漏。
func Extract(src, name string) string {
	begin := "// sample:" + name
	end := "// /sample"
	i := strings.Index(src, begin)
	if i < 0 {
		return "// 未找到标记段 " + begin
	}
	i += len(begin)
	j := strings.Index(src[i:], end)
	if j < 0 {
		return "// 标记段 " + begin + " 缺少结束标记"
	}
	seg := src[i : i+j]
	seg = strings.Trim(seg, "\n")
	if seg == "" {
		return "// 标记段 " + begin + " 为空"
	}
	return seg
}
