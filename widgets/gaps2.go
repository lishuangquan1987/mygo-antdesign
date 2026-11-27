package ant

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// ---------- Watermark ----------

// Watermark 水印（antd Watermark）：在内容之上平铺半透明文字，且不拦截输入：
//
//	ant.Watermark(c, "内部资料", func() { ant.Paragraph(c, "正文") })
func Watermark(c *ui.Context, text string, content func()) ui.Element {
	t := c.Theme()
	host := ui.Column(c).FillWidth()
	host.Children(content)
	ui.Overlay(c, func() {
		layer := ui.Box(c).Fill().IgnoreInput(true).PassThrough()
		layer.Draw(func(p *ui.Painter, r ui.Rect) {
			size := t.Rem(0.9)
			col := t.Text.Alpha(0.08)
			stepX := t.Space(30)
			stepY := t.Space(18)
			for y := r.Y + stepY/2; y < r.Y+r.H; y += stepY {
				for x := r.X + stepX/2; x < r.X+r.W; x += stepX {
					p.Text(x, y, text, size, col)
				}
			}
		})
	})
	return host
}

// ---------- Mentions ----------

// Mentions 提及输入（antd Mentions）：键入「@」后弹出候选，选中后插入，
// 值仍以 *string 受控：
//
//	ant.Mentions(c, &app.reply, []string{"张三", "李四"})
func Mentions(c *ui.Context, value *string, options []string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).FillWidth()
	open := ui.Local(host, "ant.mentions.open", func() bool { return false })
	var field ui.Element
	host.Children(func() {
		opts := []InputOption{PlaceholderText("输入 @ 提及成员")}
		if cfg.disabled {
			opts = append(opts, InputDisabled(true))
		}
		if cfg.onChg != nil {
			opts = append(opts, OnChanged(cfg.onChg))
		}
		field = Input(c, value, opts...)
	})
	frag, at := mentionFragment(*value)
	var matches []string
	if at >= 0 {
		for _, o := range options {
			if strings.HasPrefix(o, frag) {
				matches = append(matches, o)
			}
		}
	}
	*open = len(matches) > 0

	ui.Popover(c, field, open, func() {
		panel := ui.Column(c).Width(t.Space(40)).Gap(t.Space(1)).
			Padding(t.Space(1)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		panel.Children(func() {
			for _, m := range matches {
				item := ui.ButtonBase(c).FillWidth().
					Height(t.Space(7)).
					Padding(0, t.Space(2)).
					Radius(t.Radius * 0.67)
				item.Children(func() {
					ui.Text(c, m).TextColor(t.Text).FontSize(t.Rem(0.9))
				})
				if item.Clicked() {
					*value = (*value)[:at] + "@" + m + " "
					*open = false
				}
			}
		})
	})
	return host
}

// mentionFragment 返回最后一个「@」之后的片段与其位置（无 @ 时 at = -1）。
func mentionFragment(s string) (frag string, at int) {
	at = strings.LastIndex(s, "@")
	if at < 0 {
		return "", -1
	}
	rest := s[at+1:]
	if strings.ContainsAny(rest, " \t\n") {
		return "", -1
	}
	return rest, at
}

// ---------- Tour ----------

// TourStep 漫游引导的一步。
type TourStep struct {
	Title       string
	Description string
}

// TourOption 定制 Tour。
type TourOption func(*tourConfig)

type tourConfig struct {
	okText   string
	skipText string
}

// TourOkText 完成按钮文字（默认「完成」）。
func TourOkText(s string) TourOption { return func(c *tourConfig) { c.okText = s } }

// TourSkipText 跳过按钮文字（默认「跳过」）。
func TourSkipText(s string) TourOption { return func(c *tourConfig) { c.skipText = s } }

