package ant

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// DatePicker 日期选择框（antd DatePicker）。date 由调用方持有（即时模式），
// 选项复用 InputOption（尺寸/禁用/回调/无障碍标签）：
//
//	var birthday time.Time
//	ant.DatePicker(c, &birthday, ant.OnChanged(app.validate))
func DatePicker(c *ui.Context, date *time.Time, opts ...InputOption) ui.Element {
	return wrapPicker(c, func() ui.Element { return ui.DateInput(c, date) }, opts)
}

// TimePicker 时间选择框（antd TimePicker）。tm 由调用方持有：
//
//	var start time.Time
//	ant.TimePicker(c, &start, ant.InputDisabled(true))
func TimePicker(c *ui.Context, tm *time.Time, opts ...InputOption) ui.Element {
	return wrapPicker(c, func() ui.Element { return ui.TimeInput(c, tm) }, opts)
}

// wrapPicker 统一的日期/时间包装：原语自带边框与弹出交互，这里对齐
// antd 的控件高度、字号与圆角，并接上通用选项。
func wrapPicker(c *ui.Context, build func() ui.Element, opts []InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	el := build().
		Height(internal.ControlHeight(t, cfg.size)).
		FontSize(internal.ControlFont(t, cfg.size)).
		Radius(t.Radius)
	if cfg.label != "" {
		el.Label(cfg.label)
	}
	if cfg.disabled {
		el.Disabled(true)
	}
	if cfg.readonly {
		el.ReadOnly(true)
	}
	fireInputEvents(el, cfg)
	return el
}
