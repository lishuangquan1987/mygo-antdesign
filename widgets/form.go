package ant

import (
	"github.com/egoist/mygo/ui"
)

// Form 表单（antd Form）：标签在上、控件在下的纵向表单容器，内部用
// FormItem 排布：
//
//	ant.Form(c, func() {
//	    ant.FormItem(c, "用户名", func() { ant.Input(c, &app.name) },
//	        ant.FormItemRequired(true))
//	    ant.FormItem(c, "邮箱", func() { ant.Input(c, &app.email) },
//	        ant.FormItemHelp("用于接收通知"))
//	})
func Form(c *ui.Context, fn func()) ui.Element {
	return ui.Form(c, fn).FillWidth()
}

// FormItemOption 定制 FormItem。
type FormItemOption func(*formItemConfig)

type formItemConfig struct {
	required bool
	help     string
	err      string
}

// FormItemRequired 显示必填标记「*」（antd required，红色）。
func FormItemRequired(b bool) FormItemOption {
	return func(c *formItemConfig) { c.required = b }
}

// FormItemHelp 控件下方的帮助文字（antd extra / help）。
func FormItemHelp(s string) FormItemOption {
	return func(c *formItemConfig) { c.help = s }
}

// FormItemError 控件下方的错误文字（antd validateStatus="error"）。
func FormItemError(s string) FormItemOption {
	return func(c *formItemConfig) { c.err = s }
}

// FormItem 表单项（antd Form.Item，纵向布局）：标签（+ 红色必填星号）+
// 控件（+ 帮助/错误文字），项间距 24px，支持三档尺寸下的字号与间距：
//
//	ant.FormItem(c, "用户名", func() { ant.Input(c, &app.name) })
func FormItem(c *ui.Context, label string, fn func(), opts ...FormItemOption) ui.Element {
	cfg := &formItemConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	item := ui.Column(c).FillWidth().
		Gap(t.Space(1)).
		Margin(0, 0, t.Space(6), 0)
	item.Children(func() {
		// antd 标签高度 32px、右对齐（水平表单）或居左（纵向表单）。
		labelRow := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
		labelRow.Children(func() {
			ui.Text(c, label).TextColor(t.Text).FontSize(t.FontSize).SingleLine()
			if cfg.required {
				// antd 必填星号：colorError + 前置 4px，与标签同基线。
				ui.Text(c, "*").TextColor(t.Danger).FontSize(t.FontSize).SingleLine()
			}
		})
		fn()
		if cfg.help != "" {
			ui.Text(c, cfg.help).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
		}
		if cfg.err != "" {
			ui.Text(c, cfg.err).TextColor(t.Danger).FontSize(t.Rem(0.86))
		}
	})
	return item
}

// FormGroup 字段分组（antd Form 内的分组 / Fieldset）。
func FormGroup(c *ui.Context, legend string, fn func()) ui.Element {
	return ui.Fieldset(c, legend, fn)
}