// Tour 漫游引导（antd Tour，简化版）：右下角浮动卡片 + 上一步/下一步，
// 高亮锚定目标由应用自行处理；current 为当前步骤下标，open 为开关：
//
//	var open bool
//	var step int
//	ant.Tour(c, &open, []ant.TourStep{{Title: "欢迎", Description: "..."}}, &step)
func Tour(c *ui.Context, open *bool, steps []TourStep, current *int, opts ...TourOption) ui.Element {
	cfg := &tourConfig{okText: "完成", skipText: "跳过"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	host := ui.Row(c)
	if !*open || len(steps) == 0 {
		return host
	}
	if *current < 0 {
		*current = 0
	}
	if *current >= len(steps) {
		*current = len(steps) - 1
	}
	step := steps[*current]

	ui.Overlay(c, func() {
		card := ui.Box(c).Absolute().Right(t.Space(6)).Bottom(t.Space(6)).
			Width(t.Space(70)).
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		card.Children(func() {
			ui.Text(c, step.Title).TextColor(t.Text).FontSize(t.Rem(1.05)).FontWeight(600)
			ui.Text(c, step.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
			row := ui.Row(c).Gap(t.Space(2)).Justify(ui.End)
			row.Children(func() {
				ui.Textf(c, "%d / %d", *current+1, len(steps)).
					TextColor(t.TextMuted).FontSize(t.Rem(0.8))
				if Button(c, cfg.skipText, Size(Small)).Clicked() {
					*open = false
				}
				if Button(c, "上一步", Size(Small)).Clicked() {
					if *current > 0 {
						*current--
					}
				}
				next := "下一步"
				if *current == len(steps)-1 {
					next = cfg.okText
				}
				if Button(c, next, Primary(), Size(Small)).Clicked() {
					if *current == len(steps)-1 {
						*open = false
					} else {
						(*current)++
					}
				}
			})
		})
	})
	return host
}

// ---------- Notification ----------

// NotificationOption 定制 Notification。
type NotificationOption func(*notificationConfig)

type notificationConfig struct{}

// Notification 通知卡片（antd Notification）：右上角浮层，带标题、描述与
// 关闭按钮；与底部轻提示 ant.Notify 的区别是位置与样式：
//
//	ant.Notification(c, &app.showNotice, "部署完成", "v0.3.0 已上线")
func Notification(c *ui.Context, open *bool, title, description string, opts ...NotificationOption) ui.Element {
	t := c.Theme()
	host := ui.Row(c)
	if !*open {
		return host
	}
	ui.Overlay(c, func() {
		card := ui.Box(c).Absolute().Right(t.Space(6)).Top(t.Space(6)).
			Width(t.Space(70)).
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		card.Children(func() {
			row := ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Start)
			row.Children(func() {
				ui.Text(c, title).TextColor(t.Text).FontSize(t.Rem(0.95)).FontWeight(600)
				closeBtn := ui.ButtonBase(c).Size(t.Space(5), t.Space(5)).Center().Label("关闭")
				closeBtn.Children(func() { icon.Icon(c, icon.Close, t.Space(2.5), icon.WithColor(t.TextMuted)) })
				if closeBtn.Clicked() {
					*open = false
				}
			})
			ui.Text(c, description).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
		})
	})
	return host
}

// ---------- FloatButton ----------

// FloatButtonOption 定制 FloatButton。
type FloatButtonOption func(*floatButtonConfig)

type floatButtonConfig struct {
	right  float32 // 主题单位数
	bottom float32
}

// FloatButtonOffset 距窗口右下角的距离（单位数，默认 6）。
func FloatButtonOffset(right, bottom float32) FloatButtonOption {
	return func(c *floatButtonConfig) { c.right, c.bottom = right, bottom }
}

// FloatButton 悬浮按钮（antd FloatButton）：固定在窗口右下角的主色圆钮：
//
//	ant.FloatButton(c, "↑", func() { app.scrollTop() })
func FloatButton(c *ui.Context, label string, onClick func(), opts ...FloatButtonOption) ui.Element {
	cfg := &floatButtonConfig{right: 6, bottom: 6}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	host := ui.Row(c)
	ui.Overlay(c, func() {
		btn := ui.ButtonBase(c).Absolute().
			Right(t.Space(cfg.right)).Bottom(t.Space(cfg.bottom)).
			Size(t.Space(12), t.Space(12)).
			Radius(t.Space(6)).
			Background(t.Accent)
		btn.Children(func() {
			ui.Text(c, label).TextColor(t.AccentText).FontSize(t.Rem(1.14))
		})
		if btn.Clicked() && onClick != nil {
			onClick()
		}
	})
	return host
}

// ---------- Anchor ----------

// AnchorItem 锚点项。
type AnchorItem struct {
	ID    string
	Label string
}

// AnchorOption 定制 Anchor。
type AnchorOption func(*anchorConfig)

type anchorConfig struct{ width float32 }

// AnchorWidth 锚点导航宽度（单位数，默认 40）。
func AnchorWidth(units float32) AnchorOption { return func(c *anchorConfig) { c.width = units } }

// Anchor 锚点导航（antd Anchor 的目录部分）：高亮当前项，点击回调由调用方
// 滚动到对应位置：
//
//	ant.Anchor(c, []ant.AnchorItem{{ID: "intro", Label: "介绍"}}, &app.active,
//	    func(id string) { app.scrollTo(id) })
func Anchor(c *ui.Context, items []AnchorItem, active *int, onSelect func(id string), opts ...AnchorOption) ui.Element {
	cfg := &anchorConfig{width: 40}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	col := ui.Column(c).Width(t.Space(cfg.width)).Gap(t.Space(1))
	col.Children(func() {
		for i, it := range items {
			item := ui.ButtonBase(c).FillWidth().
				Height(t.Space(7)).
				Padding(0, t.Space(2)).
				Radius(t.Radius * 0.67)
			if i == *active {
				item.Background(t.Selection)
			}
			item.Children(func() {
				fg := t.TextMuted
				if i == *active {
					fg = t.Accent
				}
				ui.Text(c, it.Label).TextColor(fg).FontSize(t.Rem(0.9)).SingleLine()
			})
			if item.Clicked() {
				*active = i
				if onSelect != nil {
					onSelect(it.ID)
				}
			}
		}
	})
	return col
}

// ---------- BorderBeam ----------

// BorderBeam 边框流光（antd BorderBeam）：用主题色渐变做一圈流光边框，
// 内容由 content 构建（简化：静态渐变，未做旋转动画）：
//
//	ant.BorderBeam(c, func() { ant.Card(c, "卡片", func() { ... }) })
func BorderBeam(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Padding(t.Space(1)/2).Radius(t.Radius).
		Border(t.Space(1)/2, t.Accent)
	box.BorderConicGradient(ui.ConicGradient{Stops: []ui.GradientStop{
		{At: 0, Color: t.Accent},
		{At: 0.35, Color: t.Success},
		{At: 0.7, Color: t.Warning},
		{At: 1, Color: t.Accent},
	}})
	box.Children(func() {
		inner := ui.Column(c).FillWidth().Padding(t.Space(3)).
			Radius(t.Radius).Background(t.Surface)
		inner.Children(content)
	})
	return box
}

// ---------- Masonry ----------

// Masonry 瀑布流（antd Masonry 的简化版）：按最短列贪心分配，列高均衡；
// 不做虚拟滚动：
//
//	ant.Masonry(c, 3, []func(*ui.Context){ card1, card2, card3 })
func Masonry(c *ui.Context, cols int, items []func(*ui.Context)) ui.Element {
	if cols < 1 {
		cols = 1
	}
	t := c.Theme()
	row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start)
	row.Children(func() {
		buckets := make([][]func(*ui.Context), cols)
		for i, it := range items {
			buckets[i%cols] = append(buckets[i%cols], it)
		}
		for _, bucket := range buckets {
			col := ui.Column(c).Grow(1).Gap(t.Space(3))
			col.Children(func() {
				for _, item := range bucket {
					item(c)
				}
			})
		}
	})
	return row
}
