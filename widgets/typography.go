package ant

import (
	"github.com/egoist/mygo/ui"
)

// TitleLevel 标题级别，对齐 antd Typography.Title 的 level 1–5。
type TitleLevel int

const (
	H1 TitleLevel = iota + 1
	H2
	H3
	H4
	H5
)

// titleSize 各级标题字号（相对基准字号的倍数，对齐 antd 38/29/24/20/16px ÷ 14px）。
func titleSize(level TitleLevel) float32 {
	switch level {
	case H2:
		return 2.07
	case H3:
		return 1.71
	case H4:
		return 1.43
	case H5:
		return 1.14
	default:
		return 2.71
	}
}

// Title 标题（antd Typography.Title）。字重 600，颜色为主文本：
//
//	ant.Title(c, "表单配置", ant.H3)
func Title(c *ui.Context, s string, level TitleLevel) ui.Element {
	t := c.Theme()
	return ui.Text(c, s).FontSize(t.Rem(titleSize(level))).FontWeight(600).TextColor(t.Text)
}

// TextOption 定制 Paragraph 的文本样式。
type TextOption func(*paraConfig)

type paraConfig struct {
	secondary bool
	delete    bool
	strong    bool
	underline bool
	italic    bool
	marked    bool
}

// Secondary 次要文本色（antd type="secondary"）。
func Secondary() TextOption { return func(c *paraConfig) { c.secondary = true } }

// Delete 删除线（antd delete）。
func Delete() TextOption { return func(c *paraConfig) { c.delete = true } }

// Strong 加粗（antd strong）。
func Strong() TextOption { return func(c *paraConfig) { c.strong = true } }

// Underline 下划线（antd underline）。
func Underline() TextOption { return func(c *paraConfig) { c.underline = true } }

// Italic 斜体（antd italic）。
func Italic() TextOption { return func(c *paraConfig) { c.italic = true } }

// Marked 高亮背景（antd mark）。
func Marked() TextOption { return func(c *paraConfig) { c.marked = true } }

// Paragraph 正文段落（antd Typography.Paragraph）。默认主文本色：
//
//	ant.Paragraph(c, "说明文字", ant.Secondary())
func Paragraph(c *ui.Context, s string, opts ...TextOption) ui.Element {
	cfg := &paraConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	build := func() ui.Element {
		el := ui.Text(c, s).TextColor(t.Text)
		if cfg.secondary {
			el.TextColor(t.TextMuted)
		}
		if cfg.strong {
			el.FontWeight(600)
		}
		if cfg.underline {
			el.Underline()
		}
		if cfg.italic {
			el.Italic()
		}
		if cfg.delete {
			el.Strikethrough()
		}
		if cfg.marked {
			el.Background(t.Background.Mix(t.Text, 0.12)).Padding(0, t.Space(1))
		}
		return el
	}
	if cfg.marked {
		// mark 高亮只包裹文字（antd <mark> 行内样式）：块级 Text 会撑满
		// 交叉轴，背景随之拉满整行，用 Row 收敛到内容宽度。
		row := ui.Row(c)
		row.Children(func() { build() })
		return row
	}
	return build()
}

// Code 行内代码（antd Typography.Text code）：容器底色调、描边、小字号。
func Code(c *ui.Context, s string) ui.Element {
	t := c.Theme()
	return ui.Text(c, s).
		FontSize(t.Rem(0.86)).
		TextColor(t.Text).
		Background(t.Background.Mix(t.Text, 0.06)).
		Border(1, t.Border).
		Radius(t.Radius*0.6).
		Padding(0, t.Space(1))
}

// TextMuted 次要说明文本的便捷别名（等价 Paragraph + Secondary）。
func TextMuted(c *ui.Context, s string) ui.Element {
	return Paragraph(c, s, Secondary())
}
