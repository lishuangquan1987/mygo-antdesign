package ant

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// DrawerOption 定制 Drawer。
type DrawerOption func(*drawerConfig)

type drawerConfig struct {
	width    float32 // 主题单位数
	closable bool
}

// DrawerWidth 抽屉宽度（主题单位数，默认 94.5 ≈ antd 378px；随密度缩放）。
func DrawerWidth(units float32) DrawerOption {
	return func(c *drawerConfig) { c.width = units }
}

// DrawerClosable 显示右上角关闭按钮（默认 true）。
func DrawerClosable(b bool) DrawerOption {
	return func(c *drawerConfig) { c.closable = b }
}

// Drawer 抽屉（antd Drawer）：从窗口右侧滑入的面板，带遮罩、标题栏与
// 可滚动内容区。open 由调用方持有，点遮罩或 Esc 关闭：
//
//	var open bool
//	ant.Drawer(c, &open, "设置", func() {
//	    ant.Paragraph(c, "抽屉内容")
//	}, ant.DrawerWidth(80))
func Drawer(c *ui.Context, open *bool, title string, body func(), opts ...DrawerOption) ui.Element {
	cfg := &drawerConfig{width: 94.5, closable: true}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c)
	// 打开首帧记录一次：从屏幕右缘滑入（antd Drawer 的 open 动效）。
	opened := ui.Local(host, "ant.drawer.opened", func() bool { return false })
	if !*open {
		*opened = false
	}
	host.Children(func() {
		ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
			// 遮罩：从背景与文字色混出中性灰，再压暗，亮/暗主题都成立。
			backdrop.Background(t.Background.Mix(t.Text, 0.9).Alpha(0.45))

			if !*opened {
				panel.Jump("ant.drawer.slide", t.Space(cfg.width))
				*opened = true
			}
			slide := panel.AnimateTo("ant.drawer.slide", 0, ui.Animation{Duration: 220 * time.Millisecond})

			panel.Width(t.Space(cfg.width)).
				FillHeight().
				Radius(0).
				Background(t.Surface).
				Border(1, t.Border).
				Translate(slide, 0)
			panel.BorderWidth(0, 0, 0, 1)
			panel.AttachTo(backdrop, ui.AnchorTopRight, ui.AnchorTopRight)

			panel.Children(func() {
				// 标题栏。
				head := ui.Row(c).FillWidth().
					Height(t.Space(14)).
					Padding(0, t.Space(6)).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Border(1, t.Border)
				head.BorderWidth(0, 0, 1, 0)
				head.Children(func() {
					ui.Text(c, title).TextColor(t.Text).FontSize(t.Rem(1.14)).FontWeight(600)
					if cfg.closable {
						closeBtn := ui.ButtonBase(c).Size(t.Space(6), t.Space(6)).Radius(t.Radius).Center()
						closeBtn.Children(func() {
							ui.Text(c, "✕").TextColor(t.TextMuted).FontSize(t.Rem(0.9))
						})
						if closeBtn.Clicked() {
							*open = false
						}
					}
				})
				// 可滚动内容区。
				ui.Scroll(c).FillWidth().Children(func() {
					inner := ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(4))
					inner.Children(body)
				})
			})
		})
	})
	return host
}
