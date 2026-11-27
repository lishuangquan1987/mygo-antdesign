package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// ControlOption 选择类控件（Checkbox/Radio/Switch/Slider/Rate/Segmented）
// 的通用选项。
type ControlOption func(*controlConfig)

type controlConfig struct {
	disabled bool
	onChg    func()
}

// ControlDisabled 禁用（antd disabled）。
func ControlDisabled(b bool) ControlOption { return func(c *controlConfig) { c.disabled = b } }

// OnChange 值变化回调（antd onChange；Changed 为本帧变化语义）。
func OnChange(fn func()) ControlOption { return func(c *controlConfig) { c.onChg = fn } }

// applyControl 应用通用选项并在值变化时触发回调。
func applyControl(el ui.Element, cfg *controlConfig) {
	if cfg.disabled {
		el.Disabled(true)
	}
	if cfg.onChg != nil && el.Changed() {
		cfg.onChg()
	}
}

// Checkbox 复选框（antd Checkbox）。选中色为品牌色（来自主题 Accent）：
//
//	var agree bool
//	ant.Checkbox(c, &agree, "同意协议", ant.OnChange(app.validate))
func Checkbox(c *ui.Context, checked *bool, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Checkbox(c, checked, label)
	applyControl(el, cfg)
	return el
}

// RadioGroupItem 供 RadioGroup 内使用的单项。泛型对齐 antd 的受控值。
//
//	ant.RadioGroup(c, func() {
//	    ant.Radio(c, &size, "small", "小")
//	    ant.Radio(c, &size, "large", "大")
//	})
func RadioGroup(c *ui.Context, fn func()) ui.Element {
	return ui.RadioGroup(c, fn)
}

// Radio 单选项（antd Radio）。
func Radio[T comparable](c *ui.Context, selected *T, value T, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Radio(c, selected, value, label)
	applyControl(el, cfg)
	return el
}

// Switch 开关（antd Switch）。开为品牌色：
//
//	ant.Switch(c, &app.enabled, ant.OnChange(app.toggle))
func Switch(c *ui.Context, on *bool, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Switch(c, on)
	applyControl(el, cfg)
	return el
}

// SwitchWithLabel 带旁边文字的开关（组合 Switch + 标签文本）。
func SwitchWithLabel(c *ui.Context, on *bool, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	var sw ui.Element
	row.Children(func() {
		sw = ui.Switch(c, on)
		if cfg.disabled {
			sw.Disabled(true)
		}
		ui.Text(c, label).TextColor(t.Text)
	})
	if cfg.onChg != nil && sw.Changed() {
		cfg.onChg()
	}
	return row
}

// Slider 滑动条（antd Slider）。宽度默认随父容器，可用返回值链式
// .FillWidth() / .Width() 调整：
//
//	ant.Slider(c, &volume, 0, 100, ant.OnChange(app.setVolume)).FillWidth()
func Slider(c *ui.Context, value *float64, lo, hi float64, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Slider(c, value, lo, hi)
	applyControl(el, cfg)
	return el
}

// RangeSliderDouble 范围滑动条（antd Slider range）。
func RangeSliderDouble(c *ui.Context, low, high *float64, lo, hi, step float64, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.RangeSlider(c, low, high, lo, hi, step)
	applyControl(el, cfg)
	return el
}

// Rate 评分（antd Rate）：
//
//	ant.Rate(c, &score, 5, ant.OnChange(app.submit))
func Rate(c *ui.Context, value *int, max int, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Rating(c, value, max)
	applyControl(el, cfg)
	return el
}

// Segmented 分段控制器（antd Segmented）。selected 为选中下标：
//
//	ant.Segmented(c, &tab, []string{"全部", "进行中", "已完成"}, ant.OnChange(app.filter))
func Segmented(c *ui.Context, selected *int, labels []string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Segmented(c, selected, labels...)
	applyControl(el, cfg)
	return el
}

// Select 下拉选择（antd Select）。selected 与 options 由调用方持有；
// 原语自带下拉弹层与键盘导航，这里对齐 antd 的控件高度、字号与圆角：
//
//	var city string
//	ant.Select(c, &city, []string{"北京", "上海", "广州"}, ant.OnChange(app.reload))
func Select(c *ui.Context, selected *string, options []string, opts ...ControlOption) ui.Element {
	return wrapSelect(c, func() ui.Element { return ui.Select(c, selected, options) }, opts)
}

// Combobox 可输入下拉（antd AutoComplete / Combobox）：输入过滤 + 下拉候选。
//
//	var q string
//	ant.Combobox(c, &q, []string{"重启服务", "重启网关"}, ant.OnChange(app.suggest))
func Combobox(c *ui.Context, selected *string, options []string, opts ...ControlOption) ui.Element {
	return wrapSelect(c, func() ui.Element { return ui.Combobox(c, selected, options) }, opts)
}

// wrapSelect 统一的下拉类包装（Select/Combobox 共用）。
func wrapSelect(c *ui.Context, build func() ui.Element, opts []ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	el := build().
		Height(internal.ControlHeight(t, internal.SizeMiddle)).
		FontSize(internal.ControlFont(t, internal.SizeMiddle)).
		Radius(t.Radius)
	applyControl(el, cfg)
	return el
}

// CheckboxAll 三态复选框（全选框：全选/半选/未选），antd 表格表头场景：
//
//	var all ui.CheckState
//	ant.CheckboxAll(c, &all, "全选")
func CheckboxAll(c *ui.Context, state *ui.CheckState, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.CheckboxState(c, state, label)
	applyControl(el, cfg)
	return el
}

// AutoComplete 自动补全（antd AutoComplete）：保留用户键入内容，同时给出
// 候选建议；与 ant.Combobox 的差别是值不限于候选之一：
//
//	ant.AutoComplete(c, &app.city, []string{"北京", "上海", "广州"})
func AutoComplete(c *ui.Context, value *string, suggestions []string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.Autocomplete(c, value, suggestions)
	applyControl(el, cfg)
	return el
}

// Toggle 保持按下的按钮（antd 的 Toggle/Radio.Button 场景，例如加粗、斜体）：
//
//	ant.Toggle(c, &app.bold, "B", ant.OnChange(app.save))
func Toggle(c *ui.Context, on *bool, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	host := ui.Row(c).AlignItems(ui.Center)
	host.Children(func() {
		el := ui.Toggle(c, on, label)
		applyControl(el, cfg)
	})
	return host
}

// ToggleGroup 连体开关组（antd Radio.Button 组）：内部用 ui.ToggleGroup 组合：
//
//	ant.ToggleGroup(c, func() {
//	    ant.Toggle(c, &app.align, "左")
//	    ant.Toggle(c, &app.align, "中")
//	})
func ToggleGroup(c *ui.Context, fn func()) ui.Element {
	return ui.ToggleGroup(c, fn)
}
