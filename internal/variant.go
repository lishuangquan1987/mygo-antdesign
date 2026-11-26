package internal

import "github.com/egoist/mygo/ui"

// ButtonPalette 是一种按钮形态在各交互状态下的颜色组合。
// 全部取值来自主题，禁止硬编码（AGENTS.md 2.1：状态样式集中管理）。
type ButtonPalette struct {
	Bg, BgHover, BgPressed             ui.Color
	Fg, FgHover, FgPressed             ui.Color
	Border, BorderHover, BorderPressed ui.Color
	BorderWidth                        float32
	Dashed                             bool
}

// ButtonColors 返回 antd Button 的色板。typ 为 antd 的 type
// （default/primary/dashed/text/link），danger 为危险变体。
// 状态取值规则对齐 antd v5：
//
//   - primary：实底主色，hover/pressed 用 WithAccent 派生的色阶；
//     danger 实底由 Danger 派生（无独立 hover 槽位，用 Mix）。
//   - default / dashed：白底描边，hover 时文字与描边变主色（antd 行为）。
//   - text：透明底，hover 出浅填充（SurfaceHover）；danger 变红字。
//   - link：透明底主色文字，hover 用主色 hover 阶。
func ButtonColors(t *ui.Theme, typ string, danger bool) ButtonPalette {
	p := ButtonPalette{}

	switch typ {
	case "primary":
		if danger {
			// 危险实底：antd hover 提亮、pressed 加深；此处用 Mix 派生。
			p.Bg, p.BgHover, p.BgPressed = t.Danger, lighten(t, t.Danger), darken(t, t.Danger)
			p.Fg, p.FgHover, p.FgPressed = t.AccentText, t.AccentText, t.AccentText
			p.Border, p.BorderHover, p.BorderPressed = p.Bg, p.BgHover, p.BgPressed
		} else {
			p.Bg, p.BgHover, p.BgPressed = t.Accent, t.AccentHover, t.AccentPressed
			p.Fg, p.FgHover, p.FgPressed = t.AccentText, t.AccentText, t.AccentText
			p.Border, p.BorderHover, p.BorderPressed = p.Bg, p.BgHover, p.BgPressed
		}
		p.BorderWidth = 1

	case "text":
		p.Bg, p.BgHover, p.BgPressed = ui.Transparent, t.SurfaceHover, t.SurfacePressed
		if danger {
			p.Fg, p.FgHover, p.FgPressed = t.Danger, lighten(t, t.Danger), darken(t, t.Danger)
		} else {
			p.Fg = t.Text
			p.FgHover, p.FgPressed = t.Text, t.Text
		}
		p.BorderWidth = 0

	case "link":
		p.Bg, p.BgHover, p.BgPressed = ui.Transparent, ui.Transparent, ui.Transparent
		if danger {
			p.Fg, p.FgHover, p.FgPressed = t.Danger, lighten(t, t.Danger), darken(t, t.Danger)
		} else {
			p.Fg, p.FgHover, p.FgPressed = t.Accent, t.AccentHover, t.AccentPressed
		}
		p.BorderWidth = 0

	default: // default / dashed
		p.Bg, p.BgHover, p.BgPressed = t.Surface, t.Surface, t.Surface
		if danger {
			// 危险描边：红字红边，hover 出浅红底（antd colorErrorBg）。
			p.Fg, p.FgHover, p.FgPressed = t.Danger, lighten(t, t.Danger), darken(t, t.Danger)
			p.Border, p.BorderHover, p.BorderPressed = t.Danger, p.FgHover, p.FgPressed
			p.BgHover, p.BgPressed = t.Danger.Mix(t.Surface, 0.92), t.Danger.Mix(t.Surface, 0.85)
		} else {
			// antd：hover 时文字与描边变主色，底色不变。
			p.Fg, p.FgHover, p.FgPressed = t.Text, t.Accent, t.AccentPressed
			p.Border, p.BorderHover, p.BorderPressed = t.Border, t.Accent, t.AccentPressed
		}
		p.BorderWidth = 1
		p.Dashed = typ == "dashed"
	}
	return p
}

// lighten/danger 系的 hover 与 pressed 派生：亮色主题下 hover 提亮、
// pressed 加深；暗色主题相反（antd dark 的实底按钮 hover 更亮，这里
// 保持与亮色一致的简化处理，视觉校准在 Gallery 截图阶段微调）。
func lighten(t *ui.Theme, c ui.Color) ui.Color {
	if t.Dark {
		return c.Mix(ui.RGB(255, 255, 255), 0.12)
	}
	return c.Mix(ui.RGB(255, 255, 255), 0.12)
}

func darken(t *ui.Theme, c ui.Color) ui.Color {
	if t.Dark {
		return c.Mix(ui.RGB(255, 255, 255), 0.06)
	}
	return c.Mix(ui.RGB(0, 0, 0), 0.12)
}
