package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
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
	t := c.Theme()

	// antd 的复选框与文字间距 8px（原生控件内部间距不可调，这里自绘一行）。
	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	var box ui.Element
	row.Children(func() {
		box = ui.Checkbox(c, checked, "")
		if cfg.disabled {
			box.Disabled(true)
		}
		lbl := ui.Text(c, label).TextColor(t.Text).FontSize(t.FontSize)
		if !cfg.disabled && lbl.Clicked() {
			*checked = !*checked
			if cfg.onChg != nil {
				cfg.onChg()
			}
		}
	})
	if cfg.onChg != nil && box.Changed() {
		cfg.onChg()
	}
	return row
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
//
// 尺寸对齐 antd：44×22 轨道、18px 圆形滑块、2px 内边距；选中为主色。
func Switch(c *ui.Context, on *bool, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return switchControl(c, on, cfg)
}

// switchControl 自绘 antd Switch 轨道与滑块（原生 Switch 尺寸固定 40×22）。
func switchControl(c *ui.Context, on *bool, cfg *controlConfig) ui.Element {
	t := c.Theme()
	track := ui.ButtonBase(c).
		Size(t.Space(11), t.Space(5.5)).
		Radius(t.Space(2.75)).
		AlignItems(ui.Center)
	if *on {
		track.Background(t.Accent)
	} else {
		track.Background(t.Border)
	}
	if cfg.disabled {
		track.Disabled(true)
	}
	track.Children(func() {
		knob := ui.Box(c).
			Size(t.Space(4.5), t.Space(4.5)).
			Radius(t.Space(2.25)).
			Background(t.Surface)
		if *on {
			knob.MarginStart(t.Space(6))
		} else {
			knob.MarginStart(t.Space(0.5))
		}
		knob.Children(func() {})
	})
	if !cfg.disabled && track.Clicked() {
		*on = !*on
		if cfg.onChg != nil {
			cfg.onChg()
		}
	}
	return track
}

// SwitchWithLabel 带旁边文字的开关（组合 Switch + 标签文本）。
func SwitchWithLabel(c *ui.Context, on *bool, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	row.Children(func() {
		switchControl(c, on, cfg)
		lbl := ui.Text(c, label).TextColor(t.Text).FontSize(t.FontSize)
		if !cfg.disabled && lbl.Clicked() {
			*on = !*on
			if cfg.onChg != nil {
				cfg.onChg()
			}
		}
	})
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

// Select 下拉选择（antd Select）。selected 与 options 由调用方持有；自绘
// antd 规格的触发器（32px 高、11px 水平内边距、6px 圆角、右侧箭头）与
// 下拉面板（32px 选项、选中/悬停浅底），键盘导航由 SelectBase 提供：
//
//	var city string
//	ant.Select(c, &city, []string{"北京", "上海", "广州"}, ant.OnChange(app.reload))
func Select(c *ui.Context, selected *string, options []string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	before := *selected
	sel := ui.SelectBase(c, selected)

	// 触发器：与 Input 同高的 antd Select 外框。
	sel.Trigger.
		Height(internal.ControlHeight(t, internal.SizeMiddle)).
		Padding(0, t.Space(3)).
		Radius(t.Radius).
		Border(1, t.Border).
		Background(t.Surface).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(t.Space(1))
	if cfg.disabled {
		sel.Trigger.Disabled(true)
	}
	sel.Trigger.Children(func() {
		label := *selected
		color := t.Text
		if label == "" {
			label, color = "请选择", t.TextMuted
		}
		ui.Text(c, label).TextColor(color).
			FontSize(internal.ControlFont(t, internal.SizeMiddle)).SingleLine()
		icon.Icon(c, icon.Down, t.Space(3), icon.WithColor(t.TextMuted))
	})

	// 下拉面板：8px 内边距、6px 圆角、32px 选项高。
	sel.Popup(func(panel ui.Element) {
		panel.Padding(t.Space(1)).
			Radius(t.Radius+t.Space(0.5)).
			Background(t.Surface).
			Border(1, t.Border).
			Shadow(0, 6, 16, 0, t.Text.Alpha(0.08))
		for _, opt := range options {
			opt := opt
			item := sel.Item(opt).
				Height(internal.ControlHeight(t, internal.SizeSmall)).
				Padding(0, t.Space(2)).
				Radius(t.Radius * 0.67).
				AlignItems(ui.Center)
			if opt == *selected {
				item.Background(t.Selection)
			} else if item.Highlighted() {
				item.Background(t.SurfaceHover)
			}
			item.Children(func() {
				ui.Text(c, opt).TextColor(t.Text).
					FontSize(internal.ControlFont(t, internal.SizeMiddle)).SingleLine()
			})
		}
	})

	if *selected != before && cfg.onChg != nil {
		cfg.onChg()
	}
	return sel.Trigger
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
