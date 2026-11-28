// Package sample 是控件示例的注册表与工具（见 docs/07-samples-gallery.md）。
// 每个控件在自己的示例包（如 samples/button）里提供 Demo()，由 Gallery
// 注册并展示；示例代码用 go:embed 嵌入真实源码，保证「代码即示例」。
package sample

import (
	"slices"

	"github.com/egoist/mygo/ui"
)

// Category 示例的分类（Gallery 侧栏按此分组）。
type Category string

const (
	CatGeneral  Category = "通用"
	CatInput    Category = "数据录入"
	CatDisplay  Category = "数据展示"
	CatNav      Category = "导航"
	CatFeedback Category = "反馈"
	CatExtra    Category = "缺口控件"
	CatNative   Category = "原生补全"
)

// Sample 一个示例单元：标题 + 说明 + 可复制的真实代码 + 实时预览。
type Sample struct {
	Title       string
	Description string
	// Code 是从真实源码提取的片段（// sample:xxx 标记段），非手写字符串。
	Code string
	// View 是实时预览的构建函数，由 Gallery 在主题子树内调用。
	View func(c *ui.Context)
	// Fill 标记布局型示例（如 IDE 风格）：Gallery 与截图工具中
	// 撑满可用区域而非居中，避免四周留白。
	Fill bool
	// Page 目标控件页名（如 "Button"）；留空时沿用所属分组的 Name，
	// 或 pages.go 的 pageOverrides 映射（旧式多控件片段）。
	Page string
}

// WidgetDemo 一个控件的示例集合。
type WidgetDemo struct {
	Name     string // 控件名，如 "Button"
	Package  string // "widgets" 或 "extra"
	Category Category
	Summary  string
	Samples  []Sample
}

var registry []WidgetDemo

// Register 登记一个控件的示例集合（在 Gallery 启动时调用）。
// 控件名重复属于集成错误（两个包用了同一个 Name），直接 panic，
// 避免静默追加导致侧栏出现两个同名页。
func Register(d WidgetDemo) {
	for _, e := range registry {
		if e.Name == d.Name {
			panic("sample: 控件 " + d.Name + " 重复注册（检查 samples/*/demo.go 的 Name）")
		}
	}
	registry = append(registry, d)
}

// All 返回全部已注册示例的快照（调用方修改返回切片不影响注册表）。
func All() []WidgetDemo {
	return slices.Clone(registry)
}
