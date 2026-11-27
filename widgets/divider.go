package ant

import (
	"github.com/egoist/mygo/ui"
)

// DividerOption 定制 Divider。
type DividerOption func(*dividerConfig)

type dividerConfig struct {
	dashed bool
	text   string
	align  ui.Align // 带文字时的文字位置
}

// Dashed 虚线分割线（antd dashed）。
func DashedDivider() DividerOption { return func(c *dividerConfig) { c.dashed = true } }

// WithText 带文字的分割线（antd Divider 的 children）。
func WithText(s string) DividerOption { return func(c *dividerConfig) { c.text = s } }

// TextLeft 文字靠左（antd orientation="left"，默认）。
func TextLeft() DividerOption { return func(c *dividerConfig) { c.align = ui.Start } }

// TextCenter 文字居中。
func TextCenter() DividerOption { return func(c *dividerConfig) { c.align = ui.Center } }

// TextRight 文字靠右（antd orientation="right"）。
func TextRight() DividerOption { return func(c *dividerConfig) { c.align = ui.End } }

// Divider 水平分割线（antd Divider）。上下留白 t.Space(6)（标准密度
// 24px，对齐 antd 默认 margin；紧凑密度自动缩小）：
//
//	ant.Divider(c)
//	ant.Divider(c, ant.WithText("分组"), ant.TextLeft())
func Divider(c *ui.Context, opts ...DividerOption) ui.Element {
	cfg := &dividerConfig{align: ui.Start}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	if cfg.text == "" {
		// 纯线条：底部 1px 描边（可虚线）。
		line := ui.Box(c).FillWidth().Height(t.Space(2)).Border(1, t.Border)
		if cfg.dashed {
			line.BorderStyle(ui.BorderDashed)
		}
		line.BorderWidth(0, 0, 1, 0)
		return line
	}

	// 带文字：线 - 文 - 线。文字必须夹在两条线中间构建（此前先建两条线
	// 再建文字，导致渲染成"线-线-文字"、文字贴边且线中断开）。
	row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
		Padding(t.Space(6), 0)
	text := func() {
		ui.Text(c, cfg.text).TextColor(t.Text).FontWeight(500).Shrink(0)
	}
	row.Children(func() {
		switch cfg.align {
		case ui.Center: // 文字居中：左线 - 文字 - 右线
			left := lineEl(c, t, cfg.dashed)
			left.Grow(1)
			text()
			right := lineEl(c, t, cfg.dashed)
			right.Grow(1)
		case ui.End: // 文字靠右（antd orientation="right"）：长线 - 文字 - 短线
			left := lineEl(c, t, cfg.dashed)
			left.Grow(9)
			text()
			right := lineEl(c, t, cfg.dashed)
			right.Grow(1)
		default: // 文字靠左（antd orientation="left"，默认）：短线 - 文字 - 长线
			left := lineEl(c, t, cfg.dashed)
			left.Grow(1)
			text()
			right := lineEl(c, t, cfg.dashed)
			right.Grow(9)
		}
	})
	return row
}

// lineEl 一条水平线（Grow 由调用方设置）。
func lineEl(c *ui.Context, t *ui.Theme, dashed bool) ui.Element {
	line := ui.Box(c).FillWidth().Height(t.Space(2)).Border(1, t.Border)
	if dashed {
		line.BorderStyle(ui.BorderDashed)
	}
	line.BorderWidth(0, 0, 1, 0)
	return line
}

// VerticalDivider 垂直分割线（antd Divider type="vertical"）。左右留白
// t.Space(2)，高度撑满父容器：
//
//	ant.VerticalDivider(c)
func VerticalDivider(c *ui.Context) ui.Element {
	t := c.Theme()
	line := ui.Box(c).FillHeight().Width(t.Space(2)).Border(1, t.Border).
		Margin(0, t.Space(2))
	line.BorderWidth(0, 1, 0, 0)
	return line
}

// DashedVerticalDivider 虚线垂直分割线。
func DashedVerticalDivider(c *ui.Context) ui.Element {
	t := c.Theme()
	line := ui.Box(c).FillHeight().Width(t.Space(2)).Border(1, t.Border).
		BorderStyle(ui.BorderDashed).Margin(0, t.Space(2))
	line.BorderWidth(0, 1, 0, 0)
	return line
}
