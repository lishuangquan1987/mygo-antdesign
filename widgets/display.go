package ant

import (
	"strconv"

	"github.com/egoist/mygo/ui"
)

// TagColor 标签预设色系（antd Tag 的 color 属性常用子集）。
type TagColor string

const (
	TagColorDefault TagColor = "default"
	TagColorSuccess TagColor = "success"
	TagColorWarning TagColor = "warning"
	TagColorError   TagColor = "error"
	TagColorPrimary TagColor = "primary"
	TagColorMagenta TagColor = "magenta"
	TagColorCyan    TagColor = "cyan"
	TagColorGold    TagColor = "gold"
)

// tagColors 语义色到主题色的映射：antd 的彩色标签是「淡底 + 深字 +
// 淡描边」，全部从主题语义色派生（暗色主题自动成立）。
type tagPalette struct {
	Bg, Fg, Border ui.Color
}

func tagPal(t *ui.Theme, col TagColor) tagPalette {
	switch col {
	case TagColorSuccess:
		return tagPalette{t.Success.Mix(t.Background, 0.85), t.Success, t.Success.Mix(t.Background, 0.6)}
	case TagColorWarning:
		return tagPalette{t.Warning.Mix(t.Background, 0.85), t.Warning, t.Warning.Mix(t.Background, 0.6)}
	case TagColorError:
		return tagPalette{t.Danger.Mix(t.Background, 0.85), t.Danger, t.Danger.Mix(t.Background, 0.6)}
	case TagColorPrimary:
		return tagPalette{t.Accent.Mix(t.Background, 0.85), t.Accent, t.Accent.Mix(t.Background, 0.6)}
	case TagColorMagenta:
		// antd magenta-6 / -1；从 Danger 与 Warning 混出洋红调。
		m := t.Danger.Mix(t.Warning, 0.4)
		return tagPalette{m.Mix(t.Background, 0.85), m, m.Mix(t.Background, 0.6)}
	case TagColorCyan:
		// antd cyan-6 / -1；从 Success 与 Accent 混出青色调。
		m := t.Success.Mix(t.Accent, 0.5)
		return tagPalette{m.Mix(t.Background, 0.85), m, m.Mix(t.Background, 0.6)}
	case TagColorGold:
		m := t.Warning
		return tagPalette{m.Mix(t.Background, 0.88), m.Mix(t.Text, 0.15), m.Mix(t.Background, 0.65)}
	default:
		return tagPalette{t.Background.Mix(t.Text, 0.05), t.TextMuted, t.Border}
	}
}

// Tag 标签（antd Tag）。closable 时在标签尾部显示关闭按钮并触发 OnClose：
//
//	ant.Tag(c, "进行中", ant.TagSuccess())
//	ant.Tag(c, "可关闭", ant.TagClosable(func() { app.remove(id) }))
func Tag(c *ui.Context, text string, opts ...TagOption) ui.Element {
	cfg := &tagConfig{color: TagColorDefault}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	pal := tagPal(t, cfg.color)

	el := ui.Box(c).Height(t.Space(5)).
		Padding(0, t.Space(2)).
		Radius(t.Radius*0.67). // antd tag 圆角略小于全局
		Background(pal.Bg).
		Border(1, pal.Border).
		AlignItems(ui.Center)
	el.Children(func() {
		ui.Text(c, text).TextColor(pal.Fg).FontSize(t.Rem(0.86)).SingleLine()
		if cfg.closable {
			x := ui.ButtonBase(c).Size(t.Space(4), t.Space(4)).Center()
			x.Children(func() {
				ui.Text(c, "✕").TextColor(pal.Fg).FontSize(t.Rem(0.7))
			})
			if x.Clicked() && cfg.onClose != nil {
				cfg.onClose()
			}
		}
	})
	return el
}

// TagOption 定制 Tag。
type TagOption func(*tagConfig)

type tagConfig struct {
	color    TagColor
	closable bool
	onClose  func()
}

// TagSuccess 成功色标签。
func TagSuccess() TagOption { return func(c *tagConfig) { c.color = TagColorSuccess } }

// TagWarning 警告色标签。
func TagWarning() TagOption { return func(c *tagConfig) { c.color = TagColorWarning } }

// TagError 错误色标签。
func TagError() TagOption { return func(c *tagConfig) { c.color = TagColorError } }

// TagPrimary 品牌色标签。
func TagPrimary() TagOption { return func(c *tagConfig) { c.color = TagColorPrimary } }

// TagMagenta 洋红标签。
func TagMagenta() TagOption { return func(c *tagConfig) { c.color = TagColorMagenta } }

