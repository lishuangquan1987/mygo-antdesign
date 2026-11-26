package ant

import (
	"github.com/egoist/mygo/ui"
)

// AlertType 提示条级别（antd Alert 的 type）。
type AlertType string

const (
	AlertInfo    AlertType = "info"
	AlertSuccess AlertType = "success"
	AlertWarning AlertType = "warning"
	AlertError   AlertType = "error"
)

// alertPalette 级别 → 主题语义色。
func alertColor(t *ui.Theme, typ AlertType) ui.Color {
	switch typ {
	case AlertSuccess:
		return t.Success
	case AlertWarning:
		return t.Warning
	case AlertError:
		return t.Danger
	default:
		return t.Accent
	}
}

// alertIcons 级别 → 图标字符。
func alertIcon(typ AlertType) string {
	switch typ {
	case AlertSuccess:
		return "✓"
	case AlertWarning:
		return "!"
	case AlertError:
		return "✕"
	default:
		return "i"
	}
}

// AlertOption 定制 Alert。
type AlertOption func(*alertConfig)

type alertConfig struct {
	typ         AlertType
	description string
	closable    bool
	onClose     func()
}

// AlertTypeOpt 设置级别。
func AlertTypeOpt(typ AlertType) AlertOption { return func(c *alertConfig) { c.typ = typ } }

// AlertDescription 补充说明（第二行小字）。
func AlertDescription(s string) AlertOption { return func(c *alertConfig) { c.description = s } }

// AlertClosable 可关闭，关闭时回调。
func AlertClosable(fn func()) AlertOption {
	return func(c *alertConfig) { c.closable, c.onClose = true, fn }
}

// Alert 提示条（antd Alert）：淡底、彩色图标与标题，可带说明与关闭：
//
//	ant.Alert(c, "保存成功", ant.AlertTypeOpt(ant.AlertSuccess), ant.AlertDescription("配置已写入 server.json"))
func Alert(c *ui.Context, title string, opts ...AlertOption) ui.Element {
	cfg := &alertConfig{typ: AlertInfo}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	color := alertColor(t, cfg.typ)

	// 即时模式关闭状态：关闭后本帧起不再渲染。
	var host ui.Element
	host = ui.Column(c).FillWidth()
	closed := ui.Local(host, "ant.alert.closed", func() bool { return false })

	if *closed {
		host.Children(func() {})
		return host
	}

	host.Children(func() {
		el := ui.Row(c).FillWidth().
			Padding(t.Space(2), t.Space(3)).
			Gap(t.Space(2)).
			Radius(t.Radius).
			Background(color.Mix(t.Background, 0.92)).
			Border(1, color.Mix(t.Background, 0.6)).
			AlignItems(ui.Center)
		el.Children(func() {
			icon := ui.Box(c).Size(t.Space(5), t.Space(5)).Radius(t.Space(2.5)).
				Background(color).Center()
			icon.Children(func() {
				ui.Text(c, alertIcon(cfg.typ)).TextColor(t.AccentText).FontSize(t.Rem(0.8))
			})
			body := ui.Column(c).Grow(1).Gap(t.Space(1))
			body.Children(func() {
				ui.Text(c, title).TextColor(t.Text).FontWeight(600)
				if cfg.description != "" {
					ui.Text(c, cfg.description).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
				}
			})
			if cfg.closable {
				x := ui.ButtonBase(c).Size(t.Space(5), t.Space(5)).Center()
				x.Children(func() {
					ui.Text(c, "×").TextColor(t.TextMuted).FontSize(t.Rem(0.9))
				})
				if x.Clicked() {
					*closed = true
					if cfg.onClose != nil {
						cfg.onClose()
					}
				}
			}
		})
	})
	return host
}

// ModalOption 定制 Modal。
type ModalOption func(*modalConfig)

type modalConfig struct {
	okText     string
	cancelText string
	onOk       func()
	onCancel   func()
	hideCancel bool
}

// OkText 确认按钮文字（默认「确定」）。
func OkText(s string) ModalOption { return func(c *modalConfig) { c.okText = s } }

// CancelText 取消按钮文字（默认「取消」）。
func CancelText(s string) ModalOption { return func(c *modalConfig) { c.cancelText = s } }

// OnOk 点击确认回调（之后自动关闭）。
func OnOk(fn func()) ModalOption { return func(c *modalConfig) { c.onOk = fn } }

// OnCancel 点击取消 / 遮罩 / Esc 回调。
func OnCancel(fn func()) ModalOption { return func(c *modalConfig) { c.onCancel = fn } }

// HideCancel 隐藏取消按钮（确认对话框）。
func HideCancel() ModalOption { return func(c *modalConfig) { c.hideCancel = true } }

