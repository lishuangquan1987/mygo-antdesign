package ant

import (
	"regexp"
	"strings"
	"unicode/utf8"

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
	value    func() string
	rules    []FormRule
}

// FormRule 一条校验规则（antd Form rules 的轻量模型）：各字段可组合，
// 按声明顺序取第一条不通过的规则作为错误文案。即时模式下每帧校验。
type FormRule struct {
	// Required 必填（去掉首尾空白后非空）。
	Required bool
	// MinLen / MaxLen 字符数（按 rune 计，0 表示不限制）。
	MinLen, MaxLen int
	// Pattern 正则匹配；PatternHint 为不匹配时的文案。
	Pattern     *regexp.Regexp
	PatternHint string
	// Check 自定义判定；返回 false 即失败。
	Check func(string) bool
	// Message 统一错误文案（留空时按规则生成默认文案）。
	Message string
}

// validate 返回错误文案；通过时为空串。
func (r FormRule) validate(v string) string {
	msg := r.Message
	s := strings.TrimSpace(v)
	if r.Required && s == "" {
		if msg == "" {
			msg = "此项为必填项"
		}
		return msg
	}
	if s == "" {
		return ""
	}
	if r.MinLen > 0 && utf8.RuneCountInString(s) < r.MinLen {
		if msg == "" {
			msg = "至少 " + itoa(r.MinLen) + " 个字符"
		}
		return msg
	}
	if r.MaxLen > 0 && utf8.RuneCountInString(s) > r.MaxLen {
		if msg == "" {
			msg = "最多 " + itoa(r.MaxLen) + " 个字符"
		}
		return msg
	}
	if r.Pattern != nil && !r.Pattern.MatchString(s) {
		if msg == "" {
			msg = r.PatternHint
		}
		if msg == "" {
			msg = "格式不正确"
		}
		return msg
	}
	if r.Check != nil && !r.Check(s) {
		if msg == "" {
			msg = "校验未通过"
		}
		return msg
	}
	return ""
}

// itoa 小整数转字符串（避免为一条提示引入 strconv）。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
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

// FormItemValue 提供该字段当前值，用于按 FormItemRules 即时校验：
//
//	ant.FormItem(c, "邮箱", func() { ant.Input(c, &app.email) },
//	    ant.FormItemRequired(true),
//	    ant.FormItemValue(func() string { return app.email }),
//	    ant.FormItemRules(ant.FormRule{Pattern: emailRe, PatternHint: "邮箱格式不正确"}))
func FormItemValue(get func() string) FormItemOption {
	return func(c *formItemConfig) { c.value = get }
}

// FormItemRules 逐帧校验规则（按顺序取第一条失败文案，替代 help 显示）。
func FormItemRules(rules ...FormRule) FormItemOption {
	return func(c *formItemConfig) { c.rules = append(c.rules, rules...) }
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

	// 即时校验：有取值函数时每帧计算，错误优先于帮助文字。FormItemRequired
	// 的星号同时也是「必填」校验（对齐 antd：required 是 rules 的一项，
	// 只显示星号不校验会让「标了必填却放行空值」）。
	errText := cfg.err
	if cfg.value != nil {
		v := cfg.value()
		switch {
		case cfg.required && strings.TrimSpace(v) == "":
			errText = "此项为必填项"
		case len(cfg.rules) > 0:
			for _, r := range cfg.rules {
				if msg := r.validate(v); msg != "" {
					errText = msg
					break
				}
			}
		}
	}

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
		if cfg.help != "" && errText == "" {
			ui.Text(c, cfg.help).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
		}
		if errText != "" {
			ui.Text(c, errText).TextColor(t.Danger).FontSize(t.Rem(0.86))
		}
	})
	return item
}

// FormGroup 字段分组（antd Form 内的分组 / Fieldset）。
func FormGroup(c *ui.Context, legend string, fn func()) ui.Element {
	return ui.Fieldset(c, legend, fn)
}
