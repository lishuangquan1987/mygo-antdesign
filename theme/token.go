// Package theme 提供 Ant Design 语义 Token，并桥接到 MyGo 的 ui.Theme。
//
// 控件（widgets/ 与 extra/）只从 c.Theme() 读取颜色与度量；AntToken 是
// 构建 ui.Theme 的唯一输入。见 docs/04-theme-system.md。
package theme

import "github.com/egoist/mygo/ui"

// Density 控件密度。
type Density string

const (
	// Comfortable 是 Ant Design 默认密度（Theme.Spacing = 4）。
	Comfortable Density = "comfortable"
	// Compact 是紧凑密度（Theme.Spacing = 3），用于上位机、IDE 等
	// 信息密集场景。见 docs/05-compact-density.md。
	Compact Density = "compact"
)

// AntToken 是 Ant Design 语义 Token 的常用子集。字段名对齐 antd 的
// design token 语义（colorPrimary、colorText、colorBgContainer…），
// 数值对齐 antd v5 默认主题。
//
// 零值 token 不可用；请用 Light 或 Dark 工厂创建。
type AntToken struct {
	// Primary 品牌主色（antd colorPrimary）。
	Primary ui.Color
	// Success / Warning / Error / Info 语义色。
	Success ui.Color
	Warning ui.Color
	Error   ui.Color
	Info    ui.Color

	// Text 主文本（antd colorText）；TextSecondary 次要文本。
	Text          ui.Color
	TextSecondary ui.Color
	// TextInverse 反色文本（Tooltip/Toast 上的文字）。零值表示跟随
	// MyGo 默认（Inverse = Text、InverseText = Background），antd 的
	// 反色行为与之一致，预设通常不设。
	TextInverse ui.Color

	// BgContainer 组件容器背景（antd colorBgContainer，即按钮、输入框的底色）。
	BgContainer ui.Color
	// BgLayout 页面底色（antd colorBgLayout）。
	BgLayout ui.Color
	// Border 描边（antd colorBorder）。
	Border ui.Color

	// Radius 全局圆角（antd borderRadius）。
	Radius float32
	// FontSize 基准字号（antd fontSize）。
	FontSize float32
	// Spacing MyGo 的留白单位：标准密度 4，紧凑 3。
	Spacing float32
	// Density 当前密度。
	Density Density

	// isDark 标记来自哪个工厂，桥接时选择 ui 基底主题。
	isDark bool
	// followSystem 为 true 时，Apply/Sub 会叠加桌面的强调色、高对比与
	// 文字缩放（见 ADR 0002）。
	followSystem bool
}

// Option 定制 AntToken。
type Option func(*AntToken)

// WithPrimary 设置品牌主色（antd colorPrimary）。色阶（hover、pressed、
// 选中、焦点环）由桥接层的 WithAccent 自动派生。
func WithPrimary(c ui.Color) Option {
	return func(t *AntToken) { t.Primary = c }
}

// WithDensity 设置密度。紧凑档同时调整 Spacing、Radius 与 FontSize。
func WithDensity(d Density) Option {
	return func(t *AntToken) {
		t.Density = d
		switch d {
		case Compact:
			t.Spacing = 3
			t.Radius = 4
			t.FontSize = 13
		default:
			t.Spacing = 4
			t.Radius = 6
			t.FontSize = 14
		}
	}
}

// WithFontSize 覆盖基准字号（px）。
func WithFontSize(px float32) Option {
	return func(t *AntToken) { t.FontSize = px }
}

// WithRadius 覆盖全局圆角（px）。
func WithRadius(px float32) Option {
	return func(t *AntToken) { t.Radius = px }
}

// WithFollowSystem 让 Apply/Sub 叠加桌面偏好（强调色、高对比、文字缩放）。
// 默认关闭，以保证 Ant Design 品牌色忠实呈现（见 ADR 0002）。
func WithFollowSystem() Option { return func(t *AntToken) { t.followSystem = true } }

// WithSpacing 直接覆盖 Spacing 单位（高级用法；常规请用 WithDensity）。
func WithSpacing(px float32) Option {
	return func(t *AntToken) { t.Spacing = px }
}
