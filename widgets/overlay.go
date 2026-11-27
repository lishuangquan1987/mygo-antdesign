package ant

import (
	"github.com/egoist/mygo/ui"
)

// ---------- Tooltip ----------

// TooltipOption 定制 Tooltip。
type TooltipOption func(*tooltipConfig)

type tooltipConfig struct{}

// Tooltip 文字提示（antd Tooltip）：悬停或键盘聚焦锚点时显示反色小气泡。
// 复用 MyGo 的 TooltipBase（0.6s 延迟、键盘聚焦立即显示、Esc 关闭），
// 外观按 antd 的「反色底 + 反色文字」绘制：
//
//	ant.Tooltip(c, "删除后不可恢复", func(c *ui.Context) {
//	    ant.Button(c, "删除", ant.Danger())
//	})
func Tooltip(c *ui.Context, text string, anchor func(*ui.Context), opts ...TooltipOption) ui.Element {
	cfg := &tooltipConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c)
	var anchorEl ui.Element
	host.Children(func() {
		anchorEl = ui.Row(c)
		anchorEl.Children(func() { anchor(c) })
	})
	ui.TooltipBase(c, anchorEl, func(tip ui.Element) {
		tip.Padding(t.Space(1), t.Space(2)).
			Radius(t.Radius*0.67).
			Background(t.Inverse).
			Border(1, t.Border)
		tip.Children(func() {
			ui.Text(c, text).TextColor(t.InverseText).FontSize(t.Rem(0.86)).SingleLine()
		})
	})
	return host
}

// ---------- Popconfirm ----------

// PopconfirmOption 定制 Popconfirm。
type PopconfirmOption func(*popconfirmConfig)

type popconfirmConfig struct {
	title   string
	okText  string
	canText string
}

// PopconfirmTitle 气泡内的标题（antd title）。
func PopconfirmTitle(s string) PopconfirmOption {
	return func(c *popconfirmConfig) { c.title = s }
}

// PopconfirmOkText 确认按钮文字（默认「确定」）。
func PopconfirmOkText(s string) PopconfirmOption {
	return func(c *popconfirmConfig) { c.okText = s }
}

// PopconfirmCancelText 取消按钮文字（默认「取消」）。
func PopconfirmCancelText(s string) PopconfirmOption {
	return func(c *popconfirmConfig) { c.canText = s }
}

// Popconfirm 气泡确认（antd Popconfirm）：点击触发器弹出确认气泡，
// 确认后执行 onConfirm：
//
//	ant.Popconfirm(c, func(c *ui.Context) ui.Element {
//	    return ant.Button(c, "删除")
//	}, "确定删除这一项？", app.remove, ant.PopconfirmTitle("删除确认"))
func Popconfirm(c *ui.Context, trigger func(*ui.Context) ui.Element,
	message string, onConfirm func(), opts ...PopconfirmOption) ui.Element {
	cfg := &popconfirmConfig{okText: "确定", canText: "取消"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c)
	var anchor, trig ui.Element
	host.Children(func() {
		anchor = ui.Row(c)
		anchor.Children(func() { trig = trigger(c) })
	})
	open := ui.Local(host, "ant.popconfirm.open", func() bool { return false })
	if trig.Clicked() {
		*open = !*open
	}

	ui.Popover(c, anchor, open, func() {
		var cancel, ok bool
		panel := ui.Column(c).Width(t.Space(56)).Gap(t.Space(2)).
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		panel.Children(func() {
			if cfg.title != "" {
				ui.Text(c, cfg.title).TextColor(t.Text).FontSize(t.Rem(0.95)).FontWeight(600)
			}
			ui.Text(c, message).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
			row := ui.Row(c).Gap(t.Space(2)).Justify(ui.End)
			row.Children(func() {
				Button(c, cfg.canText, Size(Small), OnClick(func() { cancel = true }))
				Button(c, cfg.okText, Primary(), Danger(), Size(Small),
					OnClick(func() { ok = true }))
			})
		})
		if cancel {
			*open = false
		}
		if ok {
			*open = false
			if onConfirm != nil {
				onConfirm()
			}
		}
	})
	return host
}