// TagCyan 青色标签。
func TagCyan() TagOption { return func(c *tagConfig) { c.color = TagColorCyan } }

// TagGold 金色标签。
func TagGold() TagOption { return func(c *tagConfig) { c.color = TagColorGold } }

// TagClosable 可关闭标签，关闭时回调。
func TagClosable(fn func()) TagOption { return func(c *tagConfig) { c.closable, c.onClose = true, fn } }

// Badge 徽标（antd Badge）。count <= 0 不显示；数字超过 99 显示 "99+"。
// 徽标显示在内容右侧（简单行内布局；叠在内容右上角的 Overlay 定位
// 列入后续增强，见 docs/06 下一步）：
//
//	ant.Badge(c, 5, func() { ant.Button(c, "消息") })
//	ant.BadgeDot(c, func() { ant.AvatarEl(c, "LY") })
func Badge(c *ui.Context, count int, fn func()) ui.Element {
	t := c.Theme()
	row := ui.Row(c).Gap(0).AlignItems(ui.Start)
	row.Children(func() {
		fn()
		if count > 0 {
			// 徽标叠在内容的右上角。
			text := "99+"
			if count <= 99 {
				text = strconv.Itoa(count)
			}
			badge := ui.Box(c).
				MinWidth(t.Space(4)).
				Height(t.Space(4)).
				Padding(0, t.Space(1)).
				Radius(t.Space(2)).
				Background(t.Danger).
				AlignItems(ui.Center)
			badge.Children(func() {
				ui.Text(c, text).TextColor(t.Background).FontSize(t.Rem(0.7))
			})
		}
	})
	return row
}

// BadgeDot 红点徽标（antd Badge dot）。
func BadgeDot(c *ui.Context, fn func()) ui.Element {
	t := c.Theme()
	row := ui.Row(c).Gap(0).AlignItems(ui.Start)
	row.Children(func() {
		fn()
		dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(t.Space(1)).Background(t.Danger)
		dot.Children(func() {})
	})
	return row
}

// AvatarEl 头像（antd Avatar）。name 提供首字符回退，image 可为 nil：
//
//	ant.AvatarEl(c, "李")
//	ant.AvatarEl(c, "Li", img)
func AvatarEl(c *ui.Context, name string, image ...*ui.Bitmap) ui.Element {
	var img *ui.Bitmap
	if len(image) > 0 {
		img = image[0]
	}
	t := c.Theme()
	el := ui.Avatar(c, name, img).Size(t.Space(8), t.Space(8))
	return el
}

// Card 卡片（antd Card）：头部（可选）、内容区、底部操作区（可选）。
//
//	ant.Card(c, "卡片标题", ant.CardExtra(func() { ... }), func() {
//	    ant.Paragraph(c, "卡片内容")
//	})
func Card(c *ui.Context, title string, body func(), opts ...CardOption) ui.Element {
	cfg := &cardConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := ui.Column(c).FillWidth().
		Radius(t.Radius).
		Background(t.Surface).
		Border(1, t.Border)
	el.Children(func() {
		if title != "" || cfg.extra != nil {
			head := ui.Row(c).FillWidth().
				Padding(t.Space(3), t.Space(4)).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center)
			head.Children(func() {
				if title != "" {
					ui.Text(c, title).FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
				}
				if cfg.extra != nil {
					cfg.extra(c)
				}
			})
			// 头部与内容之间的一条分隔线。
			line := ui.Box(c).FillWidth().Height(t.Space(1)).Border(1, t.Border)
			line.BorderWidth(0, 0, 1, 0)
			line.Children(func() {})
		}
		content := ui.Column(c).FillWidth().Padding(t.Space(3), t.Space(4)).Gap(t.Space(2))
		content.Children(body)
		if cfg.actions != nil {
			line := ui.Box(c).FillWidth().Height(t.Space(1)).Border(1, t.Border)
			line.BorderWidth(1, 0, 0, 0)
			line.Children(func() {})
			foot := ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(4)).Gap(t.Space(2))
			foot.Children(func() { cfg.actions(c) })
		}
	})
	return el
}

// CardOption 定制 Card。
type CardOption func(*cardConfig)

type cardConfig struct {
	extra   func(*ui.Context)
	actions func(*ui.Context)
}

// CardExtra 头部右侧附加内容。
func CardExtra(fn func(*ui.Context)) CardOption { return func(c *cardConfig) { c.extra = fn } }

// CardActions 底部操作区。
func CardActions(fn func(*ui.Context)) CardOption { return func(c *cardConfig) { c.actions = fn } }
