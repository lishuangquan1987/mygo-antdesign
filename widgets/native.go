package ant

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Hyperlink 文字链接（antd Typography.Link）：主色文字，点击打开 URL。
//
//	ant.Link(c, "查看文档", "https://example.com")
func Hyperlink(c *ui.Context, label, url string) ui.Element {
	t := c.Theme()
	return ui.Link(c, label, url).TextColor(t.Accent)
}

func applyControlOpts(el ui.Element, opts []ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	applyControl(el, cfg)
	return el
}

// DecimalInput 精确金额输入（antd InputNumber 的金额场景）：值以字符串
// 受控，绝不经过浮点舍入：
//
//	ant.DecimalInput(c, &app.amount, "0", "1000000", "0.01")
func DecimalInput(c *ui.Context, value *string, lo, hi, step string, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.DecimalInput(c, value, lo, hi, step), opts)
}

// SearchField 原生搜索框（放大镜 + 清除 + Esc 清空）；与 ant.SearchInput
// 的区别是直接用 MyGo 的 SearchField 原语：
//
//	ant.SearchField(c, &app.query)
func SearchField(c *ui.Context, query *string, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.SearchField(c, query), opts)
}

// Stepper 数字步进器（antd InputNumber 的上下箭头形态；注意它不是 Steps）：
//
//	ant.Stepper(c, &app.count, 0, 100, 1)
func Stepper(c *ui.Context, value *float64, lo, hi, step float64, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.Stepper(c, value, lo, hi, step), opts)
}

// Calendar 月历（antd Calendar）：浏览 + 选择日期：
//
//	ant.Calendar(c, &app.date)
func Calendar(c *ui.Context, date *time.Time, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.Calendar(c, date), opts)
}

// TokenField 标签输入（antd Select tags 模式）：值以 *[]string 受控，
// 支持候选建议：
//
//	ant.TokenField(c, &app.tags, []string{"Go", "GUI", "Ant Design"})
func TokenField(c *ui.Context, tokens *[]string, suggestions []string, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.TokenField(c, tokens, suggestions), opts)
}

// EditableText 原位重命名文本（antd Typography editable）：
//
//	ant.EditableText(c, &app.name)
func EditableText(c *ui.Context, value *string, opts ...ControlOption) ui.Element {
	return applyControlOpts(ui.EditableText(c, value), opts)
}

// Meter 电平/阈值条（antd Progress 的仪表近似）：按阈值显示成功/警告/危险色：
//
//	ant.Meter(c, 72, 0, 100, &ui.MeterLevels{Warning: 60, Critical: 85})
func Meter(c *ui.Context, value, lo, hi float64, levels *ui.MeterLevels) ui.Element {
	return ui.Meter(c, value, lo, hi, levels)
}

// Collapsible 单个折叠块（antd Collapse 单项；带展开/收起动画）：
//
//	ant.Collapsible(c, "高级设置", &app.advanced, func() {
//	    ant.Paragraph(c, "内容")
//	})
func Collapsible(c *ui.Context, label string, open *bool, fn func()) ui.Element {
	return ui.Collapsible(c, label, open, fn)
}

// Accordion 手风琴（antd Collapse 但每项独立展开）：
//
//	ant.Accordion(c, func() {
//	    ant.AccordionItem(c, "常规", &app.general, func() { ... })
//	})
func Accordion(c *ui.Context, fn func()) ui.Element {
	return ui.Accordion(c, fn)
}

// AccordionItem 手风琴的一项。
func AccordionItem(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	return ui.AccordionItem(c, title, open, fn)
}
