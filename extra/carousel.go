package extra

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// CarouselOption 定制 Carousel。
type CarouselOption func(*carouselConfig)

type carouselConfig struct {
	height   float32       // 页面区高度；0 表示随内容
	arrows   bool          // 左右切换箭头（antd arrows）
	dots     bool          // 底部指示点（默认开启）
	autoplay time.Duration // 自动播放间隔；0 关闭
}

// CarouselHeight 页面区高度（DIP）。
func CarouselHeight(h float32) CarouselOption { return func(c *carouselConfig) { c.height = h } }

// CarouselArrows 显示左右切换箭头（antd arrows）。
func CarouselArrows(b bool) CarouselOption { return func(c *carouselConfig) { c.arrows = b } }

// CarouselDots 显示底部指示点（默认开启）。
func CarouselDots(b bool) CarouselOption { return func(c *carouselConfig) { c.dots = b } }

// CarouselAutoplay 自动播放间隔（antd autoplay，如 3*time.Second）；0 关闭。
// 系统开启「减少动态效果」时自动停止（Preferences.ReduceMotion）。
func CarouselAutoplay(d time.Duration) CarouselOption {
	return func(c *carouselConfig) { c.autoplay = d }
}

// Carousel 走马灯（antd Carousel）：一次显示一页，底部指示点切换，
// 可选左右箭头与自动播放。pages 为各页的构建函数（即时模式，每帧构建
// 当前页）：
//
//	ant.Carousel(c, []func(*ui.Context){
//	    func(c *ui.Context) { widgets.Title(c, "第一页", widgets.H3) },
//	    func(c *ui.Context) { widgets.Title(c, "第二页", widgets.H3) },
//	}, ant.CarouselArrows(true), ant.CarouselAutoplay(3*time.Second))
func Carousel(c *ui.Context, pages []func(*ui.Context), opts ...CarouselOption) ui.Element {
	cfg := &carouselConfig{dots: true}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).Gap(t.Space(2)).FillWidth()
	index := ui.Local(host, "ant.extra.carousel.index", func() int { return 0 })
	next := ui.Local(host, "ant.extra.carousel.next", func() time.Time { return time.Time{} })

	if len(pages) == 0 {
		host.Children(func() {})
		return host
	}
	if *index < 0 || *index >= len(pages) {
		*index = 0
	}

	// 自动播放：到点切下一页，并用 After 安排下一次重绘（悬停或减少动效时暂停）。
	if cfg.autoplay > 0 && len(pages) > 1 && !host.Hovered() && !c.Preferences().ReduceMotion {
		now := c.Now()
		if next.IsZero() {
			*next = now.Add(cfg.autoplay)
			c.After(cfg.autoplay)
		} else if !now.Before(*next) {
			*index = (*index + 1) % len(pages)
			*next = now.Add(cfg.autoplay)
			c.After(cfg.autoplay)
		}
	}

	host.Children(func() {
		// 页面区。
		page := ui.Box(c).FillWidth().Center()
		if cfg.height > 0 {
			page.Height(cfg.height)
		} else {
			page.MinHeight(t.Space(20))
		}
		page.Radius(t.Radius).Background(t.SurfaceHover)
		page.Children(func() {
			if *index >= 0 && *index < len(pages) {
				pages[*index](c)
			}
			// 左右切换箭头（antd arrows）：圆形半透明按钮，锚在页面两侧。
			if cfg.arrows && len(pages) > 1 {
				arrow := func(svg string, label string, delta int) {
					btn := ui.ButtonBase(c).Size(t.Space(8), t.Space(8)).
						Radius(t.Space(4)).
						Background(t.Surface.Alpha(0.9)).
						Border(1, t.Border).Center().Label(label)
					btn.Children(func() {
						icon.Icon(c, svg, t.Space(3), icon.WithColor(t.Text))
					})
					if btn.Clicked() {
						*index = (*index + delta + len(pages)) % len(pages)
					}
					if delta < 0 {
						btn.AttachTo(page, ui.AnchorLeft, ui.AnchorCenter)
					} else {
						btn.AttachTo(page, ui.AnchorRight, ui.AnchorCenter)
					}
				}
				arrow(icon.Left, "上一页", -1)
				arrow(icon.Right, "下一页", 1)
			}
		})
		// 指示点：外层 12×12 命中区，内含 8×8 视觉圆点。
		if cfg.dots && len(pages) > 1 {
			dots := ui.Row(c).Gap(t.Space(2)).Justify(ui.Center).FillWidth()
			dots.Children(func() {
				for i := range pages {
					i := i
					dot := ui.ButtonBase(c).Size(t.Space(3), t.Space(3)).Center().Label("第 " + itoa(i+1) + " 页")
					dot.Children(func() {
						inner := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(t.Space(1))
						inner.Children(func() {})
						if i == *index {
							inner.Background(t.Accent)
						} else {
							inner.Background(t.Border)
						}
					})
					if dot.Clicked() {
						*index = i
					}
				}
			})
		}
	})
	return host
}

// itoa 小整数转字符串（用于可访问标签）。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [8]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