// Modal 模态对话框（antd Modal）。open 由调用方持有；点击遮罩或 Esc 由
// MyGo 原语自动关闭（*open 置 false 并触发 OnCancel）：
//
//	ant.Modal(c, &app.confirmOpen, "确认删除", func() {
//	    ant.Paragraph(c, "删除后不可恢复。")
//	}, ant.OnOk(app.doRemove), ant.OkText("删除"))
func Modal(c *ui.Context, open *bool, title string, body func(), opts ...ModalOption) ui.Element {
	cfg := &modalConfig{okText: "确定", cancelText: "取消"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	return ui.Modal(c, open, func() {
		panel := ui.Column(c).
			Width(t.Space(90)).
			MaxWidth(t.Space(120)).
			Radius(t.Radius * 1.34). // antd modal 圆角 8
			Background(t.Surface)
		panel.Children(func() {
			// 头部：标题 + 关闭。
			head := ui.Row(c).FillWidth().Padding(t.Space(4), t.Space(6), t.Space(2), t.Space(6)).
				Justify(ui.SpaceBetween).AlignItems(ui.Center)
			head.Children(func() {
				ui.Text(c, title).FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
				x := ui.ButtonBase(c).Size(t.Space(5), t.Space(5)).Center()
				x.Children(func() {
					ui.Text(c, "×").TextColor(t.TextMuted)
				})
				if x.Clicked() {
					*open = false
					if cfg.onCancel != nil {
						cfg.onCancel()
					}
				}
			})
			// 内容。
			content := ui.Column(c).FillWidth().Padding(0, t.Space(6)).Gap(t.Space(2))
			content.Children(body)
			// 底部操作区。
			foot := ui.Row(c).FillWidth().Padding(t.Space(3), t.Space(6), t.Space(4), t.Space(6)).
				Gap(t.Space(2)).Justify(ui.End)
			foot.Children(func() {
				if !cfg.hideCancel {
					Button(c, cfg.cancelText, OnClick(func() {
						*open = false
						if cfg.onCancel != nil {
							cfg.onCancel()
						}
					}))
				}
				Button(c, cfg.okText, Primary(), OnClick(func() {
					*open = false
					if cfg.onOk != nil {
						cfg.onOk()
					}
				}))
			})
		})
	})
}

// Message 轻提示（antd message）：底部短暂出现的全局消息。
//
//	ant.Message(c, "已保存")
func Message(c *ui.Context, text string) {
	c.Toast(text)
}

// MessageTypeTyped 带类型的轻提示（success / warning / error / info）。
func MessageTypeTyped(c *ui.Context, typ, text string) {
	c.AddToast(ui.Toast{ID: text, Title: text, Type: typ})
}

// Notify 通知（antd notification）：带标题、描述与可选操作按钮。
//
//	ant.Notify(c, "部署完成", "v0.3.1 已上线", "查看", app.openLog)
func Notify(c *ui.Context, title, description, action string, onAction func()) {
	c.AddToast(ui.Toast{
		ID:          title,
		Title:       title,
		Description: description,
		Action:      action,
		OnAction:    onAction,
	})
}

// ProgressOption 定制 Progress。
type ProgressOption func(*progressConfig)

type progressConfig struct {
	showText bool
	height   float32
}

// ProgressShowText 在条尾显示百分比（默认已开启，无需显式传入）。
func ProgressShowText() ProgressOption { return func(c *progressConfig) { c.showText = true } }

// ProgressHideText 隐藏条尾百分比（antd showInfo={false}）。
func ProgressHideText() ProgressOption { return func(c *progressConfig) { c.showText = false } }

// ProgressLine 线形进度条（antd Progress）。value 0–100：
//
//	ant.ProgressLine(c, 68, ant.ProgressShowText()).FillWidth()
func ProgressLine(c *ui.Context, value float64, opts ...ProgressOption) ui.Element {
	cfg := &progressConfig{showText: true}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	row.Children(func() {
		ui.Progress(c, value).Grow(1)
		if cfg.showText {
			pct := int(value + 0.5)
			ui.Textf(c, "%d%%", pct).TextColor(t.TextMuted).FontSize(t.Rem(0.86)).SingleLine()
		}
	})
	return row
}

// Confirm 确认对话框的便捷封装：取消 + 确认两个按钮。
func Confirm(c *ui.Context, open *bool, title, message string, onOk func()) {
	// ui.AlertDialog 返回被点按钮的下标（0 = 取消，1 = 确定）。
	if ui.AlertDialog(c, open, title, message, "取消", "确定") == 1 && onOk != nil {
		onOk()
	}
}
