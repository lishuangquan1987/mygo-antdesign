// Package ant 提供 Ant Design 风格的 MyGo 原生 UI 控件。
//
// 目录名为 widgets/（见 docs/06-development-plan.md 第 3 节），包名 ant，
// 导入方式：import "github.com/lishuangquan1987/mygo-antdesign/widgets"。
//
// 所有控件遵守样式铁律：颜色只读 c.Theme()（或 internal 的集中映射），
// 间距只用 t.Space / t.Rem，状态样式集中在 internal/variant.go。
package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// ButtonType 对齐 antd Button 的 type。
type ButtonType string

const (
	DefaultButton ButtonType = "default"
	PrimaryType   ButtonType = "primary"
	DashedType    ButtonType = "dashed"
	TextType      ButtonType = "text"
	LinkType      ButtonType = "link"
)

// SizeKind 控件尺寸档（antd size）。提供 Small / Middle / Large 常量。
type SizeKind = internal.SizeKind

const (
	Small  = internal.SizeSmall
	Middle = internal.SizeMiddle
	Large  = internal.SizeLarge
)

// ButtonOption 定制 Button。
type ButtonOption func(*buttonConfig)

type buttonConfig struct {
	typ      ButtonType
	danger   bool
	loading  bool
	disabled bool
	block    bool
	size     internal.SizeKind
	onClick  func()
}

// Primary 主按钮（antd type="primary"）。
func Primary() ButtonOption { return func(c *buttonConfig) { c.typ = PrimaryType } }

// Dashed 虚线按钮（antd type="dashed"）。
func Dashed() ButtonOption { return func(c *buttonConfig) { c.typ = DashedType } }

// Text 文字按钮（antd type="text"）。
func Text() ButtonOption { return func(c *buttonConfig) { c.typ = TextType } }

// Link 链接按钮（antd type="link"）。
func Link() ButtonOption { return func(c *buttonConfig) { c.typ = LinkType } }

// Danger 危险变体（antd danger）。
func Danger() ButtonOption { return func(c *buttonConfig) { c.danger = true } }

// Size 设置尺寸档。
func Size(s SizeKind) ButtonOption { return func(c *buttonConfig) { c.size = s } }

// Loading 加载态：禁点并显示加载指示（antd loading）。
func Loading(b bool) ButtonOption { return func(c *buttonConfig) { c.loading = b } }

// Disabled 禁用。
func Disabled(b bool) ButtonOption { return func(c *buttonConfig) { c.disabled = b } }

// Block 宽度撑满父容器（antd block）。
func Block(b bool) ButtonOption { return func(c *buttonConfig) { c.block = b } }

// OnClick 点击回调。MyGo v0.10.0 的 Element 只有 Clicked() 帧内查询，
// 这里用 ui.Local 记录上一帧状态、以上升沿触发，语义对齐 React 的
// onClick：一次点击只触发一次。
func OnClick(fn func()) ButtonOption { return func(c *buttonConfig) { c.onClick = fn } }

// Button 创建 Ant Design 风格的按钮，基于 MyGo 的 ButtonBase（无外观
// 底座，自带焦点、键盘 Enter/Space 与无障碍支持）：
//
//	ant.Button(c, "保存", ant.Primary())
//	ant.Button(c, "删除", ant.Danger(), ant.OnClick(app.remove))
func Button(c *ui.Context, label string, opts ...ButtonOption) ui.Element {
	cfg := &buttonConfig{typ: DefaultButton, size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}

	t := c.Theme()
	pal := internal.ButtonColors(t, string(cfg.typ), cfg.danger)
	disabled := cfg.disabled || cfg.loading

	el := ui.ButtonBase(c).
		Height(internal.ControlHeight(t, cfg.size)).
		Padding(0, internal.ControlPaddingX(t, cfg.size)).
		Radius(t.Radius).
		Center()

	if cfg.block {
		el.FillWidth()
	}

	// 即时模式状态：每帧查询 hover / pressed 后取对应颜色。
	bg, fg, border := pal.Bg, pal.Fg, pal.Border
	if !disabled {
		switch {
		case el.Pressed():
			bg, fg, border = pal.BgPressed, pal.FgPressed, pal.BorderPressed
		case el.Hovered():
			bg, fg, border = pal.BgHover, pal.FgHover, pal.BorderHover
		}
	}

	el.Background(bg)
	if pal.BorderWidth > 0 {
		el.Border(pal.BorderWidth, border)
		if pal.Dashed {
			el.BorderStyle(ui.BorderDashed)
		}
	}

	el.Children(func() {
		if cfg.loading {
			ui.Spinner(c)
		}
		ui.Text(c, label).
			TextColor(fg).
			FontSize(internal.ControlFont(t, cfg.size)).
			SingleLine()
	})

	el.Disabled(disabled)

	if cfg.onClick != nil {
		prev := ui.Local(el, "ant.button.clicked", func() bool { return false })
		cur := el.Clicked()
		if cur && !*prev && !disabled {
			cfg.onClick()
		}
		*prev = cur
	}
	return el
}
