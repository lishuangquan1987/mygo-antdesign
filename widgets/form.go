package ant

import (
	"github.com/egoist/mygo/ui"
)

// Form 表单（antd Form）：标签列对齐的表单容器，内部用 FormItem 排布。
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

// FormItemRequired 显示必填标记「*」（antd required）。
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

// FormItem 表单项（antd Form.Item）：标签 + 控件（+ 帮助/错误文字）。
func FormItem(c *ui.Context, label string, fn func(), opts ...FormItemOption) ui.Element {
	cfg := &formItemConfig{}
	for _, o := range opts {
		o(cfg)
	}
	text := label
	if cfg.required {
		text += " *"
	}
	field := ui.Field(c, text, fn)
	if cfg.help != "" {
		field.Description(cfg.help)
	}
	if cfg.err != "" {
		field.Error(cfg.err)
	}
	return field
}

// FormGroup 字段分组（antd Form 内的分组 / Fieldset）。
func FormGroup(c *ui.Context, legend string, fn func()) ui.Element {
	return ui.Fieldset(c, legend, fn)
}
