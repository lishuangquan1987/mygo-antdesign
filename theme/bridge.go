package theme

import "github.com/egoist/mygo/ui"

// toTheme 把 AntToken 映射为 ui.Theme（docs/04-theme-system.md 的映射表）。
// Accent 相关色阶（hover、pressed、选中、焦点环、文字）由 MyGo 的
// WithAccent 派生，保证与品牌色成套。
func (t *AntToken) toTheme() *ui.Theme {
	base := ui.LightTheme()
	if t.isDark {
		base = ui.DarkTheme()
	}
	th := *base

	// —— 语义色映射 ——
	th.Accent = t.Primary
	th.Danger = t.Error
	th.Warning = t.Warning
	th.Success = t.Success
	th.Text = t.Text
	th.TextMuted = t.TextSecondary
	// TextInverse 是 Tooltip/Toast 上的文字；零值保持 MyGo 默认
	// （Inverse = Text、InverseText = Background）。
	if t.TextInverse != (ui.Color{}) {
		th.InverseText = t.TextInverse
	}
	th.Background = t.BgLayout
	th.Surface = t.BgContainer
	th.Border = t.Border

	// Surface 的 hover / pressed 由容器底色与文本色混合派生
	// （antd colorFillSecondary ≈ 文本色 6% / 12% 的叠加）。
	th.SurfaceHover = t.BgContainer.Mix(t.Text, 0.06)
	th.SurfacePressed = t.BgContainer.Mix(t.Text, 0.12)

	// —— 度量映射 ——
	th.Radius = t.Radius
	th.FontSize = t.FontSize
	th.Spacing = t.Spacing

	// 静态渲染场景禁用焦点环（见 WithoutFocusRing）。
	if t.noFocusRing {
		th.Focus = ui.Transparent
	}

	// —— 由主色派生成套色阶 ——
	// WithAccent 生成 AccentHover / AccentPressed / Selection / Focus /
	// AccentText（暗主题下色阶整体更亮）。Info 无独立槽位，跟随主色。
	return th.WithAccent(t.Primary)
}
