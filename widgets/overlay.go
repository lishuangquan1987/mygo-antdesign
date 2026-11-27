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
		// 反色小气泡；箭头在无头 Tester 中无法稳定绘制（见 docs/09），暂不实现。
		tip.Padding(t.Space(1), t.Space(2)).
			Radius(t.Radius * 0.67).
			Background(t.Inverse)
		tip.Children(func() {
			ui.Text(c, text).TextColor(t.InverseText).FontSize(t.FontSize).SingleLine()
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

// Popover 通用气泡卡片（antd Popover）：点击/由调用方控制 open 的锚定浮层，
// 内容与触发器都可由调用方构建：
//
//	var open bool
//	ant.Popover(c, func(c *ui.Context) {
//	    if ant.Button(c, "详情").Clicked() { open = true }
//	}, &open, func() {
//	    ant.Paragraph(c, "气泡内容")
//	})
func Popover(c *ui.Context, anchor func(*ui.Context), open *bool, body func()) ui.Element {
	t := c.Theme()
	host := ui.Row(c)
	var anchorEl ui.Element
	host.Children(func() {
		anchorEl = ui.Row(c)
		anchorEl.Children(func() { anchor(c) })
	})
	ui.Popover(c, anchorEl, open, func() {
		panel := ui.Column(c).Width(t.Space(56)).Gap(t.Space(2)).
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		panel.Children(body)
	})
	return host
}

// ContextMenuItem 右键菜单项。
type ContextMenuItem struct {
	// ID 选择回调的标识（留空用 Label）。
	ID string
	// Label 显示文字。
	Label string
	// Danger 危险项（红色文字，如删除）。
	Danger bool
}

// ContextMenu 右键菜单（antd Dropdown trigger="contextMenu"）：在目标元素上
// 右键或按菜单键弹出，选择后回调：
//
//	ant.ContextMenu(c, func(c *ui.Context) {
//	    ant.Paragraph(c, "右键这一行")
//	}, []ant.ContextMenuItem{
//	    {ID: "rename", Label: "重命名"},
//	    {ID: "delete", Label: "删除", Danger: true},
//	}, app.handleMenu)
func ContextMenu(c *ui.Context, target func(*ui.Context),
	items []ContextMenuItem, onSelect func(id string)) ui.Element {
	t := c.Theme()
	host := ui.Row(c)
	var targetEl ui.Element
	host.Children(func() {
		targetEl = ui.Row(c)
		targetEl.Children(func() { target(c) })
	})
	open := ui.Local(host, "ant.contextmenu.open", func() bool { return false })
	m := ui.ContextMenuBase(c, targetEl, open)
	m.Popup(func(panel ui.Element) {
		panel.Padding(t.Space(1)).Radius(t.Radius).
			Background(t.Surface).Border(1, t.Border)
		for _, it := range items {
			fg := t.Text
			if it.Danger {
				fg = t.Danger
			}
			item := m.Item().Padding(t.Space(1), t.Space(2)).Radius(t.Radius * 0.67)
			item.Children(func() {
				ui.Text(c, it.Label).TextColor(fg).FontSize(t.Rem(0.9)).SingleLine()
			})
			if item.Clicked() {
				*open = false
				if onSelect != nil {
					id := it.ID
					if id == "" {
						id = it.Label
					}
					onSelect(id)
				}
			}
		}
	})
	return host
}
